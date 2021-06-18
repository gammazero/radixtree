package radixtree

import (
	"sort"
	"strings"
)

// Bytes is a radix tree of bytes with string keys and interface{} values.
type Bytes struct {
	root bytesNode
	size int
}

type bytesNode struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix string
	edges  byteEdges
	leaf   *leaf
}

// New creates a new bytes-based radix tree
func New() *Bytes {
	return new(Bytes)
}

type leaf struct {
	key   string
	value interface{}
}

type byteEdge struct {
	radix byte
	node  *bytesNode
}

// byteEdges implements sort.Interface
type byteEdges []byteEdge

func (e byteEdges) Len() int           { return len(e) }
func (e byteEdges) Less(i, j int) bool { return e[i].radix < e[j].radix }
func (e byteEdges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// BytesIterator is a stateful iterator that traverses a Bytes radix tree one
// byte at a time.
//
// Any modification to the tree invalidates the iterator.
type BytesIterator struct {
	p    int
	node *bytesNode
}

// Len returns the number of values stored in the tree.
func (tree *Bytes) Len() int {
	return tree.size
}

// Get returns the value stored at the given key.  Returns false if there is no
// value present for the key.
func (tree *Bytes) Get(key string) (interface{}, bool) {
	node := &tree.root
	// Consume key data while mathcing edge and prefix; return if remaining key
	// data matches nothing.
	for len(key) != 0 {
		// Find edge for radix
		node = node.getEdge(key[0])
		if node == nil {
			return nil, false
		}

		// Consume key data
		key = key[1:]
		if !strings.HasPrefix(key, node.prefix) {
			return nil, false
		}
		key = key[len(node.prefix):]
	}
	if node.leaf != nil {
		return node.leaf.value, true
	}
	return nil, false
}

// Put inserts the value into the tree at the given key, replacing any existing
// items.  It returns true if it adds a new value, false if it replaces an
// existing value.
func (tree *Bytes) Put(key string, value interface{}) bool {
	var (
		p          int
		isNewValue bool
		newEdge    byteEdge
		hasNewEdge bool
	)
	node := &tree.root

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
		// Descended as far as prefixes and edges match key, and still
		// have key data, so add child that has a prefix of the unmatched
		// key data and set its value to the new value.
		newChild := &bytesNode{
			leaf: &leaf{
				key:   key,
				value: value,
			},
		}
		if i < len(key)-1 {
			newChild.prefix = string(key[i+1:])
		}
		newEdge = byteEdge{radix, newChild}
		hasNewEdge = true
		break
	}
	// Key has been consumed by traversing prefixes and/or edges, or has
	// been put into new child.

	// If key partially matches node's prefix, then need to split node.
	if p < len(node.prefix) {
		node.split(p)
		isNewValue = true
	}

	if hasNewEdge {
		node.addEdge(newEdge)
		isNewValue = true
		tree.size++
	} else {
		// Store key at existing child
		if node.leaf == nil {
			isNewValue = true
			tree.size++
		}
		node.leaf = &leaf{
			key:   key,
			value: value,
		}
	}

	return isNewValue
}

// Delete removes the value associated with the given key. Returns true if
// there was a value stored for the key. If the node or any of its ancestors
// becomes childless as a result, they are removed from the tree.
func (tree *Bytes) Delete(key string) bool {
	node := &tree.root
	var (
		parents []*bytesNode
		links   []byte
	)
	for len(key) != 0 {
		parents = append(parents, node)

		// Find edge for radix
		node = node.getEdge(key[0])
		if node == nil {
			// node does not exist
			return false
		}
		links = append(links, key[0])

		// Consume key data
		key = key[1:]
		if !strings.HasPrefix(key, node.prefix) {
			return false
		}
		key = key[len(node.prefix):]
	}

	var deleted bool
	if node.leaf != nil {
		// delete the node value, indicate that value was deleted
		node.leaf = nil
		deleted = true
		tree.size--
	}

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	node = node.prune(parents, links)

	// If node has become compressible, compress it
	if node != &tree.root {
		node.compress()
	}

	return deleted
}

// Walk visits all nodes whose keys match or are prefixed by the specified key,
// calling walkFn for each value found.  If walkFn returns true, Walk returns.
// Use empty key "" to visit all nodes.
//
// The tree is traversed in lexical order, making the output deterministic.
//
// Walk can be thought of as GetItemsWithPrefix(key)
func (tree *Bytes) Walk(key string, walkFn WalkFunc) {
	node := &tree.root
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

	// Walk down tree starting at node located at key
	node.walk(walkFn)
}

