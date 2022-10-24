package v3

import (
	"encoding/json"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

// TODO(jpower432): Move to descriptor

// FetchSchema fetches schema information from a given input.
func FetchSchema(input io.Reader) (string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return "", err
	}

	schema, ok := manifest.Annotations[descriptor.AnnotationSchema]
	if !ok {
		return "", descriptor.ErrNoKnownSchema
	}

	return schema, nil
}

// ResolveCollectionLinks finds linked collection references from a given input.
func ResolveCollectionLinks(input io.Reader) ([]uorspec.Descriptor, error) {
	var manifest uorspec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return nil, err
	}
	return manifest.Links, nil
}
