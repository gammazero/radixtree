package radixtree

// Iterator iterates all keys and values in the radix tree.
//
// Is is safe to use different Iterator instances concurrently. Any
// modification to the Tree that the Iterator was created from invalidates the
// Iterator instance.
type Iterator[V any] struct {
	nodes []*radixNode[V]
}

// NewIterator returns a new Iterator.
func (t *Tree[V]) NewIterator() *Iterator[V] {
	return &Iterator[V]{
		nodes: []*radixNode[V]{&t.root},
	}
}

// Copy creates a new Iterator at this iterator's state of iteration.
func (it *Iterator[V]) Copy() *Iterator[V] {
	nodes := make([]*radixNode[V], len(it.nodes))
	copy(nodes, it.nodes)
	return &Iterator[V]{
		nodes: nodes,
	}
}

// Next returns the next key and value stored in the Tree, and true when
// iteration is complete.
func (it *Iterator[V]) Next() (key string, value V, done bool) {
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
	var zeroV V
	return "", zeroV, true
}
