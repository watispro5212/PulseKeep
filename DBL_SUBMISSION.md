# PulseKeep DiscordBotList Submission Packet

Use this as the source copy when submitting PulseKeep to DiscordBotList.

## Core Details

- **Bot name:** PulseKeep
- **Website:** https://pulsekeep.fly.dev
- **Command page:** https://pulsekeep.fly.dev/commands.html
- **Status page:** https://pulsekeep.fly.dev/status.html
- **Privacy policy:** https://pulsekeep.fly.dev/privacy.html
- **Terms of service:** https://pulsekeep.fly.dev/terms.html
- **Prefix:** Slash commands
- **Primary command for reviewers:** `/help`
- **Backup test commands:** `/ping`, `/stats`

## Short Description

A Discord bot with moderation, tickets, economy games, and server analytics — all through clean slash commands with a web dashboard.

## Long Description

PulseKeep is a versatile Discord bot designed to help server staff teams manage communities effectively. It combines powerful moderation tools, an intuitive ticket system, a deep economy with minigames, and live server analytics — all accessible through Discord's native slash command interface.

**Moderation** includes warn/kick/ban/timeout, message purging with filters, lockdown, voice management, role management, slowmode, and announcements. Every action is logged for audit purposes.

**Tickets** use a button-based panel system that creates private channels on demand. Staff can add/remove members, close tickets with logging, and rename channels.

**Economy** features daily/weekly rewards, work shifts, gambling with configurable odds, blackjack against an AI dealer, slot machines, fishing and mining minigames, a full item shop, rob/pay interactions, and a search command. Users can vote on DiscordBotList to earn bonus Pulses.

**Configuration** is handled entirely through the `/configure` command with 9 subcommands, or via the web dashboard with Discord OAuth2 login.

**Live Monitoring** includes a public status page with real-time metrics, health checks, and graceful shutdown handling.

PulseKeep is built with TypeScript on discord.js v14, backed by PostgreSQL, and hosted on Fly.io with 99%+ uptime.

## Full Command List

### Moderation (19)
`/ban` `/kick` `/mute` `/unmute` `/warn` `/warnings` `/clearwarns` `/history` `/softban` `/purge` `/clean` `/lock` `/unlock` `/slowmode` `/move` `/vckick` `/nick` `/role add` `/role remove` `/announce`

### Economy (19)
`/balance` `/daily` `/weekly` `/work` `/gamble` `/blackjack` `/slots` `/rob` `/pay` `/search` `/fish` `/mine` `/shop` `/buy` `/use` `/inventory` `/leaderboard` `/tip` `/vote`

### Tickets (5)
`/ticketpanel` `/ticket add` `/ticket remove` `/ticket close` `/ticket rename`

### Configuration (9)
`/configure economy` `/configure tickets` `/configure modlogs` `/configure welcome` `/configure log_channel` `/configure welcome_channel` `/configure vote_channel` `/configure ticket_category` `/configure show`

### Utility (7)
`/help` `/ping` `/stats` `/invite` `/vote` `/servericon` `/roleinfo` `/channelinfo`

**Total: 60+ slash commands**

## Categories

- Moderation
- Economy
- Tickets
- Configuration
- Utility

## Webhook Setup

Set your DBL webhook in DiscordBotList dashboard → Bot Settings → Webhooks:
- **URL:** `https://pulsekeep.fly.dev/api/dbl/webhook`
- **Authorization header:** Your chosen secret (must match `DBL_WEBHOOK_SECRET` env var)

When a user votes, PulseKeep automatically:
1. Credits 500-750 Pulses to the voter
2. Posts a celebratory embed in all shared guilds with a `voteChannelId` configured
