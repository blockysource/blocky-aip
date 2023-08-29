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
	gob.Register(new(MapKeyExpr))
}

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

// Compile-time check to verify that MapKeyExpr implements Expr and FilterExpr interface.
var (
	_ FilterExpr = (*MapKeyExpr)(nil)
	_ Expr       = (*MapKeyExpr)(nil)
)

// MapKeyExpr is an expression that represents a map field - key.
// This expression might be used for filtering map key presence or
// for filtering map key value.
type MapKeyExpr struct {
	// Key is the key expression of the map field.
	Key Expr
	// Traversal is the traversal expression of the map field.
	Traversal Expr

	isAcquired bool
}

// Clone returns a copy of the MapKeyExpr.
func (e *MapKeyExpr) Clone() Expr {
	if e == nil {
		return nil
	}

	mk := AcquireMapKeyExpr()
	if e.Key != nil {
		mk.Key = e.Key.Clone().(FilterExpr)
	}

	if e.Traversal != nil {
		mk.Traversal = e.Traversal.Clone().(FilterExpr)
	}
	return mk
}

// Equals returns true if the given expression is equal to the current one.
func (e *MapKeyExpr) Equals(other Expr) bool {
	if e == nil || other == nil {
		return false
	}
	om, ok := other.(*MapKeyExpr)
	if !ok {
		return false
	}

	if e.Key != nil && !e.Key.Equals(om.Key) {
		return false
	}

	if e.Traversal != nil && !e.Traversal.Equals(om.Traversal) {
		return false
	}

	return true
}

// Complexity returns the complexity of the expression.
func (e *MapKeyExpr) Complexity() int64 {
	c := int64(1)
	if e.Key != nil {
		fe, ok := e.Key.(FilterExpr)
		if ok {
			c += fe.Complexity()
		}
	}

	if e.Traversal != nil {
		fe, ok := e.Traversal.(FilterExpr)
		if ok {
			c += fe.Complexity()
		}
	}
	return c
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

func (e *MapKeyExpr) isFilterExpr() {}
