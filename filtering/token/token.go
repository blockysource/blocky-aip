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
	TEXT // abc
	// TIMESTAMP is a special type of literal, which is not defined by the standard EBNF.
	// It is used to represent a TEXT literal, which is a valid timestamp.
	TIMESTAMP // 2021-01-01T00:00:00Z
	STRING         // "abc" or 'abc'
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
	IN    // IN extension to the standard
	comparator_end

	additional_beg
	LPAREN        // (
	RPAREN        // )
	COMMA         // ,
	PERIOD        // .
	BRACKET_OPEN  // [
	BRACKET_CLOSE // ]
	BRACE_OPEN    // {
	BRACE_CLOSE   // }
	COLON         // :
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	TEXT:           "TEXT",
	TIMESTAMP: "TIMESTAMP",
	STRING:         "STRING",

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
	IN:    "IN",

	LPAREN:        "(",
	RPAREN:        ")",
	COMMA:         ",",
	PERIOD:        ".",
	BRACKET_OPEN:  "[",
	BRACKET_CLOSE: "]",
	BRACE_OPEN:    "{",
	BRACE_CLOSE:   "}",
	COLON:         ":",
}

func (t Token) String() string     { return tokens[t] }
func (t Token) IsLiteral() bool    { return literal_beg < t && t < literal_end }
func (t Token) IsKeyword() bool    { return t > keyword_beg && t < keyword_end }
func (t Token) IsComparator() bool { return comparator_beg < t && t < comparator_end || t == COLON }
func (t Token) IsAdditional() bool { return additional_beg < t && t < additional_end }

// Position is a position of the token in the source.
type Position int
