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

package parser

import (
	"testing"

	"github.com/blockysource/blocky-aip/filtering/ast"
)

const restrictionWithEQ = "a = b"

func testRestrictionWithEQ(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}
}

const restrictionWithGE = "a >= b"

func testRestrictionWithGE(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.GE {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}
}

const restrictionWithNE = "a != b"

func testRestrictionWithNE(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.NE {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}
}

const restrictionWithNoSpaceAndComparator = "a=b"

func testRestrictionWithNoSpaceAndComparator(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}
	if rest.Comparator.Pos != 1 {
		t.Fatalf("expected position 1 got: %v", rest.Comparator.Pos)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}
	if aml.Position() != 2 {
		t.Fatalf("expected position 2 got: %v", aml.Position())
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}
}

const restrictionWithFunctionArgListNoSpace = "a = b(c,d)"

func testRestrictionWithFunctionArgListNoSpace(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	fl, ok := rest.Arg.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected function literal, got: %T", rest.Arg)
	}

	if fl.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", fl.Pos)
	}
	if fl.Lparen != 5 {
		t.Fatalf("expected position 6 got: %v", fl.Lparen)
	}
	if len(fl.Name) != 1 {
		t.Fatalf("expected one funciton name")
	}
	name, ok := fl.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", fl.Name[0])
	}
	if name.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", name.Pos)
	}
	if name.Value != "b" {
		t.Fatalf("expected 'b' got: %v", fl.Name)
	}

	if fl.ArgList == nil {
		t.Fatal("expected arg list")
	}

	if len(fl.ArgList.Args) != 2 {
		t.Fatalf("expected one arg")
	}
	aml, ok := fl.ArgList.Args[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}
	if tl.Pos != 6 {
		t.Fatalf("expected position 6 got: %v", tl.Pos)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}

	aml, ok = fl.ArgList.Args[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Pos != 8 {
		t.Fatalf("expected position 8 got: %v", tl.Pos)
	}

	if tl.Value != "d" {
		t.Fatalf("expected 'd' got: %v", tl.Value)
	}

	if fl.Rparen != 9 {
		t.Fatalf("expected position 9 got: %v", fl.Rparen)
	}
}

const restrictionWithHasArg = "m:foo"

func testRestrictionWithHasArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "m" {
		t.Fatalf("expected 'm' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.HAS {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}
	if rest.Comparator.Pos != 1 {
		t.Fatalf("expected position 1 got: %v", rest.Comparator.Pos)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "foo" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}
}

const restrictionWithStringArg = `a = "foo"`

func testRestrictionWithStringArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	sl, ok := aml.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %v", aml.Value)
	}

	if sl.Value != "foo" {
		t.Fatalf("expected 'b' got: %v", sl.Value)
	}
}

const restrictionWithSingleQuotedStringArg = `a = 'foo'`

func testRestrictionWithSingleQuotedStringArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	aml, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	sl, ok := aml.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %v", aml.Value)
	}

	if sl.Value != "foo" {
		t.Fatalf("expected 'b' got: %v", sl.Value)
	}
}

const restrictionWithAlternativeNotHastStringArg = `-file:".java"`

func testRestrictionWithAlternativeNotHastStringArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor")
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one term")
	}

	term := factor.Terms[0]
	if term.UnaryOp != "-" {
		t.Errorf("expected unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	member, ok := expr.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", member.Value)
	}

	if tl.Value != "file" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if expr.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if expr.Comparator.Type != ast.HAS {
		t.Fatalf("expected ':' got: %v", expr.Comparator.Type)
	}

	if expr.Comparator.Pos != 5 {
		t.Fatalf("expected position 5 got: %v", expr.Comparator.Pos)
	}

	if expr.Arg == nil {
		t.Fatal("expected arg")
	}

	member, ok = expr.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	sl, ok := member.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %v", member.Value)
	}

	if sl.Value != ".java" {
		t.Fatalf("expected '.java' got: %v", sl.Value)
	}
}

const restrictionWithFunctionArgList = "a = b(c, d)"

func testRestrictionWithFunctionArgList(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	fl, ok := rest.Arg.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected function literal, got: %T", rest.Arg)
	}

	if fl.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", fl.Pos)
	}
	if fl.Lparen != 5 {
		t.Fatalf("expected position 6 got: %v", fl.Lparen)
	}
	if len(fl.Name) != 1 {
		t.Fatalf("expected one funciton name")
	}
	name, ok := fl.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", fl.Name[0])
	}
	if name.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", name.Pos)
	}
	if name.Value != "b" {
		t.Fatalf("expected 'b' got: %v", fl.Name)
	}

	if fl.ArgList == nil {
		t.Fatal("expected arg list")
	}

	if len(fl.ArgList.Args) != 2 {
		t.Fatalf("expected one arg")
	}
	aml, ok := fl.ArgList.Args[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}
	if tl.Pos != 6 {
		t.Fatalf("expected position 6 got: %v", tl.Pos)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}

	aml, ok = fl.ArgList.Args[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Pos != 9 {
		t.Fatalf("expected position 8 got: %v", tl.Pos)
	}

	if tl.Value != "d" {
		t.Fatalf("expected 'd' got: %v", tl.Value)
	}

	if fl.Rparen != 10 {
		t.Fatalf("expected position 9 got: %v", fl.Rparen)
	}
}

