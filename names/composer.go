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
	"strings"
)

// Composer is a resource name composer.
type Composer struct {
	b   strings.Builder
	idx int
}

// Name returns the composed resource name.
func (c *Composer) Name() string {
	return c.b.String()
}

// WriteResource writes the resource plural name along with the identifier.
func (c *Composer) WriteResource(resource, id string) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.b.WriteString(resource)
	c.b.WriteByte('/')
	c.b.WriteString(id)
	c.idx += 2
}

// WriteIntResource writes the resource plural name along with the identifier.
func (c *Composer) WriteIntResource(resource string, id int) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.b.WriteString(resource)
	c.b.WriteByte('/')
	c.writeUint(uint64(id), id < 0)
	c.idx += 2
}

// WriteUintResource writes the resource plural name along with the identifier.
func (c *Composer) WriteUintResource(resource string, id uint) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.b.WriteString(resource)
	c.b.WriteByte('/')
	c.writeUint(uint64(id), false)
	c.idx += 2
}

// WritePart writes a resource name part.
func (c *Composer) WritePart(part string) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.b.WriteString(part)
	c.idx++
}

// WriteIntPart writes a resource name part from an integer.
func (c *Composer) WriteIntPart(part int) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.writeUint(uint64(part), part < 0)
	c.idx++
}

// WriteUintPart writes a resource name part from an unsigned integer.
func (c *Composer) WriteUintPart(part uint) {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.writeUint(uint64(part), false)
	c.idx++
}

// WriteEmptyPart writes an empty resource name part.
func (c *Composer) WriteEmptyPart() {
	if c.idx > 0 {
		c.b.WriteByte('/')
	}
	c.idx++
}

// Reset resets the composer.
func (c *Composer) Reset() {
	c.b.Reset()
	c.idx = 0
}

func (c *Composer) writeUint(iv uint64, neg bool) {
	// common case: use constants for / because
	// the compiler can optimize it into a multiply+shift
	if neg {
		c.b.WriteByte('-')
		iv = -iv
	}

	var a [64 + 1]byte // +1 for sign of 64bit value in base 2
	i := len(a)

	u := iv
	if host32bit {
		// convert the lower digits using 32bit operations
		for u >= 1e9 {
			// Avoid using r = a%b in addition to q = a/b
			// since 64bit division and modulo operations
			// are calculated by runtime functions on 32bit machines.
			q := u / 1e9
			us := uint(u - q*1e9) // u % 1e9 fits into a uint
			for j := 4; j > 0; j-- {
				is := us % 100 * 2
				us /= 100
				i -= 2
				a[i+1] = smallsString[is+1]
				a[i+0] = smallsString[is+0]
			}

			// us < 10, since it contains the last digit
			// from the initial 9-digit us.
			i--
			a[i] = smallsString[us*2+1]

			u = q
		}
		// u < 1e9
	}

	// u guaranteed to fit into a uint
	us := uint(u)
	for us >= 100 {
		is := us % 100 * 2
		us /= 100
		i -= 2
		a[i+1] = smallsString[is+1]
		a[i+0] = smallsString[is+0]
	}

	// us < 100
	is := us * 2
	i--
	a[i] = smallsString[is+1]
	if us >= 10 {
		i--
		a[i] = smallsString[is]
	}

	c.b.Write(a[i:])
}

const smallsString = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

const host32bit = ^uint(0)>>32 == 0
