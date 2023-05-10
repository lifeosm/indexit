package mapper

import (
	"testing"

	gotdpeer "github.com/gotd/td/telegram/message/peer"
	"github.com/gotd/td/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.octolab.org/toolset/indexit/internal/telegram/peers"
	"go.octolab.org/toolset/indexit/internal/telegram/uid"
)

func TestCacheEntities(t *testing.T) {
	cache := peers.New()
	entities := gotdpeer.NewEntities(
		map[int64]*tg.User{1: {ID: 1, AccessHash: 10, Username: "alice", FirstName: "Alice"}},
		map[int64]*tg.Chat{2: {ID: 2, Title: "Chat"}},
		map[int64]*tg.Channel{3: {ID: 3, AccessHash: 30, Username: "chan", Title: "Channel"}},
	)

	CacheEntities(cache, entities)

	user, ok := cache.Get(uid.KindUser, 1)
	require.True(t, ok)
	assert.Equal(t, int64(10), user.AccessHash)
	channel, ok := cache.Get(uid.KindChannel, 3)
	require.True(t, ok)
	assert.Equal(t, "chan", channel.Username)
}

func TestDialogMapsChannel(t *testing.T) {
	entities := gotdpeer.NewEntities(nil, nil, map[int64]*tg.Channel{
		3: {ID: 3, AccessHash: 30, Username: "chan", Title: "Channel", Megagroup: true, Forum: true},
	})

	record, ok := Dialog(&tg.Dialog{Peer: &tg.PeerChannel{ChannelID: 3}, UnreadCount: 2, Pinned: true}, entities, nil)
	require.True(t, ok)
	assert.Equal(t, "dialog", record.Kind)
	assert.Equal(t, "channel:3", record.UID)
	assert.Equal(t, "supergroup", record.PeerType)
	assert.True(t, record.IsForum)
}
