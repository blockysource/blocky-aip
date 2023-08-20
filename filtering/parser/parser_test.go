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

package parser

import (
	"testing"

	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/filtering/scanner"
	"github.com/blockysource/blocky-aip/filtering/token"
)

func memberTextLiteral(t *testing.T, m *ast.MemberExpr, expected string, pos int) {
	if m.Value == nil {
		t.Fatal("expected member value")
	}

	tl, ok := m.Value.(*ast.TextLiteral)
	if !ok {
		t.Fatalf("expected text literal got: %v", m.Value)
	}

	if tl.Value != expected {
		t.Fatalf("expected '%s' got: %v", expected, tl.Value)
	}

	if tl.Pos != token.Position(pos) {
		t.Fatalf("expected position %d got: %v", pos, tl.Pos)
	}
}

// TestParse tests the Parse function.
func TestParse(t *testing.T) {

	testCases := []struct {
		name            string
		src             string
		useStructs      bool
		useArray        bool
		useInComparator bool
		expectedErr     error
		checkFn         func(t *testing.T, pf *ParsedFilter)
	}{
		{
			name: "empty",
			src:  "",
			checkFn: func(t *testing.T, pf *ParsedFilter) {
				if pf.Expr != nil {
					t.Errorf("expected nil expression")
				}
			},
		},
		{
			name:    "single sequence",
			src:     singleSequenceMember,
			checkFn: testSingleSequenceMember,
		},
		{
			name:    "single sequence with string",
			src:     singleSequenceWithStringMember,
			checkFn: testSingleSequenceWithStringMember,
		},
		{
			name:    "single sequence with unary op",
			src:     singleSequenceWithUnaryOp,
			checkFn: testSingleSequenceWithUnaryOp,
		},
		{
			name:    "single sequence with unary op and string",
			src:     singleSequenceWithUnaryOpAndString,
			checkFn: testSingleSequenceWithUnaryOpAndString,
		},
		{
			name:    "single sequence with unary NOT op",
			src:     singleSequenceWithUnaryNotOp,
			checkFn: testSingleSequenceWithUnaryNotOp,
		},
		{
			name:    "restriction with comparator",
			src:     restrictionWithEQ,
			checkFn: testRestrictionWithEQ,
		},
		{
			name:    "restriction with ge comparator",
			src:     restrictionWithGE,
			checkFn: testRestrictionWithGE,
		},
		{
			name:    "restriction with ne comparator",
			src:     restrictionWithNE,
			checkFn: testRestrictionWithNE,
		},
		{
			name:    "restriction with no space and comparator",
			src:     restrictionWithNoSpaceAndComparator,
			checkFn: testRestrictionWithNoSpaceAndComparator,
		},
		{
			name:    "restriction with function arg",
			src:     restrictionWithFunctionArg,
			checkFn: testRestrictionWithFunctionArg,
		},
		{
			name:    "restriction with function arg list",
			src:     restrictionWithFunctionArgList,
			checkFn: testRestrictionWithFunctionArgList,
		},
		{
			name:    "restriction with function arg list no space",
			src:     restrictionWithFunctionArgListNoSpace,
			checkFn: testRestrictionWithFunctionArgListNoSpace,
		},
		{
			name:    "restriction with has arg",
			src:     restrictionWithHasArg,
			checkFn: testRestrictionWithHasArg,
		},
		{
			name:    "restriction with string arg",
			src:     restrictionWithStringArg,
			checkFn: testRestrictionWithStringArg,
		},
		{
			name:    "restriction with single quoted string arg",
			src:     restrictionWithSingleQuotedStringArg,
			checkFn: testRestrictionWithSingleQuotedStringArg,
		},
		{
			name:    "restriction with alternative not has string arg",
			src:     restrictionWithAlternativeNotHastStringArg,
			checkFn: testRestrictionWithAlternativeNotHastStringArg,
		},
		{
			name:    "factors with or",
			src:     factorsWithOR,
			checkFn: testFactorsWithOR,
		},
		{
			name:    "sequence with factors",
			src:     sequenceWithFactors,
			checkFn: testSequenceWithFactors,
		},
		{
			name:    "composite expression",
			src:     compositeExpression,
			checkFn: testCompositeExpression,
		},
		{
			name:    "complex func call",
			src:     complexFuncCall,
			checkFn: testComplexFuncCall,
		},
		{
			name:    "deep nested member",
			src:     deepNestedText,
			checkFn: testDeepNestedTextMember,
		},
		{
			name:    "complex",
			src:     complexExpr,
			checkFn: testComplexExpr,
		},
		{
			name:    "multi whitespaces",
			src:     exprMultiWhiteSpace,
			checkFn: testExprMultiWhiteSpace,
		},
		{
			name:    "deep nested string member",
			src:     deepNestedStringMember,
			checkFn: testDeepNestedStringMember,
		},
		{
			name:    "func call no arg",
			src:     funcCallNoArg,
			checkFn: testFuncCallNoArg,
		},
		{
			name:     "array with quote",
			useArray: true,
			src:      arrayWithQuote,
			checkFn:  testArrayWithQuote,
		},
		{
			name:     "array with quote and WS",
			src:      arrayWithQuoteAndWS,
			checkFn:  testArrayWithQuoteAndWS,
			useArray: true,
		},
		{
			name:       "struct extension",
			src:        structExpr,
			useStructs: true,
			checkFn:    testStructExpr,
		},
		{
			name:       "restriction with struct arg",
			src:        restrictionWithStructArg,
			useStructs: true,
			checkFn:    testRestrictionWithStructArg,
		},
		{
			name:       "complex restriction func call with struct and array",
			src:        complexRestrictionWithFuncCallStructAndArray,
			useArray:   true,
			useStructs: true,
			checkFn:    testComplexRestrictionWithFuncCallStructAndArray,
		},
		{
			name:       "struct with newlines",
			useStructs: true,
			src:        structExprWithNewLines,
			checkFn:    testStructExprWithNewLines,
		},
		{
			name:       "struct with newlines ended with comma",
			useStructs: true,
			src:        structExprWithNewLinesEndedWithComma,
			checkFn:    testStructExprWithNewLinesEndedWithComma,
		},
		{
			name:            "restriction with IN and array",
			useArray:        true,
			useInComparator: true,
			src:             restrictionWithIN,
			checkFn:         testRestrictionWithIN,
		},
		{
			name:    "restriction with timestamp",
			src:     restrictionWithTimestamp,
			checkFn: testRestrictionWithTimestamp,
		},
		{
			name:    "restriction with timestamp and has comparator",
			src:     restrictionWithTimestampAndHas,
			checkFn: testRestrictionWithTimestampAndHas,
		},
		{
			name:    "restriction with timestamp and timezone",
			src:     restrictionWithTimestampAndTimezone,
			checkFn: testRestrictionWithTimestampAndTimezone,
		},
		{
			name:       "map struct comparable",
			src:        mapStructComparable,
			useStructs: true,
			checkFn:    testMapStructComparable,
		},
		{
			name:    "restriction with negative int on right",
			src:     restrictionWithNegativeIntOnRight,
			checkFn: testRestrictionWithNegativeIntOnRight,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			opts := []ParserOption{
				ErrorHandlerOption(testErrHandler(t)),
			}
			if tc.useStructs {
				opts = append(opts, UseStructsOption)
			}
			if tc.useArray {
				opts = append(opts, UseArraysOption)
			}
			if tc.useInComparator {
				opts = append(opts, UseInComparatorOption)
			}
			p := NewParser(tc.src, opts...)

			pf, err := p.Parse()
			if err != nil {
				if tc.expectedErr != nil {
					if tc.expectedErr != err {
						t.Fatalf("expected error: %s got: %s", tc.expectedErr, err)
					}
					return
				} else {
					t.Fatalf("unexpected error: %s", err)
				}
			}
			defer pf.Free()

			tc.checkFn(t, pf)
		})
	}
}

