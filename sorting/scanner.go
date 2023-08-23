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

package sorting

import (
	"unicode/utf8"
)

type token int

const (
	// Special tokens
	illegal token = iota
	eof_tok
	ws

	// Literals
	field_tok

	// Keywords
	asc_tok
	desc_tok

	// Special characters
	period_tok // .
	comma_tok  // ,
)

var tokens = [...]string{
	illegal: "ILLEGAL",
	eof_tok: "EOF",
	ws:      "WS",

	field_tok: "FIELD",

	asc_tok:  "ASC",
	desc_tok: "DESC",

	period_tok: ".",
}

func (t token) String() string {
	if int(t) < len(tokens) {
		return tokens[t]
	}
	return ""
}

type position int

type scanner struct {
	src string

	// state
	ch     rune // current char
	offset int
}

func (s *scanner) scan() (pos position, tok token, lit string) {
	pos = position(s.offset)
	switch ch := s.ch; {
	case isWhitespace(ch):
		tok = ws
		lit = " "
	case isPeriod(ch):
		tok = period_tok
		lit = "."
	case ch == ',':
		tok = comma_tok
		lit = ","
	case ch == eof:
		tok = eof_tok
		lit = "EOF"
		return pos, tok, lit
	default:
		lit = s.scanField()
		switch lit {
		case "ASC", "asc":
			tok = asc_tok
		case "DESC", "desc":
			tok = desc_tok
		default:
			tok = field_tok
		}
		return pos, tok, lit
	}

	s.next()
	return pos, tok, lit
}

func (s *scanner) peekToken(fn func(pos position, tok token, lit string) bool) {
	offset := s.offset
	ch := s.ch

	pos, tok, lit := s.scan()
	if !fn(pos, tok, lit) {
		s.offset = offset
		s.ch = ch
	}
}

func (s *scanner) scanField() string {
	offset := s.offset
	sum := 0
	for {
		ch, w := s.next()
		sum += w
		if isEOF(ch) {
			return s.src[offset-1:]
		}

		if isWhitespace(ch) || isPeriod(ch) || isComma(ch) {
			break
		}
	}

	return s.src[offset-1 : offset+sum-1]
}

func (s *scanner) next() (ch rune, w int) {
	if s.offset < len(s.src) {
		ch, w = utf8.DecodeRuneInString(s.src[s.offset:])
		s.offset += w
		s.ch = ch
		return ch, w
	}
	s.ch, w = eof, 0
	return s.ch, 0
}

func (s *scanner) peek() rune {
	if s.offset < len(s.src) {
		ch, _ := utf8.DecodeRuneInString(s.src[s.offset:])
		return ch
	}
	return eof
}

func (s *scanner) skipWhitespace() int {
	var n int
	for isWhitespace(s.ch) {
		s.next()
		n++
	}
	return n
}

func (s *scanner) init(order string) {
	s.src = order
	s.offset = 0
	s.next()
	if s.ch == bom {
		s.next()
	}
}

const (
	eof = -1
	bom = 0xFEFF
)

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
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

func isComma(ch rune) bool {
	return ch == ','
}
