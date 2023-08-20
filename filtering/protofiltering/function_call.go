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
	"errors"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

type (
	// FunctionCallFn is a function call handler function.
	// It is used to handle a function call ast node, and return an expression.
	// It either can return indirect expr.FunctionCall expression (in case of indirect argument),
	// or a value expression.
	// An indirect call might be used to perform a function call on a service
	// directly outside of the filtering layer.
	FunctionCallFn func(args ...expr.FilterExpr) (FunctionCallArgument, error)

	// FunctionCallArgument is an argument of the function call.
	// It specifies the expression injected into the function, and shows information if the
	// value depends on the resource message.
	FunctionCallArgument struct {
		// Expr is the expression of the argument.
		Expr expr.FilterExpr

		// IsIndirect specifies if the argument is indirect.
		// If the argument is indirect the value provided to the argument depends on value of the filtered message.
		// This means that the function call should result in an abstract expr.FunctionCall expression,
		// and does not result with a value expression.
		IsIndirect bool
	}

	// FunctionCallDeclaration is a declaration of a function call.
	// It is used to register a function call handler.
	// A function call may either return a value or an indirect value.
	// If it returns an indirect value, then it returns an expr.FunctionCall or expr.FieldSelector expression.
	//
	// i.e. a function call geo.Distance(pt1, pt2) might return either
	// a direct value if both pt1 and pt2 are direct values, or
	// indirect value (expr.FunctionCall) if at least one of them is indirect.
	// I.e.
	//  - geo.Distance(loc, geo.Point(1, 2)) returns an indirect value.,
	// 		as the first argument is field selector, and the second is a direct value.
	//  - geo.Distance(geo.Point(1, 2), geo.Point(3, 4)) returns a direct value,
	// 		as both arguments are direct values, this would result in returning
	//      value defined by the geo.Distance function.
	//  - geo.Point(x, y) - returns an indirect value, as both arguments are indirect.
	//  - geo.Point(1, 2) - returns a direct value, as both arguments are direct.
	//
	// If the function call expected result type is a boolean, then it could be used as a sole comparable expression in the
	// restriction expression.
	// i.e.:
	//  - geo.InArea(loc, area) - returns an indirect value, as both arguments are indirect,
	// 		but it can be used as a comparable expression in the restriction expression.
	//      this function needs to define its returning declaration to be of kind boolean.
	// If the function does not contain a returning declaration, then it is service called function expression,
	// and it needs to be validated by the service.
	FunctionCallDeclaration struct {
		// Name is a unique identifier of the function call.
		Name FunctionName

		// Arguments is a list of arguments of the function call.
		// If empty then the function call has no arguments.
		Arguments []*FunctionCallArgumentDeclaration

		// Returning is an optional returning declaration of the function call.
		// If this field is undefined, this function always returns an expr.FunctionCall expression.
		// This is then named as a service called function call.
		// It can be used for any purpose of the filter, but needs to be validated on the service side.
		// A function call declaration without a returning declaration cannot be used as an argument of another function call.
		Returning *FunctionCallReturningDeclaration

		// CallFn is the execution function of the function call.
		// It is called when the function call is executed.
		CallFn FunctionCallFn
	}
	// FunctionName is the name of the function call.
	FunctionName struct {
		// PkgName is the name of the package where the function is defined.
		PkgName string
		// Name is the name of the function call.
		Name string
	}
	// FunctionCallArgumentDeclaration is a declaration of a function call argument.
	FunctionCallArgumentDeclaration struct {
		// Indirect is true if the argument might take a non value
		// indirect form of the filtered message field.
		// The function call with indirect argument returns
		// an expr.FunctionCall expression.
		Indirect bool

		// ArgName is the name of the field.
		ArgName string

		// IsRepeated is true if the argument is a repeated field.
		IsRepeated bool

		// IsNullable is true if the argument is a nullable field.
		IsNullable bool

		// AllowedServiceCallFuncs is a list of function names that can be used as indirect argument.
		// This is used to allow a function call to be used as an indirect argument.
		AllowedServiceCallFuncs []FunctionName

		// FieldKind is the kind of the argument.
		FieldKind protoreflect.Kind

		// EnumDescriptor is the enum descriptor of the argument.
		EnumDescriptor protoreflect.EnumDescriptor

		// If FieldKind is MessageKind, then this is the message descriptor.
		MessageDescriptor protoreflect.MessageDescriptor

		// MapKeyDesc is the map key descriptor of the value returned by the function call.
		// This must be set if the resulting value is a map.
		MapKeyDesc protoreflect.FieldDescriptor

		// MapValueDesc is the map value descriptor of the value returned by the function call.
		// This must be set if the resulting value is a map.
		MapValueDesc protoreflect.FieldDescriptor
	}

	// FunctionCallReturningDeclaration is a declaration of a function call returning value.
	// It either tell if the function is a Service Called abstract function call,
	// or it specifies the kind of the returning value.
	// The returning value implements the FieldDescriptor interface.
	FunctionCallReturningDeclaration struct {
		// ServiceCalled is true if the function call returns a expr.FunctionCall expression.
		// This enforces the function call to return an expr.FunctionCall expression.
		// A service called boolean excludes the possibility to return a direct value.
		ServiceCalled bool

		// FieldKind determines the kind of direct value returned by direct the function call.
		FieldKind protoreflect.Kind

		// EnumDescriptor is the enum descriptor of the direct value returned by the function call.
		EnumDescriptor protoreflect.EnumDescriptor

		// MessageDescriptor is the message descriptor of the direct value returned by the function call.
		MessageDescriptor protoreflect.MessageDescriptor

		// MapKeyDesc is the map key descriptor of the value returned by the function call.
		// This must be set if the resulting value is a map.
		MapKeyDesc protoreflect.FieldDescriptor

		// MapValueDesc is the map value descriptor of the value returned by the function call.
		// This must be set if the resulting value is a map.
		MapValueDesc protoreflect.FieldDescriptor

		// IsNullable is true if the function call returns a nullable value.
		IsNullable bool

		// IsRepeated is true if the function call returns a repeated value.
		IsRepeated bool
	}
)

