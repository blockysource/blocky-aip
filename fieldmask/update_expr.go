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

package fieldmask

import (
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/info"
	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

// ParseUpdateExpr parses a field mask, and extracts field values to update
// from the given message.
// Parsed field mask Update expressions, can be used to update selected fields of a message.
// If selected field is a map with a wildstar selector i.e: path: "map_field.*.field_name"
// This will evaluate to all keys in the mapi_field value
func (p *Parser) ParseUpdateExpr(msg proto.Message, mask *fieldmaskpb.FieldMask) (*expr.UpdateExpr, error) {
	if p.desc == nil {
		p.desc = msg.ProtoReflect().Descriptor()
		p.msgInfo = info.MapMsgInfo(p.desc)
	}
	ue := expr.AcquireUpdateExpr()
	if len(mask.Paths) == 0 {
		return ue, nil
	}

	pm := msg.ProtoReflect()
	for _, path := range mask.Paths {
		err := p.buildPathUpdateExpr(ue, pm, path)
		if err != nil {
			ue.Free()
			return nil, err
		}
	}
	return ue, nil
}

func (p *Parser) buildPathUpdateExpr(ue *expr.UpdateExpr, msgValue protoreflect.Message, path string) (err error) {
	var s scanner.Scanner
	s.Reset(path, p.errHandler)

	curMsg := msgValue

	if msgValue.Descriptor().FullName() != p.desc.FullName() {
		if p.errHandler != nil {
			p.errHandler(0, "invalid message descriptor")
		}
		return ErrInternalError
	}

	md := p.desc

	root := expr.AcquireFieldSelectorExpr()
	defer func() {
		if err != nil {
			root.Free()
		}
	}()

	fs := root
	fs.Message = msgValue.Descriptor().FullName()
	var fi info.FieldInfo
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			err = p.handleLastPathElem(ue, curMsg, fi, root, fs, pos)
			if err != nil {
				root.Free()
				return err
			}
			return nil
		}
		if tok == token.PERIOD {
			// This means it is an extra period.
			// Return an error.
			if p.errHandler != nil {
				p.errHandler(pos, "unexpected period")
			}
			root.Free()
			return ErrInvalidField
		}

		// This means it is a field selector.
		var ok bool
		fi, ok = p.msgInfo.MessageInfo(md).
			FieldByName(protoreflect.Name(lit))
		if !ok {
			if p.errHandler != nil {
				p.errHandler(pos, "unknown field name")
			}
			root.Free()
			return ErrInvalidField
		}

		// Ensure that the field is not immutable or output only.
		if fi.Immutable || fi.OutputOnly {
			if p.ignoreNonUpdatable {
				// Finish up this path without any expressions.
				root.Free()
				return nil
			}
			if p.errHandler != nil {
				p.errHandler(pos, "immutable or output only field in sub path cannot be updated")
			}
			root.Free()
			return ErrInvalidField
		}

		// Make the next field selector expression.
		fs.Field = protoreflect.Name(lit)

		// Check if the next is a period.
		var isPeriod bool
		s.Peek(func(p token.Position, t token.Token, l string) bool {
			isPeriod = t == token.PERIOD
			return isPeriod
		})

		// If the next token is a period ensure that this field is a message (or a map - which is also a message).
		if isPeriod {
			// Check if the field is a message.
			if fi.Desc.Kind() != protoreflect.MessageKind {
				if p.errHandler != nil {
					p.errHandler(pos, "expected message field in sub path")
				}
				root.Free()
				return ErrInvalidField
			}

			if fi.Desc.IsMap() {
				// Scan the map key value.
				pos, tok, lit = s.Scan()

				// Verify if the token is a valid literal.
				switch tok {
				case token.EOF:
					if p.errHandler != nil {
						p.errHandler(pos, "expected map key")
					}
					return ErrInvalidField
				case token.ASTERISK:
					// An asterisk is a wildcard selector.
					// This means we need to add all the values of the map keys recursively.
					// TODO: Implement this.
				}

				// Search for the next period to check whether the selector is a map key or it has subsequent elements.
				s.Peek(func(p token.Position, t token.Token, l string) bool {
					isPeriod = t == token.PERIOD
					return isPeriod
				})

				// Get the value of the map.
				mp := curMsg.Get(fi.Desc).Map()

				// There is no more sub paths, which means the selector is a map key.
				if !isPeriod {
					err = p.handleLastMapKeyElem(ue, root, fs, fi, mp, tok, pos, lit)
					if err != nil {
						return err
					}
					return nil
				}

				// This is a valid map key selector now.
				mke := expr.AcquireMapKeyExpr()
				var mkv *expr.ValueExpr
				switch mk := fi.Desc.MapKey(); mk.Kind() {
				case protoreflect.BoolKind:
					if !tok.IsBoolean() {
						if p.errHandler != nil {
							p.errHandler(pos, "expected boolean value")
						}

						return ErrInvalidField
					}
					mkv = expr.AcquireValueExpr()
					mkv.Value = lit == "true"
					mke.Key = mkv
				case protoreflect.StringKind:
					if tok != token.STRING && !tok.IsIdent() {
						if p.errHandler != nil {
							p.errHandler(pos, "expected string value")
						}
						return ErrInvalidField
					}
					mkv = expr.AcquireValueExpr()
					mkv.Value = lit
					mke.Key = mkv
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
					protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					if !tok.IsInteger() {
						if p.errHandler != nil {
							p.errHandler(pos, "expected integer value")
						}
						return ErrInvalidField
					}

					iv, err := strconv.ParseInt(lit, 10, 64)
					if err != nil {
						if p.errHandler != nil {
							p.errHandler(pos, "invalid integer value")
						}
						return ErrInvalidField
					}
					mkv = expr.AcquireValueExpr()
					mkv.Value = iv
					mke.Key = mkv
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					if !tok.IsInteger() {
						if p.errHandler != nil {
							p.errHandler(pos, "expected unsigned integer value")
						}
						return ErrInvalidField
					}
					iv, err := strconv.ParseUint(lit, 10, 64)
					if err != nil {
						if p.errHandler != nil {
							p.errHandler(pos, "invalid unsigned integer value")
						}
						return ErrInvalidField
					}

					mkv = expr.AcquireValueExpr()
					mkv.Value = iv
					mke.Key = mkv
				default:
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("unsupported map key type: %s", mk.Kind()))
					}
					return ErrInvalidField
				}
				// Change current message descriptor to be the map value descriptor.
				md = fi.Desc.MapValue().Message()

				// We've found the path map key, now we need to ensure if this key exists in the map value.
				mv := mp.Get(protoreflect.ValueOf(mkv.Value).MapKey())
				if !mv.IsValid() {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("field: %q, map key: %q not found in the input message", fi.Desc.Name(), mkv.Value))
					}
					return ErrInvalidField
				}

				// If it does change current context message value.
				curMsg = mv.Message()
				// Set the traversal of the last field selector to be the map key expression.
				fs.Traversal = mke
				continue
			}

			// This is a message field.
			md = fi.Desc.Message()
			// Change current field selector to a new one.
			nf := expr.AcquireFieldSelectorExpr()
			nf.Message = md.FullName()

			// Ensure if the last field selector has already a MapKey traversal.
			switch ft := fs.Traversal.(type) {
			case nil:
				// If the traversal is not already set, then set it to the new field selector.
				fs.Traversal = nf
			case *expr.MapKeyExpr:
				// If the traversal is a map key, then set the map key traversal to the new field selector.
				ft.Traversal = nf
			}

			// Change current context field selector to the new one.
			fs = nf

			// Get the value of the message field, and change current context message.
			curMsg = curMsg.Get(fi.Desc).Message()
		}
	}
}

