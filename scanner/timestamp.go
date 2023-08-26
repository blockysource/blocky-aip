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

import (
	"strconv"

	"github.com/blockysource/blocky-aip/token"
)

func (s *Scanner) scanTimestamp(used int) (token.Token, string) {
	offset := s.offset - used

	sum := used - 1
	// Format: 2006-01-02T15:04:05Z07:00
	for used < 4 {
		ch, w := s.next()
		if isBreaking(ch) {
			break
		}
		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}
	}
	year, err := strconv.Atoi(s.src[offset : offset+sum])
	if err != nil {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 4 {
		ch, w := s.next()
		if isBreaking(ch) {
			return token.INT, s.src[offset : offset+sum]
		}

		if ch != '-' {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}
		sum += w
	}

	var (
		m1 rune
		mm int
	)
	for used < 7 {
		ch, w := s.next()
		if isBreaking(ch) {
			break
		}
		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 6 {
			m1 = ch

			if m1 != '0' && m1 != '1' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			if m1 == '1' {
				mm = 10
			}
			continue
		}

		if m1 == '1' {
			if ch < '0' || ch > '2' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
		}

		mm += int(ch - '0')
	}

	if mm > 12 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 7 {
		ch, w := s.next()
		if isBreaking(ch) {
			return token.ILLEGAL, ""
		}

		if ch != '-' {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}
		sum += w
		used++
	}

	var dd int
	for used < 10 {
		ch, w := s.next()
		if isBreaking(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 9 {
			switch ch {
			case '0', '1', '2', '3':
			default:
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			dd = int(ch-'0') * 10
			continue
		}

		dd += int(ch - '0')
	}

	if dd > 31 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if daysIn(mm, year) < dd {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 10 {
		ch, w := s.next()
		if isBreaking(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if ch != 'T' && ch != ' ' {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}
		sum += w
		used++
	}

	var hh int
	for used < 13 {
		ch, w := s.next()
		if isBreaking(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 12 {
			if ch < '0' || ch > '2' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			hh = int(ch-'0') * 10
			continue
		}

		hh += int(ch - '0')
	}

	if hh > 23 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 13 {
		ch, w := s.next()
		if ch != ':' {
			if isBreaking(ch) {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
		}

		sum += w
		used++
	}

	var mi int
	for used < 16 {
		ch, w := s.next()
		if isBreaking(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 15 {
			if ch < '0' || ch > '5' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			mi = int(ch-'0') * 10
			continue
		}

		mi += int(ch - '0')
	}

	if mi > 59 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 16 {
		ch, w := s.next()
		if ch != ':' {
			if isBreaking(ch) {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
		}

		sum += w
		used++
	}

	var ss int
	for used < 19 {
		ch, w := s.next()
		if isBreaking(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 17 {
			if ch < '0' || ch > '5' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			ss = int(ch-'0') * 10
			continue
		}

		ss += int(ch - '0')
	}

	if ss > 59 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 19 {
		ch, w := s.next()
		if isBreaking(ch) {
			return token.ILLEGAL, ""
		}
		sum += w
		used++

		switch ch {
		case 'Z':
			_, w = s.next()
			sum += w
			return token.TIMESTAMP, s.src[offset : offset+sum]
		case '+', '-':
		default:
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}
	}

	var tzhh int
	for used < 22 {
		ch, w := s.next()
		if isEOF(ch) || isWhitespace(ch) || ch == ')' || ch == ',' || s.isComparator(ch) {
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 21 {
			if ch < '0' || ch > '1' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			tzhh = int(ch-'0') * 10
			continue
		}

		tzhh += int(ch - '0')
	}

	if tzhh > 23 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	if used == 22 {
		ch, w := s.next()
		if ch != ':' {
			if isBreaking(ch) {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
		}

		sum += w
		used++
	}

	var tzmm int
	for used < 25 {
		ch, w := s.next()
		if isBreaking(ch) {
			return token.ILLEGAL, ""
		}

		sum += w
		used++
		if !isDecimal(ch) {
			s.error(s.offset, "invalid timestamp")
			return token.ILLEGAL, ""
		}

		if used == 24 {
			if ch < '0' || ch > '5' {
				s.error(s.offset, "invalid timestamp")
				return token.ILLEGAL, ""
			}
			tzmm = int(ch-'0') * 10
			continue
		}

		tzmm += int(ch - '0')
	}

	if tzmm > 59 {
		s.error(s.offset, "invalid timestamp")
		return token.ILLEGAL, ""
	}

	_, w := s.next()
	sum += w

	return token.TIMESTAMP, s.src[offset : offset+sum]
}

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
