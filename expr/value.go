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

var valueExprPool = &sync.Pool{
	New: func() any {
		return &ValueExpr{
			isAcquired: true,
		}
	},
}

// AcquireValueExpr acquires a ValueExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireValueExpr() *ValueExpr {
	return valueExprPool.Get().(*ValueExpr)
}

// Free puts the ValueExpr back to the pool.
func (x *ValueExpr) Free() {
	if x == nil || !x.isAcquired {
		return
	}
	x.Value = nil
	valueExprPool.Put(x)
}

var _ FilterExpr = (*ValueExpr)(nil)

// ValueExpr is a simple value expression that contains a value.
// The value may be of any type that matches related to this expression.
// Standard field types used in the expressions are:
// - string
// - int64
// - uint64
// - bool
// - float64
// - []byte
// - time.Time
// - time.Duration
// - protoreflect.EnumNumber -- enum value
// - protoreflect.Message - message value (dynamicpb.Message for dynamic structs)
// - structpb.Value
// - nil - used for nullable fields
// This can be extended by custom types.
type ValueExpr struct {
	// Value is the value of the expression.
	Value any

	isAcquired bool
}

// Complexity of the ValueExpr is 1.
func (x *ValueExpr) Complexity() int64 {
	return 1
}

func (*ValueExpr) isFilterExpr() {}
