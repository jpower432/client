package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/manager/defaultmanager"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// PushOptions describe configuration options that can
// be set using the push subcommand.
type PushOptions struct {
	*options.Common
	options.Remote
	Destination string
	DSConfig    string
}

var clientPushExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Push artifacts."},
	CommandString: "push localhost:5000/myartifacts:latest",
}

// NewPushCmd creates a new cobra.Command for the push subcommand.
func NewPushCmd(common *options.Common) *cobra.Command {
	o := PushOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "push DST",
		Short:         "Push a UOR collection into a registry",
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

	return cmd
}

func (o *PushOptions) Complete(args []string) error {
	if len(args) < 1 {
		return errors.New("bug: expecting one argument")
	}
	o.Destination = args[0]
	return nil
}

func (o *PushOptions) Validate() error {
	return nil
}

func (o *PushOptions) Run(ctx context.Context) error {
	cache, err := layout.NewWithContext(ctx, o.CacheDir)
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
	defer func() {
		if err := client.Destroy(); err != nil {
			o.Logger.Errorf(err.Error())
		}
	}()

	manager := defaultmanager.New(cache, o.Logger)
	_, err = manager.Push(ctx, o.Destination, client)
	return err
}
