# PulseKeep

> **The Ultimate Discord Sentinel** — Fast, secure, and Go-powered.

PulseKeep is a high-performance Discord bot built with Go, designed for server monitoring, secure ticket management, and real-time health analytics.

---

## ✨ Features

- **Slash Commands** — `/ping`, `/help`, `/stats`, `/serverinfo`, `/userinfo`, `/avatar`, `/purge`, `/kick`, `/ban`, `/announce`, `/uptime`
- **Graceful Shutdown** — Handles `SIGTERM` cleanly so Fly.io updates never crash the bot or trigger Discord token suspensions
- **Live Stats API** — Gin web server exposing `/stats` and `/health` endpoints consumed by the Netlify landing page
- **In-Memory Cache** — Thread-safe `sync.RWMutex` cache for rapid data access
- **PostgreSQL Backend** — Drizzle ORM + pgx for schema management and persistence
- **CORS-Enabled API** — Ready for cross-origin consumption from the Netlify frontend

---

## 🗂 Project Structure

```
PulseKeep/
├── cmd/pulsekeep/         # Main entrypoint
├── internal/
│   ├── api/               # Gin web server (stats + health endpoints)
│   ├── bot/               # Discord bot + event handlers
│   │   └── commands/      # Slash command definitions & handlers
│   ├── cache/             # In-memory thread-safe cache
│   ├── config/            # Environment config loader
│   └── db/                # PostgreSQL connection via pgx
├── db/
│   └── schema.ts          # Drizzle ORM schema
├── web/                   # Netlify landing page
│   ├── index.html
│   ├── style.css
│   ├── app.js
│   └── assets/
├── Dockerfile             # Multi-stage Go build
├── fly.toml               # Fly.io deployment (recreate strategy)
├── netlify.toml           # Netlify publish config
├── server.md              # Discord support server blueprint
└── go.mod
```

---

## 🚀 Deployment

### Bot → Fly.io

```bash
# 1. Install flyctl
# https://fly.io/docs/hands-on/install-flyctl/

# 2. Launch the app (first time only)
fly launch --no-deploy

# 3. Set secrets (never commit these!)
fly secrets set DISCORD_TOKEN="your-token-here"
fly secrets set DATABASE_URL="postgresql://..."

# 4. Deploy
fly deploy
```

> **Note**: The `fly.toml` is configured with `strategy = "recreate"` to ensure the old instance is fully stopped before the new one starts. This prevents duplicate Discord gateway connections which would cause rate limiting or token suspension.

### Website → Netlify

```bash
# Option A: Netlify CLI
netlify deploy --dir=web --prod

# Option B: Netlify Dashboard
# Connect your GitHub repo and set Publish directory to: web
```

---

## 🛠 Local Development

```bash
# Copy and fill in your secrets
cp .env.example .env

# Run the bot
go run ./cmd/pulsekeep

# Run DB migrations (requires Node.js)
npm install
npm run db:migrate
```

---

## 📋 Environment Variables

| Variable | Description | Required |
|---|---|---|
| `DISCORD_TOKEN` | Your bot token from the Discord Developer Portal | ✅ |
| `DATABASE_URL` | PostgreSQL connection string | ✅ |
| `PORT` | HTTP port for the stats server (default: `8080`) | ❌ |

---

## 📜 License

MIT © 2026 watispro
