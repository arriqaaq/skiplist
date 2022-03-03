# skiplist
A Skiplist implementation in Go

Getting Started
===============

## Installing

To start using hash, install Go and run `go get`:

```sh
$ go get -u github.com/arriqaaq/skiplist
```

This will retrieve the library.

## Usage

```go
package main

import "github.com/arriqaaq/skiplist"

type kv struct{k,v string}

func main() {

    // Set (accepts any value)
    n := skiplist.New()
	val := []byte("test_val")

	n.Set("ec", val)
	n.Set("dc", 123)
	n.Set("ac", val)

    // Get
    node:= list.Get("ec")
	assert.Equal(t, val, node.Value())

    // Delete
	n.Delete("dc")
}
```

## Supported Commands

```go
Supported  commands

Set
Get
Delete
```
