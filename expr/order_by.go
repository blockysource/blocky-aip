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
	gob.Register(new(OrderByExpr))
	gob.Register(new(OrderByFieldExpr))
}

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

// Compile-time check to verify that OrderByExpr implements Expr interface.
var _ Expr = (*OrderByExpr)(nil)

// OrderByExpr is an expression that selects a field to be used for ordering
type OrderByExpr struct {
	// Fields is a list of fields to be used for ordering
	Fields []*OrderByFieldExpr

	isAcquired bool
}

// Merge merges the other OrderByExpr into the current one.
// If the field is already present, it is skipped and not added.
// The other OrderByExpr is not modified and is safe to use or Free after the merge.
func (o *OrderByExpr) Merge(other *OrderByExpr) {
	// Verify if no duplicated fields are present.
	ln := len(o.Fields)
	for _, field := range other.Fields {
		foundIndex := -1
		for i := 0; i < ln-1; i++ {
			f := o.Fields[i]
			if f.Field.Equals(field.Field) {
				foundIndex = i

				break
			}
		}
		if foundIndex == -1 {
			// Clone the field and add it to the list.
			o.Fields = append(o.Fields, field.Clone().(*OrderByFieldExpr))
		}
	}
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

// Equals returns true if the given expression is equal to the current one.
func (o *OrderByExpr) Equals(order Expr) bool {
	if o == nil || order == nil {
		return false
	}

	other, ok := order.(*OrderByExpr)
	if !ok {
		return false
	}

	if len(o.Fields) != len(other.Fields) {
		return false
	}

	for i := range o.Fields {
		if !o.Fields[i].Equals(other.Fields[i]) {
			return false
		}
	}
	return true
}

// Clone returns a copy of the OrderByExpr.
func (o *OrderByExpr) Clone() Expr {
	if o == nil {
		return nil
	}
	clone := AcquireOrderByExpr()
	for _, field := range o.Fields {
		clone.Fields = append(clone.Fields, field.Clone().(*OrderByFieldExpr))
	}
	return clone
}

// AcquireOrderByFieldExpr acquires an OrderByFieldExpr from the pool.
func AcquireOrderByFieldExpr() *OrderByFieldExpr {
	return orderFieldExprPool.Get().(*OrderByFieldExpr)
}

// Compile-time check to verify that OrderByFieldExpr implements Expr interface.
var _ Expr = (*OrderByFieldExpr)(nil)

// OrderByFieldExpr is an expression that selects a field to be used for ordering
type OrderByFieldExpr struct {
	// Field is the field to be used for ordering
	// It can be a traversal selector.
	Field *FieldSelectorExpr

	// Order is the order of the order by expression
	Order Order

	isAcquired bool
}

// Clone returns a copy of the OrderByFieldExpr.
func (o *OrderByFieldExpr) Clone() Expr {
	if o == nil {
		return nil
	}
	clone := AcquireOrderByFieldExpr()
	clone.Field = o.Field.Clone().(*FieldSelectorExpr)
	clone.Order = o.Order
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (o *OrderByFieldExpr) Equals(other Expr) bool {
	if o == nil || other == nil {
		return false
	}

	oe, ok := other.(*OrderByFieldExpr)
	if !ok {
		return false
	}

	return o.Field.Equals(oe.Field) && o.Order == oe.Order
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
