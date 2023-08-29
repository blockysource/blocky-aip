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

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/blockysource/blocky-aip/expr"
	"github.com/blockysource/blocky-aip/internal/info"
	"github.com/blockysource/blocky-aip/scanner"
	"github.com/blockysource/blocky-aip/token"
)

// ParseSelectExpr parses a select expression from the field mask.
// All the common sub paths are consolidated into a single sub MessageSelectExpr for
func (p *Parser) ParseSelectExpr(fm *fieldmaskpb.FieldMask) (*expr.MessageSelectExpr, error) {
	if len(fm.GetPaths()) == 0 {
		return nil, nil
	}

	paths := fm.GetPaths()

	// Consolidate equal paths.
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if paths[i] == paths[j] {
				paths[j] = paths[len(paths)-1]
				paths = paths[:len(paths)-1]
				j--
			}
		}
	}

	var err error
	se := expr.AcquireMessageSelectExpr()
	for _, path := range fm.GetPaths() {
		var s scanner.Scanner
		s.Reset(path, p.errHandler)
		if err = p.parseSelectExprPath(&s, p.desc, se); err != nil {
			se.Free()
			return nil, err
		}
	}

	return se, nil
}

func (p *Parser) parseSelectExprPath(s *scanner.Scanner, md protoreflect.MessageDescriptor, se *expr.MessageSelectExpr) error {
	pos, tok, lit := s.Scan()
	if tok == token.ASTERISK {
		// This means it is a wildcard path.
		// The next one must be the end of the path.
		pos, tok, lit = s.Scan()
		if tok != token.EOF {
			if p.errHandler != nil {
				p.errHandler(pos, "a wildcard path must be the last segment of the path")
			}
			return ErrInvalidField
		}

		mi := p.msgInfo.MessageInfo(md)
		for _, f := range mi.Fields {
			if f.InputOnly {
				continue
			}

			fs := getFieldSelectorExpr(se, f.Desc.Name())
			if fs != nil {
				continue
			}

			fs = expr.AcquireFieldSelectorExpr()
			fs.Message = md.FullName()
			fs.Field = f.Desc.Name()
			fs.FieldComplexity = f.Complexity
			se.Fields = append(se.Fields, fs)

			if f.Desc.Kind() == protoreflect.MessageKind && !f.Desc.IsMap() {
				// This means it wants to get all sub fields for given message.
				p.selectAllMsgFields(f, fs)
			}
		}
		return nil
	}

	if !tok.IsIdent() {
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("expected field name but got %q", lit))
		}
		return ErrInvalidSyntax
	}

	fi, ok := p.msgInfo.MessageInfo(md).FieldByName(protoreflect.Name(lit))
	if !ok {
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("field %q not found", lit))
		}
		return ErrInvalidField
	}

	// If the field is input only it cannot be selected.
	if fi.InputOnly {
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("field %q is marked as input only", lit))
		}
		return ErrInvalidField
	}

	// Try to get the field selector expression if it already exists.
	fs := getFieldSelectorExpr(se, fi.Desc.Name())
	if fs == nil {
		fs = expr.AcquireFieldSelectorExpr()
		fs.Message = md.FullName()
		fs.Field = fi.Desc.Name()
		fs.FieldComplexity = fi.Complexity
		se.Fields = append(se.Fields, fs)
	}

	pos, tok, lit = s.Scan()
	if tok == token.EOF {
		if fi.Desc.Kind() == protoreflect.MessageKind && !fi.Desc.IsMap() {
			// This means it wants to get all sub fields for given message.
			p.selectAllMsgFields(fi, fs)
			return nil
		}
		return nil
	}

	// If token is not EOF, than the next token must be a period.
	if tok != token.PERIOD {
		if p.errHandler != nil {
			p.errHandler(pos, fmt.Sprintf("expected '.' but got %q", lit))
		}
		return ErrInvalidSyntax
	}

	// Check if the field is a traversal field,
	// a message, map or repeated field.
	if fi.Desc.IsList() {
		// The only valid token is a wildcard now.
		pos, tok, lit = s.Scan()
		if tok != token.ASTERISK {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("field: %q is a repeated field, cannot traverse through it with non wildcard path", fi.Desc.Name()))
			}
			return ErrInvalidSyntax
		}

		pos, tok, lit = s.Scan()
		if tok == token.EOF {
			// This means it wants to get all sub fields for given message.
			if fi.Desc.Kind() == protoreflect.MessageKind {
				p.selectAllMsgFields(fi, fs)
			}
			return nil
		}

		if tok != token.PERIOD {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected '.' but got %q", lit))
			}
			return ErrInvalidSyntax
		}

		if fi.Desc.Kind() != protoreflect.MessageKind {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("field: %q is not a message, cannot traverse through it", fi.Desc.Name()))
			}
			return ErrInvalidSyntax
		}

		sub := getMessageFieldSelectExpr(se, fi.Desc.Name())
		if sub == nil {
			// This is a wildcard path selector.

			sub = expr.AcquireMessageSelectExpr()
			sub.Message = md.FullName()
			fs.Traversal = sub
		}
		return p.parseSelectExprPath(s, fi.Desc.Message(), sub)
	}

	if fi.Desc.IsMap() {
		// If the expression is a map, the next token may be either a key or wildcard.
		var mk *expr.MapSelectKeysExpr
		mk, ok = fs.Traversal.(*expr.MapSelectKeysExpr)
		if !ok {
			mk = expr.AcquireMapSelectKeysExpr()
			fs.Traversal = mk
		}
		pos, tok, lit = s.Scan()

		var found *expr.MapKeyExpr
		if tok == token.ASTERISK {
			for _, key := range mk.Keys {
				_, ok = key.Key.(*expr.WildcardExpr)
				if ok {
					found = key
					break
				}
			}

			if found == nil {
				found = expr.AcquireMapKeyExpr()
				found.Key = expr.AcquireWildcardExpr()
				mk.Keys = append(mk.Keys, found)
			}
		} else if tok.IsLiteral() || tok.IsKeyword() {
			// Extract the map key value.
			var value any
			switch fi.Desc.MapKey().Kind() {
			case protoreflect.BoolKind:
				if !tok.IsBoolean() {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected boolean but got %q", lit))
					}
					return ErrInvalidSyntax
				}

				value = tok == token.TRUE
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
				protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				if !tok.IsInteger() {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected integer but got %q", lit))
					}
					return ErrInvalidSyntax
				}

				var err error
				value, err = strconv.ParseInt(lit, 10, 64)
				if err != nil {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected integer but got %q", lit))
					}
					return ErrInvalidSyntax
				}

			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
				protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				if !tok.IsInteger() {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected integer but got %q", lit))
					}
					return ErrInvalidSyntax
				}

				var err error
				value, err = strconv.ParseUint(lit, 10, 64)
				if err != nil {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected integer but got %q", lit))
					}
					return ErrInvalidSyntax
				}

			case protoreflect.StringKind:
				if !(tok == token.STRING || tok.IsIdent()) {
					if p.errHandler != nil {
						p.errHandler(pos, fmt.Sprintf("expected string but got %q", lit))
					}
					return ErrInvalidSyntax
				}

				value = lit
			}
			for _, key := range mk.Keys {
				var ve *expr.ValueExpr
				ve, ok = key.Key.(*expr.ValueExpr)
				if ok && ve.Value == value {
					found = key
					break
				}
			}

			if found == nil {
				found = expr.AcquireMapKeyExpr()
				ve := expr.AcquireValueExpr()
				ve.Value = value
				found.Key = ve
				mk.Keys = append(mk.Keys, found)
			}
		} else {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected wildcard or literal but got %q", lit))
			}
			return ErrInvalidSyntax
		}

		mv := fi.Desc.MapValue()
		pos, tok, lit = s.Scan()
		if tok == token.EOF {
			// Check if the map value is a message, and if it is, map all fields.
			if mv.Kind() == protoreflect.MessageKind {
				mi := p.msgInfo.MessageInfo(mv.Message())

				sub, hadMessageSelect := found.Traversal.(*expr.MessageSelectExpr)
				if !hadMessageSelect {
					sub = expr.AcquireMessageSelectExpr()
					sub.Message = mv.Message().FullName()
					found.Traversal = sub
				}
				paths := len(sub.Fields)
				for _, f := range mi.Fields {
					if f.InputOnly {
						continue
					}

					var elemFS *expr.FieldSelectorExpr
					if hadMessageSelect {
						for i := 0; i < paths; i++ {
							if sub.Fields[i].Field == f.Desc.Name() {
								elemFS = sub.Fields[i]
								break
							}
						}
					} else {
						elemFS = expr.AcquireFieldSelectorExpr()
						elemFS.Message = md.FullName()
						elemFS.Field = f.Desc.Name()
						elemFS.FieldComplexity = f.Complexity
					}

					sub.Fields = append(sub.Fields, elemFS)
				}
				return nil
			}
			return nil
		}

		if mv.Kind() != protoreflect.MessageKind {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("field: %q map value is not a message, cannot traverse through it", fi.Desc.Name()))
			}
			return ErrInvalidSyntax
		}

		if tok != token.PERIOD {
			if p.errHandler != nil {
				p.errHandler(pos, fmt.Sprintf("expected '.' but got %q", lit))
			}
			return ErrInvalidSyntax
		}

		var sub *expr.MessageSelectExpr
		sub, ok = found.Traversal.(*expr.MessageSelectExpr)
		if !ok {
			sub = expr.AcquireMessageSelectExpr()
			sub.Message = mv.Message().FullName()
			found.Traversal = sub
		}

		return p.parseSelectExprPath(s, mv.Message(), sub)
	}

	sub, ok := fs.Traversal.(*expr.MessageSelectExpr)
	if !ok {
		sub = expr.AcquireMessageSelectExpr()
		sub.Message = md.FullName()
		fs.Traversal = sub
	}
	return p.parseSelectExprPath(s, fi.Desc.Message(), sub)
}

