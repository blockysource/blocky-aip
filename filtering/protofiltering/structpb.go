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
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// TryParseStructPb tries to parse a well-known structpb.Value field.
func (b *Interpreter) TryParseStructPb(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
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

		// Check validity of the JSON string.
		if !json.Valid(bv) {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Pos
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid JSON string: '%s'", in.Field.Kind(), ft.Value)
			}
			return res, ErrInvalidValue
		}
		mp := make(map[string]any)
		if err := json.Unmarshal(bv, &mp); err != nil {
			// This is internal error, return an error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Pos
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid JSON string: '%s'", in.Field.Kind(), ft.Value)
			}
			return res, ErrInternal
		}
		ve := expr.AcquireValueExpr()
		ve.Value = mp
		return TryParseValueResult{Expr: ve}, nil
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
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = elem.Position()
					res.ErrMsg = "internal error: parsed expression is nil"
				}
				ve.Free()
				return res, ErrInternal
			}

			if !in.AllowIndirect {
				switch res.Expr.(type) {
				case *expr.FunctionCallExpr, *expr.FieldSelectorExpr:
					res.Expr.Free()
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = elem.Position()
						res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not valid: '%s'", in.Field.Kind(), joinedName(elem))
					}
					ve.Free()
					return res, ErrInvalidValue
				}
			}

			ve.Elements = append(ve.Elements, res.Expr)
		}
		return TryParseValueResult{Expr: ve}, nil
	case *ast.StructExpr:
		// Parse the value as a dynamic message.
		res, err := b.TryParseMessageStructField(ctx, in)
		if err != nil {
			return res, err
		}

		x := res.Expr
		ve, ok := x.(*expr.ValueExpr)
		if !ok {
			// This is internal error, a struct expr should be represented a value.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("internal error: struct expression is not represented as a value")
			}
			x.Free()
			return res, ErrInternal
		}

		dynMsg, ok := ve.Value.(*dynamicpb.Message)
		if !ok {
			// This is internal error, a struct expr should be represented a value.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("internal error: struct expression is not represented as a dynamicpb.Message but: %T", ve.Value)
			}
			x.Free()
			return res, ErrInternal
		}

		// Convert dynamic message to structpb.Struct by marshaling and unmarshaling.
		// NOTE: This is not the most efficient way, but it is the easiest way to do it.
		bt, err := proto.Marshal(dynMsg)
		if err != nil {
			// This is internal error, a struct expr should be represented a value.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("internal error: failed to marshal dynamicpb.Message: %v", err)
			}
			x.Free()
			return res, ErrInternal
		}

		st := structpb.Struct{}
		if err := proto.Unmarshal(bt, &st); err != nil {
			// This is internal error, a struct expr should be represented a value.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = ft.Position()
				res.ErrMsg = fmt.Sprintf("internal error: failed to unmarshal structpb.Struct: %v", err)
			}
			x.Free()
			return res, ErrInternal
		}

		ve.Value = st.AsMap()
		return TryParseValueResult{Expr: ve}, nil
	default:
		// This is invalid AST node, return an error.
		if ctx.ErrHandler != nil {
			return TryParseValueResult{ErrPos: in.Value.Position(), ErrMsg: "invalid AST node"}, ErrInvalidAST
		}
		return TryParseValueResult{}, ErrInvalidAST
	}
}
