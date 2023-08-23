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

// Package parser provides AIP-160 compliant parser for filtering expressions.
// The parser is not thread-safe.
//
// The implementation satisfies all the requirements of AIP-160,
// and parses an input string filter expression into an AST.
//
// A parser is created by calling NewParser, and can be optimized for a specific
// needs by providing a set of options.
//
// Example:
//
//	func main() {
//	 p := parser.NewParser(parser.ErrorHandler(func(pos token.Position, msg string) {
//	     log.Printf("Error at %s: %s", pos, msg)
//	 }))
//
//	 expr, err := p.Parse("foo:bar")
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//
// The parser by default doesn't recognize any identifiers, and values.
// The literals are either a *ast.TextLiteral or  *ast.StringLiteral.
// What's more as defined in the ebnf grammar, if a TEXT literal contains a
// dot (.) character, it is separated into two TEXT literals.
// This behavior can be difficult to handle literals with a period (.) in their definitions.
// To solve this problem, the parser provides MemberHandler functions which can be used to
// merge and (if needed) split literals.
// An example of such a function is the ParseMemberNumber function, which tries to decode a *ast.MemberExpr
// elements, and parses their value either as a float64 or int64. What's more in case of a float64,
// it merges the two literals into one, and sets the value of the first literal to the string representation
// of the float64 value. In addition, it keeps the decoded value in a DecodedValue field of the *ast.MemberExpr,
// which could be reused later.
//
// Example:
//
//	func main() {
//	 p := parser.NewParser(parser.ArgMemberModifierOption(parser.ParseMemberNumber))
//
//	 expr, err := p.Parse("m = 1.0")
//	 if err != nil {
//	     log.Fatal(err)
//	 }
//	 // expr is a *parser.ParsedFilter, which contains a *ast.RestrictionExpr.
//	 // it has a Comparable field set to a 'm' *ast.MemberExpr, with a Value of *ast.TextLiteral.Value = "a",
//	 // in addition it has an Arg ast.ArgExpr which is a *ast.MemberExpr with a single Value of
//	 // *ast.TextLiteral.Value = "1.0", and a DecodedValue == float64(1.0).
//	 // If no ArgMemberModifierOption was provided, the Arg would be a *ast.MemberExpr with a
//	 // Value = *ast.TextLiteral.Value = "1" and a Field = *ast.TextLiteral.Value = "0".
package parser
