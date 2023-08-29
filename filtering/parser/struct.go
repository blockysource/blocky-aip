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
	"sync"

	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/token"
)

var (
	structExprPool = sync.Pool{
		New: func() any {
			return &ast.StructExpr{
				Name:     make([]ast.NameExpr, 0, 10),
				Elements: make([]*ast.StructFieldExpr, 0, 10),
			}
		},
	}
	structFieldExprPool = sync.Pool{
		New: func() any {
			return &ast.StructFieldExpr{
				Name: make([]ast.ValueExpr, 0, 10),
			}
		},
	}
)

func getStructExpr() *ast.StructExpr {
	return structExprPool.Get().(*ast.StructExpr)
}

func putStructExpr(expr *ast.StructExpr) {
	if expr == nil {
		return
	}

	for _, name := range expr.Name {
		putNameExpr(name)
	}
	expr.Name = expr.Name[:0]

	for _, field := range expr.Elements {
		putStructFieldExpr(field)
	}
	expr.Elements = expr.Elements[:0]
	expr.LBrace = 0
	expr.RBrace = 0
	structExprPool.Put(expr)
}

func getStructFieldExpr() *ast.StructFieldExpr {
	return structFieldExprPool.Get().(*ast.StructFieldExpr)
}

func putStructFieldExpr(expr *ast.StructFieldExpr) {
	if expr == nil {
		return
	}

	for _, name := range expr.Name {
		putValueExpr(name)
	}
	expr.Name = expr.Name[:0]
	putComparableExpr(expr.Value)
	structFieldExprPool.Put(expr)
}

func (p *Parser) parseStructExpr(nameParts *nameParts) (*ast.StructExpr, error) {
	st := getStructExpr()

	if nameParts != nil {
		defer putNameParts(nameParts)
	}

	p.scanner.SkipWhitespace()

	if nameParts != nil {
		for _, np := range nameParts.parts {
			switch {
			case np.tok == token.IDENT, np.tok.IsKeyword():
				// NOTE: struct name doesn't support token.TIMESTAMP.
				text := getTextLiteral()
				text.Pos = np.pos
				text.Value = np.lit
				text.Token = np.tok
				st.Name = append(st.Name, text)
			default:
				if p.err != nil {
					p.err(np.pos, "struct: TEXT expected but got: "+np.lit)
				}
				return nil, ErrInvalidFilterSyntax
			}
		}
	}

	pos, tok, lit := p.scanner.Scan()
	if tok != token.BRACE_OPEN {
		if p.err != nil {
			p.err(pos, "struct: '{' expected but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	st.LBrace = pos

	var isRBrace bool
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		isRBrace = tok == token.BRACE_CLOSE
		st.RBrace = pos
		return isRBrace
	})

	if isRBrace {
		return st, nil
	}

	i := 0
	for {
		p.scanner.SkipWhitespace()

		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			pt = tok
			if i > 0 && tok == token.COMMA {
				return true
			}
			return false
		})
		if (i > 0 && pt != token.COMMA) || (i == 0 && pt == token.BRACE_CLOSE) {
			break
		}

		p.scanner.SkipWhitespace()

		if pt == token.COMMA {
			// Peek if the next token is a brace close.
			// This is a Golang style to finish each field with a comma.
			p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
				pt = tok
				return false
			})
			if pt == token.BRACE_CLOSE {
				break
			}
		}

		sf, err := p.parseStructFieldExpr()
		if err != nil {
			putStructExpr(st)
			return nil, err
		}

		st.Elements = append(st.Elements, sf)
		i++
	}

	pos, tok, lit = p.scanner.Scan()
	if tok != token.BRACE_CLOSE {
		if p.err != nil {
			p.err(pos, "struct: '}' expected but got: "+lit)
		}
		putStructExpr(st)
		return nil, ErrInvalidFilterSyntax
	}
	st.RBrace = pos
	return st, nil
}

func (p *Parser) parseStructFieldExpr() (*ast.StructFieldExpr, error) {
	sf := getStructFieldExpr()
	pos, tok, lit := p.scanner.Scan()
	switch {
	case tok == token.STRING:
		sl := getStringLiteral()
		sl.Pos = pos
		sl.Value = lit
		sf.Name = append(sf.Name, sl)
	case tok.IsNonStringLit() || tok.IsKeyword():
		text := getTextLiteral()
		text.Pos = pos
		text.Value = lit
		text.Token = tok
		sf.Name = append(sf.Name, text)
	default:
		if p.err != nil {
			p.err(pos, "struct: TEXT, STRING or KEYWORD expected but got: '"+lit+"'")
		}
		putStructFieldExpr(sf)
		return nil, ErrInvalidFilterSyntax
	}

	var i int
	for {
		if i > 0 {
			pos, tok, lit = p.scanner.Scan()
			switch {
			case tok == token.STRING:
				sl := getStringLiteral()
				sl.Pos = pos
				sl.Value = lit
				sf.Name = append(sf.Name, sl)
			case tok.IsNonStringLit() || tok.IsKeyword():
				text := getTextLiteral()
				text.Pos = pos
				text.Value = lit
				text.Token = tok
				sf.Name = append(sf.Name, text)
			default:

				if p.err != nil {
					p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
				}
				putStructFieldExpr(sf)
				return nil, ErrInvalidFilterSyntax

			}
		}

		p.scanner.SkipWhitespace()

		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			// Expects a dot
			pt = tok
			return tok == token.PERIOD
		})
		if pt != token.PERIOD {
			break
		}
	}

	p.scanner.SkipWhitespace()

	pos, tok, lit = p.scanner.Scan()
	if tok != token.COLON {
		if p.err != nil {
			p.err(pos, "struct: ':' expected but got: "+lit)
		}
		putStructFieldExpr(sf)
		return nil, ErrInvalidFilterSyntax
	}
	sf.Colon = pos

	p.scanner.SkipWhitespace()

	value, err := p.parseComparableExpr()
	if err != nil {
		putStructFieldExpr(sf)
		return nil, err
	}

	sf.Value = value

	return sf, nil
}
