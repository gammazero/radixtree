package radixtree_test

import (
	"fmt"

	"github.com/gammazero/radixtree"
)

func ExampleTree_Iter() {
	rt := radixtree.New[int]()
	rt.Put("mercury", 1)
	rt.Put("venus", 2)
	rt.Put("earth", 3)
	rt.Put("mars", 4)

	// Find all items that that have a key that is a prefix of "tomato".
	for key, value := range rt.Iter() {
		fmt.Println(key, "=", value)
	}
	// Output:
	// earth = 3
	// mars = 4
	// mercury = 1
	// venus = 2
}

func ExampleTree_IterAt() {
	rt := radixtree.New[string]()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items whose keys start with "tom".
	for _, value := range rt.IterAt("tom") {
		fmt.Println(value)
	}
	// Output:
	// TOM
	// TOMATO
	// TOMMY
}

func ExampleTree_IterPath() {
	rt := radixtree.New[string]()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	// Find all items that that have a key that is a prefix of "tomato".
	for key, value := range rt.IterPath("tomato") {
		fmt.Println(key, "=", value)
	}
	// Output:
	// tom = TOM
	// tomato = TOMATO
}

func ExampleTree_NewStepper() {
	rt := radixtree.New[string]()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tommy", "TOMMY")
	rt.Put("tornado", "TORNADO")

	stepper := rt.NewStepper()
	word := "tomato"
	for i := range word {
		if !stepper.Next(word[i]) {
			break
		}
		if val, ok := stepper.Value(); ok {
			fmt.Println(val)
		}
	}
	// Output:
	// TOM
	// TOMATO
}
