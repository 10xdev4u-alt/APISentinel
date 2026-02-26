package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the entire application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Routes   []RouteConfig  `yaml:"routes"`
	Security SecurityConfig `yaml:"security"`
}

type ServerConfig struct {
	Port      int    `yaml:"port"`
	AdminKey  string `yaml:"admin_key"`
	RateLimit int    `yaml:"rate_limit"`
	AuditLog  string `yaml:"audit_log"`
}

type RouteConfig struct {
	Path    string   `yaml:"path"`
	Target  string   `yaml:"target"`  // Singular (backward compatibility)
	Targets []string `yaml:"targets"` // Plural (for load balancing)
}

type SecurityConfig struct {
	EnableXSS  bool   `yaml:"enable_xss"`
	EnableSQLi bool   `yaml:"enable_sqli"`
	EnableDLP  bool   `yaml:"enable_dlp"`
	DLPAction  string `yaml:"dlp_action"` // "block" or "mask"
}

// LoadConfig reads the YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set some defaults if missing
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.RateLimit == 0 {
		cfg.Server.RateLimit = 10
	}
	if cfg.Server.AuditLog == "" {
		cfg.Server.AuditLog = "audit.log"
	}

	return &cfg, nil
}
