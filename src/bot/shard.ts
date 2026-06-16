import { loadConfig } from '../config.js';
import { connect } from '../db/index.js';
import { Cache } from '../cache/index.js';
import { Bot } from './client.js';
import { commands } from './commands/index.js';
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
  }

  const sendStats = () => {
    if (!process.send) return;
    process.send({
      type: 'stats',
      data: {
        guilds: memCache.getGuildsCount(),
        users: memCache.getTotalUserCount(),
        commandsRun: memCache.getCommandsRun(),
        avgLatency: memCache.getAvgLatency(),
        uptime: Math.floor((Date.now() - memCache.getStartedAt().getTime()) / 1000),
      },
    });
  };

  sendStats();
  setInterval(sendStats, 15000);

  process.on('SIGINT', () => { discordBot.stop(); process.exit(0); });
  process.on('SIGTERM', () => { discordBot.stop(); process.exit(0); });
}

main().catch((err) => {
  console.error('[Shard] Fatal:', err);
  process.exit(1);
});
