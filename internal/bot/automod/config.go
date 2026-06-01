package automod

import (
	"context"
	"database/sql"
	"log"
	"sync"
)

type GuildConfig struct {
	GuildID           string `json:"guild_id"`
	AutomodEnabled    bool   `json:"automod_enabled"`
	SpamEnabled       bool   `json:"spam_enabled"`
	SpamMaxMessages   int    `json:"spam_max_messages"`
	SpamWindowSecs    int    `json:"spam_window_seconds"`
	SpamAction        string `json:"spam_action"`
	MentionEnabled    bool   `json:"mention_enabled"`
	MentionMax        int    `json:"mention_max"`
	MentionAction     string `json:"mention_action"`
	LinksEnabled      bool   `json:"links_enabled"`
	LinksAction       string `json:"links_action"`
	CapsEnabled       bool   `json:"caps_enabled"`
	CapsPercent       int    `json:"caps_percent"`
	CapsMinLength     int    `json:"caps_min_length"`
	CapsAction        string `json:"caps_action"`
	BannedWords       string `json:"banned_words"`
	LogChannelID      string `json:"log_channel_id"`
	ModRoleID         string `json:"mod_role_id"`
	EconomyEnabled    bool   `json:"economy_enabled"`
	TicketsEnabled    bool   `json:"tickets_enabled"`
	ModlogsEnabled    bool   `json:"modlogs_enabled"`
	WelcomeEnabled    bool   `json:"welcome_enabled"`
}

func DefaultConfig(guildID string) *GuildConfig {
	return &GuildConfig{
		GuildID:         guildID,
		AutomodEnabled:  true,
		SpamEnabled:     true,
		SpamMaxMessages: 5,
		SpamWindowSecs:  5,
		SpamAction:      "warn",
		MentionEnabled:  true,
		MentionMax:      5,
		MentionAction:   "delete",
		LinksEnabled:    false,
		LinksAction:     "delete",
		CapsEnabled:     true,
		CapsPercent:     70,
		CapsMinLength:   15,
		CapsAction:      "warn",
		EconomyEnabled:  true,
		TicketsEnabled:  true,
		ModlogsEnabled:  true,
		WelcomeEnabled:  false,
	}
}

type ConfigStore struct {
	mu    sync.RWMutex
	db    *sql.DB
	cache map[string]*GuildConfig
}

func NewConfigStore(database *sql.DB) *ConfigStore {
	cs := &ConfigStore{
		db:    database,
		cache: make(map[string]*GuildConfig),
	}
	if database != nil {
		cs.loadAll()
	}
	return cs
}

func (cs *ConfigStore) loadAll() {
	rows, err := cs.db.Query(`SELECT guild_id, automod_enabled, spam_enabled, spam_max_messages, spam_window_seconds, spam_action, mention_enabled, mention_max, mention_action, links_enabled, links_action, caps_enabled, caps_percent, caps_min_length, caps_action, banned_words, log_channel_id, mod_role_id, economy_enabled, tickets_enabled, modlogs_enabled, welcome_enabled FROM guild_config`)
	if err != nil {
		return
	}
	defer rows.Close()

	cs.mu.Lock()
	defer cs.mu.Unlock()

	for rows.Next() {
		var cfg GuildConfig
		if err := rows.Scan(&cfg.GuildID, &cfg.AutomodEnabled, &cfg.SpamEnabled, &cfg.SpamMaxMessages, &cfg.SpamWindowSecs, &cfg.SpamAction, &cfg.MentionEnabled, &cfg.MentionMax, &cfg.MentionAction, &cfg.LinksEnabled, &cfg.LinksAction, &cfg.CapsEnabled, &cfg.CapsPercent, &cfg.CapsMinLength, &cfg.CapsAction, &cfg.BannedWords, &cfg.LogChannelID, &cfg.ModRoleID, &cfg.EconomyEnabled, &cfg.TicketsEnabled, &cfg.ModlogsEnabled, &cfg.WelcomeEnabled); err != nil {
			continue
		}
		cs.cache[cfg.GuildID] = &cfg
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating guild_config rows: %v", err)
	}
}

