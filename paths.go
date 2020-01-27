package radixtree

import (
	"strings"
)

// PathSeparator splits a path key into separate segments.  The default path
// separator is forward slash.  This variable may be set to any rune to allow a
// different path separator.
var PathSeparator = '/'

// Paths is a radix tree of paths with string keys and interface{}
// values. Paths splits keys by separator, for example, "/a/b/c" splits to
// "/a", "/b", "/c". A different path separator character may be specified by
// setting PathSeparator.
//
// Non-terminal nodes have nil values so a stored nil value is no
// distinguishable and is not be included in results from GetPath or Walk.
type Paths struct {
	prefix   []string
	value    interface{}
	children map[string]*Paths
}

// PathsIterator traverses a Paths radix tree one path segment at a time.
//
// Note: Any modification to the tree invalidates the iterator.
type PathsIterator struct {
	p    int
	node *Paths
}

// NewIterator returns a new PathsIterator instance that begins iterating
// from the root of the tree.
func (tree *Paths) NewIterator() *PathsIterator {
	return &PathsIterator{
		node: tree,
	}
}

// Copy makes a copy of the current iterator.  This allows branching an
// iterator into two iterators that can take separate paths.
func (it *PathsIterator) Copy() *PathsIterator {
	return &PathsIterator{
		p:    it.p,
		node: it.node,
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
func (pi *PathsIterator) Next(part string) bool {
	if pi.p < len(pi.node.prefix) {
		if part == pi.node.prefix[pi.p] {
			pi.p++
			return true
		}
		return false
	}
	node := pi.node.children[part]
	if node == nil {
		return false
	}
	pi.p = 0
	pi.node = node
	return true
}

// Value returns the value at the current iterator position, or nil if there is
// no value at the position.
func (pi *PathsIterator) Value() interface{} {
	if pi.p != len(pi.node.prefix) {
		return nil
	}
	return pi.node.value
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (tree *Paths) Get(key string) interface{} {
	iter := tree.NewIterator()
	for part, i := pathNext(key, 0); part != ""; part, i = pathNext(key, i) {
		if !iter.Next(part) {
			return nil
		}
	}
	return iter.Value()
}

// GetPath returns all values stored in the path from the root to the node at
// the given key. Does not return values for internal nodes or for nodes with a
// value of nil. Returns a boolean indicating if there was a value stored at
// the full key.
func (tree *Paths) GetPath(key string) ([]interface{}, bool) {
	var values []interface{}
	var value interface{}
	iter := tree.NewIterator()
	for part, i := pathNext(key, 0); part != ""; part, i = pathNext(key, i) {
		if !iter.Next(part) {
			return values, false
		}
		value = iter.Value()
		if value != nil {
			values = append(values, value)
		}
	}
	return values, value != nil
}

// Put inserts the value into the tree at the given key, replacing any existing
// items.  It returns true if the put adds a new value, false if it replaces an
// existing value.
//
// Note that internal nodes have nil values so a stored nil value is not
// distinguishable and is not included in Walks.
func (tree *Paths) Put(key string, value interface{}) bool {
	var (
		p          int
		childLink  string
		newChild   *Paths
		isNewValue bool
	)
	node := tree

	for part, next := pathNext(key, 0); part != ""; part, next = pathNext(key, next) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				continue
			}
		} else if child, _ := node.children[part]; child != nil {
			node = child
			p = 0
			continue
		}
		// Descended as far as prefixes and children match key, and still
		// have key data, so add child that has a prefix of the unmatched
		// key data and set its value to the new value.
		newChild = &Paths{
			value: value,
		}
		childLink = part
		if next != -1 {
			newChild.prefix = []string{}
			for next != -1 {
				part, next = pathNext(key, next)
				newChild.prefix = append(newChild.prefix, part)
			}
		}
		value = nil // value stored in newChild, not in node
		break
	}
	// Key has been consumed by traversing prefixes and/or children, or has
	// been put into new child.

	// If key partially matches node's prefix, then need to split node.
	if p < len(node.prefix) {
		node.split(p)
		isNewValue = true
	}

	if newChild != nil {
		// Store key at new child
		if node.children == nil {
			node.children = map[string]*Paths{}
		}
		node.children[childLink] = newChild
		isNewValue = true
	} else {
		// Store key at existing child
		if node.value == nil {
			isNewValue = true
		}
		node.value = value
	}

	return isNewValue
}

