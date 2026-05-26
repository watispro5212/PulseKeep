# PulseKeep

PulseKeep is a Go-powered Discord bot and static website for moderation, audit logging, support tickets, economy commands, and live service analytics.

## What Is Included

- Discord bot runtime built with Go and Disgo.
- Gin API with `/health` and `/stats` endpoints for deploy checks and the website.
- Cloudflare Pages static website in `web/`.
- Fly.io deployment config tuned for a single always-on Discord bot machine.
- Drizzle/Postgres schema for guild settings, command logs, bot stats, and economy balances.
- Graceful shutdown on `SIGTERM` so deploys do not leave duplicate gateway sessions behind.

## Project Layout

```text
PulseKeep/
  cmd/pulsekeep/          Go entrypoint
  internal/api/           Gin API server
  internal/bot/           Discord gateway client
  internal/cache/         Thread-safe in-memory cache
  internal/config/        Environment loading
  internal/db/            Postgres connection helper
  db/schema.ts            Drizzle schema
  web/                    Cloudflare Pages static website
  fly.toml                Fly.io app config
  web/_headers            Cloudflare Pages security headers
```

## Local Development

```bash
cp .env.example .env
go run ./cmd/pulsekeep
```

For API-only local testing without a Discord token:

```bash
BOT_DISABLED=true go run ./cmd/pulsekeep
```

Then open the website from `web/index.html` or serve the folder with any static server. In production, the static site reads stats directly from the Fly.io API.

## Fly.io Deployment

```bash
fly launch --no-deploy
fly secrets set DISCORD_TOKEN="your-token-here"
fly secrets set DATABASE_URL="postgresql://..."
fly secrets set ALLOWED_ORIGIN="https://your-cloudflare-site.pages.dev"
fly deploy
```

The Fly config keeps one machine running because Discord bots should not be auto-stopped while connected to the gateway. It also checks `/health` during deploys.

## Cloudflare Pages Deployment

### Automatic (recommended)
1. Go to the [Cloudflare Dashboard](https://dash.cloudflare.com/) → Workers & Pages → Create → Pages → Connect to Git.
2. Select your PulseKeep repository.
3. Set **Build command**: leave empty (no build step).
4. Set **Build output directory**: `web`.
5. Deploy.

### Manual (CLI)
```bash
npx wrangler pages deploy web --branch main
```

The site will be available at `https://<project>.pages.dev`. For a custom domain, go to the Cloudflare Pages dashboard → your project → Custom domains.

**Note:** The `web/_headers` file configures security headers and cache rules — no extra configuration needed.

## Environment Variables

| Variable | Required | Description |
| --- | --- | --- |
| `DISCORD_TOKEN` | For full bot mode | Discord bot token from the Discord Developer Portal. |
| `DATABASE_URL` | For database features | Postgres connection string. |
| `PORT` | No | HTTP port, defaults to `8080`. |
| `ALLOWED_ORIGIN` | No | Comma-separated browser origins allowed to call the API. Use the Cloudflare Pages URL in production. |
| `BOT_DISABLED` | No | Set to `true` to run API-only mode without opening a Discord gateway connection. |

## Database

```bash
npm install
npm run db:generate
npm run db:migrate
```

## Verification

```bash
go build ./...
go test ./...
```

The static site has no build step; validate it by opening `web/index.html` or serving `web/` locally.