const restrictionWithFunctionArg = "a = b(c)"

func testRestrictionWithFunctionArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	fl, ok := rest.Arg.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected function literal, got: %T", rest.Arg)
	}

	if fl.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", fl.Pos)
	}
	if fl.Lparen != 5 {
		t.Fatalf("expected position 6 got: %v", fl.Lparen)
	}
	if len(fl.Name) != 1 {
		t.Fatalf("expected one funciton name")
	}
	name, ok := fl.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", fl.Name[0])
	}
	if name.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", name.Pos)
	}
	if name.Value != "b" {
		t.Fatalf("expected 'b' got: %v", fl.Name)
	}

	if fl.ArgList == nil {
		t.Fatal("expected arg list")
	}

	if len(fl.ArgList.Args) != 1 {
		t.Fatalf("expected one arg")
	}
	aml, ok := fl.ArgList.Args[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}
}

const restrictionWithStructArg = `a = Foo{"a": "b", c: d}`

func testRestrictionWithStructArg(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	st, ok := rest.Arg.(*ast.StructExpr)
	if !ok {
		t.Fatalf("expected struct expr got: %T", rest.Arg)
	}

	if st.Position() != 4 {
		t.Fatalf("expected position 4 got: %v", st.Position())
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f1 := st.Elements[0]
	if f1.Name == nil {
		t.Fatalf("expected field name")
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one field name, got: %v - %s", len(f1.Name), f1.Name)
	}

	sl, ok := f1.Name[0].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", f1.Name)
	}

	if sl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", sl.Value)
	}

	if f1.Value == nil {
		t.Fatalf("expected field value")
	}

	m, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", f1.Value)
	}

	if m.Value == nil {
		t.Fatalf("expected member value")
	}

	sl, ok = m.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", m.Value)
	}

	if sl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", sl.Value)
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f2 := st.Elements[1]
	if f2.Name == nil {
		t.Fatalf("expected field name")
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one field name")
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}

	if f2.Value == nil {
		t.Fatalf("expected field value")
	}

	m, ok = f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", f2.Value)
	}

	if m.Value == nil {
		t.Fatalf("expected member value")
	}

	tl, ok = m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if tl.Value != "d" {
		t.Fatalf("expected 'd' got: %v", tl.Value)
	}
}

const complexRestrictionWithFuncCallStructAndArray = "a = b(c, d) AND e = Foo{a: b, c: d} AND f = [a, b, c]"

