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

package filtering

import (
	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// HandleSequenceExpr handles an ast.SequenceExpr and returns resulting expression.
// A sequence might be composed of single or multiple factors.
// If it is composed of a single factor, the factor is handled directly.
// If it is composed of multiple factors, they are handled as an AND expression.
// This is called a 'fuzzy' AND expression, because it is not a strict AND expression.
// Read more at https://google.aip.dev/160#literals for more information.
func (b *Interpreter) HandleSequenceExpr(ctx *ParseContext, seq *ast.SequenceExpr) (TryParseValueResult, error) {
	if len(seq.Factors) == 1 {
		return b.HandleFactorExpr(ctx, seq.Factors[0])
	}

	// Fuzzy AND expression
	and := expr.AcquireAndExpr()
	var isIndirect bool
	for _, factor := range seq.Factors {
		fe, err := b.HandleFactorExpr(ctx, factor)
		if err != nil {
			and.Free()
			return fe, err
		}
		and.Expr = append(and.Expr, fe.Expr)
		isIndirect = isIndirect || fe.IsIndirect
	}

	return TryParseValueResult{Expr: and, IsIndirect: isIndirect}, nil
}
