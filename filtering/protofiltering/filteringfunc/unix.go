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

package filteringfunc

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/filtering/protofiltering"
)

var timestampDesc = new(timestamppb.Timestamp).ProtoReflect().Descriptor()

// TimeUnix is a protofiltering function call declaration,
// that converts an input int64 value into a valid google.protobuf.Timestamp - (time.Time ValueExpression).
// It may either take a direct or indirect value of the input argument.
func TimeUnix() *protofiltering.FunctionCallDeclaration {
	return &unixFunc
}

var unixFunc = protofiltering.FunctionCallDeclaration{
	Name: protofiltering.FunctionName{
		PkgName: "time",
		Name:    "Unix",
	},
	Arguments: []*protofiltering.FunctionCallArgumentDeclaration{
		{
			Indirect:   true,
			ArgName:    "unix_ts",
			IsRepeated: false,
			IsNullable: false,
			FieldKind:  protoreflect.Int64Kind,
		},
	},
	Returning: &protofiltering.FunctionCallReturningDeclaration{
		ServiceCalled:     false,
		FieldKind:         protoreflect.MessageKind,
		MessageDescriptor: timestampDesc,
		IsNullable:        false,
		IsRepeated:        false,
	},
	CallFn: func(args ...expr.FilterExpr) (protofiltering.FunctionCallArgument, error) {
		if len(args) != 1 {
			// This is internal error.
			return protofiltering.FunctionCallArgument{}, fmt.Errorf("invalid number of arguments for unix function: %v", len(args))
		}

		switch ve := args[0].(type) {
		case *expr.ValueExpr:
			var i64 int64
			switch vt := ve.Value.(type) {
			case int64:
				i64 = vt
			case uint64:
				i64 = int64(vt)
			default:
				return protofiltering.FunctionCallArgument{}, fmt.Errorf("input value is not a valid int64 value expression: %T", ve.Value)
			}

			tm := time.Unix(i64, 0)

			res := expr.AcquireValueExpr()
			res.Value = tm

			return protofiltering.FunctionCallArgument{
				Expr:       res,
				IsIndirect: false,
			}, nil
		case *expr.FieldSelectorExpr, *expr.MapKeyExpr:
			fc := expr.AcquireFunctionCallExpr()
			fc.PkgName = "time"
			fc.Name = "Unix"
			fc.Arguments = append(fc.Arguments, ve)
			fc.CallComplexity = 1
			return protofiltering.FunctionCallArgument{
				Expr:       fc,
				IsIndirect: true,
			}, nil
		default:
			return protofiltering.FunctionCallArgument{}, fmt.Errorf("input value is not a valid int64 value expression: %T", args[0])
		}
	},
}
