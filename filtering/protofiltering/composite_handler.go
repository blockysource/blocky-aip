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

package protofiltering

import (
	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// HandleCompositeExpr handles an ast.CompositeExpr and returns an expression.
func (b *Interpreter) HandleCompositeExpr(ctx *ParseContext, x *ast.CompositeExpr) (TryParseValueResult, error) {
	// In a standard way we handle the composite expression, and surround it with a group expression.
	res, err := b.HandleExpr(ctx, x.Expr)
	if err != nil {
		return res, err
	}

	// Acquire a composite expression and set the handled expression as its expression.
	cps := expr.AcquireCompositeExpr()
	cps.Expr = res.Expr

	return TryParseValueResult{Expr: cps, IsIndirect: res.IsIndirect}, nil
}
