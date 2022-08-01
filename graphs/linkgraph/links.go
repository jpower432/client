package links

import (
	"context"
	"errors"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/graphs/nodes/collection"
	"github.com/uor-framework/uor-client-go/ocimanifest"
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

// LoadFromStore traverses the content store to build a LinkedCollection.
func LoadFromStore(ctx context.Context, origin string, store content.GraphStore) (LinkedCollection, error) {
	visitedRefs := map[string]struct{}{origin: {}}
	linkedCollection := New(origin)
	linkedRefs := []string{origin}

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
