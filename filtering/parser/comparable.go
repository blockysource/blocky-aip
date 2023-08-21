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

func putComparableExpr(e ast.ComparableExpr) {
	if e == nil {
		return
	}
	switch vt := e.(type) {
	case *ast.MemberExpr:
		putMemberLiteral(vt)
	case *ast.FunctionCall:
		putFunctionLiteral(vt)
	case *ast.StructExpr:
		putStructExpr(vt)
	case *ast.ArrayExpr:
		putArrayExpr(vt)
	}
}

func (p *Parser) parseComparableExpr() (ast.ComparableExpr, error) {
	var (
		pos token.Position
		tok token.Token
		lit string
	)
	p.scanner.Peek(func(p token.Position, t token.Token, l string) bool {
		pos, tok, lit = p, t, l
		switch tok {
		case token.STRING, token.TEXT, token.TIMESTAMP, token.MINUS,
			token.NOT, token.AND, token.OR, token.IN: // (keywords)
			return true
		}
		return false
	})

	switch tok {
	case token.STRING, token.TEXT, token.TIMESTAMP, token.MINUS:
	case token.BRACE_OPEN:
		return p.parseStructExpr(nil)
	case token.BRACKET_OPEN:
		// This is returned from scanner only when arrays are enabled.
		return p.parseArrayExpr(pos)
	default:
		if !tok.IsKeyword() {
			if p.err != nil {
				p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
			}
			return nil, ErrInvalidFilterSyntax
		}
	}
	nameParts := getNameParts()
	nameParts = append(nameParts, namePart{
		pos: pos,
		lit: lit,
		tok: tok,
	})

	wasMinus := tok == token.MINUS
	var i int
	for {
		if i > 0 {
			pos, tok, lit = p.scanner.Scan()
			switch tok {
			case token.TEXT, token.TIMESTAMP:
				if wasMinus {
					wasMinus = false
				}
			case token.STRING:
				if wasMinus {
					// Invalid syntax.
					if p.err != nil {
						p.err(pos, "comparable: a '-' character is not allowed before a STRING")
					}
					return nil, ErrInvalidFilterSyntax
				}
			case token.MINUS:
				wasMinus = true
			default:
				if !tok.IsKeyword() {
					if p.err != nil {
						p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
					}
					putNameParts(nameParts)
					return nil, ErrInvalidFilterSyntax
				}
			}
			nameParts = append(nameParts, namePart{
				pos: pos,
				lit: lit,
				tok: tok,
			})
		}
		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			// Expects a dot
			pt = tok
			if tok == token.PERIOD {
				return true
			}
			return false
		})

		if p.useStructs && pt == token.BRACE_OPEN {
			// This is a struct literal
			return p.parseStructExpr(nameParts)
		}

		switch {
		case pt == token.PERIOD || wasMinus:
			i++
		case pt == token.LPAREN:
			// This is a function call.
			return p.parseFuncCall(nameParts)
		default:
			// This is the end of the member expression.
			return p.parseMemberExpr(nameParts)
		}
	}
}
