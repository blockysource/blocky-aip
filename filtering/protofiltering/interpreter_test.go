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

	"google.golang.org/protobuf/reflect/protoreflect"

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
			name:    "timestamp field IN array direct",
			filter:  tstTimestampFieldInArrayDirect,
			checkFn: testTimestampFieldInArrayDirect,
		},
		{
			name:    "i32 field EQ direct",
			filter:  tstI32FieldEQDirect,
			checkFn: testI32FieldEQDirect,
		},
		{
			name:    "i32 field GT direct",
			filter:  tstI32FieldGTDirect,
			checkFn: testI32FieldGTDirect,
		},
		{
			name:    "i32 field EQ indirect",
			filter:  tstI32FieldEQIndirect,
			checkFn: testI32FieldEQIndirect,
		},
		{
			name:   "i32 field EQ invalid string value",
			filter: `i32 = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "i32 field EQ invalid text value",
			filter: `i32 = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "i32 field IN array direct",
			filter:  tstI32FieldInArrayDirect,
			checkFn: testI32FieldInArrayDirect,
		},
		{
			name:    "i32 field IN array indirect",
			filter:  tstI32FieldInArrayIndirect,
			checkFn: testI32FieldInArrayIndirect,
		},
		{
			name:   "i32 field EQ ambiguous",
			filter: `i32 = i32`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "i32 field EQ negative direct",
			filter:  tstI32FieldEQNegativeDirect,
			checkFn: testI32FieldEQNegativeDirect,
		},
		{
			name:    "i64 field EQ direct",
			filter:  tstI64FieldEQDirect,
			checkFn: testI64FieldEQDirect,
		},
		{
			name:    "i64 field EQ indirect",
			filter:  tstI64FieldEQIndirect,
			checkFn: testI64FieldEQIndirect,
		},
		{
			name:    "i64 field LT direct",
			filter:  tstI64FieldLTDirect,
			checkFn: testI64FieldLTDirect,
		},
		{
			name:    "i64 field IN array direct",
			filter:  tstI64FieldInArrayDirect,
			checkFn: testI64FieldInArrayDirect,
		},
		{
			name:    "i64 field IN array indirect",
			filter:  tstI64FieldInArrayIndirect,
			checkFn: testI64FieldInArrayIndirect,
		},
		{
			name:    "i64 field EQ negative direct",
			filter:  tstI64FieldEQNegativeDirect,
			checkFn: testI64FieldEQNegativeDirect,
		},
		{
			name:   "i64 field EQ ambiguous",
			filter: `i64 = i64`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "s32 field EQ direct",
			filter:  tstS32FieldEQDirect,
			checkFn: testS32FieldEQDirect,
		},
		{
			name:    "s32 field EQ indirect",
			filter:  tstS32FieldEQIndirect,
			checkFn: testS32FieldEQIndirect,
		},
		{
			name:    "s32 in array direct",
			filter:  tstS32FieldInArrayDirect,
			checkFn: testS32FieldInArrayDirect,
		},
		{
			name:    "s32 in array indirect",
			filter:  tstS32FieldInArrayIndirect,
			checkFn: testS32FieldInArrayIndirect,
		},
		{
			name:   "s32 field EQ ambiguous",
			filter: `s32 = s32`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:   "s32 field EQ invalid string value",
			filter: `s32 = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "s32 field EQ negative direct",
			filter:  tstS32FieldEQNegativeDirect,
			checkFn: testS32FieldEQNegativeDirect,
		},
		{
			name:    "s64 field EQ direct",
			filter:  tstS64FieldEQDirect,
			checkFn: testS64FieldEQDirect,
		},
		{
			name:    "s64 field EQ indirect",
			filter:  tstS64FieldEQIndirect,
			checkFn: testS64FieldEQIndirect,
		},
		{
			name:    "s64 field IN array direct",
			filter:  tstS64FieldInArrayDirect,
			checkFn: testS64FieldInArrayDirect,
		},
		{
			name:    "s64 field IN array indirect",
			filter:  tstS64FieldInArrayIndirect,
			checkFn: testS64FieldInArrayIndirect,
		},
		{
			name:   "s64 field EQ ambiguous",
			filter: `s64 = s64`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:   "s64 field EQ invalid string value",
			filter: `s64 = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "s64 field EQ negative direct",
			filter:  tstS64FieldEQNegativeDirect,
			checkFn: testS64FieldEQNegativeDirect,
		},
		{
			name:    "u32 field EQ direct",
			filter:  tstU32FieldEQDirect,
			checkFn: testU32FieldEQDirect,
		},
		{
			name:    "u32 field EQ indirect",
			filter:  tstU32FieldEQIndirect,
			checkFn: testU32FieldEQIndirect,
		},
		{
			name:    "u32 field IN array direct",
			filter:  tstU32FieldINArray,
			checkFn: testU32FieldINArray,
		},
		{
			name:   "u32 field EQ invalid string value",
			filter: `u32 = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "u32 field EQ invalid text value",
			filter: `u32 = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "u32 field EQ ambiguous",
			filter: `u32 = u32`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "u32 field IN array indirect",
			filter:  tstU32FieldINArrayIndirect,
			checkFn: testU32FieldINArrayIndirect,
		},
		{
			name:    "u64 field EQ direct",
			filter:  tstU64FieldEQDirect,
			checkFn: testU64FieldEQDirect,
		},
		{
			name:    "u64 field EQ indirect",
			filter:  tstU64FieldEQIndirect,
			checkFn: testU64FieldEQIndirect,
		},
		{
			name:    "u64 field IN array direct",
			filter:  tstU64FieldINArray,
			checkFn: testU64FieldINArray,
		},
		{
			name:   "u64 field EQ invalid string value",
			filter: `u64 = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "u64 field EQ invalid text value",
			filter: `u64 = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "u64 field EQ ambiguous",
			filter: `u64 = u64`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "u64 field IN array indirect",
			filter:  tstU64FieldINArrayIndirect,
			checkFn: testU64FieldINArrayIndirect,
		},
		{
			name:    "u64 field EQ indirect u32",
			filter:  `u64 = u32`,
			checkFn: testIndirectFields("u64", "u32"),
		},
		{
			name:    "u64 field EQ indirect i32",
			filter:  `u64 = i32`,
			checkFn: testIndirectFields("u64", "i32"),
		},
		{
			name:    "u64 field EQ indirect s32",
			filter:  `u64 = s32`,
			checkFn: testIndirectFields("u64", "s32"),
		},
		{
			name:    "u64 field EQ indirect s64",
			filter:  `u64 = s64`,
			checkFn: testIndirectFields("u64", "s64"),
		},
		{
			name:    "u64 field EQ indirect i64",
			filter:  `u64 = i64`,
			checkFn: testIndirectFields("u64", "i64"),
		},
		{
			name:    "float field EQ indirect double",
			filter:  `float = double`,
			checkFn: testIndirectFields("float", "double"),
		},
		{
			name:    "double field EQ indirect float",
			filter:  `double = float`,
			checkFn: testIndirectFields("double", "float"),
		},
		{
			name:    "float field EQ direct",
			filter:  tstFloatFieldEQDirect,
			checkFn: testFloatFieldEQDirect,
		},
		{
			name:    "float field EQ indirect",
			filter:  tstFloatFieldEQIndirect,
			checkFn: testFloatFieldEQIndirect,
		},
		{
			name:   "float field EQ indirect ambiguous",
			filter: `float = float`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:   "float field EQ invalid value",
			filter: `float = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "float field EQ invalid text value",
			filter: `float = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "float field IN array direct",
			filter:  tstFloatFieldINArrayDirect,
			checkFn: testFloatFieldINArrayDirect,
		},
		{
			name:    "float field IN array indirect",
			filter:  tstFloatFieldINArrayIndirect,
			checkFn: testFloatFieldINArrayIndirect,
		},
		{
			name:    "float field EQ negative direct",
			filter:  tstFloadFieldEQNegative,
			checkFn: testFloadFieldEQNegative,
		},
		{
			name:    "double field EQ direct",
			filter:  tstDoubleFieldEQDirect,
			checkFn: testDoubleFieldEQDirect,
		},
		{
			name:    "double field EQ indirect",
			filter:  tstDoubleFieldEQIndirect,
			checkFn: testDoubleFieldEQIndirect,
		},
		{
			name:   "double field EQ invalid value",
			filter: `double = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "double field EQ invalid text value",
			filter: `double = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "double field IN array direct",
			filter:  tstDoubleFieldINArrayDirect,
			checkFn: testDoubleFieldINArrayDirect,
		},
		{
			name:    "double field IN array indirect",
			filter:  tstDoubleFieldINArrayIndirect,
			checkFn: testDoubleFieldINArrayIndirect,
		},
		{
			name:    "double field EQ negative direct",
			filter:  tstDoubleFieldEQNegative,
			checkFn: testDoubleFieldEQNegative,
		},
		{
			name:    "enum field EQ direct",
			filter:  tstEnumFieldEQDirect,
			checkFn: testEnumFieldEQDirect,
		},
		{
			name:    "enum field EQ indirect",
			filter:  tstEnumFieldEQIndirect,
			checkFn: testEnumFieldEQIndirect,
		},
		{
			name:   "enum field EQ invalid value",
			filter: `enum = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "enum field EQ invalid text value",
			filter: `enum = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "enum field IN array direct",
			filter:  tstEnumFieldInArrayDirect,
			checkFn: testEnumFieldInArrayDirect,
		},
		{
			name:    "enum field IN array indirect",
			filter:  tstEnumFieldInArrayIndirect,
			checkFn: testEnumFieldInArrayIndirect,
		},
		{
			name:   "enum field ambiguous",
			filter: `enum = enum`,
			isErr:  true,
			err:    ErrAmbiguousField,
		},
		{
			name:    "bool field EQ direct true",
			filter:  tstBoolFieldEQDirectTrue,
			checkFn: testBoolFieldEQDirectTrue,
		},
		{
			name:    "bool field EQ direct false",
			filter:  tstBoolFieldEQDirectFalse,
			checkFn: testBoolFieldEQDirectFalse,
		},
		{
			name:   "bool field EQ invalid value",
			filter: `bool = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "bool field EQ invalid text value",
			filter: `bool = 0`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "bool field EQ indirect",
			filter:  tstBoolFieldEQIndirect,
			checkFn: testBoolFieldEQIndirect,
		},
		{
			name:    "structpb field EQ direct string",
			filter:  tstStructPbFieldEQDirectString,
			checkFn: testStructPbFieldEQDirectString,
		},
		{
			name:    "structpb field EQ direct message",
			filter:  tstStructPbFieldEQDirectMessage,
			checkFn: testStructPbFieldEQDirectMessage,
		},
		{
			name:    "message field EQ direct",
			filter:  tstMsgFieldEQDirect,
			checkFn: testMsgFieldEQDirect,
		},
		{
			name:    "message field EQ indirect",
			filter:  tstMsgFieldEQIndirect,
			checkFn: testMsgFieldEQIndirect,
		},
		{
			name:    "message field EQ direct unnamed",
			filter:  tstMsgFieldEQDirectUnnamed,
			checkFn: testMsgFieldEQDirectUnnamed,
		},
		{
			name:   "message field EQ invalid struct name",
			filter: `sub = invalid.Name{}`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "message field EQ invalid struct field",
			filter: `sub = {invalid: 1}`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "message field EQ invalid string literal",
			filter: `sub = "invalid"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "message field EQ invalid text literal",
			filter: `sub = some_text`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "string field EQ direct forbidden filtering",
			filter: `no_filter = "test"`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "message field EQ direct forbidden filtering",
			filter: `no_filter_msg = {}`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:   "message field EQ sub selector forbidden filtering",
			filter: `no_filter_msg.i64 = 1`,
			isErr:  true,
			err:    ErrInvalidValue,
		},
		{
			name:    "i32_complexity field EQ direct",
			filter:  tstI32ComplexityEQDirect,
			checkFn: testI32ComplexityEQDirect,
		},
		{
			name:    "multiple expressions with AND",
			filter:  tstMultipleExpressions,
			checkFn: testMultipleExpressions,
		},
		{
			name:    "term with NOT keyword",
			filter:  tstTermWithNOTKeyword,
			checkFn: testTermWithNOTKeyword,
		},
		{
			name:    "term with NOT keyword and parentheses",
			filter:  tstTermWithNOTKeywordAndParentheses,
			checkFn: testTermWithNOTKeywordAndParentheses,
		},
		{
			name:    "factor with multiple ors",
			filter:  tstFactorWithMultipleORs,
			checkFn: testFactorWithMultipleORs,
		},
		{
			name:    "complex expression",
			filter:  tstComplexExpression,
			checkFn: testComplexExpression,
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
				if testing.Verbose() {
					t.Logf("Query: '%s' complexity: %d", tt.filter, x.Complexity())
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

func testIndirectFields(f1, f2 string) func(t *testing.T, x expr.FilterExpr) {
	return func(t *testing.T, x expr.FilterExpr) {
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

		if left.Field != md.Fields().ByName(protoreflect.Name(f1)) {
			t.Fatalf("expected field '%s' field but got %s", f1, left.Field)
		}

		right, ok := ce.Right.(*expr.FieldSelectorExpr)
		if !ok {
			t.Fatalf("expected value expression but got %T", ce.Right)
		}

		if right.Field != md.Fields().ByName(protoreflect.Name(f2)) {
			t.Fatalf("expected field '%s' field but got %s", f2, right.Field)
		}
	}
}
