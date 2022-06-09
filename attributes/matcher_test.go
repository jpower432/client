package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/model"
)

func TestString(t *testing.T) {
	expString := `kind=jpg,name=fish.jpg`
	attributes := map[string]string{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	m := NewAttributeMatcher(attributes)
	require.Equal(t, expString, m.String())
}

func TestMatches(t *testing.T) {
	attributes := map[string]string{
		"kind": "jpg",
		"name": "fish.jpg",
	}

	n := &mockNode{attributes: attributes}
	n.attributes["another"] = "attribute"
	m := NewAttributeMatcher(attributes)
	require.True(t, m.Matches(n))
}

type mockNode struct {
	attributes map[string]string
}

var _ model.Node = &mockNode{}

func (m *mockNode) ID() string {
	return "node"
}

func (m *mockNode) Address() string {
	return "address"
}

func (m *mockNode) Attributes() map[string]string {
	return m.attributes
}
