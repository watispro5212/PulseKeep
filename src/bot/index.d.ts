import { Cache } from '../cache';
import { AutomodEngine } from './automod';
export declare class Bot {
    private client;
    private cache;
    private db;
    private automod;
    private economy;
    private statusWebhookURL;
    constructor(token: string, cache: Cache, db: any, statusWebhookURL: string);
    private registerEvents;
    start(token: string): Promise<void>;
    stop(): Promise<void>;
    getAutomod(): AutomodEngine;
}
//# sourceMappingURL=index.d.ts.map