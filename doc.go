/*
Package radixtree implements multiple forms of an Adaptive Radix Tree, aka
compressed trie or prefix tree.  It is adaptive in the sense that nodes are not
constant size, having as few or many children as needed, up to the number of
different key segments to traverse to the next branch or value.

In a compressed radix tree, typically an edge is looked up by radix (single key
symbol) and the edge label contains the key symbols to get to the next prefix
branch or terminal (having a value) node, as well as a pointer the next
node. This implementation puts the edge label in the next node so that there is
not a separate structure for edge and node.  Also, the first character of the
edge label is omitted, since that is already known when descending the tree.

The implementations are optimized for Get performance and allocate 0 bytes of
heap memory per Get; therefore no garbage to collect.  Once the radix tree is
built, it can be repeatedly searched very quickly.

The implementations are optimized for Get performance and allocates 0 bytes of
heap memory per Get; therefore no garbage to collect.  Once the radix tree is
built, it can be repeatedly searched quickly. Concurrent searches are safe
since these do not modify the radixtree. Access is not synchronized (not
concurrent safe with writes), allowing the caller to synchronize, if needed, in
whatever manner works best for the application.

*/
package radixtree
