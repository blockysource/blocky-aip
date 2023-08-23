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

	"github.com/blockysource/blocky-aip/expr"
)

const tstU32FieldEQDirect = `u32 = 1`

func testU32FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u32") {
		t.Fatalf("expected field 'u32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", right.Value)
	}
}

const tstU32FieldINArray = `u32 IN [1, 2, 3]`

func testU32FieldINArray(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u32") {
		t.Fatalf("expected field 'u32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 3 {
		t.Fatalf("expected 3 items in array but got %d", len(right.Elements))
	}

	v1, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if v1.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", v1.Value)
	}

	v2, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if v2.Value != uint64(2) {
		t.Fatalf("expected value 2 but got %d", v2.Value)
	}

	v3, ok := right.Elements[2].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[2])
	}

	if v3.Value != uint64(3) {
		t.Fatalf("expected value 3 but got %d", v3.Value)
	}
}

const tstU32FieldINArrayIndirect = `u32 IN sub.rp_u32`

func testU32FieldINArrayIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u32") {
		t.Fatalf("expected field 'u32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("rp_u32") {
		t.Fatalf("expected field 'rp_u32' field but got %s", tf.Field)
	}
}

const tstU32FieldEQIndirect = `u32 = sub.u32`

func testU32FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u32") {
		t.Fatalf("expected field 'u32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("u32") {
		t.Fatalf("expected field 'u32' field but got %s", tf.Field)
	}
}

const tstU64FieldEQDirect = `u64 = 1`

func testU64FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u64") {
		t.Fatalf("expected field 'u64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", right.Value)
	}
}

const tstU64FieldINArray = `u64 IN [1, 2, 3]`

func testU64FieldINArray(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u64") {
		t.Fatalf("expected field 'u64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 3 {
		t.Fatalf("expected 3 items in array but got %d", len(right.Elements))
	}

	v1, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if v1.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", v1.Value)
	}

	v2, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if v2.Value != uint64(2) {
		t.Fatalf("expected value 2 but got %d", v2.Value)
	}

	v3, ok := right.Elements[2].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[2])
	}

	if v3.Value != uint64(3) {
		t.Fatalf("expected value 3 but got %d", v3.Value)
	}
}

const tstU64FieldINArrayIndirect = `u64 IN sub.rp_u64`

func testU64FieldINArrayIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u64") {
		t.Fatalf("expected field 'u64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("rp_u64") {
		t.Fatalf("expected field 'rp_u64' field but got %s", tf.Field)
	}
}

const tstU64FieldEQIndirect = `u64 = sub.u64`

func testU64FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("u64") {
		t.Fatalf("expected field 'u64' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("u64") {
		t.Fatalf("expected field 'u64' field but got %s", tf.Field)
	}
}

const tstF32FieldEQDirect = `f32 = 1`

func testF32FieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("f32") {
		t.Fatalf("expected field 'f32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %f", right.Value)
	}
}

const tstF32FieldINArray = `f32 IN [1, 2, 3]`

func testF32FieldINArray(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("f32") {
		t.Fatalf("expected field 'f32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 3 {
		t.Fatalf("expected 3 items in array but got %d", len(right.Elements))
	}

	v1, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if v1.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %f", v1.Value)
	}

	v2, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if v2.Value != uint64(2) {
		t.Fatalf("expected value 2 but got %f", v2.Value)
	}

	v3, ok := right.Elements[2].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[2])
	}

	if v3.Value != uint64(3) {
		t.Fatalf("expected value 3 but got %f", v3.Value)
	}
}

const tstF32FieldINArrayIndirect = `f32 IN sub.rp_f32`

func testF32FieldINArrayIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("f32") {
		t.Fatalf("expected field 'f32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("rp_f32") {
		t.Fatalf("expected field 'rp_f32' field but got %s", tf.Field)
	}
}

const tstF32FieldEQIndirect = `f32 = sub.f32`

func testF32FieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("f32") {
		t.Fatalf("expected field 'f32' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub") {
		t.Fatalf("expected field 'sub' field but got %s", right.Field)
	}

	tf, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected field selector expression but got %T", right.Traversal)
	}

	if tf.Field != md.Fields().ByName("f32") {
		t.Fatalf("expected field 'f32' field but got %s", tf.Field)
	}
}

const tstF64FieldEQDirect = `f64 = 1`

func testF64FieldEQDirect(t *testing.T, x expr.FilterExpr) {

}
