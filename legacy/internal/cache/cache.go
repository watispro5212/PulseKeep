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
	muAuth        sync.Mutex
	authorized    map[string]struct{}
	AuthorizedCnt atomic.Int64
	latencyMu     sync.Mutex
	latencies     []time.Duration
	latencyMax    int
	latencyTotal  time.Duration
	latencyCount  int64
}

func New() *Cache {
	return &Cache{
		data:       make(map[string]interface{}),
		StartedAt:  time.Now(),
		guildNames: make(map[string]string),
		authorized: make(map[string]struct{}),
		latencies:  make([]time.Duration, 0, 100),
		latencyMax: 100,
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

func (c *Cache) AddLatency(d time.Duration) {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()
	c.latencyTotal += d
	c.latencyCount++
	if len(c.latencies) >= c.latencyMax {
		c.latencies = c.latencies[1:]
	}
	c.latencies = append(c.latencies, d)
}

func (c *Cache) AvgLatency() time.Duration {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()
	if c.latencyCount == 0 {
		return 0
	}
	return c.latencyTotal / time.Duration(c.latencyCount)
}

func (c *Cache) LatencyCount() int64 {
	c.latencyMu.Lock()
	defer c.latencyMu.Unlock()
	return c.latencyCount
}

func (c *Cache) SetGuilds(count int64) {
	c.GuildCount.Store(count)
}

func (c *Cache) SetUsers(count int64) {
	c.UserCount.Store(count)
}

func (c *Cache) ResetGuilds() {
	c.muGuilds.Lock()
	defer c.muGuilds.Unlock()
	c.guildNames = make(map[string]string)
	c.GuildCount.Store(0)
}

func (c *Cache) AddGuild(guildID, name string) {
	c.muGuilds.Lock()
	defer c.muGuilds.Unlock()
	if _, exists := c.guildNames[guildID]; !exists {
		c.guildNames[guildID] = name
		c.GuildCount.Store(int64(len(c.guildNames)))
	}
}

func (c *Cache) GuildsCount() int64 {
	c.muGuilds.RLock()
	defer c.muGuilds.RUnlock()
	return int64(len(c.guildNames))
}

func (c *Cache) GetGuildIDs() []string {
	c.muGuilds.RLock()
	defer c.muGuilds.RUnlock()
	ids := make([]string, 0, len(c.guildNames))
	for id := range c.guildNames {
		ids = append(ids, id)
	}
	return ids
}

func (c *Cache) GetGuildNames() map[string]string {
	c.muGuilds.RLock()
	defer c.muGuilds.RUnlock()
	cp := make(map[string]string, len(c.guildNames))
	for k, v := range c.guildNames {
		cp[k] = v
	}
	return cp
}

func (c *Cache) AddAuthorizedUser(userID string) {
	c.muAuth.Lock()
	defer c.muAuth.Unlock()
	if _, exists := c.authorized[userID]; !exists {
		c.authorized[userID] = struct{}{}
		c.AuthorizedCnt.Add(1)
	}
}

func (c *Cache) GetTotalUserCount() int64 {
	guildUsers := c.UserCount.Load()
	authUsers := c.AuthorizedCnt.Load()
	if authUsers > guildUsers {
		return authUsers
	}
	return guildUsers
}

func (c *Cache) RemoveGuild(guildID string) {
	c.muGuilds.Lock()
	defer c.muGuilds.Unlock()
	if _, exists := c.guildNames[guildID]; exists {
		delete(c.guildNames, guildID)
	}
	c.GuildCount.Store(int64(len(c.guildNames)))
}
