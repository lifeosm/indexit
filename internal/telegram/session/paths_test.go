package session

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPaths(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "config"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(t.TempDir(), "cache"))

	paths, err := DefaultPaths()
	require.NoError(t, err)
	assert.Contains(t, paths.Session, filepath.Join("indexit", "telegram", "session.json"))
	assert.Contains(t, paths.Peers, filepath.Join("indexit", "telegram", "peers.json"))
}

func TestEnsureSessionPathPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "session.json")
	require.NoError(t, os.WriteFile(path, []byte("{}"), 0600))
	require.NoError(t, EnsureSessionPath(path))
}
