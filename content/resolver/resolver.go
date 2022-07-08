package resolver

import (
	"context"
	"fmt"
	"sync"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// Memory is a memory based resolver.
type Memory struct {
	descriptorLookup sync.Map // map[string]ocispec.Descriptor
}

func NewMemory() *Memory {
	return &Memory{
		descriptorLookup: sync.Map{},
	}
}

// Resolve resolves a reference to a descriptor.
func (m *Memory) Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error) {
	desc, ok := m.descriptorLookup.Load(reference)
	if !ok {
		return ocispec.Descriptor{}, fmt.Errorf("descriptor for reference %s is not stored", reference)
	}
	return desc.(ocispec.Descriptor), nil
}

func (m *Memory) Tag(_ context.Context, desc ocispec.Descriptor, reference string) error {
	m.descriptorLookup.Store(reference, desc)
	return nil
}

// Map dumps the memory into a built-in map structure.
// Like other operations, calling Map() is go-routine safe. However, it does not
// necessarily correspond to any consistent snapshot of the storage contents.
func (m *Memory) Map() map[string]ocispec.Descriptor {
	res := make(map[string]ocispec.Descriptor)
	m.descriptorLookup.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(ocispec.Descriptor)
		return true
	})
	return res
}
