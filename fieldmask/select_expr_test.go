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
	"testing"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/testpb"
	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

func TestParser_ParseSelectExpr(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		check func(t *testing.T, x *expr.MessageSelectExpr)
		isErr bool
		err   error
	}{
		{
			name: "single field",
			paths: []string{
				"name",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
				if f.Traversal != nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
			},
		},
		{
			name: "single field with traversal",
			paths: []string{
				"sub.name",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "sub" {
					t.Fatalf("unexpected field sub: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				subMS, ok := f.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(subMS.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
		{
			name: "multiple fields",
			paths: []string{
				"name",
				"i32",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 2 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
				if f.Traversal != nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				f = x.Fields[1]
				if f.Field != "i32" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
				if f.Traversal != nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
			},
		},
		{
			name: "map key field",
			paths: []string{
				"map_str_msg.key",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "map_str_msg" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				msk, ok := f.Traversal.(*expr.MapSelectKeysExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(msk.Keys) != 1 {
					t.Fatalf("unexpected number of keys: %d", len(msk.Keys))
				}
				k := msk.Keys[0]
				kv, ok := k.Key.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("unexpected key: %v", k.Key)
				}
				if kv.Value != "key" {
					t.Fatalf("unexpected key: %v", kv.Value)
				}
				// The traversal should be a MessageSelectExpr.
				if k.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", k.Traversal)
				}

				subMS, ok := k.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", k.Traversal)
				}

				if len(subMS.Fields) == 0 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
			},
		},
		{
			name: "double sub field consolidation",
			paths: []string{
				"sub.name",
				"sub.i32",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "sub" {
					t.Fatalf("unexpected field sub: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				subMS, ok := f.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(subMS.Fields) != 2 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}

				f = subMS.Fields[1]
				if f.Field != "i32" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
		{
			name: "input only",
			paths: []string{
				"input_only_str",
			},
			isErr: true,
			err:   ErrInvalidField,
		},
		{
			name: "unknown field",
			paths: []string{
				"unknown",
			},
			isErr: true,
			err:   ErrInvalidField,
		},
		{
			name: "unknown sub field",
			paths: []string{
				"sub.unknown",
			},
			isErr: true,
			err:   ErrInvalidField,
		},
		{
			name: "repeated wildstar selector",
			paths: []string{
				"rp_sub.*.name",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "rp_sub" {
					t.Fatalf("unexpected field sub: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				subMS, ok := f.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(subMS.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
		{
			name: "multiple fields of repeated wildstar selector",
			paths: []string{
				"rp_sub.*.name",
				"rp_sub.*.i32",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "rp_sub" {
					t.Fatalf("unexpected field sub: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				subMS, ok := f.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(subMS.Fields) != 2 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}

				f = subMS.Fields[1]
				if f.Field != "i32" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
		{
			name: "consolidated duplicates",
			paths: []string{
				"sub.name",
				"sub.name",
				"sub.name",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "sub" {
					t.Fatalf("unexpected field sub: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				subMS, ok := f.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(subMS.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}
				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
		{
			name: "map key wildcard selector",
			paths: []string{
				"map_str_msg.*.name",
			},
			check: func(t *testing.T, x *expr.MessageSelectExpr) {
				if len(x.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(x.Fields))
				}
				f := x.Fields[0]
				if f.Field != "map_str_msg" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
				if f.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}
				msk, ok := f.Traversal.(*expr.MapSelectKeysExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", f.Traversal)
				}

				if len(msk.Keys) != 1 {
					t.Fatalf("unexpected number of keys: %d", len(msk.Keys))
				}
				k := msk.Keys[0]
				if k.Key == nil {
					t.Fatalf("expected wildcard key but is nil")
				}
				_, ok = k.Key.(*expr.WildcardExpr)
				if !ok {
					t.Fatalf("expected wildcard key but is %T", k.Key)
				}
				// The traversal should be a MessageSelectExpr.
				if k.Traversal == nil {
					t.Fatalf("unexpected traversal: %v", k.Traversal)
				}

				subMS, ok := k.Traversal.(*expr.MessageSelectExpr)
				if !ok {
					t.Fatalf("unexpected traversal: %v", k.Traversal)
				}

				if len(subMS.Fields) != 1 {
					t.Fatalf("unexpected number of fields: %d", len(subMS.Fields))
				}

				f = subMS.Fields[0]
				if f.Field != "name" {
					t.Fatalf("unexpected field name: %s", f.Field)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := Parser{}
			if err := p.Reset(&testpb.Message{}, ErrHandlerOption(testErrorHandler(t, tc.isErr))); err != nil {
				t.Fatalf("failed to reset parser: %v", err)
			}

			fm := &fieldmaskpb.FieldMask{Paths: tc.paths}
			x, err := p.ParseSelectExpr(fm)
			if err != nil {
				if !tc.isErr {
					t.Fatalf("unexpected error: %v", err)
				}
				if tc.err != err {
					t.Fatalf("unexpected error: %v, expected: %v", err, tc.err)
				}
				return
			}
			defer x.Free()

			if tc.isErr {
				t.Fatalf("expected error: %v", tc.err)
			}

			tc.check(t, x)
		})
	}
}

func testErrorHandler(t *testing.T, wantsErr bool) scanner.ErrorHandler {
	return func(pos token.Position, msg string) {
		if !wantsErr {
			t.Errorf("unexpected error: %v, at: %d", msg, pos)
		}
	}
}
