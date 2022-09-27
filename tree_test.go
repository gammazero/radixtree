package radixtree

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestAddEnd(t *testing.T) {
	rt := New()
	rt.Put("tomato", "TOMATO")
	if len(rt.root.edges) != 1 {
		t.Fatal("root should have 1 child")
	}
	node := rt.root.getEdge('t')
	if node == nil {
		t.Fatal("root should have child at 't'")
	}
	if node.prefix != "omato" {
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
	if node.prefix != "om" {
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
	if node.prefix != "to" {
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

func TestAddFront(t *testing.T) {
	rt := New()
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
	if node.prefix != "om" {
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
	if node.prefix != "to" {
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

func TestAddBranch(t *testing.T) {
	rt := New()
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
	if node.prefix != "o" {
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
	if node3.prefix != "to" {
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
	if node2.prefix != "n" {
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

func TestAddBranchToBranch(t *testing.T) {
	rt := New()
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
	if node2.prefix != "g" {
		t.Fatal("expected prefix 'g', got: ", node2.prefix)
	}
	if node2.leaf == nil {
		t.Fatal("missing value at child")
	}
	if node2.leaf.value != "TAG" {
		t.Fatal("expected value 'TAG', got:", node2.leaf.value)
	}
}

func TestAddExisting(t *testing.T) {
	rt := New()
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

func TestDelete(t *testing.T) {
	rt := New()
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
	if node.prefix != "ato" {
		t.Log(dump(rt))
		t.Error("wrong prefix for compresses node: ", node.prefix)
	}

	// Test deleting key that does not exist
	if rt.Delete("xyz") {
		t.Error("expected false when deleting key 'xyz'")
	}
}

func TestBuildEdgeCases(t *testing.T) {
	tree := New()

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
	if node.prefix != "L1/L2" {
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
	if node.prefix != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if node.prefix != "L3" {
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
	if node.prefix != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if node.prefix != "L3" {
		t.Fatal("expected prefix '/L3', got ", node.prefix)
	}
	if node.leaf.value != 555 {
		t.Fatal("expected value of 555, got ", node.leaf.value)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if node.prefix != "L4" {
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
	if node.prefix != "L1/L2" {
		t.Fatal("expected prefix 'L2/L3', got ", node.prefix)
	}
	node = node.getEdge('/')
	if node == nil {
		t.Fatal("expected child at '/'")
	}
	if node.prefix != "L" {
		t.Fatal("expected prefix 'L', got ", node.prefix)
	}
	if node.leaf != nil {
		t.Fatal("expected nil value, got ", node.leaf.value)
	}
	if len(node.edges) != 2 {
		t.Fatal("expected 2 children, got ", len(node.edges))
	}

	t.Log(dump(tree))
	tree.Put("/L1/L2/L3/X", 999)
	t.Log(dump(tree))
}

func TestSimpleWalk(t *testing.T) {
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

func TestTree(t *testing.T) {
	tree := New()

	keys := []string{
		"bird",
		"/rat",
		"/bat",
		"/rats",
		"/ratatouille",
		"/rat/whiskey",
		"/rat/whiskers",
	}

	// check that keys do not exist
	for _, key := range keys {
		if val, ok := tree.Get(key); ok {
			t.Errorf("expected key %s to be missing, found value %v", key, val)
		}
	}

	// store keys
	for _, key := range keys {
		if isNew := tree.Put(key, "first"); !isNew {
			t.Errorf("expected key %s to be new", key)
		}
	}

	if tree.Len() != len(keys) {
		t.Fatalf("wrong length, expected %d, got %d", len(keys), tree.Len())
	}

	// put again, same keys new values
	for _, key := range keys {
		if isNew := tree.Put(key, strings.ToUpper(key)); isNew {
			t.Errorf("expected key %s to already have a value", key)
		}
	}

	if tree.Len() != len(keys) {
		t.Fatalf("wrong length, expected %d, got %d", len(keys), tree.Len())
	}

	// get
	for _, key := range keys {
		val, _ := tree.Get(key)
		if val.(string) != strings.ToUpper(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), val)
		}
	}

	var wvals []interface{}
	kvMap := map[string]interface{}{}
	walkFn := func(key string, value interface{}) bool {
		kvMap[key] = value
		wvals = append(wvals, value)
		return false
	}

	// walk path
	t.Log(dump(tree))
	key := "bad/key"
	tree.WalkPath(key, walkFn)
	if len(kvMap) != 0 {
		t.Error("should not have returned values, got ", kvMap)
	}
	lastKey := keys[len(keys)-1]
	var expectVals []interface{}
	for _, key := range keys {
		// If key is a prefix of lastKey, then expect value.
		if strings.HasPrefix(lastKey, key) {
			expectVals = append(expectVals, strings.ToUpper(key))
		}
	}
	kvMap = map[string]interface{}{}
	wvals = nil
	tree.WalkPath(lastKey, walkFn)
	if kvMap[lastKey] == nil {
		t.Fatalf("expected value for %s", lastKey)
	}
	if len(kvMap) != len(expectVals) {
		t.Errorf("expected %d values, got %d", len(expectVals), len(kvMap))
	} else {
		for i, expect := range expectVals {
			var found bool
			for _, v := range kvMap {
				if v == expect {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("did not find expected value %v", expectVals[i])
			}
			if wvals[i] != expect {
				t.Errorf("did not find expected value %v at pos %d, got %v", expect, i, wvals[i])
			}
		}
	}

	// delete, expect Delete to return true indicating a node was nil'd
	t.Log("Before")
	t.Log(dump(tree))
	for _, key := range keys {
		if deleted := tree.Delete(key); !deleted {
			t.Errorf("expected key %s to be deleted", key)
		}
	}
	t.Log("After")
	t.Log(dump(tree))

	if tree.Len() != 0 {
		t.Error("expected Len() to return 0 after all keys deleted")
	}

	// expect Delete to return false bc all nodes deleted
	for _, key := range keys {
		if deleted := tree.Delete(key); deleted {
			t.Errorf("expected key %s to be cleaned by delete", key)
		}
	}

	// get deleted keys
	for _, key := range keys {
		if val, ok := tree.Get(key); ok {
			t.Errorf("expected key %s to be deleted, got value %v", key, val)
		}
	}
}

func TestNilGet(t *testing.T) {
	tree := New()

	tree.Put("/rat", 1)
	tree.Put("/ratatattat", 2)
	tree.Put("/ratatouille", nil)

	val, ok := tree.Get("/ratatouille")
	if !ok {
		t.Fatal("missing value")
	}
	if val != nil {
		t.Errorf("expected nil value")
	}

	for _, key := range []string{"/", "/r", "/ra", "/other"} {
		val, ok = tree.Get(key)
		if ok {
			t.Fatal("should not indicate value present")
		}
		if val != nil {
			t.Errorf("expected key %s to have nil value, got %v", key, val)
		}
	}
}

func TestRoot(t *testing.T) {
	tree := New()

	val, ok := tree.Get("")
	if ok {
		t.Errorf("expected nil key to be missing, found value %v", val)
	}
	if !tree.Put("", "hello") {
		t.Error("expected nil key to be new")
	}
	testVal := "world"
	if tree.Put("", testVal) {
		t.Error("expected nil key to already have a value")
	}
	val, ok = tree.Get("")
	if !ok {
		t.Fatal("missing expected value")
	}
	if val != testVal {
		t.Errorf("expected nil key to have value %v, got %v", testVal, val)
	}
	if !tree.Delete("") {
		t.Error("expected nil key to be deleted")
	}
	if val, ok := tree.Get(""); ok {
		t.Errorf("expected nil key to be deleted, got value %v", val)
	}
	if tree.Delete("") {
		t.Error("expected nil key to be already deleted")
	}
}

func TestWalk(t *testing.T) {
	tree := New()

	keys := []string{
		"bird",
		"rat",
		"bat",
		"rats",
		"ratatouille",
		"rat/whis/key",           // visited by 2
		"rat/whis/kers",          // visited by 2
		"rat/whis/per/er",        // visited by 2, 3
		"rat/winks/wisely/once",  // visited by 5, 6
		"rat/winks/wisely/x/y/z", // visited by 5, 6, 7
		"rat/winks/wryly",        // visited by 5
	}

	notKeys := []string{
		"birds",                 // 0
		"rat/whiz",              // 1
		"rat/whis",              // 2
		"rat/whis/per",          // 3
		"rat/whiskey/shot",      // 4
		"rat/winks",             // 5
		"rat/winks/wisely",      // 6
		"rat/winks/wisely/x/y",  // 7
		"rat/winks/wisely/x/w",  // 8
		"rat/winks/wisely/only", // 9
	}

	visited := make(map[string]int, len(keys))

	for _, key := range keys {
		if isNew := tree.Put(key, strings.ToUpper(key)); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	var err error
	walkFn := func(key string, value interface{}) bool {
		// value for each walked key is correct
		if value != strings.ToUpper(key) {
			err = fmt.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), value)
			return true
		}
		count := visited[key]
		visited[key] = count + 1
		return false
	}

	for _, notKey := range notKeys {
		tree.Walk(notKey, walkFn)
		if err != nil {
			t.Error(err)
		}

	}
	t.Log(dump(tree))

	for _, notKey := range notKeys {
		_, ok := visited[notKey]
		if ok {
			t.Fatalf("%s should not have been visited", notKey)
		}
	}

	expectCounts := map[string]int{
		"rat/whis/key":           1, // visited by 2
		"rat/whis/kers":          1, // visited by 2
		"rat/whis/per/er":        2, // visited by 2, 3
		"rat/winks/wisely/once":  2, // visited by 5, 6
		"rat/winks/wisely/x/y/z": 3, // visited by 5, 6, 7
		"rat/winks/wryly":        1, // visited by 5
	}

	for key, count := range visited {
		expected := expectCounts[key]
		if count != expected {
			t.Fatalf("expected %s to have visited count of %d, got %d", key, expected, count)
		}
	}

	visited = make(map[string]int, len(keys))

	// Walk from root
	tree.Walk("", walkFn)
	if err != nil {
		t.Error(err)
	}

	if len(visited) != len(keys) {
		t.Error("wrong number of iterm iterated")
	}

	// each key/value visited exactly once
	for key, visitedCount := range visited {
		if visitedCount != 1 {
			t.Errorf("expected key %s to be visited exactly once, got %v", key, visitedCount)
		}
	}

	visited = make(map[string]int, len(keys))

	var iterCopy *Iterator[any]
	iter := tree.NewIterator()
	for {
		key, val, done := iter.Next()
		if done {
			break
		}
		t.Log("Iterated key:", key)
		if walkFn(key, val) {
			if err != nil {
				t.Fatal(err)
			}
			break
		}
		if key == "rat/winks/wisely/once" {
			iterCopy = iter.Copy()
		}
	}

	if len(visited) != len(keys) {
		t.Error("wrong number of iterm iterated")
	}

	// each key/value visited exactly once
	for key, visitedCount := range visited {
		if visitedCount != 1 {
			t.Errorf("expected key %s to be visited exactly once, got %v", key, visitedCount)
		}
	}

	visited = make(map[string]int, 4)

	for {
		key, val, done := iterCopy.Next()
		if done {
			break
		}
		t.Log("Iterated key:", key)
		if walkFn(key, val) {
			if err != nil {
				t.Fatal(err)
			}
			break
		}
	}

	if len(visited) != 4 {
		t.Error("Wrong iteration count of copied iterator:", len(visited))
	}

	// each key/value visited exactly once
	for key, visitedCount := range visited {
		if visitedCount != 1 {
			t.Errorf("expected key %s to be visited exactly once, got %v", key, visitedCount)
		}
	}

	visited = make(map[string]int, len(keys))

	tree.Walk("rat", walkFn)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	if visited[keys[0]] != 0 {
		t.Error(keys[0], "should not have been visited")
	}
	if visited[keys[2]] != 0 {
		t.Error(keys[2], "should not have been visited")
	}
	if visited[keys[3]] != 1 {
		t.Error(keys[3], "should have been visited")
	}
	if visited[keys[5]] != 1 {
		t.Error(keys[5], "should have been visited")
	}
	if visited[keys[6]] != 1 {
		t.Error(keys[6], "should have been visited")
	}
	if visited[keys[7]] != 1 {
		t.Error(keys[7], "should have been visited")
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}

	tree.Walk("rat/whis/kers", walkFn)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	for _, key := range keys {
		if key == "rat/whis/kers" {
			if visited[key] != 1 {
				t.Error(key, "should have been visited")
			}
			continue
		}
		if visited[key] != 0 {
			t.Error(key, "should not have been visited")
		}
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}

	testKey := "rat/winks/wryly/once"
	keys = append(keys, testKey)
	tree.Put(testKey, strings.ToUpper(testKey))

	walkPFn := func(key string, value interface{}) bool {
		// value for each walked key is correct
		v := strings.ToUpper(key)
		if value != v {
			err = fmt.Errorf("expected key %s to have value %v, got %v", key, v, value)
			return true
		}
		visited[key]++
		return false
	}

	tree.WalkPath(testKey, walkPFn)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	err = checkVisited(visited, "rat", "rat/winks/wryly", testKey)
	if err != nil {
		t.Error(err)
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}
	tree.WalkPath(testKey, func(key string, value interface{}) bool {
		pfx := "rat/winks/wryly"
		if strings.HasPrefix(key, pfx) && len(key) > len(pfx) {
			return false
		}
		visited[key]++
		return false
	})
	err = checkVisited(visited, "rat", "rat/winks/wryly")
	if err != nil {
		t.Error(err)
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}
	tree.WalkPath(testKey, func(key string, value interface{}) bool {
		visited[key]++
		if key == "rat/winks/wryly" {
			err = fmt.Errorf("error at key %s", key)
			return true
		}
		return false
	})
	if err == nil {
		t.Errorf("expected error")
	}
	err = checkVisited(visited, "rat", "rat/winks/wryly")
	if err != nil {
		t.Error(err)
	}

	var foundRoot bool
	tree.Put("", "ROOT")
	tree.WalkPath(testKey, func(key string, value interface{}) bool {
		if key == "" && value == "ROOT" {
			foundRoot = true
		}
		return false
	})
	if !foundRoot {
		t.Error("did not find root")
	}

	for _, k := range keys {
		visited[k] = 0
	}

	tree.WalkPath(testKey, func(key string, value interface{}) bool {
		if key == "" && value == "ROOT" {
			return false
		}
		return true
	})
	for k, count := range visited {
		if count != 0 {
			t.Error("should not have visited ", k)
		}
	}

	tree.WalkPath(testKey, func(key string, value interface{}) bool {
		if key == "" && value == "ROOT" {
			err = errors.New("error at root")
			return true
		}
		return false
	})
	if err == nil {
		t.Errorf("expected error")
	}
	for k, count := range visited {
		if count != 0 {
			t.Error("should not have visited ", k)
		}
	}

	var lastKey string
	tree.WalkPath("rat/winks/wisely/x/y/z/w", func(key string, value interface{}) bool {
		lastKey = key
		return false
	})
	if lastKey != "rat/winks/wisely/x/y/z" {
		t.Error("did not get expected last key")
	}
}

func checkVisited(visited map[string]int, expectVisited ...string) error {
	for _, key := range expectVisited {
		if visited[key] != 1 {
			return fmt.Errorf("%s should have been visited once", key)
		}
		delete(visited, key)
	}
	for key, count := range visited {
		if count != 0 {
			return fmt.Errorf("%s should not have been visited", key)
		}
	}
	for _, key := range expectVisited {
		visited[key] = 1
	}

	return nil
}

func TestWalkStop(t *testing.T) {
	tree := New()

	table := []struct {
		key string
		val int
	}{
		{"/L1/L2", 1},
		{"/L1/L2/L3A", 2},
		{"/L1/L2/L3B/L4", 999},
		{"/L1/L2/L3C", 4},
		{"/L1/L2/L3", 5},
	}

	for i := range table {
		tree.Put(table[i].key, table[i].val)
	}

	walkErr := errors.New("walk error")
	var walked int
	var err error
	walkFn := func(k string, value interface{}) bool {
		if value == 999 {
			err = walkErr
			return true
		}
		walked++
		return false
	}
	tree.Walk("", walkFn)
	if err != walkErr {
		t.Errorf("expected error %v, got %v", walkErr, err)
	}
	if len(table) == walked {
		t.Errorf("expected nodes walked < %d, got %d", len(table), walked)
	}
}

func TestInspectStop(t *testing.T) {
	tree := New()

	table := []struct {
		key string
		val int
	}{
		{"/L1/L2/L3/X", 999},
		{"/L1/L2", 1},
		{"/L1/L2/L3", 555},
		{"/L1/L2/L3/L4", 999},
		{"/L1/L2/L/C", 3},
	}

	for i := range table {
		tree.Put(table[i].key, table[i].val)
	}
	var keys []string
	inspectFn := func(link, prefix, key string, depth, children int, hasValue bool, value interface{}) bool {
		if !hasValue {
			// Do not count internal nodes
			return false
		}
		keys = append(keys, key)
		switch value {
		case 555:
			// Stop inspect
			return true
		case 999:
			t.Fatal("should not get here")
		}
		return false
	}
	tree.Inspect(inspectFn)
	if len(keys) != len(table)-2 {
		t.Errorf("expected nodes walked to be %d, got %d: %v", len(table)-2, len(keys), keys)
	}
}

func TestGetAfterDelete(t *testing.T) {
	tree := New()

	keys := []string{
		"bird",
		"rat",
	}

	// store keys
	for _, key := range keys {
		tree.Put(key, strings.ToUpper(key))
	}

	t.Log("Before")
	t.Log(dump(tree))

	if !tree.Delete("bird") {
		t.Fatal("should have deleted bird")
	}
	t.Log("After")
	t.Log(dump(tree))

	_, ok := tree.Get("rat")
	if !ok {
		t.Fatal("Did not get rat")
	}
}

func TestStringConvert(t *testing.T) {
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

// Use the Inspect functionality to create a function to dump the tree.
func dump(tree *Tree[any]) string {
	var b strings.Builder
	tree.Inspect(func(link, prefix, key string, depth, children int, hasValue bool, value interface{}) bool {
		for ; depth > 0; depth-- {
			b.WriteString("  ")
		}
		if hasValue {
			b.WriteString(fmt.Sprintf("%s-> (%q, [%s: %q]) children: %d\n", link, prefix, key, value, children))
		} else {
			b.WriteString(fmt.Sprintf("%s-> (%q, [%s])] children: %d\n", link, prefix, key, children))
		}
		return false
	})
	return b.String()
}
