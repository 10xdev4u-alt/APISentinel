package config

import (
	"os"
	"strconv"

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
	Target  string   `yaml:"target"`
	Targets []string `yaml:"targets"`
}

type SecurityConfig struct {
	EnableXSS  bool   `yaml:"enable_xss"`
	EnableSQLi bool   `yaml:"enable_sqli"`
	EnableDLP  bool   `yaml:"enable_dlp"`
	DLPAction  string `yaml:"dlp_action"`
}

// LoadConfig reads the YAML configuration file and applies environment overrides.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Apply Defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.RateLimit == 0 {
		cfg.Server.RateLimit = 10
	}
	if cfg.Server.AuditLog == "" {
		cfg.Server.AuditLog = "audit.log"
	}

	// Apply Environment Overrides
	cfg.applyEnvOverrides()

	return &cfg, nil
}

func (c *Config) applyEnvOverrides() {
	if val := os.Getenv("SENTINEL_PORT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			c.Server.Port = i
		}
	}
	if val := os.Getenv("SENTINEL_ADMIN_KEY"); val != "" {
		c.Server.AdminKey = val
	}
	if val := os.Getenv("SENTINEL_RATE_LIMIT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			c.Server.RateLimit = i
		}
	}
	if val := os.Getenv("SENTINEL_AUDIT_LOG"); val != "" {
		c.Server.AuditLog = val
	}
	if val := os.Getenv("SENTINEL_DLP_ACTION"); val != "" {
		c.Security.DLPAction = val
	}
}
