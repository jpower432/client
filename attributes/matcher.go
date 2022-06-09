package attributes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

// AttributeMatcher contains configuration data for searching for a node by attribute.
type AttributeMatcher struct {
	attributes map[string]string
}

var _ model.Matcher = &AttributeMatcher{}

// NewAttributeMatcher returns a new assembled matcher.
func NewAttributeMatcher(attributes map[string]string) model.Matcher {
	matcher := &AttributeMatcher{
		attributes: attributes,
	}
	return matcher
}

// String list all attributes in the Matcher in a string format.
func (m AttributeMatcher) String() string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(m.attributes))
	for k := range m.attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		line := fmt.Sprintf("%s=%s,", key, m.attributes[key])
		out.WriteString(line)
	}
	return strings.TrimSuffix(out.String(), ",")
}

// Matches determines whether a node has all required attributes.
func (m AttributeMatcher) Matches(n model.Node) bool {
	for key, value := range m.attributes {
		nodeVal, exist := n.Attributes()[key]
		if !exist {
			return false
		}
		if nodeVal != value {
			return false
		}
	}
	return true
}
