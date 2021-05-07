# radixtree

[![GoDoc](https://pkg.go.dev/badge/github.com/gammazero/radixtree)](https://pkg.go.dev/github.com/gammazero/radixtree)
[![Build Status](https://github.com/gammazero/radixtree/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/radixtree/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/radixtree)](https://goreportcard.com/report/github.com/gammazero/radixtree)
[![codecov](https://codecov.io/gh/gammazero/radixtree/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/radixtree)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Package `radixtree` implements multiple forms of an Adaptive [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree), aka compressed [trie](https://en.wikipedia.org/wiki/Trie) or compact prefix tree.  This data structure is useful to quickly lookup data, using only the portion of the key that prefixes existing data.  It is also useful for finding items whose keys are a prefix of a search key (i.e. are found along the way when retrieving an item identified by a key), or when finding items whose keys are prefixed by the serach key (i.e. are found at or after a key).  When different values are stored using keys that have a common prefix, the common part of the key is only stored once.  Consider this when keys are similar to an OID, filepath, geohash, network address, etc.

This radix tree is adaptive in the sense that nodes are not constant size, having as few or many children as needed, up to the number of different key segments to traverse to the next branch or value.

An iterator for each type of radix tree allows a tree to be traversed one key segment at a time.  This is useful for incremental lookups of partial keys.

The implementations are optimized for Get performance and allocates 0 bytes of heap memory per Get; therefore no garbage to collect.  Once the radix tree is built, it can be repeatedly searched quickly. Concurrent searches are safe since these do not modify the radixtree. Access is not synchronized (not concurrent safe with writes), allowing the caller to synchronize, if needed, in whatever manner works best for the application.

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
    rt := new(radixtree.Runes)
    rt.Put("tomato", "TOMATO")
    rt.Put("tom", "TOM")
    rt.Put("tommy", "TOMMY")
    rt.Put("tornado", "TORNADO")

    val := rt.Get("tom")
    fmt.Println("Found", val)      // output: Found TOM

    // Find all items whose keys start with "tom"
    var vals []interface{}
    rt.Walk("tom", func(key fmt.Stringer, value interface{}) error {
        vals = append(vals, value)
        return nil
    })
    fmt.Println(vals)              // output: [TOM, TOMATO, TOMMY]

    // Find all items whose keys are a prefix of "tomato"
    vals = vals[0:0]
    rt.WalkPath("tomato", func(key string, value interface{}) error {
        vals = append(vals, value)
        return nil
    })
    fmt.Println(vals)              // output: [TOM, TOMATO]

    // Find each item whose key is a prefix of "tomato", using iterator
    iter := rt.NewIterator()
    for _, r := range "tomato" {
        if !iter.Next(r) {
            break
        }
        if val := iter.Value(); val != nil {
            fmt.Println(val)       // output: TOM
        }                          // output: TOMATO
    }

    if rt.Delete("tom") {
        fmt.Println("Deleted tom") // output: Deleted tom
    }
    val = rt.Get("tom")
    fmt.Println("Found", val)      // output: Found <nil>
}
```

## License

[MIT License](LICENSE)

