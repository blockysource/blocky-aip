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

package protofiltering

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/expr"
)

const tstTermWithNOTKeyword = `NOT enum = "ONE"`

func testTermWithNOTKeyword(t *testing.T, x expr.FilterExpr) {
	ne, ok := x.(*expr.NotExpr)
	if !ok {
		t.Fatalf("expected not expression but got %T", x)
	}

	ce, ok := ne.Expr.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ne.Expr)
	}

	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("enum") {
		t.Fatalf("expected field 'enum' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != protoreflect.EnumNumber(1) {
		t.Fatalf("expected value ONE but got %v", right.Value)
	}
}

const tstTermWithNOTKeywordAndParentheses = `NOT (enum = "ONE")`

func testTermWithNOTKeywordAndParentheses(t *testing.T, x expr.FilterExpr) {
	ne, ok := x.(*expr.NotExpr)
	if !ok {
		t.Fatalf("expected not expression but got %T", x)
	}

	cm, ok := ne.Expr.(*expr.CompositeExpr)
	if !ok {
		t.Fatalf("expected composite expression but got %T", ne.Expr)
	}

	ce, ok := cm.Expr.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ne.Expr)
	}

	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("enum") {
		t.Fatalf("expected field 'enum' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != protoreflect.EnumNumber(1) {
		t.Fatalf("expected value ONE but got %v", right.Value)
	}
}
