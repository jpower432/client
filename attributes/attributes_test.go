package attributes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	attributes := Attributes{
		"kind": json.RawMessage("jpg"),
		"name": json.RawMessage("fish.jpg"),
	}
	require.True(t, attributes.Exists("kind", "jpg"))
	require.False(t, attributes.Exists("kind", "png"))
}

func TestFind(t *testing.T) {
	attributes := Attributes{
		"kind": json.RawMessage("jpg"),
		"name": json.RawMessage("fish.jpg"),
	}
	result := attributes.Find("kind")
	require.Len(t, result, 1)
	require.Contains(t, result, "jpg")
}

func TestAttributes_String(t *testing.T) {
	expString := `"kind": "jpg"`
	attributes := Attributes{
		"kind": json.RawMessage("\"jpg\""),
		"name": json.RawMessage("\"fish.jpg\""),
	}
	require.Equal(t, expString, attributes.String())
}
