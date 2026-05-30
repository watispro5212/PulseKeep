# PulseKeep

**A modern, Go-powered Discord bot built for staff teams and community engagement.**

PulseKeep brings together utility tools, moderation controls, economy games, and an interactive ticket system — all in one fast, reliable package. Built with Go 1.26 and PostgreSQL, it handles servers of every size without breaking a sweat.

---

## Quick Links

- **Website:** https://pulsekeep.williamdelilah3.workers.dev
- **Commands:** https://pulsekeep.williamdelilah3.workers.dev/commands.html
- **Status:** https://pulsekeep.williamdelilah3.workers.dev/status.html
- **Dashboard:** https://pulsekeep.williamdelilah3.workers.dev/dashboard.html
- **Changelog:** https://pulsekeep.williamdelilah3.workers.dev/changelog.html
- **Privacy:** https://pulsekeep.williamdelilah3.workers.dev/privacy.html
- **Invite:** https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands
- **GitHub:** https://github.com/watispro5212/PulseKeep

---

## Features

### Utility
- **`/help`** — Interactive command browser with category dropdown
- **`/ping`** — Check gateway connectivity
- **`/stats`** — Bot runtime snapshot and metrics
- **`/uptime`** — Time since last restart
- **`/about`** — Version, tech stack, and links
- **`/serverinfo`** — Current guild details (owner, members, boosts, roles, created)
- **`/userinfo`** — Member lookup with join dates and role count
- **`/avatar`** — Full-resolution avatar display
- **`/poll`** — Reaction-based polls with 2–4 options

### Moderation
- **`/purge`** — Bulk delete 1–100 messages
- **`/kick`** — Remove a member with optional audit-log reason
- **`/ban`** — Ban a member with optional audit-log reason
- **`/unban`** — Unban a user by ID
- **`/timeout`** — Timeout a member for 1–40320 minutes
- **`/nick`** — Change or reset a member's nickname
- **`/slowmode`** — Set channel slowmode (0–21600s)
- **`/lock`** — Lock channel for @everyone
- **`/unlock`** — Unlock channel for @everyone
- **`/announce`** — Send branded embedded announcements with optional @everyone ping
- **`/role`** — Add or remove a role from a member

Each moderation command checks the user's permissions and the bot's permissions separately, giving clear error messages when something is missing — no raw API errors exposed to users.

### Economy
- **`/daily`** — Streak-based daily reward (24h cooldown)
- **`/weekly`** — Weekly bonus (7-day cooldown)
- **`/work`** — Random job payouts (45m cooldown)
- **`/balance`** — Check wallet balance
- **`/profile`** — Full economy stats (earned, spent, streaks, games played)
- **`/pay`** — Send Pulses to another member
- **`/shop`** — Browse 10 purchasable items
- **`/buy`** — Purchase an item
- **`/sell`** — Sell an item for 60% refund
- **`/use`** — Use a consumable item
- **`/inventory`** — View owned items
- **`/gift`** — Give an item to another member
- **`/coinflip`** — 50/50 wager on heads or tails
- **`/slots`** — 3-reel slot machine (up to 10x)
- **`/gamble`** — Roll 1–100 (40+ pushes, up to 10x at 100)
- **`/blackjack`** — Interactive blackjack with hit/stand buttons against a CPU dealer
- **`/fish`** — Cast a line (requires Fishing Rod, 13 species, 7 rarities)
- **`/mine`** — Mine for ores (requires Iron Pickaxe, 13 ores, 7 rarities)
- **`/rob`** — Attempt to steal Pulses (40% base rate, 4h cooldown)
- **`/rich`** — Top 10 by balance with rank badges
- **`/lottery`** — Check the weekly lottery jackpot
- **`/lottery-claim`** — Claim your prize if you won the weekly draw

Passive 0.1% interest applied every 6 hours. All gambling has a built-in 5% house edge.

### Tickets
- **`/ticketpanel`** — Posts the interactive ticket opener
- **Open Ticket button** — Creates a private channel for 1-on-1 support
- **Close Ticket button** — Closes and auto-deletes the channel after 5 seconds

