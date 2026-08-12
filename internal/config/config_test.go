package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigStore_Load(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "kmed-test-config-load")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err.Error())
	}
	defer os.RemoveAll(tempDir)

	testConfig := &Config{DBPath: "database", AuthSecretKey: "secret"}
	data, err := json.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err.Error())
	}

	configPath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		t.Fatalf("failed to write to %v: %v", filepath.FromSlash(configPath), err.Error())
	}

	t.Run("load from existing config.json", func(t *testing.T) {
		cs := NewConfigStore(configPath)
		config, err := cs.Load()
		assert.Nil(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, testConfig, config)
	})

	t.Run("load from non-existent config.json returns error", func(t *testing.T) {
		cs := NewConfigStore(filepath.Join(tempDir, "non-existent.json"))
		config, err := cs.Load()
		assert.Nil(t, config)
		assert.NotNil(t, err)
		var expectedErr *os.PathError
		assert.ErrorAs(t, err, &expectedErr)
	})

	t.Run("load from a directory returns error", func(t *testing.T) {
		cs := NewConfigStore(tempDir)
		config, err := cs.Load()
		assert.Nil(t, config)
		assert.ErrorContains(t, err, "not directory")
	})

	t.Run("loading from a config.json with invalid data", func(t *testing.T) {
		data, err := json.Marshal("something else")
		if err != nil {
			t.Fatalf("failed to marshal data: %v", err.Error())
		}
		configPath := filepath.Join(tempDir, "invalid-config.json")
		if err := os.WriteFile(configPath, data, 0o600); err != nil {
			t.Fatalf("failed to write to %v: %v", filepath.FromSlash(configPath), err.Error())
		}

		cs := NewConfigStore(configPath)
		config, err := cs.Load()
		assert.Nil(t, config)
		assert.NotNil(t, err)
	})
}

func TestConfig_Save(t *testing.T) {
	t.Run("saving valid config", func(t *testing.T) {
		// config := &Config{DatabaseURL: "database", AuthSecretKey: "secret"}
	})
}
