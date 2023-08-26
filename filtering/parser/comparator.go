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

func (p *Parser) parseComparator() (*ast.ComparatorLiteral, error) {
	// Parse the restriction operator.
	pos, tok, lit := p.scanner.Scan()
	switch {
	case tok.IsComparator():
	default:
		if p.err != nil {
			p.err(pos, "restriction: comparator expected but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	cl := getComparatorLiteral()
	cl.Pos = pos
	switch tok {
	case token.EQUAL:
		cl.Type = ast.EQ
	case token.NEQ:
		cl.Type = ast.NE
	case token.GT:
		cl.Type = ast.GT
	case token.GEQ:
		cl.Type = ast.GE
	case token.LT:
		cl.Type = ast.LT
	case token.LEQ:
		cl.Type = ast.LE
	case token.COLON:
		cl.Type = ast.HAS
	case token.IN:
		cl.Type = ast.IN
	default:
		if p.err != nil {
			p.err(pos, "restriction: unknown comparator: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	return cl, nil
}
