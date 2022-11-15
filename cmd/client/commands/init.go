package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// InitOptions describe configuration options that can
// be set when using the init subcommand.
type InitOptions struct {
	*options.Common
}

var clientInitExamples = []examples.Example{
	{
		RootCommand:   filepath.Base(os.Args[0]),
		CommandString: "init",
		Descriptions: []string{
			"Create empty configuration files.",
		},
	},
}

// NewInitCmd creates a new cobra.Command for the init subcommand.
func NewInitCmd(common *options.Common) *cobra.Command {
	o := InitOptions{Common: common}

	cmd := &cobra.Command{
		Use:           "init",
		Short:         "Creates default UOR configuration files",
		Example:       examples.FormatExamples(clientInitExamples...),
		SilenceErrors: false,
		SilenceUsage:  false,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run())
		},
	}

	return cmd
}

func (o InitOptions) Complete(_ []string) error {
	return nil
}

func (o *InitOptions) Validate() error {
	return nil
}

func (o *InitOptions) Run() error {
	dsConfig := clientapi.DataSetConfiguration{
		TypeMeta: clientapi.TypeMeta{
			Kind:       clientapi.DataSetConfigurationKind,
			APIVersion: clientapi.GroupVersion,
		},
	}
	dsConfigJSON, err := json.Marshal(dsConfig)
	if err != nil {
		return err
	}
	dsConfigYAML, err := yaml.JSONToYAML(dsConfigJSON)
	if err != nil {
		return err
	}

	schemaConfig := clientapi.SchemaConfiguration{
		TypeMeta: clientapi.TypeMeta{
			Kind:       clientapi.SchemaConfigurationKind,
			APIVersion: clientapi.GroupVersion,
		},
	}
	schemaConfigJSON, err := json.Marshal(schemaConfig)
	if err != nil {
		return err
	}
	schemaConfigYAML, err := yaml.JSONToYAML(schemaConfigJSON)
	if err != nil {
		return err
	}

	attributeQuery := clientapi.AttributeQuery{
		TypeMeta: clientapi.TypeMeta{
			Kind:       clientapi.AttributeQueryKind,
			APIVersion: clientapi.GroupVersion,
		},
	}
	attributeQueryJSON, err := json.Marshal(attributeQuery)
	if err != nil {
		return err
	}
	attributeQueryYAML, err := yaml.JSONToYAML(attributeQueryJSON)
	if err != nil {
		return err
	}

	fmt.Fprintln(o.IOStreams.Out, string(dsConfigYAML))
	fmt.Fprintln(o.IOStreams.Out, string(schemaConfigYAML))
	fmt.Fprintln(o.IOStreams.Out, string(attributeQueryYAML))
	return nil
}
