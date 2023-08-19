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

var andExprPool = &sync.Pool{
	New: func() any {
		return &AndExpr{
			Expr:       make([]FilterExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireAndExpr acquires an AndExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireAndExpr() *AndExpr {
	return andExprPool.Get().(*AndExpr)
}

// Free puts the AndExpr back to the pool.
func (e *AndExpr) Free() {
	if e == nil {
		return
	}
	for _, sub := range e.Expr {
		if sub != nil {
			sub.Free()
		}
	}
	if e.isAcquired {
		e.Expr = e.Expr[:0]
		andExprPool.Put(e)
	}
}

var _ FilterExpr = (*AndExpr)(nil)

// AndExpr is an expression that can be evaluated.
type AndExpr struct {
	// Expr is a list of expressions that should be evaluated with AND operator.
	Expr       []FilterExpr
	isAcquired bool
}

// Complexity of the AndExpr is the sum of complexities of the inner expressions + 1.
func (e *AndExpr) Complexity() int64 {
	var complexity int64 = 1
	for _, expr := range e.Expr {
		complexity += expr.Complexity()
	}
	return complexity
}

func (e *AndExpr) isFilterExpr() {}
