package traversal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/model"
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockIterableNode{id: "node2", idx: -1, nodes: []model.Node{&mockNode{id: "node3"}}}}},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 3,
		},
		{
			name: "Failure/ExceededBudget",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 0,
				},
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
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
			visit := func(t Tracker, n model.Node) error {
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

func TestTracker_WalkRecursively(t *testing.T) {
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockIterableNode{id: "node2", idx: -1, nodes: []model.Node{&mockNode{id: "node3"}}}}},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 3,
		},
		{
			name: "Failure/ExceededBudget",

			t: Tracker{
				Budget: &Budget{
					NodeBudget: 0,
				},
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
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
			visit := func(t Tracker, n model.Node) error {
				actualInvocations++
				return nil
			}

			root, err := c.t.Tree.Root()
			require.NoError(t, err)
			err = c.t.WalkRecursively(root, visit)

			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expInvocations, actualInvocations)
			}
		})
	}
}

func TestTracker_WalkWithStop(t *testing.T) {
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
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
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockIterableNode{id: "node2", idx: -1, nodes: []model.Node{&mockNode{id: "node3"}}}}},
				},
				Seen: map[string]struct{}{},
			},
			m:              &mockMatcher{"node2"},
			expInvocations: 2,
		},
		{
			name: "Success/NilMatcher",
			t: Tracker{
				Budget: &Budget{
					NodeBudget: 0,
				},
				Tree: &mockTree{root: &mockNode{id: "node1"}, nodes: map[string][]model.Node{
					"node1": {&mockNode{id: "node2"}}},
				},
				Seen: map[string]struct{}{},
			},
			expInvocations: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var actualInvocations int
			visit := func(t Tracker, n model.Node) error {
				actualInvocations++
				return nil
			}

			root, err := c.t.Tree.Root()
			require.NoError(t, err)
			err = c.t.WalkWithStop(root, c.m, visit)

			if c.expError != nil {
				require.ErrorAs(t, err, &c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expInvocations, actualInvocations)
			}
		})
	}
}

// Mock node types
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

type mockIterableNode struct {
	id         string
	idx        int
	attributes model.Attributes
	nodes      []model.Node
}

var _ model.Node = &mockIterableNode{}

func (m *mockIterableNode) ID() string {
	return m.id
}

func (m *mockIterableNode) Address() string {
	return "address"
}

func (m *mockIterableNode) Attributes() model.Attributes {
	return m.attributes
}

func (m *mockIterableNode) Len() int {
	if m.idx >= len(m.nodes) {
		return 0
	}
	return len(m.nodes[m.idx+1:])
}

func (m *mockIterableNode) Next() bool {
	if uint(m.idx)+1 < uint(len(m.nodes)) {
		m.idx++
		return true
	}
	m.idx = len(m.nodes)
	return false
}

func (m *mockIterableNode) Node() model.Node {
	if m.idx >= len(m.nodes) || m.idx < 0 {
		return nil
	}
	return m.nodes[m.idx]
}

func (m *mockIterableNode) Reset() {
	m.idx = -1
}

func (m *mockIterableNode) Error() error {
	return nil
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

func (m *mockTree) From(n model.Node) []model.Node {
	return m.nodes[n.ID()]
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
