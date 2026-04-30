package radixtree

import (
	"iter"
	"strings"
)

// Tree is a radix tree of bytes keys and any values.
type Tree[T any] struct {
	root radixNode[T]
	size int
}

// New creates a new bytes-based radix tree
func New[T any]() *Tree[T] {
	return new(Tree[T])
}

type radixNode[T any] struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix  string
	radices []byte
	nodes   []*radixNode[T]
	leaf    *Item[T]
}

// InspectFunc is the type of the function called for each node visited by
// Inspect. The key argument contains the key at which the node is located, the
// depth is the distance from the root of the tree, and children is the number
// of children the node has.
//
// If the function returns true Inspect stops immediately and returns.
type InspectFunc[T any] func(link, prefix, key string, depth, children int, hasValue bool, value T) bool

type Item[T any] struct {
	key   string
	value T
}

func (kv *Item[T]) Key() string { return kv.key }
func (kv *Item[T]) Value() T    { return kv.value }

// Len returns the number of values stored in the tree.
func (t *Tree[T]) Len() int {
	return t.size
}

// Get returns the value stored at the given key. Returns false if there is no
// value present for the key.
func (t *Tree[T]) Get(key string) (T, bool) {
	var zero T
	node := &t.root
	// Consume key data while mathcing edge and prefix; return if remaining key
	// data matches nothing.
	for len(key) != 0 {
		// Find edge for radix.
		node = node.getEdge(key[0])
		if node == nil {
			return zero, false
		}

		// Consume key data.
		key = key[1:]
		if !strings.HasPrefix(key, node.prefix) {
			return zero, false
		}
		key = key[len(node.prefix):]
	}
	if node.leaf != nil {
		return node.leaf.value, true
	}
	return zero, false
}

// Put inserts the value into the tree at the given key, replacing any existing
// items. It returns true if it adds a new value, false if it replaces an
// existing value.
func (t *Tree[T]) Put(key string, value T) bool {
	var (
		p            int
		isNewValue   bool
		newEdgeRadix byte
		newEdgeNode  *radixNode[T]
		hasNewEdge   bool
	)
	node := &t.root

	for i := 0; i < len(key); i++ {
		radix := key[i]
		if p < len(node.prefix) {
			if radix == node.prefix[p] {
				p++
				continue
			}
		} else if child := node.getEdge(radix); child != nil {
			node = child
			p = 0
			continue
		}
		// Descended as far as prefixes and edges match key, and still have key
		// data, so add child that has a prefix of the unmatched key data and
		// set its value to the new value.
		newChild := &radixNode[T]{
			leaf: &Item[T]{
				key:   key,
				value: value,
			},
		}
		if i < len(key)-1 {
			newChild.prefix = key[i+1:]
		}
		newEdgeRadix = radix
		newEdgeNode = newChild
		hasNewEdge = true
		break
	}
	// Key has been consumed by traversing prefixes and/or edges, or has been
	// put into new child.

	// If key partially matches node's prefix, then need to split node.
	if p < len(node.prefix) {
		node.split(p)
		isNewValue = true
	}

	if hasNewEdge {
		node.addEdge(newEdgeRadix, newEdgeNode)
		isNewValue = true
		t.size++
	} else {
		// Store key at existing child
		if node.leaf == nil {
			isNewValue = true
			t.size++
		}
		node.leaf = &Item[T]{
			key:   key,
			value: value,
		}
	}

	return isNewValue
}

// Delete removes the value associated with the given key. Returns true if
// there was a value stored for the key. If the node or any of its ancestors
// becomes childless as a result, they are removed from the tree.
func (t *Tree[T]) Delete(key string) bool {
	node := &t.root
	var (
		parents []*radixNode[T]
		links   []byte
	)
	for len(key) != 0 {
		parents = append(parents, node)

		// Find edge for radix.
		node = node.getEdge(key[0])
		if node == nil {
			// node does not exist.
			return false
		}
		links = append(links, key[0])

		// Consume key data.
		key = key[1:]
		if !strings.HasPrefix(key, node.prefix) {
			return false
		}
		key = key[len(node.prefix):]
	}

	if node.leaf == nil {
		return false
	}

	// delete the node value, indicate that value was deleted.
	node.leaf = nil
	t.size--

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	node = node.prune(parents, links)

	// If node has become compressible, compress it.
	if node != &t.root {
		node.compress()
	}

	return true
}

