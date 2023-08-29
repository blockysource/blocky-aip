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

package expr_test

import (
	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/testpb"
)

func ExampleComposer_And() {
	md := new(testpb.Message).ProtoReflect().Descriptor()

	c := expr.Composer{Desc: md}

	// (i32 > 1) && (i64 < 2)
	x := c.And(
		c.Composite(
			c.Compare(
				c.MustSelect("i32"),
				expr.GT,
				c.Value(1),
			),
		),
		c.Composite(
			c.Compare(
				c.MustSelect("i64"),
				expr.LT,
				c.Value(2),
			),
		),
	)
	defer x.Free()

	// Output:
}
