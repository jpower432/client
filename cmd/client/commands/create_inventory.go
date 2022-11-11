package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/formats/spdx22json"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/components"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	collectionloader "github.com/uor-framework/uor-client-go/nodes/collection/loader"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/examples"
	"github.com/uor-framework/uor-client-go/version"
)

// InventoryOptions describe configuration options that can
// be set using the push subcommand.
type InventoryOptions struct {
	*CreateOptions
	Source string
	Format string
}

var clientInventoryExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Build inventory from artifacts."},
	CommandString: "inventory localhost:5000/myartifacts:latest",
}

// NewInventoryCmd creates a new cobra.Command for the inventory subcommand.
func NewInventoryCmd(createOpts *CreateOptions) *cobra.Command {
	o := InventoryOptions{CreateOptions: createOpts}

	cmd := &cobra.Command{
		Use:           "inventory SRC",
		Short:         "Create software inventories from UOR artifacts",
		Example:       examples.FormatExamples(clientInventoryExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	return cmd
}

func (o *InventoryOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.Source = args[0]
	return nil
}

func (o *InventoryOptions) Validate() error {
	return nil
}

func (o *InventoryOptions) Run(ctx context.Context) error {
	client, err := orasclient.NewClient(
		orasclient.SkipTLSVerify(o.Insecure),
		orasclient.WithAuthConfigs(o.Configs),
		orasclient.WithPlainHTTP(o.PlainHTTP),
	)
	if err != nil {
		return fmt.Errorf("error configuring client: %v", err)
	}
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	desc, _, err := client.GetManifest(ctx, o.Source)
	if err != nil {
		return err
	}
	co := collection.New(o.Source)
	if err := loadCollection(ctx, co, o.Source, desc, client); err != nil {
		return err
	}

	inventory, err := collectionToInventory(ctx, co, client)
	if err != nil {
		return err
	}

	formatter := spdx22json.Format()
	return formatter.Encode(o.IOStreams.Out, inventory)
}

// collectionToInventory traverses and fully resolves a collection and create a software inventory from the graph.
// This only fills out SPDX required data at this point.
func collectionToInventory(ctx context.Context, graph *collection.Collection, client registryclient.Remote) (sbom.SBOM, error) {
	inventory := sbom.SBOM{
		Artifacts: sbom.Artifacts{
			FileDigests:  map[source.Coordinates][]file.Digest{},
			FileMetadata: map[source.Coordinates]source.FileMetadata{},
		},
	}

	var packages []pkg.Package
	root, err := graph.Root()
	if err != nil {
		return inventory, err
	}

	seen := map[string]struct{}{}
	// Process and pull links before pulling the requested manifests
	tracker := traversal.NewTracker(root, nil)
	handler := traversal.HandlerFunc(func(ctx context.Context, tracker traversal.Tracker, node model.Node) ([]model.Node, error) {
		if _, ok := seen[node.ID()]; ok {
			return nil, traversal.ErrSkip
		}

		desc, ok := node.(*v2.Node)
		if !ok {
			return nil, nil
		}

		successors := graph.From(node.ID())
		props := desc.Properties
		if props == nil {
			return successors, nil
		}

		attribute := props.FindBySchema(schema.ConvertedSchemaID, ocispec.AnnotationTitle)
		var title string
		if attribute != nil {
			title, err = attribute.AsString()
			if err != nil {
				return nil, nil
			}
		}

		coordinates := source.Coordinates{
			RealPath:     title,
			FileSystemID: node.ID(),
		}

		digest := file.Digest{
			Algorithm: desc.Descriptor().Digest.Algorithm().String(),
			Value:     desc.Descriptor().Digest.Encoded(),
		}
		inventory.Artifacts.FileDigests[coordinates] = []file.Digest{digest}

		fileMeta := source.FileMetadata{
			Size:     desc.Descriptor().Size,
			MIMEType: desc.Descriptor().MediaType,
		}
		inventory.Artifacts.FileMetadata[coordinates] = fileMeta

		if props.IsAComponent() {
			locations := source.LocationSet{}
			for _, loc := range props.Descriptor.Locations {
				locations.Add(source.NewLocation(loc))
			}

			var cpes []pkg.CPE
			for _, cpe := range props.Descriptor.CPEs {
				c, err := pkg.NewCPE(cpe)
				if err != nil {
					return nil, err
				}
				cpes = append(cpes, c)
			}

			var metaType pkg.MetadataType
			var metadata interface{}
			if len(props.Descriptor.AdditionalMetadata) == 1 {
				for typ, meta := range props.Descriptor.AdditionalMetadata {
					metaType = pkg.MetadataType(typ)
					if _, ok := pkg.MetadataTypeByName[metaType]; !ok {
						continue
					}
					metadata = meta
				}
			}

			p := pkg.Package{
				Name:         props.Descriptor.Name,
				Version:      props.Descriptor.Version,
				FoundBy:      props.Descriptor.FoundBy,
				Locations:    locations,
				Licenses:     props.Descriptor.Licenses,
				Language:     pkg.Language(props.Descriptor.Language),
				Type:         pkg.Type(props.Descriptor.Type),
				CPEs:         cpes,
				PURL:         props.Descriptor.PURL,
				MetadataType: metaType,
				Metadata:     metadata,
			}
			p.OverrideID(artifact.ID(props.Descriptor.ID))
			packages = append(packages, p)
		}

		// Load link and provide access to those nodes.
		if props.IsALink() {

			constructedRef := fmt.Sprintf("%s/%s@%s", props.Link.RegistryHint, props.Link.NamespaceHint, desc.Descriptor().Digest.String())
			if err := loadCollection(ctx, graph, constructedRef, desc.Descriptor(), client); err != nil {
				return nil, err
			}

			return graph.From(node.ID()), nil
		}

		return successors, err
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return sbom.SBOM{}, err
	}

	for _, edge := range graph.Edges() {
		parent := edge.From()
		parentIdentifier := identifier{parent.ID()}
		childIdentifier := identifier{edge.To().ID()}
		relationship := artifact.Relationship{
			From: parentIdentifier,
			To:   childIdentifier,
			Type: artifact.DependencyOfRelationship,
		}
		inventory.Relationships = append(inventory.Relationships, relationship)
	}

	catalog := pkg.NewCatalog(packages...)
	inventory.Artifacts.PackageCatalog = catalog
	versionBuilder := new(strings.Builder)
	if err := version.GetVersion(versionBuilder); err != nil {
		return sbom.SBOM{}, err
	}
	inventory.Descriptor = sbom.Descriptor{
		Name:    components.ApplicationName,
		Version: versionBuilder.String(),
	}

	return inventory, nil
}

func loadCollection(ctx context.Context, graph *collection.Collection, referenece string, rootDesc ocispec.Descriptor, client registryclient.Remote) error {
	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return client.GetContent(ctx, referenece, desc)
	}
	return collectionloader.LoadFromManifest(ctx, graph, fetcherFn, rootDesc)
}

// identifier implement the syft.Identifiable interface.
type identifier struct {
	id string
}

func (i identifier) ID() artifact.ID {
	return artifact.ID(i.id)
}
