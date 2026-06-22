import express from 'express';
import crypto from 'crypto';
import path from 'path';
import { fileURLToPath } from 'url';
import { EmbedBuilder, ShardingManager } from 'discord.js';
import type { Config } from '../config.js';
import type { Cache } from '../cache/index.js';
import { sql, eq } from 'drizzle-orm';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webDir = path.resolve(__dirname, '../../web');

const STATIC_EXTENSIONS = /\.(html?|css|js|json|xml|png|jpg|jpeg|gif|ico|svg|webp|avif|woff2?|ttf|eot|map)$/i;

export class ApiServer {
  private app: express.Application;
  private config: Config;
  private cache: Cache;
  private db: any;
  private manager: ShardingManager | null;
  private server: any;

  constructor(config: Config, db: any, cache: Cache, manager: ShardingManager | null = null) {
    this.app = express();
    this.config = config;
    this.cache = cache;
    this.db = db;
    this.manager = manager;

    this.setupMiddleware();
    this.setupRoutes();
  }

  private setupMiddleware() {
    this.app.use(express.json({ limit: '32kb' }));

    this.app.use((_req, res, next) => {
      res.setHeader('X-Content-Type-Options', 'nosniff');
      res.setHeader('X-Frame-Options', 'DENY');
      res.setHeader('X-XSS-Protection', '1; mode=block');
      res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
      next();
    });

    const rateBuckets = new Map<string, { count: number; resetAt: number }>();
    setInterval(() => {
      const now = Date.now();
      for (const [ip, bucket] of rateBuckets) {
        if (now > bucket.resetAt) rateBuckets.delete(ip);
      }
    }, 120_000);
    this.app.use((req, res, next) => {
      if (req.path.startsWith('/api/')) {
        const ip = req.ip || req.socket.remoteAddress || 'unknown';
        const now = Date.now();
        let bucket = rateBuckets.get(ip);
        if (!bucket || now > bucket.resetAt) {
          bucket = { count: 0, resetAt: now + 60000 };
          rateBuckets.set(ip, bucket);
        }
        bucket.count++;
        if (bucket.count > 100) {
          res.status(429).json({ error: 'Too many requests. Try again in a minute.' });
          return;
        }
      }
      next();
    });

    this.app.use((req, res, next) => {
      const origin = req.headers.origin as string;
      const allowed = [...this.config.allowedOrigins.split(',').map(s => s.trim()), 'https://discordbotlist.com'];
      if (origin && allowed.includes(origin)) {
        res.setHeader('Access-Control-Allow-Origin', origin);
        res.setHeader('Vary', 'Origin');
      }
      res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, DELETE');
      res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
      if (req.method === 'OPTIONS') {
        res.sendStatus(204);
        return;
      }
      next();
    });
  }

