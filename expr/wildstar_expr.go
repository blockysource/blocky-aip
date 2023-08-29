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

var wildcardExprPool = &sync.Pool{
	New: func() any {
		return &WildcardExpr{
			isAcquired: true,
		}
	},
}

// AcquireWildcardExpr acquires a wildcard expression from the pool.
func AcquireWildcardExpr() *WildcardExpr {
	return wildcardExprPool.Get().(*WildcardExpr)
}

// WildcardExpr is a wildcard expression.
// It is used by the fieldmask parser.
type WildcardExpr struct {
	isAcquired bool
}

// Free frees the wildcard expression.
func (e *WildcardExpr) Free() {
	if e.isAcquired {
		wildcardExprPool.Put(e)
	}
}

// Equals returns true if the wildcard expression is equal to the other expression.
func (e *WildcardExpr) Equals(other Expr) bool {
	if other == nil {
		return false
	}
	_, ok := other.(*WildcardExpr)
	return ok
}

// Clone returns a deep copy of the wildcard expression.
func (e *WildcardExpr) Clone() Expr {
	if e == nil {
		return nil
	}
	clone := AcquireWildcardExpr()
	return clone
}
