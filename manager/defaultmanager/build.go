package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/components"
	load "github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// Build builds collection from input and store it in the underlying content store.
// If successful, the root descriptor is returned.
func (d DefaultManager) Build(ctx context.Context, space workspace.Workspace, config clientapi.DataSetConfiguration, reference string, client registryclient.Client) (string, error) {
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
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("path %q empty workspace", space.Path("."))
	}

	attributesByFile := map[string]model.AttributeSet{}
	var sets []model.AttributeSet
	for _, file := range config.Collection.Files {
		set, err := load.ConvertToModel(file.Attributes)
		if err != nil {
			return "", err
		}
		sets = append(sets, set)
		attributesByFile[file.File] = set
	}

	// Merge the sets to ensure the dataset configuration
	// meet the schema require.
	mergedSet, err := attributes.Merge(sets...)
	if err != nil {
		return "", fmt.Errorf("failed to merge attributes: %w", err)
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user. Also, collection the schema ID
	// to place in the descriptor properties.
	schemaID := schema.UnknownSchemaID
	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)

		_, _, err := client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", fmt.Errorf("error configuring client: %v", err)
		}

		schemaDoc, detectedschemaID, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
		}

		if detectedschemaID != "" {
			schemaID = detectedschemaID
		}

		valid, err := schemaDoc.Validate(mergedSet)
		if err != nil {
			return "", fmt.Errorf("schema validation error: %w", err)
		}
		if !valid {
			return "", fmt.Errorf("attributes are not valid for schema %s", config.Collection.SchemaAddress)
		}
	}

	// To allow the files to be loaded relative to the render
	// workspace, change to the render directory. This is required
	// to get path correct in the description annotations.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if err := os.Chdir(space.Path()); err != nil {
		return "", err
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			d.logger.Errorf("%v", err)
		}
	}()

	descs, err := client.AddFiles(ctx, "", files...)
	if err != nil {
		return "", err
	}

	// Gather workspace file metadata
	input := fmt.Sprintf("dir:%s", ".")
	inv, err := components.GenerateInventory(input, config)
	if err != nil {
		return "", fmt.Errorf("inventory generation for %s: %w", space.Path(), err)
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
			return "", err
		}
		node.Location = location
		if err := components.InventoryToProperties(*inv, location, node.Properties); err != nil {
			return "", err
		}
		nodes = append(nodes, *node)
	}

	// Add user provided attributes to node properties
	descs, err = v2.UpdateDescriptors(nodes, schemaID, attributesByFile)
	if err != nil {
		return "", err
	}

	linkedDescs, err := gatherLinkedCollections(ctx, config, client)
	if err != nil {
		return "", err
	}

	descs = append(descs, linkedDescs...)

	// Store the DataSetConfiguration file in the manifest config of the OCI artifact for
	// later use.
	// Artifacts don't have configs. This will have to go with the regular descriptors.
	configJSON, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	configDesc, err := client.AddContent(ctx, uorspec.MediaTypeConfiguration, configJSON, nil)
	if err != nil {
		return "", err
	}

	_, err = client.AddManifest(ctx, reference, configDesc, nil, descs...)
	if err != nil {
		return "", err
	}

	desc, err := client.Save(ctx, reference, d.store)
	if err != nil {
		return "", fmt.Errorf("client save error for reference %s: %v", reference, err)
	}
	d.logger.Infof("Artifact %s built with reference name %s\n", desc.Digest, reference)

	return desc.Digest.String(), nil
}

// fetchJSONSchema returns a schema type from a content store and a schema address.
func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, string, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, "", err
	}

	node, err := v2.NewNode(desc.Digest.String(), desc)
	schemaID := node.Properties.Schema.ID

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

// gatherLinkedCollections create null descriptors to denotes linked collections in a manifest with schema link information.
func gatherLinkedCollections(ctx context.Context, cfg clientapi.DataSetConfiguration, client registryclient.Client) ([]ocispec.Descriptor, error) {
	var linkedDescs []ocispec.Descriptor
	for _, collection := range cfg.Collection.LinkedCollections {

		annotations := map[string]string{
			"link": "true",
		}
		// The bytes contain the collection name to keep the blobs unique within the manifest
		desc, err := client.AddContent(ctx, ocispec.MediaTypeImageLayer, []byte(collection), annotations)
		if err != nil {
			return nil, err
		}
		linkedDescs = append(linkedDescs, desc)
	}
	return linkedDescs, nil
}
