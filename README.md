# radixtree

[![GoDoc](https://pkg.go.dev/badge/github.com/gammazero/radixtree)](https://pkg.go.dev/github.com/gammazero/radixtree)
[![Build Status](https://github.com/gammazero/radixtree/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/radixtree/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/radixtree)](https://goreportcard.com/report/github.com/gammazero/radixtree)
[![codecov](https://codecov.io/gh/gammazero/radixtree/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/radixtree)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Package `radixtree` implements multiple forms of an Adaptive [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree), aka compressed [trie](https://en.wikipedia.org/wiki/Trie) or compact prefix tree.  This data structure is useful to quickly lookup data by key, find find data whose keys have a common prefix, or find data whose keys are a prefix (i.e. found along the way) of a search key.

The implementations are optimized for Get performance and allocates 0 bytes of heap memory per Get or Walk; therefore no garbage to collect.  Once the radix tree is built, it can be repeatedly searched quickly. Concurrent searches are safe since these do not modify the radixtree. Access is not synchronized (not concurrent safe with writes), allowing the caller to synchronize, if needed, in whatever manner works best for the application.

This radix tree offers the following features:

- Multiple types of radix tree: Bytes, Paths, Runes
- Efficient: Operations for all types of radix tree are O(k).  Zero memory allocation for all read operations.
- Compact: When values are stored using keys that have a common prefix, the common part of the key is only stored once.  Consider this when keys are similar to an OID, filepath, geohash, network address, etc. Nodes that do not branch or contain values are compressed out of the tree.
- Adaptive: This radix tree is adaptive in the sense that nodes are not constant size, having only as many children that are needed, from zero to the maximum possible number of different key segments.
- Iterators: An iterator for each type of radix tree allows a tree to be traversed one key segment at a time.  This is useful for incremental lookup.  Iterators can be copied in order to branch a search, and iterate the copies concurrently.
- Ordered iteration: Walking and iterating the tree is done in lexical order, making the output deterministic.

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
    rt := new(radixtree.Bytes)
    rt.Put("tomato", "TOMATO")
    rt.Put("tom", "TOM")
    rt.Put("tommy", "TOMMY")
    rt.Put("tornado", "TORNADO")

    val := rt.Get("tom")
    fmt.Println("Found", val)
    // Output: Found TOM

    // Find all items whose keys start with "tom"
    rt.Walk("tom", func(key fmt.Stringer, value interface{}) error {
        fmt.Println(value)
        return nil
    })
    // Output:
    // TOM
    // TOMATO
    // TOMMY

    // Find all items whose keys are a prefix of "tomato"
    rt.WalkPath("tomato", func(key string, value interface{}) error {
        fmt.Println(value)
        return nil
    })
    // Output:
    // TOM
    // TOMATO

    // Find each item whose key is a prefix of "tomato", using iterator
    iter := rt.NewIterator()
    for _, r := range "tomato" {
        if !iter.Next(r) {
            break
        }
        if val := iter.Value(); val != nil {
            fmt.Println(val)
        }
    }
    // Output:
    // TOM
    // TOMATO

    if rt.Delete("tom") {
        fmt.Println("Deleted tom")
    }
    // Output: Deleted tom

    val = rt.Get("tom")
    fmt.Println("Found", val)
    // Output: Found <nil>
}
```

## License

[MIT License](LICENSE)

