package components

import (
	"crypto"
	"fmt"

	"github.com/anchore/syft/cmd/syft/cli/eventloop"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"

	clientapi "github.com/uor-framework/uor-client-go/api/client/v1alpha1"
)

const (
	scope             = "squash"
	skipAboveFileSize = 1048576
)

var digests = []string{"sha256"}

func getTasks(config clientapi.DataSetConfiguration) ([]eventloop.Task, error) {
	var tasks []eventloop.Task

	generators := []func(config clientapi.DataSetConfiguration) (eventloop.Task, error){
		generateCatalogPackagesTask,
		generateCatalogFileMetadataTask,
		generateCatalogFileDigestsTask,
		generateCatalogSecretsTask,
		generateCatalogFileClassificationsTask,
		generateCatalogContentsTask,
	}

	for _, generator := range generators {
		task, err := generator(config)
		if err != nil {
			return nil, err
		}

		if task != nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func generateCatalogPackagesTask(config clientapi.DataSetConfiguration) (eventloop.Task, error) {
	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		packageCatalog, relationships, theDistro, err := syft.CatalogPackages(src, toCatalogerConfig(config))
		if err != nil {
			return nil, err
		}

		results.PackageCatalog = packageCatalog
		results.LinuxDistribution = theDistro

		return relationships, nil
	}

	return task, nil
}

func generateCatalogFileMetadataTask(_ clientapi.DataSetConfiguration) (eventloop.Task, error) {
	metadataCataloger := file.NewMetadataCataloger()

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(scope)
		if err != nil {
			return nil, err
		}

		result, err := metadataCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileMetadata = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogFileDigestsTask(_ clientapi.DataSetConfiguration) (eventloop.Task, error) {
	supportedHashAlgorithms := make(map[string]crypto.Hash)
	for _, h := range []crypto.Hash{
		crypto.MD5,
		crypto.SHA1,
		crypto.SHA256,
	} {
		supportedHashAlgorithms[file.DigestAlgorithmName(h)] = h
	}

	var hashes []crypto.Hash
	for _, hashStr := range digests {
		name := file.CleanDigestAlgorithmName(hashStr)
		hashObj, ok := supportedHashAlgorithms[name]
		if !ok {
			return nil, fmt.Errorf("unsupported hash algorithm: %s", hashStr)
		}
		hashes = append(hashes, hashObj)
	}

	digestsCataloger, err := file.NewDigestsCataloger(hashes)
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(scope)
		if err != nil {
			return nil, err
		}

		result, err := digestsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileDigests = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogSecretsTask(_ clientapi.DataSetConfiguration) (eventloop.Task, error) {

	patterns, err := file.GenerateSearchPatterns(file.DefaultSecretsPatterns, nil, nil)
	if err != nil {
		return nil, err
	}

	secretsCataloger, err := file.NewSecretsCataloger(patterns, false, skipAboveFileSize)
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(scope)
		if err != nil {
			return nil, err
		}

		result, err := secretsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.Secrets = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogFileClassificationsTask(_ clientapi.DataSetConfiguration) (eventloop.Task, error) {
	classifierCataloger, err := file.NewClassificationCataloger(file.DefaultClassifiers)
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(scope)
		if err != nil {
			return nil, err
		}

		result, err := classifierCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileClassifications = result
		return nil, nil
	}

	return task, nil
}

func generateCatalogContentsTask(_ clientapi.DataSetConfiguration) (eventloop.Task, error) {
	contentsCataloger, err := file.NewContentsCataloger(nil, skipAboveFileSize)
	if err != nil {
		return nil, err
	}

	task := func(results *sbom.Artifacts, src *source.Source) ([]artifact.Relationship, error) {
		resolver, err := src.FileResolver(scope)
		if err != nil {
			return nil, err
		}

		result, err := contentsCataloger.Catalog(resolver)
		if err != nil {
			return nil, err
		}
		results.FileContents = result
		return nil, nil
	}

	return task, nil
}

func toCatalogerConfig(_ clientapi.DataSetConfiguration) cataloger.Config {
	return cataloger.Config{
		Search: cataloger.SearchConfig{
			IncludeIndexedArchives:   true,
			IncludeUnindexedArchives: false,
			Scope:                    scope,
		},
	}
}
