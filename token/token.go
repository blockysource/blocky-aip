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
	// STRING is a special type of literal, which is not defined by the standard EBNF.
	// It defines a string literal.
	STRING // "abc" or 'abc'
	non_string_literal_beg
	// IDENT is a special type of literal, which is not defined by the standard EBNF.
	// It defines either an identifier or a keyword.
	IDENT // abc
	// TIMESTAMP is a special type of literal, which is not defined by the standard EBNF.
	// It is used to represent a IDENT literal, which is a valid timestamp.
	TIMESTAMP // 2021-01-01T00:00:00Z

	numbers_beg
	// NUMERIC is a special type of literal, which is not defined by the standard EBNF.
	// It defines a numeric literal with a decimal point or an exponent.
	NUMERIC // .123 | 123. | 3.0e+123

	integers_beg
	// INT is a special type of literal, which is not defined by the standard EBNF.
	// It defines an integer literal.
	INT // -123
	// HEX is a special type of literal, which is not defined by the standard EBNF.
	// It defines a hexadecimal literal.
	HEX // 0x123
	// OCT is a special type of literal, which is not defined by the standard EBNF.
	// It defines an octal literal.
	OCT // 0o123
	integers_end
	numbers_end

	// DURATION is a special type of literal, which is not defined by the standard EBNF.
	// It defines a duration literal.
	DURATION // 1h | 1h30m | 1h30m30s

	keyword_beg
	// TRUE is a special types of literals, which are not defined by the standard EBNF,
	// that represent boolean true value.
	TRUE // true
	// FALSE is a special type of literal, which is not defined by the standard EBNF,
	// that represents a boolean false value.
	FALSE // false
	// NULL is a special type of literal, which is not defined by the standard EBNF.
	// It defines a null literal.
	NULL // null
	non_string_literal_end
	literal_end

	logical_beg
	AND // AND
	OR  // OR
	NOT // NOT
	logical_end
	ASC  // ASC
	DESC // DESC
	comparator_beg
	IN // IN extension to the standard
	keyword_end

	EQUAL // =
	LT    // <
	LEQ   // <=
	GT    // >
	GEQ   // >=
	NEQ   // !=
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
	ASTERISK      // *
	MINUS         // -
	additional_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	WS:      "WS",

	IDENT:     "IDENT",
	TIMESTAMP: "TIMESTAMP",
	NUMERIC:   "NUMERIC",
	INT:       "INT",
	HEX:       "HEX",
	OCT:       "OCT",
	DURATION:  "DURATION",
	TRUE:      "true",
	FALSE:     "false",
	NULL:      "null",

	STRING: "STRING",

	AND:  "AND",
	OR:   "OR",
	NOT:  "NOT",
	ASC:  "ASC",
	DESC: "DESC",
	IN:   "IN",

	EQUAL: "=",
	LT:    "<",
	LEQ:   "<=",
	GT:    ">",
	GEQ:   ">=",
	NEQ:   "!=",

	LPAREN:        "(",
	RPAREN:        ")",
	COMMA:         ",",
	PERIOD:        ".",
	BRACKET_OPEN:  "[",
	BRACKET_CLOSE: "]",
	BRACE_OPEN:    "{",
	BRACE_CLOSE:   "}",
	COLON:         ":",
	MINUS:         "-",
}

func (t Token) String() string        { return tokens[t] }
func (t Token) IsLiteral() bool       { return literal_beg < t && t < literal_end }
func (t Token) IsNonStringLit() bool  { return non_string_literal_beg < t && t < non_string_literal_end }
func (t Token) IsNumber() bool        { return numbers_beg < t && t < numbers_end }
func (t Token) IsInteger() bool       { return integers_beg < t && t < integers_end }
func (t Token) IsLogical() bool       { return logical_beg < t && t < logical_end }
func (t Token) IsBoolean() bool       { return t == TRUE || t == FALSE }
func (t Token) IsKeyword() bool       { return t > keyword_beg && t < keyword_end }
func (t Token) IsComparator() bool    { return comparator_beg < t && t < comparator_end || t == COLON }
func (t Token) IsAdditional() bool    { return additional_beg < t && t < additional_end }
func (t Token) IsIdent() bool         { return t == IDENT || t.IsKeyword() }
func (t Token) IsUnaryOperator() bool { return t == NOT || t == MINUS }

// Position is a position of the token in the source.
type Position int
