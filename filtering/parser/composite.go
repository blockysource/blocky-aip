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

func (p *Parser) parseCompositeExpr() (*ast.CompositeExpr, error) {
	pos, tok, lit := p.scanner.Scan()
	if tok != token.LPAREN {
		if p.err != nil {
			p.err(pos, "composite: LPAREN expected at the beginning of composite expression but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	cl := getCompositeExpr()
	cl.Lparen = pos

	// Skip possible whitespaces.
	_ = p.scanner.SkipWhitespace()

	// Parse the expression element.
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	cl.Expr = expr

	// Skip possible whitespaces.
	_ = p.scanner.SkipWhitespace()

	pos, tok, lit = p.scanner.Scan()
	if tok != token.RPAREN {
		if p.err != nil {
			p.err(pos, "composite: RPAREN expected at the end of composite expression but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}
	cl.Rparen = pos

	return cl, nil
}
