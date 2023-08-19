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

// FilterExpr is a filter expression that can be evaluated.
type FilterExpr interface {
	// Complexity returns approximate complexity of the expression.
	// It is used to estimate the cost of the expression.
	// The complexity is a number of nodes in the expression tree,
	// increased by field defined complexity.
	// Each expression evaluates its complexity based on its rules.
	Complexity() int64

	// Free releases expression resources.
	// No further calls to the expression are allowed after calling Free.
	// This should release the resource related to the expression back to the pool.
	Free()

	isFilterExpr()
}

// UndefinedFilterExpr is a filter expression that is undefined.
// It is used to embed external expressions into the expression tree.
type UndefinedFilterExpr struct{}

func (UndefinedFilterExpr) Complexity() int64 { return 0 }
func (UndefinedFilterExpr) Free()             {}
func (UndefinedFilterExpr) isFilterExpr()     {}
