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

const tstMapStringI32FieldEqDirect = `map_str_i32 = map{"test": 1, "test2": 2}`

func testMapStringI32FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_i32") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != int64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != int64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringI64FieldEqDirect = `map_str_i64 = map{"test": 1, "test2": 2}`

func testMapStringI64FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_i64") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != int64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != int64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringU32FieldEqDirect = `map_str_u32 = map{"test": 1, "test2": 2}`

func testMapStringU32FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_u32") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != uint64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringU64FieldEqDirect = `map_str_u64 = map{"test": 1, "test2": 2}`

func testMapStringU64FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_u64") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != uint64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != uint64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringS32FieldEqDirect = `map_str_s32 = map{"test": 1, "test2": 2}`

func testMapStringS32FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_s32") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != int64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != int64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringS64FieldEqDirect = `map_str_s64 = map{"test": 1, "test2": 2}`

func testMapStringS64FieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("map_str_s64") {
		t.Fatalf("expected field 'map_str_int' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.MapValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Values) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Values))
	}

	ve0 := right.Values[0]
	if ve0.Key.Value != "test" {
		t.Fatalf("expected key 'test' but got %s", ve0.Key.Value)
	}

	ve0v, ok := ve0.Value.(*expr.ValueExpr)
	if ve0v.Value != int64(1) {
		t.Fatalf("expected value 1 but got %d", ve0.Value)
	}

	ve1 := right.Values[1]
	if ve1.Key.Value != "test2" {
		t.Fatalf("expected key 'test2' but got %s", ve1.Key.Value)
	}

	ve1v, ok := ve1.Value.(*expr.ValueExpr)
	if ve1v.Value != int64(2) {
		t.Fatalf("expected value 2 but got %d", ve1.Value)
	}
}

const tstMapStringF32FieldEqDirect = `map_str_f32 = map{"test": 1, "test2": 2}`
