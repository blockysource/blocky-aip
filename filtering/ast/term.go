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
	"fmt"
	"strings"

	"github.com/blockysource/blocky-aip/filtering/token"
)

// Compile-time checks to ensure type implements desired interfaces.
var _ AnyExpr = (*TermExpr)(nil)

// TermExpr is a unary negation expression.
// The unary negation operator is a hyphen (-) that precedes a term or
// the NOT operator that precedes a term with a WS.
// TermExpr are separated in the FactorExpr by the OR operator.
// The position of the OR operator is stored in the OrOpPos field, in a TermExpr which
// stands before the operator, i.e.:
//
//		`-foo OR bar` -> OrOpPos is the position of the OR operator after `-foo`.
//	                  It is stored in the '-foo' TermExpr.OrOpPos field.
type TermExpr struct {
	Pos token.Position

	// UnaryOp is the unary operator.
	UnaryOp string

	// Expr is the expression.
	Expr SimpleExpr

	// OrOpPos is the position of the OR operator after the term.
	// If there are no more terms after this one, this value is undefined 0.
	OrOpPos token.Position
}

// HasNegation returns true if the expression has a negation.
func (e *TermExpr) HasNegation() bool {
	return e.UnaryOp == "-" || e.UnaryOp == "NOT"
}

// UnquotedString returns the unquoted string representation of the expression.
func (e *TermExpr) UnquotedString() string {
	switch e.UnaryOp {
	case "-":
		return fmt.Sprintf("-%s", e.Expr.UnquotedString())
	case "NOT":
		return fmt.Sprintf("NOT %s", e.Expr.UnquotedString())
	default:
		return e.Expr.UnquotedString()
	}
}

// String returns the string representation of the expression.
func (e *TermExpr) String() string {
	switch e.UnaryOp {
	case "-":
		return fmt.Sprintf("-%s", e.Expr.String())
	case "NOT":
		return fmt.Sprintf("NOT %s", e.Expr.String())
	default:
		return e.Expr.String()
	}
}

// WriteStringTo writes the string representation of the expression to the given builder.
func (e *TermExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	switch e.UnaryOp {
	case "-":
		sb.WriteString("-")
		e.Expr.WriteStringTo(sb, unquoted)
	case "NOT":
		sb.WriteString("NOT ")
		e.Expr.WriteStringTo(sb, unquoted)
	default:
		e.Expr.WriteStringTo(sb, unquoted)
	}
}

// Position returns the position of the expression.
func (e *TermExpr) Position() token.Position {
	return e.Pos
}

func (*TermExpr) isAstExpr() {}
