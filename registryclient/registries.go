package registryclient

import (
	"fmt"
	"regexp"
	"strings"

	"oras.land/oras-go/v2/errdef"
)

// This configuration is slightly modified and paired down version of the registries.conf.
// Source https://github.com/containers/image/blob/main/pkg/sysregistriesv2/system_registries_v2.go.
// More information on why this does not just use the `containers/system_registries_v2` library.
// While this library has a lot of overlapping functionality, it has more functionality than we
// need, and it makes sense to use the `containers` registry client which we are not. Search registries
// will eventually be a used in this library, but will be resolved and related to collection attributes
// and not short names.

// Endpoint describes a remote location of a registry.
type Endpoint struct {
	// The endpoint's remote location.
	Location string `json:"location"`
	// If true, certs verification will be skipped.
	SkipTLS bool `json:"skipTLS"`
	// If true, the client will use HTTP to
	// connect to the registry.
	PlainHTTP bool `json:"plainHTTP"`
}

// RewriteReference returns a reference for the endpoint given the original
// reference and registry prefix.
func (e Endpoint) RewriteReference(reference string) (string, error) {
	if e.Location == "" {
		return reference, nil
	}

	parts := strings.SplitN(reference, "/", 2)
	if len(parts) == 1 {
		return " ", fmt.Errorf("%w: missing repository", errdef.ErrInvalidReference)
	}
	path := parts[1]
	return fmt.Sprintf("%s/%s", e.Location, path), nil
}

// PullSource is a reference that is associated with a
// specific endpoint. This is used to generate references
// for registry mirrors and correlate them the mirror endpoint
type PullSource struct {
	Reference string
	Endpoint
}

// Registry represents a registry.
type Registry struct {
	// Prefix is used for endpoint matching.
	Prefix string `json:"prefix"`
	// A registry is an Endpoint too
	Endpoint `json:"endpoint"`
	// The registry mirrors
	Mirrors []Endpoint `json:"mirrors,omitempty"`
}

// PullSourceFromReference returns all pull source for the registry mirrors from
// a given reference.
func (r *Registry) PullSourceFromReference(ref string) ([]PullSource, error) {
	var sources []PullSource
	for _, mirror := range r.Mirrors {
		rewritten, err := mirror.RewriteReference(ref)
		if err != nil {
			return nil, err
		}
		sources = append(sources, PullSource{Endpoint: mirror, Reference: rewritten})
	}
	return sources, nil
}

// RegistryConfig is a configuration to configure multiple
// registry endpoints.
type RegistryConfig struct {
	Registries []Registry `json:"registries"`
	//AttributeSearchDomain []string
}

// FindRegistry returns the registry from the registry config that
// matches the reference.
func FindRegistry(registryConfig RegistryConfig, ref string) (*Registry, error) {
	reg := Registry{}
	prefixLen := 0

	for _, r := range registryConfig.Registries {
		prefixExp, err := regexp.Compile(validPrefix(r.Prefix))
		if err != nil {
			return nil, err
		}
		if prefixExp.MatchString(ref) {
			if len(r.Prefix) > prefixLen {
				reg = r
				prefixLen = len(r.Prefix)
			}
		}
	}
	if prefixLen != 0 {
		return &reg, nil
	}
	return nil, nil
}

// validPrefix will check the registry prefix value
// and return a valid regex.
func validPrefix(regPrefix string) string {
	if strings.HasPrefix(regPrefix, "*") {
		return strings.Replace(regPrefix, "*", ".*", -1)
	}
	return regPrefix
}