func testComplexRestrictionWithFuncCallStructAndArray(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 3 {
		t.Fatalf("expected three sequences, got: %v", len(pf.Expr.Sequences))
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	fl, ok := rest.Arg.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected function literal, got: %T", rest.Arg)
	}

	if len(fl.Name) != 1 {
		t.Fatalf("expected one function name")
	}

	name, ok := fl.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", fl.Name[0])
	}
	if name.Value != "b" {
		t.Fatalf("expected 'b' got: %v", fl.Name)
	}

	if fl.ArgList == nil {
		t.Fatal("expected arg list")
	}

	if len(fl.ArgList.Args) != 2 {
		t.Fatalf("expected two args")
	}

	aml, ok := fl.ArgList.Args[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}

	aml, ok = fl.ArgList.Args[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected arg member literal")
	}

	if aml.Value == nil {
		t.Fatal("expected arg member value")
	}

	tl, ok = aml.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", aml.Value)
	}

	if tl.Value != "d" {
		t.Fatalf("expected 'd' got: %v", tl.Value)
	}

	seq = pf.Expr.Sequences[1]
	rest = seqRestriction(t, seq)

	member, ok = rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "e" {
		t.Fatalf("expected 'e' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	st, ok := rest.Arg.(*ast.StructExpr)
	if !ok {
		t.Fatalf("expected struct expr got: %T", rest.Arg)
	}

	if len(st.Name) != 1 {
		t.Fatalf("expected one struct name")
	}

	name, ok = st.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", st.Name[0])
	}

	if name.Value != "Foo" {
		t.Fatalf("expected 'Foo' got: %v", name.Value)
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f1 := st.Elements[0]

	if f1.Name == nil {
		t.Fatalf("expected field name")
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one field name")
	}
	tl, ok = f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if f1.Value == nil {
		t.Fatalf("expected field value")
	}

	m, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", f1.Value)
	}

	if m.Value == nil {
		t.Fatalf("expected member value")
	}

	tl, ok = m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}

	f2 := st.Elements[1]
	if f2.Name == nil {
		t.Fatalf("expected field name")
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one field name")
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}

	if f2.Value == nil {
		t.Fatalf("expected field value")
	}

	m, ok = f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", f2.Value)
	}

	if m.Value == nil {
		t.Fatalf("expected member value")
	}

	tl, ok = m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if tl.Value != "d" {
		t.Fatalf("expected 'd' got: %v", tl.Value)
	}

	seq = pf.Expr.Sequences[2]
	rest = seqRestriction(t, seq)

	member, ok = rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "f" {
		t.Fatalf("expected 'f' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	arr, ok := rest.Arg.(*ast.ArrayExpr)
	if !ok {
		t.Fatalf("expected array expr got: %T", rest.Arg)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("expected three values got: %v", len(arr.Elements))
	}

	member, ok = arr.Elements[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	member, ok = arr.Elements[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}

	member, ok = arr.Elements[2].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}
}

const restrictionWithIN = "a IN [b, c]"

func testRestrictionWithIN(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}

	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}

	if rest.Comparator.Type != ast.IN {
		t.Fatalf("expected 'IN' got: %v", rest.Comparator)
	}

	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	arr, ok := rest.Arg.(*ast.ArrayExpr)
	if !ok {
		t.Fatalf("expected array expr got: %T", rest.Arg)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("expected two values got: %v", len(arr.Elements))
	}

	member, ok = arr.Elements[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}

	member, ok = arr.Elements[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}

	if member.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}

	if tl.Value != "c" {
		t.Fatalf("expected 'c' got: %v", tl.Value)
	}
}

const restrictionWithTimestamp = "a = 2018-01-01T00:00:00Z"

func testRestrictionWithTimestamp(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}
	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}
	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}
	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}
	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}
	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}
	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}
	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	me, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", rest.Arg)
	}
	if me.Value == nil {
		t.Fatal("expected text timestamp literal")
	}

	tl, ok = me.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", me.Value)
	}

	if tl.Value != "2018-01-01T00:00:00Z" {
		t.Fatalf("expected '2018-01-01T00:00:00Z' got: %v", tl.Value)
	}
	if !tl.IsTimestamp {
		t.Fatalf("expected timestamp got: %v", tl.IsTimestamp)
	}
}

const restrictionWithTimestampAndTimezone = "a = 2018-01-01T00:00:00+01:00"

func testRestrictionWithTimestampAndTimezone(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}
	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}
	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}
	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}
	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}
	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}
	if rest.Comparator.Type != ast.EQ {
		t.Fatalf("expected '=' got: %v", rest.Comparator)
	}
	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	me, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", rest.Arg)
	}
	if me.Value == nil {
		t.Fatal("expected text timestamp literal")
	}

	tl, ok = me.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", me.Value)
	}

	if tl.Value != "2018-01-01T00:00:00+01:00" {
		t.Fatalf("expected '2018-01-01T00:00:00+01:00' got: %v", tl.Value)
	}
	if !tl.IsTimestamp {
		t.Fatalf("expected timestamp got: %v", tl.IsTimestamp)
	}
}

const restrictionWithTimestampAndHas = "2018-01-01T00:00:00Z:a"

func testRestrictionWithTimestampAndHas(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter")
	}
	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence")
	}
	seq := pf.Expr.Sequences[0]
	rest := seqRestriction(t, seq)
	member, ok := rest.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	if member.Value == nil {
		t.Fatal("expected member value")
	}
	tl, ok := member.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", member.Value)
	}
	if tl.Value != "2018-01-01T00:00:00Z" {
		t.Fatalf("expected '2018-01-01T00:00:00Z' got: %v", tl.Value)
	}
	if !tl.IsTimestamp {
		t.Fatalf("expected timestamp got: %v", tl.IsTimestamp)
	}
	if rest.Comparator == nil {
		t.Fatal("expected comparator")
	}
	if rest.Comparator.Type != ast.HAS {
		t.Fatalf("expected 'HAS' got: %v", rest.Comparator)
	}
	if rest.Arg == nil {
		t.Fatal("expected arg")
	}

	me, ok := rest.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expr got: %T", rest.Arg)
	}
	if me.Value == nil {
		t.Fatal("expected text literal")
	}

	tl, ok = me.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", me.Value)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}
}
