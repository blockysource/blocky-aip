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

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// HandleRestrictionExpr handles an ast.Restriction expression and returns resulting expr.FilterExpr.
func (b *Interpreter) HandleRestrictionExpr(ctx *ParseContext, x *ast.RestrictionExpr) (TryParseValueResult, error) {
	// Try parsing the inner ComparableExpr
	var left expr.FilterExpr
	switch xt := x.Comparable.(type) {
	case *ast.MemberExpr:
		// Try to get the named selector from the left hand side.
		res, err := b.TryParseSelectorExpr(ctx, xt.Value, xt.Fields...)
		if err != nil {
			return res, err
		}
		left = res.Expr

		// The left hand side is a selector expression.
		// Check if there is a comparator.
		if x.Comparator == nil {
			var res TryParseValueResult
			// No comparator on the selector part is an error.
			if ctx.ErrHandler != nil {
				res.ErrPos = xt.Position()
				res.ErrMsg = "missing comparator in restriction expression"
			}
			left.Free() // Free the selector expression.
			return res, ErrInvalidValue
		}

		cmp, ok := parseComparator(x.Comparator)
		if !ok {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = x.Comparator.Position()
				res.ErrMsg = fmt.Sprintf("unknown comparator: %s", x.Comparator.String())
			}
			left.Free()
			return res, ErrInternal
		}

		field, mk, ok := traverseLastFieldExpr(left)
		if !ok {
			// The left hand side is not a field selector expression.
			// This is an internal error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = xt.Position()
				res.ErrMsg = "internal error: left hand side of restriction expression is not a field selector expression"
			}
			left.Free()
			return res, ErrInternal
		}

		// If the field is a map key and comparator is HAS, check if the right side is a wildcard TEXT literal.
		if cmp == expr.HAS && mk != nil {
			if me, ok := x.Arg.(*ast.MemberExpr); ok {
				if tl, ok := me.Value.(*ast.TextLiteral); ok && len(me.Fields) == 0 && tl.Value == "*" {
					// Modify the expression as a map field selector has a key expression.

					// Extract key from the map key expression.
					ke := mk.Key

					// Clear the map key expression.
					mk.Key = nil
					mk.Free()
					field.Traversal = nil

					// Return a compare expression with the field selector and a key expression.
					ce := expr.AcquireCompareExpr()
					ce.Left = field
					ce.Comparator = cmp
					ce.Right = ke
					return TryParseValueResult{Expr: ce, IsIndirect: true}, nil
				}
			}
		}

		fd := field.Field
		switch {
		case mk != nil:
			// If the left-hand side is a map key expr, set the field descriptor as map value.
			fd = field.Field.MapValue()
		case fd.Kind() == protoreflect.MessageKind && fd.IsMap() && cmp == expr.HAS:
			fd = fd.MapKey()
		}

		// Try getting the value of the right hand side.
		ve, err := b.TryParseValue(ctx, TryParseValueInput{
			Field:         fd,
			Value:         x.Arg,
			AllowIndirect: true,
			IsNullable:    IsFieldNullable(field.Field),
		})
		if err != nil {
			// The right hand side is not a value expression, try parsing it as a selector.
			switch at := x.Arg.(type) {
			case *ast.MemberExpr:
				// Try to get the named selector from the right hand side.
				right, err2 := b.TryParseSelectorExpr(ctx, at.Value, at.Fields...)
				if err2 != nil {
					// The right hand side is neither a value expression nor a selector expression.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side is not a valid value: %s", x.Arg.String())
					}
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if traversal of the right hand side types matches the left hand side types.
				rightField, rmk, ok := traverseLastFieldExpr(right.Expr)
				if !ok {
					// The right hand side is not a field selector expression.
					// This is an internal error.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = at.Position()
						res.ErrMsg = "internal error: right hand side of restriction expression is not a field selector expression"
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInternal
				}

				lf := field.Field
				rf := rightField.Field

				// Check the ambiguity of the left and right hand side.
				// This means that the left hand side ie equal to the right hand side directly.
				// I.e.: field = field
				if lf.FullName() == rf.FullName() && countTraversal(res.Expr) == countTraversal(right.Expr) {
					// This is ambiguous and is not a valid filter.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side is ambiguous: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrAmbiguousField
				}

				var leftIsMapKey bool
				switch {
				case mk != nil:
					// If the left-hand side is a map key expr, set the field descriptor as map value.
					lf = field.Field.MapValue()
				case rmk != nil:
					// If the right-hand side is a map key expr, set the field descriptor as map value.
					rf = rightField.Field.MapValue()
				case lf.Kind() == protoreflect.MessageKind && lf.IsMap():
					// If the left-hand side is a map, set the field descriptor as map key.
					leftIsMapKey = true
					lf = lf.MapKey()
				case rf.Kind() == protoreflect.MessageKind && rf.IsMap():
					// If the right-hand side is a map, set the field descriptor as map key.
					rf = rf.MapKey()
				}

				// This means that the right hand side is a value of the map.
				// We need to check the type of the map value.
				if !isKindComparable(lf.Kind(), rf.Kind()) {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// If the field is an enum or a message matching descriptors.
				if lf.Kind() == protoreflect.EnumKind && lf.Enum().FullName() != rf.Enum().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				} else if lf.Kind() == protoreflect.MessageKind && lf.Message().FullName() != rf.Message().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						ctx.ErrHandler(x.Arg.Position(), fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String()))
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if the left hand side is repeated and the right is not.
				if lf.Cardinality() == protoreflect.Repeated && rf.Cardinality() != protoreflect.Repeated {
					// If the comparator is not HAS, this is an error.
					// I.e. array_field:value
					if x.Comparator.Type != ast.HAS {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							// Invalid value.
							res.ErrPos = x.Arg.Position()
							res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
						}
						right.Expr.Free()
						left.Free()
						return res, ErrInvalidValue
					}
				}

				// Check if the left hand side is neither a map key nor repeated and the operator is HAS.
				if (!leftIsMapKey || lf.Cardinality() != protoreflect.Repeated) && x.Comparator.Type == ast.HAS {
					// If the comparator is HAS and the left hand side is not a map key, this is an error.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Comparator.Position()
						res.ErrMsg = "operator HAS (':') can only be used on map or repeated fields"
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if the right hand side is repeated and the left is not.
				if rf.Cardinality() == protoreflect.Repeated && lf.Cardinality() != protoreflect.Repeated && !lf.IsMap() {
					// If the comparator is different from IN, this is an error.
					if x.Comparator.Type != ast.IN {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							// Invalid value.
							ctx.ErrHandler(x.Comparator.Position(), fmt.Sprintf("cannot compare a repeated field: %s with a non-repeated field: %s with a comparator: %s", rf.FullName(), field.Field.FullName(), x.Comparator.String()))
							res.ErrPos = x.Comparator.Position()
							res.ErrMsg = fmt.Sprintf("cannot compare a repeated field: %s with a non-repeated field: %s with a comparator: %s", rf.FullName(), field.Field.FullName(), x.Comparator.String())
						}
						right.Expr.Free()
						left.Free()
						return res, ErrInvalidValue
					}
				}

				// The selectors should be valid now.
				ex := expr.AcquireCompareExpr()
				ex.Left = field
				ex.Comparator = cmp
				ex.Right = right.Expr
				return TryParseValueResult{Expr: ex, IsIndirect: true}, nil
			case *ast.FunctionCall:
				argFn, ok := b.getFunctionDeclaration(ctx, at)
				if !ok {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = at.Pos
						res.ErrMsg = fmt.Sprintf("function: %s undefined", at.JoinedName())
					}
					return res, ErrInvalidValue
				}

				if argFn.ServiceCall() {
					// This is not a valid value expression.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = at.Pos
						res.ErrMsg = fmt.Sprintf("function: %s can't be as a comparator argument", at.JoinedName())
					}
					left.Free()
					return res, ErrInvalidValue
				} else {
					// Try to match the kind of resulting value with the argument declaration.
					rt := argFn.Returning

					if rt.FieldKind != fd.Kind() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.FieldKind)
						}
						left.Free()
						return res, ErrInvalidValue
					}

					if rt.EnumDescriptor != nil && rt.EnumDescriptor.FullName() != fd.Enum().FullName() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.EnumDescriptor.FullName())
						}
						left.Free()
						return res, ErrInvalidValue
					}

					if fd.Message() != nil && fd.IsMap() && !rt.IsMap() {
						if cmp != expr.HAS {
							var res TryParseValueResult
							if ctx.ErrHandler != nil {
								res.ErrPos = x.Position()
								res.ErrMsg = fmt.Sprintf("function call %s does not return a map value", at.JoinedName())
							}
							left.Free()
							return res, ErrInvalidValue
						}
					}

					if fd.Message() != nil && fd.Message().FullName() != rt.Message().FullName() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.Message().FullName())
						}
						left.Free()
						return res, ErrInvalidValue
					}

					if fd.Cardinality() == protoreflect.Repeated && rt.Cardinality() != protoreflect.Repeated {
						if cmp != expr.IN {
							var res TryParseValueResult
							if ctx.ErrHandler != nil {
								res.ErrPos = x.Position()
								res.ErrMsg = fmt.Sprintf("function call %s is repeated", at.JoinedName())
							}
							left.Free()
							return res, ErrInvalidValue
						}
					}
				}

				// Call right hand side.
				rfn, err := b.tryParseAndCallFunction(ctx, at, argFn, true)
				if err != nil {
					left.Free()
					return rfn, err
				}

				ce := expr.AcquireCompareExpr()
				ce.Left = left
				ce.Comparator = cmp
				ce.Right = rfn.Expr
				return TryParseValueResult{Expr: ce, IsIndirect: res.IsIndirect || rfn.IsIndirect}, nil
			default:
				// The right hand side is not a selector expression.
				// Thus return an error.
				left.Free()
				return ve, err
			}
		}

		// The right hand side is a value expression.
		// Check if the operator matches the expression type.
		switch vt := ve.Expr.(type) {
		case *expr.ValueExpr:
			// The right hand side is a value expression,
			// if the left hand side is a map key or a repeated field, the operator must be HAS.
			if (mk != nil || field.Field.Cardinality() == protoreflect.Repeated) && cmp != expr.HAS {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated field: %s with a comparator: %s", field.Field.FullName(), x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		// The right hand side is a proper value expression.
		case *expr.ArrayExpr:
			// The right hand side is an array expression,
			// check if the left hand side is either a repeated field or a
			// single field with IN comparator.
			if fd.Cardinality() != protoreflect.Repeated && cmp != expr.IN {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated field: %s with a comparator: %s", field.Field.FullName(), x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		case *expr.MapValueExpr:
			// The right hand side is a map value expression,
			// The left hand side must be a map field (NOT a map key).
			if !fd.IsMap() || mk != nil {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Arg.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a map with a non map field: %s", field.Field.FullName())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
			// Comparator can only accept EQ or NEQ.
			if cmp != expr.EQ && cmp != expr.NE {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a map with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		case *expr.StringSearchExpr:
			// The right hand side is a string search expression,
			// The comparator needs to be EQ or IN.
			if cmp != expr.EQ && cmp != expr.IN {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a string search expression with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

			// If the left hand side is repeated field than it is an error.
			if fd.Cardinality() == protoreflect.Repeated {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated field: %s with a comparator: %s", fd.FullName(), x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

		default:
			// The right hand side is not a value expression.
			// Thus return an error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = x.Arg.Position()
				res.ErrMsg = fmt.Sprintf("the right hand side is not a valid value type: %T", ve.Expr)
			}
			left.Free()
			vt.Free()
			return res, ErrInternal
		}

		ce := expr.AcquireCompareExpr()
		ce.Left = field
		ce.Comparator = cmp
		ce.Right = ve.Expr
		return TryParseValueResult{Expr: ce, IsIndirect: true}, nil
	case *ast.FunctionCall:
		fn, ok := b.getFunctionDeclaration(ctx, xt)
		if !ok {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = xt.Pos
				res.ErrMsg = fmt.Sprintf("function: %s undefined", xt.JoinedName())
			}
			return res, ErrInvalidValue
		}

		res, err := b.tryParseAndCallFunction(ctx, xt, fn, true)
		if err != nil {
			return res, err
		}

		left = res.Expr

		if !fn.ServiceCall() && !res.IsIndirect {
			// The left hand side is not an indirect form.
			// The left hand side of the restriction needs to be indirect.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = xt.Pos
				res.ErrMsg = fmt.Sprintf("function: %s does not depend on the filtered message", xt.JoinedName())
			}
			left.Free()
			return res, ErrInvalidValue
		}

		if fn.ServiceCall() {
			// The result should be an expr.FunctionCallExpr.
			// Check if there is a comparator and an argument.
			if x.Comparator != nil || x.Arg != nil {
				// Service calls cannot have a comparator or an argument.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = xt.Pos
					res.ErrMsg = fmt.Sprintf("function: %s cannot have a comparator or an argument", xt.JoinedName())
				}
				left.Free()
				return res, ErrInvalidValue
			}
			return res, nil
		}

		ad := fn.Returning

		// The result should be an indirect expr.FunctionCallExpr,
		if x.Comparator == nil || x.Arg == nil {
			return res, nil
		}

		// Parse comparator.
		cmp, ok := parseComparator(x.Comparator)
		if !ok {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = x.Comparator.Position()
				res.ErrMsg = fmt.Sprintf("unknown comparator: %s", x.Comparator.String())
			}
			left.Free()
			return res, ErrInternal
		}

		switch at := x.Arg.(type) {
		case *ast.MemberExpr, *ast.StructExpr, *ast.ArrayExpr:
		case *ast.FunctionCall:
			argFn, ok := b.getFunctionDeclaration(ctx, at)
			if !ok {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = at.Pos
					res.ErrMsg = fmt.Sprintf("function: %s undefined", at.JoinedName())
				}
				return res, ErrInvalidValue
			}

			if argFn.ServiceCall() {
				// This is not a valid value expression.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = at.Pos
					res.ErrMsg = fmt.Sprintf("function: %s can't be as a comparator argument", at.JoinedName())
				}
				left.Free()
				return res, ErrInvalidValue
			} else {
				// Try to match the kind of resulting value with the argument declaration.
				rt := argFn.Returning

				if rt.FieldKind != ad.FieldKind {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.FieldKind)
					}
					left.Free()
					return res, ErrInvalidValue
				}

				if rt.EnumDescriptor != nil && rt.EnumDescriptor.FullName() != ad.EnumDescriptor.FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.EnumDescriptor.FullName())
					}
					left.Free()
					return res, ErrInvalidValue
				}

				if ad.Message() != nil && ad.IsMap() && !rt.IsMap() {
					if cmp != expr.HAS {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s does not return a map value", at.JoinedName())
						}
						left.Free()
						return res, ErrInvalidValue
					}
				}

				if ad.Message() != nil && ad.Message().FullName() != rt.Message().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s is not of type %s", at.JoinedName(), rt.Message().FullName())
					}
					left.Free()
					return res, ErrInvalidValue
				}

				if ad.Cardinality() == protoreflect.Repeated && rt.Cardinality() != protoreflect.Repeated {
					if cmp != expr.IN {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s is repeated", at.JoinedName())
						}
						left.Free()
						return res, ErrInvalidValue
					}
				}
			}

			// Call right hand side.
			rfn, err := b.tryParseAndCallFunction(ctx, at, argFn, true)
			if err != nil {
				left.Free()
				return rfn, err
			}

			ce := expr.AcquireCompareExpr()
			ce.Left = left
			ce.Comparator = cmp
			ce.Right = rfn.Expr
			return TryParseValueResult{Expr: ce, IsIndirect: res.IsIndirect || rfn.IsIndirect}, nil
		case *ast.CompositeExpr:
			// Handle composite if only the left hand side is a boolean expression.
			if fn.ServiceCall() || (fn.Returning.FieldKind != protoreflect.BoolKind || fn.Returning.Cardinality() == protoreflect.Repeated) {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Arg.Position()
					res.ErrMsg = fmt.Sprintf("function: %s cannot have a composite argument", xt.JoinedName())
				}
				left.Free()
				return res, ErrInvalidValue
			}

			right, err := b.HandleCompositeExpr(ctx, at)
			if err != nil {
				left.Free()
				return right, err
			}

			// The right hand side is a composite expression.
			ce := expr.AcquireCompareExpr()
			ce.Left = left
			ce.Comparator = cmp
			ce.Right = right.Expr
			return TryParseValueResult{Expr: ce, IsIndirect: res.IsIndirect || right.IsIndirect}, nil
		default:
			// Not a valid type for right hand side, internal error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = x.Arg.Position()
				res.ErrMsg = fmt.Sprintf("the right hand side is not a valid value type: %T", x.Arg)
			}
			left.Free()
			return res, ErrInternal
		}

		// Parse argument.
		ve, err := b.TryParseValue(ctx, TryParseValueInput{
			Field:         fn.Returning,
			AllowIndirect: true,
			Value:         x.Arg,
			IsNullable:    fn.Returning.IsNullable,
		})
		if err != nil {
			// If the right hand side is not a value expression,
			// try parsing it as a field selector.
			// The right hand side is not a value expression, try parsing it as a selector.
			switch at := x.Arg.(type) {
			case *ast.MemberExpr:
				// Try to get the named selector from the right hand side.
				right, err2 := b.TryParseSelectorExpr(ctx, at.Value, at.Fields...)
				if err2 != nil {
					// The right hand side is neither a value expression nor a selector expression.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side is not a valid value: %s", x.Arg.String())
					}
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if traversal of the right hand side types matches the left hand side types.
				rightField, rmk, ok := traverseLastFieldExpr(right.Expr)
				if !ok {
					// The right hand side is not a field selector expression.
					// This is an internal error.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = at.Position()
						res.ErrMsg = "internal error: right hand side of restriction expression is not a field selector expression"
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInternal
				}

				lf := fn.Returning
				rf := rightField.Field

				switch {
				case rmk != nil:
					// If the right-hand side is a map key expr, set the field descriptor as map value.
					rf = rightField.Field.MapValue()
				case rf.Kind() == protoreflect.MessageKind && rf.IsMap():
					// If the right-hand side is a map, set the field descriptor as map key.
					rf = rf.MapKey()
				}

				// This means that the right hand side is a value of the map.
				// We need to check the type of the map value.
				if isKindComparable(lf.Kind(), rf.Kind()) {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// If the field is an enum or a message matching descriptors.
				if lf.Kind() == protoreflect.EnumKind && lf.Enum().FullName() != rf.Enum().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				} else if lf.Kind() == protoreflect.MessageKind && lf.Message().FullName() != rf.Message().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if the left hand side is repeated and the right is not.
				if lf.IsRepeated && rf.Cardinality() != protoreflect.Repeated {
					// If the comparator is not HAS, this is an error.
					// I.e. array_field:value
					if x.Comparator.Type != ast.HAS {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							// Invalid value.
							res.ErrPos = x.Arg.Position()
							res.ErrMsg = fmt.Sprintf("the right hand side type of the restriction doesn't match the left hand side type: %s", x.Arg.String())
						}
						right.Expr.Free()
						left.Free()
						return res, ErrInvalidValue
					}
				}

				// Check if the left hand side is neither a map key nor repeated and the operator is HAS.
				if !lf.IsRepeated && x.Comparator.Type == ast.HAS {
					// If the comparator is HAS and the left hand side is not a map key, this is an error.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						// Invalid value.
						res.ErrPos = x.Comparator.Position()
						res.ErrMsg = "operator HAS (':') can only be used on map or repeated fields"
					}
					right.Expr.Free()
					left.Free()
					return res, ErrInvalidValue
				}

				// Check if the right hand side is repeated and the left is not.
				if rf.Cardinality() == protoreflect.Repeated && !lf.IsRepeated {
					// If the comparator is different from IN, this is an error.
					if x.Comparator.Type != ast.IN {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							// Invalid value.
							res.ErrPos = x.Comparator.Position()
							res.ErrMsg = fmt.Sprintf("cannot compare a repeated field: %s with a non-repeated field with a comparator: %s", rf.FullName(), x.Comparator.String())
						}
						right.Expr.Free()
						left.Free()
						return res, ErrInvalidValue
					}
				}

				// The selectors should be valid now.
				ex := expr.AcquireCompareExpr()
				ex.Left = left
				ex.Comparator = cmp
				ex.Right = right.Expr
				return TryParseValueResult{Expr: ex, IsIndirect: true}, nil
			default:
				// The right hand side is not a selector expression.
				// Thus return an error.
				left.Free()
				return ve, err
			}
		}

		// The right hand side is a value expression.
		// Check if the operator matches the expression type.
		switch vt := ve.Expr.(type) {
		case *expr.ValueExpr:
			// The right hand side is a value expression,
			// if the left hand side is a repeated field, the operator must be IN.
			if fn.Returning.IsRepeated && cmp != expr.IN {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated field with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

			// If the left hand side is a map return an error, as ValueExpr cannot be a map.
			if fn.Returning.IsMap() && cmp != expr.HAS {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a map with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		case *expr.ArrayExpr:
			// The right hand side is an array expression,
			// check if the left hand side is either a repeated field or a
			// single field with IN comparator.
			if !fn.Returning.IsRepeated && cmp != expr.IN {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated field with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		case *expr.MapValueExpr:
			// If the left hand side is not a map and comparator is neither EQ nor NEQ, return an error.
			if !fn.Returning.IsMap() || (cmp != expr.EQ && cmp != expr.NE) {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a map with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}
		case *expr.FunctionCallExpr:
			// The right hand side is a function call expression,
			// Check if the returning value of the function call matches the left hand side.
			vfn := b.functionCallDeclarations[vt.FullName()]
			if fn.Returning.Kind() != vfn.Returning.Kind() {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Arg.Position()
					res.ErrMsg = fmt.Sprintf("resulting function call type: %s does not match the left hand side type: %s", vfn.Returning.Kind().String(), fn.Returning.Kind().String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

			if fn.Returning.Kind() == protoreflect.EnumKind && fn.Returning.Enum().FullName() != vfn.Returning.Enum().FullName() {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Arg.Position()
					res.ErrMsg = fmt.Sprintf("resulting function call enum type: %s does not match the left hand side enum type: %s", vfn.Returning.Enum().FullName(), fn.Returning.Enum().FullName())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

			if fn.Returning.Kind() == protoreflect.MessageKind {
				if (fn.Returning.IsMap() && !vfn.Returning.IsMap()) || (!fn.Returning.IsMap() && vfn.Returning.IsMap()) {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("resulting function call message type: %s does not match the left hand side message type: %s", vfn.Returning.Message().FullName(), fn.Returning.Message().FullName())
					}
					left.Free()
					vt.Free()
					return res, ErrInvalidValue
				}
				if fn.Returning.IsMap() && (fn.Returning.MapKey().FullName() != vfn.Returning.MapKey().FullName() ||
					fn.Returning.MapValue().FullName() != vfn.Returning.MapValue().FullName()) {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("resulting function call message type: %s does not match the left hand side message type: %s", vfn.Returning.Message().FullName(), fn.Returning.Message().FullName())
					}
					left.Free()
					vt.Free()
					return res, ErrInvalidValue
				}

				if fn.Returning.IsMap() && cmp != expr.EQ && cmp != expr.NE {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Comparator.Position()
						res.ErrMsg = fmt.Sprintf("cannot compare a map with a comparator: %s", x.Comparator.String())
					}
					left.Free()
					vt.Free()
					return res, ErrInvalidValue
				}

				if !fn.Returning.IsMap() && fn.Returning.Message().FullName() != vfn.Returning.Message().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Arg.Position()
						res.ErrMsg = fmt.Sprintf("resulting function call message type: %s does not match the left hand side message type: %s", vfn.Returning.Message().FullName(), fn.Returning.Message().FullName())
					}
					left.Free()
					vt.Free()
					return res, ErrInvalidValue
				}
			}
		case *expr.StringSearchExpr:
			// The right hand side is a string search expression,
			// The comparator needs to be EQ or IN.
			if cmp != expr.EQ && cmp != expr.IN {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a string search expression with a comparator: %s", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

			// If the left hand side is repeated field than it is an error.
			if fn.Returning.Cardinality() == protoreflect.Repeated {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Comparator.Position()
					res.ErrMsg = fmt.Sprintf("cannot compare a repeated function result with a comparator: %s for string search", x.Comparator.String())
				}
				left.Free()
				vt.Free()
				return res, ErrInvalidValue
			}

		default:
			// The right hand side is not a value expression.
			// Thus return an error.
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				// This is invalid  value error ?
				res.ErrPos = x.Arg.Position()
				res.ErrMsg = "invalid value for restriction expression"
			}
			left.Free()
			vt.Free()
			return res, ErrInvalidValue
		}

		// The selectors should be valid now.
		ex := expr.AcquireCompareExpr()
		ex.Left = left
		ex.Comparator = cmp
		ex.Right = ve.Expr
		return TryParseValueResult{Expr: ex, IsIndirect: ve.IsIndirect}, nil
	default:
		// The left hand side is not a selector expression.
		// This is invalid value error.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = xt.Position()
			res.ErrMsg = fmt.Sprintf("the left hand side is not a valid selector: %s", xt.String())
		}
		return res, ErrInvalidValue
	}
}

