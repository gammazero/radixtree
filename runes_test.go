package radixtree

import (
	"testing"
)

func TestRunesAddEnd(t *testing.T) {
	rt := new(Runes)
	rt.Put("tomato", "TOMATO")
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "omato" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.value)
	}
	if len(node.children) != 0 {
		t.Fatal("child should have no children")
	}
	t.Log(dump(rt))
	// EX0: (root) t-> ("omato", TOMATO)
	//      then add "tom", TOM
	//      (root) t-> ("om", TOM) a-> ("to", TOMATO)
	//
	rt.Put("tom", "TOM")
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node = rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "om" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.value != "TOM" {
		t.Fatal("wrong value at child:", node.value)
	}
	if len(node.children) != 1 {
		t.Fatal("child should have 1 child")
	}
	node = node.children['a']
	if node == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node.prefix) != "to" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.value)
	}
	if len(node.children) != 0 {
		t.Fatal("node should have no children")
	}
	t.Log(dump(rt))
}

func TestRunesAddFront(t *testing.T) {
	rt := new(Runes)
	rt.Put("tom", "TOM")
	t.Log(dump(rt))
	// (root) t-> ("om", TOM)
	// then add "tomato", TOMATO
	// (root) t-> ("om", TOM) a-> ("to", TOMATO)
	t.Log("... add \"tomato\" TOMATO ...")
	rt.Put("tomato", "TOMATO")
	t.Log(dump(rt))
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "om" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.value != "TOM" {
		t.Fatal("wrong value at child:", node.value)
	}
	if len(node.children) != 1 {
		t.Fatal("child should have 1 child")
	}
	node = node.children['a']
	if node == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node.prefix) != "to" {
		t.Fatal("wrong prefix at child:", node.prefix)
	}
	if node.value != "TOMATO" {
		t.Fatal("wrong value at child:", node.value)
	}
	if len(node.children) != 0 {
		t.Fatal("node should have no children")
	}
}

func TestRunesAddBranch(t *testing.T) {
	rt := new(Runes)
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
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if string(node.prefix) != "o" {
		t.Fatal("expected prefix 'o', got: ", node.prefix)
	}
	if node.value != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.children) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.children['m']
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if len(node2.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node2.value != "TOM" {
		t.Fatal("wrong value at node:", node2.value)
	}
	if len(node2.children) != 1 {
		t.Fatal("node should have 1 child")
	}
	node3 := node2.children['a']
	if node3 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if string(node3.prefix) != "to" {
		t.Fatal("expected prefix 'to', got: ", node3.prefix)
	}
	if node3.value != "TOMATO" {
		t.Fatal("expected value 'TOMATO', got:", node3.value)
	}
	if len(node3.children) != 0 {
		t.Fatal("node should have no children")
	}
	node2 = node.children['r']
	if node2 == nil {
		t.Fatal("node should have child at 'r'")
	}
	if string(node2.prefix) != "n" {
		t.Fatal("wrong prefix at node: ", node2.prefix)
	}
	if node2.value != "TORN" {
		t.Fatal("wrong value at node:", node2.value)
	}
	if len(node2.children) != 0 {
		t.Fatal("node should have no children")
	}
}

func TestRunesAddBranchToBranch(t *testing.T) {
	rt := new(Runes)
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
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node.value != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.children) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.children['o']
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if len(node2.children) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 = node.children['a']
	if node2 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if len(node2.children) != 0 {
		t.Fatal("node should have no children")
	}
	if string(node2.prefix) != "g" {
		t.Fatal("expected prefix 'g', got: ", node2.prefix)
	}
	if node2.value != "TAG" {
		t.Fatal("expected value 'TAG', got:", node2.value)
	}
}

func TestRunesAddExisting(t *testing.T) {
	rt := new(Runes)
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
	if len(rt.children) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.children['t']
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("node should not have prefix")
	}
	if node.value != nil {
		t.Fatal("node should have nil value")
	}
	if len(node.children) != 2 {
		t.Fatal("node should have 2 children")
	}
	node2 := node.children['a']
	if node2 == nil {
		t.Fatal("node should have child at 'a'")
	}
	if len(node2.children) != 0 {
		t.Fatal("node should have no children")
	}
	node2 = node.children['o']
	if node2 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if node2.value != "TO" {
		t.Fatal("expected value 'TO', got:", node2.value)
	}
	if len(node2.children) != 2 {
		t.Fatal("node should have 2 children")
	}
	node3 := node2.children['m']
	if node3 == nil {
		t.Fatal("node should have child at 'm'")
	}
	if node3 = node2.children['r']; node3 == nil {
		t.Fatal("node should have child at 'r'")
	}
}

func TestRunesDelete(t *testing.T) {
	rt := new(Runes)
	rt.Put("tom", "TOM")
	rt.Put("tomato", "TOMATO")
	rt.Put("torn", "TORN")
	rt.Put("tag", "TAG")
	rt.Put("to", "TO")

	// Test that delete prunes
	if !rt.Delete("torn") {
		t.Error("did not delete \"torn\"")
	}
	node := rt.children['t']
	node = node.children['o']
	if _, ok := node.children['r']; ok {
		t.Error("deleted leaf should have been pruned")
	}

	// Test that delete compresses
	if !rt.Delete("tom") {
		t.Error("did not delete \"tom\"")
	}
	node = rt.children['t']
	node = node.children['o']
	node = node.children['m']
	if node.value == nil && len(node.children) == 1 {
		t.Log(dump(rt))
		t.Error("did not compress deleted node")
	}
	if string(node.prefix) != "ato" {
		t.Log(dump(rt))
		t.Error("worng prefix for compresses node: ", node.prefix)
	}
}

