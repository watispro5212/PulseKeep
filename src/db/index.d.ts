import { Pool } from 'pg';
import * as schema from './schema';
export declare const connect: (databaseURL: string) => (import("drizzle-orm/node-postgres").NodePgDatabase<typeof schema> & {
    $client: Pool;
}) | null;
//# sourceMappingURL=index.d.ts.map