package orasclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/errdef"
)

const (
	// MediaTypeUnknownConfig is the default mediaType used when no
	// config media type is specified.
	MediaTypeUnknownConfig = "application/vnd.unknown.config.v1+json"
	// MediaTypeUnknownArtifact is the default artifactType used when no
	// artifact type is specified.
	MediaTypeUnknownArtifact = "application/vnd.unknown.artifact.v1"
)

// PackOptions and Pack are modified version of the upstream oras.Pack.
// The main difference is that the timestamp information is not recorded in the annotations.
// To ensure digest reproducibility for this effort, timestamp information will be collected and stored
// in predicate information in the artifact attestations.

// PackOptions contains parameters for Pack.
type PackOptions struct {
	// Subject is the subject of the manifest.
	Subject *ocispec.Descriptor
	// ManifestAnnotations is the annotation map of the manifest.
	ManifestAnnotations map[string]string
	// PackImageManifest controls whether to pack an image manifest or not.
	//   - If true, pack an image manifest; artifactType will be used as the
	// config descriptor mediaType of the image manifest.
	//   - If false, pack an artifact manifest.
	// Default: false.
	PackImageManifest bool
	// ConfigDescriptor is a pointer to the descriptor of the config blob.
	// If not nil, artifactType will be implied by the mediaType of the
	// specified ConfigDescriptor. This option is valid only when
	// PackImageManifest is true.
	ConfigDescriptor *ocispec.Descriptor
	// ConfigAnnotations is the annotation map of the config descriptor.
	// This option is valid only when PackImageManifest is true
	// and ConfigDescriptor is nil.
	ConfigAnnotations map[string]string
}

func Pack(ctx context.Context, pusher content.Pusher, artifactType string, blobs []ocispec.Descriptor, opts PackOptions) (ocispec.Descriptor, error) {
	if opts.PackImageManifest {
		return packImage(ctx, pusher, artifactType, blobs, opts)
	}
	return packArtifact(ctx, pusher, artifactType, blobs, opts)
}

// packArtifact packs the given blobs, generates an artifact manifest for the
// pack, and pushes it to a content storage.
// If succeeded, returns a descriptor of the manifest.
func packArtifact(ctx context.Context, pusher content.Pusher, artifactType string, blobs []ocispec.Descriptor, opts PackOptions) (ocispec.Descriptor, error) {
	if artifactType == "" {
		artifactType = MediaTypeUnknownArtifact
	}

	manifest := ocispec.Artifact{
		MediaType:    ocispec.MediaTypeArtifactManifest,
		ArtifactType: artifactType,
		Blobs:        blobs,
		Subject:      opts.Subject,
		Annotations:  opts.ManifestAnnotations,
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	manifestDesc := content.NewDescriptorFromBytes(ocispec.MediaTypeArtifactManifest, manifestJSON)
	// populate ArtifactType and Annotations of the manifest into manifestDesc
	manifestDesc.ArtifactType = manifest.ArtifactType
	manifestDesc.Annotations = manifest.Annotations

	// push manifest
	if err := pusher.Push(ctx, manifestDesc, bytes.NewReader(manifestJSON)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push manifest: %w", err)
	}

	return manifestDesc, nil
}

// packImage packs the given blobs, generates an image manifest for the pack,
// and pushes it to a content storage. artifactType will be used as the config
// descriptor mediaType of the image manifest.
// If succeeded, returns a descriptor of the manifest.
func packImage(ctx context.Context, pusher content.Pusher, configMediaType string, layers []ocispec.Descriptor, opts PackOptions) (ocispec.Descriptor, error) {
	if configMediaType == "" {
		configMediaType = MediaTypeUnknownConfig
	}

	configDesc, err := handleConfig(ctx, pusher, configMediaType, opts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}

	if layers == nil {
		layers = []ocispec.Descriptor{} // make it an empty array to prevent potential server-side bugs
	}
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{
			SchemaVersion: 2, // historical value. does not pertain to OCI or docker version
		},
		Config:      configDesc,
		MediaType:   ocispec.MediaTypeImageManifest,
		Layers:      layers,
		Subject:     opts.Subject,
		Annotations: opts.ManifestAnnotations,
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	manifestDesc := content.NewDescriptorFromBytes(ocispec.MediaTypeImageManifest, manifestJSON)
	// populate ArtifactType and Annotations of the manifest into manifestDesc
	manifestDesc.ArtifactType = manifest.Config.MediaType
	manifestDesc.Annotations = manifest.Annotations

	// push manifest
	if err := pusher.Push(ctx, manifestDesc, bytes.NewReader(manifestJSON)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push manifest: %w", err)
	}

	return manifestDesc, nil
}

// PackCollectionOptions contains parameters for PackCollection.
type PackCollectionOptions struct {
	// Links is the descriptor for linked artifacts.
	// This is only valid for ManifestType TypeCollection.
	Links []uorspec.Descriptor
	// ManifestAnnotations is the annotation map of the manifest.
	ManifestAttributes map[string]json.RawMessage
}

func PackCollection(ctx context.Context, pusher content.Pusher, artifactType string, blobs []uorspec.Descriptor, opts PackCollectionOptions) (ocispec.Descriptor, error) {
	if artifactType == "" {
		artifactType = MediaTypeUnknownArtifact
	}

	if blobs == nil {
		blobs = []uorspec.Descriptor{} // make it an empty array to prevent potential server-side bugs
	}

	manifest := uorspec.Manifest{
		MediaType:    uorspec.MediaTypeCollectionManifest,
		ArtifactType: artifactType,
		Links:        opts.Links,
		Blobs:        blobs,
		Attributes:   opts.ManifestAttributes,
	}
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	manifestDesc := content.NewDescriptorFromBytes(uorspec.MediaTypeCollectionManifest, manifestJSON)
	// populate ArtifactType and Annotations of the manifest into manifestDesc
	manifestDesc.ArtifactType = manifest.ArtifactType
	manifestDesc.Annotations = manifest.Annotations

	// push manifest
	if err := pusher.Push(ctx, manifestDesc, bytes.NewReader(manifestJSON)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
		return ocispec.Descriptor{}, fmt.Errorf("failed to push manifest: %w", err)
	}

	return manifestDesc, nil
}

func handleConfig(ctx context.Context, pusher content.Pusher, configMediaType string, opts PackOptions) (ocispec.Descriptor, error) {
	var configDesc ocispec.Descriptor
	if opts.ConfigDescriptor != nil {
		configDesc = *opts.ConfigDescriptor
	} else {
		// Use an empty JSON object here, because some registries may not accept
		// empty config blob.
		// As of September 2022, GAR is known to return 400 on empty blob upload.
		// See https://github.com/oras-project/oras-go/issues/294 for details.
		configBytes := []byte("{}")
		configDesc = content.NewDescriptorFromBytes(configMediaType, configBytes)
		configDesc.Annotations = opts.ConfigAnnotations
		// push config
		if err := pusher.Push(ctx, configDesc, bytes.NewReader(configBytes)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
			return ocispec.Descriptor{}, fmt.Errorf("failed to push config: %w", err)
		}
	}
	return configDesc, nil
}
