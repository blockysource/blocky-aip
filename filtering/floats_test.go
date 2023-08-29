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

package filtering

import (
	"testing"

	"github.com/blockysource/blocky-aip/expr"
)

const tstFloatFieldEQDirect = `float = 1.0`

func testFloatFieldEQDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", right.Value)
	}
}

const tstFloatFieldEQIndirect = `float = sub.float`

func testFloatFieldEQIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tr, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if tr.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", tr.Field)
	}
}

const tstFloatFieldINArrayDirect = `float IN [1.0, 2.0]`

func testFloatFieldINArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	v, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if v.Value != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", v.Value)
	}

	v, ok = right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if v.Value != 2.0 {
		t.Fatalf("expected value 2.0 but got %v", v.Value)
	}
}

const tstFloatFieldINArrayIndirect = `float IN rp_float`

func testFloatFieldINArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_float").Name() {
		t.Fatalf("expected field 'rp_float' field but got %s", right.Field)
	}
}

const tstFloadFieldEQNegative = `float = -1.0`

func testFloadFieldEQNegative(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("float").Name() {
		t.Fatalf("expected field 'float' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != -1.0 {
		t.Fatalf("expected value -1.0 but got %v", right.Value)
	}
}

const tstDoubleFieldEQDirect = `double = 1.0`

func testDoubleFieldEQDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", right.Value)
	}
}

const tstDoubleFieldEQIndirect = `double = sub.double`

func testDoubleFieldEQIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tr, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if tr.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", tr.Field)
	}
}

const tstDoubleFieldINArrayDirect = `double IN [1.0, 2.0]`

func testDoubleFieldINArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	v, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if v.Value != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", v.Value)
	}

	v, ok = right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if v.Value != 2.0 {
		t.Fatalf("expected value 2.0 but got %v", v.Value)
	}
}

const tstDoubleFieldINArrayIndirect = `double IN rp_double`

func testDoubleFieldINArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_double").Name() {
		t.Fatalf("expected field 'rp_double' field but got %s", right.Field)
	}
}

const tstDoubleFieldEQNegative = `double = -1.0`

func testDoubleFieldEQNegative(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("double").Name() {
		t.Fatalf("expected field 'double' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != -1.0 {
		t.Fatalf("expected value -1.0 but got %v", right.Value)
	}
}
