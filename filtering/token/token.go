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

package token

// Token is the set of lexical tokens of the Filtering language.
type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS

	literal_beg
	TEXT   // abc
	STRING // "abc" or 'abc'
	literal_end

	keyword_beg
	AND   // AND
	OR    // OR
	NOT   // NOT
	MINUS // -
	keyword_end

	comparator_beg
	EQUAL // =
	LT    // <
	LEQ   // <=
	GT    // >
	GEQ   // >=
	NEQ   // !=
	HAS   // :
	comparator_end

	additional_beg
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	PERIOD // .
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	TEXT:   "TEXT",
	STRING: "STRING",

	AND:   "AND",
	OR:    "OR",
	NOT:   "NOT",
	MINUS: "-",

	EQUAL: "=",
	LT:    "<",
	LEQ:   "<=",
	GT:    ">",
	GEQ:   ">=",
	NEQ:   "!=",
	HAS:   ":",

	LPAREN: "(",
	RPAREN: ")",
	COMMA:  ",",
	PERIOD: ".",
}

func (t Token) String() string     { return tokens[t] }
func (t Token) IsLiteral() bool    { return literal_beg < t && t < literal_end }
func (t Token) IsKeyword() bool    { return t > keyword_beg && t < keyword_end }
func (t Token) IsComparator() bool { return comparator_beg < t && t < comparator_end }
func (t Token) IsAdditional() bool { return additional_beg < t && t < additional_end }

// Position is a position of the token in the source.
type Position int
