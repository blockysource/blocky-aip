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

const arrayWithQuote = `["a","b"]`

func testArrayWithQuote(t *testing.T, pf *ParsedFilter) {
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

	res, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if res.Comparable == nil {
		t.Fatal("expected comparable")
	}

	array, ok := res.Comparable.(*ast.ArrayExpr)
	if !ok {
		t.Fatalf("expected array expression")
	}

	if array.LBracket != 0 {
		t.Errorf("expected lbracket 0 got: %v", array.LBracket)
	}

	if len(array.Elements) != 2 {
		t.Fatalf("expected two values got: %v", len(array.Elements))
	}

	m1, ok := array.Elements[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression")
	}

	if m1.Value == nil {
		t.Fatalf("expected member value")
	}

	if len(m1.Fields) > 0 {
		t.Fatalf("expected no fields got: %v", len(m1.Fields))
	}

	sl1, ok := m1.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", m1.Value)
	}

	if sl1.Value != "a" {
		t.Fatalf("expected 'a' got: %v", sl1.Value)
	}
	if sl1.Position() != 1 {
		t.Fatalf("expected position 1 got: %v", sl1.Position())
	}

	m2, ok := array.Elements[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression")
	}

	if m2.Value == nil {
		t.Fatalf("expected member value")
	}

	if len(m2.Fields) > 0 {
		t.Fatalf("expected no fields got: %v", len(m2.Fields))
	}

	sl2, ok := m2.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", m2.Value)
	}

	if sl2.Value != "b" {
		t.Fatalf("expected 'b' got: %v", sl2.Value)
	}

	if sl2.Position() != 5 {
		t.Fatalf("expected position 5 got: %v", sl2.Position())
	}

	if array.RBracket != 8 {
		t.Errorf("expected rbracket 8 got: %v", array.RBracket)
	}
}

const arrayWithQuoteAndWS = `["a", "b"]`

func testArrayWithQuoteAndWS(t *testing.T, pf *ParsedFilter) {
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

	res, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if res.Comparable == nil {
		t.Fatal("expected comparable")
	}

	array, ok := res.Comparable.(*ast.ArrayExpr)
	if !ok {
		t.Fatalf("expected array expression")
	}

	if array.LBracket != 0 {
		t.Errorf("expected lbracket 0 got: %v", array.LBracket)
	}

	if len(array.Elements) != 2 {
		t.Fatalf("expected two values got: %v", len(array.Elements))
	}

	m1, ok := array.Elements[0].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression")
	}

	if m1.Value == nil {
		t.Fatalf("expected member value")
	}

	if len(m1.Fields) > 0 {
		t.Fatalf("expected no fields got: %v", len(m1.Fields))
	}

	sl1, ok := m1.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", m1.Value)
	}

	if sl1.Value != "a" {
		t.Fatalf("expected 'a' got: %v", sl1.Value)
	}
	if sl1.Position() != 1 {
		t.Fatalf("expected position 1 got: %v", sl1.Position())
	}

	m2, ok := array.Elements[1].(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression")
	}

	if m2.Value == nil {
		t.Fatalf("expected member value")
	}

	if len(m2.Fields) > 0 {
		t.Fatalf("expected no fields got: %v", len(m2.Fields))
	}

	sl2, ok := m2.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", m2.Value)
	}

	if sl2.Value != "b" {
		t.Fatalf("expected 'b' got: %v", sl2.Value)
	}

	if sl2.Position() != 6 {
		t.Fatalf("expected position 5 got: %v", sl2.Position())
	}

	if array.RBracket != 9 {
		t.Errorf("expected rbracket 8 got: %v", array.RBracket)
	}
}
