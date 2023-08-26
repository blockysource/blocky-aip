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
	gob.Register(new(StringSearchExpr))
}

var stringSearchExprPool = &sync.Pool{
	New: func() any {
		return &StringSearchExpr{
			isAcquired: true,
		}
	},
}

// AcquireStringSearchExpr acquires a StringSearchExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireStringSearchExpr() *StringSearchExpr {
	return stringSearchExprPool.Get().(*StringSearchExpr)
}



var _ FilterExpr = (*StringSearchExpr)(nil)

// StringSearchExpr is a restriction that searches for a string in a string field.
// The string can have a prefix or suffix wildcard.
type StringSearchExpr struct {
	// Value is the string value to search for (without wildcard characters (if present)).
	Value string

	// PrefixWildcard is true if the value has a prefix wildcard.
	PrefixWildcard bool

	// SuffixWildcard is true if the value has a suffix wildcard.
	SuffixWildcard bool

	// SearchComplexity is the complexity assigned by the parser.
	SearchComplexity int64

	isAcquired bool
}

// Clone returns a copy of the StringSearchExpr.
func (x *StringSearchExpr) Clone() Expr {
	if x == nil {
		return nil
	}
	clone := AcquireStringSearchExpr()
	clone.Value = x.Value
	clone.PrefixWildcard = x.PrefixWildcard
	clone.SuffixWildcard = x.SuffixWildcard
	clone.SearchComplexity = x.SearchComplexity
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (x *StringSearchExpr) Equals(other Expr) bool {
	if x == nil || other == nil {
		return false
	}
	if oc, ok := other.(*StringSearchExpr); ok {
		return x.Value == oc.Value &&
			x.PrefixWildcard == oc.PrefixWildcard &&
			x.SuffixWildcard == oc.SuffixWildcard
	}
	return false
}

// Free puts the StringSearchExpr back to the pool.
func (x *StringSearchExpr) Free() {
	if x == nil || !x.isAcquired {
		return
	}
	*x = StringSearchExpr{}
	stringSearchExprPool.Put(x)
}

// Complexity returns the complexity of the expression.
// The complexity is taken from the field options.
// If the value has a prefix or suffix wildcard, the complexity is multiplied by 2, by each of them.
// This means that the complexity is multiplied by 4 if both are present.
// Resultant complexity is increased by 1 for the node.
func (x *StringSearchExpr) Complexity() int64 {
	fc := x.SearchComplexity
	if fc == 0 {
		fc = 1
	}
	if x.PrefixWildcard {
		fc *= 2
	}
	if x.SuffixWildcard {
		fc *= 2
	}

	return fc + 1
}

func (x *StringSearchExpr) isFilterExpr() {}
