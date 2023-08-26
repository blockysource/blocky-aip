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

package expr

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrComposerEmptySelector = errors.New("empty selector")
)

// Composer is an expression composer that can be used to
// efficiently compose filter expressions.
// Each composed expression is allocated from the pool,
// and needs to be released back to the pool after use by
// calling Free method.
// If the expression is composed into a larger expression,
// then the larger expression is responsible for releasing
// the sub-expression, thus only a single call to Free is
// required.
type Composer struct {
	Desc protoreflect.MessageDescriptor
}

// Reset the composer.
func (c *Composer) Reset(md protoreflect.MessageDescriptor) {
	c.Desc = md
}

// Field parses the selector and returns a field selector expression.
func (c *Composer) Field(field string) (*FieldSelectorExpr, error) {
	if field == "" {
		return nil, ErrComposerEmptySelector

	}

	s, err := c.parseSelector(field)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// MustSelect parses the selector and returns a field selector expression.
// It panics if the selector is invalid.
func (c *Composer) MustSelect(selector string) *FieldSelectorExpr {
	s, err := c.Field(selector)
	if err != nil {
		panic(err)
	}
	return s
}

// And returns an AndExpr that can be used to compose a filter expression.
func (c *Composer) And(sub ...FilterExpr) *AndExpr {
	ae := AcquireAndExpr()
	ae.Expr = append(ae.Expr, sub...)
	return ae
}

// Or returns an OrExpr that can be used to compose a filter expression.
func (c *Composer) Or(sub ...FilterExpr) *OrExpr {
	oe := AcquireOrExpr()
	oe.Expr = append(oe.Expr, sub...)
	return oe
}

// Not returns a NotExpr that can be used to compose a filter expression.
func (c *Composer) Not(sub FilterExpr) *NotExpr {
	ne := AcquireNotExpr()
	ne.Expr = sub
	return ne
}

// Value returns a ValueExpr that can be used to compose a filter expression.
func (c *Composer) Value(v any) *ValueExpr {
	ve := AcquireValueExpr()
	ve.Value = v
	return ve
}

// Array returns an ArrayExpr that can be used to compose a filter expression.
func (c *Composer) Array(v ...FilterExpr) *ArrayExpr {
	ae := AcquireArrayExpr()
	ae.Elements = append(ae.Elements, v...)
	return ae
}

// Composite returns a CompositeExpr that can be used to compose a filter expression.
func (c *Composer) Composite(x FilterExpr) *CompositeExpr {
	ce := AcquireCompositeExpr()
	ce.Expr = x
	return ce
}

// Compare returns a CompareExpr that can be used to compose a filter expression.
func (c *Composer) Compare(left FilterExpr, cmp Comparator, right FilterExpr) *CompareExpr {
	ce := AcquireCompareExpr()
	ce.Left = left
	ce.Right = right
	ce.Comparator = cmp
	return ce
}

// FunctionCall returns a FunctionCallExpr that can be used to compose a filter expression.
func (c *Composer) FunctionCall(pkgName, name string, args ...FilterExpr) *FunctionCallExpr {
	fc := AcquireFunctionCallExpr()
	fc.PkgName = pkgName
	fc.Name = name
	fc.Arguments = append(fc.Arguments, args...)
	return fc
}

// MapKey returns a MapKeyExpr that can be used to compose a filter expression.
func (c *Composer) MapKey(key FilterExpr) *MapKeyExpr {
	mk := AcquireMapKeyExpr()
	mk.Key = key
	return mk
}

// MapValue returns a MapValueExpr that can be used to compose a filter expression.
func (c *Composer) MapValue(values ...MapValueExprEntry) *MapValueExpr {
	m := AcquireMapValueExpr()
	m.Values = append(m.Values, values...)
	return m
}

// OrderBy returns an OrderByExpr that can be used to compose a filter expression.
func (c *Composer) OrderBy(fields ...*OrderByFieldExpr) *OrderByExpr {
	oe := AcquireOrderByExpr()
	oe.Fields = append(oe.Fields, fields...)
	return oe
}

// OrderByField returns an OrderByFieldExpr that can be used to compose a filter expression.
func (c *Composer) OrderByField(field string, o Order) (*OrderByFieldExpr, error) {
	fd, err := c.parseSelector(field)
	if err != nil {
		return nil, err
	}

	oe := AcquireOrderByFieldExpr()
	oe.Field = fd
	oe.Order = o
	return oe, nil
}

// MustOrderByField returns an OrderByFieldExpr that can be used to compose a filter expression.
// It panics if the field is invalid.
func (c *Composer) MustOrderByField(field string, o Order) *OrderByFieldExpr {
	oe, err := c.OrderByField(field, o)
	if err != nil {
		panic(err)
	}
	return oe
}

// Pagination returns a PaginationExpr that can be used to compose a filter expression.
func (c *Composer) Pagination(pageSize, skip int32) *PaginationExpr {
	pe := AcquirePaginationExpr()
	pe.PageSize = pageSize
	pe.Skip = skip
	return pe
}

func (c *Composer) parseMapValueSelector(selector string) (protoreflect.FieldDescriptor, error) {
	var fd protoreflect.FieldDescriptor
	md := c.Desc

	split := strings.Split(selector, ".")
	for i, field := range split {
		fd = md.Fields().ByName(protoreflect.Name(field))
		if fd == nil {
			return nil, fmt.Errorf("selector: %s is not a valid field in the message: %s", field, md.FullName())
		}
		if fd.Kind() != protoreflect.MessageKind {
			return nil, fmt.Errorf("selector: %s cannot traverse through a non-message field: %s", field, fd.FullName())
		}
		if fd.Cardinality() == protoreflect.Repeated {
			return nil, fmt.Errorf("selector: %s cannot be based on a repeated field: %s", field, fd.FullName())
		}
		if i < len(split)-1 {
			md = fd.Message()
		}

		if i == len(split)-1 {
			if !fd.IsMap() {
				return nil, fmt.Errorf("selector: %s is not a map field: %s", field, fd.FullName())
			}
		}
	}
	return fd, nil
}

func (c *Composer) parseSelector(s string) (*FieldSelectorExpr, error) {
	var fd protoreflect.FieldDescriptor

	out := AcquireFieldSelectorExpr()

	out.Message = c.Desc.FullName()

	md := c.Desc

	split := strings.Split(s, ".")
	for i, field := range split {
		fd = md.Fields().ByName(protoreflect.Name(field))
		if fd == nil {
			var found bool
			for j := 0; j < md.Oneofs().Len(); j++ {
				od := md.Oneofs().Get(j)
				for j := 0; j < od.Fields().Len(); j++ {
					fd = od.Fields().Get(j)
					if fd != nil {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("selector: %s is not a valid field in the message: %s", field, c.Desc.FullName())
			}
		}
		out.Field = fd.Name()

		if i < len(split)-1 {
			if fd.Kind() != protoreflect.MessageKind {
				return nil, fmt.Errorf("selector: %s cannot traverse through a non-message field: %s", field, fd.FullName())
			}

			if fd.Cardinality() == protoreflect.Repeated {
				return nil, fmt.Errorf("selector: %s cannot be based on a repeated field: %s", field, fd.FullName())
			}
			next := AcquireFieldSelectorExpr()
			next.Message = fd.Message().FullName()
			md = fd.Message()
			out.Traversal = next
			out = next
		}
	}
	return out, nil
}
