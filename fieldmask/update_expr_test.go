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

package fieldmask

import (
	"math"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/testpb"
)

func TestParseUpdateExpr(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		msg   *testpb.Message
		check func(t *testing.T, x *expr.UpdateExpr)
	}{
		{
			name: "str",
			paths: []string{
				"str",
			},
			msg: &testpb.Message{
				Str: "test",
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "str" {
					t.Errorf("el.Field.Field = %v, want 'str'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != "test" {
					t.Errorf("el.Value = %v, want test", ev.Value)
				}
			},
		},
		{
			name: "empty str",
			paths: []string{
				"str",
			},
			msg: &testpb.Message{
				Str: "",
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "str" {
					t.Errorf("el.Field.Field = %v, want 'str'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != "" {
					t.Errorf("el.Value = %v, want ''", ev.Value)
				}
			},
		},
		{
			name: "i32",
			paths: []string{
				"i32",
			},
			msg: &testpb.Message{
				Name: "test",
				I32:  42,
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "i32" {
					t.Errorf("el.Field.Field = %v, want 'i32'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != int64(42) {
					t.Errorf("el.Value = %v, want 42", ev.Value)
				}
			},
		},
		{
			name: "s32",
			paths: []string{
				"s32",
			},
			msg: &testpb.Message{
				S32: 42,
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "s32" {
					t.Errorf("el.Field.Field = %v, want 's32'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != int64(42) {
					t.Errorf("el.Value = %v, want 42", ev.Value)
				}
			},
		},
		{
			name: "u32",
			paths: []string{
				"u32",
			},
			msg: &testpb.Message{
				U32: 42,
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "u32" {
					t.Errorf("el.Field.Field = %v, want 'u32'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != uint64(42) {
					t.Errorf("el.Value = %v, want 42", ev.Value)
				}
			},
		},
		{
			name: "i64",
			paths: []string{
				"i64",
			},
			msg: &testpb.Message{
				I64: 42,
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "i64" {
					t.Errorf("el.Field.Field = %v, want 'i64'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != int64(42) {
					t.Errorf("el.Value = %v, want 42", ev.Value)
				}
			},
		},
		{
			name: "sub",
			paths: []string{
				"sub",
			},
			msg: &testpb.Message{
				Sub: &testpb.Message{
					Name: "test",
				},
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "sub" {
					t.Errorf("el.Field.Field = %v, want 'sub'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				uev, ok := el.Value.(*expr.UpdateExpr)
				if !ok {
					t.Fatalf("el.Value is not a UpdateExpr but %T", el.Value)
				}

				// The sub message have exactly 1 value (name), 20 NULLABLE and 19 NON_EMPTY_DEFAULT fields
				if len(uev.Elements) != 40 {
					t.Errorf("len(ev.Elements) = %v, want 40", len(uev.Elements))
					return
				}

				el = uev.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				// The only selected field is name - it is also the first field by the index - thus the index is 0.
				if el.Field.Field != "name" {
					t.Errorf("el.Field.Field = %v, want 'name'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != "test" {
					t.Errorf("el.Value = %v, want test", ev.Value)
				}
			},
		},
		{
			name: "point with both fields",
			paths: []string{
				"point.x",
				"point.y",
			},
			msg: &testpb.Message{
				Point: &testpb.Point{
					X: 42,
					Y: 43.24,
				},
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 2 {
					t.Errorf("len(expr.Elements) = %v, want 2", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "point" {
					t.Errorf("el.Field.Field = %v, want 'point'", el.Field.Field)
					return
				}
				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ft, ok := el.Field.Traversal.(*expr.FieldSelectorExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal is not a FieldSelectorExpr but %T", el.Field.Traversal)
				}
				if ft.Field != "x" {
					t.Fatalf("el.Field.Traversal.Field = %v, want 'x'", ft.Field)
				}
				if ft.Message != "testpb.Point" {
					t.Fatalf("el.Field.Traversal.Message = %v, want testpb.Point", ft.Message)
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != float64(42) {
					t.Errorf("el.Value = %v, want 42", ev.Value)
				}

				el = x.Elements[1]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "point" {
					t.Errorf("el.Field.Field = %v, want 'point'", el.Field.Field)
					return
				}

				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ev, ok = el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				fv, ok := ev.Value.(float64)
				if !ok {
					t.Fatalf("el.Value is not a float64 but %T", ev.Value)
				}

				if math.Abs(fv-43.24) <= 1e-9 {
					t.Errorf("el.Value = %v, want 43.24", ev.Value)
				}
			},
		},
		{
			name: "map key selector",
			paths: []string{
				"map_str_str.key",
			},
			msg: &testpb.Message{
				MapStrStr: map[string]string{
					"key": "value",
				},
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "map_str_str" {
					t.Errorf("el.Field.Field = %v, want 'map_str_str'", el.Field.Field)
					return
				}
				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ft, ok := el.Field.Traversal.(*expr.MapKeyExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal is not a MapKeyExpr but %T", el.Field.Traversal)
				}

				fk, ok := ft.Key.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal.Key is not a ValueExpr but %T", ft.Key)
				}

				if fk.Value != "key" {
					t.Errorf("el.Field.Traversal.Key = %v, want 'key'", fk.Value)
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != "value" {
					t.Errorf("el.Value = %v, want 'value'", ev.Value)
				}
			},
		},
		{
			name: "sub map key selector",
			paths: []string{
				"sub.map_i32_str.653",
			},
			msg: &testpb.Message{
				Sub: &testpb.Message{
					MapI32Str: map[int32]string{
						653: "value",
					},
				},
			},
			check: func(t *testing.T, x *expr.UpdateExpr) {
				if x == nil {
					t.Errorf("expr is nil")
					return
				}

				if len(x.Elements) != 1 {
					t.Errorf("len(expr.Elements) = %v, want 1", len(x.Elements))
					return
				}

				el := x.Elements[0]

				if el.Field == nil {
					t.Errorf("el.Field is nil")
					return
				}

				if el.Field.Field != "sub" {
					t.Errorf("el.Field.Field = %v, want 'sub'", el.Field.Field)
					return
				}
				if el.Field.Message != "testpb.Message" {
					t.Errorf("el.Field.Message = %v, want testpb.Message", el.Field.Message)
					return
				}

				ft, ok := el.Field.Traversal.(*expr.FieldSelectorExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal is not a FieldSelectorExpr but %T", el.Field.Traversal)
				}

				if ft.Field != "map_i32_str" {
					t.Fatalf("el.Field.Traversal.Field = %v, want 'map_i32_str'", ft.Field)
				}

				fk, ok := ft.Traversal.(*expr.MapKeyExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal.Message is not a MapKeyExpr but %T", ft.Message)
				}

				kv, ok := fk.Key.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Field.Traversal.Message.Key is not a ValueExpr but %T", fk.Key)
				}

				if kv.Value != int64(653) {
					t.Errorf("el.Field.Traversal.Message.Key = %v, want 653", kv.Value)
				}

				ev, ok := el.Value.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("el.Value is not a ValueExpr but %T", el.Value)
				}

				if ev.Value != "value" {
					t.Errorf("el.Value = %v, want 'value'", ev.Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Parser
			if err := p.Reset(new(testpb.Message)); err != nil {
				t.Errorf("Reset() error = %v", err)
				return
			}

			mask := &fieldmaskpb.FieldMask{
				Paths: tt.paths,
			}
			got, err := p.ParseUpdateExpr(tt.msg, mask)
			if err != nil {
				t.Errorf("ParseUpdateExpr() error = %v", err)
				return
			}
			if got == nil {
				t.Errorf("ParseUpdateExpr() got = nil")
				return
			}

			tt.check(t, got)
		})
	}
}

func TestProtoReflectMessageFields(t *testing.T) {
	msg := testpb.Message{
		Name: "test",
	}

	ref := msg.ProtoReflect()

	ref.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if v.IsValid() {
			t.Logf("fd: %v, v: %v", fd, v)
		} else {
			t.Logf("fd: %v, v: nil", fd)
		}
		return true
	})
}
