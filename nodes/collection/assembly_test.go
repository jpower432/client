package collection

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestCollection_SubCollection(t *testing.T) {
	type spec struct {
		name       string
		nodes      []model.Node
		edges      []model.Edge
		matcher    model.Matcher
		expError   string
		assertFunc func(collection Collection) bool
	}

	cases := []spec{
		{
			name: "Success/NilMatcher",
			nodes: []model.Node{
				&testutils.FakeNode{I: "node1", A: attributes.Attributes{
					"test": attributes.NewString("test", "match"),
				}},
				&testutils.FakeNode{I: "node2", A: attributes.Attributes{
					"test": attributes.NewString("test", "notmatch"),
				}},
			},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
			},
			matcher: nil,
			assertFunc: func(collection Collection) bool {
				node1 := collection.NodeByID("node1")
				node2 := collection.NodeByID("node2")
				return node1 != nil && node2 != nil
			},
		},
		{
			name: "Success/OneNodeFiltered",
			nodes: []model.Node{
				&testutils.FakeNode{I: "node1", A: attributes.Attributes{
					"test": attributes.NewString("test", "match"),
				}},
				&testutils.FakeNode{I: "node2", A: attributes.Attributes{
					"test": attributes.NewString("test", "notmatch"),
				}},
			},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
			},
			matcher: matchers.PartialAttributeMatcher{
				"test": attributes.NewString("test", "match"),
			},
			assertFunc: func(collection Collection) bool {
				node1 := collection.NodeByID("node1")
				node2 := collection.NodeByID("node2")
				return node1 != nil && node2 == nil
			},
		},
		{
			name: "Success/AllNodesFiltered",
			nodes: []model.Node{
				&testutils.FakeNode{I: "node1", A: attributes.Attributes{
					"test": attributes.NewString("test", "match"),
				}},
				&testutils.FakeNode{I: "node2", A: attributes.Attributes{
					"test": attributes.NewString("test", "notmatch"),
				}},
			},
			edges: []model.Edge{
				&Edge{T: &testutils.FakeNode{I: "node2"}, F: &testutils.FakeNode{I: "node1"}},
			},
			matcher: matchers.PartialAttributeMatcher{
				"test": attributes.NewString("test", "nomatch"),
			},
			assertFunc: func(collection Collection) bool {
				return len(collection.Nodes()) == 0
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			collection := makeTestCollection(t, c.nodes, c.edges)
			actual, err := collection.SubCollection(c.matcher)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(actual))
			}
		})
	}
}
