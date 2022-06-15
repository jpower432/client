package attributes

import (
	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/model/traversal"
)

// FindAllPartialMatches will search the given tree and find all nodes that
// meet the attribute criteria, but may have additional attributes.
func FindAllPartialMatches(m PartialAttributeMatcher, t model.Tree) ([]model.Node, error) {
	var result []model.Node
	err := traversal.Walk(t, func(t traversal.Tracker, n model.Node) error {
		if m.Matches(n) {
			result = append(result, n)
		}
		return nil
	})
	return result, err
}

// FindAllExactMatches will search the given tree and find all nodes that
// meet the attribute criteria exactly.
func FindAllExactMatches(m ExactAttributeMatcher, t model.Tree) ([]model.Node, error) {
	var result []model.Node
	err := traversal.Walk(t, func(t traversal.Tracker, n model.Node) error {
		if m.Matches(n) {
			result = append(result, n)
		}
		return nil
	})
	return result, err
}
