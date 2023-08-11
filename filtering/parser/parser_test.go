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
	"github.com/blockysource/blocky-aip/filtering/scanner"
	"github.com/blockysource/blocky-aip/filtering/token"
)

func memberTextLiteral(t *testing.T, m *ast.MemberExpr, expected string, pos int) {
	if m.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", m.Value)
	}

	if tl.Value != expected {
		t.Fatalf("expected '%s' got: %v", expected, tl.Value)
	}

	if tl.Pos != token.Position(pos) {
		t.Fatalf("expected position %d got: %v", pos, tl.Pos)
	}
}

// TestParse tests the Parse function.
func TestParse(t *testing.T) {
	seqMember := func(t *testing.T, seq *ast.SequenceExpr) *ast.MemberExpr {
		if len(seq.Factors) != 1 {
			t.Fatalf("expected one factor")
		}

		factor := seq.Factors[0]
		if len(factor.Terms) != 1 {
			t.Fatalf("expected one term")
		}

		term := factor.Terms[0]
		if term.UnaryOp != "" {
			t.Errorf("expected no unary op")
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
		return member
	}

	factorMember := func(t *testing.T, factor *ast.FactorExpr) *ast.MemberExpr {
		if len(factor.Terms) != 1 {
			t.Fatalf("expected one term")
		}

		term := factor.Terms[0]
		if term.UnaryOp != "" {
			t.Errorf("expected no unary op")
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
		return member
	}

	seqFuncCall := func(t *testing.T, seq *ast.SequenceExpr) *ast.FunctionCall {
		if len(seq.Factors) != 1 {
			t.Fatalf("expected one factor")
		}

		factor := seq.Factors[0]
		if len(factor.Terms) != 1 {
			t.Fatalf("expected one term")
		}

		term := factor.Terms[0]
		if term.UnaryOp != "" {
			t.Errorf("expected no unary op")
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

		fnCall, ok := expr.Comparable.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("expected function call, got: %T", expr.Comparable)
		}
		return fnCall
	}

	seqRestriction := func(t *testing.T, seq *ast.SequenceExpr) *ast.RestrictionExpr {
		if len(seq.Factors) != 1 {
			t.Fatalf("expected one factor")
		}

		factor := seq.Factors[0]
		if len(factor.Terms) != 1 {
			t.Fatalf("expected one term")
		}

		term := factor.Terms[0]
		if term.UnaryOp != "" {
			t.Errorf("expected no unary op")
		}

		if term.Expr == nil {
			t.Fatalf("expected expression")
		}

		expr, ok := term.Expr.(*ast.RestrictionExpr)
		if !ok {
			t.Fatalf("expected restriction expression")
		}

		return expr
	}

	testCases := []struct {
		name    string
		src     string
		checkFn func(t *testing.T, pf *ParsedFilter)
	}{
		{
			name: "empty",
			src:  "",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr != nil {
					t.Errorf("expected nil expression")
				}
			},
		},
		{
			name: "single sequence",
			src:  "a",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter")
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence")
				}

				seq := pf.Expr.Sequences[0]
				member := seqMember(t, seq)

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
			},
		},
		{
			name: "single sequence with string",
			src:  `"a"`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter")
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence")
				}

				seq := pf.Expr.Sequences[0]
				member := seqMember(t, seq)

				if member.Value == nil {
					t.Fatal("expected member value")
				}

				sl, ok := member.Value.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected string literal got: %v", member.Value)
				}

				if sl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", sl.Value)
				}
			},
		},
		{
			name: "single sequence with unary op",
			src:  "-a",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
					t.Fatalf("expected text literal got: %v", member.Value)
				}

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "single sequence with unary op and string",
			src:  `-"a"`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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

				tl, ok := member.Value.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", member.Value)
				}

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "single sequence with unary NOT op",
			src:  "NOT a",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
				if term.UnaryOp != "NOT" {
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
					t.Fatalf("expected text literal got: %v", member.Value)
				}

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "restriction with comparator",
			src:  "a = b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with ge comparator",
			src:  "a >= b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with ne comparator",
			src:  "a != b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			"restriction with no space and comparator",
			"a=b",
			func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with function arg",
			src:  "a = b(c)",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with function arg list",
			src:  "a = b(c, d)",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with function arg list no space",
			src:  "a = b(c,d)",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with has arg",
			src:  "m:foo",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
					t.Fatalf("expected 'a' got: %v", tl.Value)
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
			},
		},
		{
			name: "restriction with string arg",
			src:  `a = "foo"`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with single quoted string arg",
			src:  `a = 'foo'`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "restriction with alternative not has string arg",
			src:  `-file:".java"`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
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
			},
		},
		{
			name: "sequence with OR factor terms",
			src:  "a OR b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}

				seq := pf.Expr.Sequences[0]
				if len(seq.Factors) != 1 {
					t.Fatalf("expected one factor got: %v", len(seq.Factors))
				}

				factor := seq.Factors[0]
				if len(factor.Terms) != 2 {
					t.Fatalf("expected two terms got: %v", factor.Terms)
				}

				term := factor.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				if term.OrOpPos != 2 {
					t.Fatalf("expected position 2 got: %v", term.OrOpPos)
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

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}

				term2 := factor.Terms[1]
				if term2.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term2.Expr == nil {
					t.Fatalf("expected expression")
				}

				if term2.OrOpPos != 0 {
					t.Fatalf("expected position 0 got: %v", term2.OrOpPos)
				}

				expr, ok = term2.Expr.(*ast.RestrictionExpr)
				if !ok {
					t.Fatalf("expected restriction expression")
				}

				if expr.Comparable == nil {
					t.Fatal("expected comparable")
				}

				member, ok = expr.Comparable.(*ast.MemberExpr)
				if !ok {
					t.Fatalf("expected member literal")
				}

				if member.Value == nil {
					t.Fatal("expected member value")
				}

				tl, ok = member.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", member.Value)
				}

				if tl.Value != "b" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "sequence with factors",
			src:  "a b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}

				seq := pf.Expr.Sequences[0]
				if len(seq.Factors) != 2 {
					t.Fatalf("expected one factor got: %v", len(seq.Factors))
				}

				factor1 := seq.Factors[0]
				if len(factor1.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor1.Terms)
				}

				term := factor1.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
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

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}

				factor2 := seq.Factors[1]
				if len(factor2.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor2.Terms)
				}

				term = factor2.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				expr, ok = term.Expr.(*ast.RestrictionExpr)
				if !ok {
					t.Fatalf("expected restriction expression")
				}

				if expr.Comparable == nil {
					t.Fatal("expected comparable")
				}

				member, ok = expr.Comparable.(*ast.MemberExpr)
				if !ok {
					t.Fatalf("expected member literal")
				}

				if member.Value == nil {
					t.Fatal("expected member value")
				}

				tl, ok = member.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", member.Value)
				}

				if tl.Value != "b" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "composite expression",
			src:  "(a b)",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}
				seq := pf.Expr.Sequences[0]
				if len(seq.Factors) != 1 {
					t.Fatalf("expected one factor got: %v", len(seq.Factors))
				}

				factor := seq.Factors[0]
				if len(factor.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor.Terms)
				}

				term := factor.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				comp, ok := term.Expr.(*ast.CompositeExpr)
				if !ok {
					t.Fatalf("expected composite expression")
				}
				if comp.Expr == nil {
					t.Fatalf("expected expression")
				}
				seq = comp.Expr.Sequences[0]
				if len(seq.Factors) != 2 {
					t.Fatalf("expected one factor got: %v", len(seq.Factors))
				}

				factor1 := seq.Factors[0]
				if len(factor1.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor1.Terms)
				}

				term = factor1.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
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

				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}

				factor2 := seq.Factors[1]
				if len(factor2.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor2.Terms)
				}

				term = factor2.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				expr, ok = term.Expr.(*ast.RestrictionExpr)
				if !ok {
					t.Fatalf("expected restriction expression")
				}

				if expr.Comparable == nil {
					t.Fatal("expected comparable")
				}

				member, ok = expr.Comparable.(*ast.MemberExpr)
				if !ok {
					t.Fatalf("expected member literal")
				}

				if member.Value == nil {
					t.Fatal("expected member value")
				}

				tl, ok = member.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", member.Value)
				}

				if tl.Value != "b" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
			},
		},
		{
			name: "complex func call",
			src:  "regex(m.key, '^.*prod.*$')",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}

				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}
				seq := pf.Expr.Sequences[0]
				if len(seq.Factors) != 1 {
					t.Fatalf("expected one factor got: %v", len(seq.Factors))
				}

				factor := seq.Factors[0]
				if len(factor.Terms) != 1 {
					t.Fatalf("expected one terms got: %v", factor.Terms)
				}

				term := factor.Terms[0]
				if term.UnaryOp != "" {
					t.Errorf("expected no unary op")
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				rest, ok := term.Expr.(*ast.RestrictionExpr)
				if !ok {
					t.Fatalf("expected restriction expression")
				}
				if rest.Comparable == nil {
					t.Fatal("expected comparable")
				}

				if rest.Comparator != nil {
					t.Fatalf("expected no comparator")
				}

				if rest.Arg != nil {
					t.Fatalf("expected no arg")
				}

				fl, ok := rest.Comparable.(*ast.FunctionCall)
				if !ok {
					t.Fatalf("expected function literal")
				}

				if fl.Pos != 0 {
					t.Fatalf("expected position 0 got: %v", fl.Pos)
				}
				if fl.Lparen != 5 {
					t.Fatalf("expected position 5 got: %v", fl.Lparen)
				}
				if len(fl.Name) != 1 {
					t.Fatalf("expected one funciton name")
				}
				if fl.Rparen != 25 {
					t.Fatalf("expected position 25 got: %v", fl.Rparen)
				}
				name, ok := fl.Name[0].(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", fl.Name[0])
				}
				if name.Pos != 0 {
					t.Fatalf("expected position 0 got: %v", name.Pos)
				}
				if name.Value != "regex" {
					t.Fatalf("expected 'regex' got: %v", fl.Name)
				}

				if fl.ArgList == nil {
					t.Fatal("expected arg list")
				}

				if len(fl.ArgList.Args) != 2 {
					t.Fatalf("expected two arg")
				}
				aml, ok := fl.ArgList.Args[0].(*ast.MemberExpr)
				if !ok {
					t.Fatalf("expected arg member literal")
				}

				if aml.Value == nil {
					t.Fatal("expected arg member value")
				}
				if aml.Position() != 6 {
					t.Fatalf("expected position 6 got: %v", aml.Position())
				}

				tl, ok := aml.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", aml.Value)
				}
				if tl.Pos != 6 {
					t.Fatalf("expected position 6 got: %v", tl.Pos)
				}
				if tl.Value != "m" {
					t.Fatalf("expected 'm' got: %v", tl.Value)
				}
				if len(aml.Fields) != 1 {
					t.Fatalf("expected one field")
				}
				field := aml.Fields[0]
				if field.Position() != 8 {
					t.Fatalf("expected position 8 got: %v", field.Position())
				}
				fn, ok := field.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if fn.Value != "key" {
					t.Fatalf("expected 'key' got: %v", fn.Value)
				}

				aml, ok = fl.ArgList.Args[1].(*ast.MemberExpr)
				if !ok {
					t.Fatalf("expected arg member literal")
				}
				if aml.Value == nil {
					t.Fatal("expected arg member value")
				}
				if len(aml.Fields) != 0 {
					t.Fatalf("expected no fields")
				}
				if aml.Position() != 13 {
					t.Fatalf("expected position 12 got: %v", aml.Position())
				}
				sl, ok := aml.Value.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", aml.Value)
				}
				if sl.Pos != 13 {
					t.Fatalf("expected position 12 got: %v", sl.Pos)
				}
				if sl.Value != "^.*prod.*$" {
					t.Fatalf("expected '^.*prod.*$' got: '%v'", sl.Value)
				}
			},
		},
		{
			name: "deep nested member",
			src:  "expr.type_map.1.type",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}
				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}
				seq := pf.Expr.Sequences[0]
				member := seqMember(t, seq)
				if member.Value == nil {
					t.Fatal("expected member value")
				}
				if len(member.Fields) != 3 {
					t.Fatalf("expected three fields")
				}
				if member.Position() != 0 {
					t.Fatalf("expected position 0 got: %v", member.Position())
				}
				tl, ok := member.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", member.Value)
				}
				if tl.Pos != 0 {
					t.Fatalf("expected position 0 got: %v", tl.Pos)
				}
				if tl.Value != "expr" {
					t.Fatalf("expected 'expr' got: %v", tl.Value)
				}
				field := member.Fields[0]
				if field.Position() != 5 {
					t.Fatalf("expected position 5 got: %v", field.Position())
				}
				tl, ok = field.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if tl.Value != "type_map" {
					t.Fatalf("expected 'type_map' got: %v", tl.Value)
				}
				field = member.Fields[1]
				if field.Position() != 14 {
					t.Fatalf("expected position 14 got: %v", field.Position())
				}
				tl, ok = field.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if tl.Value != "1" {
					t.Fatalf("expected '1' got: %v", tl.Value)
				}
				field = member.Fields[2]
				if field.Position() != 16 {
					t.Fatalf("expected position 16 got: %v", field.Position())
				}
				tl, ok = field.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if tl.Value != "type" {
					t.Fatalf("expected 'type' got: %v", tl.Value)
				}
			},
		},
		{
			name: "complex",
			src:  "(a b) AND c OR d AND (e > f OR g < h)",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}

				if len(pf.Expr.Sequences) != 3 {
					t.Fatalf("expected three sequences got: %v", pf.Expr.Sequences)
				}

				seq := pf.Expr.Sequences[0]
				if len(seq.Factors) != 1 {
					t.Fatalf("expected two factors got: %v", len(seq.Factors))
				}

				f := seq.Factors[0]
				if len(f.Terms) != 1 {
					t.Fatalf("expected one term got: %v", len(f.Terms))
				}

				term := f.Terms[0]
				if term.UnaryOp != "" {
					t.Fatalf("expected no unary op got: %v", term.UnaryOp)
				}

				if term.Expr == nil {
					t.Fatalf("expected expression")
				}

				comp, ok := term.Expr.(*ast.CompositeExpr)
				if !ok {
					t.Fatalf("expected composite expression")
				}

				if comp.Expr == nil {
					t.Fatal("expected expression")
				}

				if len(comp.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", comp.Expr.Sequences)
				}

				seq = comp.Expr.Sequences[0]
				if len(seq.Factors) != 2 {
					t.Fatalf("expected two factors got: %v", len(seq.Factors))
				}

				m1 := factorMember(t, seq.Factors[0])
				memberTextLiteral(t, m1, "a", 1)
				m2 := factorMember(t, seq.Factors[1])
				memberTextLiteral(t, m2, "b", 3)

			},
		},
		{
			name: "multi whitespaces",
			src:  "  a   AND     b",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if len(pf.Expr.Sequences) != 2 {
					t.Fatalf("expected two sequences got: %v", pf.Expr.Sequences)
				}

				seq := pf.Expr.Sequences[0]

				m := seqMember(t, seq)
				if m.Value == nil {
					t.Fatal("expected member value")
				}

				tl, ok := m.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", m.Value)
				}
				if tl.Value != "a" {
					t.Fatalf("expected 'a' got: %v", tl.Value)
				}
				if tl.Pos != 2 {
					t.Fatalf("expected position 2 got: %v", tl.Pos)
				}
				if seq.OpPos != 6 {
					t.Fatalf("expected position 6 got: %v", seq.OpPos)
				}

				seq = pf.Expr.Sequences[1]
				m = seqMember(t, seq)
				if m.Value == nil {
					t.Fatal("expected member value")
				}

				tl, ok = m.Value.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", m.Value)
				}

				if tl.Value != "b" {
					t.Fatalf("expected 'b' got: %v", tl.Value)
				}

				if tl.Pos != 14 {
					t.Fatalf("expected position 12 got: %v", tl.Pos)
				}
			},
		},
		{
			name: "deep nested string member",
			src:  `"expr".'type_map'.1."type"`,
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}
				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}
				seq := pf.Expr.Sequences[0]
				member := seqMember(t, seq)
				if member.Value == nil {
					t.Fatal("expected member value")
				}
				if len(member.Fields) != 3 {
					t.Fatalf("expected three fields")
				}
				if member.Position() != 0 {
					t.Fatalf("expected position 0 got: %v", member.Position())
				}
				sl, ok := member.Value.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", member.Value)
				}
				if sl.Pos != 0 {
					t.Fatalf("expected position 0 got: %v", sl.Pos)
				}
				if sl.Value != "expr" {
					t.Fatalf("expected 'expr' got: %v", sl.Value)
				}
				field := member.Fields[0]
				if field.Position() != 7 {
					t.Fatalf("expected position 5 got: %v", field.Position())
				}
				sl, ok = field.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if sl.Value != "type_map" {
					t.Fatalf("expected 'type_map' got: %v", sl.Value)
				}
				field = member.Fields[1]
				if field.Position() != 18 {
					t.Fatalf("expected position 14 got: %v", field.Position())
				}
				tl, ok := field.(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if tl.Value != "1" {
					t.Fatalf("expected '1' got: %v", tl.Value)
				}
				field = member.Fields[2]
				if field.Position() != 20 {
					t.Fatalf("expected position 16 got: %v", field.Position())
				}
				sl, ok = field.(*ast.StringLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %v", field)
				}
				if sl.Value != "type" {
					t.Fatalf("expected 'type' got: %v", sl.Value)
				}
			},
		},
		{
			name: "func call no arg",
			src:  "msg.has_header()",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr == nil {
					t.Fatalf("expected parsed filter got: %v", pf)
				}
				if len(pf.Expr.Sequences) != 1 {
					t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
				}
				seq := pf.Expr.Sequences[0]
				fnCall := seqFuncCall(t, seq)

				if len(fnCall.Name) != 2 {
					t.Fatalf("expected two funciton name, got: %v", len(fnCall.Name))
				}

				first, ok := fnCall.Name[0].(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", fnCall.Name[0])
				}

				if first.Value != "msg" {
					t.Fatalf("expected 'msg' got: %v", first.Value)
				}
				if first.Pos != 0 {
					t.Fatalf("expected position 0 got: %v", first.Pos)
				}

				second, ok := fnCall.Name[1].(*ast.TextLiteral)
				if !ok {
					t.Fatalf("expected text literal got: %T", fnCall.Name[1])
				}

				if second.Value != "has_header" {
					t.Fatalf("expected 'has_header' got: %v", second.Value)
				}

				if second.Pos != 4 {
					t.Fatalf("expected position 4 got: %v", second.Pos)
				}
				if fnCall.Lparen != 14 {
					t.Fatalf("expected position 12 got: %v", fnCall.Lparen)
				}
				if fnCall.Rparen != 15 {
					t.Fatalf("expected position 13 got: %v", fnCall.Rparen)
				}
				if fnCall.ArgList != nil {
					t.Fatalf("expected no arg list")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.src, ErrorHandlerOption(testErrHandler(t)))

			pf, err := p.Parse()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			defer pf.Free()

			tc.checkFn(t, pf)
		})
	}
}

func testErrHandler(t testing.TB) scanner.ErrorHandler {
	return func(pos token.Position, msg string) {
		t.Errorf("unexpected error at %d: %s", pos, msg)
	}
}

func BenchmarkParse(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		p := Parser{}
		for i := 0; i < b.N; i++ {
			p.Reset("a")
			pf, err := p.Parse()
			if err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			pf.Free()
		}
	})

	b.Run("Complex", func(b *testing.B) {
		p := Parser{}
		for i := 0; i < b.N; i++ {
			p.Reset("(a b) AND c OR d AND (e > f OR g < h)")
			pf, err := p.Parse()
			if err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			pf.Free()
		}
	})
}
