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
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

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
		if in.IsOptional && vt.Token == token.NULL {
			ve := expr.AcquireValueExpr()
			ve.Value = nil
			return TryParseValueResult{Expr: ve}, nil
		}

		if vt.Token == token.HEX {
			// Decode hex string to bytes.
			// At first trim the prefix "0x" from the hex string.
			bt, err := hex.DecodeString(vt.Value[2:])
			if err != nil {
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: vt.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), vt.Value)}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}

			ve := expr.AcquireValueExpr()
			ve.Value = bt
			return TryParseValueResult{Expr: ve}, nil
		}

		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: vt.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), vt.Value)}, ErrInvalidValue
		}
		return TryParseValueResult{}, ErrInvalidValue
	case *ast.StringLiteral:
		value = vt.Value
	case *ast.ArrayExpr:
		// An array can be parsed as a repeated field value.
		ve := expr.AcquireArrayExpr()
		for _, elem := range vt.Elements {
			// Try parsing each element as a bytes value.
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
	}

	if in.IsOptional && value == "null" {
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
