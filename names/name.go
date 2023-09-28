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

package names

// Name is the generic resource name type.
// It is used to extract the resource name parts from the resource name.
// I.e.:
// 'projects/{project}/keys/{key} ->
//
//		Name.Part(0) = projects
//		Name.Part(1) = {project}
//		Name.Part(2) = keys
//		Name.Part(3) = {key}
//	 	Name.Part(4) = ""
//		Name.Part(-1) = {key}
//		Name.Part(-2) = keys
//		Name.Part(-3) = {project}
//		Name.Part(-4) = projects
//		Name.Part(-5) = ""
type Name string

// Part returns the i-th part of the resource name.
// The index i is zero-based, so the first part of the name has index 0.
// If i is negative, the part is counted from the end of the name (i.e. -1 is the last part).
// The function call is safe for out-of-range indices, thus, it can be used to iterate over the name parts.
// If the index is out of range, the function returns an empty string.
func (n Name) Part(i int) string {
	if len(n) == 0 {
		return ""
	}

	if i < 0 {
		neg := -i
		var partStart, partEnd int
		// Iterate from the end of the string.
		partStart = len(n)
		for j := 1; j <= neg; j++ {
			partEnd = partStart
			for ; partStart > 0; partStart-- {
				if n[partStart-1] == '/' {
					if j != neg {
						partStart--
					}
					break
				}
			}
		}
		return string(n[partStart:partEnd])
	}

	var partStart, partEnd int
	for j := 0; j <= i; j++ {
		partStart = partEnd
		for ; partEnd < len(n); partEnd++ {
			if n[partEnd] == '/' {
				if j != i {
					partEnd++
				}
				break
			}
		}
	}
	return string(n[partStart:partEnd])
}

// Parts returns the number of parts in the resource name.
func (n Name) Parts() int {
	if len(n) == 0 {
		return 0
	}

	// The first part is always present.
	parts := 1
	for _, r := range n {
		if r == '/' {
			parts++
		}
	}
	return parts
}
