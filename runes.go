package radixtree

import (
	"sort"
)

// Runes is a radix tree of runes with string keys and interface{} values.
// Non-terminal nodes have nil values, so a stored nil value is not
// distinguishable and is not included in Walk or WalkPath.
type Runes struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix []rune
	edges  runeEdges
	leaf   *leaf
}

type runeEdge struct {
	radix rune
	node  *Runes
}

// runeEdges implements sort.Interface
type runeEdges []runeEdge

func (e runeEdges) Len() int           { return len(e) }
func (e runeEdges) Less(i, j int) bool { return e[i].radix < e[j].radix }
func (e runeEdges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// RunesIterator is a stateful iterator that traverses a Runes radix tree one
// character at a time.
//
// Note: Any modification to the tree invalidates the iterator.
type RunesIterator struct {
	p    int
	node *Runes
}

// NewIterator returns a new RunesIterator instance that begins iterating from
// the root of the tree.
func (tree *Runes) NewIterator() *RunesIterator {
	return &RunesIterator{
		node: tree,
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.  These iterators
// do not affect eachother and can be iterated concurrently.
func (it *RunesIterator) Copy() *RunesIterator {
	return &RunesIterator{
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
// values to be used, in subsequent calls to Next, to advance the iterator from
// its current position.
func (it *RunesIterator) Next(radix rune) bool {
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
// to indicate if there is a value at the position.
func (it *RunesIterator) Value() (interface{}, bool) {
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

// Get returns the value stored at the given key.  Returns false if the key does
// not identify a node that has a value.
func (tree *Runes) Get(k string) (interface{}, bool) {
	key := []rune(k)
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
		if !runesHasPrefix(key[1:], tree.prefix) {
			break
		}
		key = key[len(tree.prefix)+1:]
	}
	return nil, false
}

func runesHasPrefix(s, prefix []rune) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := range prefix {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

// Put inserts the value into the tree at the given key, replacing any existing
// items.  It returns true if it adds a new value, false if it replaces an
// existing value.
func (tree *Runes) Put(k string, value interface{}) bool {
	var (
		p          int
		isNewValue bool
		newEdge    runeEdge
		hasNewEdge bool
	)
	node := tree

	// Need to iterate key as slice of runes, otherwise indexes will be skipped
	// when a multibyte character is seen.
	key := []rune(k)
	for i, radix := range key {
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
		newChild := &Runes{
			leaf: &leaf{
				key:   k,
				value: value,
			},
		}
		if i < len(key)-1 {
			newChild.prefix = key[i+1:]
		}
		newEdge = runeEdge{radix, newChild}
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
			key:   k,
			value: value,
		}
	}

	return isNewValue
}

// split splits a node such that a node:
//     ("prefix", leaf, edges[])
// is split into parent branching node, and a child leaf node:
//     ("pre", nil, edges[f])--->("ix", leaf, edges[])
func (tree *Runes) split(p int) {
	split := &Runes{
		edges: tree.edges,
		leaf:  tree.leaf,
	}
	if p < len(tree.prefix)-1 {
		split.prefix = tree.prefix[p+1:]
	}
	tree.edges = nil
	tree.addEdge(runeEdge{tree.prefix[p], split})
	if p == 0 {
		tree.prefix = nil
	} else {
		tree.prefix = tree.prefix[:p]
	}
	tree.leaf = nil
}

// Delete removes the value associated with the given key.  Returns true if a
// node was found for the given key, and that node held a value.  If the node
// has no edges (is a leaf node) it is removed from the tree.  If any of of the
// node's ancestors becomes childless as a result, they are also removed from
// the tree.
func (tree *Runes) Delete(key string) bool {
	node := tree
	var (
		nodes []*Runes
		links []rune
		p     int
	)
	for _, radix := range key {
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

func (tree *Runes) prune(parents []*Runes, links []rune) *Runes {
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

func (tree *Runes) compress() {
	if len(tree.edges) != 1 || tree.leaf != nil {
		return
	}
	for _, edge := range tree.edges {
		pfx := make([]rune, len(tree.prefix)+1+len(edge.node.prefix))
		copy(pfx, tree.prefix)
		pfx[len(tree.prefix)] = edge.radix
		copy(pfx[len(tree.prefix)+1:], edge.node.prefix)
		tree.prefix = pfx
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
func (tree *Runes) Walk(k string, walkFn WalkFunc) {
	if k != "" {
		for key := []rune(k); len(key) != 0; {
			tree = tree.getEdge(key[0])
			if tree == nil {
				return
			}

			// Consume key data
			if !runesHasPrefix(key[1:], tree.prefix) {
				if runesHasPrefix(tree.prefix, key[1:]) {
					break
				}
				return
			}
			key = key[len(tree.prefix)+1:]
		}
	}

	// Walk down tree starting at node located at key
	tree.walk(walkFn)
}

func (tree *Runes) walk(walkFn WalkFunc) bool {
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

// WalkPath walks the path in tree from the root to the node at the given key,
// calling walkFn for each node that has a value.  If walkFn returns true,
// WalkPath returns.
//
// The tree is traversed in lexical order, making the output deterministic.
//
// WalkPath can be thought of as GetItemsThatArePrefixOf((key)
func (tree *Runes) WalkPath(key string, walkFn WalkFunc) {
	if tree.leaf != nil && walkFn("", tree.leaf.value) {
		return
	}
	iter := tree.NewIterator()
	for _, r := range key {
		if !iter.Next(r) {
			return
		}
		if value, ok := iter.Value(); ok {
			if walkFn(iter.node.leaf.key, value) {
				return
			}
		}
	}
	return
}

// Inspect walks every node of the tree, whether or not it holds a value,
// calling inspectFn with information about each node.  This allows the
// structure of the tree to be examined and detailed statistics to be
// collected.
//
// If inspectFn returns an error, the traversal is aborted.  If inspectFn
// returns Skip, Inspect will not descend into the node's edges.
//
// The tree is traversed in lexical order, making the output deterministic.
func (tree *Runes) Inspect(inspectFn InspectFunc) {
	tree.inspect("", "", 0, inspectFn)
}

func (tree *Runes) inspect(link, key string, depth int, inspectFn InspectFunc) bool {
	pfx := string(tree.prefix)
	key += link + pfx
	var val interface{}
	if tree.leaf != nil {
		val = tree.leaf.value
	}
	if inspectFn(link, pfx, key, depth, len(tree.edges), val) {
		return true
	}
	for _, edge := range tree.edges {
		if edge.node.inspect(string(edge.radix), key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

func (tree *Runes) getEdge(radix rune) *Runes {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= radix
	})
	if idx < count && tree.edges[idx].radix == radix {
		return tree.edges[idx].node
	}
	return nil
}

func (tree *Runes) addEdge(e runeEdge) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= e.radix
	})
	tree.edges = append(tree.edges, runeEdge{})
	copy(tree.edges[idx+1:], tree.edges[idx:])
	tree.edges[idx] = e
}

func (tree *Runes) delEdge(radix rune) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].radix >= radix
	})
	if idx < count && tree.edges[idx].radix == radix {
		copy(tree.edges[idx:], tree.edges[idx+1:])
		tree.edges[len(tree.edges)-1] = runeEdge{}
		tree.edges = tree.edges[:len(tree.edges)-1]
	}
}
