package radixtree

// Stepper traverses a Tree one byte at a time.
//
// Any modification to the tree invalidates the Stepper.
type Stepper struct {
	p    int
	node *radixNode
}

// NewStepper returns a new Stepper instance that begins at the root of the
// tree.
func (t *Tree) NewStepper() *Stepper {
	return &Stepper{
		node: &t.root,
	}
}

// Copy makes a copy of the current Stepper. This allows branching a Stepper
// into two that can take separate paths. These Steppers do not affect each
// other and can be used concurrently.
func (s *Stepper) Copy() *Stepper {
	return &Stepper{
		p:    s.p,
		node: s.node,
	}
}

// Next advances the Stepper from its current position, to the position of
// given key symbol in the tree, so long as the given symbol is next in a path
// in the tree. If the symbol allows the Stepper to advance, then true is
// returned. Otherwise false is returned.
//
// When false is returned the Stepper is not modified. This allows different
// values to be used in subsequent calls to Next.
func (s *Stepper) Next(radix byte) bool {
	// The tree.prefix represents single-edge parents without values that were
	// compressed out of the tree. Let prefix consume key symbols.
	if s.p < len(s.node.prefix) {
		if radix == s.node.prefix[s.p] {
			// Key matches prefix so far, ok to continue.
			s.p++
			return true
		}
		// Some unmatched prefix remains, node not found.
		return false
	}
	node := s.node.getEdge(radix)
	if node == nil {
		// No more prefix, no edges, so node not found.
		return false
	}
	// Key symbol matched up to this edge, ok to continue.
	s.p = 0
	s.node = node
	return true
}

// Item returns an Item containing the key and value at the current Stepper
// position, or returns nil if no value is present at the position.
func (s *Stepper) Item() *Item {
	// Only return item if all of this node's prefix was matched. Otherwise,
	// have not fully traversed into this node (edge not completely traversed).
	if s.p == len(s.node.prefix) {
		return s.node.leaf
	}
	return nil
}

// Value returns the value at the current Stepper position, and true or false
// to indicate if a value is present at the position.
func (s *Stepper) Value() (any, bool) {
	item := s.Item()
	if item == nil {
		return nil, false
	}
	return item.value, true
}
