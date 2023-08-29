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

package filtering_test

import (
	"fmt"
	"os"

	"github.com/blockysource/blocky-aip/filtering"
	"github.com/blockysource/blocky-aip/internal/testpb"
	"github.com/blockysource/blocky-aip/token"
)

func ExampleInterpreter_Parse() {
	var i filtering.Interpreter

	msg := new(testpb.Message)

	err := i.Reset(msg.ProtoReflect().Descriptor(), filtering.ErrHandlerOpt(func(pos token.Position, msg string) {
		fmt.Printf("Error: %v at: %d\n", msg, pos)
	}))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Parse a filter.
	filter := `name = "value" AND (age > 18 OR age < 10)`
	pf, err := i.Parse(filter)
	if err != nil {
		// Handle the error on request.
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	// Always Free the memory of allocated filter expressions.
	// It is safe to call Free multiple times, or on nil pointer expressions.
	defer pf.Free()

	// Evaluate the filter.
	// Output:
	//
}
