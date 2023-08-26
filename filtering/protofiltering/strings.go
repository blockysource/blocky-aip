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

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

// TryParseStringField tries to parse a string field.
// It can be a single string value or a repeated string value.
func (b *Interpreter) TryParseStringField(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	if len(in.Args) > 0 {
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: fmt.Sprintf("field is of string type, but provided value is not a valid string value: '%s'", joinedName(in.Value, in.Args...))}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	}

	switch ft := in.Value.(type) {
	case *ast.StringLiteral:
		// String literal can be a string value.
		// Check if the string literal have prefix or suffix wildcard.
		var (
			hasPrefixWildcard, hasSuffixWildcard bool
			strValue                             string
		)
		strValue = ft.Value

		if len(strValue) > 0 {
			hasPrefixWildcard = strValue[0] == '*'
			if hasPrefixWildcard {
				strValue = strValue[1:]
			}
		}

		if len(strValue) > 0 {
			hasSuffixWildcard = strValue[len(strValue)-1] == '*'
			if hasSuffixWildcard {
				strValue = strValue[:len(strValue)-1]
			}
		}

		if hasPrefixWildcard || hasSuffixWildcard {
			if !in.AllowIndirect {
				// Wildcard is not allowed for non-indirect values.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field is of string type, but provided value is not a valid string value: '%s'", ft.Value)}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}
			if len(strValue) == 0 {
				// String containing only wildcard is not allowed.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("cannot use a wildcard only string value: '%s'", ft.Value)}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}

			ve := expr.AcquireStringSearchExpr()
			ve.Value = strValue
			ve.PrefixWildcard = hasPrefixWildcard
			ve.SuffixWildcard = hasSuffixWildcard
			ve.SearchComplexity = in.Complexity
			return TryParseValueResult{Expr: ve, IsIndirect: true}, nil
		}

		ve := expr.AcquireValueExpr()
		ve.Value = ft.Value
		return TryParseValueResult{Expr: ve}, nil
	case *ast.TextLiteral:
		if in.IsNullable && ft.Token == token.NULL {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}

		// Text literal cannot be a string value.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Pos, ErrMsg: fmt.Sprintf("field cannot accept text literal as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range ft.Elements {
			// Try parsing each element as a string value.
			res, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         in.Field,
				AllowIndirect: in.AllowIndirect,
				IsNullable:    in.IsNullable,
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
						return TryParseValueResult{ErrPos: elem.Position(), ErrMsg: fmt.Sprintf("field cannot accept function call or field selector expression as a value")}, ErrInvalidValue
					}
					return TryParseValueResult{}, ErrInternal
				}
			}

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	case *ast.StructExpr:
		// A struct is not a valid string.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: fmt.Sprintf("field cannot accept struct expression as a value")}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	default:
		// This is invalid AST node, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}
}
