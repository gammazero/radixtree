package radixtree

import (
	"fmt"
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
	node := tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if strings.Join(node.prefix, "") != "L2" {
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
	node = tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if strings.Join(node.prefix, "") != "L2" {
		t.Fatal("wrong prefix:", node.prefix)
	}
	if node.value != 1 {
		t.Fatal("expected value 1, got ", node.value)
	}
	if len(node.children) != 1 {
		t.Fatal("expected 1 child")
	}
	node = node.children["L3A"]
	if node == nil {
		t.Fatal("expected child at 'L3A'")
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
	node := tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children")
	}
	node2 := node.children["L3B"]
	if node2 == nil {
		t.Fatal("expected child at 'L3B'")
	}
	if strings.Join(node2.prefix, "") != "L4" {
		t.Fatal("wrong prefix:", node2.prefix)
	}
	if node2.value != 3 {
		t.Fatal("expected value 3, got ", node2.value)
	}
	node2 = node.children["L3A"]
	if node2 == nil {
		t.Fatal("expected child at 'L3A'")
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
	node := tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
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
	node2 := node.children["L2B"]
	if node2 == nil {
		t.Fatal("Expected child at 'L2B'")
	}
	if strings.Join(node2.prefix, "") != "L3C" {
		t.Fatal("expected prefix 'L3C', got ", node2.prefix)
	}
	if node2.value != 4 {
		t.Fatal("expected value of 4, got ", node2.value)
	}
	node2 = node.children["L2"]
	if node2 == nil {
		t.Fatal("expected child at 'L2'")
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
	node := tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
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
	node := tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", node.children)
	}
	node = node.children["L2"]
	if node == nil {
		t.Fatal("expected child at 'L2'")
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
	node = tree.children["L1"]
	if node == nil {
		t.Fatal("expected child at 'L1'")
	}
	if len(node.children) != 2 {
		t.Fatal("expected 2 children, got ", node.children)
	}
	node = node.children["L2"]
	if node == nil {
		t.Fatal("expected child at 'L2'")
	}
	if len(node.children) != 0 {
		t.Fatal("expected 0 children, got ", node.children)
	}
	if strings.Join(node.prefix, "") != "L3A" {
		t.Fatal("expected prefix 'L3A', got ", node.prefix)
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
	node = tree.children["L1"]
	if _, ok := node.children["L2B"]; ok {
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
	node = tree.children["L1"]
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "/") != "L2/L3A" {
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
	node = tree.children["L1"]
	if node == nil {
		t.Fatal("expected node at \"L1\"")
	}
	if strings.Join(node.prefix, "/") != "L2/L3A" {
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
	if iter.Next("x") {
		t.Fatal("x should not have advanced iterator")
	}
	if !iter.Next("L1") {
		t.Fatal("L1 should have advanced iterator")
	}
	if iter.Value() != nil {
		t.Fatal("should not have value at /L1")
	}
	if !iter.Next("L2") {
		t.Fatal("L2 should have advanced iterator")
	}
	if iter.Value() != 1 {
		t.Fatal("expected value 1 at L2, got ", iter.Value())
	}
	if iter.Next("L4") {
		t.Fatal("L4 should not have advanced iterator")
	}

	// branch iterator
	iterB := iter.Copy()
	if !iterB.Next("L3B") {
		t.Fatal("L3B should have advanced iterator")
	}
	if iterB.Value() != nil {
		t.Fatal("should not have value at L3B")
	}
	if !iterB.Next("L4") {
		t.Fatal("L4 should have advanced iterator")
	}
	if iterB.Value() != 3 {
		t.Fatal("expected value 3 at L4, got ", iterB.Value())
	}
	if iterB.Next("L4") {
		t.Fatal("L4 should not have advanced iterator")
	}

	if !iter.Next("L3A") {
		t.Fatal("L3A should have advanced iterator")
	}
	if iter.Value() != 2 {
		t.Fatal("expected value 2 at L3A, got ", iter.Value())
	}
	if iter.Next("L3B") {
		t.Fatal("L3B should not have advanced iterator")
	}

}

func TestSimplePathWalk(t *testing.T) {
	rt := new(Paths)
	rt.Put("tom/ato", "TOMATO")
	rt.Put("tom", "TOM")
	rt.Put("torn/ad/o", "TORNADO")

	count := 0
	err := rt.Walk("tom/ato", func(key fmt.Stringer, value interface{}) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected to visit 1 key, visited %d", count)
	}

	count = 0
	err = rt.Walk("tom", func(key fmt.Stringer, value interface{}) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Errorf("expected to visit 3 keys, visited %d", count)
	}

	count = 0
	err = rt.Walk("tomx", func(key fmt.Stringer, value interface{}) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected to visit 0 keys, visited %d", count)
	}

	count = 0
	err = rt.Walk("torn", func(key fmt.Stringer, value interface{}) error {
		count++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected to visit 1 key, visited %d", count)
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

func TestPathsInspectError(t *testing.T) {
	testInspectError(t, new(Paths))
}

/*
func TestPathAjg(t *testing.T) {
	pt := new(Paths)
	pt.Put("home/abc", "my home directory")
	pt.Put("home/abc/a1.txt", "some text")
	pt.Put("home/abc/Documents", "my documents")
	pt.Put("home/abc/Documents/pic.png", "cat pic")
	pt.Put("home/abc/Documents/stuff.pdf", "story")
	pt.Put("home/abc/Documents/stuff.pdf", "story")
	pt.Put("home", "the place")

	fmt.Println("====== WalkPath =====")
	pt.WalkPath("home/abc/Documents/pic.png", func(key string, value interface{}) error {
		fmt.Println(key, "=>", value)
		return nil
	})
	fmt.Println()
	fmt.Println("====== Walk =====")
	pt.Walk("", func(key fmt.Stringer, value interface{}) error {
		fmt.Println(key, "=>", value)
		return nil
	})
	fmt.Println()

	piter := pt.NewIterator()
	//parts := strings.Split("home/abc/Documents/pic.png", "/")
	parts := []string{"home", "abc", "/Documents", "/pic.png"}
	fmt.Println(parts)
	for i, p := range parts {
		if !piter.Next(p) {
			break
		}
		value := piter.Value()
		if value != nil {
			fmt.Println(parts[i], "=>", value)
		}
	}

	fmt.Println(dump(pt))
}

func TestPathAjg2(t *testing.T) {
	pt := new(Paths)
	keys := []string{
		"bird",
		"/rat",
		"/bat",
		"/rats",
		"/ratatouille",
		"/rat/whiskey",
		"/rat/whiskers",
	}
	// put again, same keys new values
	for _, key := range keys {
		if isNew := pt.Put(key, strings.ToUpper(key)); !isNew {
			t.Fatalf("expected key %s to not already have a value", key)
		}
	}

	// get
	for _, key := range keys {
		val := pt.Get(key)
		fmt.Println("getting key:", key)
		if val == nil || val.(string) != strings.ToUpper(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), val)
		}
	}
	fmt.Println(dump(pt))
}
*/

/*
func testWalkPaths(t *testing.T, tree rtree) {
	keys := []string{
		"bird",
		"/rat",
		"/bat",
		"/rats",
		"/ratatouille",
		"/rat//whis/key",            // visited by 2
		"/rat/whis/kers",           // visited by 2
		"/rat/whis/per/er",         // visited by 2, 3
		"/rat/winks/wisely/once",  // visited by 5, 6
		"/rat/winks/wisely/x/y/z", // visited by 5, 6, 7
		"/rat/winks/wryly",        // visited by 5
	}

	notKeys := []string{
		"birds",                  // 0
		"/rat/whiz",              // 1
		"/rat/whis",              // 2
		"/rat/whis/per",           // 3
		"/rat/whis/key/shot",      // 4
		"/rat/winks",             // 5
		"/rat/winks/wisely",      // 6
		"/rat/winks/wisely/x/y",  // 7
		"/rat/winks/wisely/x/w",  // 8
		"/rat/winks/wisely/only", // 9
	}

	visited := make(map[string]int, len(keys))

	keyToValue := func(key string) string {
		f := func(r rune) rune {
			if r == '/' {
				return '-'
			}
			return r
		}
		return strings.Map(f, strings.ToUpper(strings.Trim(key, "/")))
	}

	for _, key := range keys {
		if isNew := tree.Put(key, keyToValue(key)); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	walkFn := func(k fmt.Stringer, value interface{}) error {
		// value for each walked key is correct
		key := k.String()
		if value != keyToValue(key) {
			return fmt.Errorf("expected key %s to have value %v, got %v", key, keyToValue(key), value)
		}
		fmt.Println("+++ visited:", key)
		count := visited[key]
		visited[key] = count + 1
		return nil
	}

	var err error
	for _, notKey := range notKeys {
		if err = tree.Walk(notKey, walkFn); err != nil {
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
		"/rat/whiskey":            1, // visited by 2
		"/rat/whiskers":           1, // visited by 2
		"/rat/whisperer":          2, // visited by 2, 3
		"/rat/winks/wisely/once":  2, // visited by 5, 6
		"/rat/winks/wisely/x/y/z": 3, // visited by 5, 6, 7
		"/rat/winks/wryly":        1, // visited by 5
	}

	for key, count := range visited {
		expected := expectCounts[key]
		if count != expected {
			t.Fatalf("expected %s to have visited count of %d, got %d", key, expected, count)
		}
	}

	visited = make(map[string]int, len(keys))

	// Walk from root
	if err = tree.Walk("", walkFn); err != nil {
		t.Error(err)
	}

	// each key/value visited exactly once
	for key, visitedCount := range visited {
		if visitedCount != 1 {
			t.Errorf("expected key %s to be visited exactly once, got %v", key, visitedCount)
		}
	}

	visited = make(map[string]int, len(keys))

	if err := tree.Walk("/rat", walkFn); err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	if visited[keys[0]] != 0 {
		t.Error(keys[0], " should not have been visited")
	}
	if visited[keys[2]] != 0 {
		t.Error(keys[2], " should not have been visited")
	}
	// Do not test /rats since that is visited by Runes but not Paths
	if visited[keys[5]] != 1 {
		t.Error(keys[5], " should have been visited")
	}
	if visited[keys[6]] != 1 {
		t.Error(keys[6], " should have been visited")
	}
	if visited[keys[7]] != 1 {
		t.Error(keys[7], " should have been visited")
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}

	err = tree.Walk("/rat/whiskers", walkFn)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	for i, key := range keys {
		if i == 6 {
			continue
		}
		if visited[key] != 0 {
			t.Error(key, " should not have been visited")
		}
	}
	// Do not test /rats since that is visited by Runes but not Paths
	if visited[keys[6]] != 1 {
		t.Error(keys[6], " should have been visited")
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}

	testKey := "/rat/winks/wryly/once"
	keys = append(keys, testKey)
	tree.Put(testKey, keyToValue(testKey))

	walkPFn := func(key string, value interface{}) error {
		// value for each walked key is correct
		if value != keyToValue(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, keyToValue(key), value)
		}
		visited[key]++
		return nil
	}

	err = tree.WalkPath(testKey, walkPFn)
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	err = checkVisited(visited, "/rat", "/rat/winks/wryly", testKey)
	if err != nil {
		t.Error(err)
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}
	err = tree.WalkPath(testKey, func(key string, value interface{}) error {
		visited[key]++
		if key == "/rat/winks/wryly" {
			return Skip
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	err = checkVisited(visited, "/rat", "/rat/winks/wryly")
	if err != nil {
		t.Error(err)
	}

	// Reset visited counts
	for _, k := range keys {
		visited[k] = 0
	}
	err = tree.WalkPath(testKey, func(key string, value interface{}) error {
		visited[key]++
		if key == "/rat/winks/wryly" {
			return fmt.Errorf("error at key %s", key)
		}
		return nil
	})
	if err == nil {
		t.Errorf("expected error")
	}
	err = checkVisited(visited, "/rat", "/rat/winks/wryly")
	if err != nil {
		t.Error(err)
	}

	var foundRoot bool
	tree.Put("", "ROOT")
	err = tree.WalkPath(testKey, func(key string, value interface{}) error {
		if key == "" && value == "ROOT" {
			foundRoot = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if !foundRoot {
		t.Error("did not find root")
	}

	for _, k := range keys {
		visited[k] = 0
	}

	err = tree.WalkPath(testKey, func(key string, value interface{}) error {
		if key == "" && value == "ROOT" {
			return Skip
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	for k, count := range visited {
		if count != 0 {
			t.Error("should not have visited ", k)
		}
	}

	err = tree.WalkPath(testKey, func(key string, value interface{}) error {
		if key == "" && value == "ROOT" {
			return errors.New("error at root")
		}
		return nil
	})
	if err == nil {
		t.Errorf("expected error")
	}
	for k, count := range visited {
		if count != 0 {
			t.Error("should not have visited ", k)
		}
	}
}
*/
