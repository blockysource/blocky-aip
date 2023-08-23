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
	"strings"

	"github.com/blockysource/blocky-aip/filtering/token"
)

// ArgListExpr is a list of arguments to a function call.
//
// EBNF:
//
// argList
//
//	: arg { COMMA arg}
//	;
type ArgListExpr struct {
	Args []ArgExpr
}

// Position returns the position of the first argument.
func (a *ArgListExpr) Position() token.Position {
	if len(a.Args) > 0 {
		return a.Args[0].Position()
	}
	return 0
}

func (a *ArgListExpr) String() string {
	var sb strings.Builder
	for i, arg := range a.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.String())
	}
	return sb.String()
}

func (a *ArgListExpr) UnquotedString() string {
	var sb strings.Builder
	for i, arg := range a.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.UnquotedString())
	}
	return sb.String()
}

func (a *ArgListExpr) WriteStringTo(sb *strings.Builder, unquoted bool) {
	for i, arg := range a.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		arg.WriteStringTo(sb, unquoted)
	}
}

func (a *ArgListExpr) isArgExpr() {}

func (*ArgListExpr) isAstExpr() {}
