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

package filteringfunc

import (
	"testing"
	"time"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering"
	"github.com/blockysource/blocky-aip/internal/testpb"
	"github.com/blockysource/blocky-aip/token"
)

var msgDesc = new(testpb.Message).ProtoReflect().Descriptor()

func TestUnixFunctionCall(t *testing.T) {
	testCases := []struct {
		name    string
		filter  string
		isErr   bool
		err     error
		checkFn func(t *testing.T, x expr.FilterExpr)
	}{
		{
			name:   "timestamp field EQ direct",
			filter: `timestamp = time.Unix(1614556800)`,
			checkFn: func(t *testing.T, x expr.FilterExpr) {
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

				if left.Field != msgDesc.Fields().ByName("timestamp").Name() {
					t.Fatalf("expected field 'timestamp' field but got %s", left.Field)
				}

				right, ok := ce.Right.(*expr.ValueExpr)
				if !ok {
					t.Fatalf("expected value expression but got %T", ce.Right)
				}

				expected := time.Unix(int64(1614556800), 0)

				rts, ok := right.Value.(time.Time)
				if !ok {
					t.Fatalf("expected int64 value bot got: %T", right.Value)
				}

				if !rts.Equal(expected) {
					t.Fatalf("expected value %s but got %s", expected, rts)
				}
			},
		},
		{
			name:   "timestamp field EQ indirect unix",
			filter: `timestamp = time.Unix(i64)`,
			checkFn: func(t *testing.T, x expr.FilterExpr) {
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

				if left.Field != msgDesc.Fields().ByName("timestamp").Name() {
					t.Fatalf("expected field 'timestamp' field but got %s", left.Field)
				}

				right, ok := ce.Right.(*expr.FunctionCallExpr)
				if !ok {
					t.Fatalf("expected value expression but got %T", ce.Right)
				}

				if len(right.Arguments) != 1 {
					t.Fatalf("expected 1 argument but got %d", len(right.Arguments))
				}

				arg, ok := right.Arguments[0].(*expr.FieldSelectorExpr)
				if !ok {
					t.Fatalf("expected value expression but got %T", right.Arguments[0])
				}

				if arg.Field != msgDesc.Fields().ByName("i64").Name() {
					t.Fatalf("expected field 'i64' field but got %s", arg.Field)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			it, err := filtering.NewInterpreter(msgDesc,
				filtering.RegisterFunction(TimeUnix()),
				filtering.ErrHandlerOpt(errHandler(t, tc.filter, tc.isErr)),
			)
			if err != nil {
				t.Fatalf("failed to create interpreter: %s", err)
			}

			x, err := it.Parse(tc.filter)
			if tc.isErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if tc.err != nil && tc.err != err {
					t.Fatalf("expected error %s but got %s", tc.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %s", err)
				}
				defer x.Free()
				tc.checkFn(t, x)
			}
		})
	}
}
func errHandler(t *testing.T, filter string, isErr bool) func(position token.Position, msg string) {
	return func(position token.Position, msg string) {
		if !isErr {
			t.Errorf("error at position %d: \n%s \n^ Error: %s", position, filter[position:], msg)
		}
	}
}
