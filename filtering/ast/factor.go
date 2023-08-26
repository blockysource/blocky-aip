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

// FactorExpr is a factor expression which contains one or more disjunction terms.
// The terms are separated by the OR operator.
type FactorExpr struct {
	Pos token.Position

	// Terms is a list of disjunction term expressions separated by the OR operator.
	Terms []*TermExpr
}

func (e *FactorExpr) UnquotedString() string {
	var sb strings.Builder
	for i, t := range e.Terms {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(t.UnquotedString())
	}
	return sb.String()
}

func (e *FactorExpr) String() string {
	var sb strings.Builder
	for i, t := range e.Terms {
		if i > 0 {
			sb.WriteString(" OR ")
		}
		sb.WriteString(t.String())
	}
	return sb.String()
}

func (e *FactorExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, t := range e.Terms {
		if i > 0 {
			sb.WriteString(" OR ")
		}
		t.WriteStringTo(sb, unquoted)
	}
}

func (e *FactorExpr) Position() token.Position { return e.Pos }

func (*FactorExpr) isAstExpr() {}
