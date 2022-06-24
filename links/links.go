package links

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/google/go-containerregistry/pkg/crane"
	digest "github.com/opencontainers/go-digest"
	"github.com/uor-framework/client/builder/api/v1alpha1"
	"github.com/uor-framework/client/registryclient/craneclient"
)

func FetchSchema(name string, i bool) ([]string, error) {
	// Retrieve the manifest of the collection indicated by name

	opts := craneclient.GetOpts(context.Background())
	manifest, err := crane.Manifest(name, opts...)
	if err != nil {
		return nil, err
	}
	// Unmarshal that manifest
	var man ocispec.Manifest
	json.Unmarshal([]byte(manifest), &man)

	// Write the collection schema and linkedSchema to a slice
	var schema []string
	schema = append(schema, man.Annotations["uor.schema"])
	schema = append(schema, man.Annotations["uor.linkedSchema"])

	return schema, err
}

// AddLinks creates links and returns a map of links and an error
func AddLinks(d v1alpha1.DataSetConfiguration, w string) (map[string]string, error) {
	// Create a link file with the name as the digest of the content
	links := make(map[string]string)
	for _, link := range d.LinkedCollections {
		dgst := digest.FromString(link)
		path := filepath.Join(w, dgst.String())
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		_, err = f.WriteString(link)
		if err != nil {
			return nil, err
		}
		f.Close()
		links[dgst.String()] = link
		fmt.Printf("links: %s", links)
	}
	// Return the mapping of collection addresses to their string digests
	return links, nil
}
