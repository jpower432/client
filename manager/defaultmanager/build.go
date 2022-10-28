package defaultmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2/registry"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/attributes"
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
	mergedSet, err := attributes.Merge(sets)
	if err != nil {
		return "", fmt.Errorf("failed to merge attributes :%w", err)
	}

	// If a schema is present, pull it and do the validation before
	// processing the files to get quick feedback to the user.
	collectionManifestAnnotations := map[string]string{}
	if config.Collection.SchemaAddress != "" {
		d.logger.Infof("Validating dataset configuration against schema %s", config.Collection.SchemaAddress)
		collectionManifestAnnotations[descriptor.AnnotationSchema] = config.Collection.SchemaAddress

		_, _, err = client.Pull(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", fmt.Errorf("error configuring client: %v", err)
		}

		schemaDoc, err := fetchJSONSchema(ctx, config.Collection.SchemaAddress, d.store)
		if err != nil {
			return "", err
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

	var nodes []v2.Node
	for _, desc := range descs {
		node, err := v2.NewNode(desc.Digest.String(), desc)
		if err != nil {
			return "", err
		}
		nodes = append(nodes, *node)
	}

	// TODO(jpower432): Fill into component information here
	descs, err = v2.UpdateDescriptors(nodes, attributesByFile)
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

	// Write the root collection attributes
	if len(linkedDescs) > 0 {
		collectionManifestAnnotations[descriptor.AnnotationCollectionLinks] = formatLinks(config.Collection.LinkedCollections)
	}

	ref, err := registry.ParseReference(reference)
	if err != nil {
		return "", err
	}
	coreSchema := uorspec.ManifestAttributes{
		RegistryHint: ref.Registry,
	}
	coreSchemaJSON, err := json.Marshal(coreSchema)
	if err != nil {
		return "", err
	}
	collectionManifestAnnotations[uorspec.AnnotationUORAttributes] = string(coreSchemaJSON)
	_, err = client.AddManifest(ctx, reference, configDesc, collectionManifestAnnotations, descs...)
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
func fetchJSONSchema(ctx context.Context, schemaAddress string, store content.AttributeStore) (schema.Schema, error) {
	desc, err := store.AttributeSchema(ctx, schemaAddress)
	if err != nil {
		return schema.Schema{}, err
	}
	schemaReader, err := store.Fetch(ctx, desc)
	if err != nil {
		return schema.Schema{}, fmt.Errorf("error fetching schema from store: %w", err)
	}
	schemaBytes, err := ioutil.ReadAll(schemaReader)
	if err != nil {
		return schema.Schema{}, err
	}
	loader, err := schema.FromBytes(schemaBytes)
	if err != nil {
		return schema.Schema{}, err
	}
	return schema.New(loader)
}

// gatherLinkedCollections create null descriptors to denotes linked collections in a manifest with schema link information.
func gatherLinkedCollections(ctx context.Context, cfg clientapi.DataSetConfiguration, client registryclient.Client) ([]ocispec.Descriptor, error) {
	var linkedDescs []ocispec.Descriptor
	for _, collection := range cfg.Collection.LinkedCollections {

		rootSchema, err := getSchema(ctx, collection, client)
		if err != nil {
			return nil, fmt.Errorf("collection %q: %w", collection, err)
		}

		annotations := map[string]string{
			descriptor.AnnotationSchema: rootSchema,
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

// getSchema retrieves all schema information for a given reference.
func getSchema(ctx context.Context, reference string, client registryclient.Remote) (string, error) {
	_, manBytes, err := client.GetManifest(ctx, reference)
	if err != nil {
		return "", err
	}
	defer manBytes.Close()
	return descriptor.FetchSchema(manBytes)
}

func formatLinks(links []string) string {
	n := len(links)
	switch {
	case n == 1:
		return links[0]
	case n > 1:
		dedupLinks := deduplicate(links)
		return strings.Join(dedupLinks, descriptor.Separator)
	default:
		return ""
	}
}

func deduplicate(in []string) []string {
	links := map[string]struct{}{}
	var out []string
	for _, l := range in {
		if _, ok := links[l]; ok {
			continue
		}
		links[l] = struct{}{}
		out = append(out, l)
	}
	return out
}
