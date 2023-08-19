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

func (p *Parser) parseExpr() (*ast.Expr, error) {
	// Expression is a single or 'AND' separated sequences.
	// Parse the first sequence.
	p.scanner.SkipWhitespace()

	expr := getExpr()

	// Peek for the whitespaces after the first sequence.
	for {
		// Parse the sequence.
		seq, err := p.parseSequenceExpr()
		if err != nil {
			return nil, err
		}
		expr.Sequences = append(expr.Sequences, seq)

		// Skip possible whitespaces.
		n := p.scanner.SkipWhitespace()
		if n == 0 {
			return expr, nil
		}
		if p.strictWhiteSpaces && n > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "expr: only one WS is allowed between sequence and AND operator")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// Parse the AND operator.
		var (
			andT   token.Token
			andPos token.Position
		)
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			andT = tok
			andPos = pos
			if tok == token.AND {
				return true
			}
			return false
		})

		switch andT {
		case token.AND:
			// Set the position of the AND operator.
			seq.OpPos = andPos

			n = p.scanner.SkipWhitespace()
			if n == 0 {
				if p.err != nil {
					p.err(p.scanner.Pos(), "expr: WS expected after AND operator")
				}
				return nil, ErrInvalidFilterSyntax
			}
			if p.strictWhiteSpaces && n > 1 {
				if p.err != nil {
					p.err(p.scanner.Pos(), "expr: only one WS is allowed between AND operator and sequence")
				}
				return nil, ErrInvalidFilterSyntax
			}
		case token.EOF:
			return expr, nil
		default:
			if p.err != nil {
				p.err(p.scanner.Pos(), "expr: AND operator expected but got: "+andT.String())
			}
			return nil, ErrInvalidFilterSyntax
		}
	}
}
