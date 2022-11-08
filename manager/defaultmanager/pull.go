package defaultmanager

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/registryclient"
)

// Pull pulls a single collection to a specified storage destination.
// If successful, the file locations are returned.
func (d DefaultManager) Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error) {
	descs, err := d.pullCollection(ctx, source, destination, remote)
	if err != nil {
		return nil, err
	}

	var digests []string
	for _, desc := range descs {
		digests = append(digests, desc.Digest.String())
		d.logger.Infof("Found matching digest %s", desc.Digest)
	}
	return digests, nil
}

// pullCollection pulls a single collection and returns the manifest descriptors and an error.
func (d DefaultManager) pullCollection(ctx context.Context, reference string, destination content.Store, remote registryclient.Remote) ([]ocispec.Descriptor, error) {
	rootDesc, descs, err := remote.Pull(ctx, reference, destination)
	if err != nil {
		return nil, err
	}

	// Ensure the store is tagged with the new reference.
	if len(rootDesc.Digest) != 0 {
		return descs, d.store.Tag(ctx, rootDesc, reference)
	}

	return descs, nil
}
