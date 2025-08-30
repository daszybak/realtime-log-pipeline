package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// NOTE `config` and its substructures must not be exported. Endpoints that
// require values from it should have those values injected through
// request-local context objects.
type config struct {
	PostgreSQL        struct {
		Addr     string `yaml:"addr"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		DB       string `yaml:"db"`
		MaxConns int    `yaml:"max_conns"`
	} `yaml:"postgresql"`
}

func readConfig(path string) (*config, error) {
	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %w", err)
	}

	config := &config{}
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse YAML: %w", err)
	}

	return config, nil
}

