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

package ast

import (
	"fmt"
	"strings"

	"github.com/blockysource/blocky-aip/token"
)

var (
	// Compile-time check that *CompositeExpr implements ArgExpr.
	_ ArgExpr = (*CompositeExpr)(nil)

	// Compile-time check that *CompositeExpr implements SimpleExpr.
	_ SimpleExpr = (*CompositeExpr)(nil)
)

// CompositeExpr is a composite expression with a left parenthesis, expression, and right parenthesis.
type CompositeExpr struct {
	Lparen token.Position
	Expr   *Expr
	Rparen token.Position
}

func (c *CompositeExpr) UnquotedString() string {
	if c.Expr == nil {
		return ""
	}
	return fmt.Sprintf("(%s)", c.Expr.UnquotedString())
}

func (c *CompositeExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	sb.WriteRune('(')
	if c.Expr != nil {
		c.Expr.WriteStringTo(sb, unquoted)
	}
	sb.WriteRune(')')
}

func (c *CompositeExpr) Position() token.Position { return c.Lparen }
func (c *CompositeExpr) String() string {
	return fmt.Sprintf("(%s)", c.Expr)
}
func (*CompositeExpr) isArgExpr()    {}
func (*CompositeExpr) isSimpleExpr() {}
func (*CompositeExpr) isAstExpr()    {}
