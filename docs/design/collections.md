Design: Collections
===
- [Design: Collections](#design-collections)
  - [Assumptions](#assumptions)


## Assumptions
- A collection is a node of nodes. It can be cyclic, but directed. Basically a chunk of data that represent one to many files.
- A collection represents one OCI artifact that can contain one to many descriptors that reference eachother.
- A collection must be rooted meaning it must contain a URO.
- Collections are build with the builder. They may have a "next" annotation is symbolize a connection to another collection or node.
- All attributes from child nodes are aggregated to parent nodes to allow for greedy BFS.