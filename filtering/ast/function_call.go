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

var (
	// Compile-time check that *FunctionCall implements ComparableExpr.
	_ ComparableExpr = (*FunctionCall)(nil)
)

// FunctionCall is a function call expression
// which mau use simple or qualified names with zero or more arguments.
//
// function
//
//	: name {DOT name} LPAREN [argList] RPAREN
//	;
//
// FunctionCall implements ComparableExpr.
type FunctionCall struct {
	// Pos is the position of the function.
	Pos token.Position

	// Name is a list of function name expressions.
	Name []NameExpr

	// Lparen is the left parenthesis position.
	Lparen token.Position

	// ArgList is a list of argument expressions.
	ArgList *ArgListExpr

	// Rparen is the right parenthesis position.
	Rparen token.Position
}

// JoinedPkgName returns the joined package name of the function.
// I.e. if the function Name field contains a list of TEXT or Keyword
// expressions, it will return a string representation that merges them all with a dot (.)
// separator.
// The NameExpr cannot be a StringLiteral, thus no 'unquote' parameter is needed.
func (f *FunctionCall) JoinedPkgName() string {
	var sb strings.Builder
	for i := 0; i < len(f.Name)-1; i++ {
		if i > 0 {
			sb.WriteRune('.')
		}
		f.Name[i].WriteStringTo(&sb, false)
	}
	return sb.String()
}

// FuncName returns the function name.
// A valid function name is the last element of the Name list.
// I.e.: time.Unix(1234512512) -> pkgName: time, funcName: Unix
func (f *FunctionCall) FuncName() string {
	if len(f.Name) == 0 {
		return ""
	}
	return f.Name[len(f.Name)-1].String()
}

// JoinedPkgNameEquals returns true if the joined package name of the member equals the name.
// A valid package name is the joined name of len(Name)-1 elements.
// I.e.: time.Unix(1234512512) -> pkgName: time, funcName: Unix
func (f *FunctionCall) JoinedPkgNameEquals(pkgName string) bool {
	var sb strings.Builder
	for i := 0; i < len(f.Name)-1; i++ {
		if i > 0 {
			sb.WriteRune('.')
		}
		f.Name[i].WriteStringTo(&sb, false)
	}
	return sb.String() == pkgName
}

// FuncNameEquals returns true if the function name equals the name.
// A valid function name is the last element of the Name list.
// I.e.: time.Unix(1234512512) -> pkgName: time, funcName: Unix
func (f *FunctionCall) FuncNameEquals(name string) bool {
	if len(f.Name) == 0 {
		return false
	}
	return f.Name[len(f.Name)-1].String() == name
}

// JoinedNameEquals returns true if the joined name of the member equals the name.
func (f *FunctionCall) JoinedNameEquals(name string) bool {
	var sb strings.Builder
	for i, n := range f.Name {
		if i > 0 {
			sb.WriteRune('.')
		}
		n.WriteStringTo(&sb, false)
	}
	return sb.String() == name
}

// JoinedName returns the joined name of the function.
// I.e. if the function Name field contains a list of TEXT or Keyword
// expressions, it will return a string representation that merges them all with a dot (.)
// separator.
// The NameExpr cannot be a StringLiteral, thus no 'unquote' parameter is needed.
func (f *FunctionCall) JoinedName() string {
	var sb strings.Builder
	for i, n := range f.Name {
		if i > 0 {
			sb.WriteRune('.')
		}
		n.WriteStringTo(&sb, false)
	}
	return sb.String()
}

func (f *FunctionCall) UnquotedString() string {
	var sb strings.Builder
	for i, n := range f.Name {
		if i > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(n.UnquotedString())
	}
	sb.WriteRune('(')
	if f.ArgList != nil {
		sb.WriteString(f.ArgList.UnquotedString())
	}
	sb.WriteRune(')')
	return sb.String()
}

func (f *FunctionCall) String() string {
	var sb strings.Builder
	for i, n := range f.Name {
		if i > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(n.String())
	}
	sb.WriteRune('(')
	if f.ArgList != nil {
		sb.WriteString(f.ArgList.String())
	}
	sb.WriteRune(')')
	return sb.String()
}

// WriteStringTo writes the string representation of the value to the builder.
// If unquoted argument is set to true, the StringLiterals do not write its string
// representation surrounded with quotes.
func (f *FunctionCall) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, n := range f.Name {
		if i > 0 {
			sb.WriteString(".")
		}
		n.WriteStringTo(sb, unquoted)
	}
	sb.WriteRune('(')
	if f.ArgList != nil {
		f.ArgList.WriteStringTo(sb, unquoted)
	}
	sb.WriteRune(')')
}

func (f *FunctionCall) Position() token.Position { return f.Pos }
func (*FunctionCall) isComparableExpr()          {}
func (*FunctionCall) isArgExpr()                 {}
func (*FunctionCall) isAstExpr()                 {}
