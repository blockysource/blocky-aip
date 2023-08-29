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
	"github.com/blockysource/blocky-aip/token"
)

const tstNotInAsFieldNames = `NOT.IN = "value"`

func testNotInAsFieldNames(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter got: %v", pf)
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
	}

	seq := pf.Expr.Sequences[0]
	rs := seqRestriction(t, seq)

	// NOT.IN = "value"
	m, ok := rs.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", rs.Comparable)
	}

	mv, ok := m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if mv.Value != "NOT" {
		t.Fatalf("expected 'NOT' got: %v", mv.Value)
	}

	if mv.Token != token.NOT {
		t.Fatalf("expected 'NOT' got: %v", mv.Token)
	}

	if mv.Pos != 0 {
		t.Fatalf("expected position 0 got: %v", mv.Pos)
	}

	if len(m.Fields) != 1 {
		t.Fatalf("expected one field got: %v", m.Fields)
	}

	mf, ok := m.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Fields[0])
	}

	if mf.Value != "IN" {
		t.Fatalf("expected 'IN' got: %v", mf.Value)
	}

	if mf.Token != token.IN {
		t.Fatalf("expected 'IN' got: %v", mf.Token)
	}

	if mf.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", mf.Pos)
	}

	am, ok := rs.Arg.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", rs.Arg)
	}

	sm, ok := am.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected string literal got: %T", am.Value)
	}

	if sm.Value != "value" {
		t.Fatalf("expected 'value' got: %v", sm.Value)
	}
}

const tstNotAnd = "NOT AND"

func testNotAnd(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter got: %v", pf)
	}

	if len(pf.Expr.Sequences) != 1 {
		t.Fatalf("expected one sequence got: %v", pf.Expr.Sequences)
	}

	seq := pf.Expr.Sequences[0]

	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor got: %v", seq.Factors)
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one terms got: %v", factor.Terms)
	}

	tm := factor.Terms[0]
	if tm.UnaryOp != "NOT" {
		t.Fatalf("expected 'NOT' got: %v", tm.UnaryOp)
	}

	if tm.Expr == nil {
		t.Fatalf("expected expression")
	}

	rs, ok := tm.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	// NOT AND
	m, ok := rs.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", rs.Comparable)
	}

	mv, ok := m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if mv.Value != "AND" {
		t.Fatalf("expected 'AND' got: %v", mv.Value)
	}

	if mv.Token != token.AND {
		t.Fatalf("expected 'AND' got: %v", mv.Token)
	}

	if mv.Pos != 4 {
		t.Fatalf("expected position 4 got: %v", mv.Pos)
	}
}

const tstNotAndOr = "NOT AND OR"

func testNotAndOr(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatalf("expected parsed filter got: %v", pf)
	}

	if len(pf.Expr.Sequences) != 2 {
		t.Fatalf("expected two sequences got: %v", len(pf.Expr.Sequences))
	}

	seq := pf.Expr.Sequences[0]

	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor got: %v", seq.Factors)
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected two terms got: %v", factor.Terms)
	}

	tm := factor.Terms[0]
	if tm.UnaryOp != "" {
		t.Fatalf("expected no unary op got: %v", tm.UnaryOp)
	}

	if tm.Expr == nil {
		t.Fatalf("expected expression")
	}

	rs, ok := tm.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	// NOT AND
	m, ok := rs.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", rs.Comparable)
	}

	mv, ok := m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if mv.Value != "NOT" {
		t.Fatalf("expected 'NOT' got: %v", mv.Value)
	}

	if mv.Token != token.NOT {
		t.Fatalf("expected 'NOT' got: %v", mv.Token)
	}

	if mv.Pos != 0 {
		t.Fatalf("expected position 4 got: %v", mv.Pos)
	}

	if len(m.Fields) != 0 {
		t.Fatalf("expected no fields got: %v", m.Fields)
	}

	seq = pf.Expr.Sequences[1]
	m = seqMember(t, seq)

	mv, ok = m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Value)
	}

	if mv.Value != "OR" {
		t.Fatalf("expected 'OR' got: %v", mv.Value)
	}

	if mv.Token != token.OR {
		t.Fatalf("expected 'OR' got: %v", mv.Token)
	}

	if mv.Pos != 8 {
		t.Fatalf("expected position 8 got: %v", mv.Pos)
	}

	if len(m.Fields) != 0 {
		t.Fatalf("expected no fields got: %v", m.Fields)
	}
}
