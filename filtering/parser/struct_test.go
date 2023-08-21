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

const structExpr = `geo.Point{Lat: 0.3, Lon: 15.6}`

func testStructExpr(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatal("expected parsed filter got: nil")
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
		t.Error("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatal("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatal("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	st, ok := expr.Comparable.(*ast.StructExpr)
	if !ok {
		t.Fatal("expected struct expression")
	}

	if len(st.Name) != 2 {
		t.Fatalf("expected two name got: %v", len(st.Name))
	}

	if st.Position() != 0 {
		t.Fatalf("expected position 0 got: %v", st.Position())
	}

	first, ok := st.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[0])
	}

	if first.Value != "geo" {
		t.Fatalf("expected 'geo' got: %v", first.Value)
	}

	second, ok := st.Name[1].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[1])
	}

	if second.Value != "Point" {
		t.Fatalf("expected 'Point' got: %v", second.Value)
	}

	if st.LBrace != 9 {
		t.Fatalf("expected lbrace 9 got: %v", st.LBrace)
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f1 := st.Elements[0]
	if f1.Name == nil {
		t.Fatal("expected field name")
	}

	if f1.Position() != 10 {
		t.Fatalf("expected position 10 got: %v", f1.Position())
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f1.Name))
	}

	tl, ok := f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "Lat" {
		t.Fatalf("expected 'Lat' got: %v", tl.Value)
	}

	if f1.Colon != 13 {
		t.Fatalf("expected colon 13 got: %v", f1.Colon)
	}

	if f1.Value == nil {
		t.Fatal("expected field value")
	}

	f1v, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f1.Value)
	}

	if f1v.Position() != 15 {
		t.Fatalf("expected position 15 got: %v", f1v.Position())
	}

	if f1v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f1v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Value)
	}

	if tl.Value != "0" {
		t.Fatalf("expected '0' got: %v", tl.Value)
	}

	if len(f1v.Fields) != 1 {
		t.Fatalf("expected no fields got: %v", len(f1v.Fields))
	}

	tl, ok = f1v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Fields[0])
	}

	if tl.Value != "3" {
		t.Fatalf("expected '3' got: %v", tl.Value)
	}

	f2 := st.Elements[1]
	if f2.Name == nil {
		t.Fatal("expected field name")
	}

	if f2.Position() != 20 {
		t.Fatalf("expected position 17 got: %v", f2.Position())
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f2.Name))
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "Lon" {
		t.Fatalf("expected 'Lon' got: %v", tl.Value)
	}

	if f2.Colon != 23 {
		t.Fatalf("expected colon 23 got: %v", f2.Colon)
	}

	if f2.Value == nil {
		t.Fatal("expected field value")
	}

	f2v, ok := f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f2.Value)
	}

	if f2v.Position() != 25 {
		t.Fatalf("expected position 25 got: %v", f2v.Position())
	}

	if f2v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f2v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Value)
	}

	if tl.Value != "15" {
		t.Fatalf("expected '15' got: %v", tl.Value)
	}

	if len(f2v.Fields) != 1 {
		t.Fatalf("expected one field got: %v", len(f2v.Fields))
	}

	tl, ok = f2v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Fields[0])
	}

	if tl.Value != "6" {
		t.Fatalf("expected '6' got: %v", tl.Value)
	}

	if st.RBrace != 29 {
		t.Fatalf("expected rbrace 29 got: %v", st.RBrace)
	}
}

const structExprWithNewLines = `geo.Point{
	Lat: 0.3,
	Lon: 15.6
}`

