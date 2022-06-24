package cache

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/uor-framework/client/model"
)

// Storage defines the methods for add, inspecting, and removing
// OCI content for a storage location.
type Storage interface {
	// Add adds one or more descriptors to the content store if it
	// does not exist.
	Add(context.Context, ...ocispec.Descriptor) error
	// Delete deletes one or more descriptors from the
	// content store.
	Delete(context.Context, ...ocispec.Descriptor) error
	// List returns all OCI Descriptors
	// that exists in a content store.
	List() []ocispec.Descriptor
	// Exists check whether a OCI descriptor exists
	// in the content store.
	Exists(ocispec.Descriptor) bool
	// Index returns the OCI index of the content store.
	// The index will store all aggregate manifest information
	// (e.g. schema and attributes/annotations).
	Index() (ocispec.Index, error)
	// LookupByAttribute returns OCI descriptors based
	// on the specified attributes.
	LookupByAttribute(model.Attributes) ([]ocispec.Descriptor, error)
}
