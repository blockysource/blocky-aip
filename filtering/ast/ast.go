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

// Expr may either be a conjunction (AND) of sequences or a simple sequence.
type Expr struct {
	// Pos is the position of the expression.
	Pos token.Position

	// Sequences is a list of sequence expressions.
	Sequences []*SequenceExpr
}

// Position returns the position of the expression.
func (e *Expr) Position() token.Position { return e.Pos }

func (e *Expr) UnquotedString() string {
	var sb strings.Builder
	for i, seq := range e.Sequences {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(seq.UnquotedString())
	}
	return sb.String()
}

// String returns the string representation of the expression.
func (e *Expr) String() string {
	var sb strings.Builder
	for i, seq := range e.Sequences {
		if i > 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString(seq.String())
	}
	return sb.String()
}

// SequenceExpr is composed of one or more whitespace (WS) separated factors.
// A sequence expresses a logical relationship between 'factors' where
// the ranking of a filter result may be scored according to the number
// factors that match and other such criteria as the proximity of factors
// to each other within a document.
// When filters are used with exact match semantics rather than fuzzy
// match semantics, a sequence is equivalent to AND.
// Example: `New York Giants OR Yankees`
// The expression `New York (Giants OR Yankees)` is equivalent to the
// example.
type SequenceExpr struct {
	// Pos is the position of the sequence.
	Pos token.Position

	// Factors is a list of factor expressions separated by the WS operator.
	Factors []*FactorExpr

	// OpPos is the position of the AND operator after the sequence.
	// If there are no more sequences after this one, this value is undefined 0.
	OpPos token.Position
}

func (e *SequenceExpr) UnquotedString() string {
	var sb strings.Builder
	for i, f := range e.Factors {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(f.UnquotedString())
	}
	return sb.String()
}

func (e *SequenceExpr) Position() token.Position { return e.Pos }

func (e *SequenceExpr) String() string {
	var sb strings.Builder
	for i, f := range e.Factors {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(f.String())
	}
	return sb.String()
}

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

// TermExpr is a unary negation expression.
// The unary negation operator is a hyphen (-) that precedes a term or
// the NOT operator that precedes a term with a WS.
// TermExpr are separated in the FactorExpr by the OR operator.
// The position of the OR operator is stored in the OrOpPos field, in a TermExpr which
// stands before the operator, i.e.:
// 	`-foo OR bar` -> OrOpPos is the position of the OR operator after `-foo`.
//                   It is stored in the '-foo' TermExpr.OrOpPos field.
type TermExpr struct {
	Pos token.Position

	// UnaryOp is the unary operator.
	UnaryOp string

	// Expr is the expression.
	Expr SimpleExpr

	// OrOpPos is the position of the OR operator after the term.
	// If there are no more terms after this one, this value is undefined 0.
	OrOpPos token.Position
}

func (e *TermExpr) UnquotedString() string {
	switch e.UnaryOp {
	case "-":
		return fmt.Sprintf("-%s", e.Expr.UnquotedString())
	case "NOT":
		return fmt.Sprintf("NOT %s", e.Expr.UnquotedString())
	default:
		return e.Expr.UnquotedString()
	}
}

func (e *TermExpr) String() string {
	switch e.UnaryOp {
	case "-":
		return fmt.Sprintf("-%s", e.Expr.String())
	case "NOT":
		return fmt.Sprintf("NOT %s", e.Expr.String())
	default:
		return e.Expr.String()
	}
}

// SimpleExpr is a simple expression.
type SimpleExpr interface {
	UnquotedString() string
	String() string
	Position() token.Position
	isSimpleExpr()
}

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

// Position returns the position of the restriction.
func (r *RestrictionExpr) Position() token.Position { return r.Pos }
func (*RestrictionExpr) isSimpleExpr()              {}

// ComparableExpr is either a member or a function expression.
// Comparable may either be a member or function.
//
// 	EBNF:
//
// 	comparable
// 		: member
// 		| function
// 		;
type ComparableExpr interface {
	Position() token.Position
	UnquotedString() string
	String() string
	isComparableExpr()
	isArgExpr()
}

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
	// Compile-time check that *FunctionCall implements ComparableExpr.
	_ ComparableExpr = (*FunctionCall)(nil)

	// Compile-time check that *FunctionCall implements DecodedValueExpr.
	_ DecodedValueExpr = (*FunctionCall)(nil)
)

// FunctionCall is a function call expression
// which mau use simple or qualified names with zero or more arguments.
//
// function
//    : name {DOT name} LPAREN [argList] RPAREN
//    ;
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

	// DecodedValue is a decoded value of the function call.
	// NOTE: this is not related with direct AST parsing,
	// 	but is used to embed the decoded type into the AST.
	// I.e. as a result of executing a function call.
	DecodedValue any
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

