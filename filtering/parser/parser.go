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
	"errors"
	"strings"
	"sync"

	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/filtering/scanner"
	"github.com/blockysource/blocky-aip/filtering/token"
)

// Parser is responsible for parsing the input string filter into an AST.
type Parser struct {
	src     string
	scanner scanner.Scanner

	// err is the error handler.
	err scanner.ErrorHandler

	prevCtx, curCtx   parsingContext
	strictWhiteSpaces bool

	argMemberHandlers   []MemberHandler
	restrictionHandlers []RestrictionHandler
	funcCallHandlers    []FunctionCallHandler
}

// ParserOption changes the behavior of the parser.
type ParserOption func(p *Parser)

// StrictWhitespacesOption makes the parser to fail if there are more than one whitespace between specific expresions.
func StrictWhitespacesOption() ParserOption {
	return func(p *Parser) {
		p.strictWhiteSpaces = true
	}
}

// ArgMemberModifierOption adds modifiers to the parser that are applied to the member expressions that are arguments.
func ArgMemberModifierOption(modifiers ...MemberHandler) ParserOption {
	return func(p *Parser) {
		p.argMemberHandlers = append(p.argMemberHandlers, modifiers...)
	}
}

// ErrorHandlerOption sets the error handler of the parser.
func ErrorHandlerOption(err scanner.ErrorHandler) ParserOption {
	return func(p *Parser) {
		p.err = err
	}
}

// RestrictionOption changes the behavior of the parser when parsing restriction expressions.
func RestrictionOption(modifiers ...RestrictionHandler) ParserOption {
	return func(p *Parser) {
		p.restrictionHandlers = append(p.restrictionHandlers, modifiers...)
	}
}

// FuncCallHandlerOption changes the behavior of the parser when parsing function calls.
func FuncCallHandlerOption(handlers ...FunctionCallHandler) ParserOption {
	return func(p *Parser) {
		p.funcCallHandlers = append(p.funcCallHandlers, handlers...)
	}
}

// NewParser creates a new parser with the given options.
func NewParser(src string, opts ...ParserOption) *Parser {
	p := &Parser{src: src}

	for _, opt := range opts {
		opt(p)
	}

	p.scanner.Reset(src, p.err)

	return p
}

// Reset resets the parser with the given input string.
func (p *Parser) Reset(src string, opts ...ParserOption) {
	p.src = src
	for _, opt := range opts {
		opt(p)
	}
	p.scanner.Reset(src, p.err)
}

// ErrInvalidFilterSyntax is returned when the input string filter has invalid syntax.
var ErrInvalidFilterSyntax = errors.New("invalid filter")

// Parse parses the input string filter into an AST.
// If the input was an empty string, the returned ParsedFilter will have a nil Expr.
func (p *Parser) Parse() (*ParsedFilter, error) {
	pf := getParsedFilter()
	if p.src == "" {
		return pf, nil
	}
	unsetFn := p.setCurrentContext(parsingContextFilter)
	defer unsetFn()

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	pos, tok, lit := p.scanner.Scan()
	if tok != token.EOF {
		if p.err != nil {
			p.err(pos, "expr: EOF expected but got: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	pf.Expr = expr

	return pf, nil
}

func (p *Parser) setCurrentContext(ctx parsingContext) (unset func()) {
	prevCtx := p.prevCtx
	curCtx := p.curCtx
	p.prevCtx = curCtx
	p.curCtx = ctx
	return func() {
		p.curCtx = curCtx
		p.prevCtx = prevCtx
	}
}

func (p *Parser) parseExpr() (*ast.Expr, error) {
	// Expression is a single or 'AND' separated sequences.
	// Parse the first sequence.
	unsetFn := p.setCurrentContext(parsingContextExpr)
	defer unsetFn()

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

func (p *Parser) parseSequenceExpr() (*ast.SequenceExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextSequence)
	defer unsetFn()

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

func (p *Parser) parseFactorExpr() (*ast.FactorExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextFactor)
	defer unsetFn()

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

func (p *Parser) parseTermExpr() (*ast.TermExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextTerm)
	defer unsetFn()

	// The term is either unary or simple term.
	// Check if the next token is unary operator.
	var (
		isUnary  bool
		unaryPos token.Position
		unaryTok token.Token
		unaryLit string
	)

	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		switch tok {
		case token.NOT, token.MINUS:
			isUnary = true
			unaryPos = pos
			unaryTok = tok
			unaryLit = lit
			return true
		}
		return false
	})
	te := getTermExpr()
	if !isUnary {
		// The term is a simple term.
		simple, err := p.parseSimpleExpr()
		if err != nil {
			return nil, err
		}
		te.Pos = simple.Position()
		te.Expr = simple
		return te, nil
	}
	te.Pos = unaryPos
	te.UnaryOp = unaryLit
	if unaryTok == token.NOT {
		ws := p.scanner.SkipWhitespace()
		if ws == 0 {
			if p.err != nil {
				p.err(unaryPos, "term: WS expected after NOT operator")
			}
			return nil, ErrInvalidFilterSyntax
		}

		if p.strictWhiteSpaces && ws > 1 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "term: only one WS is allowed between NOT operator and simple term")
			}
			return nil, ErrInvalidFilterSyntax
		}
	}
	// MINUS
	simple, err := p.parseSimpleExpr()
	if err != nil {
		return nil, err
	}
	te.Expr = simple
	return te, nil
}

