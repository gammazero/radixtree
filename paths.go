package radixtree

import (
	"sort"
	"strings"
)

const defaultPathSeparator = "/"

// Paths is a radix tree of paths with string keys and interface{}
// values. Paths splits keys by separator, for example, "/a/b/c" splits to
// "/a", "/b", "/c". A different path separator string may be specified by
// calling NewPaths with the separator to use.
type Paths struct {
	pathSep string
	root    pathsNode
	size    int
}

type pathsNode struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix []string
	edges  pathEdges
	leaf   *leaf
}

// NewPaths creates a new Paths instance, specifying the path separator to use.
//
// The pathSeparator splits a path key into separate segments.  The default
// path separator is forward slash.  This variable may be set to any string to
// use as a path separator, multi-character strings are OK.
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
	node  *pathsNode
}

// pathEdges implements sort.Interface
type pathEdges []pathEdge

func (e pathEdges) Len() int           { return len(e) }
func (e pathEdges) Less(i, j int) bool { return e[i].label < e[j].label }
func (e pathEdges) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

// PathsIterator traverses a Paths radix tree one path segment at a time.
//
// Any modification to the tree invalidates the iterator.
type PathsIterator struct {
	p       int
	node    *pathsNode
	pathSep string
}

// Len returns the number of values stored in the tree.
func (tree *Paths) Len() int {
	return tree.size
}

// Get returns the value stored at the given key.  Returns false if there is no
// value present for the key.
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
// items.  It returns true if it adds a new value, false if it replaces an
// existing value.
func (tree *Paths) Put(key string, value interface{}) bool {
	var (
		p          int
		isNewValue bool
		hasNewEdge bool
		newEdge    pathEdge
	)
	node := &tree.root

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
		newChild := &pathsNode{
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
func (tree *Paths) Delete(key string) bool {
	node := &tree.root
	var (
		parents []*pathsNode
		parts   []string
		p       int
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
		parents = append(parents, node)
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
		tree.size--
	}

	// If node is leaf, remove from parent.  If parent becomes leaf, repeat.
	node = node.prune(parents, parts)

	// If node has become compressible, compress it
	node.compress()

	return deleted
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
	node := &tree.root

	// Traverse tree to get to node at key
	if key != "" {
		iter := tree.NewIterator()
		for part, i := pathNext(key, pathSep, 0); part != ""; part, i = pathNext(key, pathSep, i) {
			if !iter.Next(part) {
				return
			}
		}
		node = iter.node
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
func (tree *Paths) WalkPath(key string, walkFn WalkFunc) {
	node := &tree.root
	if node.leaf != nil && walkFn("", node.leaf.value) {
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
}

// Inspect walks every node of the tree, whether or not it holds a value,
// calling inspectFn with information about each node.  This allows the
// structure of the tree to be examined and detailed statistics to be
// collected.
//
// If inspectFn returns false, the traversal is stopped and Inspect returns.
//
// The tree is traversed in lexical order, making the output deterministic.
func (tree *Paths) Inspect(inspectFn InspectFunc) {
	tree.root.inspect(tree.PathSeparator(), "", "", 0, inspectFn)
}

// NewIterator returns a new PathsIterator instance that begins iterating
// from the root of the tree.
func (tree *Paths) NewIterator() *PathsIterator {
	return &PathsIterator{
		node:    &tree.root,
		pathSep: tree.PathSeparator(),
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.  These iterators
// do not affect each other and can be iterated concurrently.
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
// values to be used in subsequent calls to Next, to advance the iterator from
// its current position.
//
// Any part subsequent to the first must begin with the PathSeparator.
func (it *PathsIterator) Next(part string) bool {
	if part == "" {
		return false
	}
	part = strings.Trim(part, it.pathSep)

	if it.p < len(it.node.prefix) {
		if part == it.node.prefix[it.p] {
			// Key matches prefix so far, ok to continue.
			it.p++
			return true
		}
		// Some unmatched prefix remains, node not found
		return false
	}
	node := it.node.getEdge(part)
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
func (it *PathsIterator) Value() (interface{}, bool) {
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
//     ("pre/fix/path", leaf, edges[])
// is split into parent branching node, and a child leaf node:
//     ("pre", nil, edges[f])--->("ix/path", leaf, edges[])
func (node *pathsNode) split(p int) {
	split := &pathsNode{
		edges: node.edges,
		leaf:  node.leaf,
	}
	if p < len(node.prefix)-1 {
		split.prefix = node.prefix[p+1:]
	}
	node.edges = nil
	node.addEdge(pathEdge{node.prefix[p], split})
	if p == 0 {
		node.prefix = nil
	} else {
		node.prefix = node.prefix[:p]
	}
	node.leaf = nil
}

func (node *pathsNode) prune(parents []*pathsNode, links []string) *pathsNode {
	if node.edges != nil {
		return node
	}
	// iterate parents towards root of node, removing the empty leaf
	for i := len(links) - 1; i >= 0; i-- {
		node = parents[i]
		node.delEdge(links[i])
		if len(node.edges) != 0 {
			// parent has other children, stop
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

func (node *pathsNode) compress() {
	if len(node.edges) != 1 || node.leaf != nil {
		return
	}
	for _, edge := range node.edges {
		pfx := make([]string, len(node.prefix)+1+len(edge.node.prefix))
		copy(pfx, node.prefix)
		pfx[len(node.prefix)] = edge.label
		copy(pfx[len(node.prefix)+1:], edge.node.prefix)
		node.prefix = pfx
		node.leaf = edge.node.leaf
		node.edges = edge.node.edges
	}
}

func (node *pathsNode) walk(walkFn WalkFunc) bool {
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

func (node *pathsNode) inspect(pathSep, link, key string, depth int, inspectFn InspectFunc) bool {
	pfx := strings.Join(node.prefix, pathSep)
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
	if node.leaf != nil {
		val = node.leaf.value
	}
	if inspectFn(link, pfx, key, depth, len(node.edges), val) {
		return true
	}
	for _, edge := range node.edges {
		if edge.node.inspect(pathSep, edge.label, key, depth+1, inspectFn) {
			return true
		}
	}
	return false
}

// getEdge binary searches for edge
func (node *pathsNode) getEdge(radix string) *pathsNode {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].label >= radix
	})
	if idx < count && node.edges[idx].label == radix {
		return node.edges[idx].node
	}
	return nil
}

// addEdge binary searches to find where to insert edge, and inserts at
func (node *pathsNode) addEdge(e pathEdge) {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].label >= e.label
	})
	node.edges = append(node.edges, pathEdge{})
	copy(node.edges[idx+1:], node.edges[idx:])
	node.edges[idx] = e
}

// delEdge binary searches for edge and removes it
func (node *pathsNode) delEdge(radix string) {
	count := len(node.edges)
	idx := sort.Search(count, func(i int) bool {
		return node.edges[i].label >= radix
	})
	if idx < count && node.edges[idx].label == radix {
		copy(node.edges[idx:], node.edges[idx+1:])
		node.edges[len(node.edges)-1] = pathEdge{}
		node.edges = node.edges[:len(node.edges)-1]
	}
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
