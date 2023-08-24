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

const tstStructPbFieldEQDirectString = `struct = "{\"field\": 1, \"field2\": \"value\"}"`

func testStructPbFieldEQDirectString(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("struct").Name() {
		t.Fatalf("expected field 'struct' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	mp, ok := right.Value.(map[string]any)
	if !ok {
		t.Fatalf("expected value map[string]any but got %T", right.Value)
	}

	fv, ok := mp["field"]
	if !ok {
		t.Fatalf("expected 'field' key in map but got %v", mp)
	}
	if fv != float64(1) {
		t.Fatalf("expected value 1 but got %v", fv)
	}

	fv, ok = mp["field2"]
	if !ok {
		t.Fatalf("expected 'field2' key in map but got %v", mp)
	}

	if fv != "value" {
		t.Fatalf("expected value 'value' but got %v", fv)
	}
}

const tstStructPbFieldEQDirectMessage = `struct = google.protobuf.Struct{fields: map{"field": google.protobuf.Value{number_value: 1}, "field2": google.protobuf.Value{string_value: "value"}}}`

func testStructPbFieldEQDirectMessage(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("struct").Name() {
		t.Fatalf("expected field 'struct' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	mp, ok := right.Value.(map[string]any)
	if !ok {
		t.Fatalf("expected value map[string]any but got %T", right.Value)
	}

	fv, ok := mp["field"]
	if !ok {
		t.Fatalf("expected 'field' key in map but got %v", mp)
	}
	if fv != float64(1) {
		t.Fatalf("expected value 1 but got %v", fv)
	}

	fv, ok = mp["field2"]
	if !ok {
		t.Fatalf("expected 'field2' key in map but got %v", mp)
	}

	if fv != "value" {
		t.Fatalf("expected value 'value' but got %v", fv)
	}
}