func (p *Parser) parseSimpleExpr() (ast.SimpleExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextSimple)
	defer unsetFn()

	var isComposite bool
	p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
		if tok == token.LPAREN {
			isComposite = true
		}
		return false
	})

	if isComposite {
		return p.parseCompositeExpr()
	}
	return p.parseRestrictionExpr()
}

func (p *Parser) parseCompositeExpr() (*ast.CompositeExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextComposite)
	defer unsetFn()

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

func (p *Parser) parseRestrictionExpr() (*ast.RestrictionExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextRestriction)
	defer unsetFn()

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
		// The restriction is a global comparable expression.
		for _, modifier := range p.restrictionHandlers {
			if modifier(re) {
				break
			}
		}
		return re, nil
	}

	if n == 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: WS expected after comparable expression")
		}
		return nil, ErrInvalidFilterSyntax
	}

	if p.strictWhiteSpaces && n > 1 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: only one WS is allowed between comparable expression and comparator")
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

	if n == 0 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: WS expected after comparator")
		}
		return nil, ErrInvalidFilterSyntax
	}

	if p.strictWhiteSpaces && n > 1 {
		if p.err != nil {
			p.err(p.scanner.Pos(), "restriction: only one WS is allowed between comparator and argument")
		}
		return nil, ErrInvalidFilterSyntax
	}

	// Parse the argument.
	arg, err := p.parseArgExpr()
	if err != nil {
		return nil, err
	}
	re.Arg = arg

	for _, modifier := range p.restrictionHandlers {
		if modifier(re) {
			break
		}
	}

	return re, nil
}

type namePart struct {
	pos token.Position
	lit string
	tok token.Token
}

var namePartsPool = sync.Pool{
	New: func() any {
		return make([]namePart, 0, 10)
	},
}

func getNameParts() []namePart {
	v := namePartsPool.Get()
	if v == nil {
		return nil
	}
	return v.([]namePart)
}

func putNameParts(nameParts []namePart) {
	if nameParts == nil {
		return
	}

	nameParts = nameParts[:0]
	namePartsPool.Put(nameParts)
}

func (p *Parser) parseComparableExpr() (ast.ComparableExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextComparable)
	defer unsetFn()

	nameParts := getNameParts()

	pos, tok, lit := p.scanner.Scan()
	switch tok {
	case token.STRING, token.TEXT:
	default:
		if !tok.IsKeyword() {
			if p.err != nil {
				p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
			}
			putNameParts(nameParts)
			return nil, ErrInvalidFilterSyntax
		}
	}

	nameParts = append(nameParts, namePart{
		pos: pos,
		lit: lit,
		tok: tok,
	})

	var i int
	for {
		if i > 0 {
			pos, tok, lit = p.scanner.Scan()
			switch tok {
			case token.TEXT, token.STRING:
			default:
				if !tok.IsKeyword() {
					if p.err != nil {
						p.err(pos, "comparable: STRING, TEXT or Keyword expected but got: '"+lit+"'")
					}
					putNameParts(nameParts)
					return nil, ErrInvalidFilterSyntax
				}
			}
			nameParts = append(nameParts, namePart{
				pos: pos,
				lit: lit,
				tok: tok,
			})
		}
		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			// Expects a dot
			pt = tok
			if tok == token.PERIOD {
				return true
			}
			return false
		})

		switch pt {
		case token.PERIOD:
			i++
		case token.LPAREN:
			// This is a function call.
			fc, err := p.parseFuncLiteral(nameParts)
			if err != nil {
				return nil, err
			}

			for _, handler := range p.funcCallHandlers {
				if handler(fc) {
					break
				}
			}
		default:
			// This is the end of the member expression.
			return p.parseMemberLiteral(nameParts, p.prevCtx == parsingContextArg)
		}
	}
}

