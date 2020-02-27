package radixtree

import "errors"

// Skip is a special error return from WalkFunc and WalkPathFunc
var Skip = errors.New("skip")

// KeyStringer is the interface to a key element accumulator, used to get the
// key of the current node during Walk.
type KeyStringer interface {
	// String returns the string form of key elements accumulated during walk.
	String() string
}

// WalkFunc is the type of the function called for each value visited by Walk.
// The key argument contains the elements of the key at which the value is
// stored, available as a string by key.String().
//
// If an error is returned, processing of Walk stops.  The sole exception is
// when the function returns the special value Skip.  When the function returns
// Skip, Walk will not descend into any children of the current node.
type WalkFunc func(key KeyStringer, value interface{}) error

// WalkPathFunc is the type of the function called for each value visited by WalkPath.
// The key argument is the key at which the value is stored.
//
// If an error is returned, processing of WalkPath stops.  The sole exception is
// when the function returns the special value Skip.  When the function returns
// Skip, WalkPath stops processing and returns a nil error.
type WalkPathFunc func(key string, value interface{}) error

// InspectFunc is the type of the function called for each node visited by
// Inspect.  The key argument contains the key at which the node is located,
// the depth is the distance from the root of the tree, and children is the
// number of children the node has.
//
// If an error is returned, processing of Inspect stops.  The sole exception is
// when the function returns the special value Skip.  When the function returns
// Skip, Inspect will not descend into any children of the current node.
type InspectFunc func(link, prefix, key string, depth, children int, value interface{}) error
