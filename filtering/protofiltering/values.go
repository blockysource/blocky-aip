// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protofiltering

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/filtering/token"
	blockyannotations "github.com/blockysource/go-genproto/blocky/api/annotations"
)

// TryParseValueInput is an input for the TryParseValue function.
// It is used to parse a value expression either directly or indirectly.
type TryParseValueInput struct {
	// Field is a required field descriptor of the value.
	Field FieldDescriptor

	// Value is a required value expression.
	Value ast.AnyExpr

	// AllowIndirect is a flag that indicates whether the value can be indirect.
	// An indirect call is defined either by the function call or the field selector expression.
	AllowIndirect bool

	// IsNullable is a flag that indicates whether the value can be null.
	IsNullable bool

	// Args are the optional arguments of the value.
	// Used mostly by the member expression fields.
	Args []ast.FieldExpr
}

// TryParseValueResult is a result of the TryParseValue function.
// It either contains an expression or an error with the error position and message.
// ArgsUsed is the number of arguments used by the value from the Args input.
type TryParseValueResult struct {
	// Expr is the parsed expression.
	Expr expr.FilterExpr

	// ErrPos is the detailed error position.
	ErrPos token.Position

	// ErrMsg is the detailed error message.
	ErrMsg string

	// ArgsUsed is the number of arguments used by the value from the Args input.
	ArgsUsed int

	// IsIndirect is a flag that indicates whether the value is indirect.
	// An indirect value means it depends on the field selector.
	IsIndirect bool
}

// IsFieldNullable checks if the input field is nullable.
func IsFieldNullable(field protoreflect.FieldDescriptor) bool {
	// At first try blockaypi.E_Nullable extension, if not found, then try google api.OPTIONAL extension.
	// If not found, then return false.
	queryOpts, ok := proto.GetExtension(field.Options(), blockyannotations.E_QueryOpt).([]blockyannotations.FieldQueryOption)
	if ok {
		for _, qo := range queryOpts {
			if qo == blockyannotations.FieldQueryOption_NULLABLE {
				return true
			}
		}
	}

	fb, ok := proto.GetExtension(field.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
	if !ok {
		return false
	}

	if field.Kind() == protoreflect.MessageKind {
		for _, b := range fb {
			if b == annotations.FieldBehavior_REQUIRED {
				return false
			}
		}
		return true
	}
	for _, b := range fb {
		switch b {
		case annotations.FieldBehavior_REQUIRED, annotations.FieldBehavior_IMMUTABLE:
			return false
		case annotations.FieldBehavior_OPTIONAL:
			return true
		}
	}
	return false
}

// GetFieldComplexity extracts the complexity of a field out of the given field descriptor.
func GetFieldComplexity(fd FieldDescriptor) int64 {
	switch fdt := fd.(type) {
	case *FunctionCallArgumentDeclaration:
		return 1
	case *FunctionCallReturningDeclaration:
		return 1
	case protoreflect.FieldDescriptor:
		return getFieldComplexity(fdt)
	default:
		return 1
	}
}

func getFieldComplexity(fdt protoreflect.FieldDescriptor) int64 {
	c, ok := proto.GetExtension(fdt.Options(), blockyannotations.E_Complexity).(int64)
	if !ok || c == 0 {
		return 1
	}
	return c
}

// TryParseValue tries to parse a value expression.
func (b *Interpreter) TryParseValue(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if in.Field == nil {
		// Internal error - no field is defined.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "internal error: no field is defined"}, ErrInternal
		}
		return TryParseValueResult{}, ErrInternal
	}
	if in.Value == nil {
		// Internal error - no value is defined.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "internal error: no input value is defined"}, ErrInternal
		}
		return TryParseValueResult{}, ErrInternal
	}
	me, ok := in.Value.(*ast.MemberExpr)
	if ok {
		in.Value = me.Value
		in.Args = me.Fields
	}
	switch in.Field.Kind() {
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		return b.TryParseFloatField(ctx, in)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return b.TryParseSignedIntField(ctx, in)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return b.TryParseSignedIntField(ctx, in)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return b.TryParseUnsignedIntField(ctx, in)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return b.TryParseUnsignedIntField(ctx, in)
	case protoreflect.BoolKind:
		return b.TryParseBooleanField(ctx, in)
	case protoreflect.StringKind:
		return b.TryParseStringField(ctx, in)
	case protoreflect.BytesKind:
		return b.TryParseBytesField(ctx, in)
	case protoreflect.EnumKind:
		return b.TryParseEnumField(ctx, in)
	case protoreflect.MessageKind:
		return b.TryParseMessageField(ctx, in)
	case protoreflect.GroupKind:
		// Group is deprecated and thus not supported.
		if ctx.ErrHandler != nil {
			ctx.ErrHandler(in.Value.Position(), fmt.Sprintf("field is of deprecated %q type", in.Field.Kind()))
		}
		return TryParseValueResult{}, ErrInvalidValue
	default:
		// No other possible field kind, return an error.
		if ctx.ErrHandler != nil {
			ctx.ErrHandler(in.Value.Position(), fmt.Sprintf("field is of unsupported %q type", in.Field.Kind()))
		}
		return TryParseValueResult{}, ErrInvalidValue
	}
}

