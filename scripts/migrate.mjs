import 'dotenv/config';
import pg from 'pg';
const { Pool } = pg;

const DATABASE_URL = process.env.DATABASE_URL;
if (!DATABASE_URL) {
  console.error('DATABASE_URL not set. Create a .env file with DATABASE_URL.');
  process.exit(1);
}

const pool = new Pool({
  connectionString: DATABASE_URL,
  max: 1,
  connectionTimeoutMillis: 10000,
});

const bootstrapSQL = `
DO $$
BEGIN
  CREATE TABLE IF NOT EXISTS bot_stats (id integer PRIMARY KEY, guild_count integer DEFAULT 0 NOT NULL, user_count integer DEFAULT 0 NOT NULL, commands_run integer DEFAULT 0 NOT NULL, updated_at timestamp DEFAULT now() NOT NULL);
  CREATE TABLE IF NOT EXISTS command_logs (id serial PRIMARY KEY, guild_id varchar(20), user_id varchar(20) NOT NULL, command_name varchar(50) NOT NULL, executed_at timestamp DEFAULT now() NOT NULL);
  CREATE TABLE IF NOT EXISTS guild_configs (guild_id varchar(20) PRIMARY KEY, prefix varchar(10) DEFAULT '!' NOT NULL, log_channel_id varchar(20), premium boolean DEFAULT false NOT NULL, economy_enabled boolean DEFAULT true NOT NULL, tickets_enabled boolean DEFAULT true NOT NULL, modlogs_enabled boolean DEFAULT true NOT NULL, welcome_enabled boolean DEFAULT false NOT NULL, ticket_category_id varchar(20), created_at timestamp DEFAULT now() NOT NULL, updated_at timestamp DEFAULT now() NOT NULL);
  CREATE TABLE IF NOT EXISTS user_inventory (id serial PRIMARY KEY, user_id varchar(20) NOT NULL, item_id varchar(50) NOT NULL, item_name varchar(100) NOT NULL, quantity integer DEFAULT 1 NOT NULL);
  CREATE TABLE IF NOT EXISTS user_warnings (id serial PRIMARY KEY, guild_id varchar(20) NOT NULL, user_id varchar(20) NOT NULL, moderator_id varchar(20) NOT NULL, reason text NOT NULL, created_at timestamp DEFAULT now() NOT NULL);
  CREATE TABLE IF NOT EXISTS user_economy (user_id varchar(20) PRIMARY KEY, balance integer DEFAULT 0 NOT NULL);
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_daily_claim timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_weekly_claim timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_work timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_rob timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_fish timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_mine timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS total_earned integer DEFAULT 0 NOT NULL;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS total_gambled integer DEFAULT 0 NOT NULL;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS transactions integer DEFAULT 0 NOT NULL;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS streak integer DEFAULT 0 NOT NULL;
END
$$;
`;

const migrationSQL = `
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS welcome_channel_id varchar(20);
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS vote_channel_id varchar(20);
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_spam_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_mention_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_link_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_caps_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_words_enabled boolean DEFAULT true NOT NULL;
  ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS automod_banned_words text;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_search timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_vote timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS xp_boost_expiry timestamp;
  ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS lucky_clover_active integer DEFAULT 0 NOT NULL;
`;

try {
  console.log('Running bootstrap migration...');
  await pool.query(bootstrapSQL);
  console.log('Bootstrap complete. Running column migrations...');
  await pool.query(migrationSQL);
  console.log('Migration complete — all columns added.');
} catch (e) {
  console.error('Migration error:', e);
  process.exit(1);
} finally {
  await pool.end();
}
