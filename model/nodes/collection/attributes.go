package collection

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uor-framework/client/model"
)

// Attributes implements the model.Attributes interface
// using a multi-map storing a set of values.
// The current implementation would allow for aggregation of the attributes
// of child nodes to the parent nodes.
// TODO(jpower432): Research alternative data structures for storing these values.
// Since this will most likely limit the collection size
// there could we a more efficient data structure for storing multiple values
// for one key for data aggregation.
type Attributes map[string]map[string]struct{}

var _ model.Attributes = &Attributes{}

// Find returns all values stored for a specified key.
func (m Attributes) Find(key string) []string {
	valSet, exists := m[key]
	if !exists {
		return nil
	}
	var vals []string
	for val := range valSet {
		vals = append(vals, val)
	}
	return vals
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (m Attributes) Exists(key, value string) bool {
	vals, exists := m[key]
	if !exists {
		return false
	}
	_, valExists := vals[value]
	return valExists
}

// Strings returns a string representation of the
// attribute set.
func (m Attributes) String() string {
	out := new(strings.Builder)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		for val := range m[key] {
			line := fmt.Sprintf("%s=%s,", key, val)
			out.WriteString(line)
		}
	}
	return strings.TrimSuffix(out.String(), ",")
}
