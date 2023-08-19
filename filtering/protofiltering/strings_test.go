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

const tstStringFieldEqDirect = `name = "test"`

func testStringFieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", right.Value)
	}
}

const tstStringFieldInArray = `name IN ["test", "test2"]`

func testStringFieldInArray(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Elements))
	}

	ve0, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if ve0.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", ve0.Value)
	}

	ve1, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if ve1.Value != "test2" {
		t.Fatalf("expected value 'test2' but got %s", ve1.Value)
	}
}

const tstStringFieldEqIndirect = `name = str`

func testStringFieldEqIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("str") {
		t.Fatalf("expected field 'str' field but got %s", right.Field)
	}
}

const tstStringFieldEqStringSearch = `name = "*test*"`

func testStringFieldEqStringSearch(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.StringSearchExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", right.Value)
	}

	if !right.PrefixWildcard {
		t.Fatalf("expected prefix wildcard")
	}

	if !right.SuffixWildcard {
		t.Fatalf("expected suffix wildcard")
	}
}

const tstStringFieldEqStringSearchPrefix = `name = "*test"`

func testStringFieldEqStringSearchPrefix(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.StringSearchExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", right.Value)
	}

	if !right.PrefixWildcard {
		t.Fatalf("expected prefix wildcard")
	}

	if right.SuffixWildcard {
		t.Fatalf("expected no suffix wildcard")
	}
}

const tstStringFieldEqStringSearchSuffix = `name = "test*"`

func testStringFieldEqStringSearchSuffix(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("name") {
		t.Fatalf("expected field 'name' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.StringSearchExpr)
	if !ok {
		t.Fatalf("expected StringSearchExpr but got %T", ce.Right)
	}

	if right.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", right.Value)
	}

	if right.PrefixWildcard {
		t.Fatalf("expected no prefix wildcard")
	}

	if !right.SuffixWildcard {
		t.Fatalf("expected suffix wildcard")
	}
}

const tstRepeatedStringFieldEqDirect = `rp_str = ["test", "test2"]`

func testRepeatedStringFieldEqDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("rp_str") {
		t.Fatalf("expected field 'rp_str' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 elements but got %d", len(right.Elements))
	}

	ve0, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[0])
	}

	if ve0.Value != "test" {
		t.Fatalf("expected value 'test' but got %s", ve0.Value)
	}

	ve1, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	if ve1.Value != "test2" {
		t.Fatalf("expected value 'test2' but got %s", ve1.Value)
	}
}


