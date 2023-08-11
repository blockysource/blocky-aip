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

package parser

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/filtering/ast"
)

// MemberHandler is a function that modifies a argument context member literal.
// It is used to modify the dot separated fields of a member literal i.e. for a need of a float value.
// All the field no longer used needs to be freed by calling FreeFieldExpr function.
// What's more MemberExpr fields needs to be emptied by calling ClearFields method.
// This way the MemberExpr can be reused from pool.
type MemberHandler func(m *ast.MemberExpr) bool

// FreeFieldExpr frees the memory of a field expression.
// NOTE: Use this function only if the field expression is not used anymore.
//       On merging multiple field expressions into one, it is required to Free unused field expressions.
//	   	 Otherwise the memory will be leaked.
func FreeFieldExpr(fe ast.FieldExpr) {
	putFieldExpr(fe)
}

// FreeNameExpr frees the memory of a name expression.
// NOTE: Use this function only if the name expression is not used anymore.
//       On merging multiple name expressions into one, it is required to Free unused name expressions.
//       Otherwise the memory will be leaked.
func FreeNameExpr(nameExpr ast.NameExpr) {
	putNameExpr(nameExpr)
}

// ParseMemberInt is a function that modifies a argument context member literal to an int value.
func ParseMemberInt(m *ast.MemberExpr) bool {
	// The int value can be represented as:
	// - the integer part value
	// - the fractional part value which is either a number or a number with an exponent.
	// The exponent is a number with a leading 'e' or 'E' and an optional sign.
	// The exponent is a decimal number.
	// This requires a member literal to have a value and at most one field (fractional).
	if len(m.Fields) > 1 || m.Value == nil {
		return false
	}

	if m.Value.HasDecodedValue() {
		return false
	}

	value := m.Value.GetStringValue()
	if value == "" {
		return false
	}

	tv, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}
	m.Value.SetDecodedValue(tv) // int64
	return true
}

// ParseMemberNumber is a function that modifies a argument context member literal to a numeric value.
func ParseMemberNumber(m *ast.MemberExpr) bool {
	// The float value like: 1.0 or 1.0e-2 can be represented as:
	// - the integer part value
	// - the fractional part value which is either a number or a number with an exponent.
	// The exponent is a number with a leading 'e' or 'E' and an optional sign.
	// The exponent is a decimal number.
	// This requires a member literal to have a value and at most one field (fractional).
	if len(m.Fields) > 1 || m.Value == nil {
		return false
	}

	if m.Value.HasDecodedValue() {
		return false
	}

	// Check if the value is a number, which must be a *ast.TextLiteral.
	var tl *ast.TextLiteral
	switch mt := m.Value.(type) {
	case *ast.StringLiteral:
		if len(m.Fields) > 0 {
			return false
		}
		if strings.ContainsRune(mt.Value, '.') {
			fl, err := strconv.ParseFloat(mt.Value, 64)
			if err != nil {
				return false
			}
			mt.DecodedValue = fl // float64
			return true
		}
		tv, err := strconv.ParseInt(mt.Value, 10, 64)
		if err != nil {
			return false
		}
		mt.DecodedValue = tv // int64
	case *ast.TextLiteral:
		tl = mt
	}

	// Check if the value in the TextLiteral is a number.
	// If it is not a number, then return false.

	for i, r := range tl.Value {
		if i == 0 && r == '-' || r == '+' {
			continue
		}
		if r < '0' || r > '9' {
			return false
		}
	}

	// The first part is a number.
	if len(m.Fields) == 0 {
		tv, err := strconv.ParseInt(tl.Value, 10, 64)
		if err != nil {
			return false
		}
		tl.DecodedValue = tv // int64
		return true
	}

	// Verify if the first field is a fractional part.
	field := m.Fields[0]
	var ftl *ast.TextLiteral
	switch ft := field.(type) {
	case *ast.StringLiteral, *ast.KeywordExpr:
		return false
	case *ast.TextLiteral:
		ftl = ft
	}

	// Check if the fractional part is a number or a number with an exponent.
	// Otherwise return false.
	var (
		prevExp  bool
		validExp bool
	)
	for i, r := range ftl.Value {
		if r < '0' || r > '9' {
			if prevExp {
				validExp = true
				prevExp = false
			}
			continue
		}

		if prevExp && r == '-' || r == '+' {
			continue
		}

		// If the fractional part doesn't start with a number, then return false.
		if i == 0 {
			return false
		}

		if r == 'e' || r == 'E' {
			prevExp = true
			continue
		}

		return false
	}

	if prevExp && !validExp {
		return false
	}

	var sb strings.Builder
	sb.WriteString(tl.Value)
	sb.WriteRune('.')
	sb.WriteString(ftl.Value)
	res := sb.String()
	tv, err := strconv.ParseFloat(res, 64)
	if err != nil {
		return false
	}
	tl.Value = res
	tl.DecodedValue = tv // float64

	// Remove the fractional part field from slice.
	m.ClearFields()

	// release the memory of the fractional part text literal.
	putTextLiteral(ftl)

	return true
}

// ParseMemberBoolean is a function that modifies a argument context member literal to a boolean value.
func ParseMemberBoolean(m *ast.MemberExpr) bool {
	if len(m.Fields) > 0 || m.Value == nil {
		return false
	}

	if m.Value.HasDecodedValue() {
		return false
	}

	value := m.Value.GetStringValue()
	switch {
	case value == "true":
		m.Value.SetDecodedValue(true) // bool
	case value == "false":
		m.Value.SetDecodedValue(false) // bool
	default:
		return false
	}
	return true
}

