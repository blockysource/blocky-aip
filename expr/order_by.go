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

var (
	orderExprPool = &sync.Pool{
		New: func() any {
			return &OrderByExpr{
				Fields:     make([]*OrderByFieldExpr, 0, 10),
				isAcquired: true,
			}
		},
	}
	orderFieldExprPool = &sync.Pool{
		New: func() any {
			return &OrderByFieldExpr{
				isAcquired: true,
			}
		},
	}
)

// AcquireOrderByExpr acquires an OrderByExpr from the pool.
func AcquireOrderByExpr() *OrderByExpr {
	return orderExprPool.Get().(*OrderByExpr)
}

// OrderByExpr is an expression that selects a field to be used for ordering
type OrderByExpr struct {
	// Fields is a list of fields to be used for ordering
	Fields []*OrderByFieldExpr

	isAcquired bool
}

// Free puts the OrderByExpr back to the pool.
func (o *OrderByExpr) Free() {
	if o == nil {
		return
	}
	for _, sub := range o.Fields {
		if sub != nil {
			sub.Free()
		}
	}
	if !o.isAcquired {
		return
	}
	o.Fields = o.Fields[:0]
	orderExprPool.Put(o)
}

// Complexity returns the complexity of the expression
func (o *OrderByExpr) Complexity() int64 {
	complexity := int64(1)
	for _, expr := range o.Fields {
		complexity += expr.Field.Complexity()
	}
	return complexity
}

func (o *OrderByExpr) isOrderExpr() {}

// AcquireOrderByFieldExpr acquires an OrderByFieldExpr from the pool.
func AcquireOrderByFieldExpr() *OrderByFieldExpr {
	return orderFieldExprPool.Get().(*OrderByFieldExpr)
}

// OrderByFieldExpr is an expression that selects a field to be used for ordering
type OrderByFieldExpr struct {
	// Field is the field to be used for ordering
	// It can be a traversal selector.
	Field *FieldSelectorExpr

	// Order is the order of the order by expression
	Order Order

	isAcquired bool
}

// Free puts the OrderByFieldExpr back to the pool.
func (o *OrderByFieldExpr) Free() {
	if o == nil {
		return
	}
	o.Field.Free()
	if !o.isAcquired {
		return
	}
	orderFieldExprPool.Put(o)
}

// Complexity returns the complexity of the expression
func (o *OrderByFieldExpr) Complexity() int64 {
	return o.Field.Complexity()
}

func (o *OrderByFieldExpr) isOrderExpr() {}

// Order is an enum for the order of the order by expression
type Order int

const (
	// ASC determines the order to be ascending.
	// This is the default order.
	ASC Order = iota
	// DESC determines the order to be descending.
	DESC
)

func (o Order) String() string {
	switch o {
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	default:
		return "UNKNOWN"
	}
}