func parseComparator(in *ast.ComparatorLiteral) (expr.Comparator, bool) {
	var comp expr.Comparator
	switch in.Type {
	case ast.EQ:
		comp = expr.EQ
	case ast.LE:
		comp = expr.LE
	case ast.LT:
		comp = expr.LT
	case ast.GE:
		comp = expr.GE
	case ast.GT:
		comp = expr.GT
	case ast.NE:
		comp = expr.NE
	case ast.HAS:
		comp = expr.HAS
	case ast.IN:
		comp = expr.IN
	default:
		return 0, false
	}
	return comp, true
}

func traverseLastFieldExpr(in expr.FilterExpr) (*expr.FieldSelectorExpr, *expr.MapKeyExpr, bool) {
	var (
		fe *expr.FieldSelectorExpr
		mk *expr.MapKeyExpr
	)
	e := in
	for {
		switch xt := e.(type) {
		case *expr.FieldSelectorExpr:
			fe = xt
			if xt.Traversal == nil {
				return fe, nil, true
			}
			e = xt.Traversal
		case *expr.MapKeyExpr:
			if xt.Traversal == nil {
				return fe, xt, true
			}
			mk = xt
			e = xt.Traversal
		default:
			return fe, mk, true
		}
	}
}

func countTraversal(in expr.FilterExpr) int {
	var count int
	e := in
	for {
		switch xt := e.(type) {
		case *expr.FieldSelectorExpr:
			if xt.Traversal == nil {
				return count
			}
			e = xt.Traversal
		case *expr.MapKeyExpr:
			if xt.Traversal == nil {
				return count
			}
			e = xt.Traversal
		default:
			return count
		}
		count++
	}
}
