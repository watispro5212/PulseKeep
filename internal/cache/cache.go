package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	mu            sync.RWMutex
	data          map[string]interface{}
	CommandsRun   atomic.Int64
	GuildCount    atomic.Int64
	UserCount     atomic.Int64
	StartedAt     time.Time
	muGuilds      sync.RWMutex
	guildNames    map[string]string
}

func New() *Cache {
	return &Cache{
		data:       make(map[string]interface{}),
		StartedAt:  time.Now(),
		guildNames: make(map[string]string),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) IncrCommands() {
	c.CommandsRun.Add(1)
}

func (c *Cache) SetGuilds(count int64) {
	c.GuildCount.Store(count)
}

func (c *Cache) SetUsers(count int64) {
	c.UserCount.Store(count)
}

func (c *Cache) AddGuild(guildID, name string) {
	c.muGuilds.Lock()
	defer c.muGuilds.Unlock()
	c.guildNames[guildID] = name
	c.GuildCount.Add(1)
}

func (c *Cache) RemoveGuild(guildID string) {
	c.muGuilds.Lock()
	defer c.muGuilds.Unlock()
	delete(c.guildNames, guildID)
	c.GuildCount.Add(-1)
}