func (n FunctionName) String() string {
	return fmt.Sprintf("%s.%s", n.PkgName, n.Name)
}

// ServiceCall returns true if the function call is a service call.
func (f *FunctionCallDeclaration) ServiceCall() bool {
	return f.Returning == nil
}

// Validate validates the function call declaration.
func (f *FunctionCallDeclaration) Validate() error {
	if f.Name.PkgName == "" && f.Name.Name == "" {
		return errors.New("undefined function name")
	}

	// Validate arguments.
	for i, arg := range f.Arguments {
		if err := arg.Validate(); err != nil {
			return fmt.Errorf("fn: %s, arg %d: %w", f.Name, i, err)
		}
	}

	// Validate returning field.
	if f.Returning != nil {
		if err := f.Returning.Validate(); err != nil {
			return fmt.Errorf("fn: %s, returning: %w", f.Name, err)
		}
	}
	return nil
}

// Compile-time check that *FunctionCallArgumentDeclaration implements FieldDescriptor.
var _ FieldDescriptor = (*FunctionCallArgumentDeclaration)(nil)

// Name returns the name of the argument.
func (f *FunctionCallArgumentDeclaration) Name() protoreflect.Name {
	return protoreflect.Name(f.ArgName)
}

// Kind returns the kind of returing value.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) Kind() protoreflect.Kind {
	return f.FieldKind
}

// Enum returns the enum descriptor of the argument.
// If returning value is not an enum, then it returns nil.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) Enum() protoreflect.EnumDescriptor {
	return f.EnumDescriptor
}

// Message returns the message descriptor of the argument.
// If returning value is not a message, then it returns nil.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) Message() protoreflect.MessageDescriptor {
	return f.MessageDescriptor
}

// Cardinality returns the cardinality of the argument.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) Cardinality() protoreflect.Cardinality {
	if f.IsRepeated {
		return protoreflect.Repeated
	}
	return protoreflect.Optional
}

// IsMap returns true if the argument is a map.
func (f *FunctionCallArgumentDeclaration) IsMap() bool {
	return f.MapKeyDesc != nil && f.MapValueDesc != nil
}

// MapKey returns the map key descriptor of the argument.
// Returns nil if the argument is not a map.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) MapKey() protoreflect.FieldDescriptor {
	if !f.IsMap() {
		return nil
	}
	return f.MapKeyDesc
}

