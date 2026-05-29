# PulseKeep

**A modern, Go-powered Discord bot built for staff teams and community engagement.**

PulseKeep brings together **utility tools, moderation controls, economy games, and an interactive ticket system** — all in one fast, reliable package. Built with Go and PostgreSQL, it handles servers of every size without breaking a sweat.

---

## Quick Links

- **Website:** https://pulsekeep.xyz
- **Support Server:** https://discord.gg/Y4uzWDyaxF
- **Dashboard:** https://pulsekeep.xyz/dashboard
- **Invite Bot:** https://pulsekeep.xyz/invite
- **Changelog:** https://pulsekeep.xyz/changelog

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

Each moderation command checks the user's permissions first and gives clear, specific error messages when something is missing — not raw API errors.

### Economy
- **`/daily`** — Streak-based daily reward (24h cooldown)
- **`/weekly`** — Weekly bonus (7-day cooldown)
- **`/work`** — Random job payouts (45m cooldown)
- **`/balance`** — Check wallet balance
- **`/profile`** — Full economy stats (earned, spent, streaks, games played)
- **`/pay`** — Send Pulses to another member
- **`/shop`** — Browse purchasable items
- **`/buy`** — Purchase an item
- **`/sell`** — Sell an item for 60% refund
- **`/use`** — Use a consumable item
- **`/inventory`** — View owned items
- **`/gift`** — Give an item to another member
- **`/coinflip`** — 50/50 wager on heads or tails
- **`/slots`** — 3-reel slot machine (up to 10x)
- **`/gamble`** — Roll 1–100 (60+ wins, up to 10x at 100)
- **`/fish`** — Cast a line (requires Fishing Rod, 13 species, 7 rarities)
- **`/mine`** — Mine for ores (requires Iron Pickaxe, 13 ores, 7 rarities)
- **`/rob`** — Attempt to steal Pulses (40% base rate, 4h cooldown)
- **`/leaderboard`** — Top 10 by balance
- **`/rich`** — Top 10 with rank badges

Passive 0.1% interest is applied every 6 hours.

### Tickets
- **`/ticketpanel`** — Post the interactive ticket opener
- **Open Ticket button** — Creates a private channel for 1-on-1 support
- **Close Ticket button** — Closes and auto-deletes the ticket channel after 5 seconds

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

PulseKeep runs on **Fly.io** with a **Cloudflare Worker** as a reverse proxy.

### Quick Deploy

1. Clone the repo:
   ```bash
   git clone https://github.com/watispro5212/PulseKeep
   cd PulseKeep
   ```
2. Set environment variables:
   - `BOT_TOKEN` — Discord bot token
   - `DATABASE_URL` — PostgreSQL connection string (Neon recommended)
   - `CLIENT_ID` — Discord OAuth2 client ID
   - `CLIENT_SECRET` — Discord OAuth2 client secret
   - `SESSION_SECRET` — Random string for session signing
3. Deploy the backend:
   ```bash
   fly deploy
   ```
4. Deploy the worker:
   ```bash
   npm install
   npm run deploy
   ```
5. Set the Cloudflare Worker as your OAuth2 redirect URI:
   ```
   https://<your-worker>.workers.dev/auth/discord/callback
   ```

### Required Bot Permissions

- **Administrator** (recommended for full functionality or)
- **Manage Channels** — Lock/unlock, slowmode, ticket creation
- **Manage Messages** — Purge, announce, poll
- **Manage Roles** — /role command
- **Manage Nicknames** — /nick command
- **Kick Members** — /kick
- **Ban Members** — /ban, /unban
- **Moderate Members** — /timeout
- **Read Message History** — Required for purge and userinfo
- **Use Slash Commands** — Required for all commands
- **Send Messages, Embed Links, Attach Files** — Required for command responses

### Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.24 |
| Discord library | disgo |
| Database | PostgreSQL (Neon) |
| Hosting | Fly.io |
| Proxy/Cache | Cloudflare Worker |
| Frontend | HTML + CSS + vanilla JS |

---

## Server Setup Guide

### Recommended Channel Layout

```
📢 INFORMATION
  #welcome              — New member landing
  #rules-and-info       — Rules and bot invite link
  #announcements        — Product updates and releases
  #status-logs          — Automated bot status posts

💬 COMMUNITY
  #general-chat         — Main community discussion
  #bot-discussion       — PulseKeep questions and feature ideas
  #showcase             — Server setups and configurations

⚙️ COMMAND CENTER
  #command-menu         — Pinned /help menu
  #bot-sandbox          — Public slash command testing
  #moderation-lab       — Staff-only moderation testing

💰 ECONOMY ZONE
  #economy-chat         — Economy commands (slowmode: 5s)
  #economy-leaderboard  — Pinned leaderboard display

🎫 CLIENT SUPPORT
  #support-faq          — Common setup and troubleshooting answers
  #open-a-ticket        — Permanent ticket panel (/ticketpanel)
  #pre-sales            — Premium and custom-work questions

📋 ACTIVE TICKETS
  #ticket-0001+         — Auto-created by /ticketpanel

🔒 STAFF OPERATIONS
  #staff-chat           — Staff coordination
  #mod-logs             — Moderation and audit events
  #ticket-archives      — Closed ticket transcripts
  #deploy-logs          — Deployment and incident notes
```

### Role Hierarchy (highest to lowest)

| Role | Color | Permissions |
|------|-------|-------------|
| Founder | `#E74C3C` | Administrator |
| Administrator | `#E67E22` | Manage Server, Manage Channels, Manage Roles, Ban/Kick Members, Moderate Members, View Audit Log |
| Moderator | `#2ECC71` | Manage Messages, Moderate Members, Kick/Ban Members, View Audit Log |
| Support Team | `#3498DB` | Manage Messages, Send Messages, Attach Files, Embed Links, Read Message History |
| PulseKeep Bot | `#9B59B6` | Send Messages, Embed Links, Attach Files, Read Message History, Manage Channels (tickets) |
| Server Booster | `#FD79A8` | Attach Files, Embed Links, External Emojis |
| Verified Member | `#95A5A6` | Send Messages, Read History, Use Slash Commands |
| @everyone | Default | View welcome area only |

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

**Q: `err.Error()` is showing in a command response.**
- This should not happen in PulseKeep v5.9.1+. All user-facing errors use clean messages. If you still see raw error text, report it in the support server.

---

## Links

- **Website:** https://pulsekeep.xyz
- **Support Server:** https://discord.gg/Y4uzWDyaxF
- **Dashboard:** https://pulsekeep.xyz/dashboard
- **Invite Bot:** https://pulsekeep.xyz/invite
- **Changelog:** https://pulsekeep.xyz/changelog
- **GitHub:** https://github.com/watispro5212/PulseKeep
- **Top.gg:** https://top.gg/bot/<bot-id>
