package radixtree

import (
	"sort"
	"strings"
)

// Bytes is a radix tree of bytes with string keys and interface{} values.
type Bytes struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix string
	edges  byteEdges
	leaf   *leaf
}

func New() *Bytes {
	return new(Bytes)
}

type leaf struct {
	key   string
	value interface{}
}

type byteEdge struct {
	radix byte
	node  *Bytes
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
	node *Bytes
}

// NewIterator returns a new BytesIterator instance that begins iterating from
// the root of the tree.
func (tree *Bytes) NewIterator() *BytesIterator {
	return &BytesIterator{
		node: tree,
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.  These iterators
// do not affect eachother and can be iterated concurrently.
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

// Get returns the value stored at the given key.  Returns false if there is no
// value present for the key.
func (tree *Bytes) Get(key string) (interface{}, bool) {
	for {
		// All key data consumed and matched against node prefix, so this is
		// the requested or an intermediate node.
		if len(key) == 0 {
			if tree.leaf != nil {
				return tree.leaf.value, true
			}
			break
		}

		// Find edge for radix
		tree = tree.getEdge(key[0])
		if tree == nil {
			break
		}

		// Consume key data
		if !strings.HasPrefix(key[1:], tree.prefix) {
			break
		}
		key = key[len(tree.prefix)+1:]
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
	node := tree

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
		newChild := &Bytes{
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
	} else {
		// Store key at existing child
		if node.leaf == nil {
			isNewValue = true
		}
		node.leaf = &leaf{
			key:   key,
			value: value,
		}
	}

	return isNewValue
}

// split splits a node such that a node:
//     ("prefix", leaf, edges[])
// is split into parent branching node, and a child leaf node:
//     ("pre", nil, edges[f])--->("ix", leaf, edges[])
func (tree *Bytes) split(p int) {
	split := &Bytes{
		edges: tree.edges,
		leaf:  tree.leaf,
	}
	if p < len(tree.prefix)-1 {
		split.prefix = tree.prefix[p+1:]
	}
	tree.edges = nil
	tree.addEdge(byteEdge{tree.prefix[p], split})
	if p == 0 {
		tree.prefix = ""
	} else {
		tree.prefix = tree.prefix[:p]
	}
	tree.leaf = nil
}

// Delete removes the value associated with the given key. Returns true if
// there was a value stored for the key. If the node or any of its ancestors
// becomes childless as a result, they are removed from the tree.
func (tree *Bytes) Delete(key string) bool {
	node := tree
	var (
		nodes []*Bytes
		links []byte
		p     int
	)
	for i := 0; i < len(key); i++ {
		radix := key[i]
		if p < len(node.prefix) {
			if radix == node.prefix[p] {
				p++
				continue
			}
			return false
		}
		nodes = append(nodes, node)
		links = append(links, radix)
		node = node.getEdge(radix)
		if node == nil {
			// node does not exist
			return false
		}
		p = 0
	}

	// Key was not completely consumed traversing tree, so tree does not
	// contain anything for key.
	if p < len(node.prefix) {
		return false
	}
	var deleted bool
	if node.leaf != nil {
		// delete the node value, indicate that value was deleted
		node.leaf = nil
		deleted = true
	}

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	node = node.prune(nodes, links)

	// If node has become compressible, compress it
	node.compress()

	return deleted
}

func (tree *Bytes) prune(parents []*Bytes, links []byte) *Bytes {
	if tree.edges != nil {
		return tree
	}
	// iterate parents towards root of tree, removing the empty leaf
	for i := len(links) - 1; i >= 0; i-- {
		tree = parents[i]
		tree.delEdge(links[i])
		if len(tree.edges) != 0 {
			// parent has other edges, stop
			break
		}
		tree.edges = nil
		if tree.leaf != nil {
			// parent has a value, stop
			break
		}
	}
	return tree
}

func (tree *Bytes) compress() {
	if len(tree.edges) != 1 || tree.leaf != nil {
		return
	}
	for _, edge := range tree.edges {
		pfx := make([]byte, len(tree.prefix)+1+len(edge.node.prefix))
		copy(pfx, tree.prefix)
		pfx[len(tree.prefix)] = edge.radix
		copy(pfx[len(tree.prefix)+1:], edge.node.prefix)
		tree.prefix = string(pfx)
		tree.leaf = edge.node.leaf
		tree.edges = edge.node.edges
	}
}

// Walk visits all nodes whose keys match or are prefixed by the specified key,
// calling walkFn for each value found.  If walkFn returns true, Walk returns.
// Use empty key "" to visit all nodes.
//
// The tree is traversed in lexical order, making the output deterministic.
//
// Walk can be thought of as GetItemsWithPrefix(key)
func (tree *Bytes) Walk(key string, walkFn WalkFunc) {
	for len(key) != 0 {
		if tree = tree.getEdge(key[0]); tree == nil {
			return
		}

		// Consume key data
		if !strings.HasPrefix(key[1:], tree.prefix) {
			if strings.HasPrefix(tree.prefix, key[1:]) {
				break
			}
			return
		}
		key = key[len(tree.prefix)+1:]
	}

	// Walk down tree starting at node located at key
	tree.walk(walkFn)
}

func (tree *Bytes) walk(walkFn WalkFunc) bool {
	if tree.leaf != nil && walkFn(tree.leaf.key, tree.leaf.value) {
		return true
	}
	for _, edge := range tree.edges {
		if edge.node.walk(walkFn) {
			return true
		}
	}
	return false
}

// WalkPath walks a path in the tree from the root to the node at the given
// key, calling walkFn for each node that has a value.  If walkFn returns true,
// WalkPath returns.
//
// The tree is traversed in lexical order, making the output deterministic.
//
// WalkPath can be thought of as GetItemsThatArePrefixOf((key)
func (tree *Bytes) WalkPath(key string, walkFn WalkFunc) {
	for {
		if tree.leaf != nil && walkFn(tree.leaf.key, tree.leaf.value) {
			return
		}

		if len(key) == 0 {
			return
		}

		if tree = tree.getEdge(key[0]); tree == nil {
			return
		}

		if !strings.HasPrefix(key[1:], tree.prefix) {
			return
		}
		key = key[len(tree.prefix)+1:]
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
	tree.inspect("", "", 0, inspectFn)
}

func (tree *Bytes) inspect(link, key string, depth int, inspectFn InspectFunc) bool {
	key += link + tree.prefix
	var val interface{}
	if tree.leaf != nil {
		val = tree.leaf.value
	}
	if inspectFn(link, tree.prefix, key, depth, len(tree.edges), val) {
		return true
	}
	for _, edge := range tree.edges {
		if edge.node.inspect(string(edge.radix), key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

// getEdge binary searchs for edge
func (tree *Bytes) getEdge(radix byte) *Bytes {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= radix
	})
	if idx < count && tree.edges[idx].radix == radix {
		return tree.edges[idx].node
	}
	return nil
}

// addEdge binary searchs to find where to insert edge, and inserts at
func (tree *Bytes) addEdge(e byteEdge) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= e.radix
	})
	tree.edges = append(tree.edges, byteEdge{})
	copy(tree.edges[idx+1:], tree.edges[idx:])
	tree.edges[idx] = e
}

// delEdge binary searches for edge and removes it
func (tree *Bytes) delEdge(radix byte) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= radix
	})
	if idx < count && tree.edges[idx].radix == radix {
		copy(tree.edges[idx:], tree.edges[idx+1:])
		tree.edges[len(tree.edges)-1] = byteEdge{}
		tree.edges = tree.edges[:len(tree.edges)-1]
	}
}
