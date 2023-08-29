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

var mapSelectKeysExprPool = &sync.Pool{
	New: func() any {
		return &MapSelectKeysExpr{
			Keys:       make([]*MapKeyExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireMapSelectKeysExpr acquires a MapSelectKeysExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireMapSelectKeysExpr() *MapSelectKeysExpr {
	return mapSelectKeysExprPool.Get().(*MapSelectKeysExpr)
}

// Compile-time check to verify that MapSelectKeysExpr implements Expr interface.
var _ Expr = (*MapSelectKeysExpr)(nil)

// MapSelectKeysExpr is an expression that represents a map field - key.
type MapSelectKeysExpr struct {
	// Keys is the list of keys to select from the map.
	Keys []*MapKeyExpr

	isAcquired bool
}

// Free puts the MapSelectKeysExpr back to the pool.
func (e *MapSelectKeysExpr) Free() {
	if e == nil || !e.isAcquired {
		return
	}
	for _, key := range e.Keys {
		key.Free()
	}
	e.Keys = e.Keys[:0]
	mapSelectKeysExprPool.Put(e)
}

// Clone returns a copy of the MapSelectKeysExpr.
func (e *MapSelectKeysExpr) Clone() Expr {
	if e == nil {
		return nil
	}
	clone := AcquireMapSelectKeysExpr()
	for _, key := range e.Keys {
		clone.Keys = append(clone.Keys, key.Clone().(*MapKeyExpr))
	}
	return clone
}

// Equals returns true if the other expression is equal to the current one.
func (e *MapSelectKeysExpr) Equals(other Expr) bool {
	if e == nil || other == nil {
		return e == other
	}
	om, ok := other.(*MapSelectKeysExpr)
	if !ok {
		return false
	}
	if len(e.Keys) != len(om.Keys) {
		return false
	}
	for i, key := range e.Keys {
		if !key.Equals(om.Keys[i]) {
			return false
		}
	}
	return true
}
