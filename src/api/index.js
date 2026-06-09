"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ApiServer = void 0;
const express_1 = __importDefault(require("express"));
const config_1 = require("../config");
const cache_1 = require("../cache");
const automod_1 = require("../bot/automod");
class ApiServer {
    app;
    config;
    cache;
    automod;
    server;
    constructor(config, db, cache, automod) {
        this.app = (0, express_1.default)();
        this.config = config;
        this.cache = cache;
        this.automod = automod;
        this.setupMiddleware();
        this.setupRoutes();
    }
    setupMiddleware() {
        this.app.use(express_1.default.json());
        this.app.use((req, res, next) => {
            const origin = req.headers.origin;
            if (this.config.allowedOrigins.split(',').includes(origin)) {
                res.setHeader('Access-Control-Allow-Origin', origin);
            }
            res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
            res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
            next();
        });
    }
    setupRoutes() {
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
    start() {
        this.server = this.app.listen(this.config.port, () => {
            console.log(`API Server listening on port ${this.config.port}`);
        });
    }
    stop() {
        if (this.server) {
            this.server.close();
        }
    }
}
exports.ApiServer = ApiServer;
//# sourceMappingURL=index.js.map