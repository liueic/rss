package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Feed struct {
	ID                     string `yaml:"id"`
	Name                   string `yaml:"name"`
	URL                    string `yaml:"url"`
	Notify                 bool   `yaml:"notify"`
	DedupeKey              string `yaml:"dedupe_key"`
	Aggregate              bool   `yaml:"aggregate"`
	AggregateWindowMinutes int    `yaml:"aggregate_window_minutes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
