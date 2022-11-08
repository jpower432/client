package v3

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

// ErrNoCollectionLinks denotes that the manifest
// does contain annotation that set collection links.
var ErrNoCollectionLinks = errors.New("no collection links")

// ResolveCollectionLinks finds linked collection references from a given input.
func ResolveCollectionLinks(reader io.Reader) ([]uorspec.Descriptor, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	var descriptor ocispec.Descriptor
	if err := json.NewDecoder(tee).Decode(&descriptor); err != nil {
		return nil, err
	}

	if descriptor.MediaType == uorspec.MediaTypeCollectionManifest {
		var manifest uorspec.Manifest
		if err := json.NewDecoder(&buf).Decode(&manifest); err != nil {
			return nil, err
		}

		return manifest.Links, nil
	}
	return nil, ErrNoCollectionLinks
}
