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
	}

	for _, q := range queries {
		if _, err := db.Conn.Exec(q); err != nil {
			log.Fatalf("Migration failed: %v\nQuery: %s", err, q)
		}
	}
	log.Println("Database migrations completed successfully")
}
