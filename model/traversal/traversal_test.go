package traversal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/util/testutils"
)

func TestTracker_Walk(t *testing.T) {
	type spec struct {
		name           string
		t              Tracker
		expError       error
		expInvocations int
	}

	cases := []spec{
		{
			name: "Success/VisitRootNode",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 3,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {&testutils.MockNode{I: "node2"}}},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 2,
		},
		{
			name: "Success/WithIterableNode",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 8,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{&testutils.MockNode{I: "node3"}}},
					},
				},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 3,
		},
		{
			name: "Success/DuplicateNodeID",
			t: Tracker{
				Budget: &Budget{
					NodeBudget: 8,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{&testutils.MockNode{I: "node1"}}},
					},
				},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 2,
		},
		{
			name: "Failure/ExceededBudget",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 0,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {&testutils.MockNode{I: "node2"}}},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 0,
			expError:       &ErrBudgetExceeded{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actualInvocations int
			visit := func(tr Tracker, n model.Node) error {
				t.Log("Visiting " + n.ID())
				actualInvocations++
				return nil
			}

			root, err := c.t.Tree.Root()
			require.NoError(t, err)
			err = c.t.Walk(root, visit)

			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expInvocations, actualInvocations)
			}
		})
	}
}

func TestTracker_WalkBFS(t *testing.T) {
	type spec struct {
		name           string
		t              Tracker
		m              model.Matcher
		expError       error
		expInvocations int
	}

	cases := []spec{
		{
			name: "Success/VisitRootNode",
			t: Tracker{
				Budget: &Budget{
					NodeBudget: 3,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {&testutils.MockNode{I: "node2"}}},
				},
				Seen: map[string]struct{}{},
			},
			m:              &mockMatcher{"node1"},
			expInvocations: 1,
		},
		{
			name: "Success/WithIterableNode",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 8,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{&testutils.MockNode{I: "node3"}}},
					},
				},
				},
				Seen: map[string]struct{}{},
			},
			m:              &mockMatcher{"node2"},
			expInvocations: 2,
		},
		{
			name: "Success/DuplicateNodeID",
			t: Tracker{
				Budget: &Budget{
					NodeBudget: 8,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockIterableNode{
							I:     "node2",
							Index: -1,
							Nodes: []model.Node{&testutils.MockNode{I: "node1"}}},
					},
				},
				},
				Seen: map[string]struct{}{},
			},
			m:              &mockMatcher{"node3"},
			expInvocations: 2,
		},
		{
			name: "Success/ShortestPath",
			t: Tracker{
				Budget: &Budget{
					NodeBudget: 20,
				},
				Tree: &mockTree{root: &testutils.MockNode{I: "node1"}, nodes: map[string][]model.Node{
					"node1": {
						&testutils.MockNode{I: "node4"},
					},
					"node2": {
						&testutils.MockNode{I: "node4"},
					},
					"node4": {},
				},
				},
				Seen: map[string]struct{}{},
			},
			m:              &mockMatcher{"node4"},
			expInvocations: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actualInvocations int
			visit := func(tr Tracker, n model.Node) error {
				t.Log("Visiting " + n.ID())
				actualInvocations++
				if c.m.Matches(n) {
					return ErrSkip
				}
				return nil
			}

			root, err := c.t.Tree.Root()
			require.NoError(t, err)
			err = c.t.WalkBFS(root, visit)

			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expInvocations, actualInvocations)
			}
		})
	}
}

// Mock tree structure
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

// Mock Matcher

type mockMatcher struct {
	criteria string
}

var _ model.Matcher = &mockMatcher{}

func (m *mockMatcher) String() string {
	return ""
}

func (m *mockMatcher) Matches(n model.Node) bool {
	return n.ID() == m.criteria
}
