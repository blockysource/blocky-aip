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

	"google.golang.org/protobuf/reflect/protoreflect"
)

var mapValueExprPool = &sync.Pool{
	New: func() any {
		return &MapValueExpr{
			Values:     make([]MapValueExprEntry, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireMapValueExpr acquires a MapValueExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireMapValueExpr() *MapValueExpr {
	return mapValueExprPool.Get().(*MapValueExpr)
}

// Free puts the MapValueExpr back to the pool.
func (e *MapValueExpr) Free() {
	if e == nil {
		return
	}
	for _, entry := range e.Values {
		if entry.Key != nil {
			entry.Key.Free()
		}
		if entry.Value != nil {
			entry.Value.Free()
		}
	}
	if !e.isAcquired {
		return
	}
	e.MapKey = nil
	e.MapValue = nil
	e.Values = e.Values[:0]
	mapValueExprPool.Put(e)
}

var _ FilterExpr = (*MapValueExpr)(nil)

type (
	// MapValueExpr is an expression that can be represented as a map of values.
	MapValueExpr struct {
		MapKey     protoreflect.FieldDescriptor
		MapValue   protoreflect.FieldDescriptor
		Values     []MapValueExprEntry
		isAcquired bool
	}
	// MapValueExprEntry is an entry of the MapValueExpr.
	MapValueExprEntry struct {
		Key   *ValueExpr
		Value FilterExpr
	}
)

// Complexity of the MapValueExpr is the sum of complexities of the inner expressions + 1.
func (e *MapValueExpr) Complexity() int64 {
	c := int64(1)
	for _, entry := range e.Values {
		if entry.Key != nil {
			c += entry.Key.Complexity()
		}
		if entry.Value != nil {
			c += entry.Value.Complexity()
		}
	}
	return c
}

func (e *MapValueExpr) isFilterExpr() {}
