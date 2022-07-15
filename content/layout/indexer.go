package layout

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"

	"github.com/uor-framework/client/model/nodes/collection"
	"github.com/uor-framework/client/model/nodes/descriptor"
)

// IndexNode indexes predecessors for each direct successor of the given node.
// There is no data consistency issue as long as deletion is not implemented
// for the underlying storage.
func (l *Layout) IndexNode(ctx context.Context, fetcher content.Fetcher, node ocispec.Descriptor) error {
	successors, err := content.Successors(ctx, fetcher, node)
	if err != nil {
		return err
	}
	return l.indexNode(ctx, node, successors)
}

// indexNode indexes predecessors for each direct successor of the given node.
// There is no data consistency issue as long as deletion is not implemented
// for the underlying storage.
func (l *Layout) indexNode(ctx context.Context, node ocispec.Descriptor, successors []ocispec.Descriptor) error {
	n := descriptor.NewNode(node.Digest.String(), node)
	if err := l.graph.AddNode(n); err != nil {
		return err
	}
	for _, successor := range successors {
		s := descriptor.NewNode(successor.Digest.String(), successor)
		if err := l.graph.AddNode(s); err != nil {
			return err
		}
		e := collection.NewEdge(n, s)
		if err := l.graph.AddEdge(e); err != nil {
			return err
		}
	}
	return nil
}
