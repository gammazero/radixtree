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
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if strings.Join(node.prefix, "") != "/L2" {
		t.Fatal("wrong prefix:", node.prefix)
	}
	if node.value != 1 {
		t.Fatal("expected value 1, got ", node.value)
	}
	if len(node.children) != 0 {
		t.Fatal("expected no children")
	}

	// (root) /L1-> ("/L2", 1)
	// add "/L1/L2/L3A", 2
	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	tree.Put("/L1/L2/L3A", 2)
	t.Log(dump(tree))
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node = tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if strings.Join(node.prefix, "") != "/L2" {
		t.Fatal("wrong prefix:", node.prefix)
	}
	if node.value != 1 {
		t.Fatal("expected value 1, got ", node.value)
	}
	if len(node.children) != 1 {
		t.Fatal("expected 1 child")
	}
	node = node.children["/L3A"]
	if node == nil {
		t.Fatal("expected child at '/L3A'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.value != 2 {
		t.Fatal("expected value 3, got ", node.value)
	}
	if len(node.children) != 0 {
		t.Fatal("expected no children")
	}
}

func TestPathsAddBranch(t *testing.T) {
	tree := new(Paths)
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)

	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	// add "/L1/L2/L3B/L4", 3
	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	//                         /L3B-> ("/L4", 3)
	tree.Put("/L1/L2/L3B/L4", 3)
	t.Log(dump(tree))
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children")
	}
	node2 := node.children["/L3B"]
	if node2 == nil {
		t.Fatal("expected child at '/L3B'")
	}
	if strings.Join(node2.prefix, "") != "/L4" {
		t.Fatal("wrong prefix:", node2.prefix)
	}
	if node2.value != 3 {
		t.Fatal("expected value 3, got ", node2.value)
	}
	node2 = node.children["/L3A"]
	if node2 == nil {
		t.Fatal("expected child at '/L3A'")
	}
}

func TestPathsAddBranchToBranch(t *testing.T) {
	tree := new(Paths)
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
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.value != nil {
		t.Fatal("expected nil value, got ", node.value)
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", len(node.children))
	}
	node2 := node.children["/L2B"]
	if node2 == nil {
		t.Fatal("Expected child at '/L2B'")
	}
	if strings.Join(node2.prefix, "") != "/L3C" {
		t.Fatal("expected prefix '/L3C', got ", node2.prefix)
	}
	if node2.value != 4 {
		t.Fatal("expected value of 4, got ", node2.value)
	}
	node2 = node.children["/L2"]
	if node2 == nil {
		t.Fatal("expected child at '/L2'")
	}
	if len(node2.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node2.prefix)
	}
	if node2.value != 1 {
		t.Fatal("expected value of 1, got ", node2.value)
	}
	if len(node2.children) != 2 {
		t.Fatal("expected 2 children, got ", len(node2.children))
	}
}

func TestPathsAddExisting(t *testing.T) {
	tree := new(Paths)
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
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if len(node.prefix) != 0 {
		t.Fatal("expected no prefix, got ", node.prefix)
	}
	if node.value != 5 {
		t.Fatal("expected value of 5, got ", node.value)
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", len(node.children))
	}
}

