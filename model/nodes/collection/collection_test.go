package collection

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/model"
)

func TestCollection_Root(t *testing.T) {
	type spec struct {
		name     string
		nodes    []model.Node
		edges    []model.Edge
		expID    string
		expError string
	}

	cases := []spec{
		{
			name:  "Success/RootExists",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			expID: "node3",
		},
		{
			name:  "Failure/NotRootExists",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
				&Edge{T: &mockNode{id: "node3"}, F: &mockNode{id: "node2"}},
			},
			expError: "no root found in graph",
		},
		{
			name:  "Failure/MultipleRootsExist",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
			},
			expError: "multiple roots found in graph: address, address",
		},
		{
			name:     "Failure/NoEdges",
			nodes:    []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges:    nil,
			expError: "multiple roots found in graph: address, address, address",
		},
		{
			name:     "Failure/NoNodes",
			nodes:    nil,
			edges:    nil,
			expError: "no root found in graph",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			root, err := collection.Root()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expID, root.ID())
			}

		})
	}
}

func TestCollection_HasEdgeToFrom(t *testing.T) {
	type spec struct {
		name  string
		nodes []model.Node
		edges []model.Edge
		to    string
		from  string
		exp   bool
	}

	cases := []spec{
		{
			name:  "Success/EdgeExists",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			to:   "node1",
			from: "node2",
			exp:  true,
		},
		{
			name:  "Success/NoEdgeExists",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			to:   "node2",
			from: "node3",
			exp:  false,
		},
		{
			name:  "Success/EdgeExitsReverse",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
				&Edge{T: &mockNode{id: "node3"}, F: &mockNode{id: "node2"}},
			},
			to:   "node2",
			from: "node1",
			exp:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			actual := collection.HasEdgeFromTo(c.from, c.to)
			require.Equal(t, c.exp, actual)
		})
	}
}

func TestCollection_To(t *testing.T) {
	type spec struct {
		name   string
		nodes  []model.Node
		edges  []model.Edge
		input  string
		expIDs []string
	}

	cases := []spec{
		{
			name:  "Success/OneNodeFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			input:  "node2",
			expIDs: []string{"node1"},
		},
		{
			name:  "Success/MultipleNodesFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node2"}},
			},
			input:  "node1",
			expIDs: []string{"node2", "node3"},
		},
		{
			name:  "Success/NoNodesFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			input:  "node3",
			expIDs: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			nodes := collection.To(c.input)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			sort.Strings(actual)
			require.Equal(t, c.expIDs, actual)
		})
	}
}

func TestCollection_From(t *testing.T) {
	type spec struct {
		name   string
		nodes  []model.Node
		edges  []model.Edge
		input  string
		expIDs []string
	}

	cases := []spec{
		{
			name:  "Success/OneNodeFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node2"}},
			},
			input:  "node1",
			expIDs: []string{"node2"},
		},
		{
			name:  "Success/MultipleNodesFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node3"}},
			},
			input:  "node3",
			expIDs: []string{"node1", "node2"},
		},
		{
			name:  "Success/NoNodesFound",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			input:  "node2",
			expIDs: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			nodes := collection.From(c.input)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			sort.Strings(actual)
			require.Equal(t, c.expIDs, actual)
		})
	}
}

func TestCollection_Attributes(t *testing.T) {
	type spec struct {
		name          string
		nodes         []model.Node
		edges         []model.Edge
		expAttributes string
	}

	cases := []spec{
		{
			// TODO(jpower432)
			name: "Success/RootExists",
			nodes: []model.Node{
				&mockNode{
					id: "node1",
					attributes: &mockAttributes{
						"title": map[string]struct{}{
							"node1": {},
						},
					},
				},
				&mockNode{
					id: "node2",
					attributes: &mockAttributes{
						"title": map[string]struct{}{
							"node2": {},
						},
					},
				},
				&mockNode{
					id: "node3",
					attributes: &mockAttributes{
						"title": map[string]struct{}{
							"node3": {},
						},
					},
				},
			},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
			},
			expAttributes: "title=node3",
		},
		{
			name:  "Failure/NotRootExists",
			nodes: []model.Node{&mockNode{id: "node1"}, &mockNode{id: "node2"}, &mockNode{id: "node3"}},
			edges: []model.Edge{
				&Edge{T: &mockNode{id: "node2"}, F: &mockNode{id: "node1"}},
				&Edge{T: &mockNode{id: "node1"}, F: &mockNode{id: "node3"}},
				&Edge{T: &mockNode{id: "node3"}, F: &mockNode{id: "node2"}},
			},
			expAttributes: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			attr := collection.Attributes()
			if attr == nil {
				require.Len(t, c.expAttributes, 0)
			} else {
				require.Equal(t, c.expAttributes, attr.String())
			}
		})
	}
}

func makeTestCollection(t *testing.T, nodes []model.Node, edges []model.Edge) Collection {
	c := NewCollection("test")
	for _, node := range nodes {
		require.NoError(t, c.AddNode(node))
	}
	for _, edge := range edges {
		require.NoError(t, c.AddEdge(edge))
	}
	return *c
}

type mockNode struct {
	id         string
	attributes model.Attributes
}

var _ model.Node = &mockNode{}

func (m *mockNode) ID() string {
	return m.id
}

func (m *mockNode) Address() string {
	return "address"
}

func (m *mockNode) Attributes() model.Attributes {
	return m.attributes
}

type mockAttributes map[string]map[string]struct{}

var _ model.Attributes = &mockAttributes{}

func (m mockAttributes) Find(key string) []string {
	valSet, exists := m[key]
	if !exists {
		return nil
	}
	var vals []string
	for val := range valSet {
		vals = append(vals, val)
	}
	return vals
}

func (m mockAttributes) Exists(key, value string) bool {
	vals, exists := m[key]
	if !exists {
		return false
	}
	_, valExists := vals[value]
	return valExists
}

func (m mockAttributes) String() string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		for val := range m[key] {
			line := fmt.Sprintf("%s=%s,", key, val)
			out.WriteString(line)
		}
	}
	return strings.TrimSuffix(out.String(), ",")
}
