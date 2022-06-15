package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/util/testutils"
)

func TestFindAllPartialMatches(t *testing.T) {
	type spec struct {
		name     string
		m        PartialAttributeMatcher
		expError error
		t        model.Tree
		expIDs   []string
	}

	cases := []spec{
		{
			name: "Success/VisitRootNode",
			t: &mockTree{root: &testutils.MockNode{
				I: "node1",
				A: testutils.MockAttributes{"title": "node1"},
			},
				nodes: map[string][]model.Node{
					"node1": {&testutils.MockNode{I: "node2"}}},
			},

			m:      PartialAttributeMatcher{"title": "node1"},
			expIDs: []string{"node1"},
		},
		{
			name: "Success/WithIterableNode",

			t: &mockTree{
				root: &testutils.MockNode{
					I: "node1",
					A: testutils.MockAttributes{
						"kind": "txt",
					},
				},
				nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{
								&testutils.MockNode{
									I: "node3",
									A: testutils.MockAttributes{
										"kind": "txt",
									},
								},
							}},
					},
				},
			},
			m: PartialAttributeMatcher{
				"kind": "txt",
			},
			expIDs: []string{"node1", "node3"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			nodes, err := FindAllPartialMatches(c.m, c.t)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expIDs, actual)
			}
		})
	}
}

func TestFindAllExactMatches(t *testing.T) {
	type spec struct {
		name     string
		m        ExactAttributeMatcher
		expError error
		t        model.Tree
		expIDs   []string
	}

	cases := []spec{
		{
			name: "Success/VisitRootNode",
			t: &mockTree{root: &testutils.MockNode{
				I: "node1",
				A: testutils.MockAttributes{"title": "node1"},
			},
				nodes: map[string][]model.Node{
					"node1": {&testutils.MockNode{I: "node2"}}},
			},

			m:      ExactAttributeMatcher{"title": "node1"},
			expIDs: []string{"node1"},
		},
		{
			name: "Success/WithIterableNode",

			t: &mockTree{
				root: &testutils.MockNode{
					I: "node1",
					A: testutils.MockAttributes{
						"kind": "txt",
					},
				},
				nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{
								&testutils.MockNode{
									I: "node3",
									A: testutils.MockAttributes{
										"kind":    "txt",
										"another": "attribute",
									},
								},
							}},
					},
				},
			},
			m: ExactAttributeMatcher{
				"kind": "txt",
			},
			expIDs: []string{"node1"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			nodes, err := FindAllExactMatches(c.m, c.t)
			var actual []string
			for _, node := range nodes {
				actual = append(actual, node.ID())
			}
			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expIDs, actual)
			}
		})
	}
}

type mockTree struct {
	root  model.Node
	nodes map[string][]model.Node
}

var _ model.Tree = &mockTree{}

func (m *mockTree) Root() (model.Node, error) {
	return m.root, nil
}

func (m *mockTree) From(id string) []model.Node {
	return m.nodes[id]
}
