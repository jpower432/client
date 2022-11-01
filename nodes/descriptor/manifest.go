package descriptor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
)

const (
	// Separator is the value used to denote a list of
	// schema or collection in a manifest.
	Separator = ","
)

var (
	// ErrNoKnownSchema denotes that no schema
	// annotation is set on the manifest.
	ErrNoKnownSchema = errors.New("no schema")
	// ErrNoCollectionLinks denotes that the manifest
	// does contain annotation that set collection links.
	ErrNoCollectionLinks = errors.New("no collection links")
)

// TODO(jpower432): Resolve into a digest for v3

// FetchSchema fetches schema information from a given input.
func FetchSchema(input io.Reader) (string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return "", err
	}

	schema, ok := manifest.Annotations[uorspec.AnnotationSchema]
	if !ok {
		return "", ErrNoKnownSchema
	}

	return schema, nil
}

// TODO(jpower432): Resolve into a query instead of slice of string

// ResolveCollectionLinks finds linked collection references from a given input.
func ResolveCollectionLinks(reader io.Reader) ([]string, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	var descriptor ocispec.Descriptor
	if err := json.NewDecoder(tee).Decode(&descriptor); err != nil {
		return nil, err
	}

	switch descriptor.MediaType {
	case uorspec.MediaTypeCollectionManifest:
		var manifest uorspec.Manifest
		if err := json.NewDecoder(&buf).Decode(&manifest); err != nil {
			return nil, err
		}

		var digests []string
		for _, l := range manifest.Links {
			digests = append(digests, l.Digest.String())
		}
		if len(digests) == 0 {
			return nil, ErrNoCollectionLinks
		}
		return digests, nil
	default:
		var manifest ocispec.Manifest
		if err := json.NewDecoder(&buf).Decode(&manifest); err != nil {
			return nil, err
		}
		links, ok := manifest.Annotations[uorspec.AnnotationCollectionLinks]
		if !ok || len(links) == 0 {
			return nil, ErrNoCollectionLinks
		}
		return strings.Split(links, Separator), nil
	}
}
