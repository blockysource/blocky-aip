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

const tstI32FieldEQDirect = `i32 = 42`

func testI32FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstI32FieldGTDirect = `i32 > 42`

func testI32FieldGTDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.GT {
		t.Fatalf("expected comparator %s but got %s", expr.GT, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstI32FieldEQIndirect = `i32 = sub.i32`

func testI32FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", right.Field)
	}
}

const tstI32FieldInArrayDirect = `i32 IN [42, 43]`

func testI32FieldInArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	if right.Elements[0].(*expr.ValueExpr).Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Elements[0].(*expr.ValueExpr).Value)
	}

	if right.Elements[1].(*expr.ValueExpr).Value != int64(43) {
		t.Fatalf("expected value 43 but got %d", right.Elements[1].(*expr.ValueExpr).Value)
	}
}

const tstI32FieldInArrayIndirect = `i32 IN rp_i32`

func testI32FieldInArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i32' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_i32").Name() {
		t.Fatalf("expected field 'rp_i32' field but got %s", right.Field)
	}
}

const tstI32FieldEQNegativeDirect = `i32 = -42`

func testI32FieldEQNegativeDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i32").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(-42) {
		t.Fatalf("expected value -42 but got %d", right.Value)
	}
}

const tstI64FieldEQDirect = `i64 = 42`

func testI64FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstI64FieldLTDirect = `i64 < 42`

func testI64FieldLTDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.LT {
		t.Fatalf("expected comparator %s but got %s", expr.LT, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstI64FieldEQIndirect = `i64 = sub.i64`

func testI64FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", right.Field)
	}
}

const tstI64FieldInArrayDirect = `i64 IN [42, 43]`

func testI64FieldInArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	if right.Elements[0].(*expr.ValueExpr).Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Elements[0].(*expr.ValueExpr).Value)
	}

	if right.Elements[1].(*expr.ValueExpr).Value != int64(43) {
		t.Fatalf("expected value 43 but got %d", right.Elements[1].(*expr.ValueExpr).Value)
	}
}

const tstI64FieldInArrayIndirect = `i64 IN rp_i64`

func testI64FieldInArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_i64").Name() {
		t.Fatalf("expected field 'rp_i64' field but got %s", right.Field)
	}
}

const tstI64FieldEQNegativeDirect = `i64 = -42`

func testI64FieldEQNegativeDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("i64").Name() {
		t.Fatalf("expected field 'i64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(-42) {
		t.Fatalf("expected value -42 but got %d", right.Value)
	}
}

const tstS32FieldEQDirect = `s32 = 42`

func testS32FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstS32FieldEQIndirect = `s32 = sub.s32`

func testS32FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", right.Field)
	}
}

const tstS32FieldInArrayDirect = `s32 IN [42, 43]`

func testS32FieldInArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	if right.Elements[0].(*expr.ValueExpr).Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Elements[0].(*expr.ValueExpr).Value)
	}

	if right.Elements[1].(*expr.ValueExpr).Value != int64(43) {
		t.Fatalf("expected value 43 but got %d", right.Elements[1].(*expr.ValueExpr).Value)
	}
}

const tstS32FieldInArrayIndirect = `s32 IN rp_s32`

func testS32FieldInArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_s32").Name() {
		t.Fatalf("expected field 'rp_s32' field but got %s", right.Field)
	}
}

const tstS32FieldEQNegativeDirect = `s32 = -42`

func testS32FieldEQNegativeDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s32").Name() {
		t.Fatalf("expected field 's32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(-42) {
		t.Fatalf("expected value -42 but got %d", right.Value)
	}
}

const tstS64FieldEQDirect = `s64 = 42`

func testS64FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Value)
	}
}

const tstS64FieldEQIndirect = `s64 = sub.s64`

func testS64FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", right.Field)
	}
}

const tstS64FieldInArrayDirect = `s64 IN [42, 43]`

func testS64FieldInArrayDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 values but got %d", len(right.Elements))
	}

	if right.Elements[0].(*expr.ValueExpr).Value != int64(42) {
		t.Fatalf("expected value 42 but got %d", right.Elements[0].(*expr.ValueExpr).Value)
	}

	if right.Elements[1].(*expr.ValueExpr).Value != int64(43) {
		t.Fatalf("expected value 43 but got %d", right.Elements[1].(*expr.ValueExpr).Value)
	}
}

const tstS64FieldInArrayIndirect = `s64 IN rp_s64`

func testS64FieldInArrayIndirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected in expression but got %T", x)
	}

	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", left.Field)
	}

	if ce.Comparator != expr.IN {
		t.Fatalf("expected comparator %s but got %s", expr.IN, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("rp_s64").Name() {
		t.Fatalf("expected field 'rp_s64' field but got %s", right.Field)
	}
}

const tstS64FieldEQNegativeDirect = `s64 = -42`

func testS64FieldEQNegativeDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("s64").Name() {
		t.Fatalf("expected field 's64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(-42) {
		t.Fatalf("expected value -42 but got %d", right.Value)
	}
}

const tstI32ComplexityEQDirect = `i32_complexity = 42`

func testI32ComplexityEQDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("Expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("Expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("Expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("i32_complexity").Name() {
		t.Fatalf("Expected field 'i32_complexity' field but got %s", left.Field)
	}
	if left.Complexity() != 44 {
		t.Fatalf("Expected complexity 44 but got %d", left.Complexity())
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("Expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("Expected value 42 but got %d", right.Value)
	}
}

const tstOneOfI32FieldEQDirect = `oneof_i32 = 42`

func testOneOfI32FieldEQDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("Expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.EQ {
		t.Fatalf("Expected comparator %s but got %s", expr.EQ, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("Expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("oneof_i32").Name() {
		t.Fatalf("Expected field 'oneof_i32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("Expected value expression but got %T", ce.Right)
	}

	if right.Value != int64(42) {
		t.Fatalf("Expected value 42 but got %d", right.Value)
	}
}