// MapValue returns the map value descriptor of the argument.
// Returns nil if the argument is not a map.
// Implements FieldDescriptor interface.
func (f *FunctionCallArgumentDeclaration) MapValue() protoreflect.FieldDescriptor {
	if !f.IsMap() {
		return nil
	}
	return f.MapValueDesc
}

// Validate validates the function call argument declaration.
func (f *FunctionCallArgumentDeclaration) Validate() error {
	if f.FieldKind == protoreflect.Kind(0) {
		return errors.New("undefined field kind")
	}

	if f.FieldKind == protoreflect.MessageKind && f.MessageDescriptor == nil {
		return errors.New("undefined message descriptor")
	}

	if f.FieldKind == protoreflect.EnumKind && f.EnumDescriptor == nil {
		return errors.New("undefined enum descriptor")
	}

	if f.MapKeyDesc != nil && f.MapValueDesc == nil || f.MapKeyDesc == nil && f.MapValueDesc != nil {
		return errors.New("undefined map key or value descriptor")
	}

	if f.IsRepeated && f.IsNullable {
		return errors.New("repeated field cannot be nullable")
	}
	return nil
}

// Compile-time check that *FunctionCallReturningDeclaration implements FieldDescriptor.
var _ FieldDescriptor = (*FunctionCallReturningDeclaration)(nil)

// Kind returns the kind of the returning value.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) Kind() protoreflect.Kind {
	return f.FieldKind
}

// Enum returns the enum descriptor of the returning value.
// If the returning value is not an enum, then it returns nil.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) Enum() protoreflect.EnumDescriptor {
	return f.EnumDescriptor
}

// Message returns the message descriptor of the returning value.
// If the returning value is not a message, then it returns nil.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) Message() protoreflect.MessageDescriptor {
	return f.MessageDescriptor
}

// Cardinality returns the cardinality of the returning value.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) Cardinality() protoreflect.Cardinality {
	if f.IsRepeated {
		return protoreflect.Repeated
	}
	if f.IsNullable {
		return protoreflect.Optional
	}
	return protoreflect.Required
}

// IsMap returns true if the returning value is a map.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) IsMap() bool {
	return f.MapKeyDesc != nil && f.MapValueDesc != nil
}

// MapKey returns the map key descriptor of the returning value.
// It returns nil if the returning value is not a map.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) MapKey() protoreflect.FieldDescriptor {
	if !f.IsMap() {
		return nil
	}
	return f.MapKeyDesc
}

// MapValue returns the map value descriptor of the returning value.
// It returns nil if the returning value is not a map.
// Implements FieldDescriptor interface.
func (f *FunctionCallReturningDeclaration) MapValue() protoreflect.FieldDescriptor {
	if !f.IsMap() {
		return nil
	}
	return f.MapValueDesc
}

func (f *FunctionCallReturningDeclaration) Validate() error {
	if f.ServiceCalled {
		return nil
	}
	if f.FieldKind == protoreflect.Kind(0) {
		return errors.New("undefined field kind")
	}

	if f.FieldKind == protoreflect.MessageKind && f.MessageDescriptor == nil {
		return errors.New("undefined message descriptor")
	}

	if f.FieldKind == protoreflect.EnumKind && f.EnumDescriptor == nil {
		return errors.New("undefined enum descriptor")
	}

	if f.MapKeyDesc != nil && f.MapValueDesc == nil || f.MapKeyDesc == nil && f.MapValueDesc != nil {
		return errors.New("undefined map key or value descriptor")
	}

	if f.IsRepeated && f.IsNullable {
		return errors.New("repeated field cannot be nullable")
	}
	return nil
}

// TryParseFunctionCall handles an ast.FunctionCall by interpreting the function call,
// and executing the function call handler.
// It returns either resulting expression value, an expr.FunctionCall expression, or an error.
func (b *Interpreter) TryParseFunctionCall(ctx *ParseContext, in TryParseValueInput) (TryParseValueResult, error) {
	// Ensure that the value is defined.
	if in.Value == nil {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrMsg = fmt.Sprintf("function call expected, got nil")
		}
		return res, ErrInternal
	}

	// Type check that the input value is a function call.
	x, ok := in.Value.(*ast.FunctionCall)
	if !ok {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = in.Value.Position()
			res.ErrMsg = fmt.Sprintf("function call expected, got %T", in.Value)
		}
		return res, ErrInternal
	}

	// Get stored function call declarations that matches given function call.
	fn, found := b.getFunctionDeclaration(ctx, x)
	if !found {
		// No matching function call declaration found.
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("function call %s not found", x.JoinedName())
		}
		return res, ErrInvalidValue
	}

	// With a valid declaration, try parsing it and call the function call expression.
	return b.tryParseAndCallFunction(ctx, x, fn, in.AllowIndirect)
}

