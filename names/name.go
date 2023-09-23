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
//	 Name.Part(4) = ""
type Name string

// Part returns the i-th part of the resource name.
func (n Name) Part(i int) string {
	if len(n) == 0 {
		return ""
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
