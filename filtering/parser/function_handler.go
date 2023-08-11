package parser

import (
	"github.com/blockysource/blocky-aip/filtering/ast"
)

// FunctionCallHandler is a function that modifies, parses or validates a function call.
// To match the name of the function call all name parts, use it's JoinedNameEquals("func.call.name") method.
// This might be used in example to cast custom types like:
// uuid("123e4567-e89b-12d3-a456-426614174000") to some UUID implementation.
// The *ast.FunctionCall contain a DecodedValue field which should be used to store the decoded value.
type FunctionCallHandler func(r *ast.FunctionCall) bool
