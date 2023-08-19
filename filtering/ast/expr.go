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

	"github.com/blockysource/blocky-aip/filtering/token"
)

// Expr may either be a conjunction (AND) of sequences or a simple sequence.
type Expr struct {
	// Pos is the position of the expression.
	Pos token.Position

	// Sequences is a list of sequence expressions.
	Sequences []*SequenceExpr
}

// Position returns the position of the expression.
func (e *Expr) Position() token.Position { return e.Pos }

func (e *Expr) UnquotedString() string {
	var sb strings.Builder
	for i, seq := range e.Sequences {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(seq.UnquotedString())
	}
	return sb.String()
}

// String returns the string representation of the expression.
func (e *Expr) String() string {
	var sb strings.Builder
	for i, seq := range e.Sequences {
		if i > 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString(seq.String())
	}
	return sb.String()
}

// WriteStringTo writes the string representation of the expression to the builder.
func (e *Expr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, seq := range e.Sequences {
		if i > 0 {
			sb.WriteString(" AND ")
		}
		seq.WriteStringTo(sb, unquoted)
	}
}

func (*Expr) isAstExpr() {}
