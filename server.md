# PulseKeep Discord Server Blueprint

This blueprint is the recommended layout for a PulseKeep support, testing, and community server. It uses category-first permissions, a clear support workflow, dedicated command channels, and isolated sandbox areas for safe bot testing.

---

## Table of Contents

- [Core Goals](#core-goals)
- [Role Hierarchy](#role-hierarchy)
- [Permission Strategy](#permission-strategy)
- [Channel Layout](#channel-layout)
  - [1. Information](#1-information)
  - [2. Community](#2-community)
  - [3. Command Center](#3-command-center)
  - [4. Economy Zone](#4-economy-zone)
  - [5. Client Support](#5-client-support)
  - [6. Active Tickets](#6-active-tickets)
  - [7. Staff Operations](#7-staff-operations)
- [Interactive Command Menu](#interactive-command-menu)
- [Economy Setup Guide](#economy-setup-guide)
- [Lock / Unlock Workflow](#lock--unlock-workflow)
- [Slowmode Best Practices](#slowmode-best-practices)
- [Webhook Recommendations](#webhook-recommendations)
- [Permission Checklist](#permission-checklist)
- [Launch Checklist](#launch-checklist)
- [FAQ](#faq)

---

## Core Goals

- Keep public channels readable and low-noise.
- Give support staff one obvious ticket workflow.
- Keep bot testing contained to sandbox channels.
- Use category permissions first, with channel overrides only when truly needed.
- Make `/help` and the ticket panel the main interaction points.
- Keep economy commands in a dedicated channel to avoid spam in general chat.

---

## Role Hierarchy

Create roles in this order from highest to lowest. Higher roles inherit permissions from lower roles — place roles carefully.

| Role | Color | Purpose | Key Permissions |
| --- | --- | --- | --- |
| **Founder** | `#E74C3C` | Project owner and final authority. | Administrator |
| **Administrator** | `#E67E22` | Senior operations and server configuration. | Manage Server, Manage Channels, Manage Roles, View Audit Log, Ban Members, Kick Members, Moderate Members |
| **Moderator** | `#2ECC71` | Community safety and chat moderation. | Manage Messages, Moderate Members, Kick Members, Ban Members, View Audit Log |
| **Support Team** | `#3498DB` | Handles tickets, setup help, and customer questions. | Manage Messages, Send Messages, Attach Files, Embed Links, Read Message History |
| **PulseKeep Bot** | `#9B59B6` | The bot account for all slash commands and ticket flows. | Manage Messages, Send Messages, Embed Links, Attach Files, Read Message History, Use Slash Commands, Manage Channels (for ticket creation) |
| **Server Booster** | `#FD79A8` | Trusted supporters with small perks. | Attach Files, Embed Links, Use External Emojis |
| **Verified Member** | `#95A5A6` | Standard trusted community member. | View Channels, Send Messages, Read Message History, Use Slash Commands |
| **@everyone** | Default | Pre-verification visitors. | View only the welcome/rules area |

**Note:** The PulseKeep Bot role should be placed above Support Team in the role list so the bot can manage ticket channels that support members use.

---

## Permission Strategy

- **Category permissions first** — set base access at the category level, not per-channel.
- **Channel overrides only when necessary** — exceptions for things like read-only announcement channels.
- **No `@everyone` access to staff or ticket categories** — use role-based access only.
- **Least privilege for PulseKeep Bot** — give it only the permissions it needs. Do **not** give it Administrator.

---

## Channel Layout

### 1. Information

Read-only server information for all members.

| Channel | Purpose |
| --- | --- |
| `#rules-and-info` | Rules, useful links, bot invite, and support expectations |
| `#announcements` | Product updates, releases, outages, and important changes |
| `#status-logs` | Automated deploy, uptime, and incident notices |
| `#welcome` | New member greeting and first steps |

**Permissions:**
- Staff can view, send, and manage messages.
- PulseKeep Bot can view, send, and embed links.
- Verified Member and @everyone can view and read history.
- Verified Member and @everyone **cannot** send messages.

**Recommended pinned message for `#rules-and-info`:**
```text
Welcome to PulseKeep.

Start here:
1. Read the rules.
2. Use /help to browse commands.
3. Use the ticket panel in #open-a-ticket when you need private setup help.

Do not post bot tokens, database URLs, or private server logs in public channels.
```

### 2. Community

Normal conversation area for members.

| Channel | Purpose |
| --- | --- |
| `#general-chat` | Main community chat |
| `#bot-discussion` | Usage questions and feature ideas |
| `#showcase` | Server setups, panels, and PulseKeep configurations |

**Permissions:**
- Verified Member can view, send messages, use slash commands, react, and read history.
- Server Booster can also attach files and embed links.
- Staff can manage messages.
- @everyone cannot view this category.

**Recommended settings:**
- Slowmode: 3 seconds in `#general-chat`.
- Disable link embeds for normal members unless you trust the community.

### 3. Command Center

The command center is where members discover and test PulseKeep commands. Keep these channels clean — no general chat here.

| Channel | Purpose |
| --- | --- |
| `#command-menu` | The permanent PulseKeep interactive command menu |
| `#bot-sandbox` | General testing for slash commands |
| `#moderation-lab` | Staff-only testing for moderation commands |

**Permissions:**
- Verified Member can view `#command-menu` and `#bot-sandbox`.
- Verified Member can use slash commands in both.
- Moderator and above can view and use `#moderation-lab`.
- PulseKeep Bot can send messages and embed links everywhere in the category.

**Setup:**
1. In `#command-menu`, run `/help`.
2. Pin the bot's interactive command browser message.
3. Tell users to use the dropdown to switch between categories.
4. Keep normal chatting out of this channel.

**Command aliases:**
- `/help`: Opens the interactive command browser privately.
- `!help`: Posts the menu publicly in a text channel.

### 4. Economy Zone

Dedicated space for PulseKeep's economy system — keeps gambling, fishing, and mining out of general chat.

| Channel | Purpose |
| --- | --- |
| `#economy-chat` | Economy commands: `/daily`, `/work`, `/balance`, `/profile`, `/coinflip`, `/leaderboard`, `/pay`, `/rob`, `/shop`, `/buy`, `/inventory`, `/slots`, `/gamble`, `/fish`, `/mine`, `/sell`, `/use` |
| `#economy-leaderboard` | Pinned leaderboard refreshes (manual or automated) |

**Permissions:**
- Verified Member can view, send messages, use slash commands, react, and read history.
- Staff can manage messages.
- Server Booster can attach files and embed links.

**Recommended settings:**
- Slowmode: 5 seconds in `#economy-chat`.
- Pin the leaderboard in `#economy-leaderboard` after running `/leaderboard`.

### 5. Client Support

Support entry points for users needing help.

| Channel | Purpose |
| --- | --- |
| `#support-faq` | Answers to common setup, deploy, and permission issues |
| `#open-a-ticket` | The permanent ticket panel |
| `#pre-sales` | Questions about premium features, setup help, and custom work |

**Permissions:**
- Verified Member can view, send messages, read history, and use slash commands.
- Support Team, Moderator, Administrator, and Founder can manage messages.
- @everyone cannot view this category.
- PulseKeep Bot can send messages, embed links, and attach files.

**Setup:**
1. In `#open-a-ticket`, run `/ticketpanel`.
2. Pin the ticket panel.
3. Keep the channel locked to support questions only.
4. Ask users to include: server name, command or feature, expected behavior, actual behavior, and any safe error text.

**Recommended `#support-faq` topics:**
- Fly.io deploy checklist
- Netlify website checklist
- Missing Discord permissions
- Bot token and secret safety
- Database connection troubleshooting
- Slash commands not appearing
- Lock/unlock not working
- Economy balance not updating

### 6. Active Tickets

Private ticket channels for 1-on-1 support.

| Channel | Purpose |
| --- | --- |
| `#ticket-0001`, `#ticket-0002`, etc. | Generated by the bot when a user opens a ticket |

**Base category permissions:**
- Founder, Administrator, and Support Team can view, send, attach files, embed links, and read history.
- PulseKeep Bot can view, send, embed links, attach files, and manage channels.
- @everyone and Verified Member cannot view.

**Dynamic ticket override:**
- The ticket creator gets: View Channel, Send Messages, Attach Files, Read Message History.
- **Do not** add broad member-role overrides inside ticket channels.
- Archive closed ticket transcripts into `#ticket-archives`.

### 7. Staff Operations

Private internal workspace for staff.

| Channel | Purpose |
| --- | --- |
| `#staff-chat` | Coordination and escalation |
| `#mod-logs` | Moderation actions, audit events, and security notes |
| `#ticket-archives` | Closed ticket transcripts and summaries |
| `#deploy-logs` | Fly.io, Netlify, database migration, and incident notes |

**Permissions:**
- Founder, Administrator, Moderator, and Support Team can view and send.
- Only Administrator and Founder should manage channels and roles.
- PulseKeep Bot can send logs and embeds.
- @everyone, Verified Member, and Server Booster cannot view.

---

## Interactive Command Menu

PulseKeep's interactive menu makes command discovery built-in.

**Entry points:**
- `/help`: Private interactive command menu.
- `!help`: Public menu reply for legacy text-command users.

**Menu categories (40+ commands):**

| Category | Commands |
| --- | --- |
| **Utility** (8) | `/ping`, `/help`, `/about`, `/stats`, `/uptime`, `/serverinfo`, `/userinfo`, `/avatar` |
| **Moderation** (11) | `/purge`, `/kick`, `/ban`, `/unban`, `/timeout`, `/nick`, `/slowmode`, `/lock`, `/unlock`, `/announce`, `/poll` |
| **Economy** (17) | `/balance`, `/profile`, `/daily`, `/work`, `/coinflip`, `/pay`, `/leaderboard`, `/rob`, `/shop`, `/buy`, `/inventory`, `/slots`, `/gamble`, `/fish`, `/mine`, `/sell`, `/use` |
| **Tickets** (1+button) | `/ticketpanel`, Open Ticket button |

---

## Economy Setup Guide

PulseKeep's economy system runs in-memory with PostgreSQL persistence. All members start with 0 Pulses.

### Getting Started for Members
1. `/daily` — Claim your first daily reward.
2. `/work` — Work a shift (45-min cooldown).
3. `/balance` — Check your wallet.
4. `/profile` — See your full economy stats.

### Earning Pulses
- **Daily rewards** — Streak-based, resets every 24 hours.
- **Work** — Random job payouts with 45-min cooldown.
- **Coinflip** — 50/50 chance to double your wager.
- **Slots** — Match 3 symbols for up to 10x payout.
- **Gamble** — Roll 60+ to win (85+ = 2x, 95+ = 4x, 100 = 10x).
- **Fishing** — Requires `Fishing Rod` from shop. 13 fish types across 7 rarities.
- **Mining** — Requires `Iron Pickaxe` from shop. 13 ore types across 7 rarities.
- **Interest** — 0.1% passive balance growth every 6 hours.

### Spending Pulses
- **Shop** — Browse with `/shop`, buy with `/buy`.
- **Sell** — 60% refund with `/sell`.
- **Use** — Consumable items with `/use`.
- **Pay** — Send Pulses to other members with `/pay`.

### Shop Items

| Item | Price | Effect |
| --- | --- | --- |
| Fishing Rod | 1,500 Pulses | Required for `/fish` |
| Iron Pickaxe | 2,000 Pulses | Required for `/mine` |
| XP Boost | 3,000 Pulses | 2x work earnings for 1 hour |
| Lucky Clover | 4,000 Pulses | +1 slot reel position |
| Lucky Pickaxe | 5,000 Pulses | +15% coinflip win chance |
| Shield Token | 6,000 Pulses | Protects from one robbery (consumable) |
| Golden Watch | 8,000 Pulses | Reduces daily cooldown by 4h |
| Health Potion | 2,500 Pulses | Restores 25% of lost rob fines (consumable) |
| Lottery Ticket | 500 Pulses | Entry for the server lottery draw |
| Treasure Map | 7,000 Pulses | Doubles next fish/mine payout |

### Rob Protection
- Members with a **Shield Token** are immune to robbery.
- The token is consumed on a successful robbery attempt.
- 40% base success rate for robbers.
- Failed robbery costs a fine.
- 4-hour cooldown between robbery attempts.

---

## Lock / Unlock Workflow

The `/lock` and `/unlock` commands control `@everyone` send permissions in a channel.

**How it works:**
- `/lock` — Sets `Deny: SendMessages` on the @everyone permission overwrite. **Preserves all other existing permission overwrites** (role-specific and member-specific).
- `/unlock` — Removes only the `SendMessages` deny from the @everyone overwrite. **Does not touch any other overwrites.**

**Best practices:**
- Use `/lock` to stop chat during raids, announcements, or incidents.
- Use `/unlock` to restore chat access.
- These commands require `Manage Channels` permission.
- Locking a channel does not affect staff — staff roles with `Manage Messages` or `Administrator` can still send messages.

---

## Slowmode Best Practices

| Channel type | Recommended slowmode |
| --- | --- |
| General chat | 3 seconds |
| Economy chat | 5 seconds |
| Bot sandbox | 0 seconds (no limit) |
| Support channels | 5 seconds |
| Announcement discussions | 10 seconds |
| Staff chat | 0 seconds |

The `/slowmode` command accepts values from 0 (disabled) to 21,600 (6 hours). Requires `Manage Channels` permission.

---

## Webhook Recommendations

Create separate webhooks for clean operational messages.

| Channel | Webhook Name | Use |
| --- | --- | --- |
| `#announcements` | PulseKeep News | Releases and major updates |
| `#status-logs` | PulseKeep Status | Uptime, deploys, incidents |
| `#mod-logs` | PulseKeep Audit | Moderation and audit events |
| `#deploy-logs` | PulseKeep Deploy | Fly.io, Netlify, and database deployment notes |

**Do not** post secrets, full database URLs, bot tokens, or private user data through webhooks.

---

## Permission Checklist

Before launching the server, confirm:

- [ ] Members cannot see staff channels.
- [ ] Members cannot see active tickets unless they own that ticket.
- [ ] Members can use slash commands in command and sandbox channels.
- [ ] PulseKeep Bot has Send Messages, Embed Links, Attach Files, Read Message History, and Use Slash Commands.
- [ ] PulseKeep Bot has Manage Channels **only** if automated ticket channel creation is enabled.
- [ ] Staff moderation commands are protected by Discord permissions.
- [ ] `/lock` and `/unlock` work correctly in text channels.
- [ ] Commands that require permissions (ban, kick, timeout) are restricted by Discord's built-in permission system.
- [ ] Economy commands are accessible in the economy channel.

---

## Launch Checklist

1. Create roles and categories.
2. Apply category permissions.
3. Create channels.
4. Invite PulseKeep with required scopes: `bot` and `applications.commands`.
5. Run `/help` in `#command-menu`.
6. Run `/ticketpanel` in `#open-a-ticket`.
7. Pin the generated bot messages.
8. Test `/help`, `!help`, and the ticket button.
9. Test member visibility with a non-staff account.
10. Publish the website and link it in `#rules-and-info`.

---

## FAQ

**Q: Why can't I see the /help menu?**
A: Make sure PulseKeep has `Use Slash Commands` permission in the channel.

**Q: The ticket panel isn't creating channels.**
A: Check that PulseKeep has `Manage Channels` permission and that the bot role is above the Support Team role.

**Q: /lock says it failed.**
A: Make sure PulseKeep has `Manage Channels` permission and that you're in a text channel.

**Q: /slowmode isn't working.**
A: The bot needs `Manage Channels` permission. Slowmode only works in text channels.

**Q: Economy commands are not showing up.**
A: Make sure the bot has `Use Slash Commands` permission. If commands were just added, wait a few minutes for Discord to sync.

**Q: The leaderboard shows no entries.**
A: Members need to use at least one economy command (`/daily` or `/work`) to appear on the leaderboard.

**Q: How do I reset a user's economy data?**
A: Economy data is stored in PostgreSQL. Contact the bot owner (watispro1 on Discord) for manual resets.
