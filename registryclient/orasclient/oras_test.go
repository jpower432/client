package orasclient

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/stretchr/testify/require"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/content/memory"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/attributes/matchers"
	"github.com/uor-framework/uor-client-go/ocimanifest"
	"github.com/uor-framework/uor-client-go/registryclient"
)

func TestAddFiles(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		ctx := context.TODO()
		expDigest := "sha256:2e30f6131ce2164ed5ef017845130727291417d60a1be6fad669bdc4473289cd"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		require.Len(t, desc, 1)
		require.Equal(t, expDigest, desc[0].Digest.String())
	})
}

func TestAddContent(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		ctx := context.TODO()
		expDigest := "sha256:cf80cd8aed482d5d1527d7dc72fceff84e6326592848447d2dc0b0e87dfc9a90"
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.AddContent(ctx, "", []byte("testing"), nil)
		require.NoError(t, err)
		require.Equal(t, expDigest, desc.Digest.String())
	})
}

func TestAddManifest(t *testing.T) {
	t.Run("Success/OneArtifact", func(t *testing.T) {
		ctx := context.TODO()
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		testdata := filepath.Join("testdata", "workspace", "fish.jpg")
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		desc, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
		require.NoError(t, err)
		mdesc, err := c.AddManifest(ctx, "localhost:5000/test:latest", configDesc, nil, desc...)
		require.NoError(t, err)
		require.Equal(t, expDigest, mdesc.Digest.String())
	})
}

func TestSave(t *testing.T) {
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	ref := fmt.Sprintf("%s/test:latest", u.Host)
	testdata := filepath.Join("testdata", "workspace", "fish.jpg")

	ctx := context.TODO()

	t.Run("Success/PushOneImage", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
		require.NoError(t, err)

		mdesc, err := c.AddManifest(ctx, ref, configDesc, nil, descs...)
		require.NoError(t, err)

		memStore := memory.New()
		desc, err := c.Save(ctx, ref, memStore)
		require.NoError(t, err)
		require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		require.Equal(t, expDigest, desc.Digest.String())
		require.NoError(t, c.Destroy())
		desc, err = memStore.Resolve(ctx, ref)
		require.NoError(t, err)
		require.Equal(t, expDigest, desc.Digest.String())
	})
}

