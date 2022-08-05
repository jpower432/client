package cli

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"

	"github.com/uor-framework/uor-client-go/util/examples"
)

// BuildOptions describe configuration options that can
// be set using the build subcommand.
type BuildOptions struct {
	*RootOptions
	Destination string
	Insecure    bool
	PlainHTTP   bool
	Configs     []string
}

var clientBuildExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build collection."},
		CommandString: "build collection my-directory localhost:5000/myartifacts:latest",
	},
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build schema."},
		CommandString: "build schema schema-config.yaml localhost:5000/myartifacts:latest",
	},
}

// NewBuildCmd creates a new cobra.Command for the build subcommand.
func NewBuildCmd(rootOpts *RootOptions) *cobra.Command {
	o := BuildOptions{RootOptions: rootOpts}

	cmd := &cobra.Command{
		Use:           "build",
		Short:         "Build and save an OCI artifact from files",
		Example:       examples.FormatExamples(clientBuildExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.Flags().StringArrayVarP(&o.Configs, "configs", "c", o.Configs, "auth config paths when contacting registries")
	cmd.Flags().BoolVarP(&o.Insecure, "insecure", "", o.Insecure, "allow connections to registries SSL registry without certs")
	cmd.Flags().BoolVarP(&o.PlainHTTP, "plain-http", "", o.PlainHTTP, "use plain http and not https when contacting registries")

	cmd.AddCommand(NewBuildSchemaCmd(&o))
	cmd.AddCommand(NewBuildCollectionCmd(&o))

	return cmd
}
