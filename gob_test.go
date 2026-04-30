package radixtree_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"testing"

	"github.com/gammazero/radixtree"
)

func TestEncode(t *testing.T) {
	rt := radixtree.New[string]()

	keys := []string{
		"bird",
		"rat",
		"bat",
		"rats",
		"ratatouille",
		"rat/whis/key",
		"rat/whis/kers",
		"rat/whis/per/er",
		"rat/winks/wisely/once",
		"rat/winks/wisely/x/y/z",
		"rat/winks/wryly",
	}

	for _, key := range keys {
		rt.Put(key, strings.ToUpper(key))
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(rt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Encoded to", buf.Len(), "bytes")

	rt2 := radixtree.New[string]()

	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(rt2)
	if err != nil {
		t.Fatal(err)
	}

	if rt.Len() != rt2.Len() {
		t.Fatalf("decoded tree has wrong size, expectd %d got %d", rt.Len(), rt2.Len())
	}
	for _, key := range keys {
		checkItem(t, rt2, key, strings.ToUpper(key))
	}
	fmt.Println("Decoded tree matches original tree")
}

func checkItem(t *testing.T, rt *radixtree.Tree[string], key, value string) {
	val, ok := rt.Get(key)
	if !ok {
		t.Fatalf("decoded tree missing value for key %q", key)
	}
	if val != value {
		t.Fatalf("wrong value in tree, expected %s got%s", value, val)
	}
}
