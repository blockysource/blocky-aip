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

func (p *Parser) parseMemberExpr(nameParts []namePart) (*ast.MemberExpr, error) {
	member := getMemberExpr()
	defer putNameParts(nameParts)

	// If the member is an argument, the nameParts might actually be a single value.
	var (
		isUnfinishedTextLiteral bool
		unfinishedTextLiteral   *ast.TextLiteral
		isValueAssigned         bool
	)
	for i, np := range nameParts {
		if i == 0 {
			switch np.tok {
			case token.TEXT, token.TIMESTAMP:
				text := getTextLiteral()
				text.Pos = np.pos
				text.Value = np.lit
				if np.tok == token.TIMESTAMP {
					text.IsTimestamp = true
				}
				member.Value = text
				isValueAssigned = true
			case token.STRING:
				sl := getStringLiteral()
				sl.Pos = np.pos
				sl.Value = np.lit
				isValueAssigned = true
				member.Value = sl
			case token.MINUS:
				// this is a character for the TEXT literal.
				if len(nameParts) == 1 {
					tl := getTextLiteral()
					tl.Pos = np.pos
					tl.Value = np.lit
					isValueAssigned = true
					member.Value = tl
					continue
				}

				isUnfinishedTextLiteral = true
				unfinishedTextLiteral = getTextLiteral()
				unfinishedTextLiteral.Pos = np.pos
				unfinishedTextLiteral.Value = np.lit
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
		switch np.tok {
		case token.TEXT, token.TIMESTAMP:
			var text *ast.TextLiteral
			if isUnfinishedTextLiteral {
				unfinishedTextLiteral.Value += np.lit
				if !isValueAssigned {
					member.Value = unfinishedTextLiteral
				} else {
					member.Fields = append(member.Fields, unfinishedTextLiteral)
				}
				isUnfinishedTextLiteral = false
				continue
			} else {
				text = getTextLiteral()
				text.Pos = np.pos
				text.Value = np.lit
			}
			fieldExpr = text
		case token.MINUS:
			if isUnfinishedTextLiteral {
				unfinishedTextLiteral.Value += np.lit
				if i == len(nameParts)-1 {
					if !isValueAssigned {
						member.Value = unfinishedTextLiteral
					} else {
						member.Fields = append(member.Fields, unfinishedTextLiteral)
					}
					isUnfinishedTextLiteral = false
				} else {
					continue
				}
			} else {
				unfinishedTextLiteral = getTextLiteral()
				unfinishedTextLiteral.Pos = np.pos
				unfinishedTextLiteral.Value = np.lit
				isUnfinishedTextLiteral = true
			}
		case token.STRING:
			sl := getStringLiteral()
			sl.Pos = np.pos
			sl.Value = np.lit
			fieldExpr = sl
		case token.AND, token.OR, token.NOT:
			// This is a keyword.
			kw := getKeywordExpr()
			kw.Pos = np.pos
			switch np.tok {
			case token.AND:
				kw.Typ = ast.AND
			case token.OR:
				kw.Typ = ast.OR
			case token.NOT:
				kw.Typ = ast.NOT
			}
			fieldExpr = kw
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
