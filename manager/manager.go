package manager

import (
	"context"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/util/workspace"
)

// Manager defines methods for building, publishing, and retrieving UOR collections.
type Manager interface {
	// Build builds collection from input and store it in the underlying content store.
	// If successful, the root descriptor is returned.
	Build(ctx context.Context, source workspace.Workspace, config clientapi.DataSetConfiguration, destination string, client registryclient.Client) (string, error)
	// Push pushes collection to a remote location from the underlying content store.
	// If successful, the root descriptor is returned.
	Push(ctx context.Context, destination string, remote registryclient.Remote) (string, error)
	// List walks a collection and returns the Node metadata
	List(ctx context.Context, source string, remote registryclient.Remote) (*collection.Collection, error)
	// Pull pulls a single collection to a specified storage destination.
	// If successful, the file locations are returned.
	Pull(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error)
	// ReadLayer reads a layer of a collection given an ocispec.AnnotationTitle value.
	// If successful, the bytes of the layer's data are returned.
	ReadLayer(ctx context.Context, source string, title string, client registryclient.Remote) ([]byte, error)
	// QueryLinks queries the attributes' endpoint for links and filters result by the provided matcher.
	QueryLinks(ctx context.Context, host string, digest string, matcher model.Matcher, client registryclient.Remote) ([]ocispec.Descriptor, error)
	// PullAll pulls linked collection to a specified storage destination.
	// If successful, the file locations are returned.
	// PullAll is similar to Pull with the exception that it walks a graph of linked collections
	// starting with the source collection reference.
	PullAll(ctx context.Context, source string, remote registryclient.Remote, destination content.Store) ([]string, error)
}
