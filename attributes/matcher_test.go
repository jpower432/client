package attributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uor-framework/client/util/testutils"
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
	mockAttributes := testutils.MockAttributes{
		"kind":    "jpg",
		"name":    "fish.jpg",
		"another": "attribute",
	}

	n := &testutils.MockNode{A: mockAttributes}
	m := NewAttributeMatcher(mockAttributes)
	require.True(t, m.Matches(n))
}
