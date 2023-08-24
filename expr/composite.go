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
	gob.Register(new(CompositeExpr))
}

var compositeExprPool = &sync.Pool{
	New: func() any {
		return &CompositeExpr{
			isAcquired: true,
		}
	},
}

// AcquireCompositeExpr acquires a CompositeExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireCompositeExpr() *CompositeExpr {
	return compositeExprPool.Get().(*CompositeExpr)
}

var _ FilterExpr = (*CompositeExpr)(nil)

// CompositeExpr is a composite expression that wraps current expression into logical group.
type CompositeExpr struct {
	// Expr is the expression to wrap.
	Expr FilterExpr

	isAcquired bool
}

// Clone returns a copy of the current expression.
func (e *CompositeExpr) Clone() FilterExpr {
	if e == nil {
		return nil
	}

	clone := AcquireCompositeExpr()
	if e.Expr != nil {
		clone.Expr = e.Expr.Clone()
	}
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (e *CompositeExpr) Equals(other FilterExpr) bool {
	if e == nil || other == nil {
		return false
	}
	if oc, ok := other.(*CompositeExpr); ok {
		return e.Expr.Equals(oc.Expr)
	}
	return false
}

// Free puts the CompositeExpr back to the pool.
func (e *CompositeExpr) Free() {
	if e == nil {
		return
	}
	if e.Expr != nil {
		e.Expr.Free()
		e.Expr = nil
	}
	if e.isAcquired {
		compositeExprPool.Put(e)
	}
}

// Complexity of the CompositeExpr is the complexity of the inner expression + 1.
func (e *CompositeExpr) Complexity() int64 {
	return 1 + e.Expr.Complexity()
}

func (e *CompositeExpr) isFilterExpr() {}
