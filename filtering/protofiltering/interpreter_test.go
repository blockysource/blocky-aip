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
	"github.com/blockysource/blocky-aip/filtering/token"
	"github.com/blockysource/blocky-aip/internal/testpb"
)

var md = new(testpb.Message).ProtoReflect().Descriptor()

func TestInterpreter_Parse(t *testing.T) {
	tc := []struct {
		name    string
		filter  string
		checkFn func(t *testing.T, x expr.FilterExpr)
		isErr   bool
		err     error
	}{
		{
			name:    "string field EQ direct",
			filter:  tstStringFieldEqDirect,
			checkFn: testStringFieldEqDirect,
		},
		{
			name:    "string field IN array",
			filter:  tstStringFieldInArray,
			checkFn: testStringFieldInArray,
		},
		{
			name:    "string field EQ indirect",
			filter:  tstStringFieldEqIndirect,
			checkFn: testStringFieldEqIndirect,
		},
		{
			name:    "string field EQ string_search",
			filter:  tstStringFieldEqStringSearch,
			checkFn: testStringFieldEqStringSearch,
		},
		{
			name:    "string field EQ string_search with prefix",
			filter:  tstStringFieldEqStringSearchPrefix,
			checkFn: testStringFieldEqStringSearchPrefix,
		},
		{
			name:    "string field EQ string_search with suffix",
			filter:  tstStringFieldEqStringSearchSuffix,
			checkFn: testStringFieldEqStringSearchSuffix,
		},
		{
			name:   "string field invalid value",
			filter: `name = 123`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "repeated string EQ direct",
			filter:  tstRepeatedStringFieldEqDirect,
			checkFn: testRepeatedStringFieldEqDirect,
		},
		{
			name:    "map string int32 field EQ direct",
			filter:  tstMapStringI32FieldEqDirect,
			checkFn: testMapStringI32FieldEqDirect,
		},
		{
			name:    "map string int64 field EQ direct",
			filter:  tstMapStringI64FieldEqDirect,
			checkFn: testMapStringI64FieldEqDirect,
		},
		{
			name:    "map string uint32 field EQ direct",
			filter:  tstMapStringU32FieldEqDirect,
			checkFn: testMapStringU32FieldEqDirect,
		},
		{
			name:    "map string uint64 field EQ direct",
			filter:  tstMapStringU64FieldEqDirect,
			checkFn: testMapStringU64FieldEqDirect,
		},
		{
			name:    "map string sint32 field EQ direct",
			filter:  tstMapStringS32FieldEqDirect,
			checkFn: testMapStringS32FieldEqDirect,
		},
		{
			name:    "map string sint64 field EQ direct",
			filter:  tstMapStringS64FieldEqDirect,
			checkFn: testMapStringS64FieldEqDirect,
		},
		{
			name:    "duration field EQ direct",
			filter:  tstDurationFieldEQDirect,
			checkFn: testDurationFieldEQDirect,
		},
		{
			name:    "duration field EQ indirect",
			filter:  tstDurationFieldEQIndirect,
			checkFn: testDurationFieldEQIndirect,
		},
		{
			name:   "duration field EQ indirect ambiguous",
			filter: "duration = duration",
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "duration field GE direct",
			filter:  tstDurationFieldGEDirect,
			checkFn: testDurationFieldGEDirect,
		},
		{
			name:    "duration field EQ fractal direct",
			filter:  tstDurationFieldEQFractalDirect,
			checkFn: testDurationFieldEQFractalDirect,
		},
		{
			name:    "duration field EQ struct direct",
			filter:  tstDurationFieldEQStructDirect,
			checkFn: testDurationFieldEQStructDirect,
		},
		{
			name:    "duration field IN array direct",
			filter:  tstDurationFieldINArrayDirect,
			checkFn: testDurationFieldINArrayDirect,
		},
		{
			name:    "map string duration map key field HAS direct",
			filter:  tstMapStringDurationFieldHasDirect,
			checkFn: testMapStringDurationFieldHasDirect,
		},
		{
			name:    "repeated duration has direct",
			filter:  tstRepeatedDurationHasDirect,
			checkFn: testRepeatedDurationHasDirect,
		},
		{
			name:   "duration field EQ invalid value",
			filter: `duration = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "duration field EQ invalid duration value",
			filter: `duration = duration{}`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "timestamp field EQ direct",
			filter:  tstTimestampFieldEQDirect,
			checkFn: testTimestampFieldEQDirect,
		},
		{
			name:    "timestamp field EQ indirect",
			filter:  tstTimestampFieldEQIndirect,
			checkFn: testTimestampFieldEQIndirect,
		},
		{
			name:   "timestamp field EQ ambiguous",
			filter: `timestamp = timestamp`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name: "timestamp field IN array direct",
			filter: tstTimestampFieldInArrayDirect,
			checkFn: testTimestampFieldInArrayDirect,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			i, err := NewInterpreter(md, ErrHandlerOpt(errHandler(t, tt.filter, tt.isErr)))
			if err != nil {
				t.Fatal(err)
			}

			x, err := i.Parse(tt.filter)
			if tt.isErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.err != nil && tt.err != err {
					t.Fatalf("expected error %s but got %s", tt.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %s", err)
				}
				defer x.Free()
				tt.checkFn(t, x)
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

func BenchmarkInterpreter_Parse(b *testing.B) {
	it, err := NewInterpreter(md)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		pf, err := it.Parse(tstStringFieldEqDirect)
		if err != nil {
			b.Fatal(err)
		}
		pf.Free()
	}
}
