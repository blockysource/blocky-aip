// Copyright (C) 2024 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/blockysource/go-genproto/blocky/api/annotations"
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

	// FilteringForbidden is true if the field filtering is forbidden.
	FilteringForbidden bool

	// OrderingForbidden is true if the field ordering is forbidden.
	OrderingForbidden bool

	// NonTraversal is true if the field is non-traversal.
	NonTraversal bool

	// NoTextSearch is true if the field is no text search.
	NoTextSearch bool

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

	// NonEmptyDefault is true if the field has a non-empty default value.
	NonEmptyDefault bool

	// IsTimestamp is true if the field is a timestamp.
	IsTimestamp bool

	// IsDuration is true if the field is a duration.
	IsDuration bool

	// IsStructpb is true if the field is a structpb.
	IsStructpb bool
}

// Undefined returns true if the descriptor is nil.
func (fi FieldInfo) Undefined() bool {
	return fi.Desc == nil
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

// MessageInfo returns the message info for the given message descriptor.
func (mi MessagesInfo) MessageInfo(md protoreflect.MessageDescriptor) *MessageInfo {
	for _, m := range mi {
		if m.Desc == md {
			return m
		}
	}
	panic("message not found")
}

// FieldByName returns the field info for the given field name.
func (mi *MessageInfo) FieldByName(name protoreflect.Name) (FieldInfo, bool) {
	for _, f := range mi.Fields {
		if f.Desc.Name() == name {
			return f, true
		}
	}
	return FieldInfo{}, false
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
			Desc:               fd,
			Complexity:         getFieldComplexity(fd),
			FilteringForbidden: isFieldFilteringForbidden(fd),
			OrderingForbidden:  isFieldOrderingForbidden(fd),
			Nullable:           isFieldOptional(fd),
			NonTraversal:       isFieldNonTraversal(fd),
			NoTextSearch:       isFieldNoTextSearch(fd),
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
				case annotations.FieldBehavior_NON_EMPTY_DEFAULT:
					fi.NonEmptyDefault = true
				}
			}
		}

		if fd.Kind() == protoreflect.MessageKind {
			switch fd.Message().FullName() {
			case "google.protobuf.Timestamp":
				fi.IsTimestamp = true
			case "google.protobuf.Duration":
				fi.IsDuration = true
			case "google.protobuf.Struct":
				fi.IsStructpb = true
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
				Desc:               fd,
				IsOneOf:            true,
				Complexity:         getFieldComplexity(fd),
				FilteringForbidden: isFieldFilteringForbidden(fd),
				Nullable:           isFieldOptional(fd),
				NonTraversal:       isFieldNonTraversal(fd),
				NoTextSearch:       isFieldNoTextSearch(fd),
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
	c, ok := proto.GetExtension(fdt.Options(), annotationspb.E_Complexity).(int64)
	if !ok || c == 0 {
		return 1
	}
	return c
}

// isFieldFilteringForbidden returns true if the field filtering is forbidden.
func isFieldFilteringForbidden(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), annotationspb.E_QueryOpt).([]annotationspb.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == annotationspb.FieldQueryOption_FORBID_FILTERING {
			return true
		}
	}
	return false
}

// isFieldOrderingForbidden returns true if the field filtering is forbidden.
func isFieldOrderingForbidden(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), annotationspb.E_QueryOpt).([]annotationspb.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == annotationspb.FieldQueryOption_FORBID_SORTING {
			return true
		}
	}
	return false
}

// isFieldNonTraversal returns true if the field is non-traversal.
func isFieldNonTraversal(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), annotationspb.E_QueryOpt).([]annotationspb.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == annotationspb.FieldQueryOption_NON_TRAVERSAL {
			return true
		}
	}
	return false
}

// isFieldNoTextSearch returns true if the field is no text search.
func isFieldNoTextSearch(field protoreflect.FieldDescriptor) bool {
	opts, ok := proto.GetExtension(field.Options(), annotationspb.E_QueryOpt).([]annotationspb.FieldQueryOption)
	if !ok {
		return false
	}
	for _, opt := range opts {
		if opt == annotationspb.FieldQueryOption_NO_TEXT_SEARCH {
			return true
		}
	}
	return false
}

// isFieldOptional checks if the input field is nullable.
func isFieldOptional(field protoreflect.FieldDescriptor) bool {
	fb, ok := proto.GetExtension(field.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior)
	if !ok {
		return false
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
