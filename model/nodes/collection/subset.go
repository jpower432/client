package collection

import "github.com/uor-framework/client/model"

// EdgeFunc defines whether to add
// the edge to the subset of data.
type EdgeFunc func(edge model.Edge) bool

// NodeFunc defines whether to add the
// node to the subset of data.
type NodeFunc func(node model.Node) bool

// EdgeSubgraph returns the directed subgraph with only the edges that match the
// provided function.
func (c *Collection) EdgeSubgraph(id string, edgeFn EdgeFunc) Collection {
	out := NewCollection(id)
	for _, node := range c.Nodes() {
		out.AddNode(node)
	}
	out.addEdges(c.Edges(), edgeFn)
	return *out
}

// Subset returns a subset of the Collection with only the nodes and edges that match the
// provided functions.
func (c *Collection) Subset(id string, nodeFn NodeFunc, edgeFn EdgeFunc) Collection {
	out := NewCollection(id)
	for _, node := range c.Nodes() {
		if nodeFn(node) {
			out.AddNode(node)
		}
	}
	out.addEdges(c.Edges(), edgeFn)
	return *out
}

// SubsetWithNodes returns a subset of the collection with only the listed nodes and edges that
// match the provided function.
func (c *Collection) SubsetWithNodes(id string, nodes []model.Node, fn EdgeFunc) Collection {
	out := NewCollection(id)
	for _, node := range nodes {
		out.AddNode(node)
	}
	out.addEdges(c.Edges(), fn)
	return *out
}

// addEdges adds the specified edges, filtered by the provided edge connection
// function.
func (c *Collection) addEdges(edges []model.Edge, fn EdgeFunc) error {
	for _, e := range edges {
		if !fn(e) {
			continue
		}
		if err := c.AddEdge(e); err != nil {
			return err
		}
	}
	return nil
}
