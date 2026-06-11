import pg from 'pg';
const { Pool } = pg;
const pool = new Pool({
  connectionString: 'postgresql://postgres.khcbyncbppidvtipaeoj:Quan52%40watispro1@aws-1-us-east-1.pooler.supabase.com:6543/postgres',
  max: 1,
  connectionTimeoutMillis: 5000,
});
const sql = `
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
try {
  await pool.query(sql);
  console.log('Migration complete');
} catch (e) {
  console.error('Migration error:', e);
}
await pool.end();
