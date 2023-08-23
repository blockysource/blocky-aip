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
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/parser"
	"github.com/blockysource/blocky-aip/filtering/token"
)

var (
	// ErrNoHandlerFound is a standard error that is returned when no handler is found for an expression.
	ErrNoHandlerFound = errors.New("no handler found")

	// ErrInvalidField is a standard error that is returned when a field is invalid.
	ErrInvalidField = errors.New("invalid field")

	// ErrInvalidValue is a standard error that is returned when a value is invalid.
	ErrInvalidValue = errors.New("invalid value")

	// ErrFieldNotFound is a standard error that is returned when a field is not found.
	ErrFieldNotFound = errors.New("field not found")

	// ErrInvalidAST is an error that is returned when the AST is invalid.
	ErrInvalidAST = errors.New("invalid AST")

	// ErrInternal is an internal error done during interpretation.
	ErrInternal = errors.New("internal error")

	// ErrAmbiguousField is an error that is returned when a field is ambiguous.
	ErrAmbiguousField = errors.New("ambiguous field selector")
)

// Interpreter is an interpreter that can parse a query string and return an expression.
type Interpreter struct {
	// msg is a message descriptor which is used to resolve field names.
	msg protoreflect.MessageDescriptor

	// error handler function used to handle errors during parsing.
	errHandlerFn func(pos token.Position, msg string)

	functionCallDeclarations map[string]*FunctionCallDeclaration

	fieldInfo struct {
		ls  []fieldInfo
		mut sync.RWMutex
	}
}

type fieldInfo struct {
	fd         protoreflect.FieldDescriptor
	complexity int64
	forbidden  bool
	nullable   bool
}

// Option is an option that can be passed to the interpreter.
type Option func(*Interpreter) error

// ErrHandlerOpt is an option that sets the error handler of the interpreter.
func ErrHandlerOpt(errorHandler func(pos token.Position, msg string)) Option {
	return func(i *Interpreter) error {
		if i.errHandlerFn != nil {
			return errors.New("error handler is already set")
		}
		i.errHandlerFn = errorHandler
		return nil
	}
}

// RegisterFunction is an Option that registers a function call declaration within the interpreter.
// Once registered, the function can be used in the filter.
func RegisterFunction(fn *FunctionCallDeclaration) Option {
	return func(i *Interpreter) error {
		if i.functionCallDeclarations == nil {
			i.functionCallDeclarations = make(map[string]*FunctionCallDeclaration)
		}

		fnFullName := fn.Name.String()
		if _, ok := i.functionCallDeclarations[fnFullName]; ok {
			return fmt.Errorf("function %q is already registered", fnFullName)
		}

		// Verify if the declaration is valid.
		if err := fn.Validate(); err != nil {
			return err
		}

		i.functionCallDeclarations[fnFullName] = fn
		return nil
	}
}

// NewInterpreter returns a new interpreter.
func NewInterpreter(msg protoreflect.MessageDescriptor, opts ...Option) (*Interpreter, error) {
	b := Interpreter{
		msg: msg,
	}

	if err := b.Reset(msg, opts...); err != nil {
		return nil, err
	}
	return &b, nil
}

func (b *Interpreter) Reset(msg protoreflect.MessageDescriptor, opts ...Option) error {
	b.msg = msg
	b.fieldInfo.ls = make([]fieldInfo, 0, 16)

	if b.msg == nil {
		return errors.New("message descriptor is not set")
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return err
		}
	}
	return nil
}

func (b *Interpreter) getFieldInfo(fd protoreflect.FieldDescriptor) fieldInfo {
	b.fieldInfo.mut.RLock()

	for _, fi := range b.fieldInfo.ls {
		if fi.fd == fd {
			b.fieldInfo.mut.RUnlock()
			return fi
		}
	}
	b.fieldInfo.mut.RUnlock()

	b.fieldInfo.mut.Lock()
	fi := fieldInfo{
		fd:         fd,
		complexity: getFieldComplexity(fd),
		forbidden:  IsFieldFilteringForbidden(fd),
		nullable:   IsFieldNullable(fd),
	}
	b.fieldInfo.ls = append(b.fieldInfo.ls, fi)
	b.fieldInfo.mut.Unlock()
	return fi
}

// Parse input filter into an expression.
// Implements filtering.Interpreter interface.
// By default, interpreter is returning a non-precise error if the parsing fails.
// For detailed error handling, provide an error handler function during initialization of the interpreter.
func (b *Interpreter) Parse(filter string) (expr.FilterExpr, error) {
	var p parser.Parser

	if b.msg == nil {
		panic("message descriptor is not set")
	}

	if filter == "" {
		return nil, nil
	}

	var errHandler parser.ParserOption
	if b.errHandlerFn != nil {
		errHandler = parser.ErrorHandlerOption(b.errHandlerFn)
	}

	p.Reset(filter,
		errHandler,
		parser.UseArraysOption,
		parser.UseInComparatorOption,
		parser.UseStructsOption)

	pf, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if pf.Expr == nil {
		return nil, status.Error(codes.Internal, "parsing filter failed")
	}
	defer pf.Free()

	ctx := contextPool.Get().(*ParseContext)
	defer ctx.Free()

	ctx.Message = b.msg
	ctx.ErrHandler = b.errHandlerFn
	ctx.Interpreter = b

	he, err := b.HandleExpr(ctx, pf.Expr)
	if err != nil {
		if b.errHandlerFn != nil {
			b.errHandlerFn(he.ErrPos, he.ErrMsg)
		}
		return nil, err
	}
	return he.Expr, nil
}

// HandledExpr is a struct that contains an expression and a flag that indicates if the expression was consumed.
type HandledExpr struct {
	Expr     expr.FilterExpr
	Consumed bool
}

var contextPool = sync.Pool{
	New: func() any {
		return &ParseContext{}
	},
}

// ParseContext is a struct that contains the context of the expression.
type ParseContext struct {
	// Message is the message descriptor of the current expression.
	Message protoreflect.MessageDescriptor

	// ErrHandler is the error handler function.
	ErrHandler func(pos token.Position, msg string)

	// Interpreter is the reference to the interpreter that parses the expression.
	// It can be used by custom handlers to reuse standard handlers for sub-expressions.
	Interpreter *Interpreter

	isAcquired bool
}

// Free frees the context.
func (c *ParseContext) Free() {
	if c == nil {
		return
	}
	if !c.isAcquired {
		return
	}

	c.Message = nil
	c.ErrHandler = nil
	c.Interpreter = nil
}
