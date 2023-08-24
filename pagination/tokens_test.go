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

package pagination

import (
	"testing"
	"time"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/testpb"
)

func TestTokens(t *testing.T) {
	c := expr.Composer{Desc: new(testpb.Message).ProtoReflect().Descriptor()}
	tests := []struct {
		name  string
		next  expr.FilterExpr
		order *expr.OrderByExpr
	}{
		{
			name: "simple",
			next: c.Compare(
				c.MustSelect("name"),
				expr.GT,
				c.Value("foo"),
			),
		},
		{
			name: "with order",
			next: c.Compare(
				c.MustSelect("i32"),
				expr.GT,
				c.Value(1),
			),
			order: c.OrderBy(
				c.MustOrderByField(
					"i32",
					expr.ASC,
				),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next := NextTokenExpr{
				Filter:  tc.next,
				OrderBy: tc.order,
			}
			token, err := TokenizeStruct(next)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			res, err := DecodeToken[NextTokenExpr](token)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !res.Filter.Equals(tc.next) {
				t.Fatalf("expected filter %v but got %v", tc.next, res.Filter)
			}

			if tc.order == nil && res.OrderBy != nil {
				t.Fatalf("expected ordering %v but got %v", tc.order, res.OrderBy)
			}
			if tc.order == nil && res.OrderBy != nil {
				t.Fatalf("expected ordering %v but got %v", tc.order, res.OrderBy)
			}
			if tc.order != nil && !res.OrderBy.Equals(tc.order) {
				t.Fatalf("expected ordering %v but got %v", tc.order, res.OrderBy)
			}

			t.Logf("Len: %d, %s", len(token), token)
		})
	}
}

func BenchmarkTokens(b *testing.B) {

	b.Run("Tokenize", func(b *testing.B) {
		c := expr.Composer{Desc: new(testpb.Message).ProtoReflect().Descriptor()}

		tm := time.Now().Unix()

		for i := 0; i < b.N; i++ {
			next := NextTokenExpr{
				Filter: c.Compare(
					c.MustSelect("timestamp"),
					expr.GT,
					c.FunctionCall("time", "Unix", c.Value(tm)),
				),
				OrderBy: c.OrderBy(
					c.MustOrderByField("timestamp", expr.ASC),
				),
			}

			_, err := TokenizeStruct(&next)
			if err != nil {
				b.Fatalf("Err: %v", err)
			}
			next.Free()
		}
	})

	b.Run("Decode", func(b *testing.B) {
		c := expr.Composer{Desc: new(testpb.Message).ProtoReflect().Descriptor()}

		tm := time.Now().Unix()

		next := NextTokenExpr{
			Filter: c.Compare(
				c.MustSelect("timestamp"),
				expr.GT,
				c.FunctionCall("time", "Unix", c.Value(tm)),
			),
			OrderBy: c.OrderBy(
				c.MustOrderByField("timestamp", expr.ASC),
			),
		}
		defer next.Free()

		token, err := TokenizeStruct(&next)
		if err != nil {
			b.Fatalf("Err: %v", err)
		}

		for i := 0; i < b.N; i++ {
			_, err = DecodeToken[NextTokenExpr](token)
			if err != nil {
				b.Fatalf("Err: %v", err)
			}
		}
	})

}
