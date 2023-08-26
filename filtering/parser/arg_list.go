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

func (p *Parser) parseArgListExpr() (*ast.ArgListExpr, error) {
	argList := getArgListExpr()

	i := 0
	for {
		// Skip possible whitespaces.
		n := p.scanner.SkipWhitespace()
		if p.strictWhiteSpaces && n > 0 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "argList: no whitespace is allowed before, between or after arguments")
			}
			return nil, ErrInvalidFilterSyntax
		}

		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			pt = tok
			if tok == token.COMMA && i > 0 {
				return true
			}
			return false
		})
		if (i > 0 && pt != token.COMMA) || (i == 0 && pt == token.RPAREN) {
			return argList, nil
		}

		// Skip possible whitespaces.
		n = p.scanner.SkipWhitespace()
		if p.strictWhiteSpaces && n > 0 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "argList: no whitespace is allowed before, between or after arguments")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// Parse the argument.
		arg, err := p.parseArgExpr()
		if err != nil {
			return nil, err
		}
		argList.Args = append(argList.Args, arg)
		i++
	}
}
