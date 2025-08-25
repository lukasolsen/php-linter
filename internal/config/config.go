package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Extends  string          `json:"extends"`
	Paths    []string        `json:"paths"`
	Excludes []string        `json:"excludes"`
	Stubs    []string        `json:"stubs"`
	Rules    map[string]bool `json:"rules"`

	PHPVersion string `json:"php_version,omitempty"`
}

func New() *Config {
	cfg, err := loadAndMergeConfig("config.json")
	if err != nil {
		return nil
	}

	phpVersion := cfg.PHPVersion
	if phpVersion == "" {
		phpVersion = "8.0"
	}

	return &Config{
		Extends:  cfg.Extends,
		Paths:    cfg.Paths,
		Excludes: cfg.Excludes,
		Stubs:    cfg.Stubs,
		Rules:    cfg.Rules,
		PHPVersion: phpVersion,
	}
}

func (cfg *Config) Validate() error {
	if cfg.Extends == "" {
		return fmt.Errorf("extends field is required")
	}
	if len(cfg.Paths) == 0 {
		return fmt.Errorf("at least one path is required")
	}
	return nil
}


func loadAndMergeConfig(path string) (*Config, error) {
	userBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var userCfg Config
	if err := json.Unmarshal(userBytes, &userCfg); err != nil {
		return nil, err
	}

	// If no extends, just return the user's config
	if userCfg.Extends == "" {
		return &userCfg, nil
	}

	// Not implemented yet.
	return &userCfg, nil
}