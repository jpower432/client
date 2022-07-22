package cli

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/uor-framework/client/cli/log"
)

// RootOptions describe global configuration options that can be set.
type RootOptions struct {
	IOStreams genericclioptions.IOStreams
	LogLevel  string
	Logger    log.Logger
	cacheDir  string
}

var clientLong = templates.LongDesc(
	`
	This client helps you build and publish UOR collections as an OCI artifact.
	`,
)

// NewRootCmd creates a new cobra.Command for the command root.
func NewRootCmd() *cobra.Command {
	o := RootOptions{}

	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         "UOR Client",
		Long:          clientLong,
		SilenceErrors: false,
		SilenceUsage:  false,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			logger, err := log.NewLogger(o.IOStreams.Out, o.LogLevel)
			if err != nil {
				return err
			}
			o.Logger = logger

			home, err := homedir.Dir()
			if err != nil {
				return err
			}
			o.cacheDir = filepath.Join(home, ".uor")

			return os.MkdirAll(o.cacheDir, 0750)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	f := cmd.PersistentFlags()
	f.StringVarP(&o.LogLevel, "loglevel", "l", "info",
		"Log level (debug, info, warn, error, fatal)")

	// TODO(sabre1041) Reenable/remove once build capability strategy determined
	//cmd.AddCommand(NewRenderCmd(&o))
	cmd.AddCommand(NewBuildCmd(&o))
	cmd.AddCommand(NewPushCmd(&o))
	cmd.AddCommand(NewPullCmd(&o))
	cmd.AddCommand(NewVersionCmd(&o))

	return cmd
}
