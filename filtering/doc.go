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

// Package filtering provides a parser for the AIP-160 filtering language,
// that matches protocol buffer messages based on the protoreflection and blocky api annotations.
// Read more at https://google.aip.dev/160.
//
// An Interpreter is used to parse a filter expression and return an implementation of the expr.Expr interface.
// This implementation supports custom function calls, that can be either direct or abstract call.
// An abstract call results in returning of expr.FunctionCall expression, that can be handled by the caller.
// A direct call works like a macro that converts input arguments to a concrete expression.
// This implementation supports extensions to the standard AIP-160 filtering language:
//   - IN operator - checks if a value is in a list of values etc.
//   - Struct value expression - allows to compose a proto.Message or a map field from a syntax like:
//     pkg.MyType{field1: value1, field2: value2} or
//     map{key1: value1, key2: value2}
//   - Array value expression - allows to use repeated value expression like: [1, 2, 3]
package filtering
