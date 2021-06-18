package radixtree

import (
	"strings"
	"testing"
)

func TestPathsAddEnd(t *testing.T) {
	// add "/L1/L2", 1
	// (root) /L1-> ("/L2", 1)
	tree := new(Paths)
	tree.Put("/L1/L2", 1)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if strings.Join(node.prefix, "") != "L2" {
		t.Fatal("wrong prefix:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("child missing value")
	}
	if node.leaf.value != 1 {
		t.Fatal("expected value 1, got ", node.leaf.value)
	}
	if len(node.edges) != 0 {
		t.Fatal("expected no children")
	}

	// (root) /L1-> ("/L2", 1)
	// add "/L1/L2/L3A", 2
	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	tree.Put("/L1/L2/L3A", 2)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node = tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if strings.Join(node.prefix, "") != "L2" {
		t.Fatal("wrong prefix:", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("child missing value")
	}
	if node.leaf.value != 1 {
		t.Fatal("expected value 1, got ", node.leaf.value)
	}
	if len(node.edges) != 1 {
		t.Fatal("expected 1 child")
	}
	node = node.getEdge("L3A")
	if node == nil {
		t.Fatal("expected child at 'L3A'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("child missing value")
	}
	if node.leaf.value != 2 {
		t.Fatal("expected value 3, got ", node.leaf.value)
	}
	if len(node.edges) != 0 {
		t.Fatal("expected no children")
	}
}

func TestPathsAddBranch(t *testing.T) {
	tree := NewPaths(".")
	tree.Put(".L1.L2", 1)
	tree.Put(".L1.L2.L3A", 2)

	// (root) .L1-> (".L2", 1) .L3A-> ("", 2)
	// add ".L1.L2.L3B.L4", 3
	// (root) .L1-> (".L2", 1) .L3A-> ("", 2)
	//                         .L3B-> (".L4", 3)
	tree.Put(".L1.L2.L3B.L4", 3)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children")
	}
	node2 := node.getEdge("L3B")
	if node2 == nil {
		t.Fatal("expected child at 'L3B'")
	}
	if strings.Join(node2.prefix, "") != "L4" {
		t.Fatal("wrong prefix:", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("child missing value")
	}
	if node2.leaf.value != 3 {
		t.Fatal("expected value 3, got ", node2.leaf.value)
	}
	node2 = node.getEdge("L3A")
	if node2 == nil {
		t.Fatal("expected child at 'L3A'")
	}
}

func TestPathsAddBranchToBranch(t *testing.T) {
	tree := NewPaths("/")
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)
	tree.Put("/L1/L2/L3B/L4", 3)

	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	//                         /L3B-> ("/L4", 3)
	// add "/L1/L2B/L3C", 4
	// (root) /L1-> ("", _) /L2-> ("", 1) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	tree.Put("/L1/L2B/L3C", 4)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.leaf != nil {
		t.Fatal("expected nil value, got ", node.leaf.value)
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node.edges))
	}
	node2 := node.getEdge("L2B")
	if node2 == nil {
		t.Fatal("Expected child at 'L2B'")
	}
	if strings.Join(node2.prefix, "") != "L3C" {
		t.Fatal("expected prefix 'L3C', got ", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("child missing value")
	}
	if node2.leaf.value != 4 {
		t.Fatal("expected value of 4, got ", node2.leaf.value)
	}
	node2 = node.getEdge("L2")
	if node2 == nil {
		t.Fatal("expected child at 'L2'")
	}
	if len(node2.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("child missing value")
	}
	if node2.leaf.value != 1 {
		t.Fatal("expected value of 1, got ", node2.leaf.value)
	}
	if len(node2.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node2.edges))
	}
}

func TestPathsAddExisting(t *testing.T) {
	tree := NewPaths("/")
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)
	tree.Put("/L1/L2/L3B/L4", 3)
	tree.Put("/L1/L2B/L3C", 4)

	// (root) /L1-> ("", _) /L2-> ("", 1) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	// add "/L1/L2/L3", 5
	// (root) /L1-> ("", 5) /L2-> ("", 1) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	tree.Put("/L1", 5)
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("child missing value")
	}
	if node.leaf.value != 5 {
		t.Fatal("expected value of 5, got ", node.leaf.value)
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node.edges))
	}
}

