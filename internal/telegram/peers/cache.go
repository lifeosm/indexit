package peers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.octolab.org/toolset/indexit/internal/telegram/uid"
)

type Entry struct {
	Kind       uid.Kind  `json:"kind"`
	ID         int64     `json:"id"`
	AccessHash int64     `json:"access_hash,omitempty"`
	Username   string    `json:"username,omitempty"`
	Title      string    `json:"title,omitempty"`
	LastSeen   time.Time `json:"last_seen"`
}

type Cache struct {
	Entries map[string]Entry `json:"entries"`
}

func New() *Cache {
	return &Cache{Entries: map[string]Entry{}}
}

func Load(path string) (*Cache, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read peer cache: %w", err)
	}
	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("decode peer cache: %w", err)
	}
	if cache.Entries == nil {
		cache.Entries = map[string]Entry{}
	}
	return &cache, nil
}

func (c *Cache) Save(path string) error {
	if c == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create peer cache directory: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode peer cache: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write peer cache: %w", err)
	}
	return nil
}

func (c *Cache) Get(kind uid.Kind, id int64) (Entry, bool) {
	if c == nil {
		return Entry{}, false
	}
	entry, ok := c.Entries[key(kind, id)]
	return entry, ok
}

func (c *Cache) Len() int {
	if c == nil {
		return 0
	}
	return len(c.Entries)
}

func (c *Cache) Put(entry Entry) {
	if c.Entries == nil {
		c.Entries = map[string]Entry{}
	}
	entry.LastSeen = time.Now().UTC()
	c.Entries[key(entry.Kind, entry.ID)] = entry
}

func key(kind uid.Kind, id int64) string {
	return fmt.Sprintf("%s:%d", kind, id)
}
