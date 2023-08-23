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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
	blockyannotations "github.com/blockysource/go-genproto/blocky/api/annotations"
)

// FieldDescriptor is an interface that describes a field.
// It can either be a protoreflect.FieldDescriptor or a function argument field descriptor.
type FieldDescriptor interface {
	// Kind returns the field kind.
	Kind() protoreflect.Kind

	// Message returns the message descriptor.
	// If the field is not a message, it returns nil.
	Message() protoreflect.MessageDescriptor

	// Enum returns the enum descriptor.
	// If the field is not an enum, it returns nil.
	Enum() protoreflect.EnumDescriptor

	// IsMap returns true if the field is a map.
	IsMap() bool

	// MapKey returns the map key field descriptor.
	// If the field is not a map, it returns nil.
	MapKey() protoreflect.FieldDescriptor

	// MapValue returns the map value field descriptor.
	// If the field is not a map, it returns nil.
	MapValue() protoreflect.FieldDescriptor

	// Cardinality returns the cardinality of the field.
	Cardinality() protoreflect.Cardinality
}

// TryParseSelectorExpr handles an ast.MemberExpr and returns an expression.
func (b *Interpreter) TryParseSelectorExpr(ctx *ParseContext, value ast.ValueExpr, args ...ast.FieldExpr) (TryParseValueResult, error) {
	// Check if the named expression is a MemberExpr.
	var field protoreflect.FieldDescriptor
	switch vt := value.(type) {
	case *ast.StringLiteral:
		// String member is not supported by default at the named selector side of a restriction expression.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = vt.Pos
			res.ErrMsg = "string literal is not supported as the first element of the named selector expression"
		}
		return res, ErrInvalidField
	case *ast.TextLiteral:
		// The text value should match the field name of the context message descriptor.
		field = ctx.Message.Fields().ByName(protoreflect.Name(vt.Value))
		if field == nil {
			// Check if the field might be in the OneOf descriptors.
			for i := 0; i < ctx.Message.Oneofs().Len(); i++ {
				ood := ctx.Message.Oneofs().Get(i)
				field = ood.Fields().ByName(protoreflect.Name(vt.Value))
				if field != nil {
					break
				}
			}
			if field == nil {
				// No field found with the given name, return error
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = vt.Pos
					res.ErrMsg = fmt.Sprintf("field: %s not found in the message: %s", vt.Value, ctx.Message.Name())
				}
				return res, ErrFieldNotFound
			}
		}

	default:
		// This either is a nil or can't happen at all for invalid ast.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = vt.Position()
			res.ErrMsg = "invalid ast"
		}
		return res, ErrInvalidAST
	}

	if len(args) > 0 && field.Cardinality() == protoreflect.Repeated && !field.IsMap() {
		// Cannot traverse through repeated fields.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = value.Position()
			res.ErrMsg = fmt.Sprintf("field: %q is a repeated field, cannot get nested field", field.Name())
		}
		return res, ErrInvalidValue
	}

	fi := b.getFieldInfo(field)

	if fi.forbidden {
		// Cannot traverse through fields that forbid filtering.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = value.Position()
			res.ErrMsg = fmt.Sprintf("field: %q forbids filtering, cannot get nested field", field.Name())
		}
		return res, ErrInvalidValue
	}

	// A member is a left hand side of a restriction expression thus it should match the field name of the
	// If member has only a Value with no Fields, then we should find non message field in the context message descriptor.
	if len(args) == 0 {
		fe := expr.AcquireFieldSelectorExpr()
		fe.Message = ctx.Message
		fe.Field = field
		fe.FieldComplexity = fi.complexity
		return TryParseValueResult{Expr: fe}, nil
	}

	root := expr.AcquireFieldSelectorExpr()
	root.Message = field.Parent().(protoreflect.MessageDescriptor)
	root.Field = field
	root.FieldComplexity = fi.complexity
	parentFieldX := root
	parent := expr.FilterExpr(root)

	for i := 0; i < len(args); i++ {
		rel := args[i]

		switch pt := parent.(type) {
		case *expr.FieldSelectorExpr:
			if pt.Field.Cardinality() == protoreflect.Repeated && !pt.Field.IsMap() {
				// Cannot traverse through repeated fields.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q is a repeated field, cannot get nested field", pt.Field.Name())
				}
				root.Free()
				return res, ErrInvalidValue
			}

			pfi := b.getFieldInfo(pt.Field)

			if pfi.forbidden {
				// Cannot traverse through fields that forbid filtering.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q forbids filtering, cannot get nested field", pt.Field.Name())
				}
				root.Free()
				return res, ErrInvalidValue
			}
			// Check if the parent field is a message or a map field.
			switch {
			case pt.Field.Kind() == protoreflect.MessageKind && pt.Field.IsMap():
				// Previous field was a map key, thus the current field should be a map value.
				// Thus, current text literal should be a map key.
				// Get the type of the map key, and try parsing the value expression matching the type.
				// If the parsing fails, then return error.
				// If the parsing succeeds, then create a map key expression and set it as the parent.

				mk := pt.Field.MapKey()

				tvi := TryParseValueInput{
					Field:      mk,
					IsNullable: fi.nullable,
					Value:      rel,
				}
				var (
					tvr TryParseValueResult
					err error
				)
				switch mk.Kind() {
				case protoreflect.StringKind:
					// String field never needs more than one value.
					// No args needed.
					tvr, err = b.TryParseStringField(ctx, tvi)
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
					protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					// Integer field never needs more than one value.
					// No args needed.
					tvr, err = b.TryParseSignedIntField(ctx, tvi)
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
					protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					// Unsigned integer field never needs more than one value.
					// No args needed.
					tvr, err = b.TryParseUnsignedIntField(ctx, tvi)
				case protoreflect.BoolKind:
					// Boolean field never needs more than one value.
					// No args needed.
					tvr, err = b.TryParseBooleanField(ctx, tvi)
				case protoreflect.FloatKind, protoreflect.DoubleKind:
					// Float and double might need one first and one arg.
					// Check if we can add the next argument.
					if i+1 < len(args) {
						tvi.Args = append(tvi.Args, args[i+1])
					}

					tvr, err = b.TryParseFloatField(ctx, tvi)
					if tvr.ArgsUsed == 1 {
						i++
					}
				case protoreflect.BytesKind:
					// Bytes field never needs more than one value.
					// No args needed.
					tvr, err = b.TryParseBytesField(ctx, tvi)
				case protoreflect.EnumKind:
					// Enum field don't need more than one value.
					// No args needed.
					tvr, err = b.TryParseEnumField(ctx, tvi)
				default:
					// This can't happen as protobuf doesn't support other types
					// to be a map key.
					// Mark it as internal error and notify
					if ctx.ErrHandler != nil {
						tvr.ErrPos = rel.Position()
						tvr.ErrMsg = fmt.Sprintf("field: %q is not a supported map key type", mk.Name())
					}
					root.Free()
					return tvr, ErrInternal
				}
				if err != nil {
					root.Free()
					return tvr, err
				}

				mke := expr.AcquireMapKeyExpr()
				mke.Key = tvr.Expr
				parentFieldX.Traversal = mke
				parent = mke
			case pt.Field.Kind() == protoreflect.MessageKind:
				// This is a message, thus we can search for the next field in the message.
				// Check if the next value is a text literal.
				tl, ok := rel.(*ast.TextLiteral)
				if !ok {
					var res TryParseValueResult
					// Traversing through a message fields requires a text literal.
					if ctx.ErrHandler != nil {
						res.ErrPos = rel.Position()
						res.ErrMsg = fmt.Sprintf("field: %q is a message type field, nested field traversal requires text literal", pt.Field.Name())
					}
					root.Free()
					return res, ErrInvalidField
				}

				// Check if the text literal value is a valid field in the message.
				field = pt.Message.Fields().ByName(protoreflect.Name(tl.Value))
				if field == nil {
					// Check if the field might be in the OneOf descriptors.
					for i := 0; i < pt.Message.Oneofs().Len(); i++ {
						ood := pt.Message.Oneofs().Get(i)
						field = ood.Fields().ByName(protoreflect.Name(tl.Value))
						if field != nil {
							break
						}
					}
					if field == nil {
						// Field was not found in the message.
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = rel.Position()
							res.ErrMsg = fmt.Sprintf("field: %q not found in the message: %s", tl.Value, pt.Message.Name())
						}
						root.Free()
						return res, ErrFieldNotFound
					}
				}

				if !field.IsMap() && i != len(args)-1 && field.Cardinality() == protoreflect.Repeated {
					// A repeated field cannot be traversed through.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = rel.Position()
						res.ErrMsg = fmt.Sprintf("field: %q is a repeated field, cannot get nested field", field.Name())
					}
					root.Free()
					return res, ErrInvalidValue
				}

				fi = b.getFieldInfo(field)

				// Create a field expression and set it as the parent.
				fe := expr.AcquireFieldSelectorExpr()
				fe.Message = pt.Message
				fe.Field = field
				fe.FieldComplexity = fi.complexity
				parentFieldX.Traversal = fe
				parent = fe
				parentFieldX = fe
			default:
				// This is not a valid field for traversing.
				// Mark as an invalid value error.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q is not a message type field, cannot get nested field", pt.Field.Name())
				}
				root.Free()
				return res, ErrInvalidValue
			}
		case *expr.MapKeyExpr:
			// Previous field was a map key, thus the current field should be a map value.
			// Check if the parent field map value is a message.
			// If it is, then the current field should be a
			msg := parentFieldX.Field.MapValue()
			if msg.Kind() != protoreflect.MessageKind {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q is not a message type field, cannot get nested field", parentFieldX.Field.Name())
				}
				root.Free()
				return res, ErrInvalidValue
			}

			// If the next value is a not a text literal than it is an error.
			tl, ok := rel.(*ast.TextLiteral)
			if !ok {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q is a message type field, nested field traversal requires text literal", parentFieldX.Field.Name())
				}
				root.Free()
				return res, ErrInvalidField
			}

			// Check the value of text literal in the map value message fields.
			field = msg.Message().(protoreflect.MessageDescriptor).
				Fields().ByName(protoreflect.Name(tl.Value))
			if field == nil {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = rel.Position()
					res.ErrMsg = fmt.Sprintf("field: %q not found in the message: %s", tl.Value, msg.Message().Name())
				}
				root.Free()
				return res, ErrFieldNotFound
			}

			fi = b.getFieldInfo(field)

			// Create a field expression and set it as the parent.
			fe := expr.AcquireFieldSelectorExpr()
			fe.Message = msg.Message().(protoreflect.MessageDescriptor)
			fe.Field = field
			fe.FieldComplexity = fi.complexity

			// Set up the traversal in the map key parent expression.
			pt.Traversal = fe

			// Set up the next parent field expression.
			parentFieldX = fe

			// Set up the next parent expression.
			parent = fe
		}
	}
	return TryParseValueResult{Expr: root}, nil
}

// IsFieldFilteringForbidden returns true if the field filtering is forbidden.
func IsFieldFilteringForbidden(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), blockyannotations.E_QueryOpt).([]blockyannotations.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == blockyannotations.FieldQueryOption_FORBID_FILTERING {
			return true
		}
	}
	return false
}