func (p *Parser) parseFuncLiteral(nameParts []namePart) (*ast.FunctionCall, error) {
	unsetFn := p.setCurrentContext(parsingContextFunction)
	defer unsetFn()

	fl := getFunctionCall()
	fl.Pos = nameParts[0].pos

	defer putNameParts(nameParts)

	for _, np := range nameParts {
		switch np.tok {
		case token.TEXT:
			text := getTextLiteral()
			text.Pos = np.pos
			text.Value = np.lit
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
			p.err(p.scanner.Pos(), "function: no WS is allowed before the first argument")
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
			p.err(p.scanner.Pos(), "function: no WS is allowed after the last argument")
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

func (p *Parser) parseMemberLiteral(nameParts []namePart, inArg bool) (*ast.MemberExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextMember)
	defer unsetFn()

	member := getMemberLiteral()
	defer putNameParts(nameParts)

	// If the member is an argument, the nameParts might actually be a single value.

	for i, np := range nameParts {
		if i == 0 {
			switch np.tok {
			case token.TEXT:
				text := getTextLiteral()
				text.Pos = np.pos
				text.Value = np.lit
				member.Value = text
			case token.STRING:
				sl := getStringLiteral()
				sl.Pos = np.pos
				sl.Value = np.lit

				// Check if the value has a wilcard prefix.
				if strings.HasPrefix(sl.Value, "*") {
					sl.IsPrefixBased = true
				}

				// Check if the value has a wilcard suffix.
				if strings.HasSuffix(sl.Value, "*") {
					sl.IsSuffixBased = true
				}
				member.Value = sl
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
		case token.TEXT:
			text := getTextLiteral()
			text.Pos = np.pos
			text.Value = np.lit
			fieldExpr = text
		case token.STRING:
			sl := getStringLiteral()
			sl.Pos = np.pos
			sl.Value = np.lit
			// Check if the value has a wilcard prefix.
			if strings.HasPrefix(sl.Value, "*") {
				sl.IsPrefixBased = true
			}
			// Check if the value has a wilcard suffix.
			if strings.HasSuffix(sl.Value, "*") {
				sl.IsSuffixBased = true
			}
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

	if inArg {
		for _, mm := range p.argMemberHandlers {
			if mm(member) {
				break
			}
		}
	}
	return member, nil
}

func (p *Parser) parseArgListExpr() (*ast.ArgListExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextArgList)
	defer unsetFn()

	argList := getArgListExpr()

	i := 0
	for {
		// Skip possible whitespaces.
		n := p.scanner.SkipWhitespace()
		if p.strictWhiteSpaces && n > 0 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "argList: no WS is allowed before, between or after arguments")
			}
			return nil, ErrInvalidFilterSyntax
		}

		var pt token.Token
		p.scanner.Peek(func(pos token.Position, tok token.Token, lit string) bool {
			pt = tok
			if tok == token.COMMA && i > 0 {
				return true
			}
			return false
		})
		if (i > 0 && pt != token.COMMA) || (i == 0 && pt == token.RPAREN) {
			return argList, nil
		}

		// Skip possible whitespaces.
		n = p.scanner.SkipWhitespace()
		if p.strictWhiteSpaces && n > 0 {
			if p.err != nil {
				p.err(p.scanner.Pos(), "argList: no WS is allowed before, between or after arguments")
			}
			return nil, ErrInvalidFilterSyntax
		}

		// Parse the argument.
		arg, err := p.parseArgExpr()
		if err != nil {
			return nil, err
		}
		argList.Args = append(argList.Args, arg)
		i++
	}
}

func (p *Parser) parseArgExpr() (ast.ArgExpr, error) {
	unsetFn := p.setCurrentContext(parsingContextArg)
	defer unsetFn()

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
	case token.HAS:
		cl.Type = ast.HAS
	default:
		if p.err != nil {
			p.err(pos, "restriction: unknown comparator: "+lit)
		}
		return nil, ErrInvalidFilterSyntax
	}

	return cl, nil
}

// ParsedFilter is a parsed filter expression.
type ParsedFilter struct {
	// Expr is a parsed filter expression, possibly nil (for empty filter).
	Expr *ast.Expr
}

// Free frees the resource associated with the parsed filter.
// This should be used in a defer statement immediately after calling Parse.
// No further use of any filter expressions is allowed after calling Free.
func (p *ParsedFilter) Free() {
	putParsedFilter(p)
}

type parsingContext int

const (
	parsingContextFilter parsingContext = iota
	parsingContextExpr
	parsingContextSequence
	parsingContextFactor
	parsingContextTerm
	parsingContextSimple
	parsingContextRestriction
	parsingContextComparable
	parsingContextMember
	parsingContextFunction
	parsingContextComposite
	parsingContextArgList
	parsingContextArg
)
