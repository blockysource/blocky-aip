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

package filtering

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

// TryParseSignedIntField tries to parse a signed int field.
// It can be a single signed int value or a repeated signed int value.
func (b *Interpreter) TryParseSignedIntField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if len(in.Args) > 0 {
		// A non-repeated signed int field cannot have nested fields.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	var tl *ast.TextLiteral
	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a signed int value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: "field cannot accept string literal as a value"}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		if !ft.Token.IsInteger() {
			// A text literal must be an int value.
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}
		tl = ft
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a signed int value.
			res, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         in.Field,
				AllowIndirect: in.AllowIndirect,
				IsOptional:    in.IsOptional,
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
						return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: "field cannot accept function call or field selector expression as a value"}, ErrInvalidValue
					}
					return TryParseValueResult{}, ErrInternal
				}
			}

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	case *ast.StructExpr:
		// A struct is not a valid signed int.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: "field cannot accept struct literal as a value"}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	default:
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}

	if in.IsOptional && tl.Token == token.NULL {
		ve := expr.AcquireValueExpr()
		ve.Value = nil
		return TryParseValueResult{Expr: ve}, nil
	}

	bs := 64
	switch in.Field.Kind() {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		bs = 32
	}
	v, err := strconv.ParseInt(tl.Value, 10, bs)
	if err != nil {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: tl.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), tl.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	ve := expr.AcquireValueExpr()
	ve.Value = v
	return TryParseValueResult{Expr: ve}, nil
}
