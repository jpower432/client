Design: Collections
===
- [Design: Collections](#design-collections)
  - [Assumptions](#assumptions)


## Assumptions
- A collection is a node of nodes. This would be a directed graph with cycle detection.
- A collection represents one OCI artifact that can contain one to many descriptors that reference each other.
- Collections are built with the builder. They may have a "next" annotation is symbolize a connection to another collection or node.
- All attributes from child nodes are aggregated to parent nodes to allow for greedy BFS searching.