package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Paths struct {
	ConfigDir string
	DataDir   string
	CacheDir  string

	ConfigFile string
	DBFile     string
}

func NewPaths(appName string) (*Paths, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dataDir, err := userDataDir()
	if err != nil {
		return nil, err
	}

	configDir = filepath.Join(configDir, appName)
	dataDir = filepath.Join(dataDir, appName)
	cacheDir = filepath.Join(cacheDir, appName)

	return &Paths{
		ConfigDir: configDir,
		DataDir:   dataDir,
		CacheDir:  cacheDir,

		ConfigFile: filepath.Join(configDir, "config.json"),
		DBFile:     filepath.Join(dataDir, "app.db"),
	}, nil
}

func userDataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		dir := os.Getenv("LOCALAPPDATA")
		if dir != "" {
			return dir, nil
		}

		return os.UserHomeDir()

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(
			home,
			"Library",
			"Application Support",
		), nil

	default:
		// Linux / Unix: respect XDG_DATA_HOME.
		if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
			return dir, nil
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, ".local", "share"), nil
	}
}
