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

// Complexity: 5
const tstMultipleExpressions = `enum = "ONE" AND i64 = 1`

func testMultipleExpressions(t *testing.T, x expr.FilterExpr) {
	ae, ok := x.(*expr.AndExpr)
	if !ok {
		t.Fatalf("expected and expression but got %T", x)
	}

	if len(ae.Expr) != 2 {
		t.Fatalf("expected 2 expressions but got %d", len(ae.Expr))
	}

	ce, ok := ae.Expr[0].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[0])
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

	ce, ok = ae.Expr[1].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[1])
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

	if right.Value != int64(1) {
		t.Fatalf("expected value 1 but got %v", right.Value)
	}
}

// Complexity: 320
const tstComplexExpression = `(enum = "ONE" AND i64 = 1) OR (enum = "TWO" AND i64 = sub.i32) OR rp_enum IN ["ONE", "TWO"]`

func testComplexExpression(t *testing.T, x expr.FilterExpr) {
	oe, ok := x.(*expr.OrExpr)
	if !ok {
		t.Fatalf("expected or expression but got %T", x)
	}

	if len(oe.Expr) != 3 {
		t.Fatalf("expected 3 expressions but got %d", len(oe.Expr))
	}

	cs, ok := oe.Expr[0].(*expr.CompositeExpr)
	if !ok {
		t.Fatalf("expected composite expression but got %T", oe.Expr[0])
	}

	if cs.Expr == nil {
		t.Fatalf("expected expression but got nil")
	}

	ae, ok := cs.Expr.(*expr.AndExpr)
	if !ok {
		t.Fatalf("expected and expression but got %T", oe.Expr[0])
	}

	if len(ae.Expr) != 2 {
		t.Fatalf("expected 2 expressions but got %d", len(ae.Expr))
	}

	ce, ok := ae.Expr[0].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[0])
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

	ce, ok = ae.Expr[1].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[1])
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

	if right.Value != int64(1) {
		t.Fatalf("expected value 1 but got %v", right.Value)
	}

	cs, ok = oe.Expr[1].(*expr.CompositeExpr)
	if !ok {
		t.Fatalf("expected composite expression but got %T", oe.Expr[1])
	}

	if cs.Expr == nil {
		t.Fatalf("expected expression but got nil")
	}

	ae, ok = cs.Expr.(*expr.AndExpr)
	if !ok {
		t.Fatalf("expected and expression but got %T", oe.Expr[1])
	}

	if len(ae.Expr) != 2 {
		t.Fatalf("expected 2 expressions but got %d", len(ae.Expr))
	}

	ce, ok = ae.Expr[0].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[0])
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

	if right.Value != protoreflect.EnumNumber(2) {
		t.Fatalf("expected value TWO but got %v", right.Value)
	}

	ce, ok = ae.Expr[1].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", ae.Expr[1])
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

	rfe, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if rfe.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'i32' field but got %s", rfe.Field)
	}

	rt, ok := rfe.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", rfe.Traversal)
	}

	if rt.Field != md.Fields().ByName("i32") {
		t.Fatalf("expected field 'i32' field but got %s", rt.Field)
	}

	ce, ok = oe.Expr[2].(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", oe.Expr[2])
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	left, ok = ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("rp_enum") {
		t.Fatalf("expected field 'rp_enum' field but got %s", left.Field)
	}

	ra, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected array expression but got %T", ce.Right)
	}

	if len(ra.Elements) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(ra.Elements))
	}

	ve, ok := ra.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ra.Elements[0])
	}

	if ve.Value != protoreflect.EnumNumber(1) {
		t.Fatalf("expected value ONE but got %v", ve.Value)
	}

	ve, ok = ra.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ra.Elements[1])
	}

	if ve.Value != protoreflect.EnumNumber(2) {
		t.Fatalf("expected value TWO but got %v", ve.Value)
	}
}
