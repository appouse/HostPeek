package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Auth       AuthConfig       `yaml:"auth"`
	Collectors CollectorsConfig `yaml:"collectors"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Listen       string        `yaml:"listen"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	Enabled bool   `yaml:"enabled"`
	APIKey  string `yaml:"api_key"`
}

// CollectorsConfig controls which collectors are enabled.
type CollectorsConfig struct {
	CPU     bool `yaml:"cpu"`
	Memory  bool `yaml:"memory"`
	Disk    bool `yaml:"disk"`
	Network bool `yaml:"network"`
	OS      bool `yaml:"os"`
	Uptime  bool `yaml:"uptime"`
}

// Defaults returns a Config with sensible default values.
func Defaults() Config {
	return Config{
		Server: ServerConfig{
			Listen:       ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Auth: AuthConfig{
			Enabled: false,
		},
		Collectors: CollectorsConfig{
			CPU:     true,
			Memory:  true,
			Disk:    true,
			Network: true,
			OS:      true,
			Uptime:  true,
		},
	}
}

// Load reads the configuration from a YAML file.
// If the file does not exist, default values are used.
func Load(path string) (*Config, error) {
	cfg := Defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &cfg, nil // use defaults when no config file
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Apply defaults for zero values
	if cfg.Server.Listen == "" {
		cfg.Server.Listen = ":8080"
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 5 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}

	return &cfg, nil
}
