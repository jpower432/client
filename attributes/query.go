package attributes

import (
	"fmt"
	"strings"

	"github.com/uor-framework/client/model"
)

func (q *AttributeQuery) List() string {
	out := new(strings.Builder)
	for key, value := range q.attributes {
		line := fmt.Sprintf("%s:%s\n", key, value)
		out.WriteString(line)
	}
	return out.String()
}

func (q *AttributeQuery) Do() (model.Node, error) {
	// Not implemented
	return nil, nil
}
