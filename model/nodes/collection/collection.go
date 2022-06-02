package collection

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

var (
	_ model.Node   = &Collection{}
	_ model.Rooted = &Collection{}
)

// Collection represents a subtree and can represent on OCI artifact.
type Collection struct {
	id       string
	Location string
	// nodes describes all nodes contained in the graph
	nodes map[string]model.Node
	// from describes all edges with the
	// origin node as the map key.
	from map[string]map[string]model.Edge
	// to describes all edges with the
	// destination node as the map key
	to map[string]map[string]model.Edge
}

// NewGraph creates an empty Graph.
func NewCollection(id string) *Collection {
	return &Collection{
		id:    id,
		nodes: map[string]model.Node{},
		from:  map[string]map[string]model.Edge{},
		to:    map[string]map[string]model.Edge{},
	}
}

// ID return the unique id of the collection.
func (c *Collection) ID() string {
	return c.id
}

// Address returns collection location.
func (c *Collection) Address() string {
	return c.Location
}

// Attributes returns a collection of all the
// attributes contained within the collection nodes.
func (c *Collection) Attributes() map[string]string {
	attributes := map[string]string{}
	for _, node := range c.Nodes() {
		for k, v := range node.Attributes() {
			attributes[k] = v
		}
	}
	return attributes
}

// Node returns the node based on the ID if the node exists.
func (c *Collection) Node(id string) model.Node {
	node, ok := c.nodes[id]
	if !ok {
		return nil
	}
	return node
}

// Nodes returns a slice containing
// all nodes in the graph.
func (c *Collection) Nodes() []model.Node {
	var nodes []model.Node
	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// Edges returns a slice containing
// all nodes in the graph.
func (c *Collection) Edges() []model.Edge {
	var edges []model.Edge
	for _, to := range c.from {
		for _, edge := range to {
			edges = append(edges, edge)
		}
	}
	return edges
}

// Edge returns the edge from the origin to the destination node if such an edge exists.
// The node from must be directly reachable from the node to as defined by the From method.
func (c *Collection) Edge(from, to string) model.Edge {
	edge, ok := c.from[from][to]
	if !ok {
		return nil
	}
	return edge
}

// HasEdgeFromTo returns whether there is an edge
// from the origin to the destination Node.
func (c *Collection) HasEdgeFromTo(from, to string) bool {
	_, ok := c.from[to][from]
	return ok
}

// HasEdgeBetween returns whether there is an edge
// from the destination to the origin Node.
func (c *Collection) HasEdgeBetween(from, to string) bool {
	if _, ok := c.to[to][from]; ok {
		return true
	}
	if _, ok := c.from[from][to]; ok {
		return true
	}
	return false
}

// From returns a list of Nodes connected
// to the node with the id.
func (c *Collection) From(id string) []model.Node {
	var connectedNodes []model.Node
	nodes, ok := c.from[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, c.nodes[id])
	}
	return connectedNodes
}

// To returns a list of Nodes connected
// to the node with the id.
func (c *Collection) To(id string) []model.Node {
	var connectedNodes []model.Node
	nodes, ok := c.to[id]
	if !ok {
		return nil
	}
	for id := range nodes {
		connectedNodes = append(connectedNodes, c.nodes[id])
	}

	return connectedNodes
}

// Root calculates to root node of the graph.
// This is calculated base on existing child nodes.
// This expects only one root node to be found.
func (c *Collection) Root() (model.Node, error) {
	childNodes := map[string]int{}
	for _, n := range c.nodes {
		for _, ch := range c.From(n.ID()) {
			childNodes[ch.ID()]++
		}
	}
	var roots []model.Node
	for _, n := range c.nodes {
		if _, found := childNodes[n.ID()]; !found {
			roots = append(roots, n)
		}
	}
	if len(roots) == 0 {
		return nil, fmt.Errorf("no root found in graph")
	}
	if len(roots) > 1 {
		var rootNames []string
		for _, root := range roots {
			rootNames = append(rootNames, root.Address())
		}
		sort.Strings(rootNames)
		return nil, fmt.Errorf("multiple roots found in graph: %s", strings.Join(rootNames, ", "))
	}
	return roots[0], nil
}
