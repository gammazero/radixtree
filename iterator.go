package radixtree

// Iterator iterates all keys and values in the radix tree.
//
// Is is safe to use different Iterator instances concurrently. Any
// modification to the Tree that the Iterator was created from invalidates the
// Iterator instance.
type Iterator struct {
	nodes []*radixNode
}

// NewIterator returns a new Iterator.
func (t *Tree) NewIterator() *Iterator {
	return &Iterator{
		nodes: []*radixNode{&t.root},
	}
}

// Copy creates a new Iterator at this iterator's state of iteration.
func (it *Iterator) Copy() *Iterator {
	nodes := make([]*radixNode, len(it.nodes))
	copy(nodes, it.nodes)
	return &Iterator{
		nodes: nodes,
	}
}

// Next returns the next key and value stored in the Tree, and true when
// iteration is complete.
func (it *Iterator) Next() (key string, value any, done bool) {
	for {
		if len(it.nodes) == 0 {
			break
		}
		node := it.nodes[len(it.nodes)-1]
		it.nodes = it.nodes[:len(it.nodes)-1]

		for i := len(node.edges) - 1; i >= 0; i-- {
			it.nodes = append(it.nodes, node.edges[i].node)
		}

		if node.leaf != nil {
			return node.leaf.key, node.leaf.value, false
		}
	}
	return "", nil, true
}
