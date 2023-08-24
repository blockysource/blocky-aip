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

package pagination

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"strings"

	"github.com/blockysource/blocky-aip/expr"
)

// TokenizeStruct encodes input value as a gob.Encoded string.
// The input should be a non exported golang struct, that contains
// exported fields, that could be decoded with a valid Next filter expression.
// The output string is a url safe base64 encoded gob of this struct.
func TokenizeStruct(in any) (string, error) {
	var buf strings.Builder

	e := base64.NewEncoder(base64.URLEncoding, &buf)
	w := gzip.NewWriter(e)
	if err := gob.NewEncoder(w).Encode(in); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	if err := e.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// DecodeToken decodes a gob encoded string into an input struct.
func DecodeToken[T any](token string) (T, error) {
	var (
		sr  strings.Reader
		out T
	)
	sr.Reset(token)
	d := base64.NewDecoder(base64.URLEncoding, &sr)
	gr, err := gzip.NewReader(d)
	if err != nil {
		return out, err
	}
	if err := gob.NewDecoder(gr).Decode(&out); err != nil {
		return out, err
	}
	if err := gr.Close(); err != nil {
		return out, err
	}
	return out, nil
}

// NextTokenExpr is a pagination token filter expression.
// This can be used directly as a token struct, or embeddable in some internal packages for open-source
// implementations.
type NextTokenExpr struct {
	// Filter is a filter expression, that should be used to filter the next page.
	Filter expr.FilterExpr

	// OrderBy is an optional ordering expression, that should be used to order the next page.
	OrderBy *expr.OrderByExpr
}

// Free returns the allocated memory for both Filter and OrderBy fields.
func (n *NextTokenExpr) Free() {
	if n.Filter != nil {
		n.Filter.Free()
	}
	n.OrderBy.Free()
}
