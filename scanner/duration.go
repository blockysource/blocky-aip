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
	"github.com/blockysource/blocky-aip/token"
)

func (s *Scanner) scanDuration(used int, isFloating, isExp bool) (token.Token, string) {
	offset := s.offset - used

	sum := used
	if isFloating && s.ch == 'n' {
		// cannot obtain fractional part of nanoseconds
		s.error(offset, "invalid duration")
		return token.ILLEGAL, ""
	}

	var topPow rune

	switch {
	case s.ch == 'm':
		if s.peek() == 's' {
			topPow = -3
		} else {
			topPow = 1
		}
	default:
		topPow = prefixPower[s.ch]
	}

	if isComposedDurationPrefix(s.ch) {
		if s.ch == 'm' {
			peek := s.peek()
			switch {
			case peek == 's':
				_, w := s.next()
				sum += w
			case isEOF(peek), isWhitespace(peek), peek == ')', peek == ',', s.isComparator(peek), peek == '}', peek == ']':
				s.next()
				return token.DURATION, s.src[offset : offset+sum]
			case !isDecimal(peek):
				s.error(offset, "invalid duration")
				return token.ILLEGAL, ""
			}
		} else {
			ch, w := s.next()
			if ch != 's' {
				s.error(offset, "invalid duration")
				return token.ILLEGAL, ""
			}
			sum += w
		}
	}

	for {
		ch, w := s.next()
		if isEOF(ch) || isWhitespace(ch) || ch == ')' || ch == ',' || s.isComparator(ch) ||
			(ch == '}') || (ch == ']') {
			break
		}

		if isFloating || topPow == -9 {
			// cannot obtain fractional part of nanoseconds, it should end up the duration
			s.error(offset, "invalid duration")
			return token.ILLEGAL, ""
		}

		// Check if there is an additional period, which is not allowed.
		// No matter if it is floating or not, the duration must finish with a suffix.
		if isPeriod(ch) {
			// cannot obtain fractional part of nanoseconds
			s.error(offset, "invalid duration")
			return token.ILLEGAL, ""
		}

		sum += w
		if isDecimal(ch) {
			continue
		}

		if !isDurationPrefix(ch) {
			s.error(offset, "invalid duration")
			return token.ILLEGAL, ""
		}

		pow := prefixPower[ch]
		if pow >= topPow {
			s.error(offset, "invalid duration")
			return token.ILLEGAL, ""
		}

		if isComposedDurationPrefix(s.ch) {
			if s.ch == 'm' {
				peek := s.peek()
				switch {
				case peek == 's':
					_, w := s.next()
					sum += w
				case isEOF(peek), isWhitespace(peek), peek == ')', peek == ',', s.isComparator(peek), peek == '}', peek == ']':
					s.next()
					return token.DURATION, s.src[offset : offset+sum]
				case !isDecimal(peek):
					s.error(offset, "invalid duration")
					return token.ILLEGAL, ""
				}
			} else {
				ch, w := s.next()
				if ch != 's' {
					s.error(offset, "invalid duration")
					return token.ILLEGAL, ""
				}
				sum += w
			}
		}
	}
	return token.DURATION, s.src[offset : offset+sum]
}

func isDurationPrefix(ch rune) bool {
	return ch == 'n' || ch == 'u' || ch == 'µ' || ch == 'μ' || ch == 'm' || ch == 's' || ch == 'h'
}

func isComposedDurationPrefix(ch rune) bool {
	return ch == 'm' || ch == 'u' || ch == 'μ' || ch == 'n'
}

var prefixPower = [...]rune{
	'n': -9,
	'u': -6,
	'µ': -6,
	's': 0,
	'h': 2,
}
