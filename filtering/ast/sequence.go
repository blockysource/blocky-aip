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

// SequenceExpr is composed of one or more whitespace (WS) separated factors.
// A sequence expresses a logical relationship between 'factors' where
// the ranking of a filter result may be scored according to the number
// factors that match and other such criteria as the proximity of factors
// to each other within a document.
// When filters are used with exact match semantics rather than fuzzy
// match semantics, a sequence is equivalent to AND.
// Example: `New York Giants OR Yankees`
// The expression `New York (Giants OR Yankees)` is equivalent to the
// example.
type SequenceExpr struct {
	// Pos is the position of the sequence.
	Pos token.Position

	// Factors is a list of factor expressions separated by the WS operator.
	Factors []*FactorExpr

	// OpPos is the position of the AND operator after the sequence.
	// If there are no more sequences after this one, this value is undefined 0.
	OpPos token.Position
}

func (e *SequenceExpr) UnquotedString() string {
	var sb strings.Builder
	for i, f := range e.Factors {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(f.UnquotedString())
	}
	return sb.String()
}

func (e *SequenceExpr) Position() token.Position { return e.Pos }

func (e *SequenceExpr) String() string {
	var sb strings.Builder
	for i, f := range e.Factors {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(f.String())
	}
	return sb.String()
}

func (e *SequenceExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, f := range e.Factors {
		if i > 0 {
			sb.WriteString(" ")
		}
		f.WriteStringTo(sb, unquoted)
	}
}

func (*SequenceExpr) isAstExpr() {}
