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
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// TryParseFloatField tries to parse a float field.
// It can be a single float value or a repeated float value.
func (b *Interpreter) TryParseFloatField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if in.Value == nil {
		var res TryParseValueResult
		// Internal error, no value in the input.
		if ctx.ErrHandler != nil {
			res.ErrMsg = fmt.Sprintf("internal error: no input value provided")
		}
		return res, ErrInternal
	}

	if me, ok := in.Value.(*ast.MemberExpr); ok {
		in.Value = me.Value
		in.Args = me.Fields
	}

	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a float value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept string literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
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
				ve.Free()
				return te, err
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
