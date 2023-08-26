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

func (p *Parser) parseMemberExpr(nameParts []namePart) (*ast.MemberExpr, error) {
	member := getMemberExpr()
	defer putNameParts(nameParts)

	for i, np := range nameParts {
		if i == 0 {
			switch {
			case np.tok == token.STRING:
				sl := getStringLiteral()
				sl.Pos = np.pos
				sl.Value = np.lit
				member.Value = sl
			case np.tok.IsNonStringLit(), np.tok.IsKeyword():
				text := getTextLiteral()
				text.Pos = np.pos
				text.Value = np.lit
				text.Token = np.tok
				member.Value = text
			default:
				if p.err != nil {
					p.err(np.pos, "comparable: TEXT or STRING expected on first element but got: "+np.lit)
				}
				return nil, ErrInvalidFilterSyntax
			}
			continue
		}

		// Others are fields (accepts a Keyword as well).
		var fieldExpr ast.FieldExpr
		switch {
		case np.tok == token.STRING:
			sl := getStringLiteral()
			sl.Pos = np.pos
			sl.Value = np.lit
			fieldExpr = sl
		case np.tok.IsNonStringLit(), np.tok.IsKeyword():
			text := getTextLiteral()
			text.Pos = np.pos
			text.Value = np.lit
			text.Token = np.tok
			fieldExpr = text
		default:
			if p.err != nil {
				p.err(np.pos, "comparable: TEXT, STRING or Keyword expected on first element but got: "+np.lit)
			}
			return nil, ErrInvalidFilterSyntax
		}
		member.Fields = append(member.Fields, fieldExpr)
	}

	return member, nil
}
