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
	"errors"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/blockysource/blocky-aip/internal/info"
	"github.com/blockysource/blocky-aip/scanner"
)

var (
	// ErrInvalidField is an error that is returned when a field is invalid or is not found.
	ErrInvalidField = errors.New("invalid field")

	// ErrInternalError is an internal error done during interpretation.
	ErrInternalError = errors.New("internal error")

	// ErrInvalidSyntax is an error returned by the parser when the field mask
	// has invalid syntax.
	ErrInvalidSyntax = errors.New("invalid syntax")
)

// Parser is a field mask to expression parser.
// It allows to match the google.protobuf.FieldMask with a message descriptor fields.
// As a result it produces an expression that can be used to patch the message.
type Parser struct {
	desc       protoreflect.MessageDescriptor
	errHandler scanner.ErrorHandler

	ignoreNonUpdatable bool
	msgInfo            info.MessagesInfo
}

// OptionFn is an option function for the Parser.
type OptionFn func(p *Parser) error

// ErrHandlerOption is an option function that sets the error handler for the parser.
func ErrHandlerOption(fn scanner.ErrorHandler) OptionFn {
	return func(p *Parser) error {
		p.errHandler = fn
		return nil
	}
}

// IgnoreNonUpdatableOption is an option function that sets the ignore non updatable option for the parser.
func IgnoreNonUpdatableOption(p *Parser) error {
	p.ignoreNonUpdatable = true
	return nil
}

// Reset the parser.
func (p *Parser) Reset(msg proto.Message, opts ...OptionFn) error {
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return err
		}
	}
	p.desc = msg.ProtoReflect().Descriptor()

	p.msgInfo = info.MapMsgInfo(p.desc)
	return nil
}

func (p *Parser) extractTimeValue(msg protoreflect.Message) time.Time {
	// The timestamppb.Timestamp contains two fields: seconds and nanos.
	var seconds, nanos int64
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch fd.Name() {
		case "seconds":
			seconds = v.Int()
		case "nanos":
			nanos = v.Int()
		}
		return true
	})
	return time.Unix(seconds, nanos)
}

func (p *Parser) extractDurationValue(msg protoreflect.Message) time.Duration {
	// The durationpb.Duration contains two fields: seconds and nanos.
	var seconds, nanos int64
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch fd.Name() {
		case "seconds":
			seconds = v.Int()
		case "nanos":
			nanos = v.Int()
		}
		return true
	})
	return time.Duration(seconds)*time.Second + time.Duration(nanos)*time.Nanosecond
}
