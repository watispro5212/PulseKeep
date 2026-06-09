import express from 'express';
import type { Config } from '../config/index.js';
import { Cache } from '../cache/index.js';
import { AutomodEngine } from '../bot/automod/index.js';

export class ApiServer {
    private app: express.Application;
    private config: Config;
    private cache: Cache;
    private automod: AutomodEngine;
    private server: any;

    constructor(config: Config, db: any, cache: Cache, automod: AutomodEngine) {
        this.app = express();
        this.config = config;
        this.cache = cache;
        this.automod = automod;

        this.setupMiddleware();
        this.setupRoutes();
    }

    private setupMiddleware() {
        this.app.use(express.json());
        this.app.use((req, res, next) => {
            const origin = req.headers.origin as string;
            if (this.config.allowedOrigins.split(',').includes(origin)) {
                res.setHeader('Access-Control-Allow-Origin', origin);
            }
            res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
            res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
            next();
        });
    }

    private setupRoutes() {
        this.app.get('/health', (req, res) => {
            res.json({ status: 'ok', uptime: process.uptime() });
        });

        this.app.get('/api/stats', (req, res) => {
            res.json({
                guilds: this.cache.getGuildsCount(),
                users: this.cache.getTotalUserCount(),
                commandsRun: this.cache.getCommandsRun(),
                avgLatency: this.cache.getAvgLatency(),
                uptime: Math.floor((Date.now() - this.cache.getStartedAt().getTime()) / 1000),
            });
        });

        this.app.get('/api/guild/:id/config', (req, res) => {
            const guildID = req.params.id;
            const config = this.automod.getConfig(guildID);
            if (!config) {
                return res.status(404).json({ error: 'Config not found' });
            }
            res.json(config);
        });

        this.app.post('/api/guild/:id/config', (req, res) => {
            const guildID = req.params.id;
            // Permission check logic would go here
            const newConfig = req.body;
            newConfig.guildID = guildID;
            this.automod.updateConfig(newConfig);
            res.json({ status: 'ok' });
        });

        // 404 handler
        this.app.use((req, res) => {
            res.status(404).send('Not Found');
        });
    }

    public start() {
        this.server = this.app.listen(this.config.port, () => {
            console.log(`API Server listening on port ${this.config.port}`);
        });
    }

    public stop() {
        if (this.server) {
            this.server.close();
        }
    }
}
