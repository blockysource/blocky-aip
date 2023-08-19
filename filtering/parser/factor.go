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

func (p *Parser) parseFactorExpr() (*ast.FactorExpr, error) {
	factor := getFactorExpr()

	// Parse the first term.
	term, err := p.parseTermExpr()
	if err != nil {
		return nil, err
	}

	factor.Terms = append(factor.Terms, term)

	for {
		bp := p.scanner.Breakpoint()

		// Skip possible whitespaces.
		n := p.scanner.SkipWhitespace()
		if n == 0 {
			return factor, nil
		}

		if p.strictWhiteSpaces && n > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "factor: only one WS is allowed between term and OR operator")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// Parse the NOT operator.
		var (
			isOR  bool
			orPos token.Position
		)
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			orPos = pos
			if tok == token.OR {
				isOR = true
				return true
			}
			return false
		})
		if !isOR {
			// Restore the break point as we've consumed the whitespaces.
			p.scanner.Restore(bp)
			// The whitespace is not followed by the NOT operator.
			// The factor is a single term.
			// The whitespace was consumed by the peek, but it doesn't chane the syntax.
			return factor, nil
		}

		// The OR operator is found.
		// We set its position in a previous TermExpr.
		term.OrOpPos = orPos

		n = p.scanner.SkipWhitespace()
		if n == 0 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "factor: WS expected after OR operator")
			}
			return nil, ErrInvalidFilterSyntax
		}
		if p.strictWhiteSpaces && n > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "factor: only one WS is allowed between OR operator and term")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// Parse the next term.
		term, err = p.parseTermExpr()
		if err != nil {
			return nil, err
		}

		factor.Terms = append(factor.Terms, term)
	}
}
