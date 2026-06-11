import pg from 'pg';
const { Pool } = pg;
const p = new Pool({ connectionString: 'postgresql://postgres.khcbyncbppidvtipaeoj:Quan52%40watispro1@aws-1-us-east-1.pooler.supabase.com:6543/postgres', max: 1, connectionTimeoutMillis: 5000 });
await p.query("ALTER TABLE user_economy ADD COLUMN IF NOT EXISTS last_vote timestamp;");
console.log('Migration OK');
await p.end();
