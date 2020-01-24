package radixtree

import (
	"strings"
)

// Default path separator is forward slash.  This variable may be set to any
// rune to allow a different path separator.
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

// Len returns the number of values stored in the tree
func (tree *Paths) Len() int {
	if tree.value == nil {
		return 0
	}
	stats := tree.value.(*treeStats)
	return stats.values
}

// Cap returns the total number of nodes, including those without values
func (tree *Paths) Cap() int {
	if tree.value == nil {
		return 0
	}
	stats := tree.value.(*treeStats)
	return stats.nodes
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (tree *Paths) Get(key string) interface{} {
	var p int
	for part, i := pathNext(key, 0); ; part, i = pathNext(key, i) {
		// The tree.prefix represents single-child parents without values that
		// were compressed out of the tree. Let prefix values consume the key.
		if p < len(tree.prefix) {
			if part == tree.prefix[p] {
				p++
				if i == -1 {
					break
				}
				continue
			}
			return nil
		}

		tree = tree.children[part]
		if tree == nil {
			return nil
		}
		if i == -1 {
			break
		}
		p = 0
	}
	return tree.value
}

// GetPath returns all values stored in the path from the root to the node at
// the given key. Does not return values for internal nodes or for nodes with a
// value of nil. Returns a boolean indicating if there was a value stored at
// the full key.
func (tree *Paths) GetPath(key string) ([]interface{}, bool) {
	var values []interface{}
	var p int
	for part, i := pathNext(key, 0); ; part, i = pathNext(key, i) {
		// The tree.prefix represents single-child parents without values that
		// were compressed out of the tree. Let prefix values consume the key.
		if p < len(tree.prefix) {
			if part == tree.prefix[p] {
				p++
				if i == -1 {
					break
				}
				continue
			}
			return values, false
		}

		tree = tree.children[part]
		if tree == nil {
			return values, false
		}
		if tree.value != nil {
			values = append(values, tree.value)
		}
		if i == -1 {
			break
		}
		p = 0
	}
	return values, true
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
		newNodes   int
	)
	node := tree

	for part, next := pathNext(key, 0); ; part, next = pathNext(key, next) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				if next == -1 {
					break
				}
				continue
			}
		} else if child, _ := node.children[part]; child != nil {
			node = child
			p = 0
			if next == -1 {
				break
			}
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
		split := &Paths{
			children: node.children,
			value:    node.value,
		}
		if p < len(node.prefix)-1 {
			split.prefix = node.prefix[p+1:]
		}
		node.children = map[string]*Paths{node.prefix[p]: split}
		if p == 0 {
			node.prefix = nil
		} else {
			node.prefix = node.prefix[:p]
		}
		node.value = nil
		isNewValue = true
		newNodes++
	}

	if newChild != nil {
		// Store key at new child
		if node.children == nil {
			node.children = map[string]*Paths{}
		}
		node.children[childLink] = newChild
		isNewValue = true
		newNodes++
	} else {
		// Store key at existing child
		if node.value == nil {
			// Filled in value of existing internal node
			isNewValue = true
		}
		node.value = value
	}

	// Update stats if any new values or nodes.
	if isNewValue || newNodes != 0 {
		var stats *treeStats
		if tree.value == nil {
			stats = &treeStats{}
			tree.value = stats
		} else {
			stats = tree.value.(*treeStats)
		}
		stats.nodes += newNodes
		if isNewValue {
			stats.values++
		}
	}

	return isNewValue
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key. If the node or any of its ancestors
// becomes childless as a result, it is removed from the tree.
func (tree *Paths) Delete(key string) bool {
	if len(key) == 0 {
		return false
	}
	node := tree
	var nodes []*Paths
	var parts []string
	var p int
	for part, i := pathNext(key, 0); ; part, i = pathNext(key, i) {
		if p < len(node.prefix) {
			if part == node.prefix[p] {
				p++
				if i == -1 {
					break // no more key parts
				}
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
		if i == -1 {
			break // no more key parts
		}
	}

	// Key was not completely consumed traversing tree, so tree does not
	// contain anything for key.
	if p < len(node.prefix) {
		return false
	}
	var deleted bool
	var removed int
	if node.value != nil {
		// delete the node value, indicate that value was deleted
		node.value = nil
		deleted = true
	}

	// If node is leaf, remove from parent.  If parent becomes leaf, repeat.
	if node.children == nil {
		// iterate parents towards root of tree, removing the empty leaf node
		for i := len(parts) - 1; i >= 0; i-- {
			node = nodes[i]
			delete(node.children, parts[i])
			removed++
			if len(node.children) != 0 {
				// parent has other children, stop
				break
			}
			node.children = nil
			if node.value != nil {
				// parent has a value, stop
				break
			}
		}
	}

	// If node has become compressible, compress it
	if len(node.children) == 1 && node.value == nil {
		for part, child := range node.children {
			node.prefix = append(node.prefix, part)
			node.prefix = append(node.prefix, child.prefix...)
			node.value = child.value
			node.children = child.children
			removed++
		}
	}

	// Update stats if anything changed.
	if deleted || removed != 0 {
		stats := tree.value.(*treeStats)
		stats.nodes -= removed
		if deleted {
			stats.values--
		}
	}
	return deleted
}

// Walk walks the tree starting at startKey ("" to start at root), calling
// walkFn for each value found, including at key. If walkFn returns an error,
// the walk is aborted. If walkFn returns Skip, Walk will not descend into the
// node's children.
//
// The tree is traversed depth-first, in no guaranteed order.
func (tree *Paths) Walk(startKey string, walkFn WalkFunc) error {
	if startKey == "" {
		for part, child := range tree.children {
			if err := child.walk(part, walkFn); err != nil {
				return err
			}
		}
		return nil
	}
	// Traverse tree to get to node at key
	var p int
	for part, i := pathNext(startKey, 0); ; part, i = pathNext(startKey, i) {
		if p < len(tree.prefix) && part == tree.prefix[p] {
			p++
			if i == -1 {
				break
			}
			continue
		}
		tree = tree.children[part]
		if tree == nil {
			return nil
		}
		if i == -1 {
			break
		}
		p = 0
	}
	// Walk down tree starting at node located at key
	return tree.walk(startKey, walkFn)
}

func (tree *Paths) walk(key string, walkFn WalkFunc) error {
	if tree.value != nil {
		if err := walkFn(key+strings.Join(tree.prefix, ""), tree.value); err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
	}
	for part, child := range tree.children {
		if err := child.walk(key+part, walkFn); err != nil {
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
	inspectFn("", "", "", 0, len(tree.children), "<root>")
	return tree.inspect("", 1, inspectFn)
}

func (tree *Paths) inspect(key string, depth int, inspectFn InspectFunc) error {
	for part, child := range tree.children {
		pfx := strings.Join(child.prefix, "")
		k := key + part + pfx
		var ik string
		if child.value == nil {
			ik = ""
		} else {
			ik = k
		}
		err := inspectFn(part, pfx, ik, depth, len(child.children), child.value)
		if err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
		if err = child.inspect(k, depth+1, inspectFn); err != nil {
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
