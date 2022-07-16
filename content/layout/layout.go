package layout

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"

	"github.com/uor-framework/client/content"
	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/model/nodes/collection"
	"github.com/uor-framework/client/ocimanifest"
)

var (
	_ content.Store          = &Layout{}
	_ content.AttributeStore = &Layout{}
	_ content.GraphStore     = &Layout{}
)

const indexFile = "index.json"

// Layout implements the storage interface by wrapping the oras
// content.Storage.
type Layout struct {
	internal orascontent.Storage
	resolver *content.Resolver
	graph    *collection.Collection
	index    *ocispec.Index
	rootPath string
}

// New initializes a new local file store in an OCI layout format.
func New(ctx context.Context, rootPath string) (*Layout, error) {
	l := &Layout{
		internal: oci.NewStorage(rootPath),
		resolver: content.NewResolver(),
		graph:    collection.New(rootPath),
		rootPath: filepath.Clean(rootPath),
	}

	return l, l.init(ctx)
}

// init performs initial layout checks and loads the index.
func (l *Layout) init(ctx context.Context) error {
	if err := l.validateOCILayoutFile(); err != nil {
		return err
	}
	return l.loadIndex(ctx)
}

// Fetch fetches the content identified by the descriptor.
func (l *Layout) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {
	return l.internal.Fetch(ctx, desc)
}

// Push pushes the content, matching the expected descriptor.
func (l *Layout) Push(ctx context.Context, desc ocispec.Descriptor, content io.Reader) error {
	return l.internal.Push(ctx, desc, content)
}

// Exists returns whether a descriptor exits in the file store.
func (l *Layout) Exists(ctx context.Context, desc ocispec.Descriptor) (bool, error) {
	return l.internal.Exists(ctx, desc)
}

// Resolve resolves a reference to a descriptor.
func (l *Layout) Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error) {
	return l.resolver.Resolve(ctx, reference)
}

// Successors returns the nodes directly pointed by the current node.
// In other words, returns the "children" of the current descriptor.
func (l *Layout) Successors(_ context.Context, node ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	fmt.Println("not implemented")
	return nil, nil
}

// Predecessors returns the nodes directly pointing to the current node.
func (l *Layout) Predecessors(_ context.Context, node ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	fmt.Println("not implemented")
	return nil, nil
}

// ResolveByAttribute returns descriptors linked to the reference that satisfy the specified matcher.
// Matcher is expected to compare attributes of nodes to set criteria.
func (l *Layout) ResolveByAttribute(ctx context.Context, reference string, matcher model.Matcher) ([]ocispec.Descriptor, error) {
	fmt.Println("not implemented")
	return nil, nil
}

// ResolveLinks returns linked collection references for a collection. If the collection
// has no links, nil is returned.
func (l *Layout) ResolveLinks(ctx context.Context, reference string) ([]string, error) {
	desc, err := l.Resolve(ctx, reference)
	if err != nil {
		return nil, err
	}
	r, err := l.Fetch(ctx, desc)
	if err != nil {
		return nil, errdef.ErrInvalidReference
	}
	var manifest ocispec.Manifest
	if err := json.NewDecoder(r).Decode(&manifest); err != nil {
		return nil, err
	}
	links, ok := manifest.Annotations[ocimanifest.AnnotationCollectionLinks]
	if !ok {
		return nil, nil
	}
	splitLinks := strings.Split(links, ocimanifest.Separator)
	return splitLinks, ocimanifest.ErrNoCollectionLinks
}

// Tag tags a descriptor with a reference string.
// A reference should be either a valid tag (e.g. "latest"),
// or a digest matching the descriptor (e.g. "@sha256:abc123").
func (l *Layout) Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error {
	if err := validateReference(reference); err != nil {
		return fmt.Errorf("invalid reference: %w", err)
	}

	exists, err := l.Exists(ctx, desc)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s: %s: %w", desc.Digest, desc.MediaType, errdef.ErrNotFound)
	}

	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}
	desc.Annotations[ocispec.AnnotationRefName] = reference

	if err := l.resolver.Tag(ctx, desc, reference); err != nil {
		return err
	}

	return l.SaveIndex()
}

// Index returns an index manifest object.
func (l *Layout) Index() (ocispec.Index, error) {
	return *l.index, nil
}

// SaveIndex writes the index.json to the file system
func (l *Layout) SaveIndex() error {
	// first need to update the index
	var descs []ocispec.Descriptor
	for name, desc := range l.resolver.Map() {
		if desc.Annotations == nil {
			desc.Annotations = map[string]string{}
		}
		desc.Annotations[ocispec.AnnotationRefName] = name
		descs = append(descs, desc)
	}

	l.index.Manifests = descs
	indexJSON, err := json.Marshal(l.index)
	if err != nil {
		return err
	}
	path := filepath.Join(l.rootPath, indexFile)
	return ioutil.WriteFile(path, indexJSON, 0640)
}

// loadIndex loads all information from the index.json
// into the resolver and graph.
func (l *Layout) loadIndex(ctx context.Context) error {
	path := filepath.Join(l.rootPath, indexFile)
	indexFile, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		l.index = &ocispec.Index{
			Versioned: specs.Versioned{
				SchemaVersion: 2,
			},
		}

		return nil
	}
	defer indexFile.Close()

	if err := json.NewDecoder(indexFile).Decode(&l.index); err != nil {
		return err
	}

	for _, d := range l.index.Manifests {
		key, ok := d.Annotations[ocispec.AnnotationRefName]
		if ok {
			if err := l.resolver.Tag(ctx, d, key); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateOCILayoutFile ensure the 'oci-layout' file exists in the
// root directory and contains a valid version.
func (l *Layout) validateOCILayoutFile() error {
	layoutFilePath := filepath.Join(l.rootPath, ocispec.ImageLayoutFile)
	layoutFile, err := os.Open(layoutFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to open OCI layout file: %w", err)
		}

		layout := ocispec.ImageLayout{
			Version: ocispec.ImageLayoutVersion,
		}
		layoutJSON, err := json.Marshal(layout)
		if err != nil {
			return fmt.Errorf("failed to marshal OCI layout file: %w", err)
		}

		return ioutil.WriteFile(layoutFilePath, layoutJSON, 0666)
	}
	defer layoutFile.Close()

	var layout *ocispec.ImageLayout
	err = json.NewDecoder(layoutFile).Decode(&layout)
	if err != nil {
		return fmt.Errorf("failed to decode OCI layout file: %w", err)
	}
	if layout.Version != ocispec.ImageLayoutVersion {
		return errdef.ErrUnsupportedVersion
	}

	return nil
}

// validateReference ensures the build reference
// contains a tag component.
func validateReference(name string) error {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 1 {
		return fmt.Errorf("reference %q: missing repository", name)
	}
	path := parts[1]
	if index := strings.Index(path, "@"); index != -1 {
		return fmt.Errorf("%w: ", errdef.ErrInvalidReference)
	} else if index := strings.Index(path, ":"); index != -1 {
		// tag found
		return nil
	} else {
		// empty reference
		return fmt.Errorf("reference %q: missing tag component", name)
	}
}
