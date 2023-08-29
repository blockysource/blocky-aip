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

var updateExprPool = &sync.Pool{
	New: func() any {
		return &UpdateExpr{
			Elements:   make([]UpdateFieldValue, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireUpdateExpr acquires an UpdateExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireUpdateExpr() *UpdateExpr {
	return updateExprPool.Get().(*UpdateExpr)
}

// UpdateExpr is an expression that contains fields to update along with their values.
// The UpdateExpr can be a value of UpdateFieldValue.
// In that case, the elements of the Value UpdateExpr, are relative to the field of the parent UpdateExpr.
type UpdateExpr struct {
	// Elements is a list of fields to update along with their values.
	Elements []UpdateFieldValue

	isAcquired bool
}

// Free puts the UpdateExpr back to the pool.
func (e *UpdateExpr) Free() {
	if e == nil {
		return
	}
	for _, sub := range e.Elements {
		if sub.Field != nil {
			sub.Field.Free()
		}
		if sub.Value != nil {
			sub.Value.Free()
		}
	}
	if e.isAcquired {
		e.Elements = e.Elements[:0]
		updateExprPool.Put(e)
	}
}

// Clone returns a copy of the UpdateExpr.
func (e *UpdateExpr) Clone() Expr {
	if e == nil {
		return nil
	}
	clone := AcquireUpdateExpr()
	for _, expr := range e.Elements {
		clone.Elements = append(clone.Elements, UpdateFieldValue{
			Field: expr.Field.Clone().(*FieldSelectorExpr),
			Value: expr.Value.Clone().(UpdateValueExpr),
		})
	}
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (e *UpdateExpr) Equals(other Expr) bool {
	if other == nil {
		return false
	}

	oa, ok := other.(*UpdateExpr)
	if !ok {
		return false
	}

	if len(e.Elements) != len(oa.Elements) {
		return false
	}

	for i := range e.Elements {
		if !e.Elements[i].Field.Equals(oa.Elements[i].Field) {
			return false
		}
		if !e.Elements[i].Value.Equals(oa.Elements[i].Value) {
			return false
		}
	}

	return true
}

// UpdateExpression can be embedded as a value of UpdateExpr.
func (e *UpdateExpr) isUpdateValueExpr() {}

// UpdateFieldValue is a field to update along with its value.
type UpdateFieldValue struct {
	// Field is a field name to update.
	Field *FieldSelectorExpr

	// Value is a value to set.
	Value UpdateValueExpr
}

// UpdateValueExpr is an expression that can be used as a value in UpdateExpr.
type UpdateValueExpr interface {
	Expr
	isUpdateValueExpr()
}

var arrayUpdateExprPool = &sync.Pool{
	New: func() any {
		return &ArrayUpdateExpr{
			Elements:   make([]*UpdateExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireArrayUpdateExpr acquires an ArrayUpdateExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireArrayUpdateExpr() *ArrayUpdateExpr {
	return arrayUpdateExprPool.Get().(*ArrayUpdateExpr)
}

// ArrayUpdateExpr is an expression that can be used as a value in UpdateExpr.
// It describes an input array of UpdateExpr.
type ArrayUpdateExpr struct {
	// Elements is a list of expression values.
	Elements []*UpdateExpr

	isAcquired bool
}

// Free puts the ArrayUpdateExpr back to the pool.
func (e *ArrayUpdateExpr) Free() {
	if e == nil {
		return
	}
	for _, sub := range e.Elements {
		sub.Free()
	}
	if e.isAcquired {
		e.Elements = e.Elements[:0]
		arrayUpdateExprPool.Put(e)
	}
}

// Clone returns a copy of the ArrayUpdateExpr.
func (e *ArrayUpdateExpr) Clone() Expr {
	if e == nil {
		return nil
	}
	clone := AcquireArrayUpdateExpr()
	for _, expr := range e.Elements {
		clone.Elements = append(clone.Elements, expr.Clone().(*UpdateExpr))
	}
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (e *ArrayUpdateExpr) Equals(other Expr) bool {
	if other == nil {
		return false
	}

	oa, ok := other.(*ArrayUpdateExpr)
	if !ok {
		return false
	}

	if len(e.Elements) != len(oa.Elements) {
		return false
	}

	for i := range e.Elements {
		if !e.Elements[i].Equals(oa.Elements[i]) {
			return false
		}
	}

	return true
}

// UpdateExpression can be embedded as a value of UpdateExpr.
func (e *ArrayUpdateExpr) isUpdateValueExpr() {}
