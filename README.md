# radixtree

[![GoDoc](https://pkg.go.dev/badge/github.com/gammazero/radixtree)](https://pkg.go.dev/github.com/gammazero/radixtree)
[![Build Status](https://github.com/gammazero/radixtree/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/radixtree/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/radixtree)](https://goreportcard.com/report/github.com/gammazero/radixtree)
[![codecov](https://codecov.io/gh/gammazero/radixtree/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/radixtree)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Package `radixtree` implements an Adaptive [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree), aka compressed [trie](https://en.wikipedia.org/wiki/Trie) or compact prefix tree.  This data structure is useful to quickly lookup data by key, find values whose keys have a common prefix, or find values whose keys are a prefix (i.e. found along the way) of a search key.

It is adaptive in the sense that nodes are not constant size, having only as many children, up to the maximum, as needed to branch to all subtrees. This package implements a radix-256 tree where each key symbol (radix) is a byte, allowing up to 256 possible branches to traverse to the next node.

The implementation is optimized for Get performance and allocate 0 bytes of heap memory for any read operation (Get, Walk, WalkPath, etc.); therefore no garbage to collect.  Once a radix tree is built, it can be repeatedly searched quickly. Concurrent searches are safe since these do not modify the data structure. Access is not synchronized (not concurrent safe with writes), allowing the caller to synchronize, if needed, in whatever manner works best for the application.

This radix tree offers the following features:

- Efficient: Operations are O(k). Zero memory allocation for all read operations.
- Ordered iteration: Walking and iterating the tree is done in lexical order, making the output deterministic.
- Store `nil` values: Read operations differentiate between missing and `nil` values.
- Compact: When values are stored using keys that have a common prefix, the common part of the key is only stored once. Consider this when keys are similar to a timestamp, OID, filepath, geohash, network address, etc. Only the minimum number of nodes are kept to branch at the points where keys differ.
- Iterators: An `Iterator` type walks every key-value pair stored in the tree. A `Stepper` type of iterator traverses the tree one specified byte at a time, and is useful for incremental lookup. A Stepper can be copied in order to branch a search and iterate the copies concurrently.

## Install

```
$ go get github.com/gammazero/radixtree
```

## Example

```go
package main

import (
    "fmt"
    "github.com/gammazero/radixtree"
)

func main() {
    rt := radixtree.New()
    rt.Put("tomato", "TOMATO")
    rt.Put("tom", "TOM")
    rt.Put("tommy", "TOMMY")
    rt.Put("tornado", "TORNADO")

    val, found := rt.Get("tom")
    if found {
        fmt.Println("Found", val)
    }
    // Output: Found TOM

    // Find all items whose keys start with "tom"
    rt.Walk("tom", func(key string, value any) bool {
        fmt.Println(value)
        return false
    })
    // Output:
    // TOM
    // TOMATO
    // TOMMY

    // Find all items whose keys are a prefix of "tomato"
    rt.WalkPath("tomato", func(key string, value any) bool {
        fmt.Println(value)
        return false
    })
    // Output:
    // TOM
    // TOMATO

    if rt.Delete("tom") {
        fmt.Println("Deleted tom")
    }
    // Output: Deleted tom

    val, found = rt.Get("tom")
    if found {
        fmt.Println("Found", val)
    } else {
        fmt.Println("not found")
    }
    // Output: not found
}
```

## License

[MIT License](LICENSE)
