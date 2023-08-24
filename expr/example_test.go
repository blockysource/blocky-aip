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
