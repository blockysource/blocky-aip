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

// SimpleExpr is a simple expression.
// It can either be a restriction or a composite expression.
//
// EBNF:
//
// simple
//
//	: restriction
//	| composite
//	;
type SimpleExpr interface {
	UnquotedString() string
	String() string
	WriteStringTo(sb *strings.Builder, unquoted bool)
	Position() token.Position
	isSimpleExpr()
	isAstExpr()
}
