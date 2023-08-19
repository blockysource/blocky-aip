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
	"fmt"
	"sync"
)

var compareExprPool = &sync.Pool{
	New: func() any {
		return &CompareExpr{
			isAcquired: true,
		}
	},
}

// AcquireCompareExpr acquires a CompareExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireCompareExpr() *CompareExpr {
	return compareExprPool.Get().(*CompareExpr)
}

// Free puts the CompareExpr back to the pool.
func (x *CompareExpr) Free() {
	if x == nil {
		return
	}
	x.Comparator = 0
	if x.Left != nil {
		x.Left.Free()
	}
	if x.Right != nil {
		x.Right.Free()
	}
	if x.isAcquired {
		compareExprPool.Put(x)
	}
}

var _ FilterExpr = (*CompareExpr)(nil)

// CompareExpr is a restriction
type CompareExpr struct {
	// Left is the left hand side of the expression, the field to compare.
	Left FilterExpr

	// Comparator is the comparator to use.
	Comparator Comparator

	// Right is the right hand side of the expression, the value to compare to.
	Right FilterExpr

	isAcquired bool
}

// Complexity returns the complexity of the expression.
// The complexity is taken from the field options.
func (x *CompareExpr) Complexity() int64 {
	if x.Left == nil || x.Right == nil {
		return 1
	}
	if x.Right == nil {
		return x.Left.Complexity() + 1
	}
	return x.Left.Complexity() + x.Right.Complexity() + 1
}

func (x *CompareExpr) isFilterExpr() {}

// Comparator is a defined type for comparators.
type Comparator int

// String returns the string representation of the comparator.
func (c Comparator) String() string {
	if c < 0 || c > Comparator(len(_ComparatorStrings)-1) {
		return fmt.Sprintf("Comparator(%d)", c)
	}
	return _ComparatorStrings[c]
}

const (
	_ Comparator = iota
	// EQ is the equal to comparator.
	EQ
	// LE is the less than or equal to comparator.
	LE
	// LT is the less than comparator.
	LT
	// GE is the greater than or equal to comparator.
	GE
	// GT is the greater than comparator.
	GT
	// NE is the not equal to comparator.
	NE
	// HAS is the has comparator.
	HAS
	// IN is the in comparator.
	IN
)

var _ComparatorStrings = [...]string{
	LE:  "<=",
	LT:  "<",
	GE:  ">=",
	GT:  ">",
	NE:  "!=",
	EQ:  "=",
	HAS: ":",
	IN:  "IN",
}
