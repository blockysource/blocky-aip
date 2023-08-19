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

func (p *Parser) parseFuncCall(nameParts []namePart) (*ast.FunctionCall, error) {
	fl := getFunctionCall()
	fl.Pos = nameParts[0].pos

	defer putNameParts(nameParts)

	for _, np := range nameParts {
		switch np.tok {
		case token.TEXT, token.TIMESTAMP:
			text := getTextLiteral()
			text.Pos = np.pos
			text.Value = np.lit
			if np.tok == token.TIMESTAMP {
				text.IsTimestamp = true
			}
			fl.Name = append(fl.Name, text)
		case token.AND, token.OR, token.NOT:
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
			fl.Name = append(fl.Name, kw)
		default:
			if p.err != nil {
				p.err(np.pos, "function: TEXT, AND, OR or NOT expected but got: "+np.lit)
			}
			return nil, ErrInvalidFilterSyntax
		}
	}

	pos, tok, lit := p.scanner.Scan()
	if tok != token.LPAREN {
		if p.err != nil {
			p.err(pos, "function: LPAREN expected at the beginning of function call but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}
	fl.Lparen = pos

	var isRParen bool
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		isRParen = tok == token.RPAREN
		fl.Rparen = pos
		return isRParen
	})
	if isRParen {
		return fl, nil
	}

	// Skip possible whitespaces before the first argument.
	n := p.scanner.SkipWhitespace()
	if p.strictWhiteSpaces && n > 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "function: no whitespace is allowed before the first argument")
		}
		return nil, ErrInvalidFilterSyntax
	}

	// Parse the first argument.
	argList, err := p.parseArgListExpr()
	if err != nil {
		return nil, err
	}
	fl.ArgList = argList

	// Skip possible whitespaces.
	n = p.scanner.SkipWhitespace()
	if p.strictWhiteSpaces && n > 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "function: no whitespace is allowed after the last argument")
		}
		return nil, ErrInvalidFilterSyntax
	}

	pos, tok, lit = p.scanner.Scan()
	if tok != token.RPAREN {
		if p.err != nil {
			p.err(pos, "function: RPAREN expected at the end of function call but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}
	fl.Rparen = pos

	return fl, nil
}
