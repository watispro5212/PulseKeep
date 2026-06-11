import express from 'express';
import crypto from 'crypto';
import path from 'path';
import { fileURLToPath } from 'url';
import type { Config } from '../config.js';
import type { Cache } from '../cache/index.js';
import { getOAuthURL, exchangeCode, fetchUser, fetchGuilds } from './oauth.js';
import { eq } from 'drizzle-orm';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webDir = path.resolve(__dirname, '../../web');

export class ApiServer {
  private app: express.Application;
  private config: Config;
  private cache: Cache;
  private db: any;
  private server: any;

  constructor(config: Config, db: any, cache: Cache) {
    this.app = express();
    this.config = config;
    this.cache = cache;
    this.db = db;

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
    // Health
    this.app.get('/health', (_req, res) => {
      res.json({ status: 'ok', uptime: process.uptime() });
    });

    // Stats
    this.app.get('/api/stats', (_req, res) => {
      res.json({
        guilds: this.cache.getGuildsCount(),
        users: this.cache.getTotalUserCount(),
        commandsRun: this.cache.getCommandsRun(),
        avgLatency: this.cache.getAvgLatency(),
        uptime: Math.floor((Date.now() - this.cache.getStartedAt().getTime()) / 1000),
      });
    });

    // OAuth login
    this.app.get('/auth/discord/login', (_req, res) => {
      const state = crypto.randomBytes(16).toString('hex');
      const url = getOAuthURL(this.config, state);
      res.redirect(url);
    });

    // OAuth callback
    this.app.get('/auth/discord/callback', async (req, res) => {
      const { code, state } = req.query;
      if (!code || typeof code !== 'string') {
        res.status(400).send('Missing authorization code.');
        return;
      }

      try {
        const tokens = await exchangeCode(this.config, code);
        const user = await fetchUser(tokens.access_token);
        const accessToken = tokens.access_token;

        res.json({ user, accessToken });
      } catch (err: any) {
        console.error('OAuth callback error:', err);
        res.status(500).send('Authentication failed.');
      }
    });

    // Guild config
    this.app.get('/api/guild/:id/config', async (req, res) => {
      if (!this.db) {
        res.status(503).json({ error: 'Database unavailable' });
        return;
      }
      const { guildConfigs } = await import('../db/schema.js');
      const rows: any[] = await this.db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, req.params.id))
        .limit(1);
      if (rows.length === 0) {
        res.status(404).json({ error: 'Config not found' });
        return;
      }
      res.json(rows[0]);
    });

    this.app.post('/api/guild/:id/config', async (req, res) => {
      if (!this.db) {
        res.status(503).json({ error: 'Database unavailable' });
        return;
      }
      const { guildConfigs } = await import('../db/schema.js');
      const guildId = req.params.id;
      const updateData = req.body;
      delete updateData.guildId;

      const existing: any[] = await this.db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);

      if (existing.length > 0) {
        await this.db
          .update(guildConfigs)
          .set(updateData)
          .where(eq(guildConfigs.guildId, guildId));
      } else {
        await this.db
          .insert(guildConfigs)
          .values({ guildId, ...updateData });
      }

      res.json({ status: 'ok' });
    });

    // Serve static files
    this.app.use(express.static(webDir));

    // SPA fallback for dashboard routes
    this.app.get('*', (req, res) => {
      const filePath = path.join(webDir, req.path);
      res.sendFile(filePath, (err) => {
        if (err) {
          res.sendFile(path.join(webDir, 'dashboard.html'));
        }
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
