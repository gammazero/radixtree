/*
Package radixtree implements an Adaptive Radix Tree, aka compressed trie or
compact prefix tree. It is adaptive in the sense that nodes are not constant
size, having as few or many children as needed to branch to all subtrees.

This package implements a radix-256 tree where each key symbol (radix) is a
byte, allowing up to 256 possible branches to traverse to the next node.

The implementation is optimized for Get performance and allocates 0 bytes of
heap memory per Get; therefore no garbage to collect. Once the radix tree is
built, it can be repeatedly searched quickly. Concurrent searches are safe
since these do not modify the radixtree. Access is not synchronized (not
concurrent safe with writes), allowing the caller to synchronize, if needed, in
whatever manner works best for the application.

The API uses string keys, since strings are immutable and therefore it is not
necessary make a copy of the key provided to the radix tree.
*/
package radixtree
