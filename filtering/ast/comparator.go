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

	"github.com/blockysource/blocky-aip/token"
)

// ComparatorLiteral is the literal of a comparator with a position and type.
//
// EBNF:
//
// comparator
//
//	: LESS_EQUALS      # <=
//	| LESS_THAN        # <
//	| GREATER_EQUALS   # >=
//	| GREATER_THAN     # >
//	| NOT_EQUALS       # !=
//	| EQUALS           # =
//	| HAS              # :
//	| IN 			     # IN (extension to the standard)
//	;
type ComparatorLiteral struct {
	Pos  token.Position
	Type ComparatorType
}

// Position returns the position of the comparator.
func (c *ComparatorLiteral) Position() token.Position { return c.Pos }

// String returns the string representation of the comparator.
func (c *ComparatorLiteral) String() string { return c.Type.String() }

// UnquotedString returns the unquoted string.
func (c *ComparatorLiteral) UnquotedString() string { return c.Type.String() }

// WriteStringTo writes the string representation of the comparator to the builder.
func (c *ComparatorLiteral) WriteStringTo(sb *strings.Builder, unquoted bool) {
	sb.WriteString(c.Type.String())
}

func (*ComparatorLiteral) isAstExpr() {}

var _ComparatorTypeStrings = [...]string{
	LE:  "<=",
	LT:  "<",
	GE:  ">=",
	GT:  ">",
	NE:  "!=",
	EQ:  "=",
	HAS: ":",
	IN:  "IN",
}

// ComparatorType is a defined type for comparators.
type ComparatorType int

// String returns the string representation of the comparator.
func (c ComparatorType) String() string {
	if c < 0 || c > ComparatorType(len(_ComparatorTypeStrings)-1) {
		return fmt.Sprintf("ComparatorType(%d)", c)
	}
	return _ComparatorTypeStrings[c]
}

const (
	_ ComparatorType = iota
	// EQ is the equal to comparator.
	EQ
	// LE is the less than or equal to comparator.
	LE
	// LT is the less than comparator.
	LT
	// GE is the greater than or equal to comparator.
	GE
	// GT is the greater than comparator.
	GT
	// NE is the not equal to comparator.
	NE
	// HAS is the has comparator.
	HAS
	// IN is the in comparator that checks if a value is in a list of values.
	// NOTE: This is an extension to the standard.
	IN
)
