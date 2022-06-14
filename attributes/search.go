package attributes

import (
	"github.com/uor-framework/client/model"
	"github.com/uor-framework/client/model/traversal"
)

// FindFirstNode will search the given tree and find the first node that
// meets the attribute criteria.
// TODO(jpower432): Use greedy approach with a priority queue here instead of using the traversal
// package possibly.
func FindFirstNode(m AttributeMatcher, t model.Tree) (model.Node, error) {
	var result model.Node
	err := traversal.WalkWithStop(t, m, func(t traversal.Tracker, n model.Node) error {
		return nil
	})
	return result, err
}

// FindAllNodes will search the given tree and find all nodes that
// meet the attribute criteria.
func FindAllNodes(m AttributeMatcher, t model.Tree) ([]model.Node, error) {
	var result []model.Node
	err := traversal.Walk(t, func(t traversal.Tracker, n model.Node) error {
		return nil
	})
	return result, err
}
