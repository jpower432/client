package testutils

import "github.com/uor-framework/client/model"

var (
	_ model.Node       = &MockNode{}
	_ model.Attributes = &MockAttributes{}
	_ model.Node       = &MockIterableNode{}
)

type MockNode struct {
	I string
	A model.Attributes
}

func (m *MockNode) ID() string {
	return m.I
}

func (m *MockNode) Address() string {
	return "address"
}

func (m *MockNode) Attributes() model.Attributes {
	return m.A
}

type MockAttributes map[string]string

func (m MockAttributes) Find(key string) []string {
	val, exists := m[key]
	if !exists {
		return nil
	}
	return []string{val}
}

func (m MockAttributes) Exists(key, value string) bool {
	val, exists := m[key]
	if !exists {
		return false
	}
	return val == value
}

func (m MockAttributes) String() string {
	return ""
}

type MockIterableNode struct {
	I     string
	Index int
	A     model.Attributes
	Nodes []model.Node
}

func (m *MockIterableNode) ID() string {
	return m.I
}

func (m *MockIterableNode) Address() string {
	return "address"
}

func (m *MockIterableNode) Attributes() model.Attributes {
	return m.A
}

func (m *MockIterableNode) Len() int {
	if m.Index >= len(m.Nodes) {
		return 0
	}
	return len(m.Nodes[m.Index+1:])
}

func (m *MockIterableNode) Next() bool {
	if uint(m.Index)+1 < uint(len(m.Nodes)) {
		m.Index++
		return true
	}
	m.Index = len(m.Nodes)
	return false
}

func (m *MockIterableNode) Node() model.Node {
	if m.Index >= len(m.Nodes) || m.Index < 0 {
		return nil
	}
	return m.Nodes[m.Index]
}

func (m *MockIterableNode) Reset() {
	m.Index = -1
}

func (m *MockIterableNode) Error() error {
	return nil
}
