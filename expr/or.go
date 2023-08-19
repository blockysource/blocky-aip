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

package expr

import (
	"sync"
)

var orExprPool = &sync.Pool{
	New: func() any {
		return &OrExpr{
			Expr:       make([]FilterExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireOrExpr acquires an OrExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireOrExpr() *OrExpr {
	return orExprPool.Get().(*OrExpr)
}

// Free puts the OrExpr back to the pool.
func (e *OrExpr) Free() {
	if e == nil {
		return
	}
	for _, sub := range e.Expr {
		if sub != nil {
			sub.Free()
		}
	}
	if !e.isAcquired {
		return
	}
	e.Expr = e.Expr[:0]
	orExprPool.Put(e)
}

var _ FilterExpr = (*OrExpr)(nil)

// OrExpr is an expression that envelops multiple expressions
// into a logical OR group.
type OrExpr struct {
	// Expr is a list of expressions to be evaluated.
	Expr []FilterExpr

	isAcquired bool
}

// Complexity of the OrExpr is the product of complexities of the inner expressions + 1.
func (e *OrExpr) Complexity() int64 {
	var complexity int64
	for _, expr := range e.Expr {
		complexity *= expr.Complexity()
	}
	return complexity + 1
}

// isFilterExpr is a marker method for expressions.
func (e *OrExpr) isFilterExpr() {}
