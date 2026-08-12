package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kmed/api/internal/httpx"
)

type ConfigStore struct {
	Paths *Paths
}

func NewConfigStore(appName string) *ConfigStore {
	paths, err := NewPaths(appName)
	if err != nil {
		log.Fatalf("failed to create paths: %v", err)
	}

	if err := os.MkdirAll(paths.DataDir, 0o700); err != nil {
		log.Fatalf("error creating the data directory: %v", err)
	}
	return &ConfigStore{Paths: paths}
}

type Config struct {
	DBPath        string `json:"db_path"         validate:"required"`
	AuthSecretKey string `json:"auth_secret_key" validate:"required"`
}

func (c *Config) Validate() error {
	return httpx.Validator.Struct(c)
}

func (cs ConfigStore) Load() (*Config, error) {
	data, err := os.ReadFile(cs.Paths.ConfigFile)
	// If file is not found, create default config.
	if errors.Is(err, os.ErrNotExist) {
		cfg, err := cs.DefaultConfig()
		if err != nil {
			return nil, err
		}
		if err := cs.Save(cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	var config *Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cs ConfigStore) Save(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("error validating config: %w", err)
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	if err := os.MkdirAll(cs.Paths.ConfigDir, 0o755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	if err := os.WriteFile(cs.Paths.ConfigFile, data, 0o644); err != nil {
		return fmt.Errorf(
			"error writing config to %s: %w",
			filepath.FromSlash(cs.Paths.ConfigFile),
			err,
		)
	}

	return nil
}

func (cs ConfigStore) DefaultConfig() (*Config, error) {
	secretKey, err := generateSecretKey()
	if err != nil {
		return nil, fmt.Errorf("error generating secret key: %w", err)
	}

	if err := os.MkdirAll(cs.Paths.DataDir, 0o700); err != nil {
		return nil, fmt.Errorf("error creating the data directory: %w", err)
	}

	return &Config{
		DBPath:        fmt.Sprintf("%s?mode=memory&cache=shared&_fk=1", cs.Paths.DBFile),
		AuthSecretKey: secretKey,
	}, nil
}

func generateSecretKey() (string, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		dir,
		"kmed",
		"internal-api",
		"config.json",
	), nil
}
