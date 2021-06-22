package radixtree

import (
	"testing"
)

func TestBytesAddEnd(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tomato", "TOMATO")
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "omato" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node.leaf.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.leaf.value)
	}
	if len(node.edges) != 0 {
		t.Fatal("child should have no children")
	}
	t.Log(dump(rt))
	// EX0: (root) t-> ("omato", TOMATO)
	//      then add "tom", TOM
	//      (root) t-> ("om", TOM) a-> ("to", TOMATO)
	//
	rt.Put("tom", "TOM")
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node = rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "om" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node.leaf.value != "TOM" {
		t.Fatal("wrong value at child:", node.leaf.value)
	}
	if len(node.edges) != 1 {
		t.Fatal("child should have 1 child")
	}
	node = node.getEdge('a')
	if node == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node.prefix) != "to" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node.leaf.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.leaf.value)
	}
	if len(node.edges) != 0 {
		t.Fatal("node should have no children")
	}
	t.Log(dump(rt))
}

func TestBytesAddFront(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	t.Log(dump(rt))
	// (root) t-> ("om", TOM)
	// then add "tomato", TOMATO
	// (root) t-> ("om", TOM) a-> ("to", TOMATO)
	t.Log("... add \"tomato\" TOMATO ...")
	rt.Put("tomato", "TOMATO")
	t.Log(dump(rt))
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "om" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node.leaf.value != "TOM" {
		t.Fatal("wrong value at child:", node.leaf.value)
	}
	if len(node.edges) != 1 {
		t.Fatal("child should have 1 child")
	}
	node = node.getEdge('a')
	if node == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node.prefix) != "to" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node.leaf.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.leaf.value)
	}
	if len(node.edges) != 0 {
		t.Fatal("node should have no children")
	}
}

