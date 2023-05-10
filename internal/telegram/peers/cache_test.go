package peers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.octolab.org/toolset/indexit/internal/telegram/uid"
)

func TestCacheRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "peers.json")
	cache := New()
	cache.Put(Entry{Kind: uid.KindChannel, ID: 42, AccessHash: 99, Username: "telegram"})
	require.NoError(t, cache.Save(path))

	loaded, err := Load(path)
	require.NoError(t, err)
	entry, ok := loaded.Get(uid.KindChannel, 42)
	require.True(t, ok)
	assert.Equal(t, int64(99), entry.AccessHash)
	assert.Equal(t, "telegram", entry.Username)
}
