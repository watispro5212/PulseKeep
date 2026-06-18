# PulseKeep

TypeScript Discord bot. Moderation, tickets, economy games, server stats — all slash commands.

## Features

- **Moderation** — warn, kick, ban, mute, purge, lockdown, voice/role management, announcements
- **Tickets** — button panels that make private channels. Add/remove people, close with logging
- **Economy** — daily streaks, blackjack, slots, fishing, mining, item shop, leaderboard, rob/pay
- **Auto-Moderation** — spam detection, mass mentions, banned words, link blocking, caps
- **Configuration** — toggle features per server via `/configure`
- **Voting** — earn coins by voting on DiscordBotList, webhook integration
- **Live Monitoring** — health checks, real-time stats API, public status page

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
fly secrets set ALLOWED_ORIGINS="https://pulsekeep.fly.dev"
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

## Commands

70+ slash commands across 5 categories:

- **Moderation** — warn, mute, kick, ban, softban, purge, clean, slowmode, lock, unlock, nick, role, move, vckick, announce, lockdown, history, warnings, clearwarns, data-deletion
- **Economy** — balance, daily, weekly, work, gamble, blackjack, slots, rob, pay, fish, mine, search, shop, buy, inventory, use, vote, leaderboard
- **Tickets** — ticketpanel, ticket add/remove/close/rename
- **Configuration** — configure (economy, tickets, modlogs, welcome, automod, automod_spam, automod_mentions, automod_caps, automod_links, automod_words, automod_banned_words, log_channel, welcome_channel, vote_channel, ticket_category, show)
- **Utility** — help, ping, stats, invite, about, userinfo, serverinfo, servericon, roleinfo, channelinfo

## Links

- Website: https://pulsekeep.fly.dev
- Commands: https://pulsekeep.fly.dev/commands.html
- Status: https://pulsekeep.fly.dev/status.html
- DiscordBotList: https://discordbotlist.com/bots/1507498795569512598
- GitHub: https://github.com/watispro5212/PulseKeep
