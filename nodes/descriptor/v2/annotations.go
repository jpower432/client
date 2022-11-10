package v2

import (
	"fmt"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

// UpdateDescriptors updates descriptors with attributes from an AttributeSet. The key in the fileAttributes
// argument can be a regular expression or the name of a single file. The descriptor and node properties are updated
// by this method and the updated descriptors are returned.
func UpdateDescriptors(nodes []Node, schemaID string, fileAttributes map[string]model.AttributeSet) ([]ocispec.Descriptor, error) {
	var updateDescs []ocispec.Descriptor

	// Process each key into a regular expression and store it.
	regexpByFilename := map[string]*regexp.Regexp{}
	for file := range fileAttributes {
		// If the config has a grouping declared, make a valid regex.
		var expression string
		if strings.Contains(file, "*") && !strings.Contains(file, ".*") {
			expression = strings.Replace(file, "*", ".*", -1)
		} else {
			expression = strings.Replace(file, file, "^"+file+"$", -1)
		}

		nameSearch, err := regexp.Compile(expression)
		if err != nil {
			return nil, err
		}
		regexpByFilename[file] = nameSearch
	}

	for _, node := range nodes {

		var sets []model.AttributeSet

		if node.Location == "" {
			continue
		}

		for file, set := range fileAttributes {
			nameSearch := regexpByFilename[file]
			if nameSearch.Match([]byte(node.Location)) {
				sets = append(sets, set)
			}
		}

		desc := node.Descriptor()
		if len(sets) > 0 {
			merged, err := attributes.Merge(sets...)
			if err != nil {
				return nil, err
			}
			if err := node.Properties.Merge(map[string]model.AttributeSet{schemaID: merged}); err != nil {
				return nil, fmt.Errorf("file %s: %w", node.Location, err)
			}
		}
		mergedJSON, err := node.Properties.MarshalJSON()
		if err != nil {
			return nil, err
		}
		desc.Annotations[uorspec.AnnotationUORAttributes] = string(mergedJSON)

		updateDescs = append(updateDescs, desc)
	}
	return updateDescs, nil
}