func (b *Interpreter) tryParseAndCallFunction(ctx *ParseContext, x *ast.FunctionCall, fn *FunctionCallDeclaration, allowIndirect bool) (TryParseValueResult, error) {
	// We have a function call handler.
	// Parse the argument fields and check if they match the function call declaration.
	// If they do, then we can call the function call handler.
	// Otherwise we return an error.
	if (x.ArgList == nil || len(x.ArgList.Args) == 0) && len(fn.Arguments) > 0 {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("function call %s has no arguments", x.JoinedName())
		}
		return TryParseValueResult{}, ErrNoHandlerFound
	}

	// Check if the function is an abstract function call, and the input does not allow indirect function call.
	if !allowIndirect && fn.ServiceCall() {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("indirect function call %s is not available in that context", x.JoinedName())
		}
		return TryParseValueResult{}, ErrNoHandlerFound
	}

	// If no arguments are provided and the function call does not need any arguments, then we can call the function call handler.
	if (x.ArgList == nil || len(x.ArgList.Args) == 0) && len(fn.Arguments) == 0 {
		// No arguments, so we can call the function call handler.
		ex, err := fn.CallFn()
		if err != nil {
			var res TryParseValueResult
			if ctx.ErrHandler != nil {
				res.ErrPos = x.Position()
				res.ErrMsg = fmt.Sprintf("function call %s failed: %v", x.JoinedName(), err)
			}
			return TryParseValueResult{}, err
		}
		return TryParseValueResult{Expr: ex.Expr, IsIndirect: ex.IsIndirect}, nil
	}

	// Check if the number of arguments match the function call declaration.
	if len(x.ArgList.Args) != len(fn.Arguments) {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("function call %s needs exactly %d arguments", x.JoinedName(), len(fn.Arguments))
		}
		return res, ErrInvalidValue
	}

	var args []expr.FilterExpr
	clearArgs := func() {
		for _, arg := range args {
			arg.Free()
		}
	}
	var isIndirect bool
	// We have arguments, so we need to parse them.
	// We need to check if the arguments match the function call declaration.
	for i, arg := range x.ArgList.Args {
		// Get ith argument declaration.
		ad := fn.Arguments[i]

		// Switch the type of the input argument.
		switch at := arg.(type) {
		case *ast.CompositeExpr:
			// A composite expression might be an argument to the function call, if the function call declaration
			// allows boolean kind of returning value.
			if ad.FieldKind != protoreflect.BoolKind || ad.Cardinality() == protoreflect.Repeated {
				// A composite expression can only be used as an argument if the function call declaration
				// allows boolean kind of returning value.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument is not of a boolean type, thus composite expression is not a valid argument", x.JoinedName())
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			// A composite expression is a boolean expression, so we can use it as an argument.
			right, err := b.HandleCompositeExpr(ctx, at)
			if err != nil {
				clearArgs()
				return right, err
			}
			isIndirect = isIndirect || right.IsIndirect
			args = append(args, right.Expr)
			continue
		case *ast.MemberExpr:
			// A member expression is either a selector or a value expression.
			// If the direct argument is a repeated field, then it is not a valid argument.
			if !ad.Indirect && ad.Cardinality() == protoreflect.Repeated {
				// A repeated field needs to be an ArrayExpr.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument %d must be an array", x.JoinedName(), i)
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			var (
				res TryParseValueResult
				err error
			)

			// If the argument is indirect, try to get the selector.
			if ad.Indirect {
				// A selector expression can be a valid repeated.
				res, err := b.TryParseSelectorExpr(ctx, at.Value, at.Fields...)
				if err == nil {
					// Ensure that the type of the selector matches the type of the argument.
					field, mk, ok := traverseLastFieldExpr(res.Expr)
					if !ok {
						// The selector is not a valid field selector.
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not a valid field selector", x.JoinedName(), i)
						}
						clearArgs()
						return res, ErrInternal
					}

					fd := field.Field
					if mk != nil {
						fd = field.Field.MapValue()
					}

					// Check if the field type matches.
					if !isKindComparable(fd.Kind(), ad.Kind()) {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, ad.FieldKind)
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					// Check if the field is repeated.
					// Check if both fields are of the same cardinality.
					if fd.Cardinality() == protoreflect.Repeated && ad.Cardinality() != protoreflect.Repeated {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not repeated", x.JoinedName(), i)
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					// Check if the input argument is repeated, but the field is not.
					if ad.Cardinality() == protoreflect.Repeated && fd.Cardinality() != protoreflect.Repeated {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s input argument %d is repeated but it should not be", x.JoinedName(), i)
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					// Check if the field is a map.
					if ad.IsMap() && !fd.IsMap() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not a map", x.JoinedName(), i)
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					if fd.IsMap() && !ad.IsMap() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is a map", x.JoinedName(), i)
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					if ad.Kind() == protoreflect.MessageKind && ad.Message().FullName() != fd.Message().FullName() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, ad.Message().FullName())
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					if ad.Kind() == protoreflect.EnumKind && ad.Enum().FullName() != fd.Enum().FullName() {
						var res TryParseValueResult
						if ctx.ErrHandler != nil {
							res.ErrPos = x.Position()
							res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, ad.Enum().FullName())
						}
						clearArgs()
						return res, ErrInvalidValue
					}

					isIndirect = true
					args = append(args, res.Expr)
					continue
				}

				// Check if the field is not repeated.
				if ad.Cardinality() == protoreflect.Repeated {
					// A repeated field needs to be an ArrayExpr.
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d must be an array", x.JoinedName(), i)
					}
					clearArgs()
					return res, ErrInvalidValue
				}
			}

			// In case of indirect argument, the selector was not a valid field selector.
			// Try to parse the member expression as a value expression.
			res, err = b.TryParseValue(ctx, TryParseValueInput{
				Field:         ad,
				Value:         at.Value,
				Args:          at.Fields,
				AllowIndirect: ad.Indirect,
			})
			if err != nil {
				// The member expression is not a valid value expression.
				return res, err
			}

			// Check if the cardinality of the field matches the cardinality of the argument.
			// A member valid value results are only value expressions.
			switch et := res.Expr.(type) {
			case *expr.ArrayExpr:
				// Ensure that the input argument is a repeated field.
				if ad.Cardinality() != protoreflect.Repeated {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not repeated", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}
			case *expr.ValueExpr:
				if ad.Cardinality() == protoreflect.Repeated {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is repeated", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}
				// Check if the argument accepts null values and the value is null.
				if !ad.IsNullable && et.Value == nil {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not nullable", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}
			case *expr.StringSearchExpr:
				if ad.Cardinality() == protoreflect.Repeated {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is repeated", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}

				if !ad.Indirect {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not indirect", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}
			case *expr.MapValueExpr:
				if ad.Cardinality() == protoreflect.Repeated {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is repeated", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}

				if !ad.IsMap() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not a map", x.JoinedName(), i)
					}
					clearArgs()
					res.Expr.Free()
					return res, ErrInvalidValue
				}
			default:
				// What else can it be?
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("unknown expression type %T as function call member argument", res.Expr)
				}
				clearArgs()
				res.Expr.Free()
				return res, ErrInternal
			}

			// Can a member be indirect value if the selector is not a valid field selector?
			// Probably not.

			isIndirect = isIndirect || res.IsIndirect
			args = append(args, res.Expr)
			continue
		case *ast.FunctionCall:
			// Ensure the function call matches the declaration resulting in a valid argument.
			argFn, ok := b.getFunctionDeclaration(ctx, at)
			if !ok {
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument %d not found", x.JoinedName(), i)
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			// An indirect argument might take both direct and indirect function call arguments
			// if only resulting type or accepted indirect function call matches with the declaration.
			if argFn.ServiceCall() {
				// ServiceCall function call means that the argument declaration need to match the function name of the argument.
				// If the argument declaration has no accepted indirect function call, then it is not a valid argument.
				found := false
				for _, aif := range ad.AllowedServiceCallFuncs {
					if aif.PkgName == argFn.Name.PkgName && aif.Name == argFn.Name.Name {
						found = true
						break
					}
				}
				if !found {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d does not allow to use indirect function call %s.%s", x.JoinedName(), i, argFn.Name.PkgName, argFn.Name.Name)
					}
					return res, ErrInvalidValue
				}
			} else {
				// Try to match the kind of resulting value with the argument declaration.
				rt := argFn.Returning

				if rt.FieldKind != ad.FieldKind {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, rt.FieldKind)
					}
					return res, ErrInvalidValue
				}

				if rt.EnumDescriptor != nil && rt.EnumDescriptor.FullName() != ad.EnumDescriptor.FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, rt.EnumDescriptor.FullName())
					}
					return res, ErrInvalidValue
				}

				if ad.Message() != nil && ad.IsMap() && !rt.IsMap() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d does not return a map value", x.JoinedName(), i)
					}
					return res, ErrInvalidValue
				}

				if ad.Message() != nil && ad.Message().FullName() != rt.Message().FullName() {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is not of type %s", x.JoinedName(), i, rt.Message().FullName())
					}
					return res, ErrInvalidValue
				}

				if ad.Cardinality() == protoreflect.Repeated && rt.Cardinality() != protoreflect.Repeated {
					var res TryParseValueResult
					if ctx.ErrHandler != nil {
						res.ErrPos = x.Position()
						res.ErrMsg = fmt.Sprintf("function call %s argument %d is repeated", x.JoinedName(), i)
					}
					return res, ErrInvalidValue
				}
			}

			// All checks passed, so we can call the function call handler.
			rfn, err := b.tryParseAndCallFunction(ctx, at, argFn, ad.Indirect)
			if err != nil {
				// If calling the function failed, then we return the error.
				return rfn, err
			}

			args = append(args, rfn.Expr)
			isIndirect = isIndirect || rfn.IsIndirect
			continue
		case *ast.StructExpr:
			// A struct expression is a valid argument if the argument declaration is a message or a map.
			if ad.Kind() != protoreflect.MessageKind {
				// A struct expression can only be used as an argument if the argument declaration
				// is a message or a map.
				// A map is also of a MessageKind, thus it is an invalid value error.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument is not of a message type, thus struct expression is not a valid argument", x.JoinedName())
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			if ad.IsMap() && !at.IsMap() {
				// A map argument needs to be a map expression.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument %d must be a map", x.JoinedName(), i)
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			if !ad.IsMap() && at.IsMap() {
				// A map argument needs to be a map expression.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument %d is a message not a map", x.JoinedName(), i)
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			// This should be a valid argument.
			av, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         ad,
				Value:         at,
				AllowIndirect: ad.Indirect,
			})
			if err != nil {
				clearArgs()
				return av, err
			}

			args = append(args, av.Expr)
		case *ast.ArrayExpr:
			// An array expression is a valid argument if the argument declaration is a repeated field.
			if ad.Cardinality() != protoreflect.Repeated {
				// An array expression can only be used as an argument if the argument declaration
				// is a repeated field.
				var res TryParseValueResult
				if ctx.ErrHandler != nil {
					res.ErrPos = x.Position()
					res.ErrMsg = fmt.Sprintf("function call %s argument is not repeated, thus array expression is not a valid argument", x.JoinedName())
				}
				clearArgs()
				return res, ErrInvalidValue
			}

			av, err := b.TryParseValue(ctx, TryParseValueInput{
				Field:         ad,
				Value:         at,
				AllowIndirect: ad.Indirect,
			})
			if err != nil {
				clearArgs()
				return av, err
			}

			isIndirect = isIndirect || av.IsIndirect
			args = append(args, av.Expr)
		default:
		}
	}

	// All arguments are parsed and checked.
	// We can call the function call handler.
	ex, err := fn.CallFn(args...)
	if err != nil {
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("function call %s failed: %v", x.JoinedName(), err)
		}
		clearArgs()
		return res, err
	}

	return TryParseValueResult{Expr: ex.Expr, IsIndirect: ex.IsIndirect || isIndirect}, nil
}

func (b *Interpreter) getFunctionDeclaration(ctx *ParseContext, x *ast.FunctionCall) (*FunctionCallDeclaration, bool) {
	fn, ok := b.functionCallDeclarations[x.JoinedName()]
	if !ok {
		return nil, false
	}
	return fn, true
}
