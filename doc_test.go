package radixtree_test

import (
	"fmt"
	"strings"

	"github.com/gammazero/radixtree"
)

//nolint errcheck // Walk only returns error if user function returns error
func ExampleRunes_Walk() {
	rt := new(radixtree.Runes)
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items whose keys start with "tom"
	rt.Walk("tom", func(key fmt.Stringer, value interface{}) error {
		fmt.Println(value)
		return nil
	})
}

//nolint errcheck // WalkPath only returns error if user function returns error
func ExampleRunes_WalkPath() {
	rt := new(radixtree.Runes)
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items that are a prefix of "tomato"
	rt.WalkPath("tomato", func(key string, value interface{}) error {
		fmt.Println(value)
		return nil
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
		val := iter.Value()
		if val != nil {
			fmt.Println(val)
		}
	}
	// Output:
	// TOM
	// TOMATO
}

//nolint errcheck // Walk only returns error if user function returns error
func ExamplePaths_Walk() {
	pt := new(radixtree.Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")

	// Find all items with keys that start with "home/abc/Documents"
	pt.Walk("home/abc/Documents", func(key fmt.Stringer, value interface{}) error {
		fmt.Println(value)
		return nil
	})
}

//nolint errcheck // WalkPath only returns error if user function returns error
func ExamplePaths_WalkPath() {
	pt := new(radixtree.Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")

	// Find item in each path segment
	pt.WalkPath("home/abc/Documents/pic.png", func(key string, value interface{}) error {
		fmt.Println(key, "=>", value)
		return nil
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
		value := iter.Value()
		if value != nil {
			key := strings.Join(parts[:i+1], "/")
			fmt.Println(key, "=>", value)
		}
	}
	// Output:
	// home/abc => my home directory
	// home/abc/Documents => my documents
	// home/abc/Documents/pic.png => cat pic
}
