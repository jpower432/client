package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	load "github.com/uor-framework/uor-client-go/builder/config"
	"github.com/uor-framework/uor-client-go/content/layout"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// BuildSchemaOptions describe configuration options that can
// be set using the build subcommand.
type BuildSchemaOptions struct {
	*BuildOptions
	SchemaConfig string
}

var clientBuildSchemaExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		Descriptions:  []string{"Build schema artifacts."},
		CommandString: "build schema schema-config.yaml localhost:5000/myartifacts:latest",
	},
}

// NewBuildSchemaCmd creates a new cobra.Command for the build subcommand.
func NewBuildSchemaCmd(buildOpts *BuildOptions) *cobra.Command {
	o := BuildSchemaOptions{BuildOptions: buildOpts}

	cmd := &cobra.Command{
		Use:           "schema DST",
		Short:         "Build and save an OCI artifact from files",
		Example:       examples.FormatExamples(clientBuildSchemaExamples...),
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

func (o *BuildSchemaOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.SchemaConfig = args[0]
	o.Destination = args[1]
	return nil
}

func (o *BuildSchemaOptions) Validate() error {
	info, err := os.Stat(o.SchemaConfig)
	if err != nil {
		return fmt.Errorf("schema configuration %q: %v", o.SchemaConfig, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("schema configuration %q: file is not regular", o.SchemaConfig)
	}
	return nil
}

func (o *BuildSchemaOptions) Run(ctx context.Context) error {

	config, err := load.ReadSchemaConfig(o.SchemaConfig)
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

	descs, err := client.AddFiles(ctx, "", config.JSONSchemaPath)
	if err != nil {
		return err
	}

	// Add the attributes from the config to their respective blocks
	configDesc, err := client.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
	if err != nil {
		return err
	}

	manifestAnnotation := map[string]string{}
	
	_, err = client.AddManifest(ctx, o.Destination, configDesc, manifestAnnotation, descs...)
	if err != nil {
		return err
	}

	desc, err := client.Save(ctx, o.Destination, cache)
	if err != nil {
		return fmt.Errorf("client save error for reference %s: %v", o.Destination, err)
	}

	o.Logger.Infof("Schema %s built with reference name %s\n", desc.Digest, o.Destination)

	return client.Destroy()
}
