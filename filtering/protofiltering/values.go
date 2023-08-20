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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		return b.TryParseWellKnownStructField(ctx, in)
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

func (b *Interpreter) TryParseMessageStructField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if in.Field.Message() == nil {
		// This is invalid AST node, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}

	// A struct field could either be null (if nullable) or a string literal of JSON format.
	if len(in.Args) > 0 {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a struct | nullable value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept string literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		if in.IsNullable && ft.Value == "null" {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}

		// Text literal cannot be a valid struct value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept text literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.StructExpr:
		// A struct can be parsed as a struct value, by setting a struct field expressions.
		desc := in.Field.Message()
		msg := dynamicpb.NewMessage(desc)

		for _, field := range ft.Elements {
			if len(field.Name) != 1 {
				var res TryParseValueResult
				// This is a map not a struct.
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(field.Name[0]))
				}
				return res, ErrInvalidValue
			}
			df := desc.Fields().ByName(protoreflect.Name(field.Name[0].UnquotedString()))
			if df == nil {
				// Field is not found within the message descriptor.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("field is not found within the message descriptor: %s", field.Name[0].UnquotedString())}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}

			// Check if the field is already set.
			if msg.Has(df) {
				// The field was duplicated.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("field %s is duplicated", field.Name[0].UnquotedString())}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}

			// Try parsing the field value.
			v, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         df,
				AllowIndirect: in.AllowIndirect,
				IsNullable:    IsFieldNullable(df),
				Value:         field.Value,
			})
			if err != nil {
				return v, err
			}

			if v.Expr == nil {
				// This is internal error, return an error.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: parsed expression is nil"}, ErrInternal
				}
				return TryParseValueResult{}, ErrInternal
			}

			switch vt := v.Expr.(type) {
			case *expr.ValueExpr:
				if df.Cardinality() == protoreflect.Repeated {
					// This is a repeated field, but we have a single value.
					// This is a syntax error.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = field.Position()
						res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), joinedName(field.Name[0]))
					}
					vt.Free()
					return res, ErrInvalidValue
				}

				var (
					pv  protoreflect.Value
					res TryParseValueResult
				)

				pv, res, err = b.exprValueToProto(ctx, df, vt, field)
				if err != nil {
					vt.Free()
					return res, err
				}

				if pv.IsValid() {
					// Set the value and free the field expression.
					msg.Set(df, pv)
				}
				vt.Free()
			case *expr.ArrayExpr:
				if df.Cardinality() != protoreflect.Repeated {
					// This is a repeated field, but we have a single value.
					// This is a syntax error.
					vt.Free()
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = field.Position()
						res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)
					}
					return res, ErrInvalidValue
				}

				ls := msg.Mutable(df).List()

				for _, elem := range vt.Elements {
					switch et := elem.(type) {
					case *expr.ValueExpr:
						var (
							pv  protoreflect.Value
							res TryParseValueResult
						)
						pv, res, err = b.exprValueToProto(ctx, df, et, field)
						if err != nil {
							vt.Free()
							return res, err
						}

						if pv.IsValid() {
							ls.Append(pv)
						}
					default:
						// This is internal error, return an error.
						vt.Free()
						if ctx.ErrHandler != nil {
							return TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: unknown value type: %T", et)}, ErrInternal
						}
						return TryParseValueResult{}, ErrInternal
					}
				}
				vt.Free()
			default:
				// This is internal error, return an error.
				vt.Free()
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: unknown expression type: %T to parse a Message", vt)}, ErrInternal
				}
				return TryParseValueResult{}, ErrInternal
			}
		}

		// Create a value expression with a dynamic message value.
		ve := expr.AcquireValueExpr()
		ve.Value = msg
		return TryParseValueResult{Expr: ve}, nil
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a message value.
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

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	default:
		// This is internal error, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("internal error: unknown value type: %T to parse a Message", ft)}, ErrInternal
		}
		return TryParseValueResult{}, ErrInternal
	}
}

