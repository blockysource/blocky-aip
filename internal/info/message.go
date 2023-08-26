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

package info

import (
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	blockyannotations "github.com/blockysource/go-genproto/blocky/api/annotations"
)

// MessageInfo is a struct that contains information about a message.
type MessageInfo struct {
	Desc   protoreflect.MessageDescriptor
	Fields []FieldInfo
}

// FieldInfo is a struct that contains information about a field.
type FieldInfo struct {
	// Desc is the field descriptor.
	Desc protoreflect.FieldDescriptor

	// Complexity is the complexity of the field.
	Complexity int64

	// Forbidden is true if the field filtering is forbidden.
	Forbidden bool

	// Nullable is true if the field is nullable.
	Nullable bool

	// InputOnly is true if the field is input only.
	InputOnly bool

	// OutputOnly is true if the field is output only.
	OutputOnly bool

	// IsOneOf is true if the field is a oneof field.
	IsOneOf bool

	// Required is true if the field is required.
	Required bool

	// Immutable is true if the field is immutable.
	Immutable bool
}

// MapMsgInfo maps a message descriptor to a MessageInfo struct.
func MapMsgInfo(desc protoreflect.MessageDescriptor) MessagesInfo {
	var b mapper
	b.mapMessage(desc)
	return b.msgInfo
}

type mapper struct {
	msgInfo []*MessageInfo
}

// MessagesInfo is a slice of MessageInfo.
type MessagesInfo []*MessageInfo

// GetFieldInfo returns the field info for the given field descriptor.
func (mi MessagesInfo) GetFieldInfo(fd protoreflect.FieldDescriptor) FieldInfo {
	for _, m := range mi {
		if m.Desc == fd.Parent() {
			for _, f := range m.Fields {
				if f.Desc == fd {
					return f
				}
			}
		}
	}
	panic("field not found")
}

func (b *mapper) mapMessage(msg protoreflect.MessageDescriptor) {
	for i := 0; i < len(b.msgInfo); i++ {
		if b.msgInfo[i].Desc.FullName() == msg.FullName() {
			return
		}
	}

	fields := msg.Fields()

	mi := &MessageInfo{Desc: msg}
	b.msgInfo = append(b.msgInfo, mi)

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		fi := FieldInfo{
			Desc:       fd,
			Complexity: getFieldComplexity(fd),
			Forbidden:  IsFieldFilteringForbidden(fd),
			Nullable:   IsFieldNullable(fd),
		}

		fb, ok := proto.GetExtension(fd.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
		if ok {
			for _, b := range fb {
				switch b {
				case annotations.FieldBehavior_INPUT_ONLY:
					fi.InputOnly = true
				case annotations.FieldBehavior_OUTPUT_ONLY:
					fi.OutputOnly = true
				case annotations.FieldBehavior_REQUIRED:
					fi.Required = true
				case annotations.FieldBehavior_IMMUTABLE:
					fi.Immutable = true
				}
			}
		}

		mi.Fields = append(mi.Fields, fi)

		if fd.Kind() == protoreflect.MessageKind {
			if fd.IsMap() {
				fd = fd.MapValue()
				if fd.Kind() != protoreflect.MessageKind {
					continue
				}
			}

			var found bool
			for k := 0; k < len(b.msgInfo); k++ {
				if b.msgInfo[k].Desc.FullName() == fd.Message().FullName() {
					found = true
					break
				}
			}

			if !found {
				b.mapMessage(fd.Message())
			}

		}
	}

	for i := 0; i < msg.Oneofs().Len(); i++ {
		o := msg.Oneofs().Get(i)
		of := o.Fields()
		for j := 0; j < of.Len(); j++ {
			fd := of.Get(j)
			fi := FieldInfo{
				Desc:       fd,
				IsOneOf:    true,
				Complexity: getFieldComplexity(fd),
				Forbidden:  IsFieldFilteringForbidden(fd),
				Nullable:   IsFieldNullable(fd),
			}

			fb, ok := proto.GetExtension(fd.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
			if ok {
				for _, b := range fb {
					switch b {
					case annotations.FieldBehavior_INPUT_ONLY:
						fi.InputOnly = true
					case annotations.FieldBehavior_OUTPUT_ONLY:
						fi.OutputOnly = true
					case annotations.FieldBehavior_REQUIRED:
						fi.Required = true
					case annotations.FieldBehavior_IMMUTABLE:
						fi.Immutable = true
					}
				}
			}

			mi.Fields = append(mi.Fields, fi)

			if fd.Kind() == protoreflect.MessageKind {
				if fd.IsMap() {
					fd = fd.MapValue()
					if fd.Kind() != protoreflect.MessageKind {
						continue
					}
				}

				fm := fd.Message()
				var found bool
				for k := 0; k < len(b.msgInfo); k++ {
					if b.msgInfo[k].Desc.FullName() == fm.FullName() {
						found = true
						break
					}
				}
				if !found {
					b.mapMessage(fm)
				}
			}
		}
	}
}

func getFieldComplexity(fdt protoreflect.FieldDescriptor) int64 {
	c, ok := proto.GetExtension(fdt.Options(), blockyannotations.E_Complexity).(int64)
	if !ok || c == 0 {
		return 1
	}
	return c
}

// IsFieldFilteringForbidden returns true if the field filtering is forbidden.
func IsFieldFilteringForbidden(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), blockyannotations.E_QueryOpt).([]blockyannotations.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == blockyannotations.FieldQueryOption_FORBID_FILTERING {
			return true
		}
	}
	return false
}

// IsFieldNullable checks if the input field is nullable.
func IsFieldNullable(field protoreflect.FieldDescriptor) bool {
	// At first try blockaypi.E_Nullable extension, if not found, then try google api.OPTIONAL extension.
	// If not found, then return false.
	queryOpts, ok := proto.GetExtension(field.Options(), blockyannotations.E_QueryOpt).([]blockyannotations.FieldQueryOption)
	if ok {
		for _, qo := range queryOpts {
			if qo == blockyannotations.FieldQueryOption_NULLABLE {
				return true
			}
		}
	}

	fb, ok := proto.GetExtension(field.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
	if !ok {
		return false
	}

	if field.Kind() == protoreflect.MessageKind {
		for _, b := range fb {
			if b == annotations.FieldBehavior_REQUIRED {
				return false
			}
		}
		return true
	}
	for _, b := range fb {
		switch b {
		case annotations.FieldBehavior_REQUIRED, annotations.FieldBehavior_IMMUTABLE:
			return false
		case annotations.FieldBehavior_OPTIONAL:
			return true
		}
	}
	return false
}
