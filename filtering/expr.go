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

// HandleExpr handles an ast.FilterExpr and returns an handled expression.
func (b *Interpreter) HandleExpr(ctx *ParseContext, x *ast.Expr) (TryParseValueResult, error) {
	// If there is only one sequence, we can handle directly its sequence expression.
	if len(x.Sequences) == 1 {
		return b.HandleSequenceExpr(ctx, x.Sequences[0])
	}

	// If there are multiple sequences, we handle them as an AND expression.
	and := expr.AcquireAndExpr()
	var isIndirect bool
	for _, seq := range x.Sequences {
		se, err := b.HandleSequenceExpr(ctx, seq)
		if err != nil {
			and.Free()
			return se, err
		}

		and.Expr = append(and.Expr, se.Expr)
		isIndirect = isIndirect || se.IsIndirect
	}
	return TryParseValueResult{Expr: and, IsIndirect: isIndirect}, nil
}
