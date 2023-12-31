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

// Expr is a generic expression interface that can be used to represent
// any expression.
type Expr interface {
	// Free releases expression resources.
	// No further calls to the expression are allowed after calling Free.
	// This should release the resource related to the expression back to the pool.
	Free()

	// Equals returns true if the expression is equal to the other expression.
	Equals(other Expr) bool

	// Clone returns a deep copy of the expression.
	Clone() Expr
}