func TestBytesAddBranch(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")

	// (root) t-> ("om", TOM) a-> ("to", TOMATO)
	// then add "torn", TORN
	// (root) t-> ("o", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                     r-> ("n", TORN)
	t.Log(dump(rt))
	t.Log("... add \"torn\", TORN ...")
	rt.Put("torn", "TORN")
	t.Log(dump(rt))
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "o" {
		t.Fatal("expected prefix 'o', got: ", node.prefix)
	}
	if node.leaf != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.edges) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.getEdge('m')
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if len(node2.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node2.leaf == nil {
		t.Fatal("missing value at node")
	}
	if node2.leaf.value != "TOM" {
		t.Fatal("wrong value at node:", node2.leaf.value)
	}
	if len(node2.edges) != 1 {
		t.Fatal("node should have 1 child")
	}
	node3 := node2.getEdge('a')
	if node3 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node3.prefix) != "to" {
		t.Fatal("expected prefix 'to', got: ", node3.prefix)
	}
	if node3.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node3.leaf.value != "TOMATO" {
		t.Fatal("expected value 'TOMATO', got:", node3.leaf.value)
	}
	if len(node3.edges) != 0 {
		t.Fatal("node should have no children")
	}
	node2 = node.getEdge('r')
	if node2 == nil {
		t.Fatal("node should have child at 'r'")
	}
	if string(node2.prefix) != "n" {
		t.Fatal("wrong prefix at node: ", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node2.leaf.value != "TORN" {
		t.Fatal("wrong value at node:", node2.leaf.value)
	}
	if len(node2.edges) != 0 {
		t.Fatal("node should have no children")
	}
}

func TestBytesAddBranchToBranch(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")

	// (root) t-> ("o", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                     r-> ("n", TORN)
	// then add "tag", TAG
	// (root) t-> ("", _) o-> ("", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                                r-> ("n", TORN)
	//                    a-> ("g", TAG)
	t.Log("... add \"tag\", TAG ...")
	rt.Put("tag", "TAG")
	t.Log(dump(rt))
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node.leaf != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.edges) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.getEdge('o')
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if len(node2.edges) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 = node.getEdge('a')
	if node2 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if len(node2.edges) != 0 {
		t.Fatal("node should have no children")
	}
	if string(node2.prefix) != "g" {
		t.Fatal("expected prefix 'g', got: ", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node2.leaf.value != "TAG" {
		t.Fatal("expected value 'TAG', got:", node2.leaf.value)
	}
}

func TestBytesAddExisting(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")
	rt.Put("tag", "TAG")

	// (root) t-> ("", _) o-> ("", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                                r-> ("n", TORN)
	//                    a-> ("g", TAG)
	// then add "to", TO
	// (root) t-> ("", _) o-> ("", TO) m-> ("", TOM) a-> ("to", TOMATO)
	//                                 r-> ("n", TORN)
	//                    a-> ("g", TAG)
	t.Log("... add \"to\", TO ...")
	rt.Put("to", "TO")
	t.Log(dump(rt))
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node.leaf != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.edges) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.getEdge('a')
	if node2 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if len(node2.edges) != 0 {
		t.Fatal("node should have no children")
	}
	node2 = node.getEdge('o')
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if node2.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node2.leaf.value != "TO" {
		t.Fatal("expected value 'TO', got:", node2.leaf.value)
	}
	if len(node2.edges) != 2 {
		t.Fatal("node should have 2 children")
	}
	node3 := node2.getEdge('m')
	if node3 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if node3 = node2.getEdge('r'); node3 == nil {
		t.Fatal("node should have child at 'r'")
	}
}

func TestBytesDelete(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")
	rt.Put("tag", "TAG")
	rt.Put("to", "TO")

	// Test that delete prunes
	if !rt.Delete("torn") {
		t.Error("did not delete \"torn\"")
	}
	node := rt.root.getEdge('t')
	node = node.getEdge('o')
	if node.getEdge('r') != nil {
		t.Error("deleted leaf should have been pruned")
	}

	// Test that delete compresses
	if !rt.Delete("tom") {
		t.Error("did not delete \"tom\"")
	}
	node = rt.root.getEdge('t')
	node = node.getEdge('o')
	node = node.getEdge('m')
	if node.leaf == nil && len(node.edges) == 1 {
		t.Log(dump(rt))
		t.Error("did not compress deleted node")
	}
	if string(node.prefix) != "ato" {
		t.Log(dump(rt))
		t.Error("wrong prefix for compresses node: ", node.prefix)
	}

	// Test deleting key that does not exist
	if rt.Delete("xyz") {
		t.Error("expected false when deleting key 'xyz'")
	}
}

func TestBytesBuildEdgeCases(t *testing.T) {
	tree := new(Bytes)

	tree.Put("ABCD", 1)
	t.Log(dump(tree))
	tree.Put("ABCDE", 2)
	t.Log(dump(tree))
	tree.Put("ABCDF", 3)
	t.Log(dump(tree))

	val, ok := tree.Get("ABCE")
	if ok || val != nil {
		t.Fatal("expected no value")
	}

	if tree.Delete("ABCE") {
		t.Fatal("should not delete non-existent value")
	}

	tree.Put("ABCE", 4)
	t.Log(dump(tree))

	tree.Put("ABCDEFGHIJK", 5)
	if tree.Delete("ABCDEFGH") {
		t.Fatal("should not delete non-existent value")
	}

	for _, k := range []string{"ABCDEFGHIJK", "ABCE", "ABCDF", "ABCDE", "ABCD"} {
		if !tree.Delete(k) {
			t.Error("failed to delete key ", k)
		}
	}

	// (root) /-> ("L1/L2", 1)
	tree.Put("/L1/L2", 1)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected 1 child, got ", len(tree.root.edges))
	}
	node := tree.root.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	if node.leaf.value != 1 {
		t.Fatal("expected value of 1, got ", node.leaf.value)
	}

	// (root) /-> ("L1/L2", 1)
	// add "/L1/L2/L3", 555
	// (root) /-> ("L1/L2", 1) /-> ("L3", 555)
	tree.Put("/L1/L2/L3", 555)
	t.Log(dump(tree))
	node = tree.root.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L3" {
		t.Fatal("expected prefix '/L3', got ", node.prefix)
	}
	if node.leaf.value != 555 {
		t.Fatal("expected value of 555, got ", node.leaf.value)
	}

	// (root) /-> ("L1/L2", 1) /-> ("L3", 555)
	// add "/L1/L2/L3/L4", 999
	// (root) /-> ("L1/L2", 1) /-> ("L3", 555) /-> ("L4", 999)
	tree.Put("/L1/L2/L3/L4", 999)
	t.Log(dump(tree))
	node = tree.root.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L3" {
		t.Fatal("expected prefix '/L3', got ", node.prefix)
	}
	if node.leaf.value != 555 {
		t.Fatal("expected value of 555, got ", node.leaf.value)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L4" {
		t.Fatal("expected prefix '/L4', got ", node.prefix)
	}
	if node.leaf.value != 999 {
		t.Fatal("expected value of 999, got ", node.leaf.value)
	}

	// (root) /-> ("L1/L2", 1) /-> ("L3", 555) /-> ("L4", 999)
	// add "/L1/L2/L/C", 3
	// (root) /-> ("L1/L2", 1) /-> ("L", _) 3-> ("L3", 555) /-> ("L4", 999)
	//                                      /-> ("C", 3)
	tree.Put("/L1/L2/L/C", 3)
	t.Log(dump(tree))
	node = tree.root.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", string(node.prefix))
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L" {
		t.Fatal("expected prefix 'L', got ", string(node.prefix))
	}
	if node.leaf != nil {
		t.Fatal("expected nil value, got ", node.leaf.value)
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node.edges))
	}
	//t.Fatal("hre")

	t.Log(dump(tree))
	tree.Put("/L1/L2/L3/X", 999)
	t.Log(dump(tree))
}

