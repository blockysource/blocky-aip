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

var (
	_ NameExpr = (*KeywordExpr)(nil)
)

// KeywordExpr is a keyword expression.
type KeywordExpr struct {
	Pos token.Position
	Typ KeywordType
}

// WriteStringTo writes the string representation of the value to the builder.
func (k *KeywordExpr) WriteStringTo(sb *strings.Builder, _ bool) {
	sb.WriteString(k.Typ.String())
}

func (k *KeywordExpr) UnquotedString() string { return k.String() }

func (k *KeywordExpr) String() string {
	return k.Typ.String()
}

func (k *KeywordExpr) Position() token.Position { return k.Pos }
func (*KeywordExpr) isNameExpr()                {}
func (*KeywordExpr) isFieldExpr()               {}

// KeywordType is a keyword type enumeration.
type KeywordType int

const (
	_ KeywordType = iota
	// NOT is a keyword type that represents the NOT operator.
	NOT

	// AND is a keyword type that represents the AND operator.
	AND

	// OR is a keyword type that represents the OR operator.
	OR
)

func (k KeywordType) String() string {
	switch k {
	case NOT:
		return "NOT"
	case AND:
		return "AND"
	case OR:
		return "OR"
	}
	return ""
}

func (*KeywordExpr) isAstExpr() {}
