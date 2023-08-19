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

var arrayExprPool = &sync.Pool{
	New: func() any {
		return &ArrayExpr{
			Elements:   make([]FilterExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireArrayExpr acquires an ArrayExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireArrayExpr() *ArrayExpr {
	return arrayExprPool.Get().(*ArrayExpr)
}

// Free puts the ArrayExpr back to the pool.
func (e *ArrayExpr) Free() {
	if e == nil {
		return
	}
	if e.isAcquired {
		e.Elements = e.Elements[:0]
		arrayExprPool.Put(e)
	}
}

var _ FilterExpr = (*ArrayExpr)(nil)

// ArrayExpr is an expression that can be
// represented as an array of expressions.
type ArrayExpr struct {
	// Elements is a list of expression values.
	Elements []FilterExpr

	isAcquired bool
}

// Complexity of the ArrayExpr is the number of values + 1.
func (e *ArrayExpr) Complexity() int64 {
	return 1 + int64(len(e.Elements))
}

func (e *ArrayExpr) isFilterExpr() {}
