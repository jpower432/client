package schema

import (
	"context"
	"encoding/json"
	"errors"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/client/registryclient"
)

const (
	AnnotationSchemaName = "uor.schema"
	AnnotationLinks      = "uor.schema.linked"
)

var (
	// ErrNoKnownSchema denotes that no schema
	// annotation is set on the manifest.
	ErrNoKnownSchema = errors.New("no schema")
)

// Fetch fetches schema information for a reference.
func Fetch(ctx context.Context, reference string, client registryclient.Client) (string, []string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return "", nil, err
	}

	var manifest ocispec.Manifest
	if err := json.NewDecoder(manBytes).Decode(&manifest); err != nil {
		return "", nil, err
	}

	schema, ok := manifest.Annotations[AnnotationSchemaName]
	if !ok {
		return "", nil, ErrNoKnownSchema
	}
	links := []string{manifest.Annotations[AnnotationLinks]}

	return schema, links, err
}