// ParseMemberDuration is a function that modifies a argument context member literal to a duration value.
func ParseMemberDuration(m *ast.MemberExpr) bool {
	if len(m.Fields) > 1 || m.Value == nil {
		return false
	}

	if m.Value.HasDecodedValue() {
		return false
	}

	// Check if the value is a TextLiteral.
	var tl *ast.TextLiteral
	switch mt := m.Value.(type) {
	case *ast.StringLiteral:
		if len(m.Fields) > 0 {
			return false
		}
		d, err := time.ParseDuration(tl.Value)
		if err != nil {
			return false
		}

		tl.DecodedValue = d
		return true
	case *ast.TextLiteral:
		tl = mt
	}

	// Check if the value is a valid duration.
	if len(m.Fields) == 0 {
		d, err := time.ParseDuration(tl.Value)
		if err != nil {
			return false
		}

		tl.DecodedValue = d
		return true
	}

	// The second field needs to be a Text literal.
	var ftl *ast.TextLiteral
	switch ft := m.Fields[0].(type) {
	case *ast.StringLiteral, *ast.KeywordExpr:
		return false
	case *ast.TextLiteral:
		ftl = ft
	}

	var sb strings.Builder
	sb.WriteString(tl.Value)
	sb.WriteRune('.')
	sb.WriteString(ftl.Value)

	res := sb.String()
	d, err := time.ParseDuration(res)
	if err != nil {
		return false
	}

	tl.Value = res
	tl.DecodedValue = d

	// Remove the fractional part field from slice.
	putTextLiteral(ftl)

	// Remove all the fields from the member literal.
	m.ClearFields()
	return true
}

// ParseMemberTimestamps is a function that modifies a argument context member literal to a timestamp value.
func ParseMemberTimestamps(m *ast.MemberExpr) bool {
	if len(m.Fields) > 0 || m.Value == nil {
		return false
	}

	if m.Value.HasDecodedValue() {
		return false
	}

	// Check if the value is a TextLiteral.
	value := m.Value.GetStringValue()

	// Parse the timestamp using RFC3339 format.
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return false
	}

	m.Value.SetDecodedValue(t) // time.Time
	return true
}

// ParseMemberEnum is a function that modifies a argument context member literal to an enum value.
// If the input field descriptor is not an enum, then it panics.
// Resultant TypedValue is of type protoreflect.EnumNumber.
func ParseMemberEnum(fieldDesc protoreflect.FieldDescriptor) MemberHandler {
	if fieldDesc.Kind() != protoreflect.EnumKind {
		panic(fmt.Errorf("unsupported field type %s", fieldDesc.Kind()))
	}

	enumType := fieldDesc.Enum()
	return func(m *ast.MemberExpr) bool {
		if len(m.Fields) > 0 || m.Value == nil {
			return false
		}

		if m.Value.HasDecodedValue() {
			return false
		}

		// Check if the value is a TextLiteral.
		var sl *ast.StringLiteral
		switch mt := m.Value.(type) {
		case *ast.StringLiteral:
			sl = mt
		case *ast.TextLiteral:
			return false
		}

		// Check if the value is a valid enum value.
		enumValue := enumType.Values().ByName(protoreflect.Name(sl.Value))
		if enumValue == nil {
			return false
		}

		sl.DecodedValue = enumValue.Number() // protoreflect.EnumNumber
		return true
	}
}

// ParseMemberBytesBase64 is a function that modifies a argument context member literal to a bytes value.
func ParseMemberBytesBase64(urlEncoded bool) MemberHandler {
	return func(m *ast.MemberExpr) bool {
		if len(m.Fields) > 0 || m.Value == nil {
			return false
		}

		if m.Value.HasDecodedValue() {
			return false
		}

		str := m.Value.GetStringValue()
		if str == "" {
			return false
		}

		var (
			// Check if the value is a valid base64 encoded string.
			res []byte
			err error
		)
		if urlEncoded {
			res, err = base64.URLEncoding.DecodeString(str)
		} else {
			res, err = base64.StdEncoding.DecodeString(str)
		}
		if err != nil {
			return false
		}

		m.Value.SetDecodedValue(res) // []byte
		return true
	}
}

// ParseMemberBytesHex is a function that modifies a argument context member literal to a bytes value.
func ParseMemberBytesHex(m *ast.MemberExpr) bool {
	if len(m.Fields) > 0 || m.Value == nil {
		return false
	}

	// Check if the value is already decoded.
	if m.Value.HasDecodedValue() {
		return false
	}

	str := m.Value.GetStringValue()
	if str == "" {
		return false
	}

	res, err := hex.DecodeString(str)
	if err != nil {
		return false
	}

	m.Value.SetDecodedValue(res) // []byte
	return true
}

// ParseMemberNull is a function that modifies a argument context member literal match a null value.
func ParseMemberNull(m *ast.MemberExpr) bool {
	if len(m.Fields) > 0 || m.Value == nil {
		return false
	}

	// Check if the value is already decoded.
	if m.Value.HasDecodedValue() {
		return false
	}

	str := m.Value.GetStringValue()
	if str == "" {
		return false
	}

	if str != "null" {
		return false
	}

	m.Value.SetDecodedValue(nil) // nil
	return true
}
