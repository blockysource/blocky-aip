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

package scanner

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/blockysource/blocky-aip/token"
)

// Scanner is a token scanner.
type Scanner struct {
	src string

	// scanning state
	ch              rune // current character
	pch             rune // previous character
	prev            token.Token
	offset          int // character offset
	err             ErrorHandler
	useStructs      bool
	useArrays       bool
	useInComparator bool
	ErrorCount      int

	initialized bool

	peeked struct {
		pos      token.Position
		tok      token.Token
		lit      string
		isPeeked bool
	}
	useOrdering bool
}

// New creates a new scanner for the src string.
func New(src string, err ErrorHandler) *Scanner {
	s := &Scanner{}
	s.Reset(src, err)

	return s
}

// ErrorHandler is an error message handler.
type ErrorHandler func(pos token.Position, msg string)

// Reset prepares the scanner s to tokenize the text src by setting the scanner at the beginning of src.
func (s *Scanner) Reset(src string, err ErrorHandler) {
	s.src = src
	s.err = err
	s.ch = ' '
	s.pch = ' '
	s.offset = 0
	s.ErrorCount = 0
	s.initialized = true

	s.next()
	if s.ch == bom {
		s.next()
	}
}

// SkipWhitespace skips spaces, tabs, and newlines, and returns the number of bytes skipped.
func (s *Scanner) SkipWhitespace() int {
	return s.skipWhitespace()
}

// Peek peeks the next token returning the token position, the token, and its literal.
// If the scanner should consume the token,
func (s *Scanner) Peek(fn func(pos token.Position, tok token.Token, lit string) (consume bool)) {
	pos, tok, lit := s.scanOrPeekToken()
	if !fn(pos, tok, lit) {
		s.peeked.isPeeked = true
		s.peeked.pos = pos
		s.peeked.tok = tok
		s.peeked.lit = lit
	}
}

// Breakpoint is a breakpoint that can be used to restore the scanner to the current state.
type Breakpoint struct {
	ch     rune
	offset int
	peeked struct {
		pos      token.Position
		tok      token.Token
		lit      string
		isPeeked bool
	}
	createdByScanner bool
}

// Breakpoint creates a breakpoint that can be used to restore the scanner to the current state.
func (s *Scanner) Breakpoint() Breakpoint {
	return Breakpoint{
		ch:               s.ch,
		offset:           s.offset,
		peeked:           s.peeked,
		createdByScanner: true,
	}
}

// Restore restores the scanner to the state of the breakpoint.
func (s *Scanner) Restore(bp Breakpoint) {
	if !bp.createdByScanner {
		panic("breakpoint not created by scanner")
	}
	s.ch = bp.ch
	s.offset = bp.offset
	s.peeked = bp.peeked
}

// Scan scans the next token and returns the token position, the token, and its literal.
func (s *Scanner) Scan() (pos token.Position, tok token.Token, lit string) {
	return s.scanOrPeekToken()
}