func testStructExprWithNewLines(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatal("expected parsed filter got: nil")
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
		t.Error("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatal("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatal("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	st, ok := expr.Comparable.(*ast.StructExpr)
	if !ok {
		t.Fatal("expected struct expression")
	}

	if len(st.Name) != 2 {
		t.Fatalf("expected two name got: %v", len(st.Name))
	}

	if st.Position() != 0 {
		t.Fatalf("expected position 0 got: %v", st.Position())
	}

	first, ok := st.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[0])
	}

	if first.Value != "geo" {
		t.Fatalf("expected 'geo' got: %v", first.Value)
	}

	second, ok := st.Name[1].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[1])
	}

	if second.Value != "Point" {
		t.Fatalf("expected 'Point' got: %v", second.Value)
	}

	if st.LBrace != 9 {
		t.Fatalf("expected lbrace 9 got: %v", st.LBrace)
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f1 := st.Elements[0]
	if f1.Name == nil {
		t.Fatal("expected field name")
	}

	if f1.Position() != 12 {
		t.Fatalf("expected position 10 got: %v", f1.Position())
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f1.Name))
	}

	tl, ok := f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "Lat" {
		t.Fatalf("expected 'Lat' got: %v", tl.Value)
	}

	if f1.Colon != 15 {
		t.Fatalf("expected colon 15 got: %v", f1.Colon)
	}

	if f1.Value == nil {
		t.Fatal("expected field value")
	}

	f1v, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f1.Value)
	}

	if f1v.Position() != 17 {
		t.Fatalf("expected position 17 got: %v", f1v.Position())
	}

	if f1v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f1v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Value)
	}

	if tl.Value != "0" {
		t.Fatalf("expected '0' got: %v", tl.Value)
	}

	if len(f1v.Fields) != 1 {
		t.Fatalf("expected no fields got: %v", len(f1v.Fields))
	}

	tl, ok = f1v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Fields[0])
	}

	if tl.Value != "3" {
		t.Fatalf("expected '3' got: %v", tl.Value)
	}

	f2 := st.Elements[1]
	if f2.Name == nil {
		t.Fatal("expected field name")
	}

	if f2.Position() != 23 {
		t.Fatalf("expected position 23 got: %v", f2.Position())
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f2.Name))
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "Lon" {
		t.Fatalf("expected 'Lon' got: %v", tl.Value)
	}

	if f2.Colon != 26 {
		t.Fatalf("expected colon 26 got: %v", f2.Colon)
	}

	if f2.Value == nil {
		t.Fatal("expected field value")
	}

	f2v, ok := f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f2.Value)
	}

	if f2v.Position() != 28 {
		t.Fatalf("expected position 28 got: %v", f2v.Position())
	}

	if f2v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f2v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Value)
	}

	if tl.Value != "15" {
		t.Fatalf("expected '15' got: %v", tl.Value)
	}

	if len(f2v.Fields) != 1 {
		t.Fatalf("expected one field got: %v", len(f2v.Fields))
	}

	tl, ok = f2v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Fields[0])
	}

	if tl.Value != "6" {
		t.Fatalf("expected '6' got: %v", tl.Value)
	}

	if st.RBrace != 33 {
		t.Fatalf("expected rbrace 33 got: %v", st.RBrace)
	}
}

const structExprWithNewLinesEndedWithComma = `geo.Point{
	Lat: 0.3,
	Lon: 15.6,
}`

