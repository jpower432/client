package schema

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uor-framework/uor-client-go/attributes"
	"github.com/uor-framework/uor-client-go/model"
)

func TestSchema_Validate(t *testing.T) {
	type spec struct {
		name       string
		properties map[string]map[string]string
		doc        map[string]model.AttributeValue
		expRes     bool
		expError   string
	}

	cases := []spec{
		{
			name: "Success/ValidAttributes",
			properties: map[string]map[string]string{
				"size": {
					"type": "number",
				},
			},
			doc: map[string]model.AttributeValue{
				"size": attributes.NewFloat(1.0),
			},
			expRes: true,
		},
		{
			name: "Failure/IncompatibleType",
			properties: map[string]map[string]string{
				"size": {
					"type": "boolean",
				},
			},
			doc: map[string]model.AttributeValue{
				"size": attributes.NewFloat(1.0),
			},
			expRes:   false,
			expError: "size: invalid type. expected: boolean, given: integer",
		},
		{
			name: "Failure/MissingKey",
			properties: map[string]map[string]string{
				"size": {
					"type": "number",
				},
			},
			doc: map[string]model.AttributeValue{
				"name": attributes.NewString("test"),
			},
			expError: "(root): size is required",
			expRes:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			loader, err := fromProperties(c.properties)
			require.NoError(t, err)

			schema, err := New(loader)
			require.NoError(t, err)

			set := attributes.NewSet(c.doc)
			result, err := schema.Validate(set)
			if c.expError != "" {
				require.EqualError(t, err, c.expError)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expRes, result)
			}
		})
	}
}

func fromProperties(properties map[string]map[string]string) (Loader, error) {
	// Build an object in json from the provided types
	type jsonSchema struct {
		Type       string                       `json:"type"`
		Properties map[string]map[string]string `json:"properties"`
		Required   []string                     `json:"required"`
	}

	// Fill in properties and required keys. At this point
	// we consider all keys as required.
	var required []string
	for key := range properties {
		required = append(required, key)
	}

	// Make the required slice order deterministic
	sort.Slice(required, func(i, j int) bool {
		return required[i] < required[j]
	})

	tmp := jsonSchema{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		return Loader{}, err
	}
	return FromBytes(b)
}
