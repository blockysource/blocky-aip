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

const singleSequenceMember = "a"

func testSingleSequenceMember(t *testing.T, pf *ParsedFilter) {
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
}

const singleSequenceWithStringMember = `"a"`

func testSingleSequenceWithStringMember(t *testing.T, pf *ParsedFilter) {
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
}



const deepNestedStringMember = `"expr".'type_map'.1."type"`

func testDeepNestedStringMember(t *testing.T, pf *ParsedFilter) {
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
}

const deepNestedText = "expr.type_map.1.type"

func testDeepNestedTextMember(t *testing.T, pf *ParsedFilter) {
	// expr.type_map.1.type
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
}
