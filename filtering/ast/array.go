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

package ast

import (
	"strings"

	"github.com/blockysource/blocky-aip/token"
)

// Compile-time checks to ensure type implements desired interfaces.
var (
	_ ComparableExpr = (*ArrayExpr)(nil)
)

// ArrayExpr is an extended expression not defined by the standard EBNF.
// It is used to represent an array of expressions.
// Its values are separated by a comma (,), and surrounded by square brackets ([]).
// The values are either a member expression or a function call (ComparableExpr).
// By default, the parser is not aware of the array expression.
// It needs to be set as an extension to the parser.
type ArrayExpr struct {
	// LBracket is the left bracket token.
	LBracket token.Position

	// Elements is a list of values.
	Elements []ComparableExpr

	// RBracket is the right bracket token.
	RBracket token.Position
}

// Position returns the position of the left bracket.
func (a *ArrayExpr) Position() token.Position { return a.LBracket }

// String returns the string representation of the array.
func (a *ArrayExpr) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, v := range a.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteRune(']')
	return sb.String()
}

// UnquotedString returns the unquoted string.
func (a *ArrayExpr) UnquotedString() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, v := range a.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.UnquotedString())
	}
	sb.WriteRune(']')
	return sb.String()
}

// WriteStringTo writes the string representation of the value to the builder.
// If unquoted argument is set to true, the StringLiterals do not write its string
// representation surrounded with quotes.
func (a *ArrayExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	sb.WriteRune('[')
	for i, v := range a.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		v.WriteStringTo(sb, unquoted)
	}
	sb.WriteRune(']')
}

// isAstExpr is a marker method for the interface.
func (*ArrayExpr) isAstExpr() {}

// isComparableExpr is a marker method for the interface.
func (*ArrayExpr) isComparableExpr() {}

// isArgExpr is a marker method for the interface.
func (*ArrayExpr) isArgExpr() {}
