import { ShardingManager } from 'discord.js';
import { loadConfig } from './config.js';
import { connect } from './db/index.js';
import { Cache } from './cache/index.js';
import { ApiServer } from './api/index.js';
import { getRedis, closeRedis } from './redis.js';
import path from 'path';
import { fileURLToPath } from 'url';

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
  const redis = getRedis(cfg.redisURL);

  const __dirname = path.dirname(fileURLToPath(import.meta.url));
  const manager = new ShardingManager(path.resolve(__dirname, './bot/shard.js'), {
    token: cfg.discordToken,
    totalShards: 'auto',
    respawn: true,
  });

  manager.on('shardCreate', (shard) => {
    console.log(`[Manager] Shard ${shard.id} spawned`);
    shard.on('message', (msg: any) => {
      if (msg?.type === 'stats') {
        memCache.setShardStats(shard.id, msg.data);
        let totalGuilds = 0, totalUsers = 0, totalCommands = 0, totalLatency = 0, shardCount = 0;
        const allStats = memCache.getAllShardStats();
        for (const s of allStats) {
          totalGuilds += s.guilds;
          totalUsers += s.users;
          totalCommands += s.commandsRun;
          totalLatency += s.avgLatency;
          shardCount++;
        }
        memCache.setGuildsCount(totalGuilds);
        memCache.setTotalUserCount(totalUsers);
        memCache.setCommandsRun(totalCommands);
        if (shardCount > 0) memCache.setAvgLatency(totalLatency / shardCount);
      }
    });
  });

  if (!cfg.botDisabled) {
    await manager.spawn({ timeout: 120000 });
    console.log(`[Manager] ${manager.totalShards} shard(s) online`);
  } else {
    console.log('Bot disabled; API server will run standalone.');
  }

  const server = new ApiServer(cfg, database, memCache, manager);
  server.start();

  console.log('PulseKeep is running.');

  const shutdown = async () => {
    console.log('Shutting down...');
    server.stop();
    closeRedis();
    if (!cfg.botDisabled) { for (const [, s] of manager.shards) s.kill(); }
    process.exit(0);
  };

  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
}

main().catch((err) => {
  console.error('Fatal:', err);
  process.exit(1);
});
