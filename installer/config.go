package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration for installation
type Config struct {
	PrivateKey   string            `yaml:"private_key_path"`
	PreInstall   []string          `yaml:"pre_install_commands"`
	Installation []string          `yaml:"installation_commands"`
	PostInstall  []string          `yaml:"post_install_commands"`
	Environment  map[string]string `yaml:"environment_variables"`
}

// LoadConfig loads and validates the installation configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yml" && ext != ".yaml" {
		return nil, fmt.Errorf("configuration file must have .yml or .yaml extension")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.PrivateKey == "" {
		return fmt.Errorf("private_key_path is required in configuration")
	}
	return nil
}
