package content

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/uor-framework/client/model"
	"oras.land/oras-go/v2/content"
)

// Store defines the methods for add, inspecting, and removing
// OCI content for a storage location. The interface wraps oras
// Storage and TagResolver interfaces for use with `oras` copy methods.
type Store interface {
	// Storage represents a content-addressable storage where contents are
	// accessed via Descriptors.
	content.Storage
	// TagResolver defines methods for indexing tags.
	content.TagResolver
	// ResolveByAttribute will return all descriptors associated
	// with a reference that match the attributes.
	ResolveByAttribute(context.Context, string, model.Matcher) ([]ocispec.Descriptor, error)
	// ResolveLinks
	//ResolveLinks(context.Context, string) ([]string, error)
}
