/*
Package radixtree implements multiple forms of an Adaptive Radix Tree, aka
compressed trie or prefix tree.  It is adaptive in the sense that nodes are not
constant size, having as few or many children as needed, up to the number of
different key segments to traverse to the next branch or value.

This package provides multiple types of radix tree: Bytes and Paths, where each
key symbol (radix) is a byte or a path segment, respectively.

The implementations are optimized for Get performance and allocates 0 bytes of
heap memory per Get; therefore no garbage to collect.  Once the radix tree is
built, it can be repeatedly searched quickly. Concurrent searches are safe
since these do not modify the radixtree. Access is not synchronized (not
concurrent safe with writes), allowing the caller to synchronize, if needed, in
whatever manner works best for the application.

*/
package radixtree
