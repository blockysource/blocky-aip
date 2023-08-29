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
		if tok.IsLiteral() || tok.IsKeyword() {
			return true
		}
		return false
	})

	switch {
	case tok.IsLiteral(), tok.IsKeyword():
	case tok == token.BRACE_OPEN:
		return p.parseStructExpr(nil)
	case tok == token.BRACKET_OPEN:
		// This is returned from scanner only when arrays are enabled.
		return p.parseArrayExpr(pos)
	default:
		if p.err != nil {
			p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
		}
		return nil, ErrInvalidFilterSyntax
	}
	np := getNameParts()
	np.parts = append(np.parts, namePart{
		pos: pos,
		lit: lit,
		tok: tok,
	})

	var i int
	for {
		if i > 0 {
			pos, tok, lit = p.scanner.Scan()
			switch {
			case tok == token.STRING:
			case tok.IsNonStringLit() || tok.IsKeyword():
			default:
				if !tok.IsKeyword() {
					if p.err != nil {
						p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
					}
					putNameParts(np)
					return nil, ErrInvalidFilterSyntax
				}
			}
			np.parts = append(np.parts, namePart{
				pos: pos,
				lit: lit,
				tok: tok,
			})
		}
		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			// Expects a dot
			pt = tok
			return tok == token.PERIOD
		})

		switch pt {
		case token.BRACE_OPEN:
			return p.parseStructExpr(np)
		case token.PERIOD:
			i++
		case token.LPAREN:
			// This is a function call.
			return p.parseFuncCall(np)
		default:
			// This is the end of the member expression.
			return p.parseMemberExpr(np)
		}
	}
}