func TestPathsDelete(t *testing.T) {
	tree := new(Paths)
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
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node := tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", node.children)
	}
	node = node.children["/L2"]
	if node == nil {
		t.Fatal("expected child at '/L2'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", node.children)
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
	if len(tree.children) != 1 {
		t.Fatal("expected one child")
	}
	node = tree.children["/L1"]
	if node == nil {
		t.Fatal("expected child at '/L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", node.children)
	}
	node = node.children["/L2"]
	if node == nil {
		t.Fatal("expected child at '/L2'")
	}
	if len(node.children) != 0 {
		t.Fatal("expected 0 children, got ", node.children)
	}
	if strings.Join(node.prefix, "") != "/L3A" {
		t.Fatal("expected prefix '/L3A', got ", node.prefix)
	}
	if node.value != 2 {
		t.Fatal("expected value of 2, got ", node.value)
	}

	// Test that Delete prunes
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	//                      /L2B-> ("L3C", 4)
	// delete /L1/L2B/L3C
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	if !tree.Delete("/L1/L2B/L3C") {
		t.Error("did not delete \"/L1/L2B/L3C\"")
	}
	node = tree.children["/L1"]
	if _, ok := node.children["/L2B"]; ok {
		t.Log(dump(tree))
		t.Error("deleted leaf should have been pruned")
	}

	// Test that Delete compresses
	// (root) /L1-> ("", 5) /L2-> ("/L3A", 2)
	// delete /L1
	// (root) /L1-> ("/L2/L3A", 2)
	if !tree.Delete("/L1") {
		t.Error("did not delete \"/L1/L2B/L3C\"")
	}
	node = tree.children["/L1"]
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "") != "/L2/L3A" {
		t.Log(dump(tree))
		t.Error("worng prefix for compresses node:", strings.Join(node.prefix, ""))
	}

	tree.Put("/L1/L2/L3A/L4", 6)

	// Check that Delete prunes up to node with value
	// (root) /L1-> ("/L2/L3A", 2) /L4->("", 6)
	// delete /L1/L2/L3A/L4
	// (root) /L1-> ("/L2/L3A", 2)
	if !tree.Delete("/L1/L2/L3A/L4") {
		t.Error("did not delete \"/L1/L2/L3A/L4\"")
	}
	node = tree.children["/L1"]
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "") != "/L2/L3A" {
		t.Log(dump(tree))
		t.Error("worng prefix for compresses node:", strings.Join(node.prefix, ""))
	}
	if len(node.children) != 0 {
		t.Log(dump(tree))
		t.Error("node should not have children")
	}
}

func TestPathsCopyIterator(t *testing.T) {
	tree := new(Paths)
	tree.Put("/L1/L2", 1)
	tree.Put("/L1/L2/L3A", 2)
	tree.Put("/L1/L2/L3B/L4", 3)

	// (root) /L1-> ("/L2", 1) /L3A-> ("", 2)
	//                         /L3B-> ("/L4", 3)

	iter := tree.NewIterator()
	if iter.Next("/x") {
		t.Fatal("/x should not have advanced iterator")
	}
	if !iter.Next("/L1") {
		t.Fatal("/L1 should have advanced iterator")
	}
	if iter.Value() != nil {
		t.Fatal("should not have value at /L1")
	}
	if !iter.Next("/L2") {
		t.Fatal("/L2 should have advanced iterator")
	}
	if iter.Value() != 1 {
		t.Fatal("expected value 1 at /L2, got ", iter.Value())
	}
	if iter.Next("/L4") {
		t.Fatal("/L4 should not have advanced iterator")
	}

	// branch iterator
	iterB := iter.Copy()
	if !iterB.Next("/L3B") {
		t.Fatal("/L3B should have advanced iterator")
	}
	if iterB.Value() != nil {
		t.Fatal("should not have value at /L3B")
	}
	if !iterB.Next("/L4") {
		t.Fatal("/L4 should have advanced iterator")
	}
	if iterB.Value() != 3 {
		t.Fatal("expected value 3 at /L4, got ", iterB.Value())
	}
	if iterB.Next("/L4") {
		t.Fatal("/L4 should not have advanced iterator")
	}

	if !iter.Next("/L3A") {
		t.Fatal("/L3A should have advanced iterator")
	}
	if iter.Value() != 2 {
		t.Fatal("expected value 2 at /L3A, got ", iter.Value())
	}
	if iter.Next("/L3B") {
		t.Fatal("/L3B should not have advanced iterator")
	}

}

func TestPaths(t *testing.T) {
	testRadixTree(t, new(Paths))
}

func TestPathsNilGet(t *testing.T) {
	testNilGet(t, new(Paths))
}

func TestPathsRoot(t *testing.T) {
	testRoot(t, new(Paths))
}

func TestPathsWalk(t *testing.T) {
	testWalk(t, new(Paths))
}

func TestPathsWalkError(t *testing.T) {
	testWalkError(t, new(Paths))
}

func TestPathsWalkSkip(t *testing.T) {
	testWalkSkip(t, new(Paths))
}

func TestPathsInspectSkip(t *testing.T) {
	testInspectSkip(t, new(Paths))
}
