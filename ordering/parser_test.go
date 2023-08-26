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

package ordering

import (
	"testing"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/testpb"
	"github.com/blockysource/blocky-aip/token"
)

var md = new(testpb.Message).ProtoReflect().Descriptor()

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		wantErr bool
		err     error
		checkFn func(t *testing.T, e *expr.OrderByExpr)
	}{
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
			err:     ErrInvalidSyntax,
		},
		{
			name:  "single field",
			input: "i64",
			checkFn: func(t *testing.T, e *expr.OrderByExpr) {
				if len(e.Fields) != 1 {
					t.Fatalf("expected 1 field but got %d", len(e.Fields))
				}

				ofe := e.Fields[0]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe := ofe.Field
				if fe.Field != md.Fields().ByName("i64").Name() {
					t.Fatalf("expected field 'i64' field but got %s", fe.Field)
				}

				if ofe.Order != expr.ASC {
					t.Fatalf("expected order %s but got %s", expr.ASC, ofe.Order)
				}
			},
		},
		{
			name:  "multiple fields",
			input: "i64, float",
			checkFn: func(t *testing.T, e *expr.OrderByExpr) {
				if len(e.Fields) != 2 {
					t.Fatalf("expected 2 fields but got %d", len(e.Fields))
				}

				ofe := e.Fields[0]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe := ofe.Field
				if fe.Field != md.Fields().ByName("i64").Name() {
					t.Fatalf("expected field 'i64' field but got %s", fe.Field)
				}

				if ofe.Order != expr.ASC {
					t.Fatalf("expected order %s but got %s", expr.ASC, ofe.Order)
				}

				ofe = e.Fields[1]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe = ofe.Field
				if fe.Field != md.Fields().ByName("float").Name() {
					t.Fatalf("expected field 'float' field but got %s", fe.Field)
				}

				if ofe.Order != expr.ASC {
					t.Fatalf("expected order %s but got %s", expr.ASC, ofe.Order)
				}
			},
		},
		{
			name:  "multiple fields with order",
			input: `i64 ASC, float DESC`,
			checkFn: func(t *testing.T, e *expr.OrderByExpr) {
				if len(e.Fields) != 2 {
					t.Fatalf("expected 2 fields but got %d", len(e.Fields))
				}

				ofe := e.Fields[0]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe := ofe.Field
				if fe.Field != md.Fields().ByName("i64").Name() {
					t.Fatalf("expected field 'i64' field but got %s", fe.Field)
				}

				if ofe.Order != expr.ASC {
					t.Fatalf("expected order %s but got %s", expr.ASC, ofe.Order)
				}

				ofe = e.Fields[1]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe = ofe.Field
				if fe.Field != md.Fields().ByName("float").Name() {
					t.Fatalf("expected field 'float' field but got %s", fe.Field)
				}

				if ofe.Order != expr.DESC {
					t.Fatalf("expected order %s but got %s", expr.DESC, ofe.Order)
				}
			},
		},
		{
			name:    "invalid field",
			input:   "invalid",
			wantErr: true,
			err:     ErrInvalidField,
		},
		{
			name:  "sub field",
			input: `sub.enum`,
			checkFn: func(t *testing.T, e *expr.OrderByExpr) {
				if len(e.Fields) != 1 {
					t.Fatalf("expected 1 field but got %d", len(e.Fields))
				}

				ofe := e.Fields[0]
				if ofe.Field == nil {
					t.Fatalf("expected field but got nil")
				}

				fe := ofe.Field
				if fe.Field != md.Fields().ByName("sub").Name() {
					t.Fatalf("expected field 'sub' field but got %s", fe.Field)
				}

				if fe.Traversal == nil {
					t.Fatalf("expected traversal but got nil")
				}

				var ok bool
				fe, ok = fe.Traversal.(*expr.FieldSelectorExpr)
				if !ok {
					t.Fatalf("expected field selector but got %T", fe.Traversal)
				}

				if fe.Field != md.Fields().ByName("enum").Name() {
					t.Fatalf("expected field 'enum' field but got %s", fe.Field)
				}

				if ofe.Order != expr.ASC {
					t.Fatalf("expected order %s but got %s", expr.ASC, ofe.Order)
				}
			},
		},
		{
			name:    "sorting forbidden",
			input:   "no_filter",
			err:     ErrSortingForbidden,
			wantErr: true,
		},
		{
			name:    "sorting sub field forbidden",
			input:   "sub.no_filter",
			err:     ErrSortingForbidden,
			wantErr: true,
		},
		{
			name:    "traversal field forbidden",
			input:   "no_filter_msg.i64",
			err:     ErrSortingForbidden,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Parser
			if err := p.Reset(md, ErrHandler(testErrHandler(t, tt.wantErr))); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := p.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFn != nil {
				tt.checkFn(t, got)
			}
		})
	}
}

func testErrHandler(t testing.TB, isErr bool) func(pos token.Position, msg string) {
	return func(pos token.Position, msg string) {
		if !isErr {
			t.Errorf("unexpected error: %s", msg)
		}
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	benchs := []struct {
		name     string
		input    string
		wantsErr bool
	}{
		{
			name:  "single field",
			input: `i64`,
		},
		{
			name:  "multiple fields",
			input: "i64 DESC, float ASC, sub.enum DESC",
		},
		{
			name:     "forbidden field",
			input:    "no_filter DESC",
			wantsErr: true,
		},
	}

	for _, bc := range benchs {
		b.Run(bc.name, func(b *testing.B) {
			var p Parser
			err := p.Reset(md, ErrHandler(testErrHandler(b, bc.wantsErr)))
			if err != nil {
				b.Fatal(err)
			}
			var x *expr.OrderByExpr
			for i := 0; i < b.N; i++ {
				x, err = p.Parse(bc.input)
				if bc.wantsErr && err == nil {
					b.Fatalf("expected error but got nil")
				}

				if !bc.wantsErr {
					x.Free()
				}
			}
		})
	}
}