### Auto-moderation
- **Spam detection** — Repeated messages in quick succession
- **Mass mention** — 5+ unique mentions per message
- **Banned words** — Configurable word list per server
- **Link spam** — Blocks messages with 3+ links
- **All-caps abuse** — Messages >80% uppercase (min 15 chars)
- **Configurable actions** — Delete only, warn, or timeout
- **Log channel** — Optional channel for auto-mod action logs

---

## Self-hosting

### Prerequisites

- Go 1.26+
- Docker (for Fly.io deployment)
- Discord bot token from the [Discord Developer Portal](https://discord.com/developers/applications)
- (Optional) Neon PostgreSQL database

### Quick Deploy

1. Clone the repo:
   ```bash
   git clone https://github.com/watispro5212/PulseKeep
   cd PulseKeep
   ```
2. Set environment variables:
   ```bash
   fly secrets set DISCORD_TOKEN="your-token-here"
   fly secrets set DATABASE_URL="postgresql://..."
   fly secrets set ALLOWED_ORIGIN="https://your-cloudflare-site.pages.dev"
   fly secrets set DISCORD_CLIENT_ID="..."
   fly secrets set DISCORD_CLIENT_SECRET="..."
   fly secrets set DISCORD_REDIRECT_URI="https://your-site.com/auth/discord/callback"
   ```
3. Deploy:
   ```bash
   fly deploy
   ```
4. Deploy the website to Cloudflare Pages:
   ```bash
   npx wrangler pages deploy web --branch main
   ```

### Required Bot Permissions

- **Administrator** is recommended for full functionality, or select:
- Manage Channels — lock/unlock, slowmode, ticket creation
- Manage Messages — purge, announce, poll
- Manage Roles — /role
- Manage Nicknames — /nick
- Kick Members — /kick
- Ban Members — /ban, /unban
- Moderate Members — /timeout
- Read Message History — purge, userinfo
- Use Slash Commands — required for all commands
- Send Messages, Embed Links, Attach Files — command responses

### Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.26 |
| Discord library | disgo v0.19.3 |
| HTTP router | Gin v1.12.0 |
| Database driver | pgx v5 |
| Database | PostgreSQL (Neon) |
| Hosting | Fly.io (iad region, 256 MB) |
| Website | Cloudflare Pages |
| Frontend | HTML + CSS + vanilla JS |

---

## Common Issues

**Q: Slash commands aren't appearing.**
- Make sure the bot has `applications.commands` scope when invited.
- Re-invite the bot if needed.
- Discord can take up to an hour to sync global commands in large servers.

**Q: The bot says I'm missing permissions, but I have them.**
- Check that your role is above the bot's role in Server Settings > Roles.
- Some commands check both your permissions and the bot's permissions separately.

**Q: Ticket panel does nothing when clicked.**
- Give the bot `Manage Channels` permission.
- Make sure the bot role is above the Support Team role in the hierarchy.

**Q: Economy balance isn't updating.**
- Check that the bot can see the channel and use slash commands.
- Verify the PostgreSQL database is reachable.
- Each economy command has a cooldown — use `/profile` to see your actual balance.

**Q: The status page shows "Service offline".**
- The website is on Cloudflare Pages and fetches from the Fly.io API. If the Fly.io app is down or the Cloudflare-to-Fly.io CORS policy isn't configured, the status page will show offline. Check `ALLOWED_ORIGIN` on Fly.io.

---

## Links

- **Website:** https://pulsekeep.williamdelilah3.workers.dev
- **Commands:** https://pulsekeep.williamdelilah3.workers.dev/commands.html
- **Status:** https://pulsekeep.williamdelilah3.workers.dev/status.html
- **Dashboard:** https://pulsekeep.williamdelilah3.workers.dev/dashboard.html
- **Changelog:** https://pulsekeep.williamdelilah3.workers.dev/changelog.html
- **GitHub:** https://github.com/watispro5212/PulseKeep