func TestBytesCopyIterator(t *testing.T) {
	rt := new(Bytes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")

	// (root) t-> ("o", _) m-> ("", TOM) a-> ("to", TOMATO)
	//                     r-> ("n", TORN)

	iter := rt.NewIterator()
	if iter.Next('x') {
		t.Fatal("'x' should not have advanced iterator")
	}
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	val, ok := iter.Value()
	if ok || val != nil {
		t.Fatal("should not have value at 't'")
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

	iter = rt.NewIterator()
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	if iter.Next('x') {
		t.Fatal("'x' should not have advanced iterator")
	}
}

func TestSimpleBytesWalk(t *testing.T) {
	rt := New()
	rt.Put("tomato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("tornado", "TORNADO")

	count := 0
	rt.Walk("tomato", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 1 {
		t.Errorf("expected to visit 1 key, visited %d", count)
	}

	count = 0
	rt.Walk("t", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 3 {
		t.Errorf("expected to visit 3 keys, visited %d", count)
	}

	count = 0
	rt.Walk("to", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 3 {
		t.Errorf("expected to visit 3 keys, visited %d", count)
	}

	count = 0
	rt.Walk("tom", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 2 {
		t.Errorf("expected to visit 2 keys, visited %d", count)
	}

	count = 0
	rt.Walk("tomx", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 0 {
		t.Errorf("expected to visit 0 keys, visited %d", count)
	}

	count = 0
	rt.Walk("torn", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 1 {
		t.Errorf("expected to visit 1 key, visited %d", count)
	}
}

func TestBytes(t *testing.T) {
	testRadixTree(t, New())
}

func TestBytesNilGet(t *testing.T) {
	testNilGet(t, New())
}

func TestBytesRoot(t *testing.T) {
	testRoot(t, New())
}

func TestBytesWalk(t *testing.T) {
	testWalk(t, New())
}

func TestBytesWalkStop(t *testing.T) {
	testWalkStop(t, New())
}

func TestBytesInspectStop(t *testing.T) {
	testInspectStop(t, New())
}

func TestGetAfterDelete(t *testing.T) {
	testGetAfterDelete(t, New())
}

func TestBytesStringConvert(t *testing.T) {
	tree := New()
	for _, w := range []string{"Bart", `Bartók`, `AbónXw`, `AbónYz`} {
		ok := tree.Put(w, w)
		if !ok {
			t.Error("did not insert new value", w)
		}

		v, _ := tree.Get(w)
		if v == nil {
			t.Log(dump(tree))
			t.Fatal("nil value returned getting", w)
		}
		s, ok := v.(string)
		if !ok {
			t.Fatal("value is not a string")
		}
		if s != w {
			t.Fatalf("returned wrong value - expected %q got %q", w, s)
		}
	}
	tree.Walk("", func(key string, val interface{}) bool {
		t.Log("Key:", key)
		s, ok := val.(string)
		if !ok {
			t.Log(dump(tree))
			t.Fatal("value is not a string")
		}
		t.Log("Val:", s)
		if key != s {
			t.Log(dump(tree))
			t.Fatal("Key and value do not match")
		}
		return false
	})
}