func testErrHandler(t testing.TB) scanner.ErrorHandler {
	return func(pos token.Position, msg string) {
		t.Errorf("unexpected error at %d: %s", pos, msg)
	}
}

func BenchmarkParse(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		p := Parser{}
		for i := 0; i < b.N; i++ {
			p.Reset("a")
			pf, err := p.Parse()
			if err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			pf.Free()
		}
	})

	b.Run("Complex", func(b *testing.B) {
		p := Parser{}
		for i := 0; i < b.N; i++ {
			p.Reset("(a b) AND c OR d AND (e > f OR g < h)")
			pf, err := p.Parse()
			if err != nil {
				b.Fatalf("unexpected error: %s", err)
			}
			pf.Free()
		}
	})
}

func seqMember(t *testing.T, seq *ast.SequenceExpr) *ast.MemberExpr {
	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor")
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one term")
	}

	term := factor.Terms[0]
	if term.UnaryOp != "" {
		t.Errorf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	member, ok := expr.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	return member
}

func factorMember(t *testing.T, factor *ast.FactorExpr) *ast.MemberExpr {
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one term")
	}

	term := factor.Terms[0]
	if term.UnaryOp != "" {
		t.Errorf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	member, ok := expr.Comparable.(*ast.MemberExpr)
	if !ok {
		t.Fatalf("expected member literal")
	}
	return member
}

func seqFuncCall(t *testing.T, seq *ast.SequenceExpr) *ast.FunctionCall {
	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor")
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one term")
	}

	term := factor.Terms[0]
	if term.UnaryOp != "" {
		t.Errorf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	if expr.Comparable == nil {
		t.Fatal("expected comparable")
	}

	fnCall, ok := expr.Comparable.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("expected function call, got: %T", expr.Comparable)
	}
	return fnCall
}

func seqRestriction(t *testing.T, seq *ast.SequenceExpr) *ast.RestrictionExpr {
	if len(seq.Factors) != 1 {
		t.Fatalf("expected one factor, got: %v", len(seq.Factors))
	}

	factor := seq.Factors[0]
	if len(factor.Terms) != 1 {
		t.Fatalf("expected one term")
	}

	term := factor.Terms[0]
	if term.UnaryOp != "" {
		t.Errorf("expected no unary op")
	}

	if term.Expr == nil {
		t.Fatalf("expected expression")
	}

	expr, ok := term.Expr.(*ast.RestrictionExpr)
	if !ok {
		t.Fatalf("expected restriction expression")
	}

	return expr
}
