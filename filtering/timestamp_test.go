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
	"time"

	"github.com/blockysource/blocky-aip/expr"
)

const tstTimestampFieldEQDirect = `timestamp = 2021-01-01T00:00:00Z`

func testTimestampFieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("timestamp").Name() {
		t.Fatalf("expected field 'timestamp' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	expected, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("expected value 2021-01-01T00:00:00Z but got %s", err)
	}

	rts, ok := right.Value.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time value bot got: %T", right.Value)
	}

	if !rts.Equal(expected) {
		t.Fatalf("expected value 2021-01-01T00:00:00Z but got %s", rts)
	}
}

const tstTimestampFieldEQIndirect = `timestamp = sub.timestamp`

func testTimestampFieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("timestamp").Name() {
		t.Fatalf("expected field 'timestamp' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("timestamp").Name() {
		t.Fatalf("expected field 'timestamp' field but got %s", right.Field)
	}
}

const tstTimestampFieldInArrayDirect = `timestamp IN [2021-01-01T00:00:00Z, 2021-02-01T00:00:00Z]`

func testTimestampFieldInArrayDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("timestamp").Name() {
		t.Fatalf("expected field 'timestamp' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ArrayExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if len(right.Elements) != 2 {
		t.Fatalf("expected 2 items in array but got %d", len(right.Elements))
	}

	expected1, err := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("expected value 2021-01-01T00:00:00Z but got %s", err)
	}

	expected2, err := time.Parse(time.RFC3339, "2021-02-01T00:00:00Z")
	if err != nil {
		t.Fatalf("expected value 2021-02-01T00:00:00Z but got %s", err)
	}

	rts1, ok := right.Elements[0].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected time.Time value bot got: %T", right.Elements[0])
	}

	rts2, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected time.Time value bot got: %T", right.Elements[1])
	}

	rtsTm1, ok := rts1.Value.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time value bot got: %T", rts1.Value)
	}

	rtsTm2, ok := rts2.Value.(time.Time)
	if !ok {
		t.Fatalf("expected time.Time value bot got: %T", rts2.Value)
	}

	if !rtsTm1.Equal(expected1) {
		t.Fatalf("expected value 2021-01-01T00:00:00Z but got %s", rtsTm1)
	}

	if !rtsTm2.Equal(expected2) {
		t.Fatalf("expected value 2021-02-01T00:00:00Z but got %s", rtsTm2)
	}

}
