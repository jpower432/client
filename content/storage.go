package content

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"

	"github.com/uor-framework/client/model"
)

// Store defines the methods for adding, inspecting, and removing
// OCI content from a storage location. The interface wraps oras
// Storage and TagResolver interfaces for use with `oras` Copy methods.
type Store interface {
	// Storage represents a content-addressable storage where contents are
	// accessed via Descriptors.
	content.Storage
	// TagResolver defines methods for indexing tags.
	content.TagResolver
}

// Attribute Store defines the methods for retrieve information
// by attribute.
type AttributeStore interface {
	Store
	// ResolveByAttribute will return all descriptors associated
	// with a reference that match the attributes.
	ResolveByAttribute(context.Context, string, model.Matcher) ([]ocispec.Descriptor, error)
}

// GraphStore defines the methods for adding, inspecting, and removing
// OCI content from a storage location. The interface wraps oras
// Storage, TagResolver, and PredecessorFinder interfaces for use with `oras` extended copy methods.
type GraphStore interface {
	Store
	// Predecessors returns the nodes directly pointing to the current node.
	content.PredecessorFinder
	// Successors returns the nodes directly pointed by the current node.
	// In other words, returns the "children" of the current descriptor.
	Successors(ctx context.Context, node ocispec.Descriptor) ([]ocispec.Descriptor, error)
	// ResolveLinks returns all sub-collections references that are linked
	// to the root node.
	ResolveLinks(context.Context, string) ([]string, error)
}
