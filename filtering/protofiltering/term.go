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

// HandleTermExpr handles an ast.TermExpr and returns resulting expression.
// If Unary operator is defined
func (b *Interpreter) HandleTermExpr(ctx *ParseContext, term *ast.TermExpr) (TryParseValueResult, error) {
	// If no unary operator is set, we can handle directly its simple expression.
	if !term.HasNegation() {
		return b.HandleSimpleExpr(ctx, term.Expr)
	}

	// If the negation operator is set, we handle the simple expression and surround it with a not expression.
	res, err := b.HandleSimpleExpr(ctx, term.Expr)
	if err != nil {
		return res, err
	}

	ne := expr.AcquireNotExpr()
	ne.Expr = res.Expr

	return TryParseValueResult{Expr: ne, IsIndirect: res.IsIndirect}, nil
}
