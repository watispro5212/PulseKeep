# PulseKeep

A Go-powered Discord bot for moderation, audit logging, support tickets, economy games, and live service analytics.

## What Is Included

- **Discord bot** — slash commands, event handling, component interactions, auto-moderation, economy loop
- **Gin API** — `/health`, `/stats`, and `/api/dashboard` endpoints with CORS support
- **Cloudflare Workers website** — static assets in `web/` plus a Worker gateway for API, auth, health, and stats
- **Fly.io deployment** — single always-on machine tuned for Discord gateway uptime
- **PostgreSQL persistence** — guild config, economy records, inventory, command logs
- **Graceful shutdown** — clean SIGTERM handling so deploys don't duplicate gateway sessions

## Project Layout

```text
PulseKeep/
  cmd/pulsekeep/          Go entrypoint
  internal/
    api/                  Gin API server (health, stats, OAuth, guild config)
    auth/                 Discord OAuth2 helpers
    bot/
      economy/            In-memory economy store with PostgreSQL persistence
      automod/            Configurable auto-moderation engine
      commands/           Slash command registration
      handlers.go         Event listeners and command dispatch
      economy_handlers.go Economy response builders
    cache/                Thread-safe atomic counters for live stats
    config/               Environment variable loader
    db/                   PostgreSQL connection and migration runner
  web/                    Cloudflare Workers static assets (HTML/CSS/JS)
  src/worker.js           Worker gateway that proxies API requests to Fly.io
  db/                     Drizzle schema and migration files
  fly.toml                Fly.io app configuration
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26 |
| Discord library | disgo v0.19.3 |
| HTTP router | Gin v1.12.0 |
| Database | PostgreSQL via pgx v5 / Neon |
| Hosting | Fly.io (iad, 256 MB, shared-cpu-1x) |
| Frontend | HTML + CSS + vanilla JS (Cloudflare Workers static assets) |
| CI/CD | GitHub Actions (build + deploy to Fly.io) |

## Local Development

```bash
cp .env.example .env
# Fill in DISCORD_TOKEN and DATABASE_URL
go run ./cmd/pulsekeep
```

For API-only mode without connecting to Discord:

```bash
BOT_DISABLED=true go run ./cmd/pulsekeep
```

The static site has no build step — open `web/index.html` or serve `web/` with any static server.

## Fly.io Deployment

```bash
fly launch --no-deploy
fly secrets set DISCORD_TOKEN="your-token-here"
fly secrets set DATABASE_URL="postgresql://..."
fly secrets set ALLOWED_ORIGIN="https://pulsekeep.williamdelilah3.workers.dev"
fly secrets set DISCORD_CLIENT_ID="..."
fly secrets set DISCORD_CLIENT_SECRET="..."
fly secrets set DISCORD_REDIRECT_URI="https://your-site.com/auth/discord/callback"
fly deploy
```

The Fly config keeps one machine running 24/7 because the Discord gateway requires a persistent connection. Health checks use `/health` and always return 200 (database status is reported in the response body).

## Cloudflare Workers Deployment

The Worker serves the website from `web/` and proxies API traffic to the Fly.io backend. The backend origin is configured in `wrangler.jsonc` and `wrangler.toml`:

```bash
npm run deploy
```

Current origin:

```text
API_ORIGIN=https://pulsekeep.fly.dev
```

Environment variables on the Fly.io backend (`ALLOWED_ORIGIN`) must include the Worker domain for CORS to work with direct API calls and OAuth redirects.

## Environment Variables

| Variable | Required | Description |
| --- | --- | --- |
| `DISCORD_TOKEN` | For full bot mode | Discord bot token from the Discord Developer Portal |
| `DATABASE_URL` | For database features | PostgreSQL connection string (Neon recommended) |
| `PORT` | No | HTTP port, defaults to `8080` |
| `ALLOWED_ORIGIN` | No | CORS origin(s) for API access, defaults to `*` |
| `BOT_DISABLED` | No | Set `true` to run API-only without Discord gateway |
| `DISCORD_CLIENT_ID` | For dashboard | Discord OAuth2 client ID |
| `DISCORD_CLIENT_SECRET` | For dashboard | Discord OAuth2 client secret |
| `DISCORD_REDIRECT_URI` | For dashboard | OAuth callback URL |
| `STATUS_WEBHOOK_URL` | No | Discord webhook for status alerts |

## Database

The project uses Drizzle for schema management:

```bash
npm install
npm run db:generate
npm run db:migrate
```

## Verification

```bash
go build ./...
go vet ./...
go test ./...
```

## Links

- Website: https://pulsekeep.williamdelilah3.workers.dev
- Commands: https://pulsekeep.williamdelilah3.workers.dev/commands.html
- Status: https://pulsekeep.williamdelilah3.workers.dev/status.html
- GitHub: https://github.com/watispro5212/PulseKeep
