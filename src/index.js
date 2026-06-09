"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const config_1 = require("./config");
const db_1 = require("./db");
const cache_1 = require("./cache");
const bot_1 = require("./bot");
const api_1 = require("./api");
async function main() {
    // 1. Load Config
    const cfg = (0, config_1.loadConfig)();
    // 2. Init Database
    const database = (0, db_1.connect)(cfg.databaseURL);
    // 3. Init Cache
    const memCache = new cache_1.Cache();
    // 4. Init Bot
    let discordBot = null;
    if (cfg.botDisabled) {
        console.log('Discord bot is disabled; API endpoints are still available.');
    }
    else {
        discordBot = new bot_1.Bot(cfg.discordToken, memCache, database, cfg.statusWebhookURL);
        await discordBot.start(cfg.discordToken);
    }
    // 5. Init Web Server
    const automod = discordBot ? discordBot.getAutomod() : (new (require('./bot/automod').AutomodEngine)());
    const server = new api_1.ApiServer(cfg, database, memCache, automod);
    server.start();
    console.log('PulseKeep (TypeScript version) is running.');
    // Graceful shutdown
    const shutdown = async () => {
        console.log('Shutdown signal received. Initiating graceful cleanup...');
        server.stop();
        if (discordBot) {
            await discordBot.stop();
        }
        process.exit(0);
    };
    process.on('SIGINT', shutdown);
    process.on('SIGTERM', shutdown);
}
main().catch(err => {
    console.error('Failed to start application:', err);
    process.exit(1);
});
//# sourceMappingURL=index.js.map