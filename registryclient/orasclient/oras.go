package orasclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	orascontent "oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/uor-framework/uor-client-go/content"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/collection"
	collectionloader "github.com/uor-framework/uor-client-go/nodes/collection/loader"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
	"github.com/uor-framework/uor-client-go/registryclient/orasclient/internal/cache"
)

type orasClient struct {
	plainHTTP     bool
	insecure      bool
	authCache     auth.Cache
	credFn        func(context.Context, string) (auth.Credential, error)
	registryConf  registryclient.RegistryConfig
	copyOpts      oras.CopyOptions
	artifactStore *file.Store
	cache         content.Store
	// collection will store a cache of
	// loaded collections from remote sources.
	collections sync.Map // map[string]collection.Collection
	destroy     func() error
	outputDir   string
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

	var packOpts oras.PackOptions
	packOpts.ConfigDescriptor = &configDesc
	packOpts.ManifestAnnotations = manifestAnnotations

	manifestDesc, err := oras.Pack(ctx, c.artifactStore, descriptors, packOpts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	return manifestDesc, c.artifactStore.Tag(ctx, manifestDesc, ref)
}

// Save saves the OCI artifact to local store location (e.g. cache)
func (c *orasClient) Save(ctx context.Context, ref string, store content.Store) (ocispec.Descriptor, error) {
	return oras.Copy(ctx, c.artifactStore, ref, store, ref, c.copyOpts)
}

// LoadCollection loads a UOR collection type from a remote registry path.
func (c *orasClient) LoadCollection(ctx context.Context, reference string) (collection.Collection, error) {
	value, exists := c.collections.Load(reference)
	if exists {
		return value.(collection.Collection), nil
	}

	repo, updatedRef, err := c.setupRepo(ctx, reference)
	if err != nil {
		return collection.Collection{}, fmt.Errorf("could not create registry target: %w", err)
	}

	graph, err := loadCollection(ctx, repo, updatedRef)
	if err != nil {
		return collection.Collection{}, err
	}
	c.collections.Store(reference, graph)

	return graph, nil
}

// Pull performs a copy of OCI artifacts to a local location from a remote location.
func (c *orasClient) Pull(ctx context.Context, reference string, store content.Store) (ocispec.Descriptor, []ocispec.Descriptor, error) {
	var allDescs []ocispec.Descriptor

	var from oras.Target
	repo, updatedRef, err := c.setupRepo(ctx, reference)
	if err != nil {
		return ocispec.Descriptor{}, allDescs, fmt.Errorf("could not create registry target: %w", err)
	}
	from = repo

	if c.cache != nil {
		from = cache.New(repo, c.cache)
	}

	// Load the collection manifest from remote or
	// the collection cache
	var graph collection.Collection
	value, exists := c.collections.Load(reference)
	if exists {
		graph = value.(collection.Collection)
	} else {
		graph, err = loadCollection(ctx, repo, updatedRef)
		if err != nil {
			return ocispec.Descriptor{}, allDescs, err
		}
		c.collections.Store(reference, graph)
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
			desc, ok := node.(*descriptor.Node)
			if !ok {
				return false, nil
			}

			switch desc.Descriptor().MediaType {
			case ocimanifest.UORSchemaMediaType:
				return true, nil
			case ocispec.MediaTypeImageConfig:
				return true, nil
			case ocimanifest.UORConfigMediaType:
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
			d, ok := s.(*descriptor.Node)
			if ok {
				result = append(result, d.Descriptor())
			}
		}
		return result, nil
	}

	// Create a copy of the options so the original copy
	// options are not modified.
	cCopyOpts := c.copyOpts
	cCopyOpts.FindSuccessors = successorFn

	desc, err := oras.Copy(ctx, from, updatedRef, store, updatedRef, cCopyOpts)
	if err != nil {
		return ocispec.Descriptor{}, allDescs, err
	}

	return desc, allDescs, nil
}

// Push performs a copy of OCI artifacts to a remote location.
func (c *orasClient) Push(ctx context.Context, store content.Store, reference string) (ocispec.Descriptor, error) {
	repo, updatedRef, err := c.setupRepo(ctx, reference)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("could not create registry target: %w", err)
	}
	return oras.Copy(ctx, store, updatedRef, repo, updatedRef, c.copyOpts)
}

// GetManifest returns the manifest the reference resolves to.
func (c *orasClient) GetManifest(ctx context.Context, reference string) (ocispec.Descriptor, io.ReadCloser, error) {
	repo, updatedRef, err := c.setupRepo(ctx, reference)
	if err != nil {
		return ocispec.Descriptor{}, nil, fmt.Errorf("could not create registry target: %w", err)
	}
	return repo.FetchReference(ctx, updatedRef)
}

