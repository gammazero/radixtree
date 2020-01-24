package radixtree

import "errors"

// Skip is a special error return
var Skip = errors.New("skip")

// WalkFunc is the type of the function called for each value visited by
// Walk.  The key argument contains the key at which the value is stored.
//
// If an error is returned, processing of Walk stops.  The sole exception is
// when the function returns the special value Skip.  When the function returns
// Skip, Walk will not descend into any children of the current node.
type WalkFunc func(startKey string, value interface{}) error

// InspectFunc is the type of the function called for each node visited by
// Inspect.  The key argument contains the key at which the node is located,
// the depth is the distance from the root of the tree, and children is the
// number of children the node has.
//
// If an error is returned, processing of Inspect stops.  The sole exception is
// when the function returns the special value Skip.  When the function returns
// Skip, Inspect will not descend into any children of the current node.
type InspectFunc func(link, prefix, key string, depth, children int, value interface{}) error
