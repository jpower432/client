package registryclient

import (
	"context"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/nodes/collection"
)

// Client defines methods to interact with OCI artifacts
// in various contexts.
type Client interface {
	Remote
	Local
}

// Remote defines methods to interact with OCI
// artifacts in remote contexts.
type Remote interface {
	// Push pushes an artifact to a remote registry from a source
	// content store and returns the root manifest digest.
	Push(context.Context, content.Store, string) (ocispec.Descriptor, error)
	// Pull pulls an artifact from a remote registry to a local
	// content store. If successful it returns the root descriptor and all the descriptors pulled.
	Pull(context.Context, string, content.Store) (ocispec.Descriptor, []ocispec.Descriptor, error)
	// GetManifest retrieves the root manifest for a reference.
	GetManifest(context.Context, string) (ocispec.Descriptor, io.ReadCloser, error)
	// GetContent retrieves the content for a specified descriptor at a specified reference.
	GetContent(context.Context, string, ocispec.Descriptor) ([]byte, error)
	// LoadCollection loads a collection from a remote reference.
	LoadCollection(context.Context, string) (collection.Collection, error)
	AttributeFinder
}

// AttributeFinder resolve attribute queries for v3 compatible registries.
type AttributeFinder interface {
	ResolveAttributeQuery(context.Context, string, ocispec.Descriptor) (ocispec.Index, error)
}

// Local defines methods to interact with OCI artifacts
// in a local context. An underlying store can be used to store
// each descriptor and is returned the Store method for use with
// Push and Pull operations.
type Local interface {
	DescriptorAdder
	// Save saves a built artifact to local store.
	Save(context.Context, string, content.Store) (ocispec.Descriptor, error)
	// Store returns the underlying content store
	// used for OCI artifact building.
	Store() (content.Store, error)
	// Destroy cleans up temporary files on-disk
	// for tracking descriptors
	Destroy() error
}

// DescriptorAdder defines methods to add OCI descriptors to an
// underlying storage type.
type DescriptorAdder interface {
	// AddFiles loads one or more files to create OCI descriptors with a specific
	// media type and pushes them into underlying storage.
	AddFiles(context.Context, string, ...string) ([]ocispec.Descriptor, error)
	// AddContent creates and stores a descriptor from content in bytes, a media type, and
	// annotations.
	AddContent(context.Context, string, []byte, map[string]string) (ocispec.Descriptor, error)
	// AddManifest creates and stores a manifest for an image reference.
	// This is generated from the config descriptor and artifact descriptors.
	AddManifest(context.Context, string, ocispec.Descriptor, map[string]string, ...ocispec.Descriptor) (ocispec.Descriptor, error)
	// AddIndex creates and stores an index manifest for an image reference.
	// This is generated from the config descriptor and artifact descriptors.
	AddIndex(context.Context, string, map[string]string, ...ocispec.Descriptor) (ocispec.Descriptor, error)
}
