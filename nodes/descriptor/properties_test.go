package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
)

func TestProperties_MarshalJSON(t *testing.T) {
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	props := &Properties{
		Manifest: &uorspec.ManifestAttributes{
			RegistryHint: "test",
		},
		Descriptor: &uorspec.DescriptorAttributes{
			Component: uorspec.Component{
				AdditionalMetadata: []byte("2"),
			},
		},
		Others: set,
	}
	propsJSON, err := props.MarshalJSON()
	require.NoError(t, err)
	t.Log(string(propsJSON))
	t.Fail()
}