// WalkPath walks a path in the tree from the root to the node at the given
// key, calling walkFn for each node that has a value.  If walkFn returns true,
// WalkPath returns.
//
// The tree is traversed in lexical order, making the output deterministic.
//
// WalkPath can be thought of as GetItemsThatArePrefixOf((key)
func (tree *Bytes) WalkPath(key string, walkFn WalkFunc) {
	node := &tree.root
	for {
		if node.leaf != nil && walkFn(node.leaf.key, node.leaf.value) {
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

// Inspect walks every node of the tree, whether or not it holds a value,
// calling inspectFn with information about each node.  This allows the
// structure of the tree to be examined and detailed statistics to be
// collected.
//
// If inspectFn returns false, the traversal is stopped and Inspect returns.
//
// The tree is traversed in lexical order, making the output deterministic.
func (tree *Bytes) Inspect(inspectFn InspectFunc) {
	tree.root.inspect("", "", 0, inspectFn)
}

// NewIterator returns a new BytesIterator instance that begins iterating from
// the root of the tree.
func (tree *Bytes) NewIterator() *BytesIterator {
	return &BytesIterator{
		node: &tree.root,
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.  These iterators
// do not affect each other and can be iterated concurrently.
func (it *BytesIterator) Copy() *BytesIterator {
	return &BytesIterator{
		p:    it.p,
		node: it.node,
	}
}

// Next advances the iterator from its current position, to the position of
// given key symbol in the tree, so long as the given symbol is next in a path
// in the tree.  If the symbol allows the iterator to advance, then true is
// returned.  Otherwise false is returned.
//
// When false is returned the iterator is not modified. This allows different
// values to be used in subsequent calls to Next, to advance the iterator from
// its current position.
func (it *BytesIterator) Next(radix byte) bool {
	// The tree.prefix represents single-edge parents without values that were
	// compressed out of the tree.  Let prefix consume key symbols.
	if it.p < len(it.node.prefix) {
		if radix == it.node.prefix[it.p] {
			// Key matches prefix so far, ok to continue.
			it.p++
			return true
		}
		// Some unmatched prefix remains, node not found
		return false
	}
	node := it.node.getEdge(radix)
	if node == nil {
		// No more prefix, no edges, so node not found
		return false
	}
	// Key symbol matched up to this edge, ok to continue.
	it.p = 0
	it.node = node
	return true
}

// Value returns the value at the current iterator position, and true or false
// to indicate if a value is present at the position.
func (it *BytesIterator) Value() (interface{}, bool) {
	// Only return value if all of this node's prefix was matched.  Otherwise,
	// have not fully traversed into this node (edge not completely traversed).
	if it.p != len(it.node.prefix) {
		return nil, false
	}
	if it.node.leaf == nil {
		return nil, false
	}
	return it.node.leaf.value, true
}

// split splits a node such that a node:
//     ("prefix", leaf, edges[])
// is split into parent branching node, and a child leaf node:
//     ("pre", nil, edges[f])--->("ix", leaf, edges[])
func (node *bytesNode) split(p int) {
	split := &bytesNode{
		edges: node.edges,
		leaf:  node.leaf,
	}
	if p < len(node.prefix)-1 {
		split.prefix = node.prefix[p+1:]
	}
	node.edges = nil
	node.addEdge(byteEdge{node.prefix[p], split})
	if p == 0 {
		node.prefix = ""
	} else {
		node.prefix = node.prefix[:p]
	}
	node.leaf = nil
}

func (node *bytesNode) prune(parents []*bytesNode, links []byte) *bytesNode {
	if node.edges != nil {
		return node
	}
	// iterate parents towards root of tree, removing the empty leaf
	for i := len(links) - 1; i >= 0; i-- {
		node = parents[i]
		node.delEdge(links[i])
		if len(node.edges) != 0 {
			// parent has other edges, stop
			break
		}
		node.edges = nil
		if node.leaf != nil {
			// parent has a value, stop
			break
		}
	}
	return node
}

func (node *bytesNode) compress() {
	if len(node.edges) != 1 || node.leaf != nil {
		return
	}
	edge := node.edges[0]
	pfx := make([]byte, len(node.prefix)+1+len(edge.node.prefix))
	copy(pfx, node.prefix)
	pfx[len(node.prefix)] = edge.radix
	copy(pfx[len(node.prefix)+1:], edge.node.prefix)
	node.prefix = string(pfx)
	node.leaf = edge.node.leaf
	node.edges = edge.node.edges
}

func (node *bytesNode) walk(walkFn WalkFunc) bool {
	if node.leaf != nil && walkFn(node.leaf.key, node.leaf.value) {
		return true
	}
	for _, edge := range node.edges {
		if edge.node.walk(walkFn) {
			return true
		}
	}
	return false
}

func (node *bytesNode) inspect(link, key string, depth int, inspectFn InspectFunc) bool {
	key += link + node.prefix
	var val interface{}
	var hasVal bool
	if node.leaf != nil {
		val = node.leaf.value
		hasVal = true
	}
	if inspectFn(link, node.prefix, key, depth, len(node.edges), hasVal, val) {
		return true
	}
	for _, edge := range node.edges {
		if edge.node.inspect(string(edge.radix), key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

// getEdge binary searches for edge
func (node *bytesNode) getEdge(radix byte) *bytesNode {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].radix >= radix
	})
	if idx < count && node.edges[idx].radix == radix {
		return node.edges[idx].node
	}
	return nil
}

// addEdge binary searches to find where to insert edge, and inserts at
func (node *bytesNode) addEdge(e byteEdge) {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].radix >= e.radix
	})
	node.edges = append(node.edges, byteEdge{})
	copy(node.edges[idx+1:], node.edges[idx:])
	node.edges[idx] = e
}

// delEdge binary searches for edge and removes it
func (node *bytesNode) delEdge(radix byte) {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].radix >= radix
	})
	if idx < count && node.edges[idx].radix == radix {
		copy(node.edges[idx:], node.edges[idx+1:])
		node.edges[len(node.edges)-1] = byteEdge{}
		node.edges = node.edges[:len(node.edges)-1]
	}
}
