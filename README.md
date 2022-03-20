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

import (
	"fmt"

	"github.com/arriqaaq/skiplist"
)

func main() {

	// Set (accepts any value)
	val := "test_val"

	n := skiplist.New()
	n.Set("ec", val)
	n.Set("dc", 123)
	n.Set("ac", val)

	// Get
	node := n.Get("ec")
	fmt.Println("value: ", node.Value())

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

