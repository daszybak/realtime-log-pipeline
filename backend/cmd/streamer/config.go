package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	PSQL struct {
		Addr     string `yaml:"addr"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		DB       string `yaml:"db"`
		MaxConns int    `yaml:"max_conns"`
	} `yaml:"psql"`
	RabbitMQ struct {
		URL string `yaml:"url"`
	} `yaml:"rabbitmq"`
	Binance struct {
		BaseURL string   `yaml:"base_url"`
		Symbols []string `yaml:"symbols"`
	} `yaml:"binance"`
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