// DeletePrefix removes all values whose key is prefixed by the given prefix.
// Returns true if any values were removed.
func (t *Tree[T]) DeletePrefix(prefix string) bool {
	node := &t.root
	var (
		parents []*radixNode[T]
		links   []byte
	)
	for len(prefix) != 0 {
		parents = append(parents, node)

		// Find edge for radix.
		node = node.getEdge(prefix[0])
		if node == nil {
			// Node does not exist.
			return false
		}
		links = append(links, prefix[0])

		// Consume prefix.
		prefix = prefix[1:]
		if !strings.HasPrefix(prefix, node.prefix) {
			if strings.HasPrefix(node.prefix, prefix) {
				// Prefix consumed, so it prefixes every key from node down.
				break
			}
			return false
		}
		prefix = prefix[len(node.prefix):]
	}

	if node.radices != nil {
		var count int
		for range node.iter() {
			count++
		}
		t.size -= count
		node.radices = nil
		node.nodes = nil
	} else {
		t.size--
	}
	node.leaf = nil

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	node = node.prune(parents, links)

	// If node has become compressible, compress it.
	if node != &t.root {
		node.compress()
	}

	return true
}

// Iter visits all nodes in the tree, yielding the key and value of each.
//
// The tree is traversed in lexical order, making the output deterministic.
func (t *Tree[T]) Iter() iter.Seq2[string, T] {
	return t.root.iter()
}

// IterAt visits all nodes whose keys match or are prefixed by the specified
// key, yielding the key and value of each. An empty key "" to visits all
// nodes, and is the same as calling Iter.
//
// The tree is traversed in lexical order, making the output deterministic.
func (t *Tree[T]) IterAt(key string) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		// Find the subtree with a matching prefix.
		node := &t.root
		for len(key) != 0 {
			if node = node.getEdge(key[0]); node == nil {
				return
			}

			// Consume key data
			key = key[1:]
			if !strings.HasPrefix(key, node.prefix) {
				if strings.HasPrefix(node.prefix, key) {
					break
				}
				return
			}
			key = key[len(node.prefix):]
		}
		// Iterate the subtree.
		node.walk(yield)
	}
}

func (node *radixNode[T]) iter() iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		node.walk(yield)
	}
}

func (node *radixNode[T]) walk(yield func(string, T) bool) bool {
	if node.leaf != nil && !yield(node.leaf.key, node.leaf.value) {
		return false
	}
	for _, child := range node.nodes {
		if !child.walk(yield) {
			return false
		}
	}
	return true
}

// IterPath returns an iterator that visits each node along the path from the
// root to the node at the given key. yielding the key and value of each.
//
// The tree is traversed in lexical order, making the output deterministic.
func (t *Tree[T]) IterPath(key string) iter.Seq2[string, T] {
	return func(yield func(string, T) bool) {
		node := &t.root
		for {
			if node.leaf != nil && !yield(node.leaf.key, node.leaf.value) {
				return
			}

			if len(key) == 0 {
				return
			}

			if node = node.getEdge(key[0]); node == nil {
				return
			}

			key = key[1:]
			if !strings.HasPrefix(key, node.prefix) {
				return
			}
			key = key[len(node.prefix):]
		}
	}
}

// Inspect walks every node of the tree, whether or not it holds a value,
// calling inspectFn with information about each node. This allows the
// structure of the tree to be examined and detailed statistics to be
// collected.
//
// If inspectFn returns false, the traversal is stopped and Inspect returns.
//
// The tree is traversed in lexical order, making the output deterministic.
func (t *Tree[T]) Inspect(inspectFn InspectFunc[T]) {
	t.root.inspect("", "", 0, inspectFn)
}

