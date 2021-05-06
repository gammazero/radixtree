package radixtree

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// rtree is an interface common to all radix tree types, used for test
type rtree interface {
	Get(key string) interface{}
	Put(key string, value interface{}) bool
	Delete(key string) bool
	Walk(key string, walkFn WalkFunc) error
	WalkPath(key string, walkFn WalkPathFunc) error
	Inspect(inspectFn InspectFunc) error
}

// Use the Inspect functionality to create a function to dump the tree.
func dump(tree rtree) string {
	var b strings.Builder
	err := tree.Inspect(func(link, prefix, key string, depth, children int, value interface{}) error {
		for ; depth > 0; depth-- {
			b.WriteString("  ")
		}
		b.WriteString(fmt.Sprintf("%s-> (%q, %v) key: %q children: %d\n", link, prefix, value, key, children))
		return nil
	})
	if err != nil {
		b.WriteString(fmt.Sprintln("ERROR:", err))
	}
	return b.String()
}

func testRadixTree(t *testing.T, tree rtree) {
	keys := []string{
		"bird",
		"/rat",
		"/bat",
		"/rats",
		"/ratatouille",
		"/rat/whiskey",
		"/rat/whiskers",
	}

	// get keys that do not exist
	for _, key := range keys {
		if val := tree.Get(key); val != nil {
			t.Errorf("expected key %s to be missing, found value %v", key, val)
		}
	}

	// store keys
	for _, key := range keys {
		if isNew := tree.Put(key, "first"); !isNew {
			t.Errorf("expected key %s to be new", key)
		}
	}

	// put again, same keys new values
	for _, key := range keys {
		if isNew := tree.Put(key, strings.ToUpper(key)); isNew {
			t.Errorf("expected key %s to already have a value", key)
		}
	}

	// get
	for _, key := range keys {
		val := tree.Get(key)
		if val.(string) != strings.ToUpper(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), val)
		}
	}

	var wvals []interface{}
	kvMap := map[string]interface{}{}
	walkFn := func(key string, value interface{}) error {
		kvMap[key] = value
		wvals = append(wvals, value)
		return nil
	}

	// walk path
	t.Log(dump(tree))
	key := "bad/key"
	err := tree.WalkPath(key, walkFn)
	if err != nil {
		t.Fatal(err)
	}
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
	err = tree.WalkPath(lastKey, walkFn)
	if err != nil {
		t.Fatal(err)
	}
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
	for _, key := range keys {
		if deleted := tree.Delete(key); !deleted {
			t.Errorf("expected key %s to be deleted", key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, key := range keys {
		if deleted := tree.Delete(string(key)); deleted {
			t.Errorf("expected key %s to be cleaned by delete", string(key[0]))
		}
	}

	// get deleted keys
	for _, key := range keys {
		if val := tree.Get(key); val != nil {
			t.Errorf("expected key %s to be deleted, got value %v", key, val)
		}
	}
}

func testNilGet(t *testing.T, tree rtree) {
	tree.Put("/rat", 1)
	tree.Put("/ratatattat", 2)
	tree.Put("/ratatouille", nil)

	for _, key := range []string{"/", "/r", "/ra", "/ratatouille", "/other"} {
		if val := tree.Get(key); val != nil {
			t.Errorf("expected key %s to have nil value, got %v", key, val)
		}
	}
}

func testRoot(t *testing.T, tree rtree) {
	if val := tree.Get(""); val != nil {
		t.Errorf("expected key '' to be missing, found value %v", val)
	}
	if !tree.Put("", "hello") {
		t.Error("expected key \"\" to be new")
	}
	testVal := "world"
	if tree.Put("", testVal) {
		t.Error("expected key \"\" to already have a value")
	}
	if val := tree.Get(""); val != testVal {
		t.Errorf("expected key \"\" to have value %v, got %v", testVal, val)
	}
	if !tree.Delete("") {
		t.Error("expected key \"\" to be deleted")
	}
	if val := tree.Get(""); val != nil {
		t.Errorf("expected key \"\" to be deleted, got value %v", val)
	}
	if tree.Delete("") {
		t.Error("expected key \"\" to be already deleted")
	}
}

func testWalk(t *testing.T, tree rtree) {
	keys := []string{
		"bird",
		"/rat",
		"/bat",
		"/rats",
		"/ratatouille",
		"/rat/whiskey",            // visited by 2
		"/rat/whiskers",           // visited by 2
		"/rat/whisperer",          // visited by 2, 3
		"/rat/winks/wisely/once",  // visited by 5, 6
		"/rat/winks/wisely/x/y/z", // visited by 5, 6, 7
		"/rat/winks/wryly",        // visited by 5
	}

	notKeys := []string{
		"birds",                  // 0
		"/rat/whiz",              // 1
		"/rat/whis",              // 2
		"/rat/whisper",           // 3
		"/rat/whiskey/shot",      // 4
		"/rat/winks",             // 5
		"/rat/winks/wisely",      // 6
		"/rat/winks/wisely/x/y",  // 7
		"/rat/winks/wisely/x/w",  // 8
		"/rat/winks/wisely/only", // 9
	}

	visited := make(map[string]int, len(keys))

	for _, key := range keys {
		if isNew := tree.Put(key, strings.ToUpper(key)); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	walkFn := func(k KeyStringer, value interface{}) error {
		// value for each walked key is correct
		key := k.String()
		if value != strings.ToUpper(key) {
			return fmt.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), value)
			//t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), value)
		}
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
	tree.Put(testKey, strings.ToUpper(testKey))

	walkPFn := func(key string, value interface{}) error {
		// value for each walked key is correct
		if value != strings.ToUpper(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), value)
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

func testWalkError(t *testing.T, tree rtree) {
	table := map[string]int{
		"/L1/L2":        1,
		"/L1/L2/L3A":    2,
		"/L1/L2/L3B/L4": 999,
		"/L1/L2/L3C":    4,
		"/L1/L2/L3":     5,
	}

	for key, value := range table {
		tree.Put(key, value)
	}

	walkErr := errors.New("walk error")
	var walked int
	walkFn := func(k KeyStringer, value interface{}) error {
		if value == 999 {
			return walkErr
		}
		walked++
		return nil
	}
	if err := tree.Walk("", walkFn); err != walkErr {
		t.Errorf("expected error %v, got %v", walkErr, err)
	}
	if len(table) == walked {
		t.Errorf("expected nodes walked < %d, got %d", len(table), walked)
	}
}

func testWalkSkip(t *testing.T, tree rtree) {
	table := map[string]int{
		"/L1/L2":       1,
		"/L1/L2/L3":    555,
		"/L1/L2/L3/L4": 999,
		"/L1/L2/L/C":   3,
		"/L1/L2/L3/X":  999,
	}

	for key, value := range table {
		tree.Put(key, value)
		t.Log(dump(tree))
	}
	var walked int
	walkFn := func(k KeyStringer, value interface{}) error {
		switch value {
		case 555:
			return Skip
		case 999:
			t.Fatal("should not get here")
		}
		walked++
		return nil
	}
	if err := tree.Walk("", walkFn); err != nil {
		t.Error(err)
	}
	if walked != len(table)-3 {
		t.Errorf("expected nodes walked to be %d, got %d", len(table)-3, walked)
	}
}

func testInspectSkip(t *testing.T, tree rtree) {
	table := map[string]int{
		"/L1/L2":       1,
		"/L1/L2/L3":    555,
		"/L1/L2/L3/L4": 999,
		"/L1/L2/L/C":   3,
		"/L1/L2/L3/X":  999,
	}

	for key, value := range table {
		tree.Put(key, value)
		t.Log(dump(tree))
	}
	var keys []string
	inspectFn := func(link, prefix, key string, depth, children int, value interface{}) error {
		if value == nil {
			// Do not count internal nodes
			return nil
		}
		keys = append(keys, key)
		switch value {
		case 555:
			// SKip all this node's children
			return Skip
		case 999:
			t.Fatal("should not get here")
		}
		return nil
	}
	if err := tree.Inspect(inspectFn); err != nil {
		t.Error(err)
	}
	if len(keys) != len(table)-2 {
		t.Errorf("expected nodes walked to be %d, got %d: %v", len(table)-2, len(keys), keys)
	}
}
