package tools

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"reflect"
)

func JsonEquals(json1, json2 string) (bool, error) {
	var data1 interface{}
	var data2 interface{}

	var err error
	err = json.Unmarshal([]byte(json1), &data1)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal json1")
	}
	err = json.Unmarshal([]byte(json2), &data2)
	if err != nil {
		return false, errors.Wrap(err, "failed to marshal json2")
	}

	return reflect.DeepEqual(data1, data2), nil
}

func Yaml2Json(yamlStr string) (jsonStr string, err error) {
	var body interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &body); err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal YAML")
	}
	body = convert(body)
	b, err := json.Marshal(body)
	return string(b), errors.Wrap(err, "Failed to marshal JSON")
}

func WrapJSONList(wrapperField, restJson string) (string, error) {
	var parsed interface{}
	if err := json.Unmarshal([]byte(restJson), &parsed); err != nil {
		return "", err
	}
	wrapped := map[string][]interface{}{
		wrapperField: {parsed},
	}
	converted, err := json.Marshal(wrapped)
	return string(converted), errors.Wrap(err, "failed to marshal JSON")
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