// split splits a node such that a node:
//
//	("prefix", leaf, edges[])
//
// is split into parent branching node, and a child leaf node:
//
//	("pre", nil, edges[f])--->("ix", leaf, edges[])
func (node *radixNode[T]) split(p int) {
	split := &radixNode[T]{
		radices: node.radices,
		nodes:   node.nodes,
		leaf:    node.leaf,
	}
	if p < len(node.prefix)-1 {
		split.prefix = node.prefix[p+1:]
	}
	node.radices = nil
	node.nodes = nil
	node.addEdge(node.prefix[p], split)
	if p == 0 {
		node.prefix = ""
	} else {
		node.prefix = node.prefix[:p]
	}
	node.leaf = nil
}

func (node *radixNode[T]) prune(parents []*radixNode[T], links []byte) *radixNode[T] {
	if node.radices != nil {
		return node
	}
	// iterate parents towards root of tree, removing the empty leaf.
	for i := len(links) - 1; i >= 0; i-- {
		node = parents[i]
		node.delEdge(links[i])
		if len(node.radices) != 0 {
			// parent has other edges, stop.
			break
		}
		node.radices = nil
		node.nodes = nil
		if node.leaf != nil {
			// parent has a value, stop.
			break
		}
	}
	return node
}

func (node *radixNode[T]) compress() {
	if len(node.radices) != 1 || node.leaf != nil {
		return
	}
	r := node.radices[0]
	child := node.nodes[0]
	var b strings.Builder
	b.Grow(len(node.prefix) + 1 + len(child.prefix))
	b.WriteString(node.prefix)
	b.WriteByte(r)
	b.WriteString(child.prefix)
	node.prefix = b.String()
	node.leaf = child.leaf
	node.radices = child.radices
	node.nodes = child.nodes
}

func (node *radixNode[T]) inspect(link, key string, depth int, inspectFn InspectFunc[T]) bool {
	key += link + node.prefix
	var val T
	var hasVal bool
	if node.leaf != nil {
		val = node.leaf.value
		hasVal = true
	}
	if inspectFn(link, node.prefix, key, depth, len(node.radices), hasVal, val) {
		return true
	}
	for i, child := range node.nodes {
		if child.inspect(string(node.radices[i]), key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

// indexEdge binary searches for the edge index.
//
// This is faster then going through sort.Interface for repeated searches.
func (node *radixNode[T]) indexEdge(radix byte) int {
	n := len(node.radices)
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if node.radices[h] < radix {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// getEdge binary searches for edge.
func (node *radixNode[T]) getEdge(radix byte) *radixNode[T] {
	idx := node.indexEdge(radix)
	if idx < len(node.radices) && node.radices[idx] == radix {
		return node.nodes[idx]
	}
	return nil
}

// addEdge binary searches to find where to insert edge, and inserts at.
func (node *radixNode[T]) addEdge(radix byte, child *radixNode[T]) {
	idx := node.indexEdge(radix)
	node.radices = append(node.radices, 0)
	copy(node.radices[idx+1:], node.radices[idx:])
	node.radices[idx] = radix
	node.nodes = append(node.nodes, nil)
	copy(node.nodes[idx+1:], node.nodes[idx:])
	node.nodes[idx] = child
}

// delEdge binary searches for edge and removes it.
func (node *radixNode[T]) delEdge(radix byte) {
	idx := node.indexEdge(radix)
	if idx < len(node.radices) && node.radices[idx] == radix {
		copy(node.radices[idx:], node.radices[idx+1:])
		node.radices[len(node.radices)-1] = 0
		node.radices = node.radices[:len(node.radices)-1]
		copy(node.nodes[idx:], node.nodes[idx+1:])
		node.nodes[len(node.nodes)-1] = nil
		node.nodes = node.nodes[:len(node.nodes)-1]
	}
}