func testStructExprWithNewLinesEndedWithComma(t *testing.T, pf *ParsedFilter) {
	if pf.Expr == nil {
		t.Fatal("expected parsed filter got: nil")
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
		t.Error("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatal("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatal("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	st, ok := expr.Comparable.(*ast.StructExpr)
	if !ok {
		t.Fatal("expected struct expression")
	}

	if len(st.Name) != 2 {
		t.Fatalf("expected two name got: %v", len(st.Name))
	}

	if st.Position() != 0 {
		t.Fatalf("expected position 0 got: %v", st.Position())
	}

	first, ok := st.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[0])
	}

	if first.Value != "geo" {
		t.Fatalf("expected 'geo' got: %v", first.Value)
	}

	second, ok := st.Name[1].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", st.Name[1])
	}

	if second.Value != "Point" {
		t.Fatalf("expected 'Point' got: %v", second.Value)
	}

	if st.LBrace != 9 {
		t.Fatalf("expected lbrace 9 got: %v", st.LBrace)
	}

	if len(st.Elements) != 2 {
		t.Fatalf("expected two fields got: %v", len(st.Elements))
	}

	f1 := st.Elements[0]
	if f1.Name == nil {
		t.Fatal("expected field name")
	}

	if f1.Position() != 12 {
		t.Fatalf("expected position 10 got: %v", f1.Position())
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f1.Name))
	}

	tl, ok := f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "Lat" {
		t.Fatalf("expected 'Lat' got: %v", tl.Value)
	}

	if f1.Colon != 15 {
		t.Fatalf("expected colon 15 got: %v", f1.Colon)
	}

	if f1.Value == nil {
		t.Fatal("expected field value")
	}

	f1v, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f1.Value)
	}

	if f1v.Position() != 17 {
		t.Fatalf("expected position 17 got: %v", f1v.Position())
	}

	if f1v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f1v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Value)
	}

	if tl.Value != "0" {
		t.Fatalf("expected '0' got: %v", tl.Value)
	}

	if len(f1v.Fields) != 1 {
		t.Fatalf("expected no fields got: %v", len(f1v.Fields))
	}

	tl, ok = f1v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Fields[0])
	}

	if tl.Value != "3" {
		t.Fatalf("expected '3' got: %v", tl.Value)
	}

	f2 := st.Elements[1]
	if f2.Name == nil {
		t.Fatal("expected field name")
	}

	if f2.Position() != 23 {
		t.Fatalf("expected position 23 got: %v", f2.Position())
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f2.Name))
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "Lon" {
		t.Fatalf("expected 'Lon' got: %v", tl.Value)
	}

	if f2.Colon != 26 {
		t.Fatalf("expected colon 26 got: %v", f2.Colon)
	}

	if f2.Value == nil {
		t.Fatal("expected field value")
	}

	f2v, ok := f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f2.Value)
	}

	if f2v.Position() != 28 {
		t.Fatalf("expected position 28 got: %v", f2v.Position())
	}

	if f2v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f2v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Value)
	}

	if tl.Value != "15" {
		t.Fatalf("expected '15' got: %v", tl.Value)
	}

	if len(f2v.Fields) != 1 {
		t.Fatalf("expected one field got: %v", len(f2v.Fields))
	}

	tl, ok = f2v.Fields[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Fields[0])
	}

	if tl.Value != "6" {
		t.Fatalf("expected '6' got: %v", tl.Value)
	}

	if st.RBrace != 34 {
		t.Fatalf("expected rbrace 33 got: %v", st.RBrace)
	}
}

const mapStructComparable = "map{1: 2}"

func testMapStructComparable(t *testing.T, pf *ParsedFilter) {
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
	if term.UnaryOp != "" {
		t.Fatalf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatalf("expected comparable")
	}

	m, ok := expr.Comparable.(*ast.StructExpr)
	if !ok {
		t.Fatalf("expected struct expression")
	}

	if len(m.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(m.Name))
	}

	if m.Position() != 0 {
		t.Fatalf("expected position 0 got: %v", m.Position())
	}

	first, ok := m.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", m.Name[0])
	}

	if first.Value != "map" {
		t.Fatalf("expected 'map' got: %v", first.Value)
	}

	if !m.IsMap() {
		t.Fatalf("expected map")
	}

	if m.LBrace != 3 {
		t.Fatalf("expected lbrace 3 got: %v", m.LBrace)
	}

	if len(m.Elements) != 1 {
		t.Fatalf("expected one field got: %v", len(m.Elements))
	}

	f1 := m.Elements[0]
	if f1.Name == nil {
		t.Fatal("expected field name")
	}

	if f1.Position() != 4 {
		t.Fatalf("expected position 4 got: %v", f1.Position())
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f1.Name))
	}

	tl, ok := f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "1" {
		t.Fatalf("expected '1' got: %v", tl.Value)
	}

	if f1.Colon != 5 {
		t.Fatalf("expected colon 5 got: %v", f1.Colon)
	}

	if f1.Value == nil {
		t.Fatal("expected field value")
	}

	f1v, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f1.Value)
	}

	if f1v.Position() != 7 {
		t.Fatalf("expected position 7 got: %v", f1v.Position())
	}

	if f1v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f1v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Value)
	}

	if tl.Value != "2" {
		t.Fatalf("expected '2' got: %v", tl.Value)
	}

	if len(f1v.Fields) != 0 {
		t.Fatalf("expected no fields got: %v", len(f1v.Fields))
	}

	if m.RBrace != 8 {
		t.Fatalf("expected rbrace 8 got: %v", m.RBrace)
	}
}

