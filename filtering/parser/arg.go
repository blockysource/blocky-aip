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

func (p *Parser) parseArgExpr() (ast.ArgExpr, error) {
	// Peek for the composite LPAREN token.
	var isComposite bool
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		isComposite = tok == token.LPAREN
		return false
	})

	if isComposite {
		// Parse the composite expression.
		return p.parseCompositeExpr()
	}
	return p.parseComparableExpr()
}
