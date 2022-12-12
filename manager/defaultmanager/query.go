package defaultmanager

import (
	"context"
	"fmt"

	specgo "github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2/registry"

	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
)

func (d DefaultManager) QueryLinks(ctx context.Context, host, digest string, matcher model.Matcher, client registryclient.Remote) ([]ocispec.Descriptor, error) {
	result, err := client.ResolveQuery(ctx, host, []string{digest}, nil, nil)
	if err != nil {
		return nil, err
	}

	var collections []collection.Collection
	for _, desc := range result.Manifests {
		if desc.Annotations == nil {
			continue
		}

		hint, ok := desc.Annotations["namespaceHint"]
		if !ok {
			continue
		}

		constructedRef := fmt.Sprintf("%s/%s@%s", host, hint, desc.Digest)
		collection, err := client.LoadCollection(ctx, constructedRef)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	filterResults, err := filterIndex(collections, matcher)
	if err != nil {
		return nil, err
	}
	return filterResults.Manifests, nil
}

func filterIndex(collections []collection.Collection, matcher model.Matcher) (ocispec.Index, error) {
	filteredIndex := ocispec.Index{
		Versioned: specgo.Versioned{
			SchemaVersion: 2,
		},
		MediaType: ocispec.MediaTypeImageIndex,
	}
	for _, currCol := range collections {
		subCollection, err := currCol.SubCollection(matcher)
		if err != nil {
			return filteredIndex, err
		}
		if len(subCollection.Nodes()) != 0 {
			rootNode, err := currCol.Root()
			if err != nil {
				return filteredIndex, err
			}

			rootDesc, ok := rootNode.(*v2.Node)
			if ok {
				rootOCIDesc := rootDesc.Descriptor()
				ref, err := registry.ParseReference(currCol.Address())
				if err != nil {
					return filteredIndex, err
				}
				props := descriptor.Properties{
					Link: &uorspec.LinkAttributes{
						RegistryHint:  ref.Registry,
						NamespaceHint: ref.Repository,
					},
				}
				propsJSON, err := props.MarshalJSON()
				if err != nil {
					return filteredIndex, err
				}
				rootOCIDesc.Annotations = map[string]string{
					uorspec.AnnotationUORAttributes: string(propsJSON),
				}
				filteredIndex.Manifests = append(filteredIndex.Manifests, rootOCIDesc)
			}
		}
	}
	return filteredIndex, nil
}
