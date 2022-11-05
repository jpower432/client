package v2

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
	"github.com/uor-framework/uor-client-go/schema"
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
		if key != uorspec.AnnotationUORAttributes {
			set[key] = attributes.NewString(key, value)
			continue
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
			return set, err
		}
		for jsonKey, jsonVal := range jsonData {
			attr, err := attributes.Reflect(jsonKey, jsonVal)
			if err != nil {
				return set, fmt.Errorf("annotation %q: error creating attribute: %w", key, err)
			}
			set[jsonKey] = attr
		}
	}
	return set, nil
}

// AnnotationsFromAttributeSet converts an AttributeSet to annotations. All annotation values
// are saved in a JSON valid syntax to allow for typing upon retrieval.
func AnnotationsFromAttributeSet(set model.AttributeSet) (map[string]string, error) {
	attrJSON, err := set.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return map[string]string{uorspec.AnnotationUORAttributes: string(attrJSON)}, nil
}

// AnnotationsToAttributes OCI descriptor annotations to collection spec attributes if
// the AnnotationsUORAttributes key is found.
func AnnotationsToAttributes(annotations map[string]string) (map[string]json.RawMessage, error) {
	specAttributes := map[string]json.RawMessage{}
	extraAnnotations := map[string]string{}
	for key, value := range annotations {

		if key != uorspec.AnnotationUORAttributes {
			extraAnnotations[key] = value
			continue
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
			return specAttributes, err
		}

		for jsonKey, iVal := range jsonData {
			jsonVal, err := json.Marshal(iVal)
			if err != nil {
				return specAttributes, err
			}
			specAttributes[jsonKey] = jsonVal
		}

	}

	if len(extraAnnotations) != 0 {
		jsonValue, err := json.Marshal(extraAnnotations)
		if err != nil {
			return specAttributes, nil
		}
		specAttributes[schema.ConvertedSchemaID] = jsonValue
	}

	return specAttributes, nil
}

// AnnotationsFromAttributes converts collection spec attributes to collection annotations
//// for OCI descriptor compatibility.
func AnnotationsFromAttributes(attributes map[string]json.RawMessage) (map[string]string, error) {
	attrJSoN, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	return map[string]string{uorspec.AnnotationUORAttributes: string(attrJSoN)}, nil
}

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
