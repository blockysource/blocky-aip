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

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

// TryParseBooleanField tries to parse a boolean field.
// It can be a single boolean value or a repeated boolean value.
func (b *Interpreter) TryParseBooleanField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal cannot be a bool value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: "field cannot accept string literal as a value"}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.TextLiteral:
		if !ft.Token.IsBoolean() {
			if ctx.ErrHandler != nil {
				return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), ft.Value)}, ErrInvalidValue
			}
			return TryParseValueResult{}, ErrInvalidValue
		}

		// Only the text literal can be a bool value.
		switch {
		case ft.Token == token.TRUE:
			ve := expr.AcquireValueExpr()
			ve.Value = true
			return TryParseValueResult{Expr: ve}, nil
		case ft.Token == token.FALSE:
			ve := expr.AcquireValueExpr()
			ve.Value = false
			return TryParseValueResult{Expr: ve}, nil
		case in.IsOptional && ft.Token == token.NULL:
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}
		// Invalid boolean value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of bool type, but provided value is not a valid bool value: '%s'", ft.Value)}, ErrInvalidValue
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
				IsOptional:    in.IsOptional,
				Value:         elem,
				Complexity:    in.Complexity,
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
		// A struct is not a valid bool.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: "field cannot accept struct expression as a value"}, ErrInvalidValue
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
		return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "field cannot accept keyword expression as a value"}, ErrInvalidValue
	}
	return TryParseValueResult{}, ErrInvalidValue
}
