package db

import "log"

func (db *Database) Migrate() {
	if db == nil || db.Conn == nil {
		return
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS user_economy (
			user_id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			balance INTEGER NOT NULL DEFAULT 250,
			total_earned INTEGER NOT NULL DEFAULT 0,
			total_spent INTEGER NOT NULL DEFAULT 0,
			daily_streak INTEGER NOT NULL DEFAULT 0,
			last_daily TIMESTAMPTZ,
			last_work TIMESTAMPTZ,
			last_rob TIMESTAMPTZ,
			last_fish TIMESTAMPTZ,
			last_mine TIMESTAMPTZ,
			last_gamble TIMESTAMPTZ,
			last_interest TIMESTAMPTZ,
			flip_wins INTEGER NOT NULL DEFAULT 0,
			flip_losses INTEGER NOT NULL DEFAULT 0,
			slot_wins INTEGER NOT NULL DEFAULT 0,
			slot_losses INTEGER NOT NULL DEFAULT 0,
			rob_wins INTEGER NOT NULL DEFAULT 0,
			rob_losses INTEGER NOT NULL DEFAULT 0,
			fish_caught INTEGER NOT NULL DEFAULT 0,
			fish_total INTEGER NOT NULL DEFAULT 0,
			mine_mined INTEGER NOT NULL DEFAULT 0,
			mine_total INTEGER NOT NULL DEFAULT 0,
			gamble_wins INTEGER NOT NULL DEFAULT 0,
			gamble_total INTEGER NOT NULL DEFAULT 0,
			lottery_wins INTEGER NOT NULL DEFAULT 0,
			weekly_streak INTEGER NOT NULL DEFAULT 0,
			last_weekly TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS user_inventory (
			id SERIAL PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES user_economy(user_id),
			item_id TEXT NOT NULL,
			item_name TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			UNIQUE(user_id, item_id)
		)`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS balance INTEGER NOT NULL DEFAULT 250`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS total_earned INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS total_spent INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS daily_streak INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_daily TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_work TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_rob TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_fish TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_mine TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_gamble TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_interest TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS flip_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS flip_losses INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS slot_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS slot_losses INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS rob_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS rob_losses INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS fish_caught INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS fish_total INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS mine_mined INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS mine_total INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS gamble_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS gamble_total INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS lottery_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS weekly_streak INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_weekly TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS xp_boost_expires TIMESTAMPTZ`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS treasure_map_active BOOLEAN NOT NULL DEFAULT false`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS gamble_boost_active BOOLEAN NOT NULL DEFAULT false`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS blackjack_wins INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS blackjack_losses INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		`ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
		`ALTER TABLE guild_config ADD COLUMN IF NOT EXISTS economy_enabled BOOLEAN NOT NULL DEFAULT true`,
		`ALTER TABLE guild_config ADD COLUMN IF NOT EXISTS tickets_enabled BOOLEAN NOT NULL DEFAULT true`,
		`ALTER TABLE guild_config ADD COLUMN IF NOT EXISTS modlogs_enabled BOOLEAN NOT NULL DEFAULT true`,
		`ALTER TABLE guild_config ADD COLUMN IF NOT EXISTS welcome_enabled BOOLEAN NOT NULL DEFAULT false`,
		`CREATE TABLE IF NOT EXISTS lottery_entries (
			id SERIAL PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES user_economy(user_id),
			entry_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			week_start TIMESTAMPTZ NOT NULL,
			claimed BOOLEAN NOT NULL DEFAULT false
		)`,
		`CREATE TABLE IF NOT EXISTS lottery_config (
			guild_id TEXT PRIMARY KEY,
			auto_draw BOOLEAN NOT NULL DEFAULT true,
			last_draw_time TIMESTAMPTZ,
			last_winner_id TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS command_logs (
			id SERIAL PRIMARY KEY,
			guild_id TEXT NOT NULL DEFAULT '',
			user_id TEXT NOT NULL DEFAULT '',
			command_name TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS guild_config (
			guild_id TEXT PRIMARY KEY,
			automod_enabled BOOLEAN NOT NULL DEFAULT true,
			spam_enabled BOOLEAN NOT NULL DEFAULT true,
			spam_max_messages INTEGER NOT NULL DEFAULT 5,
			spam_window_seconds INTEGER NOT NULL DEFAULT 5,
			spam_action TEXT NOT NULL DEFAULT 'warn',
			mention_enabled BOOLEAN NOT NULL DEFAULT true,
			mention_max INTEGER NOT NULL DEFAULT 5,
			mention_action TEXT NOT NULL DEFAULT 'delete',
			links_enabled BOOLEAN NOT NULL DEFAULT false,
			links_action TEXT NOT NULL DEFAULT 'delete',
			caps_enabled BOOLEAN NOT NULL DEFAULT true,
			caps_percent INTEGER NOT NULL DEFAULT 70,
			caps_min_length INTEGER NOT NULL DEFAULT 15,
			caps_action TEXT NOT NULL DEFAULT 'warn',
			banned_words TEXT NOT NULL DEFAULT '',
			log_channel_id TEXT NOT NULL DEFAULT '',
			mod_role_id TEXT NOT NULL DEFAULT '',
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS user_warnings (
			id SERIAL PRIMARY KEY,
			guild_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			moderator_id TEXT NOT NULL,
			reason TEXT NOT NULL DEFAULT 'No reason provided',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_warnings_guild_user ON user_warnings(guild_id, user_id)`,
	}

	for _, q := range queries {
		if _, err := db.Conn.Exec(q); err != nil {
			log.Printf("Migration warning: %v\nQuery: %s", err, q)
		}
	}
	log.Println("Database migrations completed successfully")
}
