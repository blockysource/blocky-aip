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

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/blockysource/blocky-aip/expr"
)

const tstMsgFieldEQDirect = `sub = testpb.Message{i64: 1, str: "value", enum: "ONE", bool: true, float: 1.0, rp_str: ["foo", "bar"], sub: {i64: 2}}`

func testMsgFieldEQDirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'msg' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	mp, ok := right.Value.(*dynamicpb.Message)
	if !ok {
		t.Fatalf("expected value map[string]any but got %T", right.Value)
	}

	iv := mp.Get(md.Fields().ByName("i64")).Int()
	if iv != int64(1) {
		t.Fatalf("expected value 1 but got %v", iv)
	}

	sv := mp.Get(md.Fields().ByName("str")).String()
	if sv != "value" {
		t.Fatalf("expected value 'value' but got %v", sv)
	}

	ev := mp.Get(md.Fields().ByName("enum")).Enum()
	if ev != protoreflect.EnumNumber(1) {
		t.Fatalf("expected value ONE but got %v", ev)
	}

	bv := mp.Get(md.Fields().ByName("bool")).Bool()
	if bv != true {
		t.Fatalf("expected value true but got %v", bv)
	}

	fv := mp.Get(md.Fields().ByName("float")).Float()
	if fv != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", fv)
	}

	lv := mp.Get(md.Fields().ByName("rp_str")).List()
	if lv.Len() != 2 {
		t.Fatalf("expected list of length 2 but got %v", lv.Len())
	}

	if lv.Get(0).String() != "foo" {
		t.Fatalf("expected value 'foo' but got %v", lv.Get(0).String())
	}

	if lv.Get(1).String() != "bar" {
		t.Fatalf("expected value 'bar' but got %v", lv.Get(1).String())
	}
}

const tstMsgFieldEQIndirect = `sub = sub.sub`

func testMsgFieldEQIndirect(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'msg' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.FieldSelectorExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	if right.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'msg' field but got %s", right.Field)
	}
}

const tstMsgFieldEQDirectUnnamed = `sub = {i64: 1, str: "value", enum: "ONE", bool: true, float: 1.0, rp_str: ["foo", "bar"], sub: {i64: 2}}`

func testMsgFieldEQDirectUnnamed(t *testing.T, x expr.FilterExpr) {
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

	if left.Field != md.Fields().ByName("sub").Name() {
		t.Fatalf("expected field 'msg' field but got %s", left.Field)
	}

	right, ok := ce.Right.(*expr.ValueExpr)
	if !ok {
		t.Fatalf("expected value expression but got %T", ce.Right)
	}

	mp, ok := right.Value.(*dynamicpb.Message)
	if !ok {
		t.Fatalf("expected value map[string]any but got %T", right.Value)
	}

	iv := mp.Get(md.Fields().ByName("i64")).Int()
	if iv != int64(1) {
		t.Fatalf("expected value 1 but got %v", iv)
	}

	sv := mp.Get(md.Fields().ByName("str")).String()
	if sv != "value" {
		t.Fatalf("expected value 'value' but got %v", sv)
	}

	ev := mp.Get(md.Fields().ByName("enum")).Enum()
	if ev != protoreflect.EnumNumber(1) {
		t.Fatalf("expected value ONE but got %v", ev)
	}

	bv := mp.Get(md.Fields().ByName("bool")).Bool()
	if bv != true {
		t.Fatalf("expected value true but got %v", bv)
	}

	fv := mp.Get(md.Fields().ByName("float")).Float()
	if fv != 1.0 {
		t.Fatalf("expected value 1.0 but got %v", fv)
	}

	lv := mp.Get(md.Fields().ByName("rp_str")).List()
	if lv.Len() != 2 {
		t.Fatalf("expected list of length 2 but got %v", lv.Len())
	}

	if lv.Get(0).String() != "foo" {
		t.Fatalf("expected value 'foo' but got %v", lv.Get(0).String())
	}

	if lv.Get(1).String() != "bar" {
		t.Fatalf("expected value 'bar' but got %v", lv.Get(1).String())
	}
}
