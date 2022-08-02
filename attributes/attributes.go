package attributes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/uor-framework/uor-client-go/model"
)

// Attributes implements the model.Attributes interface.
type Attributes map[string]json.RawMessage

var _ model.Attributes = &Attributes{}

// Find returns all values stored for a specified key.
func (a Attributes) Find(key string) []string {
	val, exists := a[key]
	if !exists {
		return nil
	}
	return []string{string(val)}
}

// Exists returns whether a key,value pair exists in the
// attribute set.
func (a Attributes) Exists(key, value string) bool {
	val, exists := a[key]
	if !exists {
		return false
	}
	if value == string(val) {
		return true
	}
	return false
}

// Strings returns a JSON formatted string representation of the
// attribute set. If the values are not valid, an empty string is returned.
func (a Attributes) String() string {
	var message []json.RawMessage
	// TODO (jpower432): need to incorporate key so this becomes
	// a JSON formatted dictionary.
	for _, val := range a {
		message = append(message, val)
	}
	merged, err := combine(message)
	if err != nil {
		return err.Error()
	}
	return string(merged)
}

// List will list all key, value pairs for the attributes in a
// consumable format.
func (a Attributes) List() map[string][]string {
	list := make(map[string][]string, len(a))
	for key, val := range a {
		list[key] = append(list[key], string(val))
	}
	return list
}

// Len returns the length of the attribute set.
func (a Attributes) Len() int {
	return len(a)
}

// combine will combine json messages into one
// message.
func combine(docs []json.RawMessage) (json.RawMessage, error) {
	if len(docs) == 0 {
		return []byte{}, nil
	}
	prev := docs[0]
	var err error
	for i := 1; i < len(docs); i++ {
		prev, err = mergeBytes(prev, docs[i])
		if err != nil {
			return prev, err
		}

	}

	return prev, nil
}

func mergeValue(path []string, patch map[string]interface{}, key string, value interface{}) (interface{}, error) {
	patchValue, patchHasValue := patch[key]

	if !patchHasValue {
		return value, nil
	}

	_, patchValueIsObject := patchValue.(map[string]interface{})

	path = append(path, key)
	pathStr := strings.Join(path, ".")

	if _, ok := value.(map[string]interface{}); ok {
		if !patchValueIsObject {
			return value, fmt.Errorf("patch value must be object for key \"%v\"", pathStr)
		}

		return mergeObjects(value, patchValue, path)
	}

	if _, ok := value.([]interface{}); ok && patchValueIsObject {
		return mergeObjects(value, patchValue, path)
	}

	return patchValue, nil
}

func mergeObjects(data, patch interface{}, path []string) (interface{}, error) {
	var err error
	if patchObject, ok := patch.(map[string]interface{}); ok {
		if dataArray, ok := data.([]interface{}); ok {
			ret := make([]interface{}, len(dataArray))

			for i, val := range dataArray {
				ret[i], err = mergeValue(path, patchObject, strconv.Itoa(i), val)
				if err != nil {
					return nil, err
				}
			}

			return ret, nil
		} else if dataObject, ok := data.(map[string]interface{}); ok {
			ret := make(map[string]interface{})

			for k, v := range dataObject {
				ret[k], err = mergeValue(path, patchObject, k, v)
				if err != nil {
					return nil, err
				}
			}

			return ret, nil
		}
	}

	return data, nil
}

// merge merges patch document to data document
func merge(data, patch interface{}) (interface{}, error) {
	return mergeObjects(data, patch, nil)
}

// mergeBytes merges patch document buffer to data document buffer
func mergeBytes(dataBuff, patchBuff []byte) (mergedBuff []byte, err error) {
	var data, patch, merged interface{}

	err = unmarshalJSON(dataBuff, &data)
	if err != nil {
		err = fmt.Errorf("error in data JSON: %v", err)
		return
	}

	err = unmarshalJSON(patchBuff, &patch)
	if err != nil {
		err = fmt.Errorf("error in patch JSON: %v", err)
		return
	}

	merged, err = merge(data, patch)
	if err != nil {
		return nil, err
	}

	mergedBuff, err = json.Marshal(merged)
	if err != nil {
		err = fmt.Errorf("error writing merged JSON: %v", err)
	}

	return
}

func unmarshalJSON(buff []byte, data interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(buff))
	decoder.UseNumber()
	return decoder.Decode(data)
}
