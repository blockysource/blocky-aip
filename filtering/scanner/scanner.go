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
	"unicode"
	"unicode/utf8"

	"github.com/blockysource/blocky-aip/filtering/token"
)

type Scanner struct {
	src string

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	err        ErrorHandler
	ErrorCount int

	initialized bool

	peeked struct {
		pos      token.Position
		tok      token.Token
		lit      string
		isPeeked bool
	}
}

type ErrorHandler func(pos token.Position, msg string)

func NewScanner(src string) *Scanner {
	return &Scanner{src: src}
}

// Reset prepares the scanner s to tokenize the text src by setting the scanner at the beginning of src.
func (s *Scanner) Reset(src string, err ErrorHandler) {
	s.src = src
	s.err = err
	s.ch = ' '
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
		isText   bool
		isString bool
	)
	switch s.ch {
	case ' ', '\t', '\n', '\r':
		tok = token.WS
		lit = " "
	case '=':
		tok = token.EQUAL
		lit = "="
	case ':':
		tok = token.HAS
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
	case '.':
		tok = token.PERIOD
		lit = "."
	case '-':
		tok = token.MINUS
		lit = "-"
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
	default:
		isText = true
	}

	if !isText && !isString {
		s.next() // consume the token character
		return
	}
	if isString {
		tok = token.STRING
		lit = s.scanString()
	} else {
		tok = token.TEXT
		lit = s.scanText()

		switch lit {
		case "AND":
			tok = token.AND
		case "OR":
			tok = token.OR
		case "NOT":
			tok = token.NOT
		}
	}
	return
}

func (s *Scanner) scanString() string {
	// opening quote is already consumed
	offset := s.offset
	sum := 0
	for {
		ch, w := s.next()
		sum += w
		if isEOF(ch) {
			s.error(s.offset, "unterminated string")
			break
		}
		if isQuote(ch) {
			s.next() // consume closing quote
			break
		}
	}
	return s.src[offset : offset+sum-1]
}

func isQuote(ch rune) bool {
	return ch == '\'' || ch == '"'
}

func (s *Scanner) scanText() string {
	offset := s.offset
	sum := 0
	for {
		ch, w := s.next()
		sum += w
		if isEOF(ch) {
			return s.src[offset-1:]
		}
		if isWhitespace(ch) || isPeriod(ch) || ch == '(' || ch == ')' || ch == ',' || s.isComparator(ch) {
			break
		}
	}
	return s.src[offset-1 : offset+sum-1]
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
		s.ch = ch
		return ch, w
	}
	s.ch = eof
	return s.ch, utf8.RuneLen(rune(eof))
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
