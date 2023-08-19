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

var fieldSelectorExpr = &sync.Pool{
	New: func() any {
		return &FieldSelectorExpr{
			isAcquired: true,
		}
	},
}

// AcquireFieldSelectorExpr acquires a FieldSelectorExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireFieldSelectorExpr() *FieldSelectorExpr {
	return fieldSelectorExpr.Get().(*FieldSelectorExpr)
}

// Free puts the FieldSelectorExpr back to the pool.
func (e *FieldSelectorExpr) Free() {
	if e.Traversal != nil {
		e.Traversal.Free()
		e.Traversal = nil
	}
	if e.isAcquired {
		e.Message = nil
		e.Field = nil
		e.FieldComplexity = 0
		fieldSelectorExpr.Put(e)
	}
}

var _ FilterExpr = (*FieldSelectorExpr)(nil)

// FieldSelectorExpr is a literal that represents a message field or a path of fields.
// It describes the expression "a.b.c" where b is a field of a, and c is a field of b.
type FieldSelectorExpr struct {
	// Message is the message descriptor of the literal.
	Message protoreflect.MessageDescriptor

	// Field is the field name of the literal.
	Field protoreflect.FieldDescriptor

	// Traversal is the expression related to this field literal.
	// This field is used as a linked list to traverse the field literals.
	// The whole path can be reconstructed by traversing the linked list.
	// It may be another FieldSelectorExpr or MapKeyExpr.
	Traversal FilterExpr

	// FieldComplexity is the complexity of the field, assigned by the parser.
	FieldComplexity int64

	// isAcquired is true if the field is acquired from the pool.
	isAcquired bool
}

// Parent returns the parent of the field literal.
func (e *FieldSelectorExpr) Parent() protoreflect.MessageDescriptor {
	return e.Field.Parent().(protoreflect.MessageDescriptor)
}

// Complexity returns the complexity of the field literal.
func (e *FieldSelectorExpr) Complexity() int64 {
	return e.FieldComplexity
}

// GetTraversal returns the traversal of the field literal.
func (e *FieldSelectorExpr) GetTraversal() FilterExpr {
	return e.Traversal
}

// SetTraversal sets the traversal of the field literal.
func (e *FieldSelectorExpr) SetTraversal(traversal FilterExpr) {
	e.Traversal = traversal
}

func (e *FieldSelectorExpr) isFilterExpr() {}