// TryParseMessageField tries to parse a message field.
func (b *Interpreter) TryParseMessageField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if in.Field.Message() == nil {
		// Internal error - no field message is defined even though it is a MessageKind.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "internal error: field is of MessageKind, but no message is defined"}, ErrInternal
		}
		return TryParseValueResult{}, ErrInternal
	}
	if in.Field.IsMap() {
		return b.TryParseMapField(ctx, in)
	}

	switch in.Field.Message().FullName() {
	case "google.protobuf.Timestamp":
		return b.TryParseTimestampField(ctx, in)
	case "google.protobuf.Duration":
		return b.TryParseDurationField(ctx, in)
	case "google.protobuf.Struct":
		return b.TryParseStructPb(ctx, in)
	default:
		return b.TryParseMessageStructField(ctx, in)
	}
}

// TryParseMapField tries to parse a map field.
func (b *Interpreter) TryParseMapField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if ctx.Message == nil {
		return TryParseValueResult{}, ErrInternal
	}

	if in.Value == nil {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "TryParseMapField: value is nil"
		}
		return res, ErrInternal
	}

	if in.Field == nil {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "TryParseMapField: field is nil"
			res.ErrPos = in.Value.Position()
		}
		return res, ErrInternal
	}

	if !in.Field.IsMap() {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "TryParseMapField: field is not a map"
			res.ErrPos = in.Value.Position()
		}
		return res, ErrInternal
	}

	switch vt := in.Value.(type) {
	case *ast.StructExpr:
		// This is a proper format for the map value.
		// The value is a struct, so we need to parse it as a struct.
		mve := expr.AcquireMapValueExpr()
		kd := in.Field.MapKey()
		vd := in.Field.MapValue()
		for _, elem := range vt.Elements {
			ki := TryParseValueInput{
				Field:         kd,
				IsNullable:    false,
				AllowIndirect: false,
			}
			if len(elem.Name) == 0 {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrMsg = "field is a map, but has invalid key type"
					res.ErrPos = elem.Position()
				}
				return res, ErrInvalidValue
			}

			ki.Value = elem.Name[0]
			if len(elem.Name) > 1 {
				for _, arg := range elem.Name[1:] {
					ki.Args = append(ki.Args, arg)
				}
			}
			kv, err := b.TryParseValue(ctx, ki)
			if err != nil {
				mve.Free()
				return kv, err
			}

			kve, ok := kv.Expr.(*expr.ValueExpr)
			if !ok {
				mve.Free()
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrMsg = "field is a map, but has invalid key value type"
					res.ErrPos = elem.Position()
				}
				return res, ErrInvalidValue
			}

			vv, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         vd,
				Value:         elem.Value,
				IsNullable:    false,
				AllowIndirect: false,
			})
			if err != nil {
				return vv, err
			}

			if _, ok = vv.Expr.(*expr.ArrayExpr); ok {
				mve.Free()
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrMsg = "field is a map, but has invalid value type"
					res.ErrPos = elem.Position()
				}
				return res, ErrInvalidValue
			}

			mve.Values = append(mve.Values, expr.MapValueExprEntry{
				Key:   kve,
				Value: vv.Expr,
			})
		}
		return TryParseValueResult{Expr: mve}, nil
	case *ast.TextLiteral:
		if in.IsNullable && vt.Value == "null" {
			ve := expr.AcquireValueExpr()
			ve.Value = nil

			return TryParseValueResult{Expr: ve}, nil
		}

		// Otherwise it is not a proper format for the map value.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "field is a map, but has invalid value type"
			res.ErrPos = in.Value.Position()
		}
		return res, ErrInvalidValue
	case *ast.StringLiteral, *ast.ArrayExpr:
		// Neither of these types are supported as map values.
		// From protobuf perspective a map cannot be repeated.
		// From filtering perspective no comparison operator can work on an array of maps (even IN).
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "field is a map, but has invalid value type"
			res.ErrPos = in.Value.Position()
		}
		return res, ErrInvalidValue
	default:
		// This is not a proper format for the map value.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = "field is a map, but has invalid value type"
			res.ErrPos = in.Value.Position()
		}
		return res, ErrInvalidValue
	}
}

