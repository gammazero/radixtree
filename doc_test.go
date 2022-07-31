package radixtree_test

import (
	"fmt"

	"github.com/gammazero/radixtree"
)

func Example_WalkFrom() {
	rt := radixtree.New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items whose keys start with "tom"
	rt.WalkFrom("tom", func(key string, value interface{}) bool {
		fmt.Println(value)
		return false
	})
}

func Example_WalkTo() {
	rt := radixtree.New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items that are a prefix of "tomato"
	rt.WalkTo("tomato", func(key string, value interface{}) bool {
		fmt.Println(value)
		return false
	})
	// Output:
	// TOM
	// TOMATO
}

func Example_NewIterator() {
	rt := radixtree.New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	iter := rt.NewIterator()
	for {
		key, val, done := iter.Next()
		if done {
			break
		}
		fmt.Println(key, "=", val)
	}

	// Output:
	// tom = TOM
	// tomato = TOMATO
	// tommy = TOMMY
	// tornado = TORNADO
}

func Example_NewStepper() {
	rt := radixtree.New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	iter := rt.NewStepper()
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
