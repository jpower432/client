package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/uor-framework/uor-client-go/cmd/client/commands/options"
	"github.com/uor-framework/uor-client-go/log"
)

func TestPullComplete(t *testing.T) {
	type spec struct {
		name     string
		args     []string
		opts     *PullOptions
		expOpts  *PullOptions
		expError string
	}

	cases := []spec{
		{
			name: "Valid/CorrectNumberOfArguments",
			args: []string{"test-registry.com/image:latest"},
			expOpts: &PullOptions{
				Source: "test-registry.com/image:latest",
				Output: ".",
			},
			opts: &PullOptions{},
		},
		{
			name:     "Invalid/NotEnoughArguments",
			args:     []string{},
			expOpts:  &PullOptions{},
			opts:     &PullOptions{},
			expError: "bug: expecting one argument",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Complete(c.args)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expOpts, c.opts)
			}
		})
	}
}

func TestPullValidate(t *testing.T) {
	type spec struct {
		name     string
		opts     *PullOptions
		expError string
	}

	tmp := t.TempDir()

	cases := []spec{
		{
			name: "Valid/RootDirExists",
			opts: &PullOptions{
				Output: "testdata",
			},
		},
		{
			name: "Valid/RootDirDoesNotExist",
			opts: &PullOptions{
				Output: filepath.Join(tmp, "fake"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.opts.Validate()
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				_, err = os.Stat(c.opts.Output)
				require.NoError(t, err)
			}
		})
	}
}

func TestPullRun(t *testing.T) {
	testlogr, err := log.NewLogrusLogger(io.Discard, "debug")
	require.NoError(t, err)

	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	type spec struct {
		name       string
		opts       *PullOptions
		assertFunc func(string) bool
		prepFunc   func(*testing.T, string, string)
		expError   string
	}

	cases := []spec{
		{
			name: "Success/NoAttributes",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:   fmt.Sprintf("%s/client-test:latest", u.Host),
				NoVerify: true,
			},
			prepFunc: prepTestArtifact,
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/ArtifactIndex",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:   fmt.Sprintf("%s/client-test:latest", u.Host),
				NoVerify: true,
			},
			prepFunc: prepAggregate,
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/OneMatchingAnnotation",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "testdata/configs/match.yaml",
				NoVerify:       true,
			},
			prepFunc: prepTestArtifact,
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return err == nil
			},
		},
		{
			name: "Success/NoMatchingAnnotation",
			opts: &PullOptions{
				Common: &options.Common{
					IOStreams: genericclioptions.IOStreams{
						Out:    os.Stdout,
						In:     os.Stdin,
						ErrOut: os.Stderr,
					},
					Logger: testlogr,
				},
				Remote: options.Remote{
					PlainHTTP: true,
				},
				Source:         fmt.Sprintf("%s/client-test:latest", u.Host),
				AttributeQuery: "testdata/configs/nomatch.yaml",
				NoVerify:       true,
			},
			prepFunc: prepTestArtifact,
			assertFunc: func(path string) bool {
				actual := filepath.Join(path, "hello.txt")
				_, err = os.Stat(actual)
				return errors.Is(err, os.ErrNotExist)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tmp := t.TempDir()
			c.opts.Output = tmp
			c.prepFunc(t, c.opts.Source, u.Host)

			cache := filepath.Join(t.TempDir(), "cache")
			require.NoError(t, os.MkdirAll(cache, 0750))
			c.opts.CacheDir = cache

			err := c.opts.Run(context.TODO())
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.True(t, c.assertFunc(tmp))
			}
		})
	}
}

// prepAggregate will push an aggregate into the
// registry for retrieval. Uses methods from oras-go.
func prepAggregate(t *testing.T, ref, host string) {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	ref1 := fmt.Sprintf("%s/test1:latest", host)

	desc1, err := publishFunc(fileName, ref1, fileContent, map[string]string{"test": "annotation"}, nil)
	require.NoError(t, err)
	fileName2 := "goodbye.txt"
	fileContent2 := []byte("Goodbye World!\n")
	ref2 := fmt.Sprintf("%s/test2:latest", host)
	desc2, err := publishFunc(fileName2, ref2, fileContent2, map[string]string{"test": "annotation"}, nil)

	memoryStore := memory.New()
	manifest, err := generateIndex(nil, desc1, desc2)
	require.NoError(t, err)

	manifestDesc, err := pushBlob(context.Background(), ocispec.MediaTypeImageIndex, manifest, memoryStore)
	require.NoError(t, err)

	err = memoryStore.Tag(context.Background(), manifestDesc, ref)
	require.NoError(t, err)

	repo, err := remote.NewRepository(ref)
	require.NoError(t, err)
	repo.PlainHTTP = true
	_, err = oras.Copy(context.TODO(), memoryStore, ref, repo, "", oras.DefaultCopyOptions)
	require.NoError(t, err)
}

// prepTestArtifact will push a hello.txt artifact into the
// registry for retrieval. Uses methods from oras-go.
func prepTestArtifact(t *testing.T, ref, host string) {
	fileName := "hello.txt"
	fileContent := []byte("Hello World!\n")
	_, err := publishFunc(fileName, ref, fileContent, map[string]string{"test": "annotation"}, nil)
	require.NoError(t, err)
}

func publishFunc(fileName, ref string, fileContent []byte, layerAnnotations, manifestAnnotations map[string]string) (ocispec.Descriptor, error) {
	ctx := context.TODO()
	// Push file(s) w custom mediatype to registry
	memoryStore := memory.New()
	layerDesc, err := pushBlob(ctx, ocispec.MediaTypeImageLayer, fileContent, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}
	if layerDesc.Annotations == nil {
		layerDesc.Annotations = map[string]string{}
	}
	layerDesc.Annotations = layerAnnotations
	layerDesc.Annotations[ocispec.AnnotationTitle] = fileName

	config := []byte("{}")
	configDesc, err := pushBlob(ctx, ocispec.MediaTypeImageConfig, config, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	manifest, err := generateManifest(configDesc, manifestAnnotations, layerDesc)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	manifestDesc, err := pushBlob(ctx, ocispec.MediaTypeImageManifest, manifest, memoryStore)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	err = memoryStore.Tag(ctx, manifestDesc, ref)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}

	repo, err := remote.NewRepository(ref)
	if err != nil {
		return ocispec.Descriptor{}, nil
	}
	repo.PlainHTTP = true
	return oras.Copy(context.TODO(), memoryStore, ref, repo, "", oras.DefaultCopyOptions)
}
