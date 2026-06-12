import pg from 'pg';
const { Pool } = pg;
const p = new Pool({
  connectionString: 'postgresql://postgres.khcbyncbppidvtipaeoj:Quan52%40watispro1@aws-1-us-east-1.pooler.supabase.com:6543/postgres',
  max: 1,
  connectionTimeoutMillis: 5000,
});
try {
  await p.query(`ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS welcome_channel_id varchar(20);`);
  await p.query(`ALTER TABLE guild_configs ADD COLUMN IF NOT EXISTS vote_channel_id varchar(20);`);
  console.log('Migration OK — added vote_channel_id and welcome_channel_id');
} catch (e) {
  console.error('Migration error:', e);
}
await p.end();