func (p *Parser) selectAllMsgFields(fi info.FieldInfo, fs *expr.FieldSelectorExpr) {
	mi := p.msgInfo.MessageInfo(fi.Desc.Message())
	se, ok := fs.Traversal.(*expr.MessageSelectExpr)
	if ok {
		pathsLen := len(se.Fields)
		for _, f := range mi.Fields {
			if f.InputOnly {
				continue
			}

			var found bool
			for i := 0; i < pathsLen; i++ {
				if se.Fields[i].Field == f.Desc.Name() {
					found = true
					break
				}
			}

			if found {
				continue
			}

			fse := expr.AcquireFieldSelectorExpr()
			fse.Message = mi.Desc.FullName()
			fse.Field = f.Desc.Name()
			fse.FieldComplexity = f.Complexity
			se.Fields = append(se.Fields, fse)
		}
		return
	}

	se = expr.AcquireMessageSelectExpr()
	se.Message = mi.Desc.FullName()
	fs.Traversal = se

	for _, f := range mi.Fields {
		if f.InputOnly {
			continue
		}

		fse := expr.AcquireFieldSelectorExpr()
		fse.Message = mi.Desc.FullName()
		fse.Field = f.Desc.Name()
		se.Fields = append(se.Fields, fse)
	}
}

func getMessageFieldSelectExpr(src *expr.MessageSelectExpr, fieldName protoreflect.Name) *expr.MessageSelectExpr {
	for _, path := range src.Fields {
		if path.Field == fieldName {
			f, ok := path.Traversal.(*expr.MessageSelectExpr)
			if !ok {
				return nil
			}
			return f
		}
	}
	return nil
}

func getFieldSelectorExpr(src *expr.MessageSelectExpr, fieldName protoreflect.Name) *expr.FieldSelectorExpr {
	for _, path := range src.Fields {
		if path.Field == fieldName {
			return path
		}
	}
	return nil
}