const tstUnnamedStruct = `{a: 1, b: 2}`

func testUnnamedStruct(t *testing.T, pf *ParsedFilter) {
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
	if term.UnaryOp != "" {
		t.Fatalf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatalf("expected comparable")
	}

	m, ok := expr.Comparable.(*ast.StructExpr)
	if !ok {
		t.Fatalf("expected struct expression")
	}

	if len(m.Name) != 0 {
		t.Fatalf("expected no name got: %v", len(m.Name))
	}

	if m.Position() != 0 {
		t.Fatalf("expected position 0 got: %v", m.Position())
	}

	if m.LBrace != 0 {
		t.Fatalf("expected lbrace 0 got: %v", m.LBrace)
	}

	if len(m.Elements) != 2 {
		t.Fatalf("expected two field got: %v", len(m.Elements))
	}

	f1 := m.Elements[0]
	if f1.Name == nil {
		t.Fatal("expected field name")
	}

	if f1.Position() != 1 {
		t.Fatalf("expected position 1 got: %v", f1.Position())
	}

	if len(f1.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f1.Name))
	}

	tl, ok := f1.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1.Name)
	}

	if tl.Value != "a" {
		t.Fatalf("expected 'a' got: %v", tl.Value)
	}

	if f1.Colon != 2 {
		t.Fatalf("expected colon 2 got: %v", f1.Colon)
	}

	if f1.Value == nil {
		t.Fatal("expected field value")
	}

	f1v, ok := f1.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f1.Value)
	}

	if f1v.Position() != 4 {
		t.Fatalf("expected position 4 got: %v", f1v.Position())
	}

	if f1v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f1v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f1v.Value)
	}

	if tl.Value != "1" {
		t.Fatalf("expected '1' got: %v", tl.Value)
	}

	if len(f1v.Fields) != 0 {
		t.Fatalf("expected no fields got: %v", len(f1v.Fields))
	}

	f2 := m.Elements[1]
	if f2.Name == nil {
		t.Fatal("expected field name")
	}

	if f2.Position() != 7 {
		t.Fatalf("expected position 7 got: %v", f2.Position())
	}

	if len(f2.Name) != 1 {
		t.Fatalf("expected one name got: %v", len(f2.Name))
	}

	tl, ok = f2.Name[0].(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2.Name)
	}

	if tl.Value != "b" {
		t.Fatalf("expected 'b' got: %v", tl.Value)
	}

	if f2.Colon != 8 {
		t.Fatalf("expected colon 8 got: %v", f2.Colon)
	}

	if f2.Value == nil {
		t.Fatal("expected field value")
	}

	f2v, ok := f2.Value.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member expression got: %T", f2.Value)
	}

	if f2v.Position() != 10 {
		t.Fatalf("expected position 10 got: %v", f2v.Position())
	}

	if f2v.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok = f2v.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %T", f2v.Value)
	}

	if tl.Value != "2" {
		t.Fatalf("expected '2' got: %v", tl.Value)
	}

	if len(f2v.Fields) != 0 {
		t.Fatalf("expected no fields got: %v", len(f2v.Fields))
	}

	if m.RBrace != 11 {
		t.Fatalf("expected rbrace 11 got: %v", m.RBrace)
	}
}
