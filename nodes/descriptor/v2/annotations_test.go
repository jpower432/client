package v2

import (
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

func TestUpdateLayerDescriptors(t *testing.T) {

	type spec struct {
		name           string
		fileAttributes map[string]model.AttributeSet
		expError       string
		expDesc        []ocispec.Descriptor
	}

	cases := []spec{
		{
			name: "Success/NoAttributes",
			expDesc: []ocispec.Descriptor{
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20a6",
					Size:      2,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.jpg",
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.jpg\"}}",
					},
				},
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
					Size:      8,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.json",
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.json\"}}",
					},
				},
			},
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
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.jpg\"},\"myschema\":{\"image\":true}}",
					},
				},
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
					Size:      8,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.json",
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.json\"},\"myschema\":{\"metadata\":true}}",
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
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.jpg\"},\"myschema\":{\"image\":true,\"publisher\":\"test\"}}",
					},
				},
				{
					MediaType: ocispec.MediaTypeImageLayerGzip,
					Digest:    "sha256:84f48921e4ed2e0b370fa314a78dadd499cde260032bcfcd6c1d5089d6cc20456",
					Size:      8,
					Annotations: map[string]string{
						ocispec.AnnotationTitle:         "fish.json",
						uorspec.AnnotationUORAttributes: "{\"converted\":{\"org.opencontainers.image.title\":\"fish.json\"},\"myschema\":{\"publisher\":\"test\"}}",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
				node.Location = desc.Annotations[ocispec.AnnotationTitle]
				nodes = append(nodes, *node)
			}
			res, err := UpdateDescriptors(nodes, "myschema", c.fileAttributes)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expDesc, res)
			}
		})
	}
}
