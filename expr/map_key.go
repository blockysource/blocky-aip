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

var mapKeyExprPool = &sync.Pool{
	New: func() any {
		return &MapKeyExpr{
			isAcquired: true,
		}
	},
}

// AcquireMapKeyExpr acquires a MapKeyExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireMapKeyExpr() *MapKeyExpr {
	return mapKeyExprPool.Get().(*MapKeyExpr)
}

// Free puts the MapKeyExpr back to the pool.
func (e *MapKeyExpr) Free() {
	if e == nil {
		return
	}
	if e.Key != nil {
		e.Key.Free()
		e.Key = nil
	}
	if e.Traversal != nil {
		e.Traversal.Free()
		e.Traversal = nil
	}
	if !e.isAcquired {
		return
	}
	*e = MapKeyExpr{}
	mapKeyExprPool.Put(e)
}

var _ FilterExpr = (*MapKeyExpr)(nil)

// MapKeyExpr is an expression that represents a map field - key.
// This expression might be used for filtering map key presence or
// for filtering map key value.
type MapKeyExpr struct {
	// Key is the key expression of the map field.
	Key FilterExpr
	// Traversal is the traversal expression of the map field.
	Traversal FilterExpr

	isAcquired bool
}

// Complexity returns the complexity of the expression.
func (e *MapKeyExpr) Complexity() int64 {
	c := int64(1)
	if e.Key != nil {
		c += e.Key.Complexity()
	}

	if e.Traversal != nil {
		c += e.Traversal.Complexity()
	}
	return c
}

func (e *MapKeyExpr) isFilterExpr() {}