package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Expression string   `yaml:"expression"`
	Actions    []string `yaml:"actions"` // e.g., close, approve, merge
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
