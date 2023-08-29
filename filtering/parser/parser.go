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

package parser

import (
	"errors"
	"sync"

	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

// Parser is responsible for parsing the input string filter into an AST.
type Parser struct {
	src     string
	scanner scanner.Scanner

	// err is the error handler.
	err scanner.ErrorHandler

	strictWhiteSpaces bool
}

// ParserOption changes the behavior of the parser.
type ParserOption func(p *Parser)

// StrictWhitespacesOption makes the parser to fail if there are more than one whitespace between specific expresions.
func StrictWhitespacesOption() ParserOption {
	return func(p *Parser) {
		p.strictWhiteSpaces = true
	}
}

// ErrorHandlerOption sets the error handler of the parser.
func ErrorHandlerOption(err scanner.ErrorHandler) ParserOption {
	return func(p *Parser) {
		p.err = err
	}
}

// NewParser creates a new parser with the given options.
func NewParser(src string, opts ...ParserOption) *Parser {
	p := &Parser{src: src}

	for _, opt := range opts {
		opt(p)
	}

	p.scanner.Reset(src, p.err)

	return p
}

// Reset resets the parser with the given input string.
func (p *Parser) Reset(src string, opts ...ParserOption) {
	p.src = src
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	p.scanner.Reset(src, p.err)
}

// ErrInvalidFilterSyntax is returned when the input string filter has invalid syntax.
var ErrInvalidFilterSyntax = errors.New("invalid filter")

// Parse parses the input string filter into an AST.
// If the input was an empty string, the returned ParsedFilter will have a nil Expr.
func (p *Parser) Parse() (*ParsedFilter, error) {
	pf := getParsedFilter()
	if p.src == "" {
		return pf, nil
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	pos, tok, lit := p.scanner.Scan()
	if tok != token.EOF {
		if p.err != nil {
			p.err(pos, "expr: EOF expected but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	pf.Expr = expr

	return pf, nil
}

func (p *Parser) parseSimpleExpr() (ast.SimpleExpr, error) {
	var isComposite bool
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		if tok == token.LPAREN {
			isComposite = true
		}
		return false
	})

	if isComposite {
		return p.parseCompositeExpr()
	}
	return p.parseRestrictionExpr()
}

type nameParts struct {
	parts []namePart
}

type namePart struct {
	pos token.Position
	lit string
	tok token.Token
}

var namePartsPool = sync.Pool{
	New: func() any {
		return &nameParts{
			parts: make([]namePart, 0, 10),
		}
	},
}

func getNameParts() *nameParts {
	return namePartsPool.Get().(*nameParts)
}

func putNameParts(nameParts *nameParts) {
	if nameParts == nil {
		return
	}

	nameParts.parts = nameParts.parts[:0]
	namePartsPool.Put(nameParts)
}

// ParsedFilter is a parsed filter expression.
type ParsedFilter struct {
	// Expr is a parsed filter expression, possibly nil (for empty filter).
	Expr *ast.Expr
}

// Free frees the resource associated with the parsed filter.
// This should be used in a defer statement immediately after calling Parse.
// No further use of any filter expressions is allowed after calling Free.
func (p *ParsedFilter) Free() {
	putParsedFilter(p)
}
