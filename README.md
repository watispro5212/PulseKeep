# PulseKeep

A TypeScript Discord bot for moderation, audit logging, support tickets, economy games, and server analytics — all through clean slash commands.

## Features

- **Moderation Suite** — Warn, kick, ban, timeout, purge, lockdown, voice management, role management, announcements
- **Support Tickets** — Button-based ticket panels create private channels on demand. Add/remove members, close with logging, rename channels
- **Economy System** — Daily streaks, blackjack, slots, gambling, fishing, mining, item shop, leaderboard, rob/pay, voting rewards
- **Auto-Moderation** — Spam detection, mass mention protection, banned words, link blocking, caps enforcement
- **Configuration** — Toggle features per guild, set log/welcome/vote channels via `/configure` or web dashboard
- **DiscordBotList Voting** — Users earn Pulses for voting; webhook integration with configurable announcement channels
- **Live Monitoring** — Health checks, real-time stats API, public status page

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | TypeScript (Node.js) |
| Discord library | discord.js v14 |
| HTTP framework | Express |
| Database | PostgreSQL via Drizzle ORM |
| Hosting | Fly.io (iad, 256 MB, shared-cpu-1x) |
| Frontend | HTML + CSS + vanilla JS |

## Project Layout

```
PulseKeep/
  src/
    api/          Express API server (stats, OAuth, guild config, DBL webhook)
    bot/
      commands/   Slash command implementations
      client.ts   Bot class with gateway events
      types.ts    Discord.js type extensions
    cache/        In-memory stats cache
    config.ts     Environment variable loader
    db/
      schema.ts   Drizzle schema definitions
      index.ts    Database pool
  web/            Static site (HTML, CSS, JS)
  scripts/        Utility scripts (DB migrations)
  fly.toml        Fly.io configuration
  Dockerfile      Container build
```

## Local Development

```bash
cp .env.example .env
# Fill in DISCORD_TOKEN and DATABASE_URL
npm install
npm run dev
```

For API-only mode without connecting to Discord:
```
BOT_DISABLED=true npm run dev
```

## Fly.io Deployment

```bash
fly launch --no-deploy
fly secrets set DISCORD_TOKEN="your-token"
fly secrets set DATABASE_URL="postgresql://..."
fly secrets set ALLOWED_ORIGIN="https://pulsekeep.fly.dev"
fly secrets set DISCORD_CLIENT_ID="..."
fly secrets set DISCORD_CLIENT_SECRET="..."
fly secrets set DISCORD_REDIRECT_URI="https://pulsekeep.fly.dev/auth/discord/callback"
fly secrets set DBL_API_TOKEN="your-dbl-api-token"
fly secrets set DBL_WEBHOOK_SECRET="your-webhook-secret"
fly deploy
```

## Environment Variables

| Variable | Required | Description |
| --- | --- | --- |
| `DISCORD_TOKEN` | For full bot mode | Discord bot token |
| `DATABASE_URL` | For database features | PostgreSQL connection string |
| `PORT` | No | HTTP port, defaults to 8080 |
| `ALLOWED_ORIGINS` | No | CORS origins (comma-separated) |
| `BOT_DISABLED` | No | Set `true` to run API-only |
| `DISCORD_CLIENT_ID` | For dashboard | Discord OAuth2 client ID |
| `DISCORD_CLIENT_SECRET` | For dashboard | Discord OAuth2 client secret |
| `DISCORD_REDIRECT_URI` | For dashboard | OAuth callback URL |
| `STATUS_WEBHOOK_URL` | No | Discord webhook for status alerts |
| `DBL_API_TOKEN` | For vote command | DiscordBotList API token |
| `DBL_WEBHOOK_SECRET` | For DBL webhook | Secret for vote webhook auth |

## Database

The project uses Drizzle for schema management:

```bash
npm run db:generate
npm run db:migrate
```

## Commands

50+ slash commands across 5 categories:

- **Moderation** — ban, kick, mute, warn, purge, lock, slowmode, role, nick, move, vckick, announce, clean, softban
- **Economy** — balance, daily, weekly, work, gamble, blackjack, slots, rob, pay, shop, buy, inventory, use, search, fish, mine, leaderboard, tip, vote
- **Tickets** — ticketpanel, ticket add/remove/close/rename
- **Configuration** — configure (economy, tickets, modlogs, welcome, log_channel, welcome_channel, vote_channel, ticket_category, show)
- **Utility** — help, ping, stats, invite, vote, servericon, roleinfo, channelinfo

## Links

- Website: https://pulsekeep.fly.dev
- Commands: https://pulsekeep.fly.dev/commands.html
- Dashboard: https://pulsekeep.fly.dev/dashboard.html
- Status: https://pulsekeep.fly.dev/status.html
- DiscordBotList: https://discordbotlist.com/bots/1507498795569512598
- GitHub: https://github.com/watispro5212/PulseKeep
