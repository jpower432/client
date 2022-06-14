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
	mockAttributes := mockAttributes{
		"kind":    "jpg",
		"name":    "fish.jpg",
		"another": "attribute",
	}

	n := &mockNode{attributes: mockAttributes}
	m := NewAttributeMatcher(mockAttributes)
	require.True(t, m.Matches(n))
}

type mockNode struct {
	attributes model.Attributes
}

var _ model.Node = &mockNode{}

func (m *mockNode) ID() string {
	return "node"
}

func (m *mockNode) Address() string {
	return "address"
}

func (m *mockNode) Attributes() model.Attributes {
	return m.attributes
}

type mockAttributes map[string]string

var _ model.Attributes = &mockAttributes{}

func (m mockAttributes) Find(key string) []string {
	val, exists := m[key]
	if !exists {
		return nil
	}
	return []string{val}
}

func (m mockAttributes) Exists(key, value string) bool {
	val, exists := m[key]
	if !exists {
		return false
	}
	return val == value
}

func (m mockAttributes) String() string {
	return ""
}
