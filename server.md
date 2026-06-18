# 🛡️ PulseKeep Support Server Blueprint

Owner: `watispro1` · Co‑Owner: `williamdelilah7_`  
Recommended server name: `PulseKeep Support`  
Recommended invite label: `pulsekeep-support`

> [!NOTE]
> Complete Discord support-server build plan. Create the server category by category, role by role, permission by permission.

---

## 📑 Table of Contents

| Section | Description |
|---------|-------------|
| [🎯 Server Purpose](#-server-purpose) | Why this server exists |
| [🧠 Server Style](#-server-style) | Look and feel |
| [👥 Role Stack](#-role-stack) | All roles from highest to lowest |
| [🔐 Role Permissions](#-role-permissions) | Detailed role permission tables |
| [🗂️ Category Layout](#-category-layout) | All 11 categories |
| [📌 START HERE](#-start-here) | Welcome, rules, start-here |
| [📚 PULSEKEEP INFO](#-pulsekeep-info) | Announcements, changelog, status, commands, FAQ |
| [🎫 SUPPORT](#-support) | Support info, help, bug reports, tickets |
| [🤖 BOT COMMANDS](#-bot-commands) | Command menu, status checks, economy, tickets demo |
| [🧪 TEST LAB](#-test-lab) | Slash command, moderation, economy, ticket, dashboard, automod testing |
| [💬 COMMUNITY](#-community) | General, showcase, suggestions, off-topic |
| [🧑‍💻 CONTRIBUTORS](#-contributors) | Contributor chat, docs feedback, translation help |
| [🔒 STAFF](#-staff) | Staff chat, mod chat, support notes, review queue, staff commands |
| [🧾 LOGS](#-logs) | Mod, bot, ticket, vote, automod, dashboard, deploy logs |
| [🚨 INCIDENTS](#-incidents) | Incident response, security reports, audit review |
| [📦 ARCHIVE](#-archive) | Resolved tickets, old announcements, old bugs |
| [🧩 Category Permission Defaults](#-category-permission-defaults) | Default permissions per category |
| [🤖 Bot Invite Settings](#-bot-invite-settings) | Invite links and permissions |
| [🎫 Ticket Workflow](#-ticket-workflow) | Full ticket flow guide |
| [🖥️ Dashboard Workflow](#-dashboard-workflow) | Dashboard login and config guide |
| [⏱️ Economy Cooldowns](#-economy-cooldown-reference) | Command cooldown reference |
| [🎨 Embed Colors](#-embed-color-reference) | Embed color reference |
| [📊 Status Workflow](#-status-workflow) | Status page guide |
| [✅ DBL Checklist](#-discordbotlist-readiness-checklist) | DiscordBotList submission checklist |
| [📣 First Announcement](#-first-announcement) | Server opening announcement |

---

## 🎯 Server Purpose

PulseKeep Support should exist for these core purposes:

- 👑 **Ownership hub** - Clearly shows that PulseKeep is owned by `watispro1` and co‑owned by `williamdelilah7_`.
- 🧰 **User support** - Gives server owners a clean place to ask setup questions.
- 🎫 **Private tickets** - Lets users share server-specific issues without exposing private IDs or screenshots publicly.
- 🤖 **Bot command help** - Explains every major command category and where commands should be tested.
- 🧪 **Safe testing** - Gives testers a controlled place to try moderation, economy, ticket, and dashboard features.
- 📣 **Announcements** - Publishes releases, outages, DiscordBotList news, and major updates.
- 📊 **Status updates** - Tells users whether the bot, website, dashboard, or API is healthy.
- 🔐 **Security handling** - Gives staff a private path for token leaks, abuse reports, and vulnerability reports.
- ✅ **DiscordBotList readiness** - Gives DiscordBotList reviewers a visible support server, command guide, status flow, privacy links, and owner identity.

## 🧠 Server Style

The server should feel clean, serious, and easy to navigate.

Members should immediately know:

- 📌 where to start
- 📜 where to read rules
- ⌨️ where to see commands
- 🎫 where to open a ticket
- 📊 where to check service status
- 👑 who owns PulseKeep

Staff should have:

- 🔒 private coordination channels
- 🧾 audit and action logs
- 🧪 test channels
- 🚨 incident channels
- 📋 review queues

## 👥 Role Stack

Create these roles from highest to lowest. Put the **PulseKeep Bot** role above any role it needs to manage, assign, moderate, kick, ban, timeout, or rename.

| Order | Role Name | Emoji | Purpose | Suggested Color |
| --- | --- | --- | --- | --- |
| 1 | 👑 Owner | 👑 | Full ownership role for `watispro1` only. | PulseKeep blue |
| 2 | 💎 Co-Owner | 💎 | Trusted emergency backup with full server authority. | Royal blue |
| 3 | 🛡️ Administrator | 🛡️ | Manages server structure, permissions, channels, and critical fixes. | Red |
| 4 | 🧑‍💻 Developer | 🧑‍💻 | Handles code, deploys, bot debugging, API issues, and dashboard fixes. | Indigo |
| 5 | 🔐 Security Lead | 🔐 | Handles abuse reports, leaked secrets, exploit reports, and incident response. | Dark red |
| 6 | 🧰 Support Manager | 🧰 | Leads support staff, manages ticket quality, and escalates hard issues. | Cyan |
| 7 | 🎧 Senior Support | 🎧 | Experienced support staff trusted with tickets and user troubleshooting. | Green |
| 8 | 💬 Support Team | 💬 | Normal support helpers for setup questions and tickets. | Light green |
| 9 | ⚖️ Head Moderator | ⚖️ | Leads moderation policy and supervises moderators. | Amber |
| 10 | 🔨 Moderator | 🔨 | Moderates public channels and handles rule breaks. | Orange |
| 11 | 🧹 Trial Moderator | 🧹 | Limited trainee moderator with restricted permissions. | Yellow |
| 12 | 🧪 Bot Tester | 🧪 | Tests new commands, menus, tickets, dashboard changes, and economy flows. | Purple |
| 13 | 🐞 Bug Hunter | 🐞 | Reports reproducible bugs with screenshots, logs, and steps. | Pink |
| 14 | 🧱 Contributor | 🧱 | Helps with documentation, feedback, website copy, or community ideas. | Teal |
| 15 | 🤝 Partner | 🤝 | Partner server representatives and friendly project contacts. | Light blue |
| 16 | ⭐ VIP | ⭐ | Trusted or long-term community member. | Gold |
| 17 | 🙋 Active Helper | 🙋 | Helpful community member with no staff powers. | Green |
| 18 | 👤 Member | 👤 | Normal verified member. | Default |
| 19 | 🔇 Muted | 🔇 | Restricted member who cannot chat or use commands. | Gray |
| Bot | 🤖 PulseKeep Bot | 🤖 | The bot role. Must be high enough to moderate/manage target roles. | PulseKeep blue |

## 🔐 Role Permissions

<details>
<summary><strong>Click to expand — full role permission tables</strong></summary>

### 👑 Owner

Purpose: Complete control of the support server and final authority on PulseKeep operations.

Give:

- ✅ Administrator
- ✅ Manage Server
- ✅ Manage Roles
- ✅ Manage Channels
- ✅ Manage Webhooks
- ✅ View Audit Log
- ✅ Manage Messages
- ✅ Manage Threads
- ✅ Moderate Members
- ✅ Kick Members
- ✅ Ban Members
- ✅ Manage Nicknames
- ✅ Mention Everyone
- ✅ Manage Events

Rules:

> [!WARNING]
> - Only `watispro1` should have this role.
> - Never give this role as a reward.
> - Never give this role to temporary helpers.

### 💎 Co-Owner

Role held by `williamdelilah7_`. Emergency backup if the owner is unavailable.

Give:

- ✅ Administrator
- ✅ Manage Server
- ✅ Manage Roles
- ✅ Manage Channels
- ✅ Manage Webhooks
- ✅ View Audit Log
- ✅ Manage Messages
- ✅ Moderate Members
- ✅ Kick Members
- ✅ Ban Members

Avoid:

> [!CAUTION]
> - ❌ Ownership transfer
> - ❌ Sharing Cloudflare, Discord bot, database, or hosting secrets in Discord

### 🛡️ Administrator

Purpose: Fix server-wide problems and keep categories, permissions, and support flow working.

Give:

- ✅ Administrator
- ✅ Manage Server
- ✅ Manage Roles
- ✅ Manage Channels
- ✅ Manage Webhooks
- ✅ View Audit Log
- ✅ Manage Messages
- ✅ Manage Threads
- ✅ Moderate Members
- ✅ Kick Members
- ✅ Ban Members

Use for:

- 🛠️ channel permission repairs
- 🧩 role order fixes
- 🤖 bot invite/permission troubleshooting
- 🖥️ dashboard support escalation

### 🧑‍💻 Developer

Purpose: Debug PulseKeep code, deploys, dashboard, API behavior, and bot runtime problems.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Webhooks in dev/log channels
- ✅ Manage Messages in test channels
- ✅ View Audit Log
- ✅ Manage Channels in test categories

Optional:

- ✅ Manage Server if the developer is trusted with OAuth/dashboard setup

Avoid:

- ❌ Administrator by default
- ❌ Ban Members by default
- ❌ Manage Roles by default

### 🔐 Security Lead

Purpose: Handle reports involving tokens, phishing, raids, abuse, exploit reports, and staff misconduct.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ View Audit Log
- ✅ Manage Messages
- ✅ Moderate Members
- ✅ Kick Members

Optional:

- ✅ Ban Members

Private access:

- 🔐 `#security-reports`
- 🚨 `#incident-response`
- 🧾 `#audit-review`

### 🧰 Support Manager

Purpose: Manage support flow, ticket quality, helper training, and repeated support issues.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages in support channels
- ✅ Manage Threads
- ✅ Manage Channels inside SUPPORT category

Optional:

- ✅ Moderate Members

### 🎧 Senior Support

Purpose: Trusted ticket handler for more complex setup problems.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages in support channels
- ✅ Manage Threads

Private access:

- 🎫 ticket channels
- 🧾 ticket logs
- 🧰 support notes

### 💬 Support Team

Purpose: Answer regular user questions and help with bot setup.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Add Reactions

Optional:

- ✅ Manage Messages in `#help-chat`, `#bug-reports`, and ticket channels only

Avoid:

- ❌ Manage Server
- ❌ Manage Roles
- ❌ Ban Members

### ⚖️ Head Moderator

Purpose: Leads moderation and trains moderators.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages
- ✅ Manage Threads
- ✅ Moderate Members
- ✅ Kick Members
- ✅ Ban Members
- ✅ Manage Nicknames

Private access:

- ⚖️ `#mod-chat`
- 🧾 `#mod-logs`
- 🚨 `#incident-response`

### 🔨 Moderator

Purpose: Handle public channel moderation.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages
- ✅ Moderate Members
- ✅ Kick Members

Optional:

- ✅ Ban Members if fully trusted

Avoid:

- ❌ Manage Roles
- ❌ Administrator

### 🧹 Trial Moderator

Purpose: Trainee moderation role with limited power.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages in public/community/help channels only

Avoid:

- ❌ Kick Members
- ❌ Ban Members
- ❌ Manage Channels
- ❌ Manage Roles

### 🧪 Bot Tester

Purpose: Test every PulseKeep feature without disrupting real support channels.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Add Reactions

Special access:

- 🧪 `#slash-command-testing`
- 🎰 `#economy-testing`
- 🎫 `#ticket-testing`
- 🖥️ `#dashboard-testing`

Avoid:

- ❌ Manage Messages
- ❌ Kick Members
- ❌ Ban Members

### 🐞 Bug Hunter

Purpose: Report reproducible bugs with steps and screenshots.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands

Special access:

- 🐞 `#bug-reports`
- 🧪 test channels if trusted

### 🧱 Contributor

Purpose: Help improve docs, wording, support flows, and ideas.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands

Optional:

- ✅ `#contributor-chat`
- ✅ `#docs-feedback`

### 🤝 Partner

Purpose: Partner representative role.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History

Avoid:

- ❌ staff logs
- ❌ moderation powers
- ❌ Manage Messages

### ⭐ VIP

Purpose: Trusted community recognition role.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ External Emojis and Stickers if enabled

### 🙋 Active Helper

Purpose: Recognition for helpful members.

Give:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands

Avoid:

- ❌ moderation powers
- ❌ staff logs
- ❌ Manage Messages

### 👤 Member

Purpose: Normal support server member.

Give:

- ✅ View public channels
- ✅ Send Messages in public/community/help channels
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Add Reactions
- ✅ Attach Files in support/report channels
- ✅ Embed Links if spam is under control

Avoid:

- ❌ Mention Everyone
- ❌ Manage Messages
- ❌ Manage Channels
- ❌ Manage Roles
- ❌ View staff channels

### 🔇 Muted

Purpose: Temporarily restrict disruptive users.

Deny:

- ❌ Send Messages
- ❌ Add Reactions
- ❌ Use Application Commands
- ❌ Create Public Threads
- ❌ Create Private Threads
- ❌ Speak in voice channels

Allow if desired:

- ✅ View Channels
- ✅ Read Message History

### 🤖 PulseKeep Bot

Purpose: Let the bot moderate, create tickets, send embeds, manage configured roles, and run slash commands.

Recommended granular permissions:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages
- ✅ Manage Channels
- ✅ Manage Roles
- ✅ Manage Nicknames
- ✅ Kick Members
- ✅ Ban Members
- ✅ Moderate Members

Setup notes:

- Put the bot role above `Member`, `Active Helper`, `VIP`, `Partner`, `Contributor`, `Bug Hunter`, and `Bot Tester`.
- Put the bot role above any role `/role` should add or remove.
- Put the bot role above members it needs to kick, ban, timeout, or nickname.
- Administrator is easiest for early testing, but granular permissions are cleaner for production.

</details>

---

## 🗂️ Category Layout

Create categories in this order:

```text
📌 START HERE
📚 PULSEKEEP INFO
🎫 SUPPORT
🤖 BOT COMMANDS
🧪 TEST LAB
💬 COMMUNITY
🧑‍💻 CONTRIBUTORS
🔒 STAFF
🧾 LOGS
🚨 INCIDENTS
📦 ARCHIVE
```

## 📌 START HERE

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 👋 `#welcome`

**Channel description:** Welcome to PulseKeep Support! First channel users see. Links to commands, status, dashboard, support, privacy, and terms.

Purpose:

- First channel users see.
- Explains what PulseKeep is.
- Links to the command page, status page, dashboard, support, privacy, and terms.

**Category permission sync:** Apply the 📌 START HERE category default permissions (Member: view only; Support/Mod/Admin: full access; Bot: standard). Right-click the category → Edit Category → Permissions → Copy permissions to all channels below it.

Permissions:

| Role | View | Send | React | Slash Commands |
| --- | --- | --- | --- | --- |
| 👤 Member | ✅ | ❌ | ✅ | ❌ |
| 💬 Support Team | ✅ | ✅ | ✅ | ✅ |
| 🔨 Moderator | ✅ | ✅ | ✅ | ✅ |
| 🛡️ Administrator | ✅ | ✅ | ✅ | ✅ |
| 🤖 PulseKeep Bot | ✅ | ✅ | ✅ | ✅ |

Starter message to send:

```text
👋 Welcome to PulseKeep Support!

PulseKeep is a Discord bot for moderation, tickets, economy, logging, and server operations.

👑 Owner: watispro1 · 👑 Co‑Owner: williamdelilah7_
⌨️ Commands: see #commands
🎫 Need help? Open a ticket in #ticket-panel
📊 Service status: see #status
🔐 Never share bot tokens, API keys, cookies, or database URLs.
```

### 📜 `#rules`

**Channel description:** Server rules and behavior expectations. Read before participating.

Purpose:

- Public rules.
- Safety expectations.
- Support expectations.

**Category permission sync:** Apply the 📌 START HERE category default permissions (Member: view only; Staff: can send). Sync via category → Edit Category → Permissions → Copy below.

Rules:

1. 🤝 Be respectful.
2. 🚫 No harassment, hate speech, threats, or targeted abuse.
3. 🧵 Keep support issues in the correct channels.
4. 🔐 Never post bot tokens, account tokens, API keys, database URLs, cookies, or secrets.
5. 🎫 Use tickets for private server-specific issues.
6. 🧪 Keep command spam in test channels.
7. 📣 Do not repeatedly ping staff.
8. ⚖️ Staff decisions are final.
9. ✅ Follow Discord Terms of Service.
10. 🧾 Report security issues privately in a ticket or security channel.

**Starter message to send:**

```text
📜 Please read the rules above before chatting.

1. 🤝 Be respectful
2. 🚫 No harassment
3. 🧵 Keep support in correct channels
4. 🔐 Never share tokens, keys, or secrets
5. 🎫 Private issues → #ticket-panel
6. 🧪 Command spam → #slash-command-testing
7. 📣 Don't ping staff repeatedly
8. ⚖️ Staff decisions are final
9. ✅ Follow Discord ToS
10. 🧾 Report security issues privately

Questions? Ask in #help-chat.
```

### 🧭 `#start-here`

**Channel description:** New here? Start with this onboarding checklist — invite the bot, explore commands, and check the dashboard.

Purpose:

- Simple onboarding checklist.

**Category permission sync:** Apply the 📌 START HERE category default permissions (Member: view only; Staff: can send). Sync via category → Edit Category → Permissions → Copy below.

Include:

- 🤖 bot invite link
- ⌨️ commands page
- 🖥️ dashboard page
- 📊 status page
- 🔐 privacy policy
- 📄 terms of service
- 🎫 ticket panel

**Starter message to send:**

```text
🧭 Welcome to PulseKeep!

📌 Quick links:
🤖 Invite bot: https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=2150636608&scope=bot%20applications.commands
⌨️ Commands: https://pulsekeep.fly.dev/commands.html
🖥️ Dashboard: https://pulsekeep.fly.dev/dashboard.html
📊 Status: https://pulsekeep.fly.dev/status.html
🔐 Privacy: https://pulsekeep.fly.dev/privacy.html
📄 Terms: https://pulsekeep.fly.dev/terms.html

Need help? Open a ticket in #ticket-panel or ask in #help-chat.
```

</details>

---

## 📚 PULSEKEEP INFO

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 📣 `#announcements`

**Channel description:** Official PulseKeep announcements — releases, outages, and updates. Members can read but only staff and developers can post.

Purpose:

- Major updates
- Outages
- DiscordBotList approval updates
- New feature releases

**Category permission sync:** Apply the 📚 PULSEKEEP INFO category default permissions (Member: view only; Staff: can send). Sync via category → Edit Category → Permissions → Copy below.

Permissions:

| Role | View | Send |
| --- | --- | --- |
| 👤 Member | ✅ | ❌ |
| 💬 Support Team | ✅ | ❌ |
| 🧑‍💻 Developer | ✅ | ✅ |
| 🛡️ Administrator | ✅ | ✅ |

**Starter message to send:**

```text
📣 Welcome to PulseKeep Announcements!

This channel is for official updates only. You cannot send messages here.

Expect:
- New feature releases
- Scheduled maintenance notices
- DiscordBotList status updates
- Outage reports and resolutions

Discussion → #general or #off-topic
```

### 📝 `#changelog`

**Channel description:** Release notes, version history, and command/dashboard changes. Mirrors the website changelog.

Purpose:

- Release notes.
- Website changelog mirrors.
- Command changes.
- Dashboard changes.

**Category permission sync:** Apply the 📚 PULSEKEEP INFO category default permissions (Member: view only). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
📝 **Changelog**

This channel tracks every PulseKeep release and update.

Latest: v7.1.0 — DBL voting, new commands, configure system, cache invalidation, bug fixes.
Full history: https://pulsekeep.fly.dev/changelog.html

Only staff can post here.
```

### 📊 `#status`

**Channel description:** Live service health updates, maintenance notices, and deploy logs. For real-time metrics use the website status page.

Purpose:

- Public bot and website health updates.
- Maintenance windows.
- Deploy notes.

**Category permission sync:** Apply the 📚 PULSEKEEP INFO category default permissions (Member: view only). Sync via category → Edit Category → Permissions → Copy below.

Starter message to pin:

```text
📊 PulseKeep Status

This channel is for service health updates, deploy notes, and outage notices.
For live details, use the website status page: https://pulsekeep.fly.dev/status.html
```

### ⌨️ `#commands`

**Channel description:** Full command reference. Run `/help` to browse interactively or see the pinned list below for all 50+ commands.

Purpose:

- Public command reference with full listing.
- Tell users to run `/help` to browse interactively.
- Link the website command page at `https://pulsekeep.fly.dev/commands.html`.
- Pin a full command list for quick reference.

**Category permission sync:** Apply the 📚 PULSEKEEP INFO category default permissions (Member: view only). Sync via category → Edit Category → Permissions → Copy below.

Starter message to pin — full command listing:

<details>
<summary><strong>📋 Full Command List (50+)</strong></summary>

**🎮 Economy** — `/balance` `/daily` `/weekly` `/work` `/gamble` `/blackjack` `/slots` `/rob` `/pay` `/search` `/fish` `/mine` `/shop` `/buy` `/inventory` `/use` `/leaderboard` `/tip` `/vote`

**🛡️ Moderation** — `/warn` `/warnings` `/clearwarns` `/history` `/mute` `/unmute` `/kick` `/ban` `/softban` `/purge` `/clean` `/slowmode` `/lock` `/unlock` `/nick` `/role add` `/role remove` `/move` `/vckick` `/announce` `/data-deletion`

**🎫 Tickets** — `/ticketpanel` `/ticket add` `/ticket remove` `/ticket close` `/ticket rename`

**⚙️ Configuration** — `/configure economy` `/configure tickets` `/configure modlogs` `/configure welcome` `/configure log_channel` `/configure welcome_channel` `/configure vote_channel` `/configure ticket_category` `/configure show`

**🔧 Utility** — `/help` `/ping` `/stats` `/invite` `/vote` `/servericon` `/roleinfo` `/channelinfo` `/tip`
</details>

Run `/help <category>` to filter (Economy, Moderation, Tickets, Configuration, Utility).

Starter commands:

```text
/help
/ping
/stats
/ticketpanel
/daily
/work
/balance
/vote
```

### ❓ `#faq`

**Channel description:** Frequently asked questions about PulseKeep — commands, economy, moderation, tickets, and troubleshooting.

Purpose:

- Public FAQ reference.
- Answer common questions without opening tickets.
- Reduce repeated support requests.

**Category permission sync:** Apply the 📚 PULSEKEEP INFO category default permissions (Member: view only). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
❓ **Frequently Asked Questions**

Browse the pinned FAQ entries below for answers to common questions about PulseKeep.

Topics covered:
- Slash commands not appearing
- Ticket creation issues
- Moderation command failures
- Dashboard login problems
- Economy and earning Pulses
- Vote rewards and cooldowns
- Permissions and setup
- Bug reporting
- Data deletion requests

If your question isn't answered here, ask in #help-chat or open a ticket in #ticket-panel.
```

Below are detailed Q&A entries to paste into the channel.

<details>
<summary><strong>❓ Why are slash commands not appearing?</strong></summary>

- Bot needs `applications.commands` scope on invite. Commands can take up to 1 hour to sync. Re-invite or restart Discord.
</details>

<details>
<summary><strong>❓ Why can't the bot create tickets?</strong></summary>

- Bot needs **Manage Channels**. Check `/configure tickets enabled:true` and that the ticket category is set correctly.
</details>

<details>
<summary><strong>❓ Why do moderation commands fail?</strong></summary>

- Bot's role must be **above** the target's highest role. Bot needs the relevant permission (Kick/Ban/Moderate Members).
</details>

<details>
<summary><strong>❓ Why does dashboard login fail?</strong></summary>

- Must use an account that's a member of the server with **Manage Server** permission. Clear browser cache and retry.
</details>

<details>
<summary><strong>❓ What permissions does PulseKeep need?</strong></summary>

- **Administrator** (recommended). Otherwise: View/Send Messages, Manage Channels, Kick/Ban/Moderate Members, Manage Nicknames, Manage Roles.
</details>

<details>
<summary><strong>❓ How do I request data deletion?</strong></summary>

- Open a support ticket with your user ID and guild ID(s). Deletion processed within 48 hours. Or ask the owner/co-owner to run `/data-deletion`.
</details>

<details>
<summary><strong>❓ How do I report a bug?</strong></summary>

- Post in `#bug-reports` with command, what happened, expected behavior, and screenshots if possible. Check for duplicates first.
</details>

<details>
<summary><strong>❓ Why does the status page show offline?</strong></summary>

- API at `pulsekeep.fly.dev` may be down. Check `/health` directly. Deployments cause brief restarts — wait 1-2 min.
</details>

<details>
<summary><strong>❓ How do I earn Pulses?</strong></summary>

- `/daily` `/weekly` `/work` `/gamble` `/blackjack` `/slots` `/search` `/fish` `/mine` `/vote`. Fish/mine need tools from `/shop`. Streaks reset if you miss a day.
</details>

<details>
<summary><strong>❓ How does blackjack work?</strong></summary>

- Bet Pulses, get 2 cards (one dealer hidden). Hit/stand to reach 21. Blackjack pays 2.5x, regular win 2x, push returns bet.
</details>

<details>
<summary><strong>❓ What items are in the shop?</strong></summary>

- Fishing Rod (2,500), Mining Pick (5,000), Treasure Map (1,000), Lucky Clover (3,000), EXP Boost (4,000). Run `/shop`.
</details>

<details>
<summary><strong>❓ Economy disabled but commands still work?</strong></summary>

- Run `/configure economy enabled:false` again and verify with `/configure show`.
</details>

<details>
<summary><strong>❓ How do I set up welcome messages?</strong></summary>

- Create a welcome channel, run `/configure welcome enabled:true`, then `/configure welcome_channel channel:#channel`.
</details>

<details>
<summary><strong>❓ How does /rob work?</strong></summary>

- Attempt to steal Pulses. Success depends on target balance. Failure pays a fine. Both need ≥100 Pulses.
</details>

<details>
<summary><strong>❓ What happens when I vote?</strong></summary>

- Vote via `/vote`, get 500-750 Pulses via webhook. Can vote every 12 hours. Announcement in vote channel if configured.
</details>

<details>
<summary><strong>❓ My question isn't listed here</strong></summary>

- Open a ticket in `#ticket-panel`, ask in `#help-chat`, or check the website.</details>

---

## 🎫 SUPPORT

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🧰 `#support-info`

**Channel description:** How PulseKeep support works — what details to include when asking for help and how to open a ticket.

Purpose:

- Explain how support works.
- Tell users what details to include.

**Category permission sync:** Apply the 🎫 SUPPORT category default permissions (Member: view and send). Sync via category → Edit Category → Permissions → Copy below.

Starter message to pin:

```text
Need help? Please include:

1. What command or feature you used
2. What happened
3. What you expected
4. Screenshots or error text
5. Whether PulseKeep has the needed permissions
6. Whether the bot role is above the target role/member

Never share bot tokens, API keys, cookies, account tokens, or database URLs.
```

### 💬 `#help-chat`

**Channel description:** Public help channel for quick setup questions, command help, and general PulseKeep support. Do not share private server info here.

Purpose:

- Public non-private help.
- Quick setup questions.
- General questions about commands.

**Category permission sync:** Apply the 🎫 SUPPORT category default permissions (Member: view and send). Sync via category → Edit Category → Permissions → Copy below.

Permissions:

- 👤 Members can send.
- 💬 Support Team can manage messages.
- 🔨 Moderators can manage messages.

**Starter message to send:**

```text
💬 **Help Chat**

Ask your PulseKeep questions here!

Before asking:
1. Check #faq for common answers
2. Run /help to browse commands
3. Check #commands for the full list

For private server-specific issues, open a ticket in #ticket-panel.
Never share bot tokens, API keys, or database URLs in public chat.
```

### 🐞 `#bug-reports`

**Channel description:** Report bugs you find in PulseKeep. Include reproduction steps, screenshots, and permission info.

Purpose:

- Public reproducible bug reports.

**Category permission sync:** Apply the 🎫 SUPPORT category default permissions (Member: view and send). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🐞 **Bug Reports**

Found a bug? Use the template below to report it.

Template:
Command/feature:
What happened:
Expected behavior:
Server permissions checked? yes/no
Bot role above target? yes/no
Screenshot/error:
Can staff reproduce it? yes/no

Check if your bug was already reported before posting.
For security issues (token leaks, exploits), open a ticket instead.
```

Template:

```text
🐞 Bug Report

Command/feature:
What happened:
Expected behavior:
Server permissions checked? yes/no
Bot role above target role/member? yes/no
Screenshot/error:
Can staff reproduce it? yes/no
```

### 💡 `#feature-requests`

**Channel description:** Suggest new features for PulseKeep. Describe what you want, why it helps, and who would use it.

Purpose:

- Community feature suggestions.

**Category permission sync:** Apply the 🎫 SUPPORT category default permissions (Member: view and send). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
💡 **Feature Requests**

Have an idea for PulseKeep? Use the template below.

Template:
Feature:
Why it helps:
Who would use it:
Priority (low/medium/high):
Example:

Upvote suggestions you like — staff review popular requests first.
```

Template:

```text
💡 Feature Request

Feature:
Why it helps:
Who would use it:
Priority:
Example:
```

### 🎫 `#ticket-panel`

**Channel description:** Open a support ticket here. Click the button below to create a private ticket channel for staff to help you.

Purpose:

- Permanent PulseKeep ticket panel.

**Category permission sync:** Apply the 🎫 SUPPORT category default permissions. Members can see the panel and click buttons but cannot send messages in this channel (only in the ticket channel). Sync via category → Edit Category → Permissions → Copy below.

Permissions:

| Role | View | Send | Use Buttons |
| --- | --- | --- | --- |
| 👤 Member | ✅ | ❌ | ✅ |
| 💬 Support Team | ✅ | ✅ | ✅ |
| 🤖 PulseKeep Bot | ✅ | ✅ | ✅ |

Setup:

1. Create the channel.
2. Make sure PulseKeep Bot can send messages and embeds.
3. Run `/ticketpanel`.
4. Pin the panel.
5. Open a test ticket.
6. Confirm the ticket is private.
7. Confirm staff can view the ticket.
8. Confirm the close button deletes/closes the ticket.

</details>

---

## 🤖 BOT COMMANDS

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🧭 `#command-menu`

**Channel description:** Browse and test PulseKeep's commands here. Run `/help` to explore categories or try any command freely.

Purpose:

- Let members run `/help`.
- Let members browse commands.
- Let staff test the command menu.

**Category permission sync:** Apply the 🤖 BOT COMMANDS category default permissions (Member: view and send commands; Staff: full access; Bot: required). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🧭 **Command Menu**

Welcome! This channel is for browsing and testing PulseKeep commands.

Start here:
/help             — Browse all commands by category
/help Economy     — Economy commands only
/help Moderation  — Moderation commands only
/ping             — Check bot latency
/stats            — View bot statistics

All commands are slash commands — just type / and browse.
```

### 🩺 `#bot-status-checks`

**Channel description:** Quick bot health checks. Run `/ping`, `/stats`, and `/help` to verify the bot is responding properly.

Purpose:

- Lightweight health checks.

**Category permission sync:** Apply the 🤖 BOT COMMANDS category default permissions (Member: view and send commands; Staff: full access; Bot: required). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🩺 **Bot Status Checks**

Use this channel to verify PulseKeep is healthy:

/ping    — Check latency
/stats   — View bot metrics
/help    — Verify slash commands work

If these work, the bot is running properly.
For service-level status, see #status or https://pulsekeep.fly.dev/status.html
```

Test commands:

```text
/ping
/stats
/help
```

### 🎰 `#economy`

**Channel description:** Economy gameplay and discussion — earn Pulses, gamble, fish, mine, and shop. Run economy commands here.

Purpose:

- Economy gameplay, questions about earning Pulses, cooldowns, and mechanics.
- Explain how the economy works to new users.

**Category permission sync:** Apply the 🤖 BOT COMMANDS category default permissions (Member: view and send commands; Staff: full access; Bot: required). Sync via category → Edit Category → Permissions → Copy below.

Economy overview to pin:

```text
💰 PulseKeep Economy Guide

PULSES
Pulses are the currency. Earn them through commands, daily streaks, and voting.

DAILY COMMANDS
/daily    — Claim once per day. Streak resets if you miss 24h.
/weekly   — Claim once per week.
/work     — Every 30 minutes. Base payout scales slightly with RNG.
/search   — Every 15 minutes. Find random amounts of Pulses.

GAMBLING
/gamble   — 55% win rate. Choose risk level (2x to 10x payout).
/blackjack — Standard blackjack. Dealer hits on 16, stands on 17.
/slots    — Three reels. Match 2 for small win, 3 for jackpot.

MINIGAMES
/fish     — Requires Fishing Rod (2,500 from shop). Catch fish for varying rewards.
/mine     — Requires Mining Pick (5,000 from shop). Mine gems and ores.

SOCIAL
/pay      — Transfer Pulses to another user.
/rob      — Attempt to rob another user. Fail = pay a fine.
/leaderboard — See the richest users.
/tip      — Random economy tip.

VOTING
/vote     — Vote on DiscordBotList for 500-750 Pulses. 12h cooldown.

SHOP
/buy      — Buy items from the shop.
/inventory — View your items.
/use      — Use items (Treasure Map, Lucky Clover, EXP Boost).
/shop     — View all available items and prices.
```

Allowed commands:

```text
/daily
/weekly
/work
/balance
/pay
/rob
/slots
/gamble
/blackjack
/fish
/mine
/search
/shop
/buy
/inventory
/use
/leaderboard
/tip
/vote
```

### 🎟️ `#tickets-demo`

**Channel description:** Try the ticket system here. Click the ticket panel button to see how tickets work before opening a real support ticket.

Purpose:

- Demo the ticket flow without touching real support.

**Category permission sync:** Apply the 🤖 BOT COMMANDS category default permissions (Member: view and send; Staff: full access; Bot: required). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🎟️ **Tickets Demo**

Want to see how PulseKeep tickets work? A ticket panel may be set up here for testing.

- Click the button to create a test ticket
- See how the private channel is created
- Close the ticket to see it deleted

Demo tickets are deleted automatically. Real support tickets open in #ticket-panel.
```

</details>

---

## 🧪 TEST LAB

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🧪 `#slash-command-testing`

**Channel description:** Safe space to test any slash command. Try commands, verify outputs, and report unexpected behavior.

Purpose:

- Test all slash commands.

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Recommended access:

- 🧑‍💻 Developer
- 🧪 Bot Tester
- 🐞 Bug Hunter
- 🛡️ Administrator
- 🤖 PulseKeep Bot

**Starter message to send:**

```text
🧪 **Slash Command Testing**

This channel is for testing PulseKeep commands. Try anything here.

Available testers: Developers, Bot Testers, Bug Hunters, and Administrators.

Report bugs in #bug-reports with reproduction steps.
```

### 🔨 `#moderation-testing`

**Channel description:** Test moderation commands safely. Use purge, mute, warn, kick, ban, lock, and more. Only use kick/ban against consenting test users.

Purpose:

- Test moderation commands safely.

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🔨 **Moderation Testing**

Test all moderation commands here. See the pinned message for the full command list and permissions table.

⚠️ Rules:
- Use /purge and /clean freely
- Use /mute, /warn, /lock, /slowmode freely
- /kick and /ban only against consenting test users
- /lock will lock THIS channel — use /unlock after

Report moderation bugs in #bug-reports.
```

All commands and their required permissions:

| Command | Bot Permission Required | User Permission Required |
|---------|------------------------|-------------------------|
| `/warn` | — | Moderate Members |
| `/warnings` | — | Moderate Members |
| `/clearwarns` | — | Moderate Members |
| `/history` | — | Moderate Members |
| `/mute` | Moderate Members | Moderate Members |
| `/unmute` | Moderate Members | Moderate Members |
| `/kick` | Kick Members | Kick Members |
| `/ban` | Ban Members | Ban Members |
| `/softban` | Ban Members | Ban Members |
| `/purge` | Manage Messages | Manage Messages |
| `/clean` | Manage Messages | Manage Messages |
| `/slowmode` | Manage Channels | Manage Channels |
| `/lock` | Manage Channels | Manage Channels |
| `/unlock` | Manage Channels | Manage Channels |
| `/nick` | Manage Nicknames | Manage Nicknames |
| `/role add/remove` | Manage Roles | Manage Roles |
| `/move` | Move Members | Move Members |
| `/vckick` | Move Members | Move Members |
| `/announce` | — | Manage Server |

Allowed test commands:

- `/purge` `/clean` — Test deleting messages with and without filters
- `/mute` `/unmute` — Test timing out and removing timeout
- `/warn` `/warnings` `/clearwarns` — Test warning system locally
- `/slowmode` — Test setting and removing slowmode
- `/lock` `/unlock` — Test channel lockdown (will lock the test channel!)
- `/announce` — Test announcement embed formatting
- `/role add` `/role remove` — Test role management
- `/nick` — Test nickname changes
- `/move` — Test moving members between voice channels
- `/vckick` — Test disconnecting from voice
- `/kick` — Only against consenting test users
- `/ban` — Only against consenting test users or alternate test accounts
- `/softban` — Only against alternate test accounts

### 🎰 `#economy-testing`

**Channel description:** Test economy commands — balances, gambling, fishing, mining, shop, and item usage without affecting real data.

Purpose:

- Test balances, items, cooldowns, gambling, fishing, mining, lottery, and shop behavior.

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🎰 **Economy Testing**

Test all economy features here:
- /daily, /weekly, /work, /search
- /gamble, /blackjack, /slots
- /fish (needs Fishing Rod from /shop)
- /mine (needs Mining Pick from /shop)
- /shop, /buy, /inventory, /use
- /pay, /rob, /leaderboard

Report economy bugs in #bug-reports.
```

### 🎫 `#ticket-testing`

**Channel description:** Test ticket creation, closing, user add/remove, and rename without creating real support tickets.

Purpose:

- Test ticket creation, close buttons, private permissions, and ticket logs.

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🎫 **Ticket Testing**

Test the ticket system here:
- Create a ticket using the panel button
- Test /ticket add @user and /ticket remove @user
- Test /ticket rename new-name
- Test the close button
- Check #ticket-logs for the closing summary

Real support tickets should be opened in #ticket-panel.
```

### 🖥️ `#dashboard-testing`

**Channel description:** Test the dashboard — OAuth login, guild loading, config saves, and permission states.

Purpose:

- Test OAuth login.
- Test guild loading.
- Test configuration saving.
- Test log channel IDs.
- Test welcome settings.
- Test permission-denied states.

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🖥️ **Dashboard Testing**

Test the web dashboard here: https://pulsekeep.fly.dev/dashboard.html

Test scenarios:
1. Login with Discord OAuth
2. Verify guild list loads (needs Manage Server permission)
3. Toggle economy/tickets/modlogs/welcome
4. Save changes
5. Log out and log in again
6. Test with an account that lacks Manage Server (should show 0 guilds)

Report dashboard bugs in #bug-reports with screenshots.
```

### 🚧 `#automod-testing`

**Channel description:** Test auto-moderation filters — spam, banned words, mass mentions, link blocking, and caps detection. Messages here may be silently deleted.

Purpose:

- Test auto-moderation filters (banned words, spam, mass mentions, link blocking, caps detection).

**Category permission sync:** Apply the 🧪 TEST LAB category default permissions (Testers and up: full access; Bot: required; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🚧 **Auto-Mod Testing**

Test auto-moderation filters in this channel.

⚠️ Messages you send here may be silently deleted by the bot.

Test scenarios:
- Send repeated short messages (spam test)
- Use a configured banned word
- Mass-mention 5+ users
- Post a link with link blocking enabled
- Send an ALL CAPS message

If auto-mod doesn't trigger, check that the filter is configured on your test server.
```

Testing scenarios:

```text
🧪 Spam test — Send 6+ short messages rapidly → check if spam filter triggers
🧪 Mass mention — Send a message with 5+ mentions (@user @user2 @user3 ...)
🧪 Banned word — Send a message containing a configured banned word
🧪 Link test — Send `https://example.com` with link blocking enabled
🧪 Caps test — Send an ALL CAPS MESSAGE THAT IS VERY LONG
```

Expected behavior:

- Auto-mod silently deletes the offending message.
- A log entry appears in `#automod-logs` (if configured).
- The user is not warned or muted unless auto-mod punishment is configured.
- Repeated offenses may trigger escalating actions.

</details>

---

## 💬 COMMUNITY

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 💬 `#general`

**Channel description:** General community chat — talk about PulseKeep, Discord bots, server management, or anything on topic.

Purpose:

- Public community chat.

**Category permission sync:** Apply the 💬 COMMUNITY category default permissions (Member: view and send; Staff: manage; Bot: standard). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
💬 **General Chat**

Welcome! Chat about PulseKeep, ask questions, share tips, or discuss Discord bots.

Rules still apply — be respectful, no harassment, no spam.

For support: #help-chat or #ticket-panel
For off-topic: #off-topic
```

### 🖼️ `#showcase`

**Channel description:** Share your PulseKeep setup — custom ticket panels, command menus, welcome messages, or server layouts.

Purpose:

- Users show their PulseKeep setup, ticket panels, command menus, or server layouts.

**Category permission sync:** Apply the 💬 COMMUNITY category default permissions (Member: view and send; Staff: manage; Bot: standard). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🖼️ **Showcase**

Show off your PulseKeep setup!

Post screenshots of:
- Custom ticket panels
- Economy leaderboards
- Welcome embeds
- Moderation log setups
- Dashboard configurations

Include what commands/settings you used so others can learn!
```

### 💡 `#suggestions`

**Channel description:** Quick ideas and suggestions for PulseKeep. For detailed feature requests use #feature-requests.

Purpose:

- Lightweight suggestions that are not full feature requests.

**Category permission sync:** Apply the 💬 COMMUNITY category default permissions (Member: view and send; Staff: manage; Bot: standard). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
💡 **Suggestions**

Have a quick idea for PulseKeep? Drop it here!

For detailed feature requests with examples, use #feature-requests instead.

Staff review suggestions regularly. Popular ideas get prioritized.
```

### ☕ `#off-topic`

**Channel description:** Casual chat that doesn't fit anywhere else. Keep it respectful and follow the rules.

Purpose:

- Optional casual chat.

**Category permission sync:** Apply the 💬 COMMUNITY category default permissions (Member: view and send; Staff: manage; Bot: standard). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
☕ **Off-Topic**

Chat about anything here (within reason). Games, music, tech, life — go for it.

Rules still apply:
- Be respectful
- No harassment
- No spam
- Follow Discord ToS

Support questions should go in #help-chat or #ticket-panel.
```

</details>

---

## 🧑‍💻 CONTRIBUTORS

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🧱 `#contributor-chat`

**Channel description:** Chat for contributors — discuss documentation, website ideas, and community improvements.

Purpose:

- Community contributors.
- Documentation ideas.
- Website feedback.

**Category permission sync:** Apply the 🧑‍💻 CONTRIBUTORS category default permissions (Contributor role and up: view and send; Bot: standard; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🧱 **Contributor Chat**

Welcome contributors! This space is for:
- Improving documentation
- Website feedback
- Community growth ideas
- Translation help

Check pinned messages for current projects and priorities.
```

### 📚 `#docs-feedback`

**Channel description:** Feedback on PulseKeep documentation — README, support docs, command descriptions, website copy, and DBL text.

Purpose:

- Improve README, support docs, command descriptions, DiscordBotList text, and website copy.

**Category permission sync:** Apply the 🧑‍💻 CONTRIBUTORS category default permissions (Contributor role and up: view and send; Bot: standard; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
📚 **Docs Feedback**

Help improve PulseKeep documentation!

What we need feedback on:
- Is the README clear?
- Are command descriptions accurate?
- Is the website copy helpful?
- Are support docs easy to follow?

Found a typo or confusing section? Post it here or suggest an edit.
```

### 🌐 `#translation-help`

**Channel description:** Help translate PulseKeep into other languages. Future localization efforts will be coordinated here.

Purpose:

- Optional future localization help.

**Category permission sync:** Apply the 🧑‍💻 CONTRIBUTORS category default permissions (Contributor role and up: view and send; Bot: standard; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🌐 **Translation Help**

This channel is for future localization efforts.

When translation support is added, we'll coordinate here for:
- Translating command descriptions
- Translating website content
- Translating docs

Languages planned: Spanish, French, German, Portuguese, Japanese.
Interested? Let us know what languages you speak!
```

</details>

---

## 🔒 STAFF

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🔒 `#staff-chat`

**Channel description:** Private staff coordination — discuss server management, policies, and internal matters. Staff only.

Purpose:

- Private staff coordination.

**Category permission sync:** Apply the 🔒 STAFF category default permissions (Staff roles: view and send based on role; Members: no access; Bot: as needed). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🔒 **Staff Chat**

Private channel for PulseKeep support staff.

Use this for:
- Server management discussions
- Policy decisions
- Internal announcements
- General staff coordination

Moderation-specific discussions → #mod-chat
Support-specific notes → #support-notes
```

### ⚖️ `#mod-chat`

**Channel description:** Moderation coordination — discuss rule enforcement, user reports, and moderation actions. Moderators and up.

Purpose:

- Moderation coordination.

**Category permission sync:** Apply the 🔒 STAFF category default permissions (Moderator and up: view and send; Support: no access; Members: no access; Bot: as needed). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
⚖️ **Mod Chat**

Private channel for moderation discussions.

Use this for:
- Rule enforcement decisions
- User behavior discussions
- Mod action reviews
- Policy clarification

Support-specific discussions → #staff-chat or #support-notes
```

### 🧰 `#support-notes`

**Channel description:** Internal knowledge base for support staff — document common issues, known fixes, and tricky dashboard/OAuth scenarios.

Purpose:

- Repeated support issues.
- Known fixes.
- Dashboard/OAuth gotchas.

**Category permission sync:** Apply the 🔒 STAFF category default permissions (Support Team and up: view and send; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🧰 **Support Notes**

Internal knowledge base for support staff.

Post common issues and their fixes here:
- "User says slash commands don't appear → They need to re-invite with applications.commands scope"
- "Dashboard login fails → Check OAuth redirect URI and third-party cookies"
- "Ticket button does nothing → Bot needs Manage Channels permission"

Keep this organized — prefix with [FIX], [GOTCHA], or [KNOWN ISSUE].
```

### 📋 `#review-queue`

**Channel description:** Track open bugs, pending features, DBL tasks, dashboard issues, and support escalations. Staff only.

Purpose:

- Track open items.

**Category permission sync:** Apply the 🔒 STAFF category default permissions (Support and up: view and send; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
📋 **Review Queue**

Track all open items here.

Use prefixes:
🐞 [BUG] Description — Assigned to @user
💡 [FEATURE] Description — Priority: medium
✅ [DBL] Task description
🖥️ [DASHBOARD] Issue description
🎫 [ESCALATION] User ID — Issue summary

Staff should update items as they're resolved. Move resolved items to archive channels.
```

Track:

- 🐞 open bugs
- 💡 pending features
- ✅ DiscordBotList tasks
- 🖥️ dashboard issues
- 🎫 support escalations

### 🛠️ `#staff-commands`

**Channel description:** Staff-only command testing and verification. Use moderation, configuration, and utility commands here.

Purpose:

- Staff-only command use and moderation command checks.

**Category permission sync:** Apply the 🔒 STAFF category default permissions (Staff roles: view and send; Members: no access; Bot: required). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🛠️ **Staff Commands**

Staff-only channel for testing and using commands that require elevated permissions.

Useful commands:
/configure show    — View current server config
/configure economy enabled:true/false
/configure tickets enabled:true/false
/purge 10         — Test message cleanup
/announce         — Test announcement formatting
/ticketpanel      — Set up ticket panel
/warn @user test  — Test warning system

For moderation testing with full command permissions, use #moderation-testing.
```

</details>

---

## 🧾 LOGS

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🔨 `#mod-logs`

**Channel description:** Moderation action log — kicks, bans, unbans, timeouts, purges, slowmode, lock/unlock, and nick changes. Staff and bot only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- kicks
- bans
- unbans
- timeouts
- purges
- slowmode changes
- lock/unlock actions
- nickname changes

### 🤖 `#bot-logs`

**Channel description:** Bot runtime log — startup, shutdown, command errors, warnings, and gateway events. Developers and admins only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- bot startup
- bot shutdown
- command errors
- runtime warnings
- gateway reconnects

### 🎫 `#ticket-logs`

**Channel description:** Ticket activity log — opened, closed, creator info, and staff actions. Staff and bot only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- ticket opened
- ticket closed
- ticket creator
- staff actions

### 📊 `#vote-logs`

**Channel description:** DBL vote event log — vote rewards, webhook results, cooldown hits, and errors. Staff and bot only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- DiscordBotList vote events — user ID, username, timestamp
- Vote reward credited (amount + new balance)
- Vote channel announcements sent to configured guild channels
- Failed webhook attempts (auth failures, missing user ID)
- Vote cooldown hits (user tried to vote again before 12h)

### 🗳️ Voting explained

PulseKeep integrates with DiscordBotList for voting rewards.

**How voting works:**
1. User runs `/vote` — gets the DBL vote link.
2. User votes on DiscordBotList (once every 12 hours).
3. DBL sends a webhook POST to `https://pulsekeep.fly.dev/api/dbl/webhook`.
4. PulseKeep verifies the webhook signature and credits 500-750 Pulses.
5. If the user is in any server with a `voteChannelId` configured, an announcement embed is posted.

**What staff need to configure:**
- Set a vote channel: `/configure vote_channel channel:#channel-name`
- Enable economy (voting requires economy): `/configure economy enabled:true`

**Verification endpoints:**
- Webhook: `POST /api/dbl/webhook` (called by DBL)
- Status check: `GET /api/dbl/webhook` (returns readiness)

### 🚧 `#automod-logs`

**Channel description:** Auto-moderation log — spam, banned words, links, caps, and mass mentions. Staff and bot only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- spam detection
- banned words
- link spam
- caps abuse
- mass mentions

### 🖥️ `#dashboard-logs`

**Channel description:** Dashboard activity log — config saves, login issues, permission errors, and failed saves. Developers and admins only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- config changes
- login issues
- permission errors
- failed saves

### 🚀 `#deploy-logs`

**Channel description:** Deployment log — deploy start, completion, failures, and health check issues. Developers and admins only.

**Category permission sync:** Apply the 🧾 LOGS category default permissions (Staff: view; Bot: view and send; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

Log:

- deploy started
- deploy completed
- deploy failed
- health check failures

</details>

---

## 🚨 INCIDENTS

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### 🚨 `#incident-response`

**Channel description:** Active incident response — security events, outages, bot down, dashboard failures, and abuse waves. Emergency staff only.

Purpose:

- Active security or outage response.

**Category permission sync:** Apply the 🚨 INCIDENTS category default permissions (Security Lead, Admin, Developer: full access; Other staff: view only; Members: no access; Bot: view and send). Sync via category → Edit Category → Permissions → Copy below.

Use for:

- leaked tokens
- service outages
- bot offline incidents
- dashboard auth failures
- raids or abuse waves

**Starter message to send:**

```text
🚨 **Incident Response**

This channel is for active incidents only.

When an incident occurs:
1. Post a summary: what happened, when, impact level
2. Assign a lead responder
3. Provide status updates as available
4. Post resolution summary with root cause

When the incident is resolved, move notes to #audit-review.
Non-urgent discussion goes in #staff-chat.
```

### 🔐 `#security-reports`

**Channel description:** Private security reports — vulnerabilities, token leaks, and abuse. Security Lead, Admin, and Developer only.

Purpose:

- Private vulnerability reports.
- Token leak reports.
- Abuse reports.

**Category permission sync:** Apply the 🚨 INCIDENTS category default permissions (Security Lead, Admin, Developer: view and send; Other staff: no access; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🔐 **Security Reports**

This channel is for reporting security issues.

Use for:
- Bot token leaks
- API key exposure
- Database vulnerabilities
- Exploit discoveries
- Abuse reports (raids, harassment campaigns)

Do NOT post security issues in public channels. Report them here.
Do NOT share leaked tokens outside this channel.
```

### 🧾 `#audit-review`

**Channel description:** Post-incident review — audit staff actions, suspicious activity, and permission changes. Admin and Security Lead only.

Purpose:

- Review suspicious activity.
- Review staff actions.
- Review permission changes.

**Category permission sync:** Apply the 🚨 INCIDENTS category default permissions (Admin, Security Lead: view and send; Other staff: view only; Members: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🧾 **Audit Review**

Post-incident review channel.

Use this for:
- Reviewing staff actions during incidents
- Documenting what went wrong
- Suggesting process improvements
- Reviewing permission changes

Format:
[INCIDENT] Brief summary
What happened:
Root cause:
Actions taken:
Prevention:
```

</details>

---

## 📦 ARCHIVE

<details>
<summary><strong>Click to expand — channels in this category</strong></summary>

### ✅ `#resolved-tickets`

**Channel description:** Archive of resolved ticket summaries. Do not store private user data unless covered by your privacy policy.

Purpose:

- Optional summaries of resolved tickets.

**Category permission sync:** Apply the 📦 ARCHIVE category default permissions (Staff: view; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
✅ **Resolved Tickets**

This channel archives summaries of resolved support tickets.

Format:
[TICKET] Username#0000
Issue: Brief description
Resolution: What was done
Closed by: @staff

⚠️ Do not store private user data (tokens, IDs, API keys) unless your privacy policy covers it.
```

Do not store private user data unless your privacy policy covers it.

### 🗃️ `#old-announcements`

**Channel description:** Archive of older announcements. Read-only reference for past updates.

Purpose:

- Archive older announcements.

**Category permission sync:** Apply the 📦 ARCHIVE category default permissions (Staff: view; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🗃️ **Old Announcements**

This channel stores past announcements for reference.

Older announcements are moved here from #announcements to keep that channel current.
```

### 🐞 `#old-bugs`

**Channel description:** Archive of fixed and closed bug reports. Reference for known resolved issues.

Purpose:

- Archive fixed bug reports.

**Category permission sync:** Apply the 📦 ARCHIVE category default permissions (Staff: view; Members: no access; Bot: no access). Sync via category → Edit Category → Permissions → Copy below.

**Starter message to send:**

```text
🐞 **Old Bugs**

This channel archives resolved bug reports for reference.

Fixed bugs are moved here from #bug-reports.
Check here before reporting a new bug — it may already be fixed!
```

</details>

---

## 🧩 Category Permission Defaults

> [!IMPORTANT]
> The tables below mirror Discord's channel permission interface — each row is a permission, each column is a role.  
> ✅ = **Allowed** (green check), ❌ = **Denied** (red X), — = **Neutral** (gray — uses role default, typically inherited from @everyone / role permissions).  
> Set these as **category-level permission overwrites** so all child channels inherit them.  
> The **PulseKeep Bot** role should always be above any member roles in the role list.

---

### 📌 Public Info Categories (START HERE · PULSEKEEP INFO)

**Used for:** Welcome, rules, announcements, commands list, FAQ, changelog, status.

| Permission | 👤 Member | 🛡️ Staff | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ✅ | ✅ | ✅ |
| Send Messages | ❌¹ | ✅ | ✅ |
| Send Messages in Threads | ❌ | ✅ | ✅ |
| Create Public Threads | ❌ | ✅ | ✅ |
| Create Private Threads | ❌ | ✅ | ✅ |
| Embed Links | — | ✅ | ✅ |
| Attach Files | — | ✅ | ✅ |
| Add Reactions | ✅ | ✅ | ✅ |
| Use External Emoji | — | ✅ | ✅ |
| Use External Stickers | — | — | — |
| Mention @everyone/@here/@role | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅ | ✅ |
| Manage Threads | ❌ | ✅ | — |
| Read Message History | ✅ | ✅ | ✅ |
| Send TTS Messages | ❌ | ❌ | ❌ |
| Use Application Commands | — | — | ✅ |
| Send Voice Messages | — | — | — |

¹ *Member can only view content; posting is staff/bot only in announcement channels.*

---

### 🎫 Support Category

**Used for:** Support info, help channel, bug reports, ticket panel.

| Permission | 👤 Member | 🧰 Support Team | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ✅² | ✅ | ✅ |
| Send Messages | ✅³ | ✅ | ✅ |
| Send Messages in Threads | ❌ | ✅ | ✅ |
| Create Public Threads | ❌ | ✅ | ❌ |
| Create Private Threads | ❌ | ✅ | ❌ |
| Embed Links | — | ✅ | ✅ |
| Attach Files | ✅ | ✅ | ✅ |
| Add Reactions | ✅ | ✅ | ✅ |
| Use External Emoji | — | ✅ | ✅ |
| Mention @everyone/@here/@role | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅ | ✅ |
| Manage Channels | ❌ | ❌ | ✅⁴ |
| Read Message History | ✅ | ✅ | ✅ |
| Send TTS Messages | ❌ | ❌ | ❌ |
| Use Application Commands | — | — | ✅ |

² *Member can see support channels but NOT `#ticket-panel` after panels are posted (hide that channel after setup).*  
³ *Member can send in help/report channels but NOT in `#ticket-panel`.*  
⁴ *Bot needs Manage Channels to create/delete ticket channels.*

---

### 🔒 Staff Category

**Used for:** Staff chat, mod chat, support notes, review queue, staff commands.

| Permission | 👤 Member | 💬 Support | 🔨 Moderator | 🛡️ Admin | 🤖 PulseKeep Bot |
|---|---|---|---|---|---|
| View Channels | ❌ | ✅⁵ | ✅⁶ | ✅ | ✅⁷ |
| Send Messages | ❌ | ✅ | ✅ | ✅ | ✅ |
| Send Messages in Threads | ❌ | ✅ | ✅ | ✅ | ✅ |
| Create Private Threads | ❌ | ✅ | ✅ | ✅ | ❌ |
| Embed Links | ❌ | ✅ | ✅ | ✅ | ✅ |
| Attach Files | ❌ | ✅ | ✅ | ✅ | ✅ |
| Add Reactions | ❌ | ✅ | ✅ | ✅ | ✅ |
| Mention @everyone/@here/@role | ❌ | ❌ | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅⁵ | ✅⁶ | ✅ | ✅ |
| Manage Channels | ❌ | ❌ | ❌ | ✅ | ❌ |
| Read Message History | ❌ | ✅ | ✅ | ✅ | ✅ |
| Use Application Commands | ❌ | — | — | — | ✅ |

⁵ *Support sees `#staff-chat`, `#support-notes`, `#review-queue` but NOT `#mod-chat` or admin-only channels.*  
⁶ *Moderator sees `#mod-chat`, `#staff-chat`, `#review-queue` but NOT admin-only channels.*  
⁷ *Bot sees only channels it needs to post to (logs, incident-response).*

---

### 🤖 Bot Commands Category

**Used for:** Command menu, status checks, economy interaction, ticket demo.

| Permission | 👤 Member | 🛡️ Staff | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ✅ | ✅ | ✅ |
| Send Messages | ✅ | ✅ | ✅ |
| Send Messages in Threads | ❌ | ✅ | ✅ |
| Create Public Threads | ❌ | ✅ | ❌ |
| Create Private Threads | ❌ | ✅ | ❌ |
| Embed Links | — | ✅ | ✅ |
| Attach Files | ✅ | ✅ | ✅ |
| Add Reactions | ✅ | ✅ | ✅ |
| Use External Emoji | — | ✅ | ✅ |
| Mention @everyone/@here/@role | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅ | ✅ |
| Read Message History | ✅ | ✅ | ✅ |
| Send TTS Messages | ❌ | ❌ | ❌ |
| Use Application Commands | ✅ | ✅ | ✅ |
| Send Voice Messages | — | — | — |

---

### 🧪 Test Lab Category

**Used for:** Slash command testing, moderation testing, economy testing, ticket/dashboard/automod testing.

| Permission | 👤 Member | 🧪 Bot Tester · 🐞 Bug Hunter | 🧑‍💻 Developer · 🛡️ Admin | 🤖 PulseKeep Bot |
|---|---|---|---|---|
| View Channels | ❌ | ✅ | ✅ | ✅ |
| Send Messages | ❌ | ✅ | ✅ | ✅ |
| Send Messages in Threads | ❌ | ✅ | ✅ | ✅ |
| Create Public Threads | ❌ | ✅ | ✅ | ❌ |
| Create Private Threads | ❌ | ✅ | ✅ | ❌ |
| Embed Links | ❌ | ✅ | ✅ | ✅ |
| Attach Files | ❌ | ✅ | ✅ | ✅ |
| Add Reactions | ❌ | ✅ | ✅ | ✅ |
| Mention @everyone/@here/@role | ❌ | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅ | ✅ | ✅ |
| Kick Members | ❌ | ❌ | ✅⁸ | ✅ |
| Ban Members | ❌ | ❌ | ✅⁸ | ✅ |
| Moderate Members | ❌ | ❌ | ✅⁸ | ✅ |
| Move Members | ❌ | ❌ | ✅ | ❌ |
| Deafen Members | ❌ | ❌ | ✅ | ❌ |
| Read Message History | ❌ | ✅ | ✅ | ✅ |
| Use Application Commands | ❌ | ✅ | ✅ | ✅ |

⁸ *Only in moderation test channels — restrict these permissions to specific test channels via channel-level overwrites.*

---

### 💬 Community Category

**Used for:** General chat, showcase, suggestions, off-topic.

| Permission | 👤 Member | 🛡️ Staff | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ✅ | ✅ | ✅ |
| Send Messages | ✅ | ✅ | ✅ |
| Send Messages in Threads | ✅ | ✅ | ✅ |
| Create Public Threads | ✅ | ✅ | ❌ |
| Create Private Threads | ❌ | ✅ | ❌ |
| Embed Links | ✅ | ✅ | ✅ |
| Attach Files | ✅ | ✅ | ✅ |
| Add Reactions | ✅ | ✅ | ✅ |
| Use External Emoji | ✅ | ✅ | ✅ |
| Use External Stickers | ✅ | ✅ | — |
| Mention @everyone/@here/@role | ❌ | ✅ | ❌ |
| Manage Messages | ❌ | ✅ | ✅ |
| Manage Threads | ❌ | ✅ | — |
| Read Message History | ✅ | ✅ | ✅ |
| Send TTS Messages | ❌ | ❌ | ❌ |
| Use Application Commands | — | — | ✅ |
| Send Voice Messages | — | — | — |

---

### 🧑‍💻 Contributors Category

**Used for:** Contributor chat, docs feedback, translation help.

| Permission | 👤 Member | 🧱 Contributor · 🐞 Bug Hunter · ⭐ VIP | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ❌ | ✅ | ❌⁹ |
| Send Messages | ❌ | ✅ | ❌ |
| Send Messages in Threads | ❌ | ✅ | ❌ |
| Create Public Threads | ❌ | ✅ | ❌ |
| Create Private Threads | ❌ | ✅ | ❌ |
| Embed Links | ❌ | ✅ | ❌ |
| Attach Files | ❌ | ✅ | ❌ |
| Add Reactions | ❌ | ✅ | ❌ |
| Mention @everyone/@here/@role | ❌ | ❌ | ❌ |
| Manage Messages | ❌ | ❌ | ❌ |
| Read Message History | ❌ | ✅ | ❌ |
| Use Application Commands | ❌ | — | ❌ |

⁹ *Bot generally doesn't need access here — skip the category entirely unless the bot needs to post automated contributor updates.*

---

### 🧾 Logs Category

**Used for:** Mod logs, bot logs, ticket logs, vote logs, automod logs, deploy logs.

| Permission | 👤 Member | 💬 Support | 🔨 Moderator | 🔐 Security | 🧑‍💻 Developer | 🛡️ Admin | 🤖 PulseKeep Bot |
|---|---|---|---|---|---|---|---|
| View Channels | ❌ | ✅¹⁰ | ✅¹⁰ | ✅ | ✅ | ✅ | ✅¹¹ |
| Send Messages | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Send Messages in Threads | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Embed Links | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Attach Files | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |
| Add Reactions | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Mention @everyone/@here/@role | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Manage Messages | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅¹² |
| Read Message History | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Use Application Commands | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |

¹⁰ *Support/Mod can VIEW logs (read-only) for audit purposes but cannot send or manage messages.*  
¹¹ *Bot has full write access to log channels to post embed-rich logs.*  
¹² *Bot may need Manage Messages to edit/clean up log embeds (optional).*

---

### 🚨 Incidents Category

**Used for:** Incident response, security reports, audit review.

| Permission | 👤 Member | 💬 Support | 🔨 Moderator | 🔐 Security Lead | 🧑‍💻 Developer | 🛡️ Admin | 🤖 PulseKeep Bot |
|---|---|---|---|---|---|---|---|
| View `#incident-response` | ❌ | ✅¹³ | ✅¹³ | ✅ | ✅ | ✅ | ✅¹⁴ |
| View `#security-reports` | ❌ | ❌ | ❌ | ✅ | ✅¹⁵ | ✅ | ❌ |
| View `#audit-review` | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ |
| Send Messages (all) | ❌ | ❌¹⁶ | ❌¹⁶ | ✅ | ✅ | ✅ | ✅¹⁴ |
| Embed Links | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Attach Files | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |
| Manage Messages | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |
| Read Message History | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Mention @everyone/@here/@role | ❌ | ❌ | ❌ | ✅ | ✅ | ✅ | ❌ |

¹³ *Support/Mod can monitor `#incident-response` and `#audit-review` but NOT `#security-reports`.*  
¹⁴ *Bot posts automated alerts to `#incident-response` only.*  
¹⁵ *Developer access to `#security-reports` is on a per-incident basis — grant/revoke as needed.*  
¹⁶ *Support/Mod are read-only during active incidents unless explicitly authorized by Security Lead.*

---

### 📦 Archive Category

**Used for:** Resolved tickets, old announcements, old bugs.

| Permission | 👤 Member | 🛡️ Staff | 🤖 PulseKeep Bot |
|---|---|---|---|
| View Channels | ❌ | ✅ | ❌¹⁷ |
| Send Messages | ❌ | ❌ | ❌ |
| Embed Links | ❌ | ❌ | ❌ |
| Attach Files | ❌ | ❌ | ❌ |
| Add Reactions | ❌ | ❌ | ❌ |
| Manage Messages | ❌ | ❌ | ❌ |
| Manage Channels | ❌ | ❌ | ❌ |
| Read Message History | ❌ | ✅ | ❌ |
| Use Application Commands | ❌ | ❌ | ❌ |

¹⁷ *Bot does not need archive access — archives are for staff reference only.*

---

## 🤖 PulseKeep Bot Invite Settings

Scopes:

```text
bot
applications.commands
```

Invite link (Administrator - quick setup):

```text
https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands
```

Production granular permissions value:

```text
https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=2150636608&scope=bot%20applications.commands
```

Granular permissions included:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Read Message History
- ✅ Use Application Commands
- ✅ Manage Messages
- ✅ Manage Channels
- ✅ Manage Roles
- ✅ Manage Nicknames
- ✅ Kick Members
- ✅ Ban Members
- ✅ Moderate Members

---

## 🎫 Ticket Workflow

### User flow:

1. 👤 User navigates to `#ticket-panel`.
2. 🎫 User clicks the **Open Ticket** button.
3. 🤖 PulseKeep creates a private channel named `ticket-username` under the configured ticket category.
4. 🚪 The ticket channel is only visible to the user and support staff.
5. 💬 Support Team joins and helps the user via the private channel.
6. ✅ When resolved, staff (or the user) clicks the **Close Ticket** button.
7. 🧾 PulseKeep sends a closing summary to `#ticket-logs`.
8. 🗑️ After 3 seconds, the ticket channel is deleted.

### Staff commands inside a ticket:

| Command | Description |
|---------|-------------|
| `/ticket add @user` | Add a user to the ticket |
| `/ticket remove @user` | Remove a user from the ticket |
| `/ticket close` | Close the current ticket |
| `/ticket rename new-name` | Rename the ticket channel |

### Setting up tickets:

1. (Optional) Create a category for tickets (e.g., `🎫 TICKETS`).
2. Run `/configure tickets enabled:true`.
3. Run `/configure ticket_category category:#your-category` (select the category, or skip to create tickets in the same channel category).
4. In your chosen channel, run `/ticketpanel` to create the button panel.
5. Pin the panel message so users can always find it.
6. Test by clicking the button — a private channel should appear.
7. Verify only staff and the user can view the new ticket.

### Troubleshooting tickets:

| Issue | Fix |
|-------|-----|
| Button does nothing | Check PulseKeep has Manage Channels permission |
| Ticket channel not private | Check the category permissions — PulseKeep needs to manage overwrites |
| Staff can't see tickets | Add staff roles to the ticket category permission overwrites |
| Ticket not deleted | PulseKeep deletes after 3 seconds — if it fails, check Manage Channels permission |
| `/ticket add` fails | The bot role must be above the target user's highest role |

### Recommended ticket channel names:

```text
ticket-username
ticket-userid
support-username
```

---

## 🖥️ Dashboard Workflow

### Access:

Dashboard is available at `https://pulsekeep.fly.dev/dashboard.html`

### Login flow:

1. 👤 User opens `https://pulsekeep.fly.dev/dashboard.html`.
2. 🔐 User clicks **Login with Discord**.
3. 📋 Discord OAuth screen asks for authorization (guilds and identify scopes).
4. 🔄 User is redirected back to the dashboard with an access token stored in localStorage.
5. 🧭 Dashboard loads all servers where PulseKeep is present and the user has Manage Server permission.
6. 🏠 User selects a guild from the sidebar.
7. ✅ Dashboard fetches the current config from the API.
8. ⚙️ User toggles features or updates channel IDs.
9. 💾 User clicks **Save Changes**.
10. 🤖 PulseKeep immediately applies the new configuration.
11. 📌 A save confirmation toast appears.

### Configurable settings:

| Setting | Type | Description |
|---------|------|-------------|
| Economy System | Toggle | Enable/disable all economy commands |
| Tickets | Toggle | Enable/disable ticket button and commands |
| Moderation Logs | Toggle | Enable/disable moderation action logging |
| Welcome Messages | Toggle | Enable/disable welcome embeds for new members |
| Ticket Category | Channel ID | Category where ticket channels are created |
| Log Channel | Channel ID | Where moderation logs are posted |
| Welcome Channel | Channel ID | Where welcome embeds are posted |
| Vote Channel | Channel ID | Where vote announcements from DBL are posted |

### Troubleshooting:

> [!TIP]
> | Issue | Cause | Fix |
> |-------|-------|-----|
> | ❌ **Login fails** | OAuth redirect URI mismatch | Check it matches `https://pulsekeep.fly.dev/auth/discord/callback` |
> | ❌ **Guilds do not load** | No mutual servers or no Manage Server permission | Make sure PulseKeep is in your server and you have Manage Server |
> | ❌ **Guild list empty** | Token expired | Log out and log in again |
> | ❌ **Save fails** | API server error or DB down | Check `https://pulsekeep.fly.dev/health` — retry in a minute |
> | ❌ **Permission denied** | User lacks Manage Server | Only users with Manage Server or Administrator can configure |
> | ❌ **Config does not persist** | Database connection error | Check database connectivity — report in support if persistent |
> | ❌ **Logged out unexpectedly** | Token expired or deleted | Tokens expire after ~7 days — log in again |

---

## ⏱️ Economy Cooldown Reference

| Command | Cooldown | Notes |
|---------|----------|-------|
| `/daily` | 24 hours | Streak resets if missed |
| `/weekly` | 7 days | Resets Monday midnight UTC |
| `/work` | 30 minutes | Base 50-150 Pulses |
| `/search` | 15 minutes | 20-100 Pulses |
| `/gamble` | None | Limited by your balance |
| `/blackjack` | None | Limited by your balance |
| `/slots` | None | Limited by your balance |
| `/rob` | 10 minutes | Fail = pay fine |
| `/fish` | 30 minutes | Needs Fishing Rod |
| `/mine` | 30 minutes | Needs Mining Pick |
| `/vote` | 12 hours | DBL cooldown, not bot cooldown |

Items that affect cooldowns:
- **EXP Boost** (from `/shop`) — 1.5x earnings for 30 minutes, does not reduce cooldowns
- **Lucky Clover** — Doubles next `/gamble` win only

---

## 🎨 Embed Color Reference

PulseKeep uses consistent embed colors across all commands:

| Color | Hex | Used For |
|-------|-----|----------|
| Purple | `#7c5cfc` | Default / Utility / Help |
| Green | `#22c55e` | Success / Welcome / Economy positive actions |
| Red | `#ef4444` | Errors / Failures / Warnings |
| Amber | `#eab308` | Warnings / Timeouts / Gambling |
| Cyan | `#22c8e5` | Tickets / Info commands |
| Blue-gray | `#5865f2` | Moderation actions (warn, kick, ban, mute) |

All embeds include a footer with the bot name and a timestamp.

---

## 📊 Status Workflow

Status page shows:

- 🩺 Discord Bot status (online/offline + guild count)
- 🗄️ API Server health
- 🗄️ Database connectivity
- ⏱️ Average latency
- 👥 Total users across all guilds
- ⏱️ Bot uptime

The status page auto-refreshes every 30 seconds.

If status shows offline:

1. Check the bot backend is running (`flyctl logs`).
2. Check `/health` endpoint directly: `curl https://pulsekeep.fly.dev/health`
3. Check deploy logs: `flyctl logs`
4. Check the Fly.io machine is running: `flyctl status`
5. Check the database is reachable: `flyctl ssh -- pg_isready`

---

## ✅ DiscordBotList Readiness Checklist

Before submitting:

- ✅ Bot is online at `https://pulsekeep.fly.dev`
- ✅ Bot can be invited (Administrator or granular permissions)
- ✅ 50+ slash commands registered globally
- ✅ `/help` works with category filtering
- ✅ `/ping` works with latency display
- ✅ `/stats` works with real-time bot metrics
- ✅ `/ticketpanel` works with button-based ticket creation
- ✅ `/vote` works with DBL API verification
- ✅ Support server invite works
- ✅ Website loads (`https://pulsekeep.fly.dev`)
- ✅ Commands page loads with search
- ✅ Status page shows live service health
- ✅ Dashboard login with Discord OAuth2
- ✅ Privacy policy exists at `/privacy.html`
- ✅ Terms of service exists at `/terms.html`
- ✅ DBL webhook configured at `/api/dbl/webhook`
- ✅ Owner is listed as `watispro1` · Co‑Owner listed as `williamdelilah7_`
- ✅ No private backend origin is shown publicly
- ✅ No bot tokens or secrets in client-side code
- ✅ Support server has clear rules
- ✅ Support server has a ticket channel
- ✅ Support server has a command guide
- ✅ Support server has a status channel
- ✅ Vote rewards automatically credit Pulses

---

## 📣 First Announcement

Post this in `#announcements`:

```text
📣 PulseKeep Support is open!

Welcome to the official PulseKeep support server.

🤖 PulseKeep is a Discord bot for moderation, tickets, economy, and server analytics.
👑 Owner: watispro1 · 👑 Co‑Owner: williamdelilah7_
⌨️ 50+ slash commands — run /help to browse
🎫 Need help? Open a ticket in #ticket-panel
📊 Service status: https://pulsekeep.fly.dev/status.html
🗳️ Vote for Pulses: /vote
🔐 Never share bot tokens, API keys, or database URLs.

Thanks for helping test and improve PulseKeep.
```