func (p *Parser) handleLastPathElem(ue *expr.UpdateExpr, curMsg protoreflect.Message, fi info.FieldInfo, root, fs *expr.FieldSelectorExpr, pos token.Position) (err error) {
	// If this is the last element of the path, then we need to extract the value of the field.
	fv := curMsg.Get(fi.Desc)

	// Extract the value out of the message.
	switch fi.Desc.Kind() {
	case protoreflect.MessageKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			if fi.Desc.IsMap() {
				mve := expr.AcquireMapValueExpr()
				fvm := fv.Map()

				mk := fi.Desc.MapKey()
				mv := fi.Desc.MapValue()
				var err error
				fvm.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
					var mkv *expr.ValueExpr
					switch mk.Kind() {
					case protoreflect.BoolKind:
						ve := expr.AcquireValueExpr()
						ve.Value = k.Bool()
						mkv = ve
					case protoreflect.StringKind:
						ve := expr.AcquireValueExpr()
						ve.Value = k.String()
						mkv = ve
					case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
						protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
						ve := expr.AcquireValueExpr()
						ve.Value = k.Int()
						mkv = ve
					case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
						ve := expr.AcquireValueExpr()
						ve.Value = k.Uint()
						mkv = ve
					}

					fkc := fs.Clone().(*expr.FieldSelectorExpr)

					mke := expr.AcquireMapKeyExpr()
					mke.Key = mkv
					fkc.Traversal = mke

					var mvv expr.FilterExpr
					switch mv.Kind() {
					case protoreflect.MessageKind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Message()
						mvv = ve
					case protoreflect.BoolKind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Bool()
						mvv = ve
					case protoreflect.StringKind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.String()
						mvv = ve
					case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
						protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Int()
						mvv = ve
					case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Uint()
						mvv = ve
					case protoreflect.FloatKind, protoreflect.DoubleKind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Float()
						mvv = ve
					case protoreflect.BytesKind:
						ve := expr.AcquireValueExpr()
						v.IsValid()
						ve.Value = v.Bytes()
						mvv = ve
					case protoreflect.EnumKind:
						ve := expr.AcquireValueExpr()
						ve.Value = v.Enum()
						mvv = ve
					default:
						if p.errHandler != nil {
							p.errHandler(pos, "unsupported field type")
						}
						mkv.Free()
						err = ErrInvalidField
						return false
					}

					mve.Values = append(mve.Values, expr.MapValueExprEntry{
						Key:   mkv,
						Value: mvv,
					})
					return true
				})
				if err != nil {
					return err
				}
				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: root,
					Value: mve,
				})
				return nil
			}
			ae := expr.AcquireArrayUpdateExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				elem := ls.Get(i)
				if !elem.IsValid() {
					ae.Elements = append(ae.Elements, nil)
					continue
				}

				subUe := expr.AcquireUpdateExpr()
				if err = p.addMsgAllFieldsExpr(subUe, elem.Message()); err != nil {
					ae.Free()
					return err
				}
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})

		} else {
			// If the field value is not valid, it means that the field is nil.
			if !fv.IsValid() {
				// If the field is not nullable, then this is an error.
				if !fi.Nullable {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("field: %s is not nullable", fi.Desc.Name()))
					}
					return ErrInvalidField
				}

				// Otherwise make an update value to be nil.
				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: root,
					Value: expr.AcquireValueExpr(),
				})
			}

			// If the field is a valid message, then we create a sub UpdateExpression.
			// This is done recursively.
			switch {
			case fi.IsTimestamp:
				// If the field is a timestamp, then we need to convert it to a time.Time.
				ve := expr.AcquireValueExpr()
				ve.Value = p.extractTimeValue(fv.Message())
				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: root,
					Value: ve,
				})

			case fi.IsDuration:
				// If the field is a duration, then we need to convert it to a time.Duration.
				ve := expr.AcquireValueExpr()
				ve.Value = p.extractDurationValue(fv.Message())
				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: root,
					Value: ve,
				})
			default:
				subUe := expr.AcquireUpdateExpr()

				if err = p.addMsgAllFieldsExpr(subUe, fv.Message()); err != nil {
					subUe.Free()
					return err
				}

				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: root,
					Value: subUe,
				})
			}
		}
	case protoreflect.BoolKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Bool()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Bool()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.StringKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).String()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.String()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Int()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Int()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Uint()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Uint()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Float()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Float()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.BytesKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Bytes()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Bytes()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	case protoreflect.EnumKind:
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			ae := expr.AcquireArrayExpr()
			ls := fv.List()
			for i := 0; i < ls.Len(); i++ {
				ve := expr.AcquireValueExpr()
				ve.Value = ls.Get(i).Enum()
				ae.Elements = append(ae.Elements, ve)
			}
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ae,
			})
		} else {
			ve := expr.AcquireValueExpr()
			ve.Value = fv.Enum()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: root,
				Value: ve,
			})
		}
	default:
		if p.errHandler != nil {
			p.errHandler(pos, "unsupported field type")
		}
		return ErrInvalidField
	}
	return nil
}

