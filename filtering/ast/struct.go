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

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/filtering/token"
)

// Compile-time check for StructExpr implementing ComparableExpr.
var _ ComparableExpr = (*StructExpr)(nil)

// StructExpr is an extension to the standard AST expressions to support structs.
// It is used to represent a message creation expression.
// By default, the parser is not aware of structs and will not parse them.
// To enable struct parsing, the parser must be configured with the struct definitions.
type StructExpr struct {
	// Name is the name of the struct.
	Name []NameExpr

	// LBrace is the position of the left brace.
	LBrace token.Position

	// Elements is a list of fields.
	Elements []*StructFieldExpr

	// RBrace is the position of the right brace.
	RBrace token.Position
}

// FullName returns a protoreflect.FullName equivalent of the struct expression.
// I.e.: "google.protobuf.Timestamp{}" -> "google.protobuf.Timestamp".
func (e *StructExpr) FullName() protoreflect.FullName {
	if len(e.Name) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, name := range e.Name {
		if i > 0 {
			sb.WriteString(".")
		}
		sb.WriteString(name.String())
	}
	return protoreflect.FullName(sb.String())
}

// IsMap returns true if the struct expression is a map.
func (e *StructExpr) IsMap() bool {
	if len(e.Name) != 1 {
		return false
	}

	return e.Name[0].String() == "map"
}

// Position returns the position of the struct expression.
func (e *StructExpr) Position() token.Position {
	if len(e.Name) > 0 {
		return e.Name[0].Position()
	}
	return e.LBrace
}

// UnquotedString returns the unquoted string representation of the struct expression.
func (e *StructExpr) UnquotedString() string {
	var sb strings.Builder
	e.WriteStringTo(&sb, true)
	return sb.String()
}

// String returns the string representation of the struct expression.
func (e *StructExpr) String() string {
	var sb strings.Builder
	e.WriteStringTo(&sb, false)
	return sb.String()
}

// WriteStringTo writes the string representation of the struct expression to the given builder.
func (e *StructExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, name := range e.Name {
		name.WriteStringTo(sb, unquoted)
		if i < len(e.Name)-1 {
			sb.WriteString(".")
		}
	}

	sb.WriteString("{")

	for i, field := range e.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		field.WriteStringTo(sb, unquoted)
	}

	sb.WriteString("}")
}

func (e *StructExpr) isAstExpr()        {}
func (e *StructExpr) isComparableExpr() {}
func (e *StructExpr) isArgExpr()        {}

// StructFieldExpr is an extension to the standard AST expressions to support structs.
// It is used to represent a field in a message creation expression.
// By default, the parser is not aware of structs and will not parse them.
type StructFieldExpr struct {
	// Name is the required name of the field.
	Name []ValueExpr

	// Colon is the position of the colon.
	Colon token.Position

	// Value is the value of the field.
	Value ComparableExpr
}

// Position returns the position of the struct field expression.
func (e *StructFieldExpr) Position() token.Position {
	if len(e.Name) > 0 {
		return e.Name[0].Position()
	}
	return 0
}

func (e *StructFieldExpr) UnquotedString() string {
	var sb strings.Builder
	e.WriteStringTo(&sb, true)
	return sb.String()
}

func (e *StructFieldExpr) String() string {
	var sb strings.Builder
	e.WriteStringTo(&sb, false)
	return sb.String()
}

func (e *StructFieldExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, name := range e.Name {
		name.WriteStringTo(sb, unquoted)
		if i < len(e.Name)-1 {
			sb.WriteString(".")
		}
	}
	sb.WriteString(": ")
	e.Value.WriteStringTo(sb, unquoted)
}

func (e *StructFieldExpr) isAstExpr() {}
