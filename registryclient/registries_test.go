package registryclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindRegistry(t *testing.T) {
	type spec struct {
		name     string
		cfg      RegistryConfig
		inRef    string
		expError string
		expReg   Registry
	}
	cases := []spec{
		{
			name: "Success/OneMatch",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.example.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*.not.com",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "*.example.com",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/MultipleMatches",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.example.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "*.example.com",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/SubDomainWildcard",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "reg.example.*",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
					{
						Prefix: "*",
						Endpoint: Endpoint{
							SkipTLS: false,
						},
					},
				},
			},
			inRef: "reg.example.com",
			expReg: Registry{
				Prefix: "reg.example.*",
				Endpoint: Endpoint{
					SkipTLS: true,
				},
			},
		},
		{
			name: "Success/NotMatch",
			cfg: RegistryConfig{
				Registries: []Registry{
					{
						Prefix: "*.not.com",
						Endpoint: Endpoint{
							SkipTLS: true,
						},
					},
				},
			},
			inRef:  "reg.example.com",
			expReg: Registry{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg, err := FindRegistry(c.cfg, c.inRef)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				if c.expReg.Prefix == "" {
					require.Equal(t, (*Registry)(nil), reg)
				} else {
					require.Equal(t, c.expReg, *reg)
				}
			}
		})
	}
}

func TestEndpoint_RewriteReference(t *testing.T) {
	type spec struct {
		name     string
		expError string
		endpoint Endpoint
		inRef    string
		expRef   string
	}

	cases := []spec{
		{
			name: "Success/MatchingPrefix",
			endpoint: Endpoint{
				Location: "alt.example.com",
			},
			inRef:  "reg.example.com/test:latest",
			expRef: "alt.example.com/test:latest",
		},
		{
			name: "Success/EmptyLocation",
			endpoint: Endpoint{
				Location: "",
			},
			inRef:  "reg.example.com/test:latest",
			expRef: "reg.example.com/test:latest",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ref, err := c.endpoint.RewriteReference(c.inRef)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expRef, ref)
			}
		})
	}
}

func TestRegistry_PullSourceFromReference(t *testing.T) {
	type spec struct {
		name       string
		expError   string
		registry   Registry
		inRef      string
		expSources []PullSource
	}
	cases := []spec{
		{
			name: "Success/NoMirrors",
			registry: Registry{
				Prefix: "reg.example.com",
				Endpoint: Endpoint{
					SkipTLS: false,
				},
			},
			inRef: "reg.example.com/test:latest",
		},
		{
			name: "Success/OneMirror",
			registry: Registry{
				Prefix: "reg.example.com",
				Endpoint: Endpoint{
					SkipTLS: false,
				},
				Mirrors: []Endpoint{
					{
						SkipTLS:  true,
						Location: "alt.registry.com",
					},
				},
			},
			inRef: "reg.example.com/test:latest",
			expSources: []PullSource{
				{
					Reference: "alt.registry.com/test:latest",
					Endpoint: Endpoint{
						SkipTLS:   true,
						PlainHTTP: false,
						Location:  "alt.registry.com",
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sources, err := c.registry.PullSourceFromReference(c.inRef)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expSources, sources)
			}
		})
	}
}
