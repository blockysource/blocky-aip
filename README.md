# The Golang Blocky AIP implementation

[![GoDoc](https://godoc.org/github.com/blockysource/blockysql?status.svg)](https://godoc.org/github.com/blockysource/blocky-aip)
[![Go Report Card](https://goreportcard.com/badge/github.com/blockysource/blocky-aip)](https://goreportcard.com/report/github.com/blockysource/blocky-aip)
[![Build Status](https://travis-ci.org/blockysource/blocky-aip.svg?branch=master)](https://travis-ci.org/blockysource/blocky-aip)
[![codecov](https://codecov.io/gh/blockysource/blocky-aip/branch/master/graph/badge.svg)](https://codecov.io/gh/blockysource/blocky-aip)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The Golang Blocky AIP is a library that provides [API Improvement Proposals](https://google.aip.dev/).
implementation for gRPC services.

It provides a simple set of tools to parse and validate the request and response messages.

## Filtering

The library provides a way to filter the request and response messages based on the
[Google AIP-160](https://google.aip.dev/160) specification.

### Abstract Syntax Tree (AST)

The library provides a way to parse the filter string into an AST (specified in
the [EBNF Grammar](https://google.aip.dev/assets/misc/ebnf-filtering.txt))

- `github.com/blockysource/blocky-aip/filtering/ast` defines AST nodes, which can be used to traverse the AST.
- `github.com/blockysource/blocky-aip/parser` - parses an input filtering string into a valid AST expression.

As an extension to the EBNF Grammar, parser has an ability to parse structures and arrays.

A structure is defined as:

- name - period separated list of structure name and field name, e.g. `foo.bar.baz`
- fields: a list of field names, e.g. `foo.bar.baz` has fields `foo`, `bar`, `baz`
- each field can be any valid comparable expression i.e.: member, function, struct or array.
- a structure fields definition is opened by the `{` and closed by the `}`

example: `foo.bar.baz{foo: 1, bar: "bar", baz: func(), qux: [1, 2, 3]}`

An array is defined as:

- opening bracket `[`
- list of expressions separated by comma `,`
- closing bracket `]`

### Proto Filtering

The library provides a way to filter the request and response messages based on the AST.

It maps the field names from provided `protoreflect.MessageDescriptor`, parses the filter string into AST
and converts it into some simple form of `expr.FilterExpr`.



