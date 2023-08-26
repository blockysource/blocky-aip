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
	"fmt"
	"sync"
)

func init() {
	gob.Register(new(FunctionCallExpr))
}

var functionCallExprPool = &sync.Pool{
	New: func() any {
		return &FunctionCallExpr{
			Arguments:  make([]FilterExpr, 0, 10),
			isAcquired: true,
		}
	},
}

// AcquireFunctionCallExpr acquires a FunctionCallExpr from the pool.
// Once acquired it must be released via Free method.
func AcquireFunctionCallExpr() *FunctionCallExpr {
	return functionCallExprPool.Get().(*FunctionCallExpr)
}

// Compile-time check to verify that FunctionCallExpr implements Expr and FilterExpr interface.
var (
	_ FilterExpr = (*FunctionCallExpr)(nil)
	_ Expr       = (*FunctionCallExpr)(nil)
)

// FunctionCallExpr is an expression that represents a function call.
// It should be used by the service that handles the function call.
// It may be used by the Database, filtering service, etc.
type FunctionCallExpr struct {
	// PkgName is the name of the package where the function is defined.
	PkgName string

	// Name is the name of the function call.
	Name string

	// Arguments is a list of arguments of the function call.
	// If empty then the function call has no arguments.
	Arguments []FilterExpr

	// CallComplexity is the complexity of the function call,
	// predefined by the parser or the function call handler.
	CallComplexity int64

	// isAcquired is true if the Expression was isAcquired from the pool.
	isAcquired bool
}

// Clone returns a copy of the current expression.
func (x *FunctionCallExpr) Clone() Expr {
	if x == nil {
		return nil
	}
	clone := AcquireFunctionCallExpr()
	clone.PkgName = x.PkgName
	clone.Name = x.Name
	clone.CallComplexity = x.CallComplexity
	for _, a := range x.Arguments {
		clone.Arguments = append(clone.Arguments, a.Clone().(FilterExpr))
	}
	return clone
}

// Equals returns true if the given expression is equal to the current one.
func (x *FunctionCallExpr) Equals(other Expr) bool {
	if other == nil {
		return false
	}
	oc, ok := other.(*FunctionCallExpr)
	if !ok {
		return false
	}
	if x.PkgName != oc.PkgName || x.Name != oc.Name || len(x.Arguments) != len(oc.Arguments) {
		return false
	}
	for i := range x.Arguments {
		if !x.Arguments[i].Equals(oc.Arguments[i]) {
			return false
		}
	}
	return true
}

// Complexity returns the complexity of the expression.
func (x *FunctionCallExpr) Complexity() int64 {
	c := x.CallComplexity
	for _, a := range x.Arguments {
		// The complexity of the arguments is taken as complexity of parsing them.
		c += a.Complexity()
	}
	return c + 1
}

// Free puts the FunctionCallExpr back to the pool.
func (x *FunctionCallExpr) Free() {
	if x == nil {
		return
	}
	for _, a := range x.Arguments {
		a.Free()
	}
	if x.isAcquired {
		x.PkgName = ""
		x.Name = ""
		x.CallComplexity = 0
		x.Arguments = x.Arguments[:0]
		functionCallExprPool.Put(x)
	}
}

func (x *FunctionCallExpr) FullName() string {
	if x.PkgName == "" {
		return x.Name
	}
	return fmt.Sprintf("%s.%s", x.PkgName, x.Name)
}

func (x *FunctionCallExpr) isFilterExpr() {}
