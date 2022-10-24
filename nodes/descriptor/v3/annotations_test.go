package v3

import (
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/model"
)

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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expError != "" {

			} else {

			}
		})
	}
}
