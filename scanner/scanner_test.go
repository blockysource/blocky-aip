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

package scanner_test

import (
	"testing"

	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		check     func(t *testing.T, s *scanner.Scanner)
		useStruct bool
		useArray  bool
		useIn     bool
		isErr     bool
	}{
		{
			name: "text",
			src:  `text`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %s", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "text with space",
			src:  `text with space`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %s", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %s", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.WS {
					t.Errorf("unexpected token: %s", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 5 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.IDENT {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "with" {
					t.Errorf("unexpected literal: %s", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 9 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.WS {
					t.Errorf("unexpected token: %s", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 10 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.IDENT {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "space" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "timestamp",
			src:  `2021-01-01T00:00:00Z`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.TIMESTAMP {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "2021-01-01T00:00:00Z" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "timestamp with TZ",
			src:  `2021-01-01T00:00:00+01:00`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.TIMESTAMP {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "2021-01-01T00:00:00+01:00" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "integer",
			src:  `123`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.INT {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "123" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "numeric",
			src:  `123.456`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.NUMERIC {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "123.456" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "numeric with exponent",
			src:  `123.456e+10`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.NUMERIC {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "123.456e+10" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "string",
			src:  `"string"`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.STRING {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != `string` {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "string with escaped quote",
			src:  `"string with \"escaped quote\""`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.STRING {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != `string with "escaped quote"` {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "text ws string",
			src:  `text "string"`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %s", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %s", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.WS {
					t.Errorf("unexpected token: %s", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 5 {
					t.Errorf("unexpected position: %d", pos)
				}

				if tok != token.STRING {
					t.Errorf("unexpected token: %s", tok)
				}

				if lit != "string" {
					t.Errorf("unexpected literal: %s", lit)
				}
			},
		},
		{
			name: "text.int",
			src:  `text.123`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 5 {
					t.Errorf("unexpected position: %v", pos)
				}

				if tok != token.INT {
					t.Errorf("unexpected token: %v", tok)
				}

				if lit != "123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "int.text",
			src:  `123.text`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.INT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "123" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 3 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %v", pos)
				}

				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}

				if lit != "text" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "text.text",
			src:  `text.text2`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 5 {
					t.Errorf("unexpected position: %v", pos)
				}

				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}

				if lit != "text2" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "path.int.text",
			src:  `path.123.text`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "path" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 4 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 5 {
					t.Errorf("unexpected position: %v", pos)
				}

				if tok != token.INT {
					t.Errorf("unexpected token: %v", tok)
				}

				if lit != "123" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if pos != 8 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if pos != 9 {
					t.Errorf("unexpected position: %v", pos)
				}

				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}

				if lit != "text" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "-123",
			src:  `-123`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.INT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "-123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "path.-123",
			src:  `text.-123`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.IDENT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "text" {
					t.Errorf("unexpected literal: %v", lit)
				}

				pos, tok, lit = s.Scan()
				if tok != token.PERIOD {
					t.Errorf("unexpected token: %v", tok)
				}

				pos, tok, lit = s.Scan()
				if tok != token.INT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "-123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int hour",
			src:  `1h`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1h" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int minute",
			src:  `1m`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1m" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int second",
			src:  `1s`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1s" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int millisecond",
			src:  `1ms`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1ms" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int microsecond",
			src:  `1us`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %v", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1us" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration int nanosecond",
			src:  `1ns`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1ns" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration float hour",
			src:  `1.5h`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position: %d", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1.5h" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration hour and minute",
			src:  `1h30m`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position %d", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1h30m" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "duration fractional second",
			src:  `1.5s`,
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position %d", pos)
				}
				if tok != token.DURATION {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "1.5s" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "hex int",
			src:  "0x123",
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position %d", pos)
				}
				if tok != token.HEX {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "0x123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "octal int",
			src:  "0123",
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position %d", pos)
				}
				if tok != token.OCT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "0123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
		{
			name: "octal with 'o' int",
			src:  "0o123",
			check: func(t *testing.T, s *scanner.Scanner) {
				pos, tok, lit := s.Scan()
				if pos != 0 {
					t.Errorf("unexpected position %d", pos)
				}
				if tok != token.OCT {
					t.Errorf("unexpected token: %v", tok)
				}
				if lit != "0o123" {
					t.Errorf("unexpected literal: %v", lit)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := scanner.New(tc.src, errHandler(t, tc.src, tc.isErr))
			tc.check(t, s)
		})
	}
}

func errHandler(t *testing.T, src string, wantsErr bool) func(pos token.Position, msg string) {
	return func(pos token.Position, msg string) {
		if !wantsErr {
			t.Error(src[pos:])
			t.Errorf("^ %v", msg)
		}
	}
}
