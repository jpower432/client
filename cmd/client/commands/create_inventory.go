package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anchore/syft/syft/formats/spdx22json"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
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

	collection, err := client.LoadCollection(ctx, o.Source)
	if err != nil {
		return err
	}

	inventory, err := collectionToInventory(ctx, collection, client)
	if err != nil {
		return err
	}

	formatter := spdx22json.Format()
	return formatter.Encode(o.IOStreams.Out, inventory)
}

func collectionToInventory(ctx context.Context, graph collection.Collection, client registryclient.Remote) (sbom.SBOM, error) {
	var inventory sbom.SBOM
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

		// Load link and provide access to those nodes.
		if desc.Properties != nil && desc.Properties.IsALink() {
			constructedRef := fmt.Sprintf("%s/%s@%s", desc.Properties.Link.RegistryHint, desc.Properties.Link.NamespaceHint, desc.ID())
			linkedCollection, err := client.LoadCollection(ctx, constructedRef)
			if err != nil {
				return nil, err
			}
			return linkedCollection.Nodes(), nil
		}

		successors := graph.From(node.ID())
		return successors, err
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return sbom.SBOM{}, err
	}

	catalog := pkg.NewCatalog(packages...)
	inventory.Artifacts.PackageCatalog = catalog

	return inventory, nil
}
