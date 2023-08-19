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
	"sync"

	"github.com/blockysource/blocky-aip/filtering/ast"
)

var (
	// parsedFilterPool is a pool of parsed filters.
	parsedFilterPool = sync.Pool{
		New: func() any {
			return &ParsedFilter{}
		},
	}
	// exprPool is a pool of AST expressions.
	exprPool = sync.Pool{
		New: func() any {
			return &ast.Expr{
				Sequences: make([]*ast.SequenceExpr, 0, 10),
			}
		},
	}
	sequenceExprPool = sync.Pool{
		New: func() any {
			return &ast.SequenceExpr{
				Factors: make([]*ast.FactorExpr, 0, 10),
			}
		},
	}

	factorExprPool = sync.Pool{
		New: func() any {
			return &ast.FactorExpr{
				Terms: make([]*ast.TermExpr, 0, 10),
			}
		},
	}

	termExprPool = sync.Pool{
		New: func() any { return &ast.TermExpr{} },
	}

	restrictionExprPool = sync.Pool{
		New: func() any { return &ast.RestrictionExpr{} },
	}

	memberLiteralPool = sync.Pool{
		New: func() any {
			return &ast.MemberExpr{
				Fields: make([]ast.FieldExpr, 0, 10),
			}
		},
	}

	funcCallPool = sync.Pool{
		New: func() any {
			fc := &ast.FunctionCall{
				Name: make([]ast.NameExpr, 0, 10),
			}
			return fc
		},
	}

	comparatorLiteralPool = sync.Pool{
		New: func() any { return &ast.ComparatorLiteral{} },
	}

	compositeExprPool = sync.Pool{
		New: func() any { return &ast.CompositeExpr{} },
	}

	textLiteralPool = sync.Pool{
		New: func() any {
			tl := &ast.TextLiteral{}
			return tl
		},
	}

	stringLiteralPool = sync.Pool{
		New: func() any {
			sl := &ast.StringLiteral{}
			return sl
		},
	}

	argListExprPool = sync.Pool{
		New: func() any {
			return &ast.ArgListExpr{
				Args: make([]ast.ArgExpr, 0, 10),
			}
		},
	}

	keywordExprPool = sync.Pool{
		New: func() any { return &ast.KeywordExpr{} },
	}
)

func getParsedFilter() *ParsedFilter {
	return parsedFilterPool.Get().(*ParsedFilter)
}

func putParsedFilter(f *ParsedFilter) {
	if f == nil {
		return
	}
	putExpr(f.Expr)
	f.Expr = nil
}

func getExpr() *ast.Expr {
	return exprPool.Get().(*ast.Expr)
}

func putExpr(e *ast.Expr) {
	if e == nil {
		return
	}
	e.Pos = 0
	for _, v := range e.Sequences {
		putSequenceExpr(v)
	}
	e.Sequences = e.Sequences[:0]
	exprPool.Put(e)
}

func putSequenceExpr(e *ast.SequenceExpr) {
	if e == nil {
		return
	}
	e.Pos = 0
	for _, v := range e.Factors {
		putFactorExpr(v)
	}
	e.Factors = e.Factors[:0]
	sequenceExprPool.Put(e)
}

func getSequenceExpr() *ast.SequenceExpr {
	return sequenceExprPool.Get().(*ast.SequenceExpr)
}

func putFactorExpr(e *ast.FactorExpr) {
	if e == nil {
		return
	}
	for _, v := range e.Terms {
		putTermExpr(v)
	}
	e.Terms = e.Terms[:0]
	e.Pos = 0
	factorExprPool.Put(e)
}

func getFactorExpr() *ast.FactorExpr {
	return factorExprPool.Get().(*ast.FactorExpr)
}

func putTermExpr(e *ast.TermExpr) {
	if e == nil {
		return
	}
	putSimpleExpr(e.Expr)
	e.Expr = nil
	e.UnaryOp = ""
	e.Pos = 0
	termExprPool.Put(e)
}

func getTermExpr() *ast.TermExpr {
	return termExprPool.Get().(*ast.TermExpr)
}

func putSimpleExpr(e ast.SimpleExpr) {
	if e == nil {
		return
	}
	switch vt := e.(type) {
	case *ast.CompositeExpr:
		putCompositeLiteral(vt)
	case *ast.RestrictionExpr:
		putRestrictionExpr(vt)
	}
}

