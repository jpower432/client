package v3

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

func AttributesToAttributeSet(specAttributes map[string]json.RawMessage, skip func(string) bool) (model.AttributeSet, error) {
	set := attributes.Attributes{}

	for key, value := range specAttributes {
		if skip != nil && skip(key) {
			continue
		}

		// Handle key collision. This should only occur if
		// an annotation is set and the key also exists in the UOR
		// specific attributes.
		if _, exists := set[key]; exists {
			continue
		}

		// FIXME(jpower432): Is there a better way to do this
		var data map[string]interface{}
		if err := json.Unmarshal(value, &data); err != nil {
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

func AttributesFromAttributeSet(set model.AttributeSet) (map[string]json.RawMessage, error) {
	attributes := map[string]json.RawMessage{}
	for _, a := range set.List() {
		valueJSON, err := json.Marshal(a.AsAny())
		if err != nil {
			return nil, err
		}
		attributes[a.Key()] = valueJSON
	}
	return attributes, nil
}

func AttributesFromAnnotations(annotations map[string]string, skip func(string) bool) (map[string]json.RawMessage, error) {
	specAttributes := map[string]json.RawMessage{}

	value, found := annotations[descriptor.AnnotationUORAttributes]
	if !found {
		return specAttributes, nil
	}

	// TODO(jpower432): Custom unmarsharling

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return specAttributes, err
	}
	for key, val := range data {
		jVal, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		specAttributes[key] = jVal
	}

	return specAttributes, nil
}

func AttributesToAnnotations(attributes map[string]json.RawMessage) (map[string]string, error) {
	attrJSoN, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	return map[string]string{descriptor.AnnotationUORAttributes: string(attrJSoN)}, nil
}

// TODO(jpower432): Deduplicate

// UpdateLayerDescriptors updates layers descriptor annotations with attributes from an AttributeSet. The key in the fileAttributes
// argument can be a regular expression or the name of a single file.
func UpdateLayerDescriptors(descs []uorspec.Descriptor, fileAttributes map[string]model.AttributeSet) ([]uorspec.Descriptor, error) {

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

	var updateDescs []uorspec.Descriptor
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
			mergedAttributes, err := AttributesFromAttributeSet(mergedSet)
			if err != nil {
				return nil, err
			}
			desc.Attributes = mergedAttributes
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
