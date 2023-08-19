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
	"github.com/blockysource/blocky-aip/filtering/token"
)

var arrayExprPool = &sync.Pool{
	New: func() any {
		return &ast.ArrayExpr{
			Elements: make([]ast.ComparableExpr, 0, 10),
		}
	},
}

func getArrayExpr() *ast.ArrayExpr {
	return arrayExprPool.Get().(*ast.ArrayExpr)
}

func putArrayExpr(a *ast.ArrayExpr) {
	if a == nil {
		return
	}

	for i := range a.Elements {
		putComparableExpr(a.Elements[i])
	}
	a.Elements = a.Elements[:0]
	arrayExprPool.Put(a)
}

func (p *Parser) parseArrayExpr(pos token.Position) (*ast.ArrayExpr, error) {
	a := getArrayExpr()
	a.LBracket = pos

	i := 0
	for {
		p.scanner.SkipWhitespace()

		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) (consume bool) {
			pt = tok
			if tok == token.COMMA && i > 0 {
				return true
			}
			return false
		})

		if (i > 0 && pt != token.COMMA) || (i == 0 && pt == token.BRACKET_CLOSE) {
			break
		}
		p.scanner.SkipWhitespace()

		c, err := p.parseComparableExpr()
		if err != nil {
			putArrayExpr(a)
			return nil, err
		}
		a.Elements = append(a.Elements, c)

		i++
	}

	pos, tok, lit := p.scanner.Scan()
	if tok != token.BRACKET_CLOSE {
		if p.err != nil {
			p.err(pos, "array: ']' expected but got: "+lit)
		}
		putArrayExpr(a)
		return nil, ErrInvalidFilterSyntax
	}

	a.RBracket = pos
	return a, nil
}
