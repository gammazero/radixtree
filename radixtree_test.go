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
	GetPath(key string) ([]interface{}, bool)
	Put(key string, value interface{}) bool
	Delete(key string) bool
	Walk(startKey string, walkFn WalkFunc) error
	Inspect(inspectFn InspectFunc) error
}

// Use the Inspect functionality to create a function to dump the tree.
func dump(tree rtree) string {
	var b strings.Builder
	tree.Inspect(func(link, prefix, key string, depth, children int, value interface{}) error {
		for ; depth > 0; depth-- {
			b.WriteString("  ")
		}
		b.WriteString(fmt.Sprintf("%s-> (%q, %v) key: %q children: %d\n", link, prefix, value, key, children))
		return nil
	})
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

	// get path
	t.Log(dump(tree))
	vals, ok := tree.GetPath("bad/key")
	if ok {
		t.Error("should not have found node")
	}
	if len(vals) != 0 {
		t.Error("should not have returned values, got ", vals)
	}
	lastKey := keys[len(keys)-1]
	var expectVals []interface{}
	for _, key := range keys {
		// If key is a prefix of lastKey, then expect value.
		if strings.HasPrefix(lastKey, key) {
			expectVals = append(expectVals, strings.ToUpper(key))
		}
	}
	vals, ok = tree.GetPath(lastKey)
	if !ok {
		t.Fatalf("expected value for %s", lastKey)
	}
	if len(vals) != len(expectVals) {
		t.Errorf("expected %d values, got %d", len(expectVals), len(vals))
	} else {
		for i := range expectVals {
			if vals[i] != expectVals[i] {
				t.Errorf("expected value %v at position %d, got %v", expectVals[i], i, vals[i])
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
		"/rat/whiskey",
		"/rat/whiskers",
		"/rat/whisperer",
		"/rat/winks/wisely/once",
		"/rat/winks/wisely/x/y/z",
		"/rat/winks/wryly",
	}
	// key -> times visited
	visited := make(map[string]int, len(keys))
	for _, key := range keys {
		visited[key] = 0
	}

	for _, key := range keys {
		if isNew := tree.Put(key, strings.ToUpper(key)); !isNew {
			t.Errorf("expected key %s to be missing", key)
		}
	}

	walkFn := func(key string, value interface{}) error {
		// value for each walked key is correct
		if value != strings.ToUpper(key) {
			t.Errorf("expected key %s to have value %v, got %v", key, strings.ToUpper(key), value)
		}
		visited[key]++
		return nil
	}
	// Walk key that is not stored
	err := tree.Walk("/rat/whiz", walkFn)
	if err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/whis", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/whisper", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/whiskey/shot", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/winks", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/winks/wisely", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/winks/wisely/x/y", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/winks/wisely/x/w", walkFn); err != nil {
		t.Error(err)
	}
	if err = tree.Walk("/rat/winks/wisely/only", walkFn); err != nil {
		t.Error(err)
	}
	t.Log(dump(tree))

	for key, visitedCount := range visited {
		if visitedCount != 0 {
			t.Log(dump(tree))
			t.Fatalf("expected key %s to not be visited", key)
		}
	}

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

	// Reset visited counts
	for _, key := range keys {
		visited[key] = 0
	}

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
	for _, key := range keys {
		visited[key] = 0
	}

	if err := tree.Walk("/rat/whiskers", walkFn); err != nil {
		t.Errorf("expected error nil, got %v", err)
	}
	for i, key := range keys {
		if i == 6 {
			continue
		}
		if visited[key] != 0 {
			t.Log(dump(tree))
			t.Error(key, " should not have been visited")
		}
	}
	// Do not test /rats since that is visited by Runes but not Paths
	if visited[keys[6]] != 1 {
		t.Log(dump(tree))
		t.Error(keys[6], " should have been visited")
	}
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
	walkFn := func(key string, value interface{}) error {
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
	t.Log(dump(tree))
	var walked int
	walkFn := func(key string, value interface{}) error {
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
