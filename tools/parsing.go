package tools

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type idHolder struct {
	Id string `yaml:"id" json:"id"`
}

func RouteIdFromYaml(route string) (string, error) {
	var holder idHolder

	if err := yaml.Unmarshal([]byte(route), &holder); err != nil {
		return "", errors.Wrap(err, "failed to parse YAML")
	}
	if holder.Id == "" {
		return "", errors.Errorf("empty route ID")
	}
	return holder.Id, nil
}
