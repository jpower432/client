package orasclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2"
	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/model/traversal"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	collectionloader "github.com/uor-framework/uor-client-go/nodes/collection/loader"
	"github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient/internal/cache"
)

type orasClient struct {
	plainHTTP  bool
	authClient *auth.Client
	// oras-specific copy options
	copyOpts oras.CopyOptions
	// User specified pre-pull actions
	// (e.g. signature verification).
	prePullFn func(context.Context, string) error
	// underlying store for collection
	// building on disk
	artifactStore *file.Store
	// Location of cached blobs
	cache content.Store
	// collection will store a cache of
	// loaded collections from remote sources.
	collections sync.Map // map[string]collection.Collection
	destroy     func() error
	// attributes is set to filter
	// collections by attribute.
	attributes model.Matcher
}

var _ registryclient.Client = &orasClient{}

// AddFiles loads one or more files to create OCI descriptors with a specific
// media type and pushes them into underlying storage.
func (c *orasClient) AddFiles(ctx context.Context, mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return nil, err
	}
	descs, err := loadFiles(ctx, c.artifactStore, mediaType, files...)
	if err != nil {
		return nil, fmt.Errorf("unable to load files: %w", err)
	}
	return descs, nil
}

// AddContent creates and stores a descriptor from content in bytes, a media type, and
// annotations.
func (c *orasClient) AddContent(ctx context.Context, mediaType string, content []byte, annotations map[string]string) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}
	configDesc := ocispec.Descriptor{
		MediaType:   mediaType,
		Digest:      digest.FromBytes(content),
		Size:        int64(len(content)),
		Annotations: annotations,
	}

	return configDesc, c.artifactStore.Push(ctx, configDesc, bytes.NewReader(content))
}

// AddManifest creates and stores a manifest.
// This is generated from the config descriptor and artifact descriptors.
func (c *orasClient) AddManifest(ctx context.Context, ref string, configDesc ocispec.Descriptor, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}
	if descriptors == nil {
		descriptors = []ocispec.Descriptor{}
	}

	// Keep descriptor order deterministic
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Digest < descriptors[j].Digest
	})

	var packOpts PackOptions
	packOpts.ManifestAnnotations = manifestAnnotations
	packOpts.DisableTimestamp = true
	packOpts.PackImageManifest = true
	packOpts.ConfigDescriptor = &configDesc

	manifestDesc, err := Pack(ctx, c.artifactStore, "", descriptors, packOpts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, c.artifactStore.Tag(ctx, manifestDesc, ref)
}

// AddIndex creates and stores am index manifest.
func (c *orasClient) AddIndex(ctx context.Context, ref string, manifestAnnotations map[string]string, descriptors ...ocispec.Descriptor) (ocispec.Descriptor, error) {
	if err := c.checkFileStore(); err != nil {
		return ocispec.Descriptor{}, err
	}
	if descriptors == nil {
		descriptors = []ocispec.Descriptor{}
	}

	// Keep descriptor order deterministic
	sort.Slice(descriptors, func(i, j int) bool {
		return descriptors[i].Digest < descriptors[j].Digest
	})

	var packOpts PackOptions
	packOpts.ManifestAnnotations = manifestAnnotations
	packOpts.DisableTimestamp = true

	manifestDesc, err := PackAggregate(ctx, c.artifactStore, descriptors, packOpts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, c.artifactStore.Tag(ctx, manifestDesc, ref)
}

// Save saves the OCI artifact to local store location (e.g. cache)
func (c *orasClient) Save(ctx context.Context, ref string, store content.Store) (ocispec.Descriptor, error) {
	// Create a copy of the options so the original copy
	// options are not modified.
	cCopyOpts := c.copyOpts
	cCopyOpts.FindSuccessors = successorFnWithSparseManifests
	return oras.Copy(ctx, c.artifactStore, ref, store, ref, cCopyOpts)
}

// LoadCollection loads a UOR collection type from a remote registry path.
func (c *orasClient) LoadCollection(ctx context.Context, reference string) (collection.Collection, error) {
	value, exists := c.collections.Load(reference)
	if exists {
		return value.(collection.Collection), nil
	}

	rootDesc, manifestBytes, err := c.GetManifest(ctx, reference)
	if err != nil {
		return collection.Collection{}, err
	}

	// Get manifest information to obtain annotations
	var manifest ocispec.Descriptor
	if err := json.NewDecoder(manifestBytes).Decode(&manifest); err != nil {
		return collection.Collection{}, err
	}
	rootDesc.Annotations = manifest.Annotations
	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return c.GetContent(ctx, reference, desc)
	}

	co := collection.New(reference)
	if err := collectionloader.LoadFromManifest(ctx, co, fetcherFn, rootDesc); err != nil {
		return collection.Collection{}, err
	}
	co.Location = reference
	c.collections.Store(reference, *co)
	return *co, nil
}

