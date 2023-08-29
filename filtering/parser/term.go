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
	"github.com/blockysource/blocky-aip/token"
)

func (p *Parser) parseTermExpr() (*ast.TermExpr, error) {
	// The term is either unary or simple term.
	// Check if the next token is unary operator.
	var (
		pos token.Position
		tok token.Token
		lit string
	)

	bp := p.scanner.Breakpoint()
	p.scanner.Peek(func(p token.Position, t token.Token, l string) bool {
		pos, tok, lit = p, t, l
		return tok.IsUnaryOperator()
	})

	switch tok {
	case token.NOT:
		// If the token is NOT, it is probably a unary operator.
		// However, we need to check edge cases when the message field is name.
		nm, err := p.isKeywordMember()
		if err != nil {
			return nil, err
		}
		p.scanner.Restore(bp)

		if nm {
			// NOT is a member.
			te := getTermExpr()
			simple, err := p.parseSimpleExpr()
			if err != nil {
				putTermExpr(te)
				return nil, err
			}
			te.Pos = simple.Position()
			te.Expr = simple
			return te, nil
		}
		p.scanner.Scan() // Scan the NOT token.

		n := p.scanner.SkipWhitespace()
		if p.strictWhiteSpaces && n > 1 {
			if p.err != nil {
				p.err(pos, "term: invalid syntax")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// NOT is a unary operator.
		te := getTermExpr()
		te.Pos = pos
		te.UnaryOp = lit
		simple, err := p.parseSimpleExpr()
		if err != nil {
			putTermExpr(te)
			return nil, err
		}
		te.Expr = simple
		return te, nil
	case token.MINUS:
		// Check if this is a negative number or a unary operator.
		n := p.scanner.SkipWhitespace()
		if n > 0 {
			if p.err != nil {
				p.err(pos, "whitespace after '-' unary operator is not allowed")
			}
			return nil, ErrInvalidFilterSyntax
		}

		te := getTermExpr()
		te.Pos = pos
		te.UnaryOp = lit
		simple, err := p.parseSimpleExpr()
		if err != nil {
			putTermExpr(te)
			return nil, err
		}
		te.Expr = simple
		return te, nil
	default:
		// Not a unary operator.
		// Restore the breakpoint and parse a simple term.
		te := getTermExpr()
		simple, err := p.parseSimpleExpr()
		if err != nil {
			return nil, err
		}
		te.Pos = simple.Position()
		te.Expr = simple
		return te, nil
	}
}

func (p *Parser) isKeywordMember() (bool, error) {
	n := p.scanner.SkipWhitespace()
	if n == 0 {
		pos, tok, _ := p.scanner.Scan()
		switch tok {
		case token.PERIOD, token.LPAREN, token.BRACE_OPEN, token.BRACKET_CLOSE, token.BRACE_CLOSE, token.COMMA, token.EOF:
			return true, nil
		default:
			if tok.IsComparator() {
				return true, nil
			}
			// Invalid syntax.
			if p.err != nil {
				p.err(pos, "term: invalid syntax")
			}
			return false, ErrInvalidFilterSyntax
		}
	}

	// If there is more than zero whitespace, then check the next token.
	pos, tok, _ := p.scanner.Scan()

	// If the token is a comparator or the next token is LPAREN, then it is a member term.
	if tok.IsComparator() {
		return true, nil
	}

	switch tok {
	case token.LPAREN:
		return false, nil
	case token.RPAREN, token.BRACKET_CLOSE, token.BRACE_CLOSE, token.COMMA:
		return true, nil
	case token.AND, token.OR, token.NOT:
		is, err := p.isKeywordMember()
		if err != nil {
			return false, err
		}
		return !is, nil
	case token.EOF:
		// EOF means that the NOT is a unary operator.
		return true, nil
	}

	// If the next token is an identifier
	if tok.IsIdent() {
		return false, nil
	}

	// Invalid syntax.
	if p.err != nil {
		p.err(pos, "term: invalid syntax")
	}
	return false, ErrInvalidFilterSyntax
}
