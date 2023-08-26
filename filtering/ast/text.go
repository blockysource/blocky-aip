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
	"strconv"
	"strings"

	"github.com/blockysource/blocky-aip/token"
)

// Compile-time checks that *TextLiteral implements ValueExpr, FieldExpr, and NameExpr.
var (
	_ ValueExpr = (*TextLiteral)(nil)
	_ FieldExpr = (*TextLiteral)(nil)
	_ NameExpr  = (*TextLiteral)(nil)
)

// TextLiteral is a string literal
// TEXT is a free-form set of characters without whitespace (WS)
// or . (DOT) within it. The text may represent a variable, string,
// number, boolean, or alternative literal value and must be handled
// in a manner consistent with the service's intention.
type TextLiteral struct {
	// Pos is the position of the text literal.
	Pos token.Position

	// Value is a raw value of the text literal.
	Value string

	// Token is the token type of the text literal.
	Token token.Token
}

// WriteStringTo writes the string representation of the value to the builder.
func (t *TextLiteral) WriteStringTo(sb *strings.Builder, _ bool) {
	sb.WriteString(t.Value)
}

// GetStringValue returns the string value.
func (t *TextLiteral) GetStringValue() string {
	return t.Value
}

// DecodeInt64 decodes the literal text as an int64 value.
func (t *TextLiteral) DecodeInt64() (int64, error) {
	return strconv.ParseInt(t.Value, 10, 64)
}

// DecodeUint64 decodes the literal text as an uint64 value.
func (t *TextLiteral) DecodeUint64() (uint64, error) {
	return strconv.ParseUint(t.Value, 10, 64)
}

func (t *TextLiteral) UnquotedString() string   { return t.Value }
func (t *TextLiteral) String() string           { return t.Value }
func (t *TextLiteral) Position() token.Position { return t.Pos }
func (*TextLiteral) isNameExpr()                {}
func (*TextLiteral) isFieldExpr()               {}
func (*TextLiteral) isValueExpr()               {}
func (*TextLiteral) isAstExpr()                 {}
func (*TextLiteral) isArgExpr()                 {}
