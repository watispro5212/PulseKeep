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

Discord bot with moderation, tickets, economy, and server stats — all slash commands, zero prefixes.

## Long Description

Does the stuff you'd expect: moderation, tickets, economy with mini-games, and a status page. All through slash commands so nobody has to remember weird prefixes.

**Moderation** — warn, kick, ban, timeout, purge with filters, lockdown, voice management, role management, slowmode, announcements. Everything gets logged.

**Tickets** — button panels that create private channels. Add or remove people, close with logging, rename channels.

**Economy** — daily/weekly rewards, work shifts, gamble, blackjack, slots, fishing, mining, item shop, rob/pay, search. Vote on DiscordBotList to earn bonus coins.

**Configuration** — `/configure` with 16 subcommands. Toggle features, set channels, configure auto-mod. All per-server.

**Live Monitoring** — public status page that updates every 30 seconds. Health checks, graceful shutdown.

Built with TypeScript + discord.js v14, PostgreSQL, hosted on Fly.io.

## Full Command List

### Moderation (20)
`/ban` `/kick` `/mute` `/unmute` `/warn` `/warnings` `/clearwarns` `/history` `/softban` `/purge` `/clean` `/lock` `/unlock` `/slowmode` `/move` `/vckick` `/nick` `/role` `/announce` `/data-deletion`

### Economy (18)
`/balance` `/daily` `/weekly` `/work` `/gamble` `/blackjack` `/slots` `/rob` `/pay` `/search` `/fish` `/mine` `/shop` `/buy` `/use` `/inventory` `/leaderboard` `/vote`

### Tickets (5)
`/ticketpanel` `/ticket add` `/ticket remove` `/ticket close` `/ticket rename`

### Configuration (16)
`/configure economy` `/configure tickets` `/configure modlogs` `/configure welcome` `/configure automod` `/configure automod_spam` `/configure automod_mentions` `/configure automod_caps` `/configure automod_links` `/configure automod_words` `/configure automod_banned_words` `/configure log_channel` `/configure welcome_channel` `/configure vote_channel` `/configure ticket_category` `/configure show`

### Utility (11)
`/help` `/ping` `/stats` `/invite` `/about` `/userinfo` `/serverinfo` `/servericon` `/roleinfo` `/channelinfo`

**Total: 70 slash commands**

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
