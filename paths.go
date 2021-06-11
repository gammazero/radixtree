package radixtree

import (
	"sort"
	"strings"
)

const defaultPathSeparator = "/"

// Paths is a radix tree of paths with string keys and interface{}
// values. Paths splits keys by separator, for example, "/a/b/c" splits to
// "/a", "/b", "/c". A different path separator string may be specified by
// setting by calling NewPaths with a the string to use as a separator.
//
// Non-terminal nodes have nil values, so a stored nil value is not
// distinguishable and is not be included in Walk or WalkPath.
type Paths struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix  []string
	edges   pathEdges
	pathSep string
	leaf    *leaf
}

// NewPaths creates a new Paths instance allowin the path separator to be set.
//
// The pathSeparator splits a path key into separate segments.  The default
// path separator is forward slash.  This variable may be set to any string to
// allow a different path separator, multi-character strings are OK.
func NewPaths(pathSeparator string) *Paths {
	return &Paths{
		pathSep: pathSeparator,
	}
}

// PathSeparator returns this Paths instance's path separator
func (tree *Paths) PathSeparator() string {
	if tree.pathSep == "" {
		tree.pathSep = defaultPathSeparator
	}
	return tree.pathSep
}

type pathEdge struct {
	label string
	node  *Paths
}

// pathEdges implements sort.Interface
type pathEdges []pathEdge

