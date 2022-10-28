package v3

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

func TestAttributesFromAttributeSet(t *testing.T) {
	expAttrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	attrs, err := AttributesFromAttributeSet(set)
	require.NoError(t, err)
	require.Equal(t, expAttrs, attrs)
}

func TestAttributesToAttributeSet(t *testing.T) {
	expJSON := `{"kind":"jpg","name":"fish.jpg","size":2}`
	attrs := map[string]json.RawMessage{
		"kind": []byte("\"jpg\""),
		"name": []byte("\"fish.jpg\""),
		"size": []byte("2"),
	}
	set, err := AttributesToAttributeSet(attrs, nil)
	require.NoError(t, err)
	require.Equal(t, expJSON, string(set.AsJSON()))
	// JSON standard lib will unmarshal all numbers as float64
	exists, err := set.Exists(attributes.NewFloat("size", 2))
	require.NoError(t, err)
	require.True(t, exists)
}

func TestAttributesFromAnnotations(t *testing.T) {
	annotations := map[string]string{
		descriptor.AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	expAttrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	attrs, err := AttributesFromAnnotations(annotations)
	require.NoError(t, err)
	require.Equal(t, expAttrs, attrs)
}

func TestAttributesToAnnotations(t *testing.T) {
	expMap := map[string]string{
		descriptor.AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	attrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	annotations, err := AttributesToAnnotations(attrs)
	require.NoError(t, err)
	require.Equal(t, expMap, annotations)
}
