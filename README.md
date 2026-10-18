# radixtree

[![GoDoc](https://pkg.go.dev/badge/github.com/gammazero/radixtree)](https://pkg.go.dev/github.com/gammazero/radixtree)
[![Build Status](https://github.com/gammazero/radixtree/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/radixtree/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gammazero/radixtree/graph/badge.svg?token=RT5GYXGZ9T)](https://codecov.io/gh/gammazero/radixtree)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Package `radixtree` implements a [**radix tree**](https://en.wikipedia.org/wiki/Radix_tree), also known as a **compressed** [**trie**](https://en.wikipedia.org/wiki/Trie) or **compact prefix tree**. Use it to look up values by key, find values whose keys share a common prefix, or find values whose keys lie along the path to a given key.

The tree uses a radix-256 structure where each key symbol is a byte, giving up to 256 branches per node. Nodes hold only as many children as needed, keeping memory proportional to the data stored.

Read operations (`Get`, `Iter`, `IterPrefix`, `IterPath`, and `Stepper` methods) allocate no heap memory and are safe to call concurrently. Write operations are not synchronized; callers that mix reads and writes must coordinate access themselves.

## Features

- **Efficient**: All operations are O(key-length). Reads allocate no heap memory.
- **Ordered**: Iteration visits keys in lexical order, making output deterministic.
- **Nil-safe**: `Get` distinguishes between a missing key and a key whose value is `nil`.
- **Compact**: Keys with a common prefix share storage. Well-suited for timestamps, file paths, geohashes, and network addresses.
- **Iterators**: Range iterators cover all key-value pairs (`Iter`), pairs whose keys have a given prefix (`IterPrefix`), or pairs whose keys are a prefix of a given key (`IterPath`).
- **Inspectable**: `Nodes` iterates over every node of the tree, exposing its structure for debugging and statistics.
- **Stepper**: Walk the tree one byte at a time for incremental lookup. Copy a `Stepper` to branch a search and use the copies concurrently.
- **Generics**: Store any value type without interface conversions.

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
    rt := radixtree.New[string]()
    rt.Put("tomato", "TOMATO")
    rt.Put("tom", "TOM")
    rt.Put("tommy", "TOMMY")
    rt.Put("tornado", "TORNADO")

    if val, found := rt.Get("tom"); found {
        fmt.Println("Found", val)
    }

    // Find all items whose keys start with "tom".
    for key, value := range rt.IterPrefix("tom") {
        fmt.Println(key, "=", value)
    }

    // Find all items whose keys are a prefix of "tomato".
    for key, value := range rt.IterPath("tomato") {
        fmt.Println(key, "=", value)
    }

    if rt.Delete("tom") {
        fmt.Println("Deleted tom")
    }

    if _, found := rt.Get("tom"); !found {
        fmt.Println("tom not found")
    }
}
```

Output:

```
Found TOM
tom = TOM
tomato = TOMATO
tommy = TOMMY
tom = TOM
tomato = TOMATO
Deleted tom
tom not found
```

## License

[MIT License](LICENSE)
