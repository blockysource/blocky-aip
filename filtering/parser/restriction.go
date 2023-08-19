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

func (p *Parser) parseRestrictionExpr() (*ast.RestrictionExpr, error) {
	re := getRestrictionExpr()

	// Parse comparable expression.
	comp, err := p.parseComparableExpr()
	if err != nil {
		return nil, err
	}

	bp := p.scanner.Breakpoint()
	re.Pos = comp.Position()
	re.Comparable = comp

	// Skip possible whitespaces.
	n := p.scanner.SkipWhitespace()

	// Peek if there is a comparator.
	var (
		isComparator bool
		eof          bool
	)
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		isComparator = tok.IsComparator()
		eof = tok == token.EOF
		return false
	})

	if !isComparator || eof {
		p.scanner.Restore(bp)
		return re, nil
	}

	if p.strictWhiteSpaces && n == 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: whitespace expected after comparable expression")
		}
		return nil, ErrInvalidFilterSyntax
	}

	if p.strictWhiteSpaces && n > 1 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: only one whitespace is allowed between comparable expression and comparator")
		}
		return nil, ErrInvalidFilterSyntax
	}

	compOp, err := p.parseComparator()
	if err != nil {
		return nil, err
	}
	re.Comparator = compOp

	// Skip possible whitespaces.
	n = p.scanner.SkipWhitespace()

	if p.strictWhiteSpaces && n == 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: whitespace expected after comparator")
		}
		return nil, ErrInvalidFilterSyntax
	}

	if p.strictWhiteSpaces && n > 1 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: only one whitespace is allowed between comparator and argument")
		}
		return nil, ErrInvalidFilterSyntax
	}

	// Parse the argument.
	arg, err := p.parseArgExpr()
	if err != nil {
		return nil, err
	}
	re.Arg = arg

	return re, nil
}
