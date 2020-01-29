package radixtree

// Runes is a radix tree of runes with string keys and interface{} values.
// Non-terminal nodes have nil values, so a stored nil value is not
// distinguishable and is not included in Walk or WalkPath.
type Runes struct {
	// prefix is the edge label between this node and the parent, minus the key
	// segment used in the parent to index this child.
	prefix   []rune
	value    interface{}
	children map[rune]*Runes
}

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
// iterator into two iterators that can take separate paths.
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
func (ri *RunesIterator) Next(r rune) bool {
	// The tree.prefix represents single-child parents without values that were
	// compressed out of the tree.  Let prefix consume key symbols.
	if ri.p < len(ri.node.prefix) {
		if r == ri.node.prefix[ri.p] {
			// Key matches prefix so far, ok to continue.
			ri.p++
			return true
		}
		// Some unmatched prefix remains, node not found
		return false
	}
	node := ri.node.children[r]
	if node == nil {
		// No more prefix, no children, so node not found
		return false
	}
	// Key symbol matched up to this child, ok to continue.
	ri.p = 0
	ri.node = node
	return true
}

// Value returns the value at the current iterator position, or nil if there is
// no value at the position.
func (ri *RunesIterator) Value() interface{} {
	// Only return value if all of this node's prefix was matched.  Otherwise,
	// have not fully traversed into this node (edge not completely traversed).
	if ri.p != len(ri.node.prefix) {
		return nil
	}
	return ri.node.value
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (tree *Runes) Get(key string) interface{} {
	iter := tree.NewIterator()
	for _, r := range key {
		if !iter.Next(r) {
			return nil
		}
	}
	return iter.Value()
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
		}
		childLink = r
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
		if node.children == nil {
			node.children = map[rune]*Runes{}
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
//     ("prefix", "value", child[..])
// is split into parent branching node, and a child value node:
//     ("pre", "", [-])--->("fix", "value", [..])
func (tree *Runes) split(p int) {
	split := &Runes{
		children: tree.children,
		value:    tree.value,
	}
	if p < len(tree.prefix)-1 {
		split.prefix = tree.prefix[p+1:]
	}
	tree.children = map[rune]*Runes{tree.prefix[p]: split}
	if p == 0 {
		tree.prefix = nil
	} else {
		tree.prefix = tree.prefix[:p]
	}
	tree.value = nil
}

// Delete removes the value associated with the given key. Returns true if a
// node was found for the given key, and that node had a non-nil value. If the
// node has no children (is a leaf node) it is removed from the tree. If any of
// of the node's ancestors becomes childless as a result, they are also removed
// from the tree.
func (tree *Runes) Delete(key string) bool {
	node := tree
	var (
		nodes []*Runes
		runes []rune
		p     int
	)
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
	if node.value != nil {
		// delete the node value, indicate that value was deleted
		node.value = nil
		deleted = true
	}

	// If node is leaf, remove from parent. If parent becomes leaf, repeat.
	node = node.prune(nodes, runes)

	// If node has become compressible, compress it
	node.compress()

	return deleted
}

func (tree *Runes) prune(parents []*Runes, links []rune) *Runes {
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

func (tree *Runes) compress() {
	if len(tree.children) != 1 || tree.value != nil {
		return
	}
	for r, child := range tree.children {
		tree.prefix = append(tree.prefix, r)
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
func (tree *Runes) Walk(root string, walkFn WalkFunc) error {
	if root != "" {
		iter := tree.NewIterator()
		// Traverse tree to get to node at key
		for _, r := range root {
			if !iter.Next(r) {
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

func (tree *Runes) walk(key string, walkFn WalkFunc) error {
	if tree.value != nil {
		if err := walkFn(key, tree.value); err != nil {
			if err == Skip {
				// Ignore current node's children.
				return nil
			}
			return err
		}
	}
	for r, child := range tree.children {
		if err := child.walk(key+string(r)+string(child.prefix), walkFn); err != nil {
			return err
		}
	}
	return nil
}

// WalkPath walks the path in tree from the root to the node at the given key,
// calling walkFn for each node that has a value.
func (tree *Runes) WalkPath(key string, walkFn WalkFunc) error {
	if tree.value != nil {
		if err := walkFn("", tree.value); err != nil {
			if err == Skip {
				return nil
			}
			return err
		}
	}
	iter := tree.NewIterator()
	for i, r := range key {
		if !iter.Next(r) {
			return nil
		}
		value := iter.Value()
		if value != nil {
			err := walkFn(string(key[0:i+1]), value)
			if err != nil {
				if err == Skip {
					return nil
				}
				return err
			}
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
	return tree.inspect("", "", 0, inspectFn)
}

func (tree *Runes) inspect(link, key string, depth int, inspectFn InspectFunc) error {
	pfx := string(tree.prefix)
	key += link + pfx
	err := inspectFn(link, pfx, key, depth, len(tree.children), tree.value)
	if err != nil {
		if err == Skip {
			// Ignore current node's children.
			return nil
		}
		return err
	}
	for r, child := range tree.children {
		if err = child.inspect(string(r), key, depth+1, inspectFn); err != nil {
			return err
		}
	}
	return nil
}
