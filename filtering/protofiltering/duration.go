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
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

var durationMsgDesc = new(durationpb.Duration).ProtoReflect().Descriptor()

// TryParseDurationField tries to parse the provided value as a duration.
// It returns an error if the value is not a valid duration.
func (b *Interpreter) TryParseDurationField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if in.Value == nil {
		// A nil value is an internal error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "nil duration TryParseValueInput.Value"}, ErrInternal
		}
		return TryParseValueResult{}, ErrInternal
	}

	// If the value is a member expression, use the value of the member expression.
	if me, ok := in.Value.(*ast.MemberExpr); ok {
		in.Value = me.Value
		in.Args = me.Fields
	}

	// Duration can be a single or two fields separated by a period.
	if len(in.Args) > 1 {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a duration value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept string literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		if len(in.Args) == 0 && in.IsNullable && ft.Value == "null" {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}

		if len(in.Args) == 0 {
			// The duration probably don't have fractal part.
			// Try parsing it as an integer with unit.
			d, err := time.ParseDuration(ft.Value)
			if err != nil {
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}

			ve := expr.AcquireValueExpr()
			ve.Value = d
			return TryParseValueResult{Expr: ve}, nil
		}

		// There cannot be more than one argument for period separated duration.
		if len(in.Args) > 1 {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = in.Value.Position()
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
			}
			return res, ErrInvalidValue
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

		d, err := time.ParseDuration(sb.String())
		if err != nil {
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value+"."+fractal)}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}

		ve := expr.AcquireValueExpr()
		ve.Value = d
		return TryParseValueResult{Expr: ve, ArgsUsed: 1}, nil
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ae := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a duration value.
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
				ae.Free()
				return TryParseValueResult{}, ErrInternal
			}

			if !in.AllowIndirect {
				switch res.Expr.(type) {
				case *expr.FunctionCallExpr, *expr.FieldSelectorExpr:
					ae.Free()
					res.Expr.Free()
					if ctx.ErrHandler != nil {
						return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: fmt.Sprintf("field cannot accept function call or field selector expression as a value")}, ErrInvalidValue
					}
					return TryParseValueResult{}, ErrInternal
				}
			}
			ae.Elements = append(ae.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ae}, nil
	case *ast.StructExpr:
		// Check if the struct is a message or a map.
		if ft.IsMap() {
			// A map is not a valid duration value.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
			}
			return res, ErrInvalidValue
		}

		// If the struct name is a google.protobuf.Duration, parse it as a duration.
		var sb strings.Builder
		for i, nm := range ft.Name {
			if i > 0 {
				sb.WriteRune('.')
			}
			sb.WriteString(nm.String())
		}
		if sb.String() != "google.protobuf.Duration" {
			// This is not a duration.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
			}
			return res, ErrInvalidValue
		}

		// This is a duration.
		// Extract duration values from the struct.
		var seconds, nanos int64
		for _, field := range ft.Elements {
			if len(field.Name) != 1 {
				// This is not a valid durationpb. Invalid value.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
				}
				return res, ErrInvalidValue
			}

			var isNanos bool
			fieldName := field.Name[0].String()
			switch fieldName {
			case "seconds":
			case "nanos":
				isNanos = true
			default:
				// This is not a valid durationpb. Invalid value.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
				}
				return res, ErrInvalidValue
			}
			sec := durationMsgDesc.Fields().ByName(protoreflect.Name(fieldName))
			vi := TryParseValueInput{
				Field:         sec,
				AllowIndirect: in.AllowIndirect,
				IsNullable:    false,
				Value:         field.Value,
			}

			res, err := b.TryParseValue(ctx, vi)
			if err != nil {
				return res, err
			}

			if res.Expr == nil {
				// This is internal error, return an error.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: field.Position(), ErrMsg: "internal error: parsed expression is nil"}, ErrInternal
				}
				return TryParseValueResult{}, ErrInternal
			}

			ve, ok := res.Expr.(*expr.ValueExpr)
			if !ok {
				// Invalid seconds value.
				res.Expr.Free()
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
				}
				return res, ErrInvalidValue
			}

			var v int64
			v, ok = ve.Value.(int64)
			if !ok {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(in.Value, in.Args...))
				}
				return res, ErrInvalidValue
			}

			if isNanos {
				nanos = v
			} else {
				seconds = v
			}

			// Free the expression.
			res.Expr.Free()
		}

		d := durationpb.Duration{
			Seconds: seconds,
			Nanos:   int32(nanos),
		}

		ve := expr.AcquireValueExpr()
		ve.Value = d.AsDuration()
		return TryParseValueResult{Expr: ve}, nil
	default:
		// This is invalid AST node, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}
}