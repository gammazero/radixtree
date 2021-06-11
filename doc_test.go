package radixtree_test

import (
	"fmt"
	"strings"

	"github.com/gammazero/radixtree"
)

func ExampleRunes_Walk() {
	rt := new(radixtree.Bytes)
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items whose keys start with "tom"
	rt.Walk("tom", func(key string, value interface{}) bool {
		fmt.Println(value)
		return false
	})
}

func ExampleRunes_WalkPath() {
	rt := new(radixtree.Bytes)
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items that are a prefix of "tomato"
	rt.WalkPath("tomato", func(key string, value interface{}) bool {
		fmt.Println(value)
		return false
	})
	// Output:
	// TOM
	// TOMATO
}

func ExampleRunes_NewIterator() {
	rt := new(radixtree.Runes)
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	iter := rt.NewIterator()
	for _, r := range "tomato" {
		if !iter.Next(r) {
			break
		}
		if val, ok := iter.Value(); ok {
			fmt.Println(val)
		}
	}
	// Output:
	// TOM
	// TOMATO
}

func ExamplePaths_Walk() {
	pt := new(radixtree.Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")

	// Find all items with keys that start with "home/abc/Documents"
	pt.Walk("home/abc/Documents", func(key string, value interface{}) bool {
		fmt.Println(value)
		return false
	})
}

func ExamplePaths_WalkPath() {
	pt := new(radixtree.Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")

	// Find item in each path segment
	pt.WalkPath("home/abc/Documents/pic.png", func(key string, value interface{}) bool {
		fmt.Println(key, "=>", value)
		return false
	})
	// Output:
	// home/abc => my home directory
	// home/abc/Documents => my documents
	// home/abc/Documents/pic.png => cat pic
}

func ExamplePaths_NewIterator() {
	pt := new(radixtree.Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")

	// Find item in each path segment
	parts := strings.Split("home/abc/Documents/pic.png", "/")
	iter := pt.NewIterator()
	for i, p := range parts {
		if !iter.Next(p) {
			break
		}
		if value, ok := iter.Value(); ok {
			key := strings.Join(parts[:i+1], "/")
			fmt.Println(key, "=>", value)
		}
	}
	// Output:
	// home/abc => my home directory
	// home/abc/Documents => my documents
	// home/abc/Documents/pic.png => cat pic
}
