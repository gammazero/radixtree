package radixtree_test

import (
	"fmt"
	"strings"

	"github.com/gammazero/radixtree"
)

func ExampleBytes_Walk() {
	rt := radixtree.New()
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

func ExampleBytes_WalkPath() {
	rt := radixtree.New()
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

func ExampleBytes_NewIterator() {
	rt := radixtree.New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	iter := rt.NewIterator()
	word := "tomato"
	for i := range word {
		if !iter.Next(word[i]) {
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
	pt := radixtree.NewPaths("/")
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
	pt := radixtree.NewPaths("/")
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
	pt := radixtree.NewPaths("/")
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
