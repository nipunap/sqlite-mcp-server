package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Database struct {
		RegistryPath string `json:"registry_path"`
		DataDir      string `json:"data_dir"`
	} `json:"database"`
	Auth struct {
		Secret      string `json:"secret"`
		TokenExpiry int    `json:"token_expiry"` // in hours
	} `json:"auth"`
}

var DefaultConfig = Config{
	Server: struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{
		Host: "localhost",
		Port: 8080,
	},
	Database: struct {
		RegistryPath string `json:"registry_path"`
		DataDir      string `json:"data_dir"`
	}{
		RegistryPath: "data/registry.db",
		DataDir:      "data/databases",
	},
	Auth: struct {
		Secret      string `json:"secret"`
		TokenExpiry int    `json:"token_expiry"`
	}{
		Secret:      "change-me-in-production",
		TokenExpiry: 24,
	},
}

func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig

	if path == "" {
		return &config, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(config.Database.RegistryPath), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(config.Database.DataDir, 0755); err != nil {
		return nil, err
	}

	return &config, nil
}
