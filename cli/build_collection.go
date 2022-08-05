package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
	load "github.com/uor-framework/uor-client-go/builder/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// BuildCollectionOptions describe configuration options that can
// be set using the build subcommand.
type BuildCollectionOptions struct {
	*BuildOptions
	RootDir     string
	BuildConfig string
}

// TODO(jpower432): Update to pull a validation schema attributes

var clientBuildCollectionExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build artifacts."},
		CommandString: "build collection my-directory localhost:5000/myartifacts:latest",
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build artifacts with custom annotations."},
		CommandString: "build collection my-directory localhost:5000/myartifacts:latest --build-config dataset-config.yaml",
	},
}

// NewBuildCollectionCmd creates a new cobra.Command for the build subcommand.
func NewBuildCollectionCmd(buildOpts *BuildOptions) *cobra.Command {
	o := BuildCollectionOptions{BuildOptions: buildOpts}

	cmd := &cobra.Command{
		Use:           "collection SRC DST",
		Short:         "Build and save an OCI artifact from files",
		Example:       examples.FormatExamples(clientBuildCollectionExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVarP(&o.BuildConfig, "build-config", "", o.BuildConfig, "config path for artifact building")

	return cmd
}

func (o *BuildCollectionOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RootDir = args[0]
	o.Destination = args[1]
	return nil
}

func (o *BuildCollectionOptions) Validate() error {
	if _, err := os.Stat(o.RootDir); err != nil {
		return fmt.Errorf("workspace directory %q: %v", o.RootDir, err)
	}

	return nil
}

func (o *BuildCollectionOptions) Run(ctx context.Context) error {
	space, err := workspace.NewLocalWorkspace(o.RootDir)
	if err != nil {
		return err
	}

	cache, err := layout.NewWithContext(ctx, o.cacheDir)
	if err != nil {
		return err
	}

	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}

	var files []string
	err = space.Walk(func(path string, info os.FileInfo, err error) error {
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
		return err
	}

	var config v1alpha1.DataSetConfiguration
	if len(o.BuildConfig) > 0 {
		config, err = load.ReadCollectionConfig(o.BuildConfig)
		if err != nil {
			return err
		}
	}

	// To allow the files to be loaded relative to the render
	// workspace, change to the render directory. This is required
	// to get path correct in the description annotations.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(space.Path()); err != nil {
		return err
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			o.Logger.Errorf("%v", err)
		}
	}()

	descs, err := client.AddFiles(ctx, "", files...)
	if err != nil {
		return err
	}

	descs, err = ocimanifest.UpdateLayerDescriptors(descs, config)
	if err != nil {
		return err
	}

	linkedDescs, linkedSchemas, err := gatherLinkedCollections(ctx, config, client)
	if err != nil {
		return err
	}

	descs = append(descs, linkedDescs...)

	// Add the attributes from the config to their respective blocks
	configDesc, err := client.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
	if err != nil {
		return err
	}

	// Write the root collection attributes
	manifestAnnotations := map[string]string{}
	if config.Collection.SchemaAddress != "" {
		manifestAnnotations[ocimanifest.AnnotationSchema] = config.Collection.SchemaAddress
	}

	if len(linkedDescs) > 0 {
		manifestAnnotations[ocimanifest.AnnotationSchemaLinks] = formatLinks(linkedSchemas)
		manifestAnnotations[ocimanifest.AnnotationCollectionLinks] = formatLinks(config.LinkedCollections)
	}

	_, err = client.AddManifest(ctx, o.Destination, configDesc, manifestAnnotations, descs...)
	if err != nil {
		return err
	}

	desc, err := client.Save(ctx, o.Destination, cache)
	if err != nil {
		return fmt.Errorf("client save error for reference %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Artifact %s built with reference name %s\n", desc.Digest, o.Destination)

	return client.Destroy()
}

// gatherLinkedCollections create null descriptors to denotes linked collections in a manifest with schema link information.
func gatherLinkedCollections(ctx context.Context, cfg v1alpha1.DataSetConfiguration, client registryclient.Client) ([]ocispec.Descriptor, []string, error) {
	var allLinkedSchemas []string
	var linkedDescs []ocispec.Descriptor
	for _, collection := range cfg.LinkedCollections {
		schema, linkedSchemas, err := ocimanifest.FetchSchema(ctx, collection, client)
		if err != nil {
			return nil, nil, err
		}

		if len(linkedSchemas) != 0 {
			allLinkedSchemas = append(allLinkedSchemas, linkedSchemas...)
		}

		allLinkedSchemas = append(allLinkedSchemas, schema)

		annotations := map[string]string{
			ocimanifest.AnnotationSchema:      schema,
			ocimanifest.AnnotationSchemaLinks: formatLinks(linkedSchemas),
		}
		// The bytes contain the collection name to keep the blobs unique within the manifest
		desc, err := client.AddContent(ctx, ocispec.MediaTypeImageLayer, []byte(collection), annotations)
		if err != nil {
			return nil, nil, err
		}
		linkedDescs = append(linkedDescs, desc)
	}
	return linkedDescs, allLinkedSchemas, nil
}

func formatLinks(links []string) string {
	n := len(links)
	switch {
	case n == 1:
		return links[0]
	case n > 1:
		dedupLinks := deduplicate(links)
		return strings.Join(dedupLinks, ocimanifest.Separator)
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
