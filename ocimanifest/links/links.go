package links

import (
	"context"
	"errors"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model/nodes/collection"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
)

// LinkedCollection is a type of UOR collection that defines relationships
// between sub-collections.
type LinkedCollection struct {
	*collection.Collection
}

// New returns a new LinkedCollection.
func New(origin string) LinkedCollection {
	return LinkedCollection{
		Collection: collection.New(origin),
	}
}

// ToManifest will build an OCI Image manifest and store it in a local client from the LinkedCollection.
func (l LinkedCollection) ToManifest(ctx context.Context, schemaAddress string, client registryclient.Local) (ocispec.Descriptor, error) {
	return client.AddManifest(ctx, "", ocispec.Descriptor{}, nil)
}

// Build traverses the content store to build a LinkedCollection.
func Build(ctx context.Context, origin string, store content.GraphStore) (LinkedCollection, error) {
	visitedRefs := map[string]struct{}{origin: {}}
	linkedCollection := New(origin)

	linkedRefs, err := store.ResolveLinks(ctx, origin)
	if err := checkResolvedLinksError(origin, err); err != nil {
		return linkedCollection, err
	}

	for len(linkedRefs) != 0 {
		currRef := linkedRefs[0]
		linkedRefs = linkedRefs[1:]
		if _, ok := visitedRefs[currRef]; ok {
			continue
		}
		visitedRefs[currRef] = struct{}{}

		currLinks, err := store.ResolveLinks(ctx, currRef)
		if err := checkResolvedLinksError(currRef, err); err != nil {
			return linkedCollection, err
		}
		linkedRefs = append(linkedRefs, currLinks...)
	}
	return linkedCollection, nil
}

// checkResolvedLinksError logs errors when no collection is
// found.
func checkResolvedLinksError(ref string, err error) error {
	if err == nil {
		return nil
	}
	if !errors.Is(err, ocimanifest.ErrNoCollectionLinks) {
		return err
	}
	return nil
}