func TestPathsDelete(t *testing.T) {
	tree := NewPaths("/")
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)
	tree.Put("/L1/L2/L3B/L4", 3)
	tree.Put("/L1/L2B/L3C", 4)
	tree.Put("/L1", 5)

	// (root) /L1-> ("", 5) /L2-> ("", 1) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	// delete "/L1/L2"
	// (root) /L1-> ("", 5) /L2-> ("", _) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	tree.Delete("/L1/L2")
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", node.edges)
	}
	node = node.getEdge("L2")
	if node == nil {
		t.Fatal("expected child at 'L2'")
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", node.edges)
	}

	// Delete a key that does not exist
	if tree.Delete("/L1/L2/L2B/L4") {
		t.Fatal("should not have deleted non-existent key")
	}

	// (root) /L1-> ("", 5) /L2-> ("", _) /L3A-> ("", 2)
	//                                    /L3B-> ("/L4", 3)
	//
	//                      /L2B-> ("L3C", 4)
	// delete "/L1/L2/L2B/L4"
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	//                      /L2B-> ("L3C", 4)
	if !tree.Delete("/L1/L2/L3B/L4") {
		t.Fatal("should have deleted key")
	}
	t.Log(dump(tree))
	if len(tree.root.edges) != 1 {
		t.Fatal("expected one child")
	}
	node = tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node.edges))
	}
	node = node.getEdge("L2")
	if node == nil {
		t.Fatal("expected child at 'L2'")
	}
	if len(node.edges) != 0 {
		t.Fatal("expected 0 children, got ", len(node.edges))
	}
	if strings.Join(node.prefix, "") != "L3A" {
		t.Fatal("expected prefix 'L3A', got ", node.prefix)
	}
	if node.leaf == nil {
		t.Fatal("child missing value")
	}
	if node.leaf.value != 2 {
		t.Fatal("expected value of 2, got ", node.leaf.value)
	}

	// Test that Delete prunes
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	//                      /L2B-> ("L3C", 4)
	// delete /L1/L2B/L3C
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	if !tree.Delete("/L1/L2B/L3C") {
		t.Error("did not delete \"/L1/L2B/L3C\"")
	}
	node = tree.root.getEdge("L1")
	if node.getEdge("L2B") != nil {
		t.Log(dump(tree))
		t.Error("deleted leaf should have been pruned")
	}

	// Test that Delete compresses
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	// delete /L1
	// (root) /L1-> ("/L2/L3A", 2)
	if !tree.Delete("L1") {
		t.Error("did not delete \"/L1/L2B/L3C\"")
	}
	node = tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "/") != "L2/L3A" {
		t.Log(dump(tree))
		t.Error("wrong prefix for compresses node:", strings.Join(node.prefix, ""))
	}

	tree.Put("/L1/L2/L3A/L4", 6)

	// Check that Delete prunes up to node with value
	// (root) /L1-> ("/L2/L3A", 2) /L4->("", 6)
	// delete /L1/L2/L3A/L4
	// (root) /L1-> ("/L2/L3A", 2)
	if !tree.Delete("/L1/L2/L3A/L4") {
		t.Error("did not delete \"/L1/L2/L3A/L4\"")
	}
	node = tree.root.getEdge("L1")
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "/") != "L2/L3A" {
		t.Log(dump(tree))
		t.Error("wrong prefix for compresses node:", strings.Join(node.prefix, ""))
	}
	if len(node.edges) != 0 {
		t.Log(dump(tree))
		t.Error("node should not have children")
	}
}

func TestPathsCopyIterator(t *testing.T) {
	tree := NewPaths("/")
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)
	tree.Put("/L1/L2/L3B/L4", 3)

	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	//                         /L3B-> ("/L4", 3)

	iter := tree.NewIterator()
	if iter.Next("") {
		t.Fatal("\"\" should not have advanced iterator")
	}
	if iter.Next("x") {
		t.Fatal("x should not have advanced iterator")
	}
	if !iter.Next("L1") {
		t.Fatal("L1 should have advanced iterator")
	}
	val, ok := iter.Value()
	if ok || val != nil {
		t.Fatal("should not have value at /L1")
	}
	if !iter.Next("L2") {
		t.Fatal("L2 should have advanced iterator")
	}
	val, ok = iter.Value()
	if !ok || val != 1 {
		t.Fatal("expected value 1 at L2, got ", val)
	}
	if iter.Next("L4") {
		t.Fatal("L4 should not have advanced iterator")
	}

	// branch iterator
	iterB := iter.Copy()
	if !iterB.Next("L3B") {
		t.Fatal("L3B should have advanced iterator")
	}
	if _, ok = iterB.Value(); ok {
		t.Fatal("should not have value at L3B")
	}
	if !iterB.Next("L4") {
		t.Fatal("L4 should have advanced iterator")
	}
	val, ok = iterB.Value()
	if !ok || val != 3 {
		t.Fatal("expected value 3 at L4, got ", val)
	}
	if iterB.Next("L4") {
		t.Fatal("L4 should not have advanced iterator")
	}

	if !iter.Next("L3A") {
		t.Fatal("L3A should have advanced iterator")
	}
	val, ok = iter.Value()
	if !ok || val != 2 {
		t.Fatal("expected value 2 at L3A, got ", val)
	}
	if iter.Next("L3B") {
		t.Fatal("L3B should not have advanced iterator")
	}

}

func TestSimplePathWalk(t *testing.T) {
	rt := NewPaths("/")
	rt.Put("tom/ato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("torn/ad/o", "TORNADO")

	count := 0
	rt.Walk("tom/ato", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 1 {
		t.Errorf("expected to visit 1 key, visited %d", count)
	}

	count = 0
	rt.Walk("tom", func(key string, value interface{}) bool {
		count++
		return false
	})
	if count != 2 {
		t.Errorf("expected to visit 3 keys, visited %d", count)
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

func TestPathsEdgeSort(t *testing.T) {
	var edges pathEdges = []pathEdge{pathEdge{"xyz/987", nil}, pathEdge{"abc/123", nil}}

	if edges.Len() != 2 {
		t.Fatal("bad Len")
	}
	if edges.Less(0, 1) {
		t.Fatal("bad Less")
	}
	if !edges.Less(1, 0) {
		t.Fatal("bad Less")
	}
	edges.Swap(0, 1)
	if !edges.Less(0, 1) {
		t.Fatal("bad Swap")
	}
	if edges.Less(1, 0) {
		t.Fatal("bad Swap")
	}
}

func TestPaths(t *testing.T) {
	testRadixTree(t, NewPaths("/"))
}

func TestPathsNilGet(t *testing.T) {
	testNilGet(t, NewPaths("/"))
}

func TestPathsRoot(t *testing.T) {
	testRoot(t, NewPaths("/"))
}

func TestPathsWalk(t *testing.T) {
	testWalk(t, NewPaths("/"))
}

func TestPathsWalkStop(t *testing.T) {
	testWalkStop(t, NewPaths("/"))
}

func TestPathsInspectStop(t *testing.T) {
	testInspectStop(t, NewPaths("/"))
}

func TestPathsGetAfterDelete(t *testing.T) {
	testGetAfterDelete(t, NewPaths("/"))
}
