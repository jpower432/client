package descriptor

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	// AnnotationSchema is the reference to the
	// default schema of the collection.
	AnnotationSchema = "uor.schema"
	// AnnotationCollectionLinks references the collections
	// that are linked to a collection node. They will only
	// reference adjacent collection and will not descend
	// into sub-collections.
	AnnotationCollectionLinks = "uor.collections.linked"
	// AnnotationUORAttributes references the collection attributes in a
	// JSON format.
	AnnotationUORAttributes = "uor.attributes"
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

// TODO(jpower432): Resolve into a digest

// FetchSchema fetches schema information from a given input.
func FetchSchema(input io.Reader) (string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return "", err
	}

	schema, ok := manifest.Annotations[AnnotationSchema]
	if !ok {
		return "", ErrNoKnownSchema
	}

	return schema, nil
}

// TODO(jpower432): Resolve into a query

// ResolveCollectionLinks finds linked collection references from a given input.
func ResolveCollectionLinks(input io.Reader) ([]string, error) {
	var manifest ocispec.Manifest
	if err := json.NewDecoder(input).Decode(&manifest); err != nil {
		return nil, err
	}
	links, ok := manifest.Annotations[AnnotationCollectionLinks]
	if !ok || len(links) == 0 {
		return nil, ErrNoCollectionLinks
	}
	return strings.Split(links, Separator), nil
}
