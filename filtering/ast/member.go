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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blockysource/blocky-aip/filtering/token"
)

// Compile-time check that *MemberExpr implements ComparableExpr.
var _ ComparableExpr = (*MemberExpr)(nil)

// MemberExpr is a member expression which either is a value or
// DOR qualified field references.
//
// EBNF:
//
// member
//    : value {DOT field}
//    ;
//
// MemberExpr implements ComparableExpr.
type MemberExpr struct {
	// Value is the value expression.
	Value ValueExpr

	// Fields is a list of field expressions, DOT separated.
	Fields []FieldExpr
}

// JoinedNameEquals returns true if the joined name of the member equals the name.
func (m *MemberExpr) JoinedNameEquals(name string, unquoted bool) bool {
	var sb strings.Builder

	if m.Value != nil {
		m.Value.WriteStringTo(&sb, unquoted)
		sb.WriteRune('.')
	}

	for i, f := range m.Fields {
		if i > 0 {
			sb.WriteRune('.')
		}
		f.WriteStringTo(&sb, unquoted)
	}

	return sb.String() == name
}

// JoinedName returns a result of joining the value and fields,
// with a dot (.) separator.
// If the unquote parameter is true, it will return the unquoted string for
// the StringLiteral values, otherwise a string literal is surrounded by double quotes (").
func (m *MemberExpr) JoinedName(unquote bool) string {
	var sb strings.Builder

	if m.Value != nil {
		m.Value.WriteStringTo(&sb, unquote)
		sb.WriteRune('.')
	}

	for i, f := range m.Fields {
		if i > 0 {
			sb.WriteRune('.')
		}
		f.WriteStringTo(&sb, unquote)
	}
	return sb.String()
}

// ClearFields clears the fields.
func (m *MemberExpr) ClearFields() {
	if m.Fields == nil {
		return
	}
	m.Fields = m.Fields[:0]
}

func (m *MemberExpr) String() string {
	if m.Value == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(m.Value.String())
	for _, f := range m.Fields {
		sb.WriteRune('.')
		sb.WriteString(f.String())
	}
	return sb.String()
}

func (m *MemberExpr) UnquotedString() string {
	if m.Value == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(m.Value.UnquotedString())
	for _, f := range m.Fields {
		sb.WriteRune('.')
		sb.WriteString(f.UnquotedString())
	}
	return sb.String()
}

// WriteStringTo writes the string representation of the value to the builder.
// If unquoted argument is set to true, the StringLiterals do not write its string
// representation surrounded with quotes.
func (m *MemberExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	if m.Value == nil {
		return
	}
	m.Value.WriteStringTo(sb, unquoted)
	for _, f := range m.Fields {
		sb.WriteRune('.')
		f.WriteStringTo(sb, unquoted)
	}
}

// Position returns the position of the member.
func (m *MemberExpr) Position() token.Position {
	if m.Value == nil {
		return 0
	}
	return m.Value.Position()
}
func (*MemberExpr) isComparableExpr() {}
func (*MemberExpr) isArgExpr()        {}

var (
	ErrMemberNotFloat    = errors.New("member expression is not a float64")
	ErrMemberNotDuration = errors.New("member expression is not a time.Duration")
	ErrMemberNotTime     = errors.New("member expression is not a time.Time")
)

// DecodeFloat64 tries to decode a member expression as a float64 value.
// This can be satisfied in two cases:
// 1. The member expression contains only a Value expression which is an integer,
//    that type can be converted to float64.
// 2. The member expression contains a Value expression along with a single field
//    which is a TextLiteral and along with the Value expression forms a float64 value.
func (m *MemberExpr) DecodeFloat64() (float64, error) {
	if len(m.Fields) == 0 {
		v, ok := m.Value.(*TextLiteral)
		if !ok {
			return 0, ErrMemberNotFloat
		}

		f, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	}

	if len(m.Fields) != 1 {
		return 0, ErrMemberNotFloat
	}

	v, ok := m.Value.(*TextLiteral)
	if !ok {
		return 0, ErrMemberNotFloat
	}

	f, ok := m.Fields[0].(*TextLiteral)
	if !ok {
		return 0, ErrMemberNotFloat
	}

	fl, err := strconv.ParseFloat(fmt.Sprintf("%s.%s", v.Value, f.Value), 64)
	if err != nil {
		return 0, err
	}
	return fl, nil
}

// DecodeDuration tries to decode a member expression as a time.Duration value.
// This can be satisfied in two cases:
// 1. The member expression contains only a Value expression which is an integer with a 's', 'm', 'h', 'd' etc. suffix,
//    that type can be converted to time.Duration.
// 2. The member expression contains a Value expression along with a single field
//    which is a TextLiteral and along with the Value expression (fractional) forms a time.Duration value.
func (m *MemberExpr) DecodeDuration() (time.Duration, error) {
	if len(m.Fields) == 0 {
		v, ok := m.Value.(*TextLiteral)
		if !ok {
			return 0, ErrMemberNotDuration
		}

		d, err := time.ParseDuration(v.Value)
		if err != nil {
			return 0, err
		}
		return d, nil
	}

	if len(m.Fields) != 1 {
		return 0, ErrMemberNotDuration
	}

	v, ok := m.Value.(*TextLiteral)
	if !ok {
		return 0, ErrMemberNotDuration
	}

	f, ok := m.Fields[0].(*TextLiteral)
	if !ok {
		return 0, ErrMemberNotDuration
	}

	d, err := time.ParseDuration(fmt.Sprintf("%s.%s", v.Value, f.Value))
	if err != nil {
		return 0, err
	}
	return d, nil
}

func (*MemberExpr) isAstExpr() {}
