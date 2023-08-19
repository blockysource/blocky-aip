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

func (p *Parser) parseSequenceExpr() (*ast.SequenceExpr, error) {
	seq := getSequenceExpr()

	// Parse the first factor.
	factor, err := p.parseFactorExpr()
	if err != nil {
		return nil, err
	}
	seq.Factors = append(seq.Factors, factor)

	for {
		bp := p.scanner.Breakpoint()
		n := p.scanner.SkipWhitespace()
		// If the whitespace is not found the sequence is a single factor.
		if n == 0 {
			p.scanner.Restore(bp)
			return seq, nil
		}
		if p.strictWhiteSpaces && n > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "sequence: only one WS is allowed between factors")
			}
			return nil, ErrInvalidFilterSyntax
		}
		var isAND bool
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			if tok == token.AND {
				isAND = true
			}
			return false
		})
		if isAND {
			// Restore the break point as we've consumed the whitespaces.
			p.scanner.Restore(bp)
			// The whitespace is not followed by the AND operator.
			// The sequence is a single factor.
			// The whitespace was consumed by the peek, but it doesn't chane the syntax.
			return seq, nil
		}

		// Parse the next factor.
		factor, err = p.parseFactorExpr()
		if err != nil {
			return nil, err
		}

		seq.Factors = append(seq.Factors, factor)
	}
}
