package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
	Sleep int    `yaml:"sleep,omitempty"` // optional sleep time in seconds between actions
}

type Rule struct {
	Id         string   `yaml:"id,omitempty"`
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
