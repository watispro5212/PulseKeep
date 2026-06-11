import { drizzle } from 'drizzle-orm/node-postgres';
import pg from 'pg';
import * as schema from './schema.js';

const { Pool } = pg;

export function connect(databaseURL: string) {
  if (!databaseURL) {
    console.warn('DATABASE_URL not set; running without database.');
    return null;
  }

  const pool = new Pool({
    connectionString: databaseURL,
    max: 10,
    idleTimeoutMillis: 30000,
    connectionTimeoutMillis: 2000,
  });

  pool.on('error', (err) => {
    console.error('Unexpected pool error:', err);
  });

  console.log('Database pool initialized.');
  return drizzle(pool, { schema });
}