// GetContent retrieves the content for a specified descriptor at a specified reference.
func (c *orasClient) GetContent(ctx context.Context, reference string, desc ocispec.Descriptor) ([]byte, error) {
	repo, _, err := c.setupRepo(ctx, reference)
	if err != nil {
		return nil, fmt.Errorf("could not create registry target: %w", err)
	}
	return getContent(ctx, desc, repo)
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

// authClient configures a new auth client with a given configuration.
func (c *orasClient) authClient(insecure bool) *auth.Client {
	return &auth.Client{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
		Cache:      c.authCache,
		Credential: c.credFn,
	}
}

// setupRepo configures the client to access the remote repository.
func (c *orasClient) setupRepo(ctx context.Context, reference string) (registry.Repository, string, error) {
	registryConfig, err := registryclient.FindRegistry(c.registryConf, reference)
	if err != nil {
		return nil, reference, err
	}

	repo, err := remote.NewRepository(reference)
	if err != nil {
		return nil, reference, fmt.Errorf("could not create registry target: %w", err)
	}

	// If the incoming reference does not match any registry prefixes,
	// use the client configuration. If there is a match and the registry
	// has mirrors configured, try each one before attempting to contact the
	// input reference (the default).
	switch {
	case registryConfig == nil:
		repo.PlainHTTP = c.plainHTTP
		repo.Client = c.authClient(c.insecure)
		return repo, reference, nil
	case len(registryConfig.Mirrors) != 0:
		mirrorRepo, mirrorReference, err := c.pickMirror(ctx, *registryConfig, reference)
		if err == nil {
			return mirrorRepo, mirrorReference, nil
		}

		var merr *registryclient.ErrNoAvailableMirrors
		if err != nil && !errors.As(err, &merr) {
			return nil, reference, err
		}

		fallthrough
	default:
		repo.PlainHTTP = registryConfig.PlainHTTP
		repo.Client = c.authClient(registryConfig.SkipTLS)
		return repo, reference, nil
	}
}

// pickMirror is used if the reference is linked to a registry in the registry configuration that has mirrors. It returns
// a configured remote.Repository and rewritten reference per the mirror configuration.
func (c *orasClient) pickMirror(ctx context.Context, reg registryclient.Registry, ref string) (registry.Repository, string, error) {
	pullSources, err := reg.PullSourceFromReference(ref)
	if err != nil {
		return nil, ref, err
	}

	for _, ps := range pullSources {
		mirror := ps.Endpoint
		mirrorReg, err := remote.NewRegistry(mirror.Location)
		if err != nil {
			return nil, ref, fmt.Errorf("could not create registry target: %w", err)
		}
		mirrorReg.PlainHTTP = mirror.PlainHTTP
		mirrorReg.Client = c.authClient(mirror.SkipTLS)

		if err := mirrorReg.Ping(ctx); err == nil {
			mirrorReference, err := registry.ParseReference(ps.Reference)
			if err != nil {
				return nil, ref, fmt.Errorf("reference %q: %w", mirrorReference, err)
			}
			mirrorRepo, err := mirrorReg.Repository(ctx, mirrorReference.Repository)
			if err != nil {
				return nil, ref, err
			}
			return mirrorRepo, mirrorReference.String(), nil
		}
	}

	return nil, ref, &registryclient.ErrNoAvailableMirrors{Registry: reg.Location}
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

// loadCollection is a helper function that allows a collection to be loaded with a given repository.
func loadCollection(ctx context.Context, repo registry.Repository, reference string) (collection.Collection, error) {
	desc, err := repo.Resolve(ctx, reference)
	if err != nil {
		return collection.Collection{}, err
	}
	fetcherFn := func(ctx context.Context, desc ocispec.Descriptor) ([]byte, error) {
		return getContent(ctx, desc, repo)
	}
	co := collection.New(reference)
	if err := collectionloader.LoadFromManifest(ctx, co, fetcherFn, desc); err != nil {
		return collection.Collection{}, err
	}
	co.Location = reference
	return *co, nil
}

// getContent is a helper function that allows content to be retrieved with a given repository.
func getContent(ctx context.Context, desc ocispec.Descriptor, repo registry.Repository) ([]byte, error) {
	r, err := repo.Fetch(ctx, desc)
	if err != nil {
		return nil, err
	}
	return orascontent.ReadAll(r, desc)
}