func (p *Parser) handleLastMapKeyElem(ue *expr.UpdateExpr, root, fs *expr.FieldSelectorExpr, fi info.FieldInfo, mp protoreflect.Map, tok token.Token, pos token.Position, lit string) error {
	var (
		mkv protoreflect.MapKey
		mvv protoreflect.Value
	)

	ke := expr.AcquireMapKeyExpr()
	mk := fi.Desc.MapKey()
	fs.Field = fi.Desc.Name()
	// A map key can only be a string, Int, Uint, Bool.
	// It cannot be a float, double, message, bytes or enum.
	switch mk.Kind() {
	case protoreflect.BoolKind:
		if !tok.IsBoolean() {
			if p.errHandler != nil {
				p.errHandler(pos, "expected boolean value")
			}

			return ErrInvalidField
		}
		v := lit == "true"
		ve := expr.AcquireValueExpr()
		ve.Value = v
		ke.Key = ve
		mkv = protoreflect.ValueOf(v).MapKey()
		mvv = mp.Get(mkv)
	case protoreflect.StringKind:
		if tok != token.STRING && !tok.IsIdent() {
			if p.errHandler != nil {
				p.errHandler(pos, "expected string value")
			}
			return ErrInvalidField
		}
		ve := expr.AcquireValueExpr()
		ve.Value = lit
		ke.Key = ve
		mkv = protoreflect.ValueOf(lit).MapKey()
		mvv = mp.Get(mkv)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if !tok.IsInteger() {
			if p.errHandler != nil {
				p.errHandler(pos, "expected integer value")
			}
			return ErrInvalidField
		}
		iv, err := strconv.ParseInt(lit, 10, 64)
		if err != nil {
			if p.errHandler != nil {
				p.errHandler(pos, "invalid integer value")
			}
			return ErrInvalidField
		}
		ve := expr.AcquireValueExpr()
		ve.Value = iv
		ke.Key = ve
		mkv = protoreflect.ValueOf(iv).MapKey()
		mvv = mp.Get(mkv)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if !tok.IsInteger() {
			if p.errHandler != nil {
				p.errHandler(pos, "expected unsigned integer value")
			}
			return ErrInvalidField
		}
		iv, err := strconv.ParseUint(lit, 10, 64)
		if err != nil {
			if p.errHandler != nil {
				p.errHandler(pos, "invalid unsigned integer value")
			}
			return ErrInvalidField
		}
		ve := expr.AcquireValueExpr()
		ve.Value = iv
		ke.Key = ve
		mkv = protoreflect.ValueOf(iv).MapKey()
		mvv = mp.Get(mkv)
	default:
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("unsupported map key type: %s", mk.Kind()))
		}
		return ErrInvalidField
	}
	fs.Traversal = ke

	if !mvv.IsValid() {
		if p.errHandler != nil {
			p.errHandler(pos, "map key not found in the input message")
		}
		return ErrInvalidField
	}

	var fv expr.UpdateValueExpr
	switch fi.Desc.MapValue().Kind() {
	case protoreflect.MessageKind:
		// This is a special case where the field traversal contains a map key, and each message field is a different
		// update expression.
		mke := expr.AcquireMapKeyExpr()
		kve := expr.AcquireValueExpr()
		switch mk.Kind() {
		case protoreflect.BoolKind:
			kve.Value = mkv.Bool()
		case protoreflect.StringKind:
			kve.Value = mkv.String()
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			kve.Value = mkv.Int()
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			kve.Value = mkv.Uint()
		}
		mke.Key = kve
		fs.Traversal = mke

		subUe := expr.AcquireUpdateExpr()
		if err := p.addMsgAllFieldsExpr(subUe, mvv.Message()); err != nil {
			subUe.Free()
			return err
		}

		// Add the sub update expression to the update expression.
		ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
			Field: root,
			Value: subUe,
		})
	case protoreflect.BoolKind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Bool()
		fv = ve
	case protoreflect.StringKind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.String()
		fv = ve
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Int()
		fv = ve
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Uint()
		fv = ve
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Float()
		fv = ve
	case protoreflect.BytesKind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Bytes()
		fv = ve
	case protoreflect.EnumKind:
		ve := expr.AcquireValueExpr()
		ve.Value = mvv.Enum()
		fv = ve
	default:
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("unsupported map value type: %s", fi.Desc.MapValue().Kind()))
		}
		return ErrInternalError
	}

	ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
		Field: root,
		Value: fv,
	})
	return nil
}

