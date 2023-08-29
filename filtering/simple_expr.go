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
	"fmt"

	"github.com/blockysource/blocky-aip/filtering/ast"
)

// HandleSimpleExpr handles an ast.SimpleExpr and returns an handled expression.
// The handled expression can be a composite expression or a restriction expression.
func (b *Interpreter) HandleSimpleExpr(ctx *ParseContext, x ast.SimpleExpr) (TryParseValueResult, error) {
	// No standard simple expression handler found, we need to check the type of the simple expression
	switch e := x.(type) {
	case *ast.CompositeExpr:
		return b.HandleCompositeExpr(ctx, e)
	case *ast.RestrictionExpr:
		return b.HandleRestrictionExpr(ctx, e)
	default:
		var res TryParseValueResult
		if ctx.ErrHandler != nil {
			res.ErrPos = x.Position()
			res.ErrMsg = fmt.Sprintf("unknown simple expression type %T", x)
		}
		return res, ErrInternal
	}

}
