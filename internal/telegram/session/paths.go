package session

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const appDir = "indexit"

type Paths struct {
	ConfigDir string
	CacheDir  string
	Session   string
	Peers     string
}

func DefaultPaths() (Paths, error) {
	config, err := xdgDir("XDG_CONFIG_HOME", ".config")
	if err != nil {
		return Paths{}, err
	}
	cache, err := xdgDir("XDG_CACHE_HOME", ".cache")
	if err != nil {
		return Paths{}, err
	}

	configDir := filepath.Join(config, appDir, "telegram")
	cacheDir := filepath.Join(cache, appDir, "telegram")
	return Paths{
		ConfigDir: configDir,
		CacheDir:  cacheDir,
		Session:   filepath.Join(configDir, "session.json"),
		Peers:     filepath.Join(cacheDir, "peers.json"),
	}, nil
}

func EnsureSessionPath(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create session directory: %w", err)
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat session file: %w", err)
	}
	if info.Mode().Perm() != 0600 {
		return fmt.Errorf("session file %s has permissions %s, want 0600", path, info.Mode().Perm())
	}
	return nil
}

func CredentialsFromEnv() (int, string, error) {
	idText := os.Getenv("TELEGRAM_API_ID")
	hash := os.Getenv("TELEGRAM_API_HASH")
	if idText == "" {
		return 0, "", fmt.Errorf("TELEGRAM_API_ID is not set")
	}
	if hash == "" {
		return 0, "", fmt.Errorf("TELEGRAM_API_HASH is not set")
	}
	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, "", fmt.Errorf("parse TELEGRAM_API_ID: %w", err)
	}
	return id, hash, nil
}

func xdgDir(env string, fallback string) (string, error) {
	if value := os.Getenv(env); value != "" {
		return value, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, fallback), nil
}
