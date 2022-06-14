package basic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	attributes := Attributes{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	require.True(t, attributes.Exists("kind", "jpg"))
}

func TestFind(t *testing.T) {
	expVals := []string{"fish.jpg"}
	attributes := Attributes{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	require.Equal(t, expVals, attributes.Find("name"))
}

func TestString(t *testing.T) {
	expString := `kind=jpg,name=fish.jpg`
	attributes := Attributes{
		"kind": "jpg",
		"name": "fish.jpg",
	}
	require.Equal(t, expString, attributes.String())
}
