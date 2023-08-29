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

package ordering

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/info"
	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

var (
	// ErrInternalError is an internal sorting error.
	ErrInternalError = errors.New("internal error")

	// ErrInvalidField is an error returned by the parser when the sorting expression
	// contains invalid field name.
	ErrInvalidField = errors.New("invalid field")

	// ErrInvalidSyntax is an error returned by the parser when the sorting expression
	// has invalid syntax.
	ErrInvalidSyntax = errors.New("invalid syntax")

	// ErrSortingForbidden is an error returned by the parser when the sorting
	// is forbidden of specific field.
	ErrSortingForbidden = errors.New("sorting forbidden")
)

// Parser parses an order by expression
type Parser struct {
	msgDesc    protoreflect.MessageDescriptor
	errHandler scanner.ErrorHandler

	msgInfo info.MessagesInfo
}

// ParserOpt is an option function for the parser.
type ParserOpt func(p *Parser) error

// ErrHandler sets the error handler for the parser.
func ErrHandler(fn scanner.ErrorHandler) ParserOpt {
	return func(p *Parser) error {
		p.errHandler = fn
		return nil
	}
}

// NewParser creates a new parser with a message descriptor and optional error handler.
func NewParser(msg protoreflect.MessageDescriptor, opts ...ParserOpt) (*Parser, error) {
	p := &Parser{msgDesc: msg}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	p.msgInfo = info.MapMsgInfo(msg)

	return p, nil
}

// Reset resets the parser with a new message descriptor
// and optional error handler.
// If the error handler is nil, the parser will handling errors.
func (p *Parser) Reset(msgDesc protoreflect.MessageDescriptor, opts ...ParserOpt) error {
	p.msgDesc = msgDesc
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return err
		}
	}

	p.msgInfo = info.MapMsgInfo(msgDesc)

	return nil
}

// ErrHandlerFn is a function that handles errors
type ErrHandlerFn func(pos token.Position, msg string)

// Parse parses a sorting order option and returns an expression.
func (p *Parser) Parse(orderBy string) (oe *expr.OrderByExpr, err error) {
	var s scanner.Scanner
	s.Reset(orderBy, p.errHandler)

	// Check if the input is empty.
	s.SkipWhitespace()

	var tk token.Token
	s.Peek(func(_ token.Position, t token.Token, _ string) bool {
		tk = t
		return tk == token.EOF
	})
	if tk == token.EOF {
		if p.errHandler != nil {
			p.errHandler(0, "empty input")
		}
		return nil, ErrInvalidSyntax
	}

	// Parse the order by expression.
	oe = expr.AcquireOrderByExpr()
	for {
		// Scan next token, for the EOF or next field.
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			// This means the end of the field order by expression
			break
		}

		if !tok.IsIdent() {
			oe.Free()
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected field name but got %q", lit))
			}
			return nil, ErrInvalidSyntax
		}

		// Set up current field context.
		cur := expr.AcquireOrderByFieldExpr()

		// Parse the field name literal.
		fd, c, err := p.parseField(p.msgDesc, pos, lit)
		if err != nil {
			oe.Free()
			return nil, err
		}

		fi := p.msgInfo.GetFieldInfo(fd)
		if fi.OrderingForbidden {
			if p.errHandler != nil {
				p.errHandler(pos, "ordering by given field is forbidden")
			}
			oe.Free()
			return nil, ErrSortingForbidden
		}

		fe := expr.AcquireFieldSelectorExpr()
		fe.Field = fd.Name()
		fe.FieldComplexity = c
		cur.Field = fe

		for {
			var tk token.Token
			// Peek the token to see if there is traversal selector.
			s.Peek(func(_ token.Position, t token.Token, _ string) bool {
				tk = t
				return tk == token.PERIOD
			})

			// If there is no traversal selector, break the loop.
			if tk != token.PERIOD {
				break
			}

			// We're after a traversal 'dot' selector,
			// so we need to parse the next field name.
			pos, tok, lit = s.Scan()
			if !tok.IsIdent() {
				oe.Free()
				if p.errHandler != nil {
					p.errHandler(pos, fmt.Sprintf("expected field name but got %s", tok))
				}
				return nil, ErrInvalidSyntax
			}

			fd, c, err = p.parseField(fd.Message(), pos, lit)
			if err != nil {
				oe.Free()
				return nil, err
			}

			fe := expr.AcquireFieldSelectorExpr()
			fe.Field = fd.Name()
			fe.FieldComplexity = c

			setLatestTraverseField(cur, fe)
		}

		s.SkipWhitespace()

		// Scan next token.
		// It may either be a comma, order or EOF.
		pos, tok, lit = s.Scan()
		switch tok {
		case token.COMMA:
			// This means the end of the field order by expression
			oe.Fields = append(oe.Fields, cur)

			s.SkipWhitespace()
			continue
		case token.ASC:
			cur.Order = expr.ASC
		case token.DESC:
			cur.Order = expr.DESC
		case token.EOF:
			oe.Fields = append(oe.Fields, cur)
			return oe, nil
		default:
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected comma, order or EOF but got %q", lit))
			}
			oe.Free()
			return nil, ErrInvalidSyntax
		}

		s.SkipWhitespace()

		// Scan next token, for the comma or EOF.
		pos, tok, _ = s.Scan()
		switch tok {
		case token.COMMA:
			// This means the end of the field order by expression
			oe.Fields = append(oe.Fields, cur)

			s.SkipWhitespace()
			continue
		case token.EOF:
			oe.Fields = append(oe.Fields, cur)
			return oe, nil
		default:
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected comma or EOF but got %s", tok))
			}
			oe.Free()
			return nil, ErrInvalidSyntax
		}
	}

	return oe, nil
}

func setLatestTraverseField(obfe *expr.OrderByFieldExpr, fs *expr.FieldSelectorExpr) {
	if obfe.Field == nil {
		obfe.Field = fs
		return
	}

	prev := obfe.Field
	for {
		switch pt := prev.Traversal.(type) {
		case *expr.FieldSelectorExpr:
			prev = pt
		case nil:
			prev.Traversal = fs
			return
		}
	}
}

func (p *Parser) parseField(md protoreflect.MessageDescriptor, pos token.Position, lit string) (protoreflect.FieldDescriptor, int64, error) {
	if lit == "" {
		if p.errHandler != nil {
			p.errHandler(pos, "expected field name")
		}
		return nil, 0, ErrInternalError
	}

	fd := md.Fields().ByName(protoreflect.Name(lit))
	if fd == nil {
		var found bool
		for i := 0; i < md.Oneofs().Len(); i++ {
			of := md.Oneofs().Get(i)
			fd = of.Fields().ByName(protoreflect.Name(lit))
			if fd != nil {
				found = true
				break
			}
		}
		if !found {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("field: %s is not a valid field", lit))
			}
			return nil, 0, ErrInvalidField
		}
	}

	fi := p.msgInfo.GetFieldInfo(fd)

	if fi.FilteringForbidden {
		if p.errHandler != nil {
			p.errHandler(pos, "ordering by given field is forbidden")
		}
		return nil, 0, ErrSortingForbidden
	}

	return fd, fi.Complexity, nil
}
