package protofiltering

import (
	"testing"
	"time"

	"github.com/blockysource/blocky-aip/expr"
)

const tstDurationFieldEQDirect = `duration = 1s`

func testDurationFieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second {
		t.Fatalf("expected value 1s but got %d", dv)
	}
}

const tstDurationFieldEQIndirect = `duration = sub.duration`

func testDurationFieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	rt, ok := right.Traversal.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Traversal)
	}

	if rt.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", rt.Field)
	}
}

const tstDurationFieldGEDirect = `duration >= 1s`

func testDurationFieldGEDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected compare expression but got %T", x)
	}
	if ce.Comparator != expr.GE {
		t.Fatalf("expected comparator %s but got %s", expr.GE, ce.Comparator)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second {
		t.Fatalf("expected value 1s but got %d", dv)
	}
}

const tstDurationFieldEQFractalDirect = `duration = 1.5s`

func testDurationFieldEQFractalDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second+500*time.Millisecond {
		t.Fatalf("expected value 1.5s but got %d", dv)
	}
}

const tstDurationFieldEQStructDirect = `duration = google.protobuf.Duration{seconds: 1, nanos: 500000000}`

func testDurationFieldEQStructDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second+500*time.Millisecond {
		t.Fatalf("expected value 1.5s but got %d", dv)
	}
}

const tstDurationFieldINArrayDirect = `duration IN [1s, 2.5s]`

func testDurationFieldINArrayDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("duration").Name() {
		t.Fatalf("expected field 'duration' field but got %s", left.Field)
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

	dv, ok := ve0.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", ve0.Value)
	}

	if dv != time.Second {
		t.Fatalf("expected value 1s but got %d", dv)
	}

	ve1, ok := right.Elements[1].(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", right.Elements[1])
	}

	dv, ok = ve1.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", ve1.Value)
	}

	if dv != 2*time.Second+500*time.Millisecond {
		t.Fatalf("expected value 1.5s but got %d", dv)
	}
}

const tstMapStringDurationFieldHasDirect = `map_str_duration."key":1s`

func testMapStringDurationFieldHasDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected has expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("map_str_duration").Name() {
		t.Fatalf("expected field 'map_str_duration' field but got %s", left.Field)
	}

	// Field selector has a Map key selector in its Traversal.
	mk, ok := left.Traversal.(*expr.MapKeyExpr)
	if !ok {
		t.Fatalf("expected map key expression but got %T", left.Traversal)
	}

	mkv, ok := mk.Key.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", mk.Key)
	}

	if mkv.Value != "key" {
		t.Fatalf("expected value 'key' but got %s", mkv.Value)
	}

	if ce.Comparator != expr.HAS {
		t.Fatalf("expected comparator %s but got %s", expr.HAS, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second {
		t.Fatalf("expected value 1s but got %d", dv)
	}
}

const tstRepeatedDurationHasDirect = `rp_duration:1s`

func testRepeatedDurationHasDirect(t *testing.T, x expr.FilterExpr) {
	ce, ok := x.(*expr.CompareExpr)
	if !ok {
		t.Fatalf("expected has expression but got %T", x)
	}
	left, ok := ce.Left.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Left)
	}

	if left.Field != md.Fields().ByName("rp_duration").Name() {
		t.Fatalf("expected field 'rp_duration' field but got %s", left.Field)
	}

	if ce.Comparator != expr.HAS {
		t.Fatalf("expected comparator %s but got %s", expr.HAS, ce.Comparator)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	dv, ok := right.Value.(time.Duration)
	if !ok {
		t.Fatalf("expected value 1s but got %T", right.Value)
	}

	if dv != time.Second {
		t.Fatalf("expected value 1s but got %d", dv)
	}
}