func (e pathEdges) Len() int           { return len(e) }
func (e pathEdges) Less(i, j int) bool { return e[i].label < e[j].label }
func (e pathEdges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// PathsIterator traverses a Paths radix tree one path segment at a time.
//
// Note: Any modification to the tree invalidates the iterator.
type PathsIterator struct {
	p       int
	node    *Paths
	pathSep string
}

// NewIterator returns a new PathsIterator instance that begins iterating
// from the root of the tree.
func (tree *Paths) NewIterator() *PathsIterator {
	return &PathsIterator{
		node:    tree,
		pathSep: tree.PathSeparator(),
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.  These iterators
// do not affect eachother and can be iterated concurrently.
func (it *PathsIterator) Copy() *PathsIterator {
	return &PathsIterator{
		p:       it.p,
		node:    it.node,
		pathSep: it.pathSep,
	}
}

// Next advances the iterator from its current position, to the position of
// given path segment in the tree, so long as the given segment is next in a
// path in the tree.  If the segment allows the iterator to advance, then true
// is returned.  Otherwise false is returned.
//
// When false is returned the iterator is not modified. This allows different
// values to be used, in subsequent calls to Next, to advance the iterator from
// its current position.
//
// Any part subsequent to the first, must begin with the PathSeparator.
func (it *PathsIterator) Next(part string) bool {
	if part == "" {
		return false
	}
	part = strings.Trim(part, it.pathSep)

	if it.p < len(it.node.prefix) {
		if part == it.node.prefix[it.p] {
			it.p++
			return true
		}
		return false
	}
	node := it.node.getEdge(part)
	if node == nil {
		return false
	}
	it.p = 0
	it.node = node
	return true
}

// Value returns the value at the current iterator position, and true or false
// to indicate if there is a value at the position.
func (it *PathsIterator) Value() (interface{}, bool) {
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
func (tree *Paths) Get(key string) (interface{}, bool) {
	pathSep := tree.PathSeparator()
	iter := tree.NewIterator()
	for part, i := pathNext(key, pathSep, 0); part != ""; part, i = pathNext(key, pathSep, i) {
		if !iter.Next(part) {
			return nil, false
		}
	}
	return iter.Value()
}

// Put inserts the value into the tree at the given key, replacing any existing
// items.  It returns true if the put adds a new value, false if it replaces an
// existing value.
func (tree *Paths) Put(key string, value interface{}) bool {
	var (
		p          int
		isNewValue bool
		hasNewEdge bool
		newEdge    pathEdge
	)
	node := tree

	pathSep := tree.PathSeparator()
	for part, next := pathNext(key, pathSep, 0); part != ""; part, next = pathNext(key, pathSep, next) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				continue
			}
		} else if child := node.getEdge(part); child != nil {
			node = child
			p = 0
			continue
		}
		// Descended as far as prefixes and children match key, and still
		// have key data, so add child that has a prefix of the unmatched
		// key data and set its value to the new value.
		newChild := &Paths{
			leaf: &leaf{
				key:   key,
				value: value,
			},
		}
		childLink := part
		if next != -1 {
			newChild.prefix = []string{}
			for next != -1 {
				part, next = pathNext(key, pathSep, next)
				newChild.prefix = append(newChild.prefix, part)
			}
		}
		newEdge = pathEdge{childLink, newChild}
		hasNewEdge = true
		break
	}
	// Key has been consumed by traversing prefixes and/or children, or has
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
//     ("pre/fix/path", leaf, edges[])
// is split into parent branching node, and a child leaf node:
//     ("pre", nil, edges[f])--->("ix/path", leaf, edges[])
func (tree *Paths) split(p int) {
	split := &Paths{
		edges:   tree.edges,
		leaf:    tree.leaf,
		pathSep: tree.PathSeparator(),
	}
	if p < len(tree.prefix)-1 {
		split.prefix = tree.prefix[p+1:]
	}
	tree.edges = nil
	tree.addEdge(pathEdge{tree.prefix[p], split})
	if p == 0 {
		tree.prefix = nil
	} else {
		tree.prefix = tree.prefix[:p]
	}
	tree.leaf = nil
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the tree.
func (tree *Paths) Delete(key string) bool {
	node := tree
	var (
		nodes []*Paths
		parts []string
		p     int
	)
	pathSep := tree.PathSeparator()
	for part, i := pathNext(key, pathSep, 0); part != ""; part, i = pathNext(key, pathSep, i) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				continue
			}
			return false
		}
		nodes = append(nodes, node)
		parts = append(parts, part)
		node = node.getEdge(part)
		if node == nil {
			// no child for key segment so node does not exist
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

	// If node is leaf, remove from parent.  If parent becomes leaf, repeat.
	node = node.prune(nodes, parts)

	// If node has become compressible, compress it
	node.compress()

	return deleted
}

func (tree *Paths) prune(parents []*Paths, links []string) *Paths {
	if tree.edges != nil {
		return tree
	}
	// iterate parents towards root of tree, removing the empty leaf
	for i := len(links) - 1; i >= 0; i-- {
		tree = parents[i]
		tree.delEdge(links[i])
		if len(tree.edges) != 0 {
			// parent has other children, stop
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

func (tree *Paths) compress() {
	if len(tree.edges) != 1 || tree.leaf != nil {
		return
	}
	for _, edge := range tree.edges {
		pfx := make([]string, len(tree.prefix)+1+len(edge.node.prefix))
		copy(pfx, tree.prefix)
		pfx[len(tree.prefix)] = edge.label
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
func (tree *Paths) Walk(key string, walkFn WalkFunc) {
	pathSep := tree.PathSeparator()
	// Traverse tree to get to node at key
	if key != "" {
		iter := tree.NewIterator()
		for part, i := pathNext(key, pathSep, 0); part != ""; part, i = pathNext(key, pathSep, i) {
			if !iter.Next(part) {
				return
			}
		}
		tree = iter.node
	}

	// Walk down tree starting at node located at key
	tree.walk(walkFn)
}

func (tree *Paths) walk(walkFn WalkFunc) bool {
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
func (tree *Paths) WalkPath(key string, walkFn WalkFunc) {
	if tree.leaf != nil && walkFn("", tree.leaf.value) {
		return
	}
	pathSep := tree.PathSeparator()
	iter := tree.NewIterator()
	for part, i := pathNext(key, pathSep, 0); part != ""; part, i = pathNext(key, pathSep, i) {
		if !iter.Next(part) {
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
// returns Skip, Inspect will not descend into the node's children.
//
// The tree is traversed in lexical order, making the output deterministic.
func (tree *Paths) Inspect(inspectFn InspectFunc) {
	tree.inspect("", "", 0, inspectFn)
}

func (tree *Paths) inspect(link, key string, depth int, inspectFn InspectFunc) bool {
	pathSep := tree.PathSeparator()
	pfx := strings.Join(tree.prefix, pathSep)
	var keyParts []string
	if key != "" {
		keyParts = append(keyParts, key)
	}
	if link != "" {
		keyParts = append(keyParts, link)
	}
	if pfx != "" {
		keyParts = append(keyParts, pfx)
	}
	key = strings.Join(keyParts, pathSep)
	var val interface{}
	if tree.leaf != nil {
		val = tree.leaf.value
	}
	if inspectFn(link, pfx, key, depth, len(tree.edges), val) {
		return true
	}
	for _, edge := range tree.edges {
		if edge.node.inspect(edge.label, key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

// pathNext splits path strings by a path separator character. For
// example, "/a/b/c" -> ("a", 3), ("b", 5), ("c", -1) in successive
// calls. It does not allocate any heap memory.
func pathNext(path, pathSep string, start int) (string, int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	sepLen := len(pathSep)

	// Advance start past separators
	for strings.HasPrefix(path[start:], pathSep) {
		start += sepLen
	}
	if start == len(path) {
		return "", -1
	}

	// Segment ends at next separator or end of path
	end := strings.Index(path[start+1:], pathSep)
	if end == -1 || end == len(path)-sepLen {
		return path[start:], -1
	}
	end += start + 1 // change from relative to absolute offset

	// Next place to look is next non-separator past end, if any
	next := end + sepLen
	for strings.HasPrefix(path[next:], pathSep) {
		next += sepLen
	}
	if next == len(path) {
		next = -1
	}

	return path[start:end], next
}

func (tree *Paths) getEdge(radix string) *Paths {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].label >= radix
	})
	if idx < count && tree.edges[idx].label == radix {
		return tree.edges[idx].node
	}
	return nil
}

func (tree *Paths) addEdge(e pathEdge) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].label >= e.label
	})
	tree.edges = append(tree.edges, pathEdge{})
	copy(tree.edges[idx+1:], tree.edges[idx:])
	tree.edges[idx] = e
}

func (tree *Paths) delEdge(radix string) {
	count := len(tree.edges)
	idx := sort.Search(count, func(i int) bool {
		return tree.edges[i].label >= radix
	})
	if idx < count && tree.edges[idx].label == radix {
		copy(tree.edges[idx:], tree.edges[idx+1:])
		tree.edges[len(tree.edges)-1] = pathEdge{}
		tree.edges = tree.edges[:len(tree.edges)-1]
	}
}
