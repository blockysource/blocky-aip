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
	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/filtering/token"
)

func (p *Parser) parseTermExpr() (*ast.TermExpr, error) {
	// The term is either unary or simple term.
	// Check if the next token is unary operator.
	var (
		isUnary  bool
		unaryPos token.Position
		unaryTok token.Token
		unaryLit string
	)

	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		switch tok {
		case token.NOT, token.MINUS:
			isUnary = true
			unaryPos = pos
			unaryTok = tok
			unaryLit = lit
			return true
		}
		return false
	})
	te := getTermExpr()
	if !isUnary {
		// The term is a simple term.
		simple, err := p.parseSimpleExpr()
		if err != nil {
			return nil, err
		}
		te.Pos = simple.Position()
		te.Expr = simple
		return te, nil
	}
	te.Pos = unaryPos
	te.UnaryOp = unaryLit
	if unaryTok == token.NOT {
		ws := p.scanner.SkipWhitespace()
		if ws == 0 {
			if p.err != nil {
				p.err(unaryPos, "term: WS expected after NOT operator")
			}
			return nil, ErrInvalidFilterSyntax
		}

		if p.strictWhiteSpaces && ws > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "term: only one WS is allowed between NOT operator and simple term")
			}
			return nil, ErrInvalidFilterSyntax
		}
	}
	// MINUS
	simple, err := p.parseSimpleExpr()
	if err != nil {
		return nil, err
	}
	te.Expr = simple
	return te, nil
}
