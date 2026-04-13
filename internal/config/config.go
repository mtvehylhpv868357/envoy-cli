package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds global CLI configuration.
type Config struct {
	DefaultShell   string `json:"default_shell,omitempty"`
	ProfilesDir    string `json:"profiles_dir,omitempty"`
	SnapshotsDir   string `json:"snapshots_dir,omitempty"`
	AutoExport     bool   `json:"auto_export"`
	ColorOutput    bool   `json:"color_output"`
}

const configFileName = "config.json"

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		AutoExport:  false,
		ColorOutput: true,
	}
}

// configPath returns the path to the config file inside the given base dir.
func configPath(baseDir string) string {
	return filepath.Join(baseDir, configFileName)
}

// Load reads the config file from baseDir. If it does not exist, the default
// config is returned without error.
func Load(baseDir string) (*Config, error) {
	cfg := DefaultConfig()
	path := configPath(baseDir)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config to baseDir, creating the directory if necessary.
func Save(baseDir string, cfg *Config) error {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath(baseDir), data, 0o644)
}
