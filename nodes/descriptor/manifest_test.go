package descriptor

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

func TestResolveCollectionLinks(t *testing.T) {
	t.Run("Success/ArtifactManifest", func(t *testing.T) {
		manifest := ocispec.Artifact{
			MediaType: ocispec.MediaTypeArtifactManifest,
			Annotations: map[string]string{
				uorspec.AnnotationCollectionLinks: "alink",
			},
		}
		manifestJSON, err := json.Marshal(manifest)
		require.NoError(t, err)
		var buf bytes.Buffer
		buf.Write(manifestJSON)
		actual, err := ResolveCollectionLinks(&buf)
		require.NoError(t, err)
		require.Equal(t, []string{"alink"}, actual)
	})
	t.Run("Success/CollectionManifest", func(t *testing.T) {
		digest, err := digest.Parse("sha256:a078fbb3a7d1b312d0334ea261fb8d97ac2d95a0eb56f70b975d258dff486352")
		require.NoError(t, err)
		manifest := uorspec.Manifest{
			MediaType: uorspec.MediaTypeCollectionManifest,
			Links: []uorspec.Descriptor{
				{
					Digest: digest,
				},
			},
		}
		manifestJSON, err := json.Marshal(manifest)
		require.NoError(t, err)
		var buf bytes.Buffer
		buf.Write(manifestJSON)
		actual, err := ResolveCollectionLinks(&buf)
		require.NoError(t, err)
		require.Equal(t, []string{"sha256:a078fbb3a7d1b312d0334ea261fb8d97ac2d95a0eb56f70b975d258dff486352"}, actual)
	})

	t.Run("Success/NoLink", func(t *testing.T) {
		manifest := ocispec.Artifact{
			MediaType: ocispec.MediaTypeArtifactManifest,
		}
		manifestJSON, err := json.Marshal(manifest)
		require.NoError(t, err)
		var buf bytes.Buffer
		buf.Write(manifestJSON)
		_, err = ResolveCollectionLinks(&buf)
		require.ErrorIs(t, err, ErrNoCollectionLinks)
	})
}
