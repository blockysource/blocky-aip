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

// HandleFactorExpr handles an ast.FactorExpr and returns resulting expression.
// If there is only one term, the term is handled directly.
// If there are multiple terms, they are handled as an OR expression.
func (b *Interpreter) HandleFactorExpr(ctx *ParseContext, factor *ast.FactorExpr) (TryParseValueResult, error) {
	// If there is only one term, we can handle directly its term expression.
	if len(factor.Terms) == 1 {
		return b.HandleTermExpr(ctx, factor.Terms[0])
	}

	// If there are multiple terms, we handle them as an OR expression.
	var isIndirect bool
	or := expr.AcquireOrExpr()
	for _, term := range factor.Terms {
		te, err := b.HandleTermExpr(ctx, term)
		if err != nil {
			or.Free()
			return te, err
		}
		or.Expr = append(or.Expr, te.Expr)
		isIndirect = isIndirect || te.IsIndirect
	}
	return TryParseValueResult{Expr: or, IsIndirect: isIndirect}, nil
}