func putRestrictionExpr(e *ast.RestrictionExpr) {
	if e == nil {
		return
	}
	e.Pos = 0
	putComparableExpr(e.Comparable)
	e.Comparable = nil
	putComparatorLiteral(e.Comparator)
	e.Comparator = nil
	putArgExpr(e.Arg)
	e.Arg = nil
}

func getRestrictionExpr() *ast.RestrictionExpr {
	return restrictionExprPool.Get().(*ast.RestrictionExpr)
}

func putMemberLiteral(e *ast.MemberExpr) {
	if e == nil {
		return
	}

	putValueExpr(e.Value)
	e.Value = nil

	for _, v := range e.Fields {
		putFieldExpr(v)
	}
	e.Fields = e.Fields[:0]
	memberLiteralPool.Put(e)
}

func getMemberLiteral() *ast.MemberExpr {
	return memberLiteralPool.Get().(*ast.MemberExpr)
}

func putFunctionLiteral(e *ast.FunctionCall) {
	if e == nil {
		return
	}
	e.Pos = 0
	for _, v := range e.Name {
		putNameExpr(v)
	}
	e.Name = e.Name[:0]

	e.Lparen = 0
	putArgListExpr(e.ArgList)
	e.ArgList = nil
	e.Rparen = 0
	funcCallPool.Put(e)
}

func getFunctionCall() *ast.FunctionCall {
	return funcCallPool.Get().(*ast.FunctionCall)
}

func putComparatorLiteral(e *ast.ComparatorLiteral) {
	if e == nil {
		return
	}
	e.Pos = 0
	e.Type = 0
	comparatorLiteralPool.Put(e)
}

func getComparatorLiteral() *ast.ComparatorLiteral {
	return comparatorLiteralPool.Get().(*ast.ComparatorLiteral)
}

func putCompositeLiteral(e *ast.CompositeExpr) {
	if e == nil {
		return
	}
	e.Lparen = 0
	e.Rparen = 0

	putExpr(e.Expr)
	e.Expr = nil
}

func getCompositeExpr() *ast.CompositeExpr {
	return compositeExprPool.Get().(*ast.CompositeExpr)
}

func putValueExpr(e ast.ValueExpr) {
	switch vt := e.(type) {
	case *ast.TextLiteral:
		putTextLiteral(vt)
	case *ast.StringLiteral:
		putStringLiteral(vt)
	}
}

func putTextLiteral(e *ast.TextLiteral) {
	if e == nil {
		return
	}
	e.Pos = 0
	e.Value = ""
	e.IsTimestamp = false
	textLiteralPool.Put(e)
}

func getTextLiteral() *ast.TextLiteral {
	return textLiteralPool.Get().(*ast.TextLiteral)
}

func putStringLiteral(e *ast.StringLiteral) {
	if e == nil {
		return
	}
	e.Pos = 0
	e.Value = ""
	stringLiteralPool.Put(e)
}

func getStringLiteral() *ast.StringLiteral {
	return stringLiteralPool.Get().(*ast.StringLiteral)
}

func putArgExpr(e ast.ArgExpr) {
	if e == nil {
		return
	}

	switch vt := e.(type) {
	case *ast.CompositeExpr:
		putCompositeLiteral(vt)
	case ast.ComparableExpr:
		putComparableExpr(vt)
	}
}

func putArgListExpr(e *ast.ArgListExpr) {
	if e == nil {
		return
	}
	for _, v := range e.Args {
		putArgExpr(v)
	}
	e.Args = e.Args[:0]
	argListExprPool.Put(e)
}

func getArgListExpr() *ast.ArgListExpr {
	return argListExprPool.Get().(*ast.ArgListExpr)
}

func putNameExpr(e ast.NameExpr) {
	if e == nil {
		return
	}
	switch vt := e.(type) {
	case *ast.TextLiteral:
		putTextLiteral(vt)
	case *ast.KeywordExpr:
		putKeywordExpr(vt)
	}
}

func putKeywordExpr(e *ast.KeywordExpr) {
	if e == nil {
		return
	}
	e.Pos = 0
	e.Typ = 0
	keywordExprPool.Put(e)
}

func getKeywordExpr() *ast.KeywordExpr {
	return keywordExprPool.Get().(*ast.KeywordExpr)
}

func putFieldExpr(e ast.FieldExpr) {
	if e == nil {
		return
	}
	switch vt := e.(type) {
	case *ast.TextLiteral:
		putTextLiteral(vt)
	case *ast.KeywordExpr:
		putKeywordExpr(vt)
	}
}