func (s *Scanner) scanOrPeekToken() (pos token.Position, tok token.Token, lit string) {
	if !s.initialized {
		s.error(0, "scanner not initialized")
	}

	if s.peeked.isPeeked {
		s.peeked.isPeeked = false
		pos = s.peeked.pos
		s.peeked.pos = 0
		tok = s.peeked.tok
		s.peeked.tok = 0
		lit = s.peeked.lit
		s.peeked.lit = ""
		return pos, tok, lit
	}

	pos = s.pos()
	var (
		isText, isString, isNumeric bool
	)
	switch s.ch {
	case ' ', '\t', '\n', '\r':
		tok = token.WS
		lit = " "
	case '=':
		tok = token.EQUAL
		lit = "="
	case ':':
		tok = token.COLON
		lit = ":"
	case '(':
		tok = token.LPAREN
		lit = "("
	case ')':
		tok = token.RPAREN
		lit = ")"
	case ',':
		tok = token.COMMA
		lit = ","
	case '*':
		tok = token.ASTERISK
		lit = "*"
	case '.':
		tok = token.PERIOD
		lit = "."
		// Check if this is a path separator by checking if previous token was literal.
		if !s.prev.IsLiteral() {
			// This is a numeric after a period.
			isNumeric = true
		}
	case '-':
		tok = token.MINUS
		lit = "-"
		if p := s.peek(); isDecimal(p) || p == '.' {
			isNumeric = true
		}
	case '<':
		if s.peek() == '=' {
			s.next()
			pos = s.pos()
			tok = token.LEQ
			lit = "<="
		} else {
			tok = token.LT
			lit = "<"
		}
	case '>':
		if s.peek() == '=' {
			s.next()
			pos = s.pos()
			tok = token.GEQ
			lit = ">="
		} else {
			tok = token.GT
			lit = ">"
		}
	case '!':
		if s.peek() == '=' {
			s.next()
			s.pos()
			tok = token.NEQ
			lit = "!="
		} else {
			isText = true
		}
	case '"', '\'':
		isString = true
	case eof:
		tok = token.EOF
	case '[':
		tok = token.BRACKET_OPEN
		lit = "["
	case ']':
		tok = token.BRACKET_CLOSE
		lit = "]"
	case '{':
		tok = token.BRACE_OPEN
		lit = "{"
	case '}':
		tok = token.BRACE_CLOSE
		lit = "}"
	case '0':
		// This is number or numeric.
		isNumeric = true
	default:
		if isDecimal(s.ch) {
			isNumeric = true
		} else {
			isText = true
		}
	}

	if !isText && !isString && !isNumeric {
		s.next() // consume the token character
		return
	}

	switch {
	case isString:
		tok = token.STRING
		lit = s.scanString()
	case isNumeric:
		tok, lit = s.scanNumber()
	case isText:
		tok, lit = s.scanText()

		switch lit {
		case "true", "false":
			tok = token.BOOLEAN
		case "AND":
			tok = token.AND
		case "OR":
			tok = token.OR
		case "NOT":
			tok = token.NOT
		case "IN":
			tok = token.IN
		case "ASC", "asc":
			tok = token.ASC
		case "DESC", "desc":
			tok = token.DESC
		case "null":
			tok = token.NULL
		}
	}

	s.prev = tok
	return
}

func (s *Scanner) scanString() string {
	// opening quote is already consumed
	isEscape := false
	open := s.ch
	var sb strings.Builder
	for {
		ch, _ := s.next()
		if isEOF(ch) {
			s.error(s.offset, "unterminated string")
			break
		}

		// Handle escape characters.
		if ch == '\\' {
			isEscape = true
			continue
		}
		if ch == open {
			if isEscape {
				sb.WriteRune(ch)
				isEscape = false
				continue
			}
			s.next() // consume closing quote
			break
		}
		sb.WriteRune(ch)
	}
	return sb.String()
}

func (s *Scanner) scanText() (token.Token, string) {
	sum := utf8.RuneLen(s.ch)
	offset := s.offset - sum

	for {
		ch, w := s.next()
		if isBreaking(ch) {
			break
		}
		sum += w

		if isLetter(ch) || isDecimal(ch) || ch == '_' {
			continue
		}

		s.error(s.offset, fmt.Sprintf("invalid character: %q", ch))
		return token.ILLEGAL, ""
	}

	lit := s.src[offset : offset+sum]

	// This is a simple text literal.
	return token.IDENT, lit
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

func (s *Scanner) next() (ch rune, w int) {
	if s.offset < len(s.src) {
		ch, w := utf8.DecodeRuneInString(s.src[s.offset:])
		s.offset += w
		s.pch = s.ch
		s.ch = ch
		return ch, w
	}
	s.ch = eof
	return s.ch, 1
}

func (s *Scanner) peek() rune {
	offset := s.offset
	if offset < len(s.src) {
		ch, _ := utf8.DecodeRuneInString(s.src[offset:])
		return ch
	}
	return eof
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1
)

func (s *Scanner) backup() {
	s.offset -= utf8.RuneLen(s.ch)
}

func (s *Scanner) skipWhitespace() int {
	var n int
	if s.peeked.isPeeked {
		if s.peeked.tok != token.WS {
			return 0
		}
		s.peeked.isPeeked = false
		s.peeked.pos = 0
		s.peeked.tok = 0
		s.peeked.lit = ""
		n++
	}

	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
		n++
	}
	return n
}

func (s *Scanner) pos() token.Position {
	return token.Position(s.offset) - 1
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(token.Position(offs), msg)
	}
	s.ErrorCount++

}

