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
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

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
		if in.IsOptional && ft.Token == token.NULL {
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

		if len(ft.Name) > 0 {
			// Check matching name of the message.
			if desc.FullName() != ft.FullName() {
				// The message name doesn't match.
				if ctx.ErrHandler != nil {
					return TryParseValueResult{ErrPos: ft.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid %q value: '%s'", in.Field.Kind(), in.Field.Kind(), joinedName(ft.Name[0]))}, ErrInvalidValue
				}
				return TryParseValueResult{}, ErrInvalidValue
			}
		}

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

			fi := b.msgInfo.GetFieldInfo(df)

			// Try parsing the field value.
			v, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         df,
				AllowIndirect: in.AllowIndirect,
				IsOptional:    fi.Nullable,
				Value:         field.Value,
				Complexity:    fi.Complexity,
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

			pv, res, err := b.exprFilterToProto(ctx, msg, df, v.Expr, field)
			if err != nil {
				v.Expr.Free()
				return res, err
			}

			if !pv.IsValid() {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)
				}
				v.Expr.Free()
				return res, ErrInternal
			}

			// Set the value and free the field expression.
			msg.Set(df, pv)
			v.Expr.Free()
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

func (b *Interpreter) exprFilterToProto(ctx *ParseContext, msg *dynamicpb.Message, df protoreflect.FieldDescriptor, ft expr.FilterExpr, field *ast.StructFieldExpr) (protoreflect.Value, TryParseValueResult, error) {
	switch vt := ft.(type) {
	case *expr.ValueExpr:
		pv, res, err := b.exprValueToProto(ctx, df, vt, field)
		if err != nil {
			return protoreflect.Value{}, res, err
		}
		return pv, TryParseValueResult{}, nil
	case *expr.ArrayExpr:
		if df.Cardinality() != protoreflect.Repeated {
			// This is a repeated field, but we have a single value.
			// This is a syntax error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = field.Position()
				res.ErrMsg = fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)
			}
			return protoreflect.Value{}, res, ErrInvalidValue
		}

		ls := msg.Mutable(df).List()

		for _, elem := range vt.Elements {
			switch et := elem.(type) {
			case *expr.ValueExpr:
				ev, res, err := b.exprValueToProto(ctx, df, et, field)
				if err != nil {
					return protoreflect.Value{}, res, err
				}

				ls.Append(ev)
			default:
				// Array values cannot be different then value expressions.
				// The definition of prutobufs doesn't allow to use arrays of arrays,
				// or arrays of maps.
				// This is internal error, return an error.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("internal error: unknown value type: %T to parse a Message", et)
				}
				return protoreflect.Value{}, res, ErrInternal
			}
		}
		return protoreflect.ValueOfList(ls), TryParseValueResult{}, nil
	case *expr.MapValueExpr:
		if df.Cardinality() != protoreflect.Repeated && df.Kind() != protoreflect.MessageKind {
			// This is a repeated field, but we have a single value.
			// This is a syntax error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = field.Position()
			}
		}

		mv := msg.NewField(df).Map()
		for _, elem := range vt.Values {
			kv, res, err := b.exprValueToProto(ctx, df.MapKey(), elem.Key, field)
			if err != nil {
				return protoreflect.Value{}, res, err
			}

			// A value of a map entry cannot be different then value expression,
			// as protobuf definition doesn't allow to use maps of maps or maps of arrays.
			eve, ok := elem.Value.(*expr.ValueExpr)
			if !ok {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = field.Position()
					res.ErrMsg = fmt.Sprintf("internal error: unknown value type: %T to parse a Message", elem.Value)
				}
				return protoreflect.Value{}, res, ErrInternal
			}
			vv, res, err := b.exprValueToProto(ctx, df.MapValue(), eve, field)
			if err != nil {
				return protoreflect.Value{}, res, err
			}

			mv.Set(protoreflect.MapKey(kv), vv)
		}

		return protoreflect.ValueOfMap(mv), TryParseValueResult{}, nil
	default:
		// This is internal error, return an error.
		if ctx.ErrHandler != nil {
			return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("internal error: unknown expression type: %T to parse a Message", vt)}, ErrInternal
		}
		return protoreflect.Value{}, TryParseValueResult{}, ErrInternal
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
	case map[string]any:
		// If it is a structpb.Struct, then we need to parse it as a struct.
		if df.Kind() != protoreflect.MessageKind || df.Message().FullName() != "google.protobuf.Struct" {
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)}, ErrInvalidValue
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInvalidValue
		}
		if df.Cardinality() == protoreflect.Repeated {
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)}, ErrInvalidValue
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInvalidValue
		}

		st, err := structpb.NewStruct(et)
		if err != nil {
			if ctx.ErrHandler != nil {
				return protoreflect.Value{}, TryParseValueResult{ErrPos: field.Position(), ErrMsg: fmt.Sprintf("field is of %q type, but provided value is not a valid value: '%s'", df.Kind(), field.Name)}, ErrInvalidValue
			}
			return protoreflect.Value{}, TryParseValueResult{}, ErrInvalidValue
		}

		pv = protoreflect.ValueOfMessage(st.ProtoReflect())
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
