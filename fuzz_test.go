package radixtree

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

// FuzzTree applies a sequence of operations, decoded from the fuzz input, to
// both a Tree and a map, and checks that they agree after every operation.
//
// Each operation is encoded as an opcode byte, a key length byte, and that
// many key bytes. Keys are mapped onto a small alphabet so that they share
// prefixes often, which exercises node splitting and compression.
func FuzzTree(f *testing.F) {
	f.Add([]byte{0, 3, 't', 'o', 'm', 0, 6, 't', 'o', 'm', 'a', 't', 'o', 1, 3, 't', 'o', 'm'})
	f.Add([]byte{0, 2, 'a', 'b', 0, 2, 'a', 'c', 1, 2, 'a', 'b', 1, 2, 'a', 'c', 2, 0})
	f.Add([]byte{0, 0, 0, 1, 'a', 2, 0, 3, 1, 'a'})
	f.Add([]byte{0, 4, 'a', 'b', 'c', 'd', 0, 2, 'a', 'b', 2, 3, 'a', 'b', 'c'})

	f.Fuzz(func(t *testing.T, data []byte) {
		const alphabet = "abc/\x00\xff"
		tree := New[int]()
		model := map[string]int{}

		for step := 0; len(data) >= 2; step++ {
			op, n := data[0]%4, int(data[1]%8)
			data = data[2:]
			n = min(n, len(data))
			var b strings.Builder
			for _, c := range data[:n] {
				b.WriteByte(alphabet[int(c)%len(alphabet)])
			}
			key := b.String()
			data = data[n:]

			switch op {
			case 0:
				_, exists := model[key]
				if tree.Put(key, step) == exists {
					t.Fatalf("Put(%q) returned %v, key existed: %v", key, !exists, exists)
				}
				model[key] = step
			case 1:
				_, exists := model[key]
				if tree.Delete(key) != exists {
					t.Fatalf("Delete(%q) returned %v, key existed: %v", key, !exists, exists)
				}
				delete(model, key)
			case 2:
				var deleted bool
				for k := range model {
					if strings.HasPrefix(k, key) {
						delete(model, k)
						deleted = true
					}
				}
				if tree.DeletePrefix(key) != deleted {
					t.Fatalf("DeletePrefix(%q) returned %v, expected %v", key, !deleted, deleted)
				}
			case 3:
				checkQueries(t, tree, model, key)
			}
			checkTree(t, tree, model)
		}
	})
}

// checkTree verifies that the tree holds exactly the contents of model, and
// that the internal structure of the tree is valid.
func checkTree(t *testing.T, tree *Tree[int], model map[string]int) {
	t.Helper()
	if tree.Len() != len(model) {
		t.Fatalf("Len() = %d, expected %d", tree.Len(), len(model))
	}

	keys := slices.Sorted(maps.Keys(model))
	var i int
	for k, v := range tree.All() {
		if i >= len(keys) {
			t.Fatalf("All yielded extra key %q", k)
		}
		if k != keys[i] || v != model[k] {
			t.Fatalf("All yielded (%q, %d), expected (%q, %d)", k, v, keys[i], model[keys[i]])
		}
		i++
	}
	if i != len(keys) {
		t.Fatalf("All yielded %d items, expected %d", i, len(keys))
	}

	for k, v := range model {
		got, ok := tree.Get(k)
		if !ok || got != v {
			t.Fatalf("Get(%q) = (%d, %v), expected (%d, true)", k, got, ok, v)
		}
	}

	var values int
	for node := range tree.Nodes() {
		if node.HasValue {
			values++
			if v, ok := model[node.Key]; !ok || v != node.Value {
				t.Fatalf("node at %q has value %d, expected %d (present: %v)", node.Key, node.Value, v, ok)
			}
		} else if node.Depth != 0 && node.Children < 2 {
			t.Fatalf("node at %q has no value and %d children, should be compressed", node.Key, node.Children)
		}
	}
	if values != len(model) {
		t.Fatalf("Nodes yielded %d values, expected %d", values, len(model))
	}

	checkNode(t, &tree.root, "", true)
}

// checkNode verifies that the parallel radices and nodes slices are consistent
// and sorted, and that each leaf is stored at its own key.
func checkNode(t *testing.T, node *radixNode[int], key string, isRoot bool) {
	t.Helper()
	key += node.prefix
	if isRoot && node.prefix != "" {
		t.Fatalf("root has prefix %q", node.prefix)
	}
	if len(node.radices) != len(node.nodes) {
		t.Fatalf("node at %q has %d radices and %d nodes", key, len(node.radices), len(node.nodes))
	}
	if !slices.IsSorted(node.radices) || len(slices.Compact(slices.Clone(node.radices))) != len(node.radices) {
		t.Fatalf("node at %q has unsorted or duplicate radices %q", key, node.radices)
	}
	if !isRoot && node.leaf == nil && len(node.nodes) == 0 {
		t.Fatalf("node at %q has no value and no children", key)
	}
	if node.leaf != nil && node.leaf.key != key {
		t.Fatalf("node at %q holds key %q", key, node.leaf.key)
	}
	for i, child := range node.nodes {
		checkNode(t, child, key+string(node.radices[i:i+1]), false)
	}
}

// checkQueries verifies WithPrefix, PathTo, and Stepper against the model.
func checkQueries(t *testing.T, tree *Tree[int], model map[string]int, key string) {
	t.Helper()
	keys := slices.Sorted(maps.Keys(model))

	var expect []string
	for _, k := range keys {
		if strings.HasPrefix(k, key) {
			expect = append(expect, k)
		}
	}
	var got []string
	for k, v := range tree.WithPrefix(key) {
		if v != model[k] {
			t.Fatalf("WithPrefix(%q) yielded value %d at %q, expected %d", key, v, k, model[k])
		}
		got = append(got, k)
	}
	if !slices.Equal(got, expect) {
		t.Fatalf("WithPrefix(%q) yielded %q, expected %q", key, got, expect)
	}

	expect = expect[:0]
	for _, k := range keys {
		if strings.HasPrefix(key, k) {
			expect = append(expect, k)
		}
	}
	got = got[:0]
	for k := range tree.PathTo(key) {
		got = append(got, k)
	}
	if !slices.Equal(got, expect) {
		t.Fatalf("PathTo(%q) yielded %q, expected %q", key, got, expect)
	}

	// Step through the key one byte at a time, checking for a value at each
	// position.
	s := tree.NewStepper()
	for i := 0; ; i++ {
		v, ok := s.Value()
		k, iv, iok := s.Item()
		mv, mok := model[key[:i]]
		if ok != mok || v != mv || iok != mok || iv != mv || (iok && k != key[:i]) {
			t.Fatalf("Stepper at %q: Value() = (%d, %v), Item() = (%q, %d, %v), expected (%d, %v)", key[:i], v, ok, k, iv, iok, mv, mok)
		}
		if i == len(key) {
			break
		}
		if !s.Next(key[i]) {
			// No key in the tree may extend past this point.
			for _, mk := range keys {
				if strings.HasPrefix(mk, key[:i+1]) {
					t.Fatalf("Stepper could not advance to %q, but key %q exists", key[:i+1], mk)
				}
			}
			break
		}
	}
}
