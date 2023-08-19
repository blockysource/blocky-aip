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

// Compile-time check that *RestrictionExpr implements SimpleExpr.
var (
	_ SimpleExpr = (*RestrictionExpr)(nil)
)

// RestrictionExpr express a relationship between a comparable value and a
// single argument. When the restriction only specifies a comparable
// without an operator, this is a global restriction.
//
// 	EBNF:
//
// 	restriction
// 		: comparable [comparator arg]
// 		;
//
// RestrictionExpr implements SimpleExpr.
type RestrictionExpr struct {
	// Pos is the position of the restriction.
	Pos token.Position

	// Comparable is the comparable expression.
	Comparable ComparableExpr

	// Comparator is the comparator expression.
	Comparator *ComparatorLiteral

	// Arg is the argument expression.
	Arg ArgExpr
}

// IsGlobal returns true if the restriction is global.
func (r *RestrictionExpr) IsGlobal() bool {
	return r.Comparator == nil && r.Arg == nil
}

// String returns the string representation of the restriction.
func (r *RestrictionExpr) String() string {
	if r.Comparator == nil && r.Arg == nil {
		return r.Comparable.String()
	}
	return fmt.Sprintf("%s %s %s", r.Comparable, r.Comparator, r.Arg)
}

// UnquotedString returns the unquoted string.
func (r *RestrictionExpr) UnquotedString() string {
	if r.Comparator == nil && r.Arg == nil {
		return r.Comparable.UnquotedString()
	}
	return fmt.Sprintf("%s %s %s", r.Comparable.UnquotedString(), r.Comparator.String(), r.Arg.UnquotedString())
}

func (r *RestrictionExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	if r.Comparator == nil && r.Arg == nil {
		r.Comparable.WriteStringTo(sb, unquoted)
	} else {
		r.Comparable.WriteStringTo(sb, unquoted)
		sb.WriteRune(' ')
		r.Comparator.WriteStringTo(sb, unquoted)
		sb.WriteRune(' ')
		r.Arg.WriteStringTo(sb, unquoted)
	}
}

// Position returns the position of the restriction.
func (r *RestrictionExpr) Position() token.Position { return r.Pos }
func (*RestrictionExpr) isSimpleExpr()              {}
func (*RestrictionExpr) isAstExpr()                 {}