// GetDecodedValue returns the decoded value, which may be a result of a function call.
func (f *FunctionCall) GetDecodedValue() (any, bool) {
	return f.DecodedValue, f.DecodedValue != undefinedDecodedValueInstance
}

// SetDecodedValue sets the decoded value.
func (f *FunctionCall) SetDecodedValue(v any) {
	f.DecodedValue = v
}

// SetUndefinedDecodedValue sets the undefined decoded value.
func (f *FunctionCall) SetUndefinedDecodedValue() {
	f.DecodedValue = undefinedDecodedValueInstance
}

// HasDecodedValue returns true if the function call has a decoded value.
// If the function call has a decoded value, it means that the function
// was executed and the result is stored in the DecodedValue field.
func (f *FunctionCall) HasDecodedValue() bool {
	return f.DecodedValue != undefinedDecodedValueInstance
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

func (f *FunctionCall) Position() token.Position { return f.Pos }
func (*FunctionCall) isComparableExpr()          {}
func (*FunctionCall) isArgExpr()                 {}

// ComparatorLiteral is the literal of a comparator with a position and type.
//
// EBNF:
//
// comparator
//    : LESS_EQUALS      # <=
//    | LESS_THAN        # <
//    | GREATER_EQUALS   # >=
//    | GREATER_THAN     # >
//    | NOT_EQUALS       # !=
//    | EQUALS           # =
//    | HAS              # :
//    ;
type ComparatorLiteral struct {
	Pos  token.Position
	Type ComparatorType
}

// Position returns the position of the comparator.
func (c *ComparatorLiteral) Position() token.Position { return c.Pos }

// String returns the string representation of the comparator.
func (c *ComparatorLiteral) String() string { return c.Type.String() }

var _ComparatorTypeStrings = [...]string{
	LE:  "<=",
	LT:  "<",
	GE:  ">=",
	GT:  ">",
	NE:  "!=",
	EQ:  "=",
	HAS: ":",
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
)

var (
	// Compile-time check that *CompositeExpr implements ArgExpr.
	_ ArgExpr = (*CompositeExpr)(nil)

	// Compile-time check that *CompositeExpr implements SimpleExpr.
	_ SimpleExpr = (*CompositeExpr)(nil)
)

// CompositeExpr is a composite expression with a left parenthesis, expression, and right parenthesis.
type CompositeExpr struct {
	Lparen token.Position
	Expr   *Expr
	Rparen token.Position
}

func (c *CompositeExpr) UnquotedString() string {
	if c.Expr == nil {
		return ""
	}
	return fmt.Sprintf("(%s)", c.Expr.UnquotedString())
}

func (c *CompositeExpr) Position() token.Position { return c.Lparen }
func (c *CompositeExpr) String() string {
	return fmt.Sprintf("(%s)", c.Expr)
}
func (*CompositeExpr) isArgExpr()    {}
func (*CompositeExpr) isSimpleExpr() {}
func (*CompositeExpr) isTermExpr()   {}

// NameExpr is a name expression which may be either
// a Text or KeyWord.
//
// EBNF:
//
// name
//    : TEXT
//    | KEYWORD
//    ;
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
}

// ValueExpr is a value expression.
// Value may either be a TEXT or STRING.
//
// EBNF:
//
// value
//    : TEXT
//    | STRING
//    ;
type ValueExpr interface {
	DecodedValueExpr

	// String returns the string representation of the value.
	String() string

	// UnquotedString returns the unquoted string, used for StringLiterals
	UnquotedString() string

	// Position returns the position of the value.
	Position() token.Position

	// GetStringValue returns the string value.
	GetStringValue() string

	// WriteStringTo writes the string representation of the value to the builder.
	// If unquoted argument is set to true, the StringLiterals do not write its string
	// representation surrounded with quotes.
	WriteStringTo(sb *strings.Builder, unquoted bool)

	// isValueExpr is a marker method for the interface.
	isValueExpr()
}

// DecodedValueExpr is an expression used for the reason of decoding
// the value of the expression.
// It is not directly related with the AST parsing, but is used
// to embed and share the decoded value between consumers.
type DecodedValueExpr interface {
	// HasDecodedValue returns true if the value has a decoded value.
	// NOTE: if DecodedValue is nil, it will return true.
	// 		This is a reason of decoding a 'null' value.
	// In order to reset the decoded value, use SetUndefinedDecodedValue().
	HasDecodedValue() bool

	// GetDecodedValue returns the decoded value.
	GetDecodedValue() (any, bool)

	// SetDecodedValue sets the decoded value.
	SetDecodedValue(v any)

	// SetUndefinedDecodedValue sets the undefined decoded value.
	SetUndefinedDecodedValue()
}

