package craneclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

func CreateRT() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			// By default we wrap the transport in retries, so reduce the
			// default dial timeout to 5s to avoid 5x 30s of connection
			// timeouts when doing the "ping" on certain http registries.
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func GetOpts(ctx context.Context) []crane.Option {
	options := []crane.Option{
		crane.WithAuthFromKeychain(authn.DefaultKeychain),
		crane.WithContext(ctx),
		crane.WithTransport(CreateRT()),
	}

	options = append(options, crane.Insecure)
	return options
}