// TryParseBytesField tries to parse a bytes field.
// It can be a single bytes value or a repeated bytes value.
func (b *Interpreter) TryParseBytesField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	// Check if no more Fields are present in the input *x.MemberExpr.
	// If there are, then return an error.
	if len(in.Args) > 0 {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	var value string
	switch vt := in.Value.(type) {
	case *ast.TextLiteral:
		value = vt.Value
	case *ast.StringLiteral:
		value = vt.Value
	case *ast.KeywordExpr:
		// KeywordExpr is not supported for bytes field.
		if ctx.ErrHandler != nil {
			ctx.ErrHandler(vt.Position(), "keyword expression is not supported for bytes field")
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range vt.Elements {
			// Try parsing each element as a bytes value.
			res, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         in.Field,
				AllowIndirect: in.AllowIndirect,
				IsNullable:    in.IsNullable,
				Value:         elem,
			})
			if err != nil {
				return res, err
			}

			if res.Expr == nil {
				// This is internal error, return an error.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: "internal error: parsed expression is nil"}, ErrInternal
				}
				return TryParseValueResult{}, ErrInternal
			}

			if !in.AllowIndirect {
				switch res.Expr.(type) {
				case *expr.FunctionCallExpr, *expr.FieldSelectorExpr:
					ve.Free()
					res.Expr.Free()
					if ctx.ErrHandler != nil {
						return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: fmt.Sprintf("field cannot accept function call or field selector expression as a value")}, ErrInvalidValue
					}
					return TryParseValueResult{}, ErrInternal
				}
			}

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	}

	if in.IsNullable && value == "null" {
		ve := expr.AcquireValueExpr()
		ve.Value = nil
		return TryParseValueResult{Expr: ve}, nil
	}
	dec, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	ve := expr.AcquireValueExpr()
	ve.Value = dec

	return TryParseValueResult{Expr: ve}, nil
}

func joinedName(v ast.AnyExpr, args ...ast.FieldExpr) string {
	var sb strings.Builder
	v.WriteStringTo(&sb, false)
	for _, arg := range args {
		sb.WriteRune('.')
		arg.WriteStringTo(&sb, false)
	}
	return sb.String()
}

func isKindComparable(k1, k2 protoreflect.Kind) bool {
	if k1 == k2 {
		return true
	}

	switch k1 {
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		switch k2 {
		case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			return true
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
			return true
		default:
			return false
		}
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		switch k2 {
		case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
			return true
		case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind, protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
			return true
		default:
			return false
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		switch k2 {
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			return true
		default:
			return false
		}
	default:
		return false
	}
}
