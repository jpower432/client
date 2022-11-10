package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/config"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient"
	"github.com/uor-framework/uor-client-go/schema"
	"github.com/uor-framework/uor-client-go/util/examples"
)

// AggregateOptions describe configuration options that can
// be set using the push subcommand.
type AggregateOptions struct {
	*CreateOptions
	AttributeQuery string
	RegistryHost   string
	SchemaID       string
}

var clientAggregateExamples = examples.Example{
	RootCommand:   filepath.Base(os.Args[0]),
	Descriptions:  []string{"Build aggregate from a query."},
	CommandString: "aggregate localhost:5001 myquery.yaml",
}

// NewAggregateCmd creates a new cobra.Command for the aggregate subcommand.
func NewAggregateCmd(createOps *CreateOptions) *cobra.Command {
	o := AggregateOptions{CreateOptions: createOps}

	cmd := &cobra.Command{
		Use:           "aggregate HOST QUERY",
		Short:         "Create an artifact aggregate from an attribute query",
		Example:       examples.FormatExamples(clientAggregateExamples),
		SilenceErrors: false,
		SilenceUsage:  false,
		Args:          cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Run(cmd.Context()))
		},
	}

	cmd.Flags().StringVarP(&o.SchemaID, "schema-id", "s", schema.UnknownSchemaID, "Schema ID to scope attribute query. Default is \"unknown\"")

	return cmd
}

func (o *AggregateOptions) Complete(args []string) error {
	if len(args) < 2 {
		return errors.New("bug: expecting two arguments")
	}
	o.RegistryHost = args[0]
	o.AttributeQuery = args[1]
	return nil
}

func (o *AggregateOptions) Validate() error {
	return nil
}

func (o *AggregateOptions) Run(ctx context.Context) error {
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

	userQuery, err := config.ReadAttributeQuery(o.AttributeQuery)
	if err != nil {
		return err
	}

	if len(userQuery.Attributes) != 0 {
		// Based on the way descriptor are written the schema is always the root key.
		constructedQuery := map[string]v1alpha1.Attributes{
			o.SchemaID: userQuery.Attributes,
		}
		queryJSON, err := json.Marshal(constructedQuery)
		if err != nil {
			return err
		}

		result, err := client.ResolveAttributeQuery(ctx, o.RegistryHost, queryJSON)
		if err != nil {
			return err
		}

		resultsJSON, err := json.MarshalIndent(result, "", "")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	if len(userQuery.Digests) != 0 {
		result, err := client.ResolveDigestQuery(ctx, o.RegistryHost, userQuery.Digests)
		if err != nil {
			return err
		}
		resultsJSON, err := json.MarshalIndent(result, "", "")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	if len(userQuery.Links) != 0 {
		result, err := client.ResolveLinkQuery(ctx, o.RegistryHost, userQuery.Links)
		if err != nil {
			return err
		}
		resultsJSON, err := json.MarshalIndent(result, "", "")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.IOStreams.Out, string(resultsJSON))
	}

	return nil
}
