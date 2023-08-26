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

// NameExpr is a name expression which may be either
// a Text or KeyWord.
//
// EBNF:
//
// name
//
//	: TEXT
//	| KEYWORD
//	;
type NameExpr interface {
	// String returns the string representation of the name.
	String() string

	// UnquotedString returns the unquoted string, used for StringLiterals
	UnquotedString() string

	// Position returns the position of the name.
	Position() token.Position

	// WriteStringTo writes the string representation of the value to the builder.
	// If unquoted argument is set to true, the StringLiterals do not write its string
	// representation surrounded with quotes.
	WriteStringTo(sb *strings.Builder, unquoted bool)

	// isNameExpr is a marker method for the interface.
	isNameExpr()
	isAstExpr()
}
