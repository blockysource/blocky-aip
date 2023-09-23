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

import (
	"testing"
)

func TestName_Part(t *testing.T) {
	tests := []struct {
		name string
		n    Name
		i    int
		want string
	}{
		{
			name: "empty",
			n:    "",
			i:    0,
			want: "",
		},
		{
			name: "one part",
			n:    "projects",
			i:    0,
			want: "projects",
		},
		{
			name: "two parts/resource",
			n:    "projects/{project}",
			i:    0,
			want: "projects",
		},
		{
			name: "two parts/id",
			n:    "projects/{project}",
			i:    1,
			want: "{project}",
		},
		{
			name: "three parts/parent id",
			n:    "projects/{project}/keys",
			i:    1,
			want: "{project}",
		},
		{
			name: "three parts/resource",
			n:    "projects/{project}/keys",
			i:    2,
			want: "keys",
		},
		{
			name: "three parts/out of range",
			n:    "projects/{project}/keys",
			i:    5,
			want: "",
		},
		{
			name: "three parts/-1",
			n:    "projects/{project}/keys",
			i:    -1,
			want: "keys",
		},
		{
			name: "three parts/-2",
			n:    "projects/{project}/keys",
			i:    -2,
			want: "{project}",
		},
		{
			name: "three parts/-3",
			n:    "projects/{project}/keys",
			i:    -3,
			want: "projects",
		},
		{
			name: "three parts/-4",
			n:    "projects/{project}/keys",
			i:    -4,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Part(tt.i); got != tt.want {
				t.Errorf("Name.Part() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkName_Part(b *testing.B) {
	n := Name("projects/{project}/keys/{key}")
	b.Run("Part(2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n.Part(2)
		}
	})

	b.Run("Part(-2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n.Part(-2)
		}
	})
}
