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

## What's New — v7.6.0

Stuff that changed recently:

- Shop command tells you when the DB is down instead of making up a fake balance
- API rate limiter won't leak memory anymore — old IPs get cleaned out every 2 minutes
- Auto-mod spam cleaner won't stack duplicate timers if the module reloads
- Rob cooldown now applies even when you fail with zero balance (no more infinite retries)
- HSTS header added, robots.txt so search engines find the sitemap
- Status page retries API calls with backoff (2s, 4s) instead of 3 instant spams
- Mobile hero panel no longer puts stats above the headline
- OG tags added to changelog, privacy, terms, and 404 pages so link previews work
- `prefers-reduced-motion` respected for accessibility

## What's New — v7.5.0

The previous batch:

- Blackjack went from auto-play to actual interactive Hit/Stand buttons
- Cooldown display fixed for fish/mine — shows seconds when under a minute
- Every economy command wrapped in DB transactions to kill race conditions
- Pay, rob, and gamble all re-fetch balances inside transactions now
- Migration scripts cleaned up — no more hardcoded DB passwords
- Empty catch blocks across API + bot code now log errors
- Pay and use respect the economy toggle now
- Version unified to v7.0.0 across the board
- Drizzle migration generated — adds 13 missing columns to the database (automod toggles, channel IDs, xp boost fields)
- npm run db:migrate and db:generate scripts added
- /about cleaned up — removed an unused variable that was doing nothing
- Help command crash fixed — was using a bad non-null assertion that blew up when no commands matched
- Team page redesigned with proper cards and social links
- Privacy + terms pages got staggered scroll animations

## Links

- Website: https://pulsekeep.fly.dev
- Commands: https://pulsekeep.fly.dev/commands.html
- Status: https://pulsekeep.fly.dev/status.html
- DiscordBotList: https://discordbotlist.com/bots/1507498795569512598
- GitHub: https://github.com/watispro5212/PulseKeep