// PullWithLinks performs a copy of OCI artifacts to a local location from a remote location and follow links to
// other artifacts.
func (c *orasClient) PullWithLinks(ctx context.Context, ref string, store content.Store) ([]ocispec.Descriptor, error) {
	seen := map[string]struct{}{}
	var allDescs []ocispec.Descriptor

	graph, err := c.LoadCollection(ctx, ref)
	if err != nil {
		return nil, err
	}
	root, err := graph.Root()
	if err != nil {
		return nil, err
	}

	processFunc := func(currRef string) error {
		rootDesc, descs, err := c.Pull(ctx, currRef, store)
		if err != nil {
			return err
		}
		if len(rootDesc.Digest) != 0 {
			if err := store.Tag(ctx, rootDesc, ref); err != nil {
				return err
			}
		}
		allDescs = append(allDescs, descs...)
		return nil
	}

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
			linkedCollection, err := c.LoadCollection(ctx, constructedRef)
			if err != nil {
				return nil, err
			}
			if err := processFunc(constructedRef); err != nil {
				return nil, err
			}
			return linkedCollection.Nodes(), nil
		}

		successors := graph.From(node.ID())
		return successors, err
	})

	if err := tracker.Walk(ctx, handler, root); err != nil {
		return nil, err
	}

	if err := processFunc(ref); err != nil {
		return nil, err
	}

	return allDescs, nil
}

// Pull performs a copy of OCI artifacts to a local location from a remote location.
func (c *orasClient) Pull(ctx context.Context, ref string, store content.Store) (ocispec.Descriptor, []ocispec.Descriptor, error) {
	var allDescs []ocispec.Descriptor

	if c.prePullFn != nil {
		if err := c.prePullFn(ctx, ref); err != nil {
			return ocispec.Descriptor{}, nil, err
		}
	}

	var from oras.Target
	repo, err := c.setupRepo(ref)
	if err != nil {
		return ocispec.Descriptor{}, allDescs, fmt.Errorf("could not create registry target: %w", err)
	}
	from = repo

	if c.cache != nil {
		from = cache.New(repo, c.cache)
	}

	graph, err := c.LoadCollection(ctx, ref)
	if err != nil {
		return ocispec.Descriptor{}, allDescs, err
	}

	// Filter the collection per the matcher criteria
	if c.attributes != nil {
		var matchedLeaf int
		matchFn := model.MatcherFunc(func(node model.Node) (bool, error) {
			// This check ensure we are not weeding out any manifests needed
			// for OCI DAG traversal.
			if len(graph.From(node.ID())) != 0 {
				return true, nil
			}

			// Check that this is a descriptor node and the blob is
			// not a config or schema resource.
			desc, ok := node.(*v2.Node)
			if !ok {
				return false, nil
			}

			switch desc.Descriptor().MediaType {
			case uorspec.MediaTypeSchemaDescriptor:
				return true, nil
			case ocispec.MediaTypeImageConfig:
				return true, nil
			case uorspec.MediaTypeConfiguration:
				return true, nil
			}

			match, err := c.attributes.Matches(node)
			if err != nil {
				return false, err
			}

			if match {
				matchedLeaf++
			}

			return match, nil
		})

		var err error
		graph, err = graph.SubCollection(matchFn)
		if err != nil {
			return ocispec.Descriptor{}, allDescs, err
		}

		if matchedLeaf == 0 {
			return ocispec.Descriptor{}, allDescs, nil
		}
	}

	var mu sync.Mutex
	successorFn := func(_ context.Context, fetcher orascontent.Fetcher, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		mu.Lock()
		successors := graph.From(desc.Digest.String())
		allDescs = append(allDescs, desc)
		mu.Unlock()

		var result []ocispec.Descriptor
		for _, s := range successors {
			d, ok := s.(*v2.Node)
			if ok {
				// Skip any attempts to pull a link as they could
				// be outside the repository.
				if d.Properties == nil || !d.Properties.IsALink() {
					result = append(result, d.Descriptor())
				}
			}
		}
		return result, nil
	}

	// Create a copy of the options so the original copy
	// options are not modified.
	cCopyOpts := c.copyOpts
	cCopyOpts.FindSuccessors = successorFn

	desc, err := oras.Copy(ctx, from, ref, store, ref, cCopyOpts)
	if err != nil {
		return ocispec.Descriptor{}, allDescs, err
	}

	return desc, allDescs, nil
}

