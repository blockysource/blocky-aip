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
	"fmt"
	"math"
	"testing"
)

func TestComposer(t *testing.T) {
	tests := []struct {
		name string
		fn   func(c *Composer)
		want string
	}{
		{
			name: "empty",
			fn:   func(c *Composer) {},
			want: "",
		},
		{
			name: "one part",
			fn: func(c *Composer) {
				c.WritePart("projects")
			},
			want: "projects",
		},
		{
			name: "two parts",
			fn: func(c *Composer) {
				c.WritePart("projects")
				c.WriteIntPart(123)
			},
			want: "projects/123",
		},
		{
			name: "three parts middle empty",
			fn: func(c *Composer) {
				c.WritePart("projects")
				c.WriteEmptyPart()
				c.WritePart("keys")
			},
			want: "projects//keys",
		},
		{
			name: "resource and id",
			fn: func(c *Composer) {
				c.WriteResource("projects", "123")
			},
			want: "projects/123",
		},
		{
			name: "resource and int64",
			fn: func(c *Composer) {
				c.WriteIntResource("projects", math.MaxInt64)
				c.WriteIntResource("keys", math.MinInt64)
			},
			want: "projects/9223372036854775807/keys/-9223372036854775808",
		},
		{
			name: "resource and uint64",
			fn: func(c *Composer) {
				c.WriteUintResource("projects", math.MaxUint64)
			},
			want: "projects/18446744073709551615",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Composer
			tt.fn(&c)
			if got := c.b.String(); got != tt.want {
				t.Errorf("Composer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkComposer(b *testing.B) {
	b.Run("PartByPart", func(b *testing.B) {
		var c Composer
		for i := 0; i < b.N; i++ {
			c.Reset()
			c.WritePart("projects")
			c.WritePart("123")
			c.WritePart("keys")
			c.WriteIntPart(456)
			c.WritePart("versions")
			c.WriteUintPart(789)
			_ = c.Name()
		}
	})

	b.Run("ResourceByResource", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var c Composer
			c.WriteResource("projects", "123")
			c.WriteIntResource("keys", 456)
			c.WriteUintResource("versions", 789)
			_ = c.Name()
		}
	})
}

func BenchmarkFmtCompose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("projects/%s/keys/%d/versions/%d", "123", 456, 789)
	}
}
