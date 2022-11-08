package attributequeries

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"oras.land/oras-go/v2/registry/remote/auth"

	"github.com/uor-framework/uor-client-go/util/testutils"
)

// TODO(jpower432): Mock the attributes API for testing.
func TestQueryForAttributes(t *testing.T) {
	manifest := []byte("hello world")
	server := httptest.NewServer(testutils.NewRegistry(t, nil, [][]byte{manifest}))
	t.Cleanup(server.Close)
	u, err := url.Parse(server.URL)
	require.NoError(t, err)

	client := auth.DefaultClient

	results, err := QueryForAttributes(context.Background(), u.Host, []byte(`{"test": "me}`), client, true)
	require.NoError(t, err)

	var index ocispec.Index
	err = json.Unmarshal(results, &index)
}
