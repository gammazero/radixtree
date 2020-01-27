# radixtree

[![Build Status](https://travis-ci.com/gammazero/radixtree.svg)](https://travis-ci.com/gammazero/radixtree)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/radixtree)](https://goreportcard.com/report/github.com/gammazero/radixtree)
[![codecov](https://codecov.io/gh/gammazero/radixtree/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/radixtree)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Package `radixtree` implements multiple forms of an Adaptive [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree), aka compressed [trie](https://en.wikipedia.org/wiki/Trie) or compact prefix tree.  It is adaptive in the sense that nodes are not constant size, having as few or many children as needed, up to the number of different key segments to traverse to the next branch or value.

The implementations are optimized for Get performance and allocate 0 bytes of heap memory per Get; therefore no garbage to collect.  Once the radix tree is build, it can be repeatedly searched very quickly.

Access is not synchronized (not concurrent safe), allowing the caller to synchronize, if needed, in whatever manner works best for the application.

[![GoDoc](https://godoc.org/github.com/gammazero/radixtree?status.svg)](https://godoc.org/github.com/gammazero/radixtree)

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
	rt.Put("tornado", "TORNADO")

	val := rt.Get("tom")
	fmt.Println("Found", val)      // output: Found TOM

	vals, _ := rt.GetPath("tomato")
	fmt.Println("Values:", vals)   // output: [TOM, TOMATO]

    iter := rt.NewIterator()
    for _, r := range "tomato" {
        if !iter.Next(r) {
            break
        }
        val := iter.Value()
        if val == nil {
            continue
        }
        fmt.Println(val)           // output: tom
    }                              // output: tomato
       
    if rt.Delete("tom") {
        fmt.Println("Deleted tom") // output: Deleted tom
    }
	val = rt.Get("tom")
	fmt.Println("Found", val)      // output: Found <nil>

}
```

## License

[MIT License](LICENSE)