  private setupRoutes() {
    this.app.get('/health', (_req, res) => {
      res.json({ status: 'ok', uptime: process.uptime() });
    });

    let dbConnected = false;
    let lastDbCheck = 0;
    this.app.get('/api/stats', async (_req, res) => {
      if (this.db && Date.now() - lastDbCheck > 30000) {
        try {
          await this.db.execute(sql`SELECT 1`);
          dbConnected = true;
        } catch { dbConnected = false; }
        lastDbCheck = Date.now();
      }
      res.json({
        guilds: this.cache.getGuildsCount(),
        users: this.cache.getTotalUserCount(),
        commandsRun: this.cache.getCommandsRun(),
        avgLatency: Math.round(this.cache.getAvgLatency()),
        uptime: Math.floor((Date.now() - this.cache.getStartedAt().getTime()) / 1000),
        dbConnected: this.db ? dbConnected : false,
        shards: this.cache.getShardCount(),
      });
    });

    this.app.get('/api/bot/info', (_req, res) => {
      res.json({
        name: 'PulseKeep',
        version: '7.0.0',
        library: 'discord.js v14',
        website: 'https://pulsekeep.fly.dev',
        support: 'https://discord.gg/b9HBphyeuP',
        dbl: 'https://discordbotlist.com/bots/1507498795569512598',
        discords: 'https://discords.com/bots/bot/1507498795569512598',
      });
    });

    // dbl vote webhook
    const dblCors = (req: any, res: any, next: any) => {
      res.setHeader('Access-Control-Allow-Origin', '*');
      res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
      res.setHeader('Access-Control-Allow-Headers', 'Authorization, Content-Type');
      if (req.method === 'OPTIONS') { res.sendStatus(204); return; }
      next();
    };
    this.app.get('/api/dbl/webhook', dblCors, (_req, res) => {
      res.json({ status: 'ok', message: 'DBL webhook endpoint ready' });
    });
    this.app.options('/api/dbl/webhook', dblCors);
    this.app.post('/api/dbl/webhook', dblCors, async (req, res) => {
      console.log('[DBL] Webhook received:', JSON.stringify({ id: req.body?.id, username: req.body?.username }));
      const auth = req.headers.authorization || req.body?.auth;
      if (!auth || auth !== this.config.dblWebhookSecret) {
        console.warn('[DBL] Auth mismatch from', req.ip);
        res.status(401).json({ error: 'Unauthorized' });
        return;
      }
      const userId = req.body?.id;
      if (!userId) {
        console.warn('[DBL] Missing user field in payload');
        res.json({ status: 'ignored', reason: 'missing user' });
        return;
      }
      const reward = 500 + Math.floor(Math.random() * 251);
      try {
        if (this.db) {
          const { userEconomy } = await import('../db/schema.js');
          const rows: any[] = await this.db
            .select()
            .from(userEconomy)
            .where(eq(userEconomy.userId, userId))
            .limit(1);
          if (rows[0]) {
            await this.db
              .update(userEconomy)
              .set({
                balance: (rows[0].balance ?? 0) + reward,
                lastVote: new Date(),
                totalEarned: (rows[0].totalEarned ?? 0) + reward,
                transactions: (rows[0].transactions ?? 0) + 1,
              })
              .where(eq(userEconomy.userId, userId));
          } else {
            await this.db
              .insert(userEconomy)
              .values({ userId, balance: reward, lastVote: new Date(), totalEarned: reward, transactions: 1 });
          }
          console.log(`[DBL] Credited ${reward} pulses to ${userId}`);
        }
      } catch (err) {
        console.error('[DBL] DB error:', err);
      }
      if (this.manager && this.db) {
        try {
          const { guildConfigs } = await import('../db/schema.js');
          const allConfigs: any[] = await this.db.select().from(guildConfigs);
          const guildsWithVoteChannel = (allConfigs || []).filter((c: any) => c.voteChannelId);
          for (const cfg of guildsWithVoteChannel) {
            try {
              const guilds = await this.manager.fetchClientValues('guilds.cache');
              const g = (guilds as any[]).find((g: any) => g.id === cfg.guildId);
              if (g) {
                const channel = g.channels.cache.get(cfg.voteChannelId);
                if (channel?.isTextBased()) {
                  const emb = new EmbedBuilder()
                    .setTitle('🗳️ New Vote!')
                    .setDescription(`<@${userId}> just voted for PulseKeep!\nThanks for the support! 💜`)
                    .setColor(0x7c5cfc)
                    .setTimestamp();
                  await channel.send({ embeds: [emb] }).catch(() => {});
                }
              }
            } catch (innerErr) {
              console.error('[DBL] Channel send error:', innerErr);
            }
          }
        } catch (err) {
          console.error('[DBL] Channel announce error:', err);
        }
      }
      console.log(`[DBL] Vote processed for ${userId}, reward: ${reward}`);
      res.json({ status: 'ok', reward });
    });

    // discords.com vote webhook
    const discordsCors = (req: any, res: any, next: any) => {
      res.setHeader('Access-Control-Allow-Origin', '*');
      res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
      res.setHeader('Access-Control-Allow-Headers', 'Authorization, Content-Type');
      if (req.method === 'OPTIONS') { res.sendStatus(204); return; }
      next();
    };
    this.app.get('/api/discords/webhook', discordsCors, (_req, res) => {
      res.json({ status: 'ok', message: 'Discords.com webhook endpoint ready' });
    });
    this.app.options('/api/discords/webhook', discordsCors);
    this.app.post('/api/discords/webhook', discordsCors, async (req, res) => {
      const auth = req.headers.authorization || '';
      const botId = req.body?.bot;
      const userId = req.body?.user;
      console.log('[Discords] Webhook received:', JSON.stringify({ bot: botId, user: userId }));
      if (!this.config.discordsWebhookSecret || auth !== this.config.discordsWebhookSecret) {
        console.warn('[Discords] Auth mismatch from', req.ip);
        res.status(401).json({ error: 'Unauthorized' });
        return;
      }
      if (!userId || botId !== this.config.discordBotID) {
        res.json({ status: 'ignored', reason: 'invalid payload' });
        return;
      }
      const reward = 500 + Math.floor(Math.random() * 251);
      try {
        if (this.db) {
          const { userEconomy } = await import('../db/schema.js');
          const rows: any[] = await this.db
            .select()
            .from(userEconomy)
            .where(eq(userEconomy.userId, userId))
            .limit(1);
          if (rows[0]) {
            await this.db
              .update(userEconomy)
              .set({
                balance: (rows[0].balance ?? 0) + reward,
                lastVote: new Date(),
                totalEarned: (rows[0].totalEarned ?? 0) + reward,
                transactions: (rows[0].transactions ?? 0) + 1,
              })
              .where(eq(userEconomy.userId, userId));
          } else {
            await this.db
              .insert(userEconomy)
              .values({ userId, balance: reward, lastVote: new Date(), totalEarned: reward, transactions: 1 });
          }
          console.log(`[Discords] Credited ${reward} pulses to ${userId}`);
        }
      } catch (err) {
        console.error('[Discords] DB error:', err);
      }
      if (this.manager && this.db) {
        try {
          const { guildConfigs } = await import('../db/schema.js');
          const allConfigs: any[] = await this.db.select().from(guildConfigs);
          const guildsWithVoteChannel = (allConfigs || []).filter((c: any) => c.voteChannelId);
          for (const cfg of guildsWithVoteChannel) {
            try {
              const guilds = await this.manager.fetchClientValues('guilds.cache');
              const g = (guilds as any[]).find((g: any) => g.id === cfg.guildId);
              if (g) {
                const channel = g.channels.cache.get(cfg.voteChannelId);
                if (channel?.isTextBased()) {
                  const emb = new EmbedBuilder()
                    .setTitle('🗳️ New Vote!')
                    .setDescription(`<@${userId}> just voted for PulseKeep!\nThanks for the support! 💜`)
                    .setColor(0x7c5cfc)
                    .setTimestamp();
                  await channel.send({ embeds: [emb] }).catch(() => {});
                }
              }
            } catch (innerErr) {
              console.error('[Discords] Channel send error:', innerErr);
            }
          }
        } catch (err) {
          console.error('[Discords] Channel announce error:', err);
        }
      }
      console.log(`[Discords] Vote processed for ${userId}, reward: ${reward}`);
      res.json({ status: 'ok', reward });
    });

    // static files
    this.app.use(express.static(webDir));

    // spa fallback
    this.app.use((req, res, next) => {
      if (req.path.startsWith('/api/')) {
        res.status(404).json({ error: 'Not found' });
        return;
      }
      if (STATIC_EXTENSIONS.test(req.path)) {
        res.status(404).send('Not found');
        return;
      }
      const htmlPath = path.join(webDir, req.path + '.html');
      res.sendFile(htmlPath, (err) => {
        if (err) res.sendFile(path.join(webDir, '404.html'));
      });
    });
  }

  start() {
    const port = this.config.port;
    this.server = this.app.listen(port, () => {
      console.log(`API server listening on port ${port}`);
    });
  }

  stop() {
    if (this.server) this.server.close();
  }
}
