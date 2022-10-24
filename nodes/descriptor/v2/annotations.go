package v2

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

// AnnotationsToAttributeSet converts annotations from descriptors
// to an AttributeSet. This also performs annotation validation.
func AnnotationsToAttributeSet(annotations map[string]string, skip func(string) bool) (model.AttributeSet, error) {
	set := attributes.Attributes{}

	for key, value := range annotations {
		if skip != nil && skip(key) {
			continue
		}

		// Handle key collision. This should only occur if
		// an annotation is set and the key also exists in the UOR
		// specific attributes.
		if _, exists := set[key]; exists {
			continue
		}

		// Since annotations are in the form of map[string]string, we
		// can just assume it is a string attribute at this point. Incorporating
		// this into thr attribute set allows, users to pull by filename or reference name (cache).
		if key != descriptor.AnnotationUORAttributes {
			set[key] = attributes.NewString(key, value)
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			return set, err
		}
		for jKey, jVal := range data {
			attr, err := attributes.Reflect(jKey, jVal)
			if err != nil {
				return set, fmt.Errorf("annotation %q: error creating attribute: %w", key, err)
			}
			set[jKey] = attr
		}
	}
	return set, nil
}

// AnnotationsFromAttributeSet converts an AttributeSet to annotations. All annotation values
// are saved in a JSON valid syntax to allow for typing upon retrieval.
func AnnotationsFromAttributeSet(set model.AttributeSet) (map[string]string, error) {
	return map[string]string{descriptor.AnnotationUORAttributes: string(set.AsJSON())}, nil
}

// UpdateLayerDescriptors updates layers descriptor annotations with attributes from an AttributeSet. The key in the fileAttributes
// argument can be a regular expression or the name of a single file.
func UpdateLayerDescriptors(descs []ocispec.Descriptor, fileAttributes map[string]model.AttributeSet) ([]ocispec.Descriptor, error) {

	// Fail fast
	if len(fileAttributes) == 0 {
		return descs, nil
	}

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

	var updateDescs []ocispec.Descriptor
	for _, desc := range descs {

		var sets []model.AttributeSet
		if desc.Annotations == nil {
			desc.Annotations = map[string]string{}
		}

		filename, ok := desc.Annotations[ocispec.AnnotationTitle]
		if !ok {
			// skip any descriptor with no name attached
			continue
		}

		for file, set := range fileAttributes {
			nameSearch := regexpByFilename[file]
			if nameSearch.Match([]byte(filename)) {
				sets = append(sets, set)
			}
		}

		if len(sets) > 0 {
			mergedSet := mergeAttributes(sets)
			desc.Annotations[descriptor.AnnotationUORAttributes] = string(mergedSet.AsJSON())
		}

		updateDescs = append(updateDescs, desc)
	}
	return updateDescs, nil
}

func mergeAttributes(sets []model.AttributeSet) model.AttributeSet {
	newSet := attributes.Attributes{}

	if len(sets) == 0 {
		return newSet
	}

	if len(sets) == 1 {
		return sets[0]
	}

	for _, set := range sets {
		for key, value := range set.List() {
			newSet[key] = value
		}
	}

	return newSet
}