func (p *Parser) addMsgAllFieldsExpr(ue *expr.UpdateExpr, subV protoreflect.Message) error {
	msg := subV.Descriptor()

	mi := p.msgInfo.MessageInfo(msg)
	for _, fi := range mi.Fields {
		// We don't want to update immutable or output only fields.
		if fi.Immutable || fi.OutputOnly {
			continue
		}

		noValue := !subV.Has(fi.Desc)
		if noValue {
			if fi.IsOneOf {
				continue
			}

			if !fi.Nullable && !fi.NonEmptyDefault {
				continue
			}

		}

		v := subV.Get(fi.Desc)

		var uv expr.UpdateValueExpr
		if fi.Desc.Cardinality() == protoreflect.Repeated {
			if fi.Nullable && !v.IsValid() {
				fs := expr.AcquireFieldSelectorExpr()
				fs.Message = msg.FullName()
				fs.Field = fi.Desc.Name()
				ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
					Field: fs,
					Value: expr.AcquireValueExpr(),
				})
				continue
			}

			switch fi.Desc.Kind() {
			case protoreflect.MessageKind:
				if fi.IsTimestamp || fi.IsDuration {
					ae := expr.AcquireArrayExpr()
					ls := v.List()
					for i := 0; i < ls.Len(); i++ {
						elem := ls.Get(i)
						if !elem.IsValid() {
							ve := expr.AcquireValueExpr()
							switch {
							case fi.IsTimestamp:
								ve.Value = time.Time{}
							case fi.IsDuration:
								ve.Value = time.Duration(0)
							}
							ae.Elements = append(ae.Elements, ve)
							continue
						}
						ve := expr.AcquireValueExpr()
						switch {
						case fi.IsTimestamp:
							ve.Value = p.extractTimeValue(elem.Message())
						case fi.IsDuration:
							ve.Value = p.extractDurationValue(elem.Message())
						}
						ae.Elements = append(ae.Elements, ve)
					}

					fs := expr.AcquireFieldSelectorExpr()
					fs.Message = msg.FullName()
					fs.Field = fi.Desc.Name()
					ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
						Field: fs,
						Value: ae,
					})
				}

				if !fi.Desc.IsMap() {
					if noValue {
						// Add the null
						fs := expr.AcquireFieldSelectorExpr()
						fs.Message = msg.FullName()
						fs.Field = fi.Desc.Name()
						ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
							Field: fs,
							Value: expr.AcquireValueExpr(),
						})
						continue
					}
					aue := expr.AcquireArrayUpdateExpr()
					ls := v.List()

					for i := 0; i < ls.Len(); i++ {

						sub := ls.Get(i)
						if !sub.IsValid() {
							// Provide a nil, element into expressions.
							aue.Elements = append(aue.Elements, nil)
							continue
						}

						subUe := expr.AcquireUpdateExpr()
						if err := p.addMsgAllFieldsExpr(subUe, sub.Message()); err != nil {
							aue.Free()
							return err
						}
						aue.Elements = append(aue.Elements, subUe)
					}

					fs := expr.AcquireFieldSelectorExpr()
					fs.Message = msg.FullName()
					fs.Field = fi.Desc.Name()

					ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
						Field: fs,
						Value: aue,
					})
					continue
				} else {
					if noValue {
						continue
					}
					mp := v.Map()
					var err error
					mp.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
						var mkv *expr.ValueExpr
						switch mk := fi.Desc.MapKey(); mk.Kind() {
						case protoreflect.BoolKind:
							ve := expr.AcquireValueExpr()
							ve.Value = k.Bool()
							mkv = ve
						case protoreflect.StringKind:
							ve := expr.AcquireValueExpr()
							ve.Value = k.String()
							mkv = ve
						case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
							protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
							ve := expr.AcquireValueExpr()
							ve.Value = k.Int()
							mkv = ve
						case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
							ve := expr.AcquireValueExpr()
							ve.Value = k.Uint()
							mkv = ve
						}

						fs := expr.AcquireFieldSelectorExpr()
						fs.Message = msg.FullName()
						fs.Field = fi.Desc.Name()

						mke := expr.AcquireMapKeyExpr()
						mke.Key = mkv
						fs.Traversal = mke

						switch mv := fi.Desc.MapValue(); mv.Kind() {
						case protoreflect.MessageKind:
							switch mv.Message().FullName() {
							case "google.protobuf.Timestamp":
								ve := expr.AcquireValueExpr()
								ve.Value = p.extractTimeValue(v.Message())
								ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
									Field: fs,
									Value: ve,
								})
							case "google.protobuf.Duration":
								ve := expr.AcquireValueExpr()
								ve.Value = p.extractDurationValue(v.Message())
								ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
									Field: fs,
								})
							default:
								subUe := expr.AcquireUpdateExpr()
								if err = p.addMsgAllFieldsExpr(subUe, v.Message()); err != nil {
									subUe.Free()
									fs.Free()
									return false
								}

								ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
									Field: fs,
									Value: subUe,
								})
							}
						case protoreflect.BoolKind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Bool()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						case protoreflect.StringKind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.String()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
							protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Int()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Uint()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						case protoreflect.FloatKind, protoreflect.DoubleKind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Float()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						case protoreflect.BytesKind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Bytes()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
							})
						case protoreflect.EnumKind:
							ve := expr.AcquireValueExpr()
							ve.Value = v.Enum()
							ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
								Field: fs,
								Value: ve,
							})
						}
						return true
					})
					if err != nil {
						return err
					}
					continue
				}
			case protoreflect.BoolKind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Bool()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.StringKind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.String()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
				protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Int()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Uint()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.FloatKind, protoreflect.DoubleKind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Float()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.BytesKind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Bytes()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			case protoreflect.EnumKind:
				ve := expr.AcquireArrayExpr()
				for i := 0; i < v.List().Len(); i++ {
					elem := v.List().Get(i)
					vee := expr.AcquireValueExpr()
					vee.Value = elem.Enum()
					ve.Elements = append(ve.Elements, vee)
				}
				uv = ve
			default:
				if p.errHandler != nil {
					p.errHandler(0, fmt.Sprintf("unsupported field type: %s", fi.Desc.Kind()))
				}
				return ErrInternalError
			}
			fs := expr.AcquireFieldSelectorExpr()
			fs.Message = msg.FullName()
			fs.Field = fi.Desc.Name()
			ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
				Field: fs,
				Value: uv,
			})
			continue
		}

		// We treat a regular field to update.
		switch fi.Desc.Kind() {
		case protoreflect.MessageKind:
			if !v.IsValid() {
				ve := expr.AcquireValueExpr()
				if fi.IsTimestamp {
					ve.Value = time.Time{}
				} else if fi.IsDuration {
					ve.Value = time.Duration(0)
				}
				uv = ve
			} else {

				switch {
				case fi.IsTimestamp:
					ve := expr.AcquireValueExpr()
					if v.IsValid() {
						ve.Value = p.extractTimeValue(v.Message())
					} else {
						if fi.Nullable {
							ve.Value = nil
						} else {
							ve.Value = time.Time{}
						}
					}
					uv = ve
				case fi.IsDuration:
					ve := expr.AcquireValueExpr()
					if v.IsValid() {
						ve.Value = p.extractDurationValue(v.Message())
					} else {
						if fi.Nullable {
							ve.Value = nil
						} else {
							ve.Value = time.Duration(0)
						}
					}
					uv = ve
				default:
					if noValue {
						// Add the null
						fs := expr.AcquireFieldSelectorExpr()
						fs.Message = msg.FullName()
						fs.Field = fi.Desc.Name()
						ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
							Field: fs,
							Value: expr.AcquireValueExpr(),
						})
						continue
					}
					subUe := expr.AcquireUpdateExpr()
					if err := p.addMsgAllFieldsExpr(subUe, v.Message()); err != nil {
						subUe.Free()
						return err
					}
					uv = subUe
				}
			}
		case protoreflect.BoolKind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Bool()
			uv = ve
		case protoreflect.StringKind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.String()
			uv = ve
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Int()
			uv = ve
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Uint()
			uv = ve
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Float()
			uv = ve
		case protoreflect.BytesKind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Bytes()
			uv = ve
		case protoreflect.EnumKind:
			ve := expr.AcquireValueExpr()
			ve.Value = v.Enum()
			uv = ve
		default:
			if p.errHandler != nil {
				p.errHandler(0, fmt.Sprintf("unsupported field type: %s", fi.Desc.Kind()))
			}
			return ErrInternalError
		}

		// This is a top level of msg field.
		rc := expr.AcquireFieldSelectorExpr()
		rc.Message = msg.FullName()
		rc.Field = fi.Desc.Name()

		ue.Elements = append(ue.Elements, expr.UpdateFieldValue{
			Field: rc,
			Value: uv,
		})
	}
	return nil
}
