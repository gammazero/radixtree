package radixtree

import (
	"testing"
)

func TestStepper(t *testing.T) {
	rt := new(Tree[string])
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")

	// (root) t-> ("o", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                     r-> ("n", TORN)

	iter := rt.NewStepper()
	if iter.Next('x') {
		t.Fatal("'x' should not have advanced iterator")
	}
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	val, ok := iter.Value()
	if ok || val != "" {
		t.Fatal("should not have value at 't'")
	}
	item := iter.Item()
	if item != nil {
		t.Fatal("should not have item at 't'")
	}
	if !iter.Next('o') {
		t.Fatal("'o' should have advanced iterator")
	}
	if _, ok = iter.Value(); ok {
		t.Fatal("should not have value at 'o'")
	}
	if iter.Next('o') {
		t.Fatal("'o' should not have advanced iterator")
	}

	// branch iterator
	iterR := iter.Copy()

	if !iter.Next('m') {
		t.Fatal("'m' should have advanced iterator")
	}
	val, ok = iter.Value()
	if !ok || val != "TOM" {
		t.Fatalf("expected \"TOM\" at 'm', got %q", val)
	}
	item = iter.Item()
	if item == nil || item.Value() != "TOM" {
		t.Fatalf("expected value \"TOM\" at 'm'")
	}
	if item.Key() != "tom" {
		t.Fatalf("expected key \"tom\" at 'm'")
	}
	if !iter.Next('a') {
		t.Fatal("'a' should have advanced iterator")
	}
	if _, ok = iter.Value(); ok {
		t.Fatal("should not have value at 'a'")
	}
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	if _, ok = iter.Value(); ok {
		t.Fatal("should not have value at 't'")
	}
	if !iter.Next('o') {
		t.Fatal("'o' should have advanced iterator")
	}
	val, ok = iter.Value()
	if !ok || val != "TOMATO" {
		t.Fatal("expected \"TOMATO\" 'o'")
	}

	if !iterR.Next('r') {
		t.Fatal("'r' should have advanced iterator")
	}
	if val, ok = iterR.Value(); ok {
		t.Fatal("should not have value at 'r', got ", val)
	}
	if !iterR.Next('n') {
		t.Fatal("'n' should have advanced iterator")
	}
	val, ok = iterR.Value()
	if !ok || val != "TORN" {
		t.Fatal("expected \"TORN\" 'n'")
	}
	if iterR.Next('n') {
		t.Fatal("'n' should not have advanced iterator")
	}

	iter = rt.NewStepper()
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	if iter.Next('x') {
		t.Fatal("'x' should not have advanced iterator")
	}
}
