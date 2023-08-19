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

const funcCallNoArg = "msg.has_header()"

func testFuncCallNoArg(t *testing.T, pf *ParsedFilter) {
	// msg.has_header()
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
}

const complexFuncCall = "regex(m.key, '^.*prod.*$')"

func testComplexFuncCall(t *testing.T, pf *ParsedFilter) {
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
}
