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

package parser_test

import (
	"fmt"
	"os"

	"github.com/blockysource/blocky-aip/filtering/ast"
	"github.com/blockysource/blocky-aip/filtering/parser"
)

func ExampleParse() {
	p := parser.NewParser("m = 10")

	parsed, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer parsed.Free()

	r := parsed.Expr.Sequences[0].Factors[0].Terms[0].Expr.(*ast.RestrictionExpr)

	left := r.Comparable.(*ast.MemberExpr).Value.(*ast.TextLiteral)

	right := r.Arg.(*ast.MemberExpr).Value.(*ast.TextLiteral)

	fmt.Printf("left: %s, right: %s", left.Value, right.Value)

	// Output:
	// left: m, right: 10
}
