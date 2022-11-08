package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2/registry"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/components"
	load "github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	"github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// Build builds collection from input and store it in the underlying content store.
// If successful, the root descriptor is returned.
func (d DefaultManager) Build(ctx context.Context, space workspace.Workspace, config clientapi.DataSetConfiguration, reference string, client registryclient.Client) ([]string, error) {
	var files []string
	err := space.Walk(func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traversing %s: %v", path, err)
		}
		if info == nil {
			return fmt.Errorf("no file info")
		}

		if info.Mode().IsRegular() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("path %q empty workspace", space.Path("."))
	}

	attributesByFile := map[string]model.AttributeSet{}
	var sets []model.AttributeSet
	for _, file := range config.Collection.Files {
		set, err := load.ConvertToModel(file.Attributes)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
		attributesByFile[file.File] = set
	}

	// Merge the sets to ensure the dataset configuration
	// meet the schema require.
	mergedSet, err := attributes.Merge(sets...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge attributes: %w", err)
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user. Also, collection the schema ID
	// to place in the descriptor properties.
	schemaID := schema.UnknownSchemaID
	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)

		_, _, err := client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return nil, fmt.Errorf("error configuring client: %v", err)
		}

		schemaDoc, detectedSchemaID, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return nil, err
		}

		if detectedSchemaID != "" {
			schemaID = detectedSchemaID
		}

		valid, err := schemaDoc.Validate(mergedSet)
		if err != nil {
			return nil, fmt.Errorf("schema validation error: %w", err)
		}
		if !valid {
			return nil, fmt.Errorf("attributes are not valid for schema %s", config.Collection.SchemaAddress)
		}
	}

	// To allow the files to be loaded relative to the render
	// workspace, change to the render directory. This is required
	// to get path correct in the description annotations.
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err := os.Chdir(space.Path()); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			d.logger.Errorf("%v", err)
		}
	}()

	descs, err := client.AddFiles(ctx, "", files...)
	if err != nil {
		return nil, err
	}

	// Gather workspace file metadata
	input := fmt.Sprintf("dir:%s", ".")
	inv, err := components.GenerateInventory(input, config)
	if err != nil {
		return nil, fmt.Errorf("inventory generation for %s: %w", space.Path(), err)
	}

	// Create nodes and update node properties
	var nodes []v2.Node
	for _, desc := range descs {
		location, ok := desc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			continue
		}
		// Using location as ID in this case because it is unique and
		// the digest may not be.
		node, err := v2.NewNode(location, desc)
		if err != nil {
			return nil, err
		}
		node.Location = location
		if err := components.InventoryToProperties(*inv, location, node.Properties); err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	}

	// Add user provided attributes to node properties
	descs, err = v2.UpdateDescriptors(nodes, schemaID, attributesByFile)
	if err != nil {
		return nil, err
	}

	// Store the DataSetConfiguration file in the manifest config of the OCI artifact for
	// later use.
	// Artifacts don't have configs. This will have to go with the regular descriptors.
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	configDesc, err := client.AddContent(ctx, uorspec.MediaTypeConfiguration, configJSON, nil)
	if err != nil {
		return nil, err
	}

	_, err = client.AddManifest(ctx, reference, configDesc, nil, descs...)
	if err != nil {
		return nil, err
	}

	desc, err := client.Save(ctx, reference, d.store)
	if err != nil {
		return nil, fmt.Errorf("client save error for reference %s: %v", reference, err)
	}
	d.logger.Infof("Artifact %s built with reference name %s\n", desc.Digest, reference)

	// Build index manifest
	if len(config.Collection.LinkedCollections) != 0 {
		d.logger.Infof("Processing %s links for aggregate", len(config.Collection.LinkedCollections))
		linkedDesc := []ocispec.Descriptor{desc}
		for _, l := range config.Collection.LinkedCollections {
			desc, _, err := client.GetManifest(ctx, l)
			if err != nil {
				return nil, fmt.Errorf("link %q: %w", l, err)
			}
			if desc.Annotations == nil {
				desc.Annotations = map[string]string{}
			}

			ref, err := registry.ParseReference(l)
			if err != nil {
				return nil, fmt.Errorf("link %q: %w", l, err)
			}
			linkAttr := descriptor.Properties{
				Manifest: &uorspec.ManifestAttributes{
					RegistryHint:  ref.Registry,
					NamespaceHint: ref.Repository,
					Transitive:    true,
				},
			}
			linkJSON, err := json.Marshal(linkAttr)
			if err != nil {
				return nil, err
			}
			desc.Annotations[uorspec.AnnotationUORAttributes] = string(linkJSON)
			linkedDesc = append(linkedDesc, desc)
		}
		aggregateRef := fmt.Sprintf("%s-aggregate", reference)
		desc, err := client.AddIndex(ctx, aggregateRef, nil, linkedDesc...)
		if err != nil {
			return nil, fmt.Errorf("aggregate reference %s: %w", aggregateRef, err)
		}
		d.logger.Infof("Aggregate %s built with reference name %s\n", desc.Digest, aggregateRef)
	}

	return []string{desc.Digest.String()}, nil
}

// fetchJSONSchema returns a schema type from a content store and a schema address.
func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, string, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, "", err
	}

	var schemaID string
	node, err := v2.NewNode(desc.Digest.String(), desc)
	props := node.Properties
	if props.Schema != nil {
		schemaID = props.Schema.ID
	}

	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return schema.Schema{}, "", fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return schema.Schema{}, "", err
	}
	loader, err := schema.FromBytes(schemaBytes)
	if err != nil {
		return schema.Schema{}, "", err
	}

	sc, err := schema.New(loader)
	return sc, schemaID, err
}
