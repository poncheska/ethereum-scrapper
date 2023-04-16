package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	EthereumNodeURL string `yaml:"ethereumNodeURL"`

	Address string `yaml:"address"`

	MinBlock int64 `yaml:"minBlock"`
	MaxBlock int64 `yaml:"maxBlock"`
}

func ParseConfig(configPath string) (*Config, error) {
	fileBody, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(fileBody, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
