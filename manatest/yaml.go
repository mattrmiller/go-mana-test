// Package manatest provides the inner workings of go-mana-test.
package manatest

// Imports
import (
	"encoding/json"
)

// Converts yaml to json
func ConvertYamlToJson(i interface{}) (string, error) {

	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			tmp, err := ConvertYamlToJson(v)
			if err != nil {
				return "", err
			}
			m2[k.(string)] = tmp
		}

		ret, err := json.Marshal(m2)
		if err != nil {
			return "", err
		}
		return string(ret), nil
	case []interface{}:
		for i, v := range x {
			tmp, err := ConvertYamlToJson(v)
			if err != nil {
				return "", err
			}
			x[i] = tmp
		}
	}

	return i.(string), nil
}
