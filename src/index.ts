import { loadConfig } from './config.js';
import { connect } from './db/index.js';
import { Cache } from './cache/index.js';
import { Bot } from './bot/client.js';
import { commands } from './bot/commands/index.js';
import { ApiServer } from './api/index.js';

process.on('unhandledRejection', (reason) => {
  console.error('Unhandled rejection:', reason);
});
process.on('uncaughtException', (err) => {
  console.error('Uncaught exception:', err);
  process.exit(1);
});

async function main() {
  const cfg = loadConfig();
  const database = connect(cfg.databaseURL);
  const memCache = new Cache();

  const discordBot = new Bot(cfg.discordToken, memCache, database, cfg.statusWebhookURL, cfg);

  for (const cmd of commands) {
    discordBot.registerCommand(cmd);
  }

  if (!cfg.botDisabled) {
    await discordBot.start(cfg.discordToken);
    await discordBot.registerSlashCommands();
  } else {
    console.log('Bot disabled; API server will run standalone.');
  }

  const server = new ApiServer(cfg, database, memCache, discordBot);
  server.start();

  console.log('PulseKeep is running.');

  const shutdown = async () => {
    console.log('Shutting down...');
    server.stop();
    if (!cfg.botDisabled) await discordBot.stop();
    process.exit(0);
  };

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
}

main().catch((err) => {
  console.error('Fatal:', err);
  process.exit(1);
});
