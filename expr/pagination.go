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
	gob.Register(new(PaginationExpr))
}

var paginationExprPool = &sync.Pool{
	New: func() any {
		return &PaginationExpr{
			isAcquired: true,
		}
	},
}

// AcquirePaginationExpr acquires a PaginationExpr from the pool.
// Once acquired it must be released via Free method.
func AcquirePaginationExpr() *PaginationExpr {
	return paginationExprPool.Get().(*PaginationExpr)
}

// Compile-time check to verify that PaginationExpr implements Expr interface.
var _ Expr = (*PaginationExpr)(nil)

// PaginationExpr is an expression that defines pagination.
type PaginationExpr struct {
	// PageSize is the number of items to return per page.
	PageSize int32

	// Skip is the number of items to skip.
	Skip int32

	isAcquired bool
}

// Free puts the PaginationExpr back to the pool.
func (x *PaginationExpr) Free() {
	if x == nil || !x.isAcquired {
		return
	}
	x.PageSize = 0
	x.Skip = 0
	paginationExprPool.Put(x)
}

// Equals returns true if the other expression is equal to the current one.
func (x *PaginationExpr) Equals(other Expr) bool {
	if x == nil || other == nil {
		return x == other
	}

	o, ok := other.(*PaginationExpr)
	if !ok {
		return false
	}

	return x.PageSize == o.PageSize && x.Skip == o.Skip
}

// Clone returns a copy of the current expression.
func (x *PaginationExpr) Clone() Expr {
	if x == nil {
		return nil
	}

	clone := AcquirePaginationExpr()
	clone.PageSize = x.PageSize
	clone.Skip = x.Skip
	return clone
}
