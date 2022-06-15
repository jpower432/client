Design: Model
===
- [Design: Model](#design-model)


The model package and the sub-package contains all types and methods that can be used to defines and work with UOR data.

- Tree: This structure defines relationship between different UOR node types. This structure would be an implementation of a B-Tree since this can contain one of more UOR collection. The constraint is that one UOR collection can fit into main memory, but multiple collections can result in huge amount of data that would require disk access.
- Iterator: Nodes can be iterable (e.g. a UOR collection). Using the iterator interface allows these structures to be iterated over during tree traversal.
- Matcher: Defines criteria for node searching in a tree or subtree.