// split splits a node such that a node:
//     ("pre/fix", "value", child[..])
// is split into parent branching node, and a child value node:
//     ("pre/", "", [-])--->("fix/", "value", [..])
func (tree *Paths) split(p int) {
	split := &Paths{
		children: tree.children,
		value:    tree.value,
	}
	if p < len(tree.prefix)-1 {
		split.prefix = tree.prefix[p+1:]
	}
	tree.children = map[string]*Paths{tree.prefix[p]: split}
	if p == 0 {
		tree.prefix = nil
	} else {
		tree.prefix = tree.prefix[:p]
	}
	tree.value = nil
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
	for part, i := pathNext(key, 0); part != ""; part, i = pathNext(key, i) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				continue
			}
			return false
		}
		nodes = append(nodes, node)
		parts = append(parts, part)
		node = node.children[part]
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
	if node.value != nil {
		// delete the node value, indicate that value was deleted
		node.value = nil
		deleted = true
	}

	// If node is leaf, remove from parent.  If parent becomes leaf, repeat.
	node = node.prune(nodes, parts)

	// If node has become compressible, compress it
	node.compress()

	return deleted
}

func (tree *Paths) prune(parents []*Paths, links []string) *Paths {
	if tree.children != nil {
		return tree
	}
	// iterate parents towards root of tree, removing the empty leaf
	for i := len(links) - 1; i >= 0; i-- {
		tree = parents[i]
		delete(tree.children, links[i])
		if len(tree.children) != 0 {
			// parent has other children, stop
			break
		}
		tree.children = nil
		if tree.value != nil {
			// parent has a value, stop
			break
		}
	}
	return tree
}

func (tree *Paths) compress() {
	if len(tree.children) != 1 || tree.value != nil {
		return
	}
	for part, child := range tree.children {
		tree.prefix = append(tree.prefix, part)
		tree.prefix = append(tree.prefix, child.prefix...)
		tree.value = child.value
		tree.children = child.children
	}
}

// Walk walks the radix tree rooted at root ("" to start at root or tree),
// calling walkFn for each value found. If walkFn returns an error, the walk is
// aborted. If walkFn returns Skip, Walk will not descend into the node's
// children.
//
// The tree is traversed depth-first, in no guaranteed order.
func (tree *Paths) Walk(root string, walkFn WalkFunc) error {
	// Traverse tree to get to node at key
	if root != "" {
		iter := tree.NewIterator()
		for part, i := pathNext(root, 0); part != ""; part, i = pathNext(root, i) {
			if !iter.Next(part) {
				return nil
			}
		}
		// Root is not valid unless a value was stored there
		if iter.Value() == nil {
			return nil
		}
		tree = iter.node
	}
	// Walk down tree starting at node located at root
	return tree.walk(root, walkFn)
}

func (tree *Paths) walk(key string, walkFn WalkFunc) error {
	if tree.value != nil {
		if err := walkFn(key, tree.value); err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
	}
	for part, child := range tree.children {
		if err := child.walk(key+part+strings.Join(child.prefix, ""), walkFn); err != nil {
			return err
		}
	}
	return nil
}

// Inspect walks every node of the tree, whether or not it holds a value,
// calling inspectFn with information about each node.  This allows the
// structure of the tree to be examined and detailed statistics to be
// collected.
//
// If inspectFn returns an error, the traversal is aborted.  If inspectFn
// returns Skip, Inspect will not descend into the node's children.
//
// The tree is traversed depth-first, in no guaranteed order.
func (tree *Paths) Inspect(inspectFn InspectFunc) error {
	return tree.inspect("", "", 0, inspectFn)
}

func (tree *Paths) inspect(link, key string, depth int, inspectFn InspectFunc) error {
	pfx := strings.Join(tree.prefix, "")
	key += link + pfx
	err := inspectFn(link, pfx, key, depth, len(tree.children), tree.value)
	if err != nil {
		if err == Skip {
			// Ignore current node's children.
			return nil
		}
		return err
	}
	for part, child := range tree.children {
		if err = child.inspect(part, key, depth+1, inspectFn); err != nil {
			return err
		}
	}
	return nil
}

// pathNext splits path strings by a path separator character. For
// example, "/a/b/c" -> ("/a", 2), ("/b", 4), ("/c", -1) in successive
// calls. It does not allocate any heap memory.
func pathNext(path string, start int) (string, int) {
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}
	end := strings.IndexRune(path[start+1:], PathSeparator)
	if end == -1 {
		return path[start:], -1
	}
	next := start + end + 1
	return path[start:next], next
}