func TestPushPull(t *testing.T) {
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	ref := fmt.Sprintf("%s/test:latest", u.Host)
	notExistTag := "latest"
	notExistRef := fmt.Sprintf("%s/notexist:%s", u.Host, notExistTag)
	images := []string{fmt.Sprintf("%s/test:latest", u.Host), fmt.Sprintf("%s/test2:latest", u.Host)}
	testdata := filepath.Join("testdata", "workspace", "fish.jpg")

	ctx := context.TODO()

	t.Run("Success/PushOneCollection", func(t *testing.T) {
		cache := memory.New()
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		c, err := NewClient(WithPlainHTTP(true), WithCache(cache))
		require.NoError(t, err)
		descs, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
		require.NoError(t, err)

		mdesc, err := c.AddManifest(ctx, ref, configDesc, nil, descs...)
		require.NoError(t, err)
		source, err := c.Store()
		require.NoError(t, err)

		desc, err := c.Push(context.TODO(), source, ref)
		require.NoError(t, err)
		require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		require.Equal(t, expDigest, desc.Digest.String())
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullOneCollection", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		root, descs, err := c.Pull(context.TODO(), ref, memory.New())
		require.NoError(t, err)
		require.Equal(t, expDigest, root.Digest.String())
		require.Len(t, descs, 4)
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/FilteredCollection", func(t *testing.T) {
		expDigest := ""
		matcher := matchers.PartialAttributeMatcher{"test": attributes.NewString("test", "fail")}
		c, err := NewClient(WithPlainHTTP(true), WithPullableAttributes(matcher))
		require.NoError(t, err)
		root, descs, err := c.Pull(context.TODO(), ref, memory.New())
		require.NoError(t, err)
		require.Equal(t, expDigest, root.Digest.String())
		require.Len(t, descs, 0)
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullWithCache", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		cache := memory.New()
		c, err := NewClient(WithPlainHTTP(true), WithCache(cache))
		require.NoError(t, err)
		root, descs, err := c.Pull(context.TODO(), ref, memory.New())
		require.NoError(t, err)
		require.Equal(t, expDigest, root.Digest.String())
		require.Len(t, descs, 3)
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PushMultipleCollections", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		descs, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
		require.NoError(t, err)

		source, err := c.Store()
		require.NoError(t, err)

		for _, ref := range images {
			mdesc, err := c.AddManifest(ctx, ref, configDesc, nil, descs...)
			require.NoError(t, err)
			desc, err := c.Push(context.TODO(), source, ref)
			require.NoError(t, err)
			require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		}
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullMultipleCollections", func(t *testing.T) {
		tmp := t.TempDir()
		destination := file.New(tmp)
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		for _, ref := range images {
			_, _, err := c.Pull(context.TODO(), ref, destination)
			require.NoError(t, err)
			_, err = os.Stat(filepath.Join(tmp, testdata))
			require.NoError(t, err)
		}
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PushWithRegistryConfig", func(t *testing.T) {
		cache := memory.New()
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		config := registryclient.RegistryConfig{
			Registries: []registryclient.Registry{
				{
					Prefix: u.Host,
					Endpoint: registryclient.Endpoint{
						PlainHTTP: true,
						Location:  u.Host,
					},
				},
			},
		}
		c, err := NewClient(WithRegistryConfig(config), WithCache(cache))
		require.NoError(t, err)
		descs, err := c.AddFiles(ctx, "", testdata)
		require.NoError(t, err)
		configDesc, err := c.AddContent(ctx, ocimanifest.UORConfigMediaType, []byte("{}"), nil)
		require.NoError(t, err)

		mdesc, err := c.AddManifest(ctx, ref, configDesc, nil, descs...)
		require.NoError(t, err)
		source, err := c.Store()
		require.NoError(t, err)

		desc, err := c.Push(context.TODO(), source, ref)
		require.NoError(t, err)
		require.Equal(t, mdesc.Digest.String(), desc.Digest.String())
		require.Equal(t, expDigest, desc.Digest.String())
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullWithRegistryConfig", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		config := registryclient.RegistryConfig{
			Registries: []registryclient.Registry{
				{
					Prefix: u.Host,
					Endpoint: registryclient.Endpoint{
						PlainHTTP: true,
						Location:  u.Host,
					},
				},
			},
		}
		c, err := NewClient(WithRegistryConfig(config))
		require.NoError(t, err)
		root, descs, err := c.Pull(context.TODO(), ref, memory.New())
		require.NoError(t, err)
		require.Equal(t, expDigest, root.Digest.String())
		require.Len(t, descs, 4)
		require.NoError(t, c.Destroy())
	})

	t.Run("Success/PullWithRegistryConfigNoMatch", func(t *testing.T) {
		expDigest := "sha256:98f36e12e9dbacfbb10b9d1f32a46641eb42de588e54cfd7e8627d950ae8140a"
		config := registryclient.RegistryConfig{
			Registries: []registryclient.Registry{
				{
					Prefix: "anotherhost",
					Endpoint: registryclient.Endpoint{
						PlainHTTP: false,
						Location:  "another-host",
					},
				},
			},
		}
		c, err := NewClient(WithRegistryConfig(config), WithPlainHTTP(true))
		require.NoError(t, err)
		root, descs, err := c.Pull(context.TODO(), ref, memory.New())
		require.NoError(t, err)
		require.Equal(t, expDigest, root.Digest.String())
		require.Len(t, descs, 4)
		require.NoError(t, c.Destroy())
	})

	t.Run("Failure/ImageDoesNotExist", func(t *testing.T) {
		c, err := NewClient(WithPlainHTTP(true))
		require.NoError(t, err)
		_, _, err = c.Pull(context.TODO(), notExistRef, memory.New())
		require.EqualError(t, err, fmt.Sprintf("%s: not found", notExistRef))
		require.NoError(t, c.Destroy())
	})
}
