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

// Compile-
var (
	_ ValueExpr = (*StringLiteral)(nil)
	_ FieldExpr = (*StringLiteral)(nil)
)

// StringLiteral is a string literal. It is enclosed in double quotes.
// The string may or may not contain a special wildcard `*` character
// at the beginning or end of the string to indicate a prefix or
// suffix-based search within a restriction.
// If both are present, the string is a substring-based search.
type StringLiteral struct {
	// Pos is the position of the string literal.
	Pos token.Position

	// Value is the raw string value without quotes.
	Value string
}

// WriteStringTo writes the string representation of the value to the builder.
// If unquoted argument is set to true, the StringLiterals do not write its string
// representation surrounded with quotes.
func (s *StringLiteral) WriteStringTo(sb *strings.Builder, unquoted bool) {
	if unquoted {
		sb.WriteString(s.Value)
		return
	}

	sb.WriteRune('"')
	sb.WriteString(s.Value)
	sb.WriteRune('"')
}

// GetStringValue returns the string value.
func (s *StringLiteral) GetStringValue() string {
	return s.Value
}

// Position returns the position of the string literal.
func (s *StringLiteral) String() string {
	return fmt.Sprintf("%q", s.Value)
}

// UnquotedString returns the unquoted string.
func (s *StringLiteral) UnquotedString() string {
	return s.Value
}

func (s *StringLiteral) Position() token.Position { return s.Pos }
func (*StringLiteral) isFieldExpr()               {}
func (*StringLiteral) isValueExpr()               {}
func (*StringLiteral) isAstExpr()                 {}