// Push performs a copy of OCI artifacts to a remote location.
func (c *orasClient) Push(ctx context.Context, store content.Store, ref string) (ocispec.Descriptor, error) {
	repo, err := c.setupRepo(ref)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create registry target: %w", err)
	}

	// Create a copy of the options so the original copy
	// options are not modified.
	cCopyOpts := c.copyOpts
	cCopyOpts.FindSuccessors = successorFnWithSparseManifests

	return oras.Copy(ctx, store, ref, repo, ref, cCopyOpts)
}

// GetManifest returns the manifest the reference resolves to.
func (c *orasClient) GetManifest(ctx context.Context, reference string) (ocispec.Descriptor, io.ReadCloser, error) {
	repo, err := c.setupRepo(reference)
	if err != nil {
		return ocispec.Descriptor{}, nil, fmt.Errorf("could not create registry target: %w", err)
	}
	return repo.FetchReference(ctx, reference)
}

// GetContent retrieves the content for a specified descriptor at a specified reference.
func (c *orasClient) GetContent(ctx context.Context, reference string, desc ocispec.Descriptor) ([]byte, error) {
	repo, err := c.setupRepo(reference)
	if err != nil {
		return nil, fmt.Errorf("could not create registry target: %w", err)
	}
	r, err := repo.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}
	return orascontent.ReadAll(r, desc)
}

// Store returns the source storage being used to store
// the OCI artifact.
func (c *orasClient) Store() (content.Store, error) {
	return c.artifactStore, nil
}

// Destroy cleans up any temporary on-disk resources used to track descriptors.
func (c *orasClient) Destroy() error {
	return c.destroy()
}

// checkFileStore ensures that the file store
// has been initialized.
func (c *orasClient) checkFileStore() error {
	if c.artifactStore == nil {
		return errors.New("file store uninitialized")
	}
	return nil
}

// setupRepo configures the client to access the remote repository.
func (c *orasClient) setupRepo(ref string) (*remote.Repository, error) {
	repo, err := remote.NewRepository(ref)
	if err != nil {
		return nil, fmt.Errorf("could not create registry target: %w", err)
	}
	repo.PlainHTTP = c.plainHTTP
	repo.Client = c.authClient
	return repo, nil
}

// TODO(jpower432): PR upstream so that this can be done with the pre-copy option. Currently the error to skip descriptors is
// private https://github.com/oras-project/oras-go/blob/9e5b1419cdedd6240a5bf836c83f75270ba9d74b/copy.go#L49.

// successorFnWithSparseManifest defines a successor function to use with oras.Copy that will skip any expected linked content (i.e. sparse manifests)
func successorFnWithSparseManifests(ctx context.Context, fetcher orascontent.Fetcher, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	successors, err := orascontent.Successors(ctx, fetcher, desc)
	if err != nil {
		return nil, err
	}

	var filtered []ocispec.Descriptor
	for _, s := range successors {
		node, err := v2.NewNode(s.Digest.String(), s)
		if err != nil {
			return nil, err
		}
		if node.Properties == nil || !node.Properties.IsALink() {
			filtered = append(filtered, s)
		}

	}
	return filtered, nil
}

// loadFiles stores files in a file store and creates descriptors representing each file in the store.
func loadFiles(ctx context.Context, store *file.Store, mediaType string, files ...string) ([]ocispec.Descriptor, error) {
	var descs []ocispec.Descriptor
	var skipMediaTypeDetection bool
	var err error

	if mediaType != "" {
		skipMediaTypeDetection = true
	}
	for _, fileRef := range files {
		name := filepath.Clean(fileRef)
		if !filepath.IsAbs(name) {
			// convert to slash-separated path unless it is absolute path
			name = filepath.ToSlash(name)
		}

		if !skipMediaTypeDetection {
			mediaType, err = getDefaultMediaType(fileRef)
			if err != nil {
				return nil, fmt.Errorf("file %q: error dectecting media type: %v", name, err)
			}
		}

		desc, err := store.Add(ctx, name, mediaType, fileRef)
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}
	return descs, nil
}

// getDefaultMediaType detects the media type of the
// file based on content.
func getDefaultMediaType(file string) (string, error) {
	mType, err := mimetype.DetectFile(file)
	if err != nil {
		return "", err
	}
	return mType.String(), nil
}