func TestRunesBuildEdgeCases(t *testing.T) {
	tree := new(Runes)

	tree.Put("ABCD", 1)
	t.Log(dump(tree))
	tree.Put("ABCDE", 2)
	t.Log(dump(tree))
	tree.Put("ABCDF", 3)
	t.Log(dump(tree))

	val := tree.Get("ABCE")
	if val != nil {
		t.Fatal("expected nil value")
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
	if len(tree.children) != 1 {
		t.Fatal("expected 1 child, got ", len(tree.children))
	}
	node := tree.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	if node.value != 1 {
		t.Fatal("expected value of 1, got ", node.value)
	}

	// (root) /-> ("L1/L2", 1)
	// add "/L1/L2/L3", 555
	// (root) /-> ("L1/L2", 1) /-> ("L3", 555)
	tree.Put("/L1/L2/L3", 555)
	t.Log(dump(tree))
	node = tree.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L3" {
		t.Fatal("expected prefix '/L3', got ", node.prefix)
	}
	if node.value != 555 {
		t.Fatal("expected value of 555, got ", node.value)
	}

	// (root) /-> ("L1/L2", 1) /-> ("L3", 555)
	// add "/L1/L2/L3/L4", 999
	// (root) /-> ("L1/L2", 1) /-> ("L3", 555) /-> ("L4", 999)
	tree.Put("/L1/L2/L3/L4", 999)
	t.Log(dump(tree))
	node = tree.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L3" {
		t.Fatal("expected prefix '/L3', got ", node.prefix)
	}
	if node.value != 555 {
		t.Fatal("expected value of 555, got ", node.value)
	}
	node = node.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L4" {
		t.Fatal("expected prefix '/L4', got ", node.prefix)
	}
	if node.value != 999 {
		t.Fatal("expected value of 999, got ", node.value)
	}

	// (root) /-> ("L1/L2", 1) /-> ("L3", 555) /-> ("L4", 999)
	// add "/L1/L2/L/C", 3
	// (root) /-> ("L1/L2", 1) /-> ("L", _) 3-> ("L3", 555) /-> ("L4", 999)
	//                                      /-> ("C", 3)
	tree.Put("/L1/L2/L/C", 3)
	t.Log(dump(tree))
	node = tree.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", string(node.prefix))
	}
	node = node.children['/']
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if string(node.prefix) != "L" {
		t.Fatal("expected prefix 'L', got ", string(node.prefix))
	}
	if node.value != nil {
		t.Fatal("expected nil value, got ", node.value)
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", len(node.children))
	}
	//t.Fatal("hre")

	t.Log(dump(tree))
	tree.Put("/L1/L2/L3/X", 999)
	t.Log(dump(tree))
}

func TestRunesCopyIterator(t *testing.T) {
	rt := new(Runes)
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
	if iter.Value() != nil {
		t.Fatal("should not have value at 't'")
	}
	if !iter.Next('o') {
		t.Fatal("'o' should have advanced iterator")
	}
	if iter.Value() != nil {
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
	if iter.Value() != "TOM" {
		t.Fatalf("expected \"TOM\" at 'm', got %q", iter.Value())
	}
	if !iter.Next('a') {
		t.Fatal("'a' should have advanced iterator")
	}
	if iter.Value() != nil {
		t.Fatal("should not have value at 'a'")
	}
	if !iter.Next('t') {
		t.Fatal("'t' should have advanced iterator")
	}
	if iter.Value() != nil {
		t.Fatal("should not have value at 't'")
	}
	if !iter.Next('o') {
		t.Fatal("'o' should have advanced iterator")
	}
	if iter.Value() != "TOMATO" {
		t.Fatal("expected \"TOMATO\" 'o'")
	}

	if !iterR.Next('r') {
		t.Fatal("'r' should have advanced iterator")
	}
	if iterR.Value() != nil {
		t.Fatal("should not have value at 'r', got ", iterR.Value())
	}
	if !iterR.Next('n') {
		t.Fatal("'n' should have advanced iterator")
	}
	if iterR.Value() != "TORN" {
		t.Fatal("expected \"TORN\" 'n'")
	}
	if iterR.Next('n') {
		t.Fatal("'n' should not have advanced iterator")
	}
}

func TestRunes(t *testing.T) {
	testRadixTree(t, new(Runes))
}

func TestRunesNilGet(t *testing.T) {
	testNilGet(t, new(Runes))
}

func TestRunesRoot(t *testing.T) {
	testRoot(t, new(Runes))
}

func TestRunesWalk(t *testing.T) {
	testWalk(t, new(Runes))
}

func TestRunesWalkError(t *testing.T) {
	testWalkError(t, new(Runes))
}

func TestRunesWalkSkip(t *testing.T) {
	testWalkSkip(t, new(Runes))
}

func TestRunesInspectSkip(t *testing.T) {
	testInspectSkip(t, new(Runes))
}
