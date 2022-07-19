package content

import (
	"context"
	"fmt"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Resolver is a memory based resolver.
type Resolver struct {
	descriptorLookup sync.Map // map[string]ocispec.Descriptor
}

func NewResolver() *Resolver {
	return &Resolver{
		descriptorLookup: sync.Map{},
	}
}

// Resolve resolves a reference to a descriptor.
func (r *Resolver) Resolve(_ context.Context, reference string) (ocispec.Descriptor, error) {
	desc, ok := r.descriptorLookup.Load(reference)
	if !ok {
		return ocispec.Descriptor{}, fmt.Errorf("descriptor for reference %s is not stored", reference)
	}
	return desc.(ocispec.Descriptor), nil
}

func (r *Resolver) Tag(_ context.Context, desc ocispec.Descriptor, reference string) error {
	r.descriptorLookup.Store(reference, desc)
	return nil
}

// Map dumps the memory into a built-in map structure.
// Like other operations, calling Map() is go-routine safe. However, it does not
// necessarily correspond to any consistent snapshot of the storage contents.
func (r *Resolver) Map() map[string]ocispec.Descriptor {
	res := make(map[string]ocispec.Descriptor)
	r.descriptorLookup.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(ocispec.Descriptor)
		return true
	})
	return res
}
