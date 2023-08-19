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

package scanner

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func daysIn(month, year int) int {
	if month == 2 && isLeap(year) {
		return 29
	}
	return int(daysBefore[month] - daysBefore[month-1])
}

func parseUint(s string, min, max int) (x int, ok bool) {
	for _, c := range []byte(s) {
		if c < '0' || '9' < c {
			ok = false
			return min, ok
		}
		x = x*10 + int(c) - '0'
	}
	if x < min || max < x {
		ok = false
		return min, ok
	}
	return x, true
}

func isValidTimestamp(lit string) bool {
	if len(lit) < len("2006-01-02T15:04:05") {
		return false
	}

	if !(lit[4] == '-' && lit[7] == '-' && lit[10] == 'T' && lit[13] == ':' && lit[16] == ':') {
		return false
	}

	year, ok := parseUint(lit[0:4], 0, 9999) // e.g., 2006
	if !ok {
		return false
	}

	m, ok := parseUint(lit[5:7], 1, 12) // e.g., 01
	if !ok {
		return false
	}

	_, ok = parseUint(lit[8:10], 1, daysIn(m, year)) // e.g., 02
	if !ok {
		return false
	}

	_, ok = parseUint(lit[11:13], 0, 23) // e.g., 15
	if !ok {
		return false
	}

	_, ok = parseUint(lit[14:16], 0, 59) // e.g., 04
	if !ok {
		return false
	}

	lit = lit[19:]

	// Parse the fractional second.
	if len(lit) >= 2 && lit[0] == '.' && isDecimal(rune(lit[1])) {
		n := 2
		for ; n < len(lit) && isDecimal(rune(lit[n])); n++ {
		}
		lit = lit[n:]
	}

	// Parse the time zone.
	if len(lit) != 1 || lit[0] != 'Z' {
		if len(lit) != len("-07:00") {
			return false
		}
		_, ok = parseUint(lit[1:3], 0, 23) // e.g., 07
		if !ok {
			return false
		}

		_, ok = parseUint(lit[4:6], 0, 59) // e.g., 00
		if !ok {
			return false
		}

		if !((lit[0] == '-' || lit[0] == '+') && lit[3] == ':') {
			return false
		}
	}

	return true
}