type undefinedDecodedValue struct{}

var undefinedDecodedValueInstance = &undefinedDecodedValue{}

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

	// DecodedValue is a decoded value parsed by the modifier.
	// NOTE: this is not related with direct AST parsing,
	// 	but is used to embed the decoded value and pass
	//  to other consumers.
	DecodedValue any
}

// WriteStringTo writes the string representation of the value to the builder.
func (t *TextLiteral) WriteStringTo(sb *strings.Builder, _ bool) {
	sb.WriteString(t.Value)
}

// GetStringValue returns the string value.
func (t *TextLiteral) GetStringValue() string {
	return t.Value
}

// HasDecodedValue returns true if the string literal has a decoded value.
func (t *TextLiteral) HasDecodedValue() bool {
	return t.DecodedValue != undefinedDecodedValueInstance
}

// GetDecodedValue returns the decoded value.
func (t *TextLiteral) GetDecodedValue() (any, bool) {
	return t.DecodedValue, t.DecodedValue != undefinedDecodedValueInstance
}

// SetDecodedValue sets the decoded value.
func (t *TextLiteral) SetDecodedValue(v any) {
	t.DecodedValue = v
}

// SetUndefinedDecodedValue sets the undefined decoded value.
func (t *TextLiteral) SetUndefinedDecodedValue() {
	t.DecodedValue = undefinedDecodedValueInstance
}

func (t *TextLiteral) UnquotedString() string   { return t.Value }
func (t *TextLiteral) String() string           { return t.Value }
func (t *TextLiteral) Position() token.Position { return t.Pos }
func (*TextLiteral) isNameExpr()                {}
func (*TextLiteral) isFieldExpr()               {}
func (*TextLiteral) isValueExpr()               {}

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

	// IsPrefixBased is true if the string is prefix-based search,
	// i.e. it starts with a wildcard. (e.g. "*foo")
	IsPrefixBased bool

	// IsSuffixBased is true if the string is suffix-based search.
	// i.e. it ends with a wildcard. (e.g. "foo*")
	IsSuffixBased bool

	// DecodedValue is a decoded value of the string literal.
	// NOTE: this is not related with direct AST parsing,
	// 	but is used to embed the decoded type into the AST.
	// This might be used for parsing enums, uuid, bytes, etc.
	DecodedValue any
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

// HasDecodedValue returns true if the string literal has a decoded value.
func (s *StringLiteral) HasDecodedValue() bool {
	return s.DecodedValue != undefinedDecodedValueInstance
}

// GetDecodedValue returns the decoded value.
func (s *StringLiteral) GetDecodedValue() (any, bool) {
	return s.DecodedValue, s.DecodedValue != undefinedDecodedValueInstance
}

// SetDecodedValue sets the decoded value.
func (s *StringLiteral) SetDecodedValue(v any) {
	s.DecodedValue = v
}

// SetUndefinedDecodedValue sets the undefined decoded value.
func (s *StringLiteral) SetUndefinedDecodedValue() {
	s.DecodedValue = undefinedDecodedValueInstance
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

// FieldExpr may be either a value or a keyword.
//
// EBNF:
//
// field
//    : value
//    | keyword
//    ;
type FieldExpr interface {
	// String returns the string representation of the field.
	String() string

	// UnquotedString returns the unquoted string.
	// If the field is a StringLiteral, it will return the unquoted string.
	UnquotedString() string

	// Position returns the position of the field.
	Position() token.Position

	// WriteStringTo writes the string representation of the value to the builder.
	// If unquoted argument is set to true, the StringLiterals do not write its string
	// representation surrounded with quotes.
	WriteStringTo(sb *strings.Builder, unquoted bool)

	// isFieldExpr is a marker method for the interface.
	isFieldExpr()
}

// ArgListExpr is a list of arguments to a function call.
//
// EBNF:
//
// argList
//    : arg { COMMA arg}
//    ;
type ArgListExpr struct {
	Args []ArgExpr
}

// Position returns the position of the first argument.
func (a *ArgListExpr) Position() token.Position {
	if len(a.Args) > 0 {
		return a.Args[0].Position()
	}
	return 0
}

func (a *ArgListExpr) String() string {
	var sb strings.Builder
	for i, arg := range a.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.String())
	}
	return sb.String()
}

func (a *ArgListExpr) UnquotedString() string {
	var sb strings.Builder
	for i, arg := range a.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.UnquotedString())
	}
	return sb.String()
}

// ArgExpr is either a Comparable or Composite
//
// EBNF:
//
// arg
//    : comparable
//    | composite
//    ;
type ArgExpr interface {
	String() string
	UnquotedString() string
	Position() token.Position
	isArgExpr()
}

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