func (s *Scanner) Pos() token.Position {
	return s.pos()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isPeriod(ch rune) bool {
	return ch == '.'
}

func isEOF(ch rune) bool {
	return ch == eof
}

func (s *Scanner) isComparator(ch rune) bool {
	if ch == '=' || ch == '<' || ch == '>' || ch == ':' {
		return true
	}

	if ch == '!' {
		if s.peek() == '=' {
			return true
		}
	}
	return false
}

func (s *Scanner) scanNumber() (token.Token, string) {
	if s.ch == '.' {
		return s.scanNumeric(1)
	}

	// Check if it is not a timestamp in RFC3339 format.
	offset := s.offset - 1
	sum := 1

	isNegative := s.ch == '-'

	tok := token.INT
	if s.ch == '0' {
		peek := s.peek()
		switch {
		case peek == 'x', peek == 'X':
			_, w := s.next()
			sum += w
			tok = token.HEX
		case peek == 'o':
			_, w := s.next()
			sum += w
			tok = token.OCT
		case isOctalDigit(peek):
			tok = token.OCT
		case peek == '.':
			return s.scanNumeric(1)
		case isDurationPrefix(peek):
			return s.scanDuration(1, false, false)
		case isBreaking(peek):
			s.next()
			return token.INT, s.src[offset : offset+sum]
		default:
			// This might be a timestamp in RFC3339 format.
			// i.e.: 0001-01-01T00:00:00Z
			return s.scanTimestamp(1)
		}
	}

	count := 1
	for {
		ch, w := s.next()
		if isNonPeriodBreaking(ch) {
			break
		}
		count++

		if isPeriod(ch) {
			return s.scanNumeric(count)
		}

		sum += w

		if tok == token.HEX {
			if !isHexDigit(ch) {
				s.error(s.offset, "invalid hexadecimal")
				return token.ILLEGAL, ""
			}
			continue
		}

		if count == 5 {
			if ch == '-' {
				// Check if the '-' is a timestamp separator of a year.
				if tok == token.HEX || isNegative {
					return token.INT, s.src[offset : offset+sum]
				}
				return s.scanTimestamp(count)
			}
		}

		if tok == token.OCT {
			if !isOctalDigit(ch) {
				s.error(s.offset, "invalid octal")
				return token.ILLEGAL, ""
			}
			continue
		}

		if !isDecimal(ch) {
			if isDurationPrefix(ch) {
				return s.scanDuration(count, false, false)
			}
			s.error(s.offset, "invalid decimal")
			return token.ILLEGAL, ""
		}
	}

	return tok, s.src[offset : offset+sum]
}

func (s *Scanner) scanNumeric(used int) (token.Token, string) {
	offset := s.offset
	offset -= used
	sum := used

	if s.ch != '.' {
		ch, w := s.next()
		sum += w
		if ch != '.' {
			panic("expected '.'")
		}
	}

	// If the next character is not a digit, then this is not a numeric,
	// but rather a text literal.
	// We need to revert the offset and return the INT token.
	if !isDecimal(s.peek()) {
		return token.INT, s.src[offset : offset+sum-1]
	}

	var (
		isExp     bool
		isExpSign bool
		hasExp    bool
	)
	for {
		ch, w := s.next()

		if isBreaking(ch) {
			break
		}
		sum += w

		switch ch {
		case 'e', 'E':
			if isExp || isExpSign {
				s.error(offset, "invalid numeric")
				return token.ILLEGAL, ""
			}
			isExp = true
			hasExp = true
			continue
		case '+', '-':
			if isExpSign || !isExp {
				s.error(offset, "invalid numeric")
				return token.ILLEGAL, ""
			}
			isExp = false
			isExpSign = true
			continue
		}
		isExpSign = false

		if !isDecimal(ch) {
			if isDurationPrefix(ch) {
				return s.scanDuration(sum, true, hasExp)
			}
			s.error(offset, "invalid decimal")
			return token.ILLEGAL, ""
		}

		if isExp {
			isExp = false
		}

	}

	if isExp || isExpSign {
		s.error(offset, "invalid numeric")
		return token.ILLEGAL, ""
	}

	return token.NUMERIC, s.src[offset : offset+sum]
}

func isHexDigit(ch rune) bool {
	return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f'
}

func isOctalDigit(ch rune) bool {
	return '0' <= ch && ch <= '7'
}

func isBreaking(ch rune) bool {
	return isEOF(ch) || isWhitespace(ch) || isPeriod(ch) || ch == '(' || ch == ')' || ch == ',' || isComparator(ch) ||
		ch == ']' || ch == '}' || ch == '[' || ch == '{'
}

func isComparator(ch rune) bool {
	return ch == '=' || ch == '<' || ch == '>' || ch == ':' || ch == '!'
}

func isNonPeriodBreaking(ch rune) bool {
	return isEOF(ch) || isWhitespace(ch) || ch == '(' || ch == ')' || ch == ',' || isComparator(ch) ||
		ch == ']' || ch == '}' || ch == '[' || ch == '{'
}
