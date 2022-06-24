package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/content/local"
	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/uor-framework/client/model"
)

var (
	_ Storage          = &Local{}
	_ remotes.Resolver = &Local{}
)

const indexFile = "index.json"

// Local implements the storage interface by wrapping the containerd
// content.Store.
// https://pkg.go.dev/github.com/containerd/containerd@v1.6.1/content#Store
// an OCI index.json manifest is used to track manifest metadata for attribute
// addressing use.
type Local struct {
	content.Store
	descriptorLookup map[string]ocispec.Descriptor
	index            ocispec.Index
	rootPath         string
}

func NewLocal(rootPath string) (*Local, error) {
	fileStore, err := local.NewStore(rootPath)
	if err != nil {
		return nil, err
	}
	l := &Local{
		Store:            fileStore,
		descriptorLookup: make(map[string]ocispec.Descriptor),
		rootPath:         rootPath,
	}
	return l, l.init()
}

func (l *Local) init() error {
	if err := l.ensureOCILayoutFile(); err != nil {
		return err
	}
	ii, err := l.Index()
	if err != nil {
		return err
	}
	for _, d := range ii.Manifests {
		key := d.Annotations[ocispec.AnnotationTitle]
		l.descriptorLookup[key] = d
	}
	return nil
}

// Storage methods

func (l *Local) Add(ctx context.Context, descs ...ocispec.Descriptor) error {
	for _, desc := range descs {
		key := desc.Annotations[ocispec.AnnotationTitle]
		l.descriptorLookup[key] = desc
		w, err := l.Store.Writer(ctx, content.WithDescriptor(desc))
		if err != nil {
			return err
		}
		if err := w.Commit(ctx, 0, "", nil); err != nil {
			return err
		}
	}

	return l.saveIndex()
}

func (l *Local) Delete(ctx context.Context, descs ...ocispec.Descriptor) error {
	for _, desc := range descs {
		delete(l.descriptorLookup, desc.Annotations[ocispec.AnnotationTitle])
		if err := l.Store.Delete(ctx, desc.Digest); err != nil {
			return err
		}
	}
	return l.saveIndex()
}

func (l *Local) Index() (ocispec.Index, error) {
	return l.index, nil
}

func (l *Local) List() []ocispec.Descriptor {
	var descs []ocispec.Descriptor
	for _, desc := range l.descriptorLookup {
		descs = append(descs, desc)
	}
	return descs
}

func (l *Local) Exists(desc ocispec.Descriptor) bool {
	key := desc.Annotations[ocispec.AnnotationTitle]
	_, exits := l.descriptorLookup[key]
	return exits
}

func (l *Local) LookupByAttribute(attributes model.Attributes) ([]ocispec.Descriptor, error) {
	return nil, nil
}

// remotes.Resolver methods

func (l *Local) Fetcher(ctx context.Context, ref string) (remotes.Fetcher, error) {
	return l, nil
}

func (l *Local) Fetch(ctx context.Context, desc ocispec.Descriptor) (io.ReadCloser, error) {
	return nil, nil
}

func (l *Local) Pusher(ctx context.Context, ref string) (remotes.Pusher, error) {
	return l, nil
}

func (l *Local) Push(ctx context.Context, d ocispec.Descriptor) (content.Writer, error) {
	return l.Store.Writer(ctx, content.WithDescriptor(d))
}

func (l *Local) Resolve(ctx context.Context, ref string) (name string, desc ocispec.Descriptor, err error) {
	return "", ocispec.Descriptor{}, err
}

// ensureOCILayoutFile ensures the `oci-layout` file.
func (l *Local) ensureOCILayoutFile() error {
	layoutFilePath := filepath.Join(l.rootPath, indexFile)
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
		return errors.New("unsupported version")
	}

	return nil
}

// saveIndex writes the index.json to the file system
func (l *Local) saveIndex() error {
	// first need to update the index
	var descs []ocispec.Descriptor
	for name, desc := range l.descriptorLookup {
		if desc.Annotations == nil {
			desc.Annotations = map[string]string{}
		}
		desc.Annotations[ocispec.AnnotationTitle] = name
		descs = append(descs, desc)
	}
	l.index.Manifests = descs
	indexJSON, err := json.Marshal(l.index)
	if err != nil {
		return err
	}

	path := filepath.Join(l.rootPath, indexFile)
	return ioutil.WriteFile(path, indexJSON, 0644)
}
