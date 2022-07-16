package ocimanifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/client/model/nodes/collection"
	"github.com/uor-framework/client/registryclient"
)

const (
	// AnnotationsSchemaName is the reference to the
	// default schema of the collection.
	AnnotationSchema = "uor.schema"
	// AnnotationSchemaLinks is the reference to linked
	// schemas for a collection. This will define all referenced
	// schemas for the collection and sub-collection. The tree will
	// be fully resolved.
	AnnotationSchemaLinks = "uor.schema.linked"
	// AnnotationCollectionLinks references the collections
	// that are linked to a collection node. The will only
	// reference adjacent collection and will not descend
	// into sub-collections.
	AnnotationCollectionLinks = "uor.collections.linked"
	// Separator is the value used to denote a list of
	// schema or collection in a manifest.
	Separator = ","
)

var (
	// ErrNoKnownSchema denotes that no schema
	// annotation is set on the manifest.
	ErrNoKnownSchema = errors.New("no schema")
	// ErrNoCollectionLinkes denotes that the manifest
	// does contain annotation that set collection links.
	ErrNoCollectionLinks = errors.New("no collection links")
)

// FetchSchema fetches schema information for a reference.
func FetchSchema(ctx context.Context, reference string, client registryclient.Remote) (string, []string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return "", nil, err
	}

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manBytes).Decode(&manifest); err != nil {
		return "", nil, err
	}

	schema, ok := manifest.Annotations[AnnotationSchema]
	if !ok {
		return "", nil, ErrNoKnownSchema
	}
	links := []string{manifest.Annotations[AnnotationSchemaLinks]}

	return schema, links, err
}

// ManifestToCollection converts a UOR managed OCI manifest to a Collection.
func ManifestToCollection(_ ocispec.Manifest) (collection.Collection, error) {
	fmt.Println("not implemented")
	return collection.Collection{}, nil
}

// CollectionToManifest converts a Collection to a UOR managed OCI manifest.
func CollectionToManifest(_ collection.Collection) (ocispec.Manifest, error) {
	fmt.Println("not implemented")
	return ocispec.Manifest{}, nil
}