func (b *Interpreter) exprValueToProto(ctx *ParseContext, df protoreflect.FieldDescriptor, vt *expr.ValueExpr, field *ast.StructFieldExpr) (protoreflect.Value, TryParseValueResult, error) {
	var pv protoreflect.Value
	switch et := vt.Value.(type) {
	case time.Time:
		if df.Kind() != protoreflect.MessageKind {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a message kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
		if df.Message().FullName() != "google.protobuf.Timestamp" {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a timestamp message"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}

		pv = protoreflect.ValueOfMessage(timestamppb.New(et).ProtoReflect())
	case time.Duration:
		if df.Kind() != protoreflect.MessageKind {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a message kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
		if df.Message().FullName() != "google.protobuf.Duration" {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a duration message"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
		pv = protoreflect.ValueOfMessage(durationpb.New(et).ProtoReflect())
	case int64:
		switch df.Kind() {
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			pv = protoreflect.ValueOfInt32(int32(et))
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			pv = protoreflect.ValueOfInt64(et)
		default:
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not an integer kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case uint64:
		switch df.Kind() {
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			pv = protoreflect.ValueOfUint32(uint32(et))
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			pv = protoreflect.ValueOfUint64(et)
		default:
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not an unsigned integer kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case float64:
		switch df.Kind() {
		case protoreflect.DoubleKind:
			pv = protoreflect.ValueOfFloat64(et)
		case protoreflect.FloatKind:
			pv = protoreflect.ValueOfFloat32(float32(et))
		default:
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a float kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case bool:
		if df.Kind() != protoreflect.BoolKind {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a bool kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
		pv = protoreflect.ValueOfBool(et)
	case string:
		switch df.Kind() {
		case protoreflect.StringKind:
			pv = protoreflect.ValueOfString(et)
		case protoreflect.BytesKind:
			pv = protoreflect.ValueOfBytes([]byte(et))
		default:
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a string kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case []byte:
		switch df.Kind() {
		case protoreflect.StringKind:
			pv = protoreflect.ValueOfString(string(et))
		case protoreflect.BytesKind:
			pv = protoreflect.ValueOfBytes(et)
		default:
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a bytes kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case protoreflect.EnumNumber:
		if df.Kind() != protoreflect.EnumKind {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not an enum kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
		pv = protoreflect.ValueOfEnum(et)
		if df.Enum().Values().ByNumber(et) == nil {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: enum value %d is not found", et)}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}
	case proto.Message:
		// This works for both protoreflect.Message and structpb.Value.
		if df.Kind() != protoreflect.MessageKind {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: field is not a message kind"}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}

		if df.Message().FullName() != et.ProtoReflect().Descriptor().FullName() {
			// This is internal error, return an error.
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: message type %q is not expected", et.ProtoReflect().Descriptor().FullName())}, ErrInternal
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
		}

		pv = protoreflect.ValueOfMessage(et.ProtoReflect())
	case nil:
		if df.Cardinality() == protoreflect.Repeated {
			// This is a repeated field, but we have a single value.
			// This is a syntax error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = field.Position()
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)
			}
			return protoreflect.Value{}, res, ErrInvalidValue
		}
	default:
		// This is internal error, return an error.
		if ctx.ErrHandler != nil {
			return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: unknown value type: %T", et)}, ErrInternal
		}
		return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
	}
	return pv, TryParseValueResult{}, nil
}

// TryParseWellKnownStructField tries to parse a well-known structpb.Value field.
func (b *Interpreter) TryParseWellKnownStructField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	// A struct field could either be null (if nullable) or a string literal of JSON format.
	if len(in.Args) > 0 {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal can be a nullable value.
		// Check if the value is a valid JSON string.
		bv := json.RawMessage(ft.Value)

		if json.Valid(bv) {
			ve := expr.AcquireValueExpr()
			ve.Value = bv
			return TryParseValueResult{Expr: ve}, nil
		}

		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		if in.IsNullable && ft.Value == "null" {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}

		// Text literal cannot be a struct value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept text literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a struct value.
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
					res.Expr.Free()
					if ctx.ErrHandler != nil {
						return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(elem))}, ErrInvalidValue
					}
					return TryParseValueResult{}, ErrInvalidValue
				}
			}

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	case *ast.StructExpr:
		return b.TryParseMessageStructField(ctx, in)
	default:
		// This is invalid AST node, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}
}

// TryParseBooleanField tries to parse a boolean field.
// It can be a single boolean value or a repeated boolean value.
func (b *Interpreter) TryParseBooleanField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a bool value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept string literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		// Only the text literal can be a bool value.
		switch {
		case ft.Value == "true":
			ve := expr.AcquireValueExpr()
			ve.Value = true
			return TryParseValueResult{Expr: ve}, nil
		case ft.Value == "false":
			ve := expr.AcquireValueExpr()
			ve.Value = false
			return TryParseValueResult{Expr: ve}, nil
		case in.IsNullable && ft.Value == "null":
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}
		// Invalid boolean value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of bool type, but provided value is not a valid bool value: '%s'", ft.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.KeywordExpr:
		// Keyword expression cannot be a bool value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept keyword expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a bool value.
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
	case *ast.StructExpr:
		// A struct is not a valid bool.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: fmt.Sprintf("field cannot accept struct expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}
	_, ok := in.Value.(ast.FieldExpr)
	if !ok {
		// Invalid AST syntax, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "internal error: invalid AST syntax"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}

	// A FieldSelectorExpr can either be a value or keyword. ValueExpr is either string literal or text literal.
	// This means that the FieldSelectorExpr is a keyword expression.
	if ctx.ErrHandler != nil {
		return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field cannot accept keyword expression as a value")}, ErrInvalidValue
	}
	return TryParseValueResult{}, ErrInvalidValue
}

// TryParseFloatField tries to parse a float field.
// It can be a single float value or a repeated float value.
func (b *Interpreter) TryParseFloatField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal can be a float value.
		ve := expr.AcquireValueExpr()
		ve.Value = ft.Value
		return TryParseValueResult{Expr: ve}, nil
	case *ast.KeywordExpr:
		// Keyword expression cannot be a float value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept keyword expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		// Only the text literal can be a float value.
		if len(in.Args) == 0 {
			if in.IsNullable && ft.Value == "null" {
				ve := expr.AcquireValueExpr()
				ve.Value = nil
				return TryParseValueResult{Expr: ve}, nil
			}
			// This is a non fractial numeric value.
			// Try parsing it as an integer.
			v, err := strconv.ParseInt(ft.Value, 10, 64)
			if err != nil {
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}
			ve := expr.AcquireValueExpr()
			ve.Value = float64(v)
			return TryParseValueResult{Expr: ve}, nil
		}

		// There cannot be more than one argument for period separated float.
		if len(in.Args) > 1 {
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}

		var fractal string
		// This is a fractal numeric value.
		// Try parsing it as a float.
		switch at := in.Args[0].(type) {
		case *ast.TextLiteral:
			fractal = at.Value
		default:
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value+"."+at.String())}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}

		var sb strings.Builder
		sb.WriteString(ft.Value)
		sb.WriteRune('.')
		sb.WriteString(fractal)

		bs := 64
		if in.Field.Kind() == protoreflect.FloatKind {
			bs = 32
		}
		v, err := strconv.ParseFloat(sb.String(), bs)
		if err != nil {
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value+"."+fractal)}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}
		ve := expr.AcquireValueExpr()
		ve.Value = v
		return TryParseValueResult{Expr: ve, ArgsUsed: 1}, nil
	case *ast.ArrayExpr:
		// Parse each element of the array.
		// If any element is not a valid float value, return an error.
		ve := expr.AcquireArrayExpr()
		for _, e := range ft.Elements {
			te, err := b.TryParseFloatField(ctx, TryParseValueInput{
				Field:         in.Field,
				AllowIndirect: in.AllowIndirect,
				IsNullable:    in.IsNullable,
				Value:         e,
			})
			if err != nil {
				return TryParseValueResult{}, err
			}
			ve.Elements = append(ve.Elements, te.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	case *ast.StructExpr:
		// A struct value cannot be a float value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: fmt.Sprintf("field cannot accept struct expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.FunctionCall:
		// Call the function.
		res, err := b.TryParseFunctionCall(ctx, in)
		if err != nil {
			return TryParseValueResult{}, err
		}

		// If the input does not allow indirect value, and result is a FunctionCall or FieldSelectorExpr,
		// then return an error.
		if !in.AllowIndirect {
			switch res.Expr.(type) {
			case *expr.FunctionCallExpr, *expr.FieldSelectorExpr:
				res.Expr.Free()
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = in.Value.Position()
					res.ErrMsg = fmt.Sprintf("field does not allow indirect value")
				}
				return res, ErrInvalidValue
			}
		}
		return res, nil
	}
	// Invalid AST syntax, return an error.
	if ctx.ErrHandler != nil {
		return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "internal error: invalid AST syntax"}, ErrInvalidAST
	}
	return TryParseValueResult{}, ErrInvalidAST
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

// TryParseEnumField tries to parse an enum field.
// It can be a single enum value or a repeated enum value.
func (b *Interpreter) TryParseEnumField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if len(in.Args) > 0 {
		// A non-repeated enum field cannot have nested fields.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	if in.Field.Enum() == nil {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is not an enum field")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	var sl *ast.StringLiteral
	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		sl = ft
	case *ast.TextLiteral:
		if in.IsNullable && ft.Value == "null" {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'. String literal required", in.Field.Enum().FullName(), ft.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.KeywordExpr:
		// Keyword expression cannot be a enum value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept keyword expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a enum value.
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
	case *ast.StructExpr:
		// A struct is not a valid enum.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: fmt.Sprintf("field cannot accept struct expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	default:
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}

	enumValue := in.Field.Enum().Values().ByName(protoreflect.Name(sl.Value))
	if enumValue == nil {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: sl.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Enum().FullName(), sl.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	ve := expr.AcquireValueExpr()
	ve.Value = enumValue.Number()
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

func GetFieldComplexity(fd FieldDescriptor) int64 {
	switch fdt := fd.(type) {
	case *FunctionCallArgumentDeclaration:
		return 1
	case *FunctionCallReturningDeclaration:
		return 1
	case protoreflect.FieldDescriptor:
		c, ok := proto.GetExtension(fdt.Options(), blockyannotations.E_Complexity).(int64)
		if ok {
			return c
		}
		return 1
	default:
		return 1
	}
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
