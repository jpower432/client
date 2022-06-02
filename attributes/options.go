package attributes

import (
	"fmt"

	"github.com/uor-framework/client/model"
)

type QueryOption func(q *AttributeQuery) error

// Query contains configuration data for an attribute query.
type AttributeQuery struct {
	attributes map[string]string
}

func (q *AttributeQuery) apply(options []QueryOption) error {
	for _, option := range options {
		if err := option(q); err != nil {
			return err
		}
	}
	return nil
}

// NewQuery returns a new assembled query
func NewQuery(options ...QueryOption) (model.Query, error) {
	query := &AttributeQuery{}
	if err := query.apply(options); err != nil {
		return nil, err
	}

	return query, nil
}

// WithAttributes adds attribute to query for node searching.
func WithAuthConfigs(attributes map[string]string) QueryOption {
	return func(q *AttributeQuery) error {
		for k, v := range attributes {
			if _, exists := q.attributes[k]; exists {
				return fmt.Errorf("duplicates attributes in query %s", k)
			}
			q.attributes[k] = v
		}
		return nil
	}
}


