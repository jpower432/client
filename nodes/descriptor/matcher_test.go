package descriptor

import (
	"testing"

	"github.com/uor-framework/uor-client-go/attributes"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/util/testutils"
)

func TestJSONSubsetMatcher_Matches(t *testing.T) {
	mockAttributes := attributes.Attributes{
		"kind":    attributes.NewString("kind", "jpg"),
		"name":    attributes.NewString("name", "fish.jpg"),
		"another": attributes.NewString("another", "attribute"),
	}

	n := &testutils.FakeNode{A: mockAttributes}
	m := JSONSubsetMatcher(`{"name":"fish.jpg"}`)
	match, err := m.Matches(n)
	require.NoError(t, err)
	require.True(t, match)
}
