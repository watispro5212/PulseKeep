import { Config } from '../config';
import { Cache } from '../cache';
import { AutomodEngine } from '../bot/automod';
export declare class ApiServer {
    private app;
    private config;
    private cache;
    private automod;
    private server;
    constructor(config: Config, db: any, cache: Cache, automod: AutomodEngine);
    private setupMiddleware;
    private setupRoutes;
    start(): void;
    stop(): void;
}
//# sourceMappingURL=index.d.ts.map