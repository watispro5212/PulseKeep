import { loadConfig } from './config/index.js';
import { connect } from './db/index.js';
import { Cache } from './cache/index.js';
import { Bot } from './bot/index.js';
import { ApiServer } from './api/index.js';

async function main() {
    // 1. Load Config
    const cfg = loadConfig();

    // 2. Init Database
    const database = connect(cfg.databaseURL);

    // 3. Init Cache
    const memCache = new Cache();

    // 4. Init Bot
    let discordBot: Bot | null = null;
    if (cfg.botDisabled) {
        console.log('Discord bot is disabled; API endpoints are still available.');
    } else {
        discordBot = new Bot(cfg.discordToken, memCache, database, cfg.statusWebhookURL);
        await discordBot.start(cfg.discordToken);
    }

    // 5. Init Web Server
    const automod = discordBot ? discordBot.getAutomod() : (new (require('./bot/automod').AutomodEngine)());
    const server = new ApiServer(cfg, database, memCache, automod);
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
