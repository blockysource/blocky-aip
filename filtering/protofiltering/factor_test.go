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

const tstFactorWithMultipleORs = `enum = "ONE" OR i64 = 2 OR enum = "THREE"`

func testFactorWithMultipleORs(t *testing.T, x expr.FilterExpr) {
	oe, ok := x.(*expr.OrExpr)
	if !ok {
		t.Fatalf("expected or expression but got %T", x)
	}

	if len(oe.Expr) != 3 {
		t.Fatalf("expected 3 expressions but got %d", len(oe.Expr))
	}

	ce, ok := oe.Expr[0].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", oe.Expr[0])
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

	ce, ok = oe.Expr[1].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", oe.Expr[1])
	}

	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}

	left, ok = ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i64") {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok = ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(2) {
		t.Fatalf("expected value 2 but got %v", right.Value)
	}

	ce, ok = oe.Expr[2].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", oe.Expr[2])
	}

	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}

	left, ok = ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("enum") {
		t.Fatalf("expected field 'enum' field but got %s", left.Field)
	}

	right, ok = ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != protoreflect.EnumNumber(3) {
		t.Fatalf("expected value THREE but got %v", right.Value)
	}
}
