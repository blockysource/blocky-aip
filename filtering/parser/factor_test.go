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

const factorsWithOR = `a OR b`

func testFactorsWithOR(t *testing.T, pf *ParsedFilter) {
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
}
