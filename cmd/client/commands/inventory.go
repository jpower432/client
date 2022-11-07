package commands

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/components"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// InventoryOptions describe configuration options that can
// be set using the push subcommand.
type InventoryOptions struct {
	*options.Common
	options.Remote
	options.RemoteAuth
	Source string
	Format string
}

var clientInventoryExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Build inventory from artifacts."},
	CommandString: "inventory localhost:5000/myartifacts:latest",
}

// NewInventoryCmd creates a new cobra.Command for the push subcommand.
func NewInventoryCmd(common *options.Common) *cobra.Command {
	o := PushOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "inventory SRC",
		Short:         "Create software inventories from UOR artifacts",
		Example:       examples.FormatExamples(clientPushExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	o.Remote.BindFlags(cmd.Flags())
	o.RemoteAuth.BindFlags(cmd.Flags())

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

	collection, err := client.LoadCollection(ctx, o.Source)
	if err != nil {
		return err
	}

	var props []descriptor.Properties
	for _, node := range collection.Nodes() {
		desc, ok := node.(*v2.Node)
		if ok {
			props = append(props, *desc.Properties)
		}
	}

	inventory := components.PropertiesToInventory(props)

	return components.InventoryToSPDXJSON(o.IOStreams.Out, inventory)
}
