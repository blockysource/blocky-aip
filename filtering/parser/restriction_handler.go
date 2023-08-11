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

package parser

import (
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// RestrictionHandler is a function that handles parsed *ast.RestrictionExpr.
// It can be useful for customization of the parser.
// It might be used for validating the type of the restriction
// where Comparable equals to predefined value and Arg should match some pattern.
// i.e.: with the requirements such as: Comparable = msg.email, the Arg should be a valid email.
//       msg.email = "john.doe@gmail.com" should be rejected.
// To match the name of the restriction's *ast.MemberExpr, use it's JoinedNameEquals("msg.name", false) method.
type RestrictionHandler func(r *ast.RestrictionExpr) bool
