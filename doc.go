// Package radixtree implements an Adaptive Radix Tree, also called a
// compressed trie or compact prefix tree. Use it to look up values by key,
// find values whose keys share a common prefix, or find values whose keys
// lie along the path to a given key.
//
// The tree uses a radix-256 structure where each key symbol is a byte,
// giving up to 256 branches per node. Nodes hold only as many children
// as needed, keeping memory proportional to the data stored.
//
// Read operations (Get, Iter, IterAt, IterPath) allocate no heap memory
// and are safe to call concurrently. Write operations are not synchronized;
// callers that mix reads and writes must coordinate access themselves.
//
// The API accepts string keys. Because strings are immutable, the tree
// stores them directly without copying.
package radixtree
