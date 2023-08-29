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

var messageSelectExprPool = &sync.Pool{
	New: func() any {
		return &MessageSelectExpr{
			Fields:     make([]*FieldSelectorExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireMessageSelectExpr acquires a select expression from the pool.
func AcquireMessageSelectExpr() *MessageSelectExpr {
	return messageSelectExprPool.Get().(*MessageSelectExpr)
}

// MessageSelectExpr is a select expression.
// It provides a way to select specific fields from a message.
// It is used by the fieldmask parser.
type MessageSelectExpr struct {
	// Message is the full name of the message paths.
	Message protoreflect.FullName

	// Fields is a list of field selector expressions.
	Fields []*FieldSelectorExpr

	isAcquired bool
}

// Free frees the select expression.
func (e *MessageSelectExpr) Free() {
	for _, path := range e.Fields {
		path.Free()
	}
	if e.isAcquired {
		e.Fields = e.Fields[:0]
		messageSelectExprPool.Put(e)
	}
}

// Equals returns true if the select expressions are equal.
func (e *MessageSelectExpr) Equals(other Expr) bool {
	if e == nil || other == nil {
		return false
	}
	oe, ok := other.(*MessageSelectExpr)
	if !ok {
		return false
	}

	if len(e.Fields) != len(oe.Fields) {
		return false
	}

	for i := range e.Fields {
		if !e.Fields[i].Equals(oe.Fields[i]) {
			return false
		}
	}

	return true
}

// Clone returns a deep copy of the select expression.
func (e *MessageSelectExpr) Clone() Expr {
	clone := AcquireMessageSelectExpr()
	for _, path := range e.Fields {
		clone.Fields = append(clone.Fields, path.Clone().(*FieldSelectorExpr))
	}
	return clone
}
