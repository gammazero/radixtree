package radixtree

// WalkFunc is the type of the function called for each value visited by Walk
// or WalkPath.  The key argument contains the elements of the key at which the
// value is stored.
//
// If the function returns true Walk stops immediately and returns.  This
// applies to WalkPath as well.
type WalkFunc func(key string, value interface{}) bool

// InspectFunc is the type of the function called for each node visited by
// Inspect.  The key argument contains the key at which the node is located,
// the depth is the distance from the root of the tree, and children is the
// number of children the node has.
//
// If the function returns true Inspect stops immediately and returns.
type InspectFunc func(link, prefix, key string, depth, children int, hasValue bool, value interface{}) bool

// Iterator iterates all keys and values in the radixtree
//
// Is is safe to use different iterators concurrently.  Any modification to the
// tree that the iterator was created from invalidates the iterator.
type Iterator interface {
	// Next returns the next key and value stored in the tree, and true when
	// iteration is complete
	Next() (key string, value interface{}, more bool)
	// Copy creates a new Iterator at this iterator's state of iteration
	Copy() Iterator
}
