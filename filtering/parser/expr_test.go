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

const exprMultiWhiteSpace = "  a   AND     b"

func testExprMultiWhiteSpace(t *testing.T, pf *ParsedFilter) {
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
}

const complexExpr = "(a b) AND c OR d AND (e > f OR g < h)"

func testComplexExpr(t *testing.T, pf *ParsedFilter) {
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
}
