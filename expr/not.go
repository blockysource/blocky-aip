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
	"encoding/gob"
	"sync"
)

func init() {
	gob.Register(new(NotExpr))
}

var notExprPool = &sync.Pool{
	New: func() any {
		return &NotExpr{
			isAcquired: true,
		}
	},
}

// AcquireNotExpr acquires a NotExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireNotExpr() *NotExpr {
	return notExprPool.Get().(*NotExpr)
}

// Compile-time check to verify that NotExpr implements Expr and FilterExpr interface.
var (
	_ FilterExpr = (*NotExpr)(nil)
	_ Expr       = (*NotExpr)(nil)
)

// NotExpr is an expression that returns a negated result.
type NotExpr struct {
	// Expr is an expression that should be negated.
	Expr FilterExpr

	isAcquired bool
}

// Clone returns a copy of the current expression.
func (e *NotExpr) Clone() Expr {
	if e == nil {
		return nil
	}
	ne := AcquireNotExpr()
	ne.Expr = e.Expr.Clone().(FilterExpr)
	return ne
}

// Equals returns true if the given expression is equal to the current one.
func (e *NotExpr) Equals(other Expr) bool {
	if e == nil || other == nil {
		return false
	}
	if oc, ok := other.(*NotExpr); ok {
		return e.Expr.Equals(oc.Expr)
	}
	return false
}

// Free puts the NotExpr back to the pool.
func (e *NotExpr) Free() {
	if e == nil {
		return
	}
	if e.Expr != nil {
		e.Expr.Free()
		e.Expr = nil
	}
	if e.isAcquired {
		notExprPool.Put(e)
	}
}

// Complexity of the NotExpr is 1 + complexity of the inner expression.
func (e *NotExpr) Complexity() int64 {
	return 1 + e.Expr.Complexity()
}

func (e *NotExpr) isFilterExpr() {}
