package radixtree

// Runes is a radix tree of runes with string keys and interface{} values.
// Non-terminal nodes have nil values so a stored nil value is no
// distinguishable and is not be included in results from GetPath or Walk.
type Runes struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix   []rune
	value    interface{}
	children map[rune]*Runes
}

// Len returns the number of values stored in the tree
func (tree *Runes) Len() int {
	if tree.value == nil {
		return 0
	}
	stats := tree.value.(*treeStats)
	return stats.values
}

// Cap returns the total number of nodes, including those without values
func (tree *Runes) Cap() int {
	if tree.value == nil {
		return 0
	}
	stats := tree.value.(*treeStats)
	return stats.nodes
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (tree *Runes) Get(key string) interface{} {
	var p int
	for _, r := range key {
		// The tree.prefix represents single-child parents without values that
		// were compressed out of the tree. Let prefix values consume the key.
		if p < len(tree.prefix) {
			if r == tree.prefix[p] {
				p++
				continue
			}
			// Some unmatched prefix remains, node not found
			return nil
		}

		tree = tree.children[r]
		if tree == nil {
			// No more prefix, no children, so node not found
			return nil
		}

		p = 0
	}
	// Key has been consumed by traversing prefixes and/or children.  If key
	// did not match all of this node's prefix, then did not find value.
	if p < len(tree.prefix) {
		return nil
	}
	return tree.value
}

// GetPath returns all values stored in the path from the root to the node at
// the given key. Does not return values for internal nodes or for nodes with a
// value of nil. Returns a boolean indicating if there was a value stored at
// the full key.
func (tree *Runes) GetPath(key string) ([]interface{}, bool) {
	var values []interface{}
	var p int
	for _, r := range key {
		if p < len(tree.prefix) {
			if r == tree.prefix[p] {
				p++
				continue
			}
			return values, false
		}

		tree = tree.children[r]
		if tree == nil {
			return values, false
		}

		if tree.value != nil {
			values = append(values, tree.value)
		}
		p = 0
	}
	// If key did not match all of this node's prefix, then did not find value.
	return values, p == len(tree.prefix)
}

// Put inserts the value into the tree at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
//
// Note that internal nodes have nil values so a stored nil value is not
// distinguishable and is not included in Walks.
func (tree *Runes) Put(key string, value interface{}) bool {
	var (
		p          int
		childLink  rune
		newChild   *Runes
		isNewValue bool
		newNodes   int
	)
	node := tree

	for i, r := range key {
		if p < len(node.prefix) {
			if r == node.prefix[p] {
				p++
				continue
			}
		} else if child, _ := node.children[r]; child != nil {
			node = child
			p = 0
			continue
		}
		// Descended as far as prefixes and children match key, and still
		// have key data, so add child that has a prefix of the unmatched
		// key data and set its value to the new value.
		newChild = &Runes{
			value: value,
		}
		if i < len(key)-1 {
			newChild.prefix = []rune(key[i+1:])
		} else {
		}
		childLink = r
		value = nil // value stored in newChild, not in node
		break
	}
	// Key has been consumed by traversing prefixes and/or children, or has
	// been put into new child.

	// If key partially matches node's prefix, then need to split node.
	if p < len(node.prefix) {
		split := &Runes{
			children: node.children,
			value:    node.value,
		}
		if p < len(node.prefix)-1 {
			split.prefix = node.prefix[p+1:]
		}
		node.children = map[rune]*Runes{node.prefix[p]: split}
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
		if node.children == nil {
			node.children = map[rune]*Runes{}
		}
		node.children[childLink] = newChild
		isNewValue = true
		newNodes++
	} else {
		if node.value == nil {
			// Filled in value of existing internal node
			isNewValue = true
		}
		node.value = value
	}

	// Update stats is any new values or nodes.
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
// node was found for the given key, and that node had a non-nil value. If the
// node has no children (is a leaf node) it is removed from the tree. If any of
// of the node's ancestors becomes childless as a result, they are also removed
// from the tree.
func (tree *Runes) Delete(key string) bool {
	if len(key) == 0 {
		return false
	}
	node := tree
	var nodes []*Runes
	var runes []rune
	var p int
	for _, r := range key {
		if p < len(node.prefix) {
			if r == node.prefix[p] {
				p++
				continue
			}
			return false
		}
		nodes = append(nodes, node)
		runes = append(runes, r)
		node = node.children[r]
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
	var removed int
	if node.value != nil {
		// delete the node value, indicate that value was deleted
		node.value = nil
		deleted = true
	}

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	if node.children == nil {
		// iterate parents towards root of tree, removine the empty leaf node
		for i := len(runes) - 1; i >= 0; i-- {
			node = nodes[i]
			delete(node.children, runes[i])
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
		for r, child := range node.children {
			node.prefix = append(node.prefix, r)
			node.prefix = append(node.prefix, child.prefix...)
			node.value = child.value
			node.children = child.children
			removed++
		}
	}

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
func (tree *Runes) Walk(startKey string, walkFn WalkFunc) error {
	// Traverse tree to get to node at key
	var p int
	for _, r := range startKey {
		if p < len(tree.prefix) && r == tree.prefix[p] {
			p++
			continue
		}
		tree = tree.children[r]
		if tree == nil {
			return nil
		}
		p = 0
	}
	// Walk down tree starting at node located at key
	return tree.walk(startKey, walkFn)
}

func (tree *Runes) walk(key string, walkFn WalkFunc) error {
	if len(key) != 0 && tree.value != nil {
		if err := walkFn(key, tree.value); err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
	}
	for r, child := range tree.children {
		k := key + string(r) + string(child.prefix)
		err := child.walk(k, walkFn)
		if err != nil {
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
func (tree *Runes) Inspect(inspectFn InspectFunc) error {
	inspectFn("", "", "", 0, len(tree.children), "<root>")
	return tree.inspect("", 1, inspectFn)
}

func (tree *Runes) inspect(key string, depth int, inspectFn InspectFunc) error {
	for r, child := range tree.children {
		k := key + string(r) + string(child.prefix)
		var ik string
		if child.value == nil {
			ik = ""
		} else {
			ik = k
		}
		if err := inspectFn(string(r), string(child.prefix), ik, depth, len(child.children), child.value); err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
		if err := child.inspect(k, depth+1, inspectFn); err != nil {
			return err
		}
	}
	return nil
}
