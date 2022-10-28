package v2

import (
	"encoding/json"
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

func TestAnnotationsFromAttributeSet(t *testing.T) {
	expMap := map[string]string{
		uorspec.AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	set := attributes.Attributes{
		"name": attributes.NewString("name", "test"),
		"size": attributes.NewInt("size", 2),
	}
	annotations, err := AnnotationsFromAttributeSet(set)
	require.NoError(t, err)
	require.Equal(t, expMap, annotations)
}

func TestAnnotationsToAttributeSet(t *testing.T) {
	expJSON := `{"kind":"jpg","name":"fish.jpg","ref":"example","size":2}`
	annotations := map[string]string{
		"ref":                           "example",
		uorspec.AnnotationUORAttributes: `{"kind":"jpg","name":"fish.jpg","size":2}`,
	}
	set, err := AnnotationsToAttributeSet(annotations, nil)
	require.NoError(t, err)
	setJSON, err := set.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, expJSON, string(setJSON))
	// JSON standard lib will unmarshal all numbers as float64
	exists, err := set.Exists(attributes.NewFloat("size", 2))
	require.NoError(t, err)
	require.True(t, exists)
}

func TestAnnotationsToAttributes(t *testing.T) {
	annotations := map[string]string{
		uorspec.AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	expAttrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	attrs, err := AnnotationsToAttributes(annotations)
	require.NoError(t, err)
	require.Equal(t, expAttrs, attrs)
}

func TestAnnotationsFromAttributes(t *testing.T) {
	expMap := map[string]string{
		uorspec.AnnotationUORAttributes: "{\"name\":\"test\",\"size\":2}",
	}
	attrs := map[string]json.RawMessage{
		"name": []byte("\"test\""),
		"size": []byte("2"),
	}
	annotations, err := AnnotationsFromAttributes(attrs)
	require.NoError(t, err)
	require.Equal(t, expMap, annotations)
}

func TestUpdateLayerDescriptors(t *testing.T) {

	descs := []ocispec.Descriptor{
		{
			MediaType: ocispec.MediaTypeImageLayerGzip,
			Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20a6",
			Size:      2,
			Annotations: map[string]string{
				ocispec.AnnotationTitle: "fish.jpg",
			},
		},
		{
			MediaType: ocispec.MediaTypeImageLayerGzip,
			Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
			Size:      8,
			Annotations: map[string]string{
				ocispec.AnnotationTitle: "fish.json",
			},
		},
	}
	var nodes []Node
	for _, desc := range descs {
		node, err := NewNode(desc.Digest.String(), desc)
		require.NoError(t, err)
		nodes = append(nodes, *node)
	}
	type spec struct {
		name           string
		fileAttributes map[string]model.AttributeSet
		expError       string
		expDesc        []ocispec.Descriptor
	}

	cases := []spec{
		{
			name:    "Success/NoAttributes",
			expDesc: descs,
		},
		{
			name: "Success/SeperatedAttributes",
			fileAttributes: map[string]model.AttributeSet{
				"*.jpg": attributes.Attributes{
					"image": attributes.NewBool("image", true),
				},
				"*.json": attributes.Attributes{
					"metadata": attributes.NewBool("metadata", true),
				},
			},
			expDesc: []ocispec.Descriptor{
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20a6",
					Size:      2,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.jpg",
						uorspec.AnnotationUORAttributes: "{\"image\":true}",
					},
				},
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
					Size:      8,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.json",
						uorspec.AnnotationUORAttributes: "{\"metadata\":true}",
					},
				},
			},
		},
		{
			name: "Success/OverlappingAttributes",
			fileAttributes: map[string]model.AttributeSet{
				"*.jpg": attributes.Attributes{
					"image": attributes.NewBool("image", true),
				},
				"*": attributes.Attributes{
					"publisher": attributes.NewString("publisher", "test"),
				},
			},
			expDesc: []ocispec.Descriptor{
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20a6",
					Size:      2,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.jpg",
						uorspec.AnnotationUORAttributes: "{\"image\":true,\"publisher\":\"test\"}",
					},
				},
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
					Size:      8,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.json",
						uorspec.AnnotationUORAttributes: "{\"publisher\":\"test\"}",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res, err := UpdateDescriptors(nodes, c.fileAttributes)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expDesc, res)
			}
		})
	}
}
