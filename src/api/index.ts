import express from 'express';
import crypto from 'crypto';
import path from 'path';
import { fileURLToPath } from 'url';
import type { Config } from '../config.js';
import type { Cache } from '../cache/index.js';
import { getOAuthURL, exchangeCode, fetchUser, fetchGuilds, fetchMutualGuilds } from './oauth.js';
import { eq } from 'drizzle-orm';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webDir = path.resolve(__dirname, '../../web');

// Permission bit flags
const PERMISSION_ADMINISTRATOR = 0x8n;
const PERMISSION_MANAGE_GUILD = 0x20n;

function hasAdminPerms(permissions: string): boolean {
  try {
    const perms = BigInt(permissions);
    return (perms & PERMISSION_ADMINISTRATOR) === PERMISSION_ADMINISTRATOR ||
           (perms & PERMISSION_MANAGE_GUILD) === PERMISSION_MANAGE_GUILD;
  } catch {
    return false;
  }
}

const STATIC_EXTENSIONS = /\.(html?|css|js|json|xml|png|jpg|jpeg|gif|ico|svg|webp|avif|woff2?|ttf|eot|map)$/i;

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
      if (origin && this.config.allowedOrigins.split(',').map(s => s.trim()).includes(origin)) {
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

  private async verifyGuildAccess(guildId: string, accessToken: string): Promise<{ ok: boolean; error?: string; status?: number }> {
    try {
      const userGuilds = await fetchGuilds(accessToken);
      const guild = userGuilds.find((g: any) => g.id === guildId);
      if (!guild) return { ok: false, error: 'You are not a member of this guild', status: 403 };
      if (!hasAdminPerms(guild.permissions)) return { ok: false, error: 'You need Manage Server or Administrator permission', status: 403 };
      return { ok: true };
    } catch {
      return { ok: false, error: 'Invalid token', status: 401 };
    }
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

    // Bot info endpoint
    this.app.get('/api/bot/info', (_req, res) => {
      res.json({
        name: 'PulseKeep',
        version: '7.0.0',
        library: 'discord.js v14',
        website: 'https://pulsekeep.fly.dev',
        support: 'https://discord.gg/pulsekeep',
        topgg: 'https://top.gg/bot/1507498795569512598',
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
      const { code } = req.query;
      if (!code || typeof code !== 'string') {
        res.status(400).send('Missing authorization code.');
        return;
      }

      try {
        const tokens = await exchangeCode(this.config, code);
        const accessToken = tokens.access_token;
        res.redirect('/dashboard.html?token=' + encodeURIComponent(accessToken));
      } catch (err: any) {
        console.error('OAuth callback error:', err);
        res.status(500).send('Authentication failed.');
      }
    });

    // User guilds (requires Bearer token from OAuth)
    this.app.get('/api/user/guilds', async (req, res) => {
      const auth = req.headers.authorization;
      if (!auth || !auth.startsWith('Bearer ')) {
        res.status(401).json({ error: 'Missing authorization' });
        return;
      }
      const accessToken = auth.slice(7);
      try {
        const userGuilds = await fetchGuilds(accessToken);
        const botGuilds = this.cache.getBotGuilds();
        const mutual = fetchMutualGuilds(userGuilds, botGuilds);
        // Add hasAdmin flag for each guild
        const withAdmin = mutual.map((g: any) => ({
          ...g,
          hasAdmin: hasAdminPerms(g.permissions),
        }));
        res.json(withAdmin);
      } catch {
        res.status(401).json({ error: 'Invalid token' });
      }
    });

    // Guild config (GET - requires admin)
    this.app.get('/api/guild/:id/config', async (req, res) => {
      const auth = req.headers.authorization;
      if (!auth || !auth.startsWith('Bearer ')) {
        res.status(401).json({ error: 'Missing authorization' });
        return;
      }
      const accessToken = auth.slice(7);
      const access = await this.verifyGuildAccess(req.params.id, accessToken);
      if (!access.ok) {
        res.status(access.status!).json({ error: access.error });
        return;
      }

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
        res.json({ guildId: req.params.id });
        return;
      }
      res.json(rows[0]);
    });

    // Guild config (POST - requires admin)
    this.app.post('/api/guild/:id/config', async (req, res) => {
      const auth = req.headers.authorization;
      if (!auth || !auth.startsWith('Bearer ')) {
        res.status(401).json({ error: 'Missing authorization' });
        return;
      }
      const accessToken = auth.slice(7);
      const access = await this.verifyGuildAccess(req.params.id, accessToken);
      if (!access.ok) {
        res.status(access.status!).json({ error: access.error });
        return;
      }

      if (!this.db) {
        res.status(503).json({ error: 'Database unavailable' });
        return;
      }
      const { guildConfigs } = await import('../db/schema.js');
      const guildId = req.params.id;
      const updateData = req.body;
      delete updateData.guildId;
      delete updateData.createdAt;
      delete updateData.updatedAt;

      const existing: any[] = await this.db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);

      if (existing.length > 0) {
        await this.db
          .update(guildConfigs)
          .set({ ...updateData, updatedAt: new Date() })
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

    // SPA fallback - redirect /dashboard to /dashboard.html, handle other paths
    this.app.use((req, res, next) => {
      // Skip API routes
      if (req.path.startsWith('/api/') || req.path.startsWith('/auth/')) {
        res.status(404).json({ error: 'Not found' });
        return;
      }
      // Skip requests with file extensions
      if (STATIC_EXTENSIONS.test(req.path)) {
        res.status(404).send('Not found');
        return;
      }
      // Try appending .html
      const htmlPath = path.join(webDir, req.path + '.html');
      res.sendFile(htmlPath, (err) => {
        if (err) {
          res.sendFile(path.join(webDir, '404.html'));
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
