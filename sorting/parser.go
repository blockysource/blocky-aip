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

package sorting

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
	blockyannnotations "github.com/blockysource/go-genproto/blocky/api/annotations"
)

var (
	ErrInternalError    = errors.New("internal error")
	ErrInvalidField     = errors.New("invalid field")
	ErrInvalidSyntax    = errors.New("invalid syntax")
	ErrSortingForbidden = errors.New("sorting forbidden")
)

// Parser parses an order by expression
type Parser struct {
	msgDesc protoreflect.MessageDescriptor

	errHandler ErrHandlerFn
}

// Reset resets the parser with a new message descriptor
// and optional error handler.
// If the error handler is nil, the parser will handling errors.
func (p *Parser) Reset(msgDesc protoreflect.MessageDescriptor, errHandler ErrHandlerFn) {
	p.msgDesc = msgDesc
	p.errHandler = errHandler
}

// ErrHandlerFn is a function that handles errors
type ErrHandlerFn func(pos int, msg string)

// Parse parses a sorting order option and returns an expression.
func (p *Parser) Parse(orderBy string) (oe *expr.OrderByExpr, err error) {
	var s scanner
	s.init(orderBy)

	// Check if the input is empty.
	s.skipWhitespace()

	var tk token
	s.peekToken(func(p position, t token, l string) bool {
		tk = t
		return tk == eof_tok
	})
	if tk == eof_tok {
		if p.errHandler != nil {
			p.errHandler(0, "empty input")
		}
		return nil, ErrInvalidSyntax
	}

	// Parse the order by expression.
	oe = expr.AcquireOrderByExpr()
	for {
		// Scan next token, for the EOF or next field.
		pos, tok, lit := s.scan()
		if tok == eof_tok {
			// This means the end of the field order by expression
			break
		}

		if tok != field_tok {
			oe.Free()
			if p.errHandler != nil {
				p.errHandler(int(pos), fmt.Sprintf("expected field name but got %q", lit))
			}
			return nil, ErrInvalidSyntax
		}

		// Set up current field context.
		cur := expr.AcquireOrderByFieldExpr()

		// Parse the field name literal.
		fd, err := p.parseField(p.msgDesc, int(pos), lit)
		if err != nil {
			oe.Free()
			return nil, err
		}

		if isFieldSortingForbidden(fd) {
			oe.Free()
			if p.errHandler != nil {
				p.errHandler(int(pos), fmt.Sprintf("field %q is forbidden for sorting", fd.Name()))
			}
			return nil, ErrSortingForbidden
		}

		fe := expr.AcquireFieldSelectorExpr()
		fe.Field = fd
		fe.FieldComplexity = getFieldComplexity(fd)
		cur.Field = fe

		for {
			var tk token
			// Peek the token to see if there is traversal selector.
			s.peekToken(func(_ position, t token, _ string) bool {
				tk = t
				return tk == period_tok
			})

			// If there is no traversal selector, break the loop.
			if tk != period_tok {
				break
			}

			// We're after a traversal 'dot' selector,
			// so we need to parse the next field name.
			pos, tok, lit = s.scan()
			if tok != field_tok {
				oe.Free()
				if p.errHandler != nil {
					p.errHandler(int(pos), fmt.Sprintf("expected field name but got %s", tok))
				}
				return nil, ErrInvalidSyntax
			}

			fd, err = p.parseField(fd.Message(), int(pos), lit)
			if err != nil {
				oe.Free()
				return nil, err
			}

			if isFieldSortingForbidden(fd) {
				oe.Free()
				if p.errHandler != nil {
					p.errHandler(int(pos), fmt.Sprintf("field %q is forbidden for sorting", fd.Name()))
				}
				return nil, ErrSortingForbidden
			}

			fe := expr.AcquireFieldSelectorExpr()
			fe.Field = fd
			fe.FieldComplexity = getFieldComplexity(fd)

			setLatestTraverseField(cur, fe)
		}

		s.skipWhitespace()

		// Scan next token.
		// It may either be a comma, order or EOF.
		pos, tok, lit = s.scan()
		switch tok {
		case comma_tok:
			// This means the end of the field order by expression
			oe.Fields = append(oe.Fields, cur)

			s.skipWhitespace()
			continue
		case asc_tok:
			cur.Order = expr.ASC
		case desc_tok:
			cur.Order = expr.DESC
		case eof_tok:
			oe.Fields = append(oe.Fields, cur)
			return oe, nil
		default:
			if p.errHandler != nil {
				p.errHandler(int(pos), fmt.Sprintf("expected comma, order or EOF but got %q", lit))
			}
			oe.Free()
			return nil, ErrInvalidSyntax
		}

		s.skipWhitespace()

		// Scan next token, for the comma or EOF.
		pos, tok, lit = s.scan()
		switch tok {
		case comma_tok:
			// This means the end of the field order by expression
			oe.Fields = append(oe.Fields, cur)

			s.skipWhitespace()
			continue
		case eof_tok:
			oe.Fields = append(oe.Fields, cur)
			return oe, nil
		default:
			if p.errHandler != nil {
				p.errHandler(int(pos), fmt.Sprintf("expected comma or EOF but got %s", tok))
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

func (p *Parser) parseField(md protoreflect.MessageDescriptor, pos int, lit string) (protoreflect.FieldDescriptor, error) {
	if lit == "" {
		if p.errHandler != nil {
			p.errHandler(pos, "expected field name")
		}
		return nil, ErrInternalError
	}

	fd := md.Fields().ByName(protoreflect.Name(lit))
	if fd == nil {
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("field: %s is not a valid field", lit))
		}
		return nil, ErrInvalidField
	}

	return fd, nil
}

func getFieldComplexity(fd protoreflect.FieldDescriptor) int64 {
	c, ok := proto.GetExtension(fd.Options(), blockyannnotations.E_Complexity).(int64)
	if !ok {
		return 1
	}
	return c
}

func isFieldSortingForbidden(fd protoreflect.FieldDescriptor) bool {
	qp, ok := proto.GetExtension(fd.Options(), blockyannnotations.E_QueryOpt).([]blockyannnotations.FieldQueryOption)
	if !ok {
		return false
	}

	for _, p := range qp {
		if p == blockyannnotations.FieldQueryOption_FORBID_SORTING {
			return true
		}
	}
	return false
}