func (cs *ConfigStore) Get(guildID string) *GuildConfig {
	cs.mu.RLock()
	cfg, ok := cs.cache[guildID]
	cs.mu.RUnlock()
	if ok {
		return cfg
	}
	cfg = DefaultConfig(guildID)
	if cs.db != nil {
		if err := cs.load(guildID, cfg); err == nil {
			cs.mu.Lock()
			cs.cache[guildID] = cfg
			cs.mu.Unlock()
		}
	}
	return cfg
}

func (cs *ConfigStore) load(guildID string, cfg *GuildConfig) error {
	return cs.db.QueryRow(`SELECT automod_enabled, spam_enabled, spam_max_messages, spam_window_seconds, spam_action, mention_enabled, mention_max, mention_action, links_enabled, links_action, caps_enabled, caps_percent, caps_min_length, caps_action, banned_words, log_channel_id, mod_role_id, economy_enabled, tickets_enabled, modlogs_enabled, welcome_enabled FROM guild_config WHERE guild_id = $1`, guildID).Scan(&cfg.AutomodEnabled, &cfg.SpamEnabled, &cfg.SpamMaxMessages, &cfg.SpamWindowSecs, &cfg.SpamAction, &cfg.MentionEnabled, &cfg.MentionMax, &cfg.MentionAction, &cfg.LinksEnabled, &cfg.LinksAction, &cfg.CapsEnabled, &cfg.CapsPercent, &cfg.CapsMinLength, &cfg.CapsAction, &cfg.BannedWords, &cfg.LogChannelID, &cfg.ModRoleID, &cfg.EconomyEnabled, &cfg.TicketsEnabled, &cfg.ModlogsEnabled, &cfg.WelcomeEnabled)
}

func (cs *ConfigStore) Update(cfg *GuildConfig) error {
	if cs.db == nil {
		cs.mu.Lock()
		cs.cache[cfg.GuildID] = cfg
		cs.mu.Unlock()
		return nil
	}

	_, err := cs.db.ExecContext(context.Background(), `
		INSERT INTO guild_config (guild_id, automod_enabled, spam_enabled, spam_max_messages, spam_window_seconds, spam_action, mention_enabled, mention_max, mention_action, links_enabled, links_action, caps_enabled, caps_percent, caps_min_length, caps_action, banned_words, log_channel_id, mod_role_id, economy_enabled, tickets_enabled, modlogs_enabled, welcome_enabled, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,NOW())
		ON CONFLICT (guild_id) DO UPDATE SET
			automod_enabled=$2, spam_enabled=$3, spam_max_messages=$4, spam_window_seconds=$5, spam_action=$6,
			mention_enabled=$7, mention_max=$8, mention_action=$9, links_enabled=$10, links_action=$11,
			caps_enabled=$12, caps_percent=$13, caps_min_length=$14, caps_action=$15,
			banned_words=$16, log_channel_id=$17, mod_role_id=$18,
			economy_enabled=$19, tickets_enabled=$20, modlogs_enabled=$21, welcome_enabled=$22, updated_at=NOW()`,
		cfg.GuildID, cfg.AutomodEnabled, cfg.SpamEnabled, cfg.SpamMaxMessages, cfg.SpamWindowSecs, cfg.SpamAction,
		cfg.MentionEnabled, cfg.MentionMax, cfg.MentionAction, cfg.LinksEnabled, cfg.LinksAction,
		cfg.CapsEnabled, cfg.CapsPercent, cfg.CapsMinLength, cfg.CapsAction,
		cfg.BannedWords, cfg.LogChannelID, cfg.ModRoleID,
		cfg.EconomyEnabled, cfg.TicketsEnabled, cfg.ModlogsEnabled, cfg.WelcomeEnabled)
	if err != nil {
		return err
	}

	cs.mu.Lock()
	cs.cache[cfg.GuildID] = cfg
	cs.mu.Unlock()
	return nil
}
