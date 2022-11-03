package descriptor

import (
	"testing"

	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
)

func TestProperties_MarshalJSON(t *testing.T) {
	expJSON := `{"uor.core.manifest":{"registryHint":"test"},"uor.core.descriptor":{"id":"id","name":"","version":"","type":"","foundBy":"","locations":null,"licenses":null,"language":"","cpes":null,"purl":""},"uor.user.attributes":{"name":"test","size":2}}`
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
				ID: "id",
			},
		},
		Others: set,
	}
	propsJSON, err := props.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expJSON, string(propsJSON))
}
