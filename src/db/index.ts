import { drizzle } from 'drizzle-orm/node-postgres';
import { Pool } from 'pg';
import * as schema from './schema.js';

export const connect = (databaseURL: string) => {
    if (!databaseURL) {
        console.warn('DATABASE_URL is not set; running without database.');
        return null;
    }

    const pool = new Pool({
        connectionString: databaseURL,
        max: 10,
        idleTimeoutMillis: 30000,
        connectionTimeoutMillis: 2000,
    });

    pool.on('error', (err) => {
        console.error('Unexpected error on idle client', err);
    });

    console.log('Database connection pool initialized.');
    
    return drizzle(pool, { schema });
};
