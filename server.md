# 🛡️ PulseKeep Support Server Blueprint

Owner: `watispro1`  
Recommended server name: `PulseKeep Support`  
Recommended permanent invite label: `pulsekeep-support`

> [!NOTE]
> This file is a complete Discord support-server build plan for PulseKeep. It is written so you can create the server category by category, role by role, and permission by permission without guessing.

## 🎯 Server Purpose

PulseKeep Support should exist for these core purposes:

- 👑 **Ownership hub** - Clearly shows that PulseKeep is owned by `watispro1`.
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

Purpose: Emergency backup if the owner is unavailable.

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

### 👋 `#welcome`

Purpose:

- First channel users see.
- Explains what PulseKeep is.
- Links to the command page, status page, dashboard, support, privacy, and terms.

Permissions:

| Role | View | Send | React | Slash Commands |
| --- | --- | --- | --- | --- |
| 👤 Member | ✅ | ❌ | ✅ | ❌ |
| 💬 Support Team | ✅ | ✅ | ✅ | ✅ |
| 🔨 Moderator | ✅ | ✅ | ✅ | ✅ |
| 🛡️ Administrator | ✅ | ✅ | ✅ | ✅ |
| 🤖 PulseKeep Bot | ✅ | ✅ | ✅ | ✅ |

Suggested message:

```text
👋 Welcome to PulseKeep Support!

PulseKeep is a Discord bot for moderation, tickets, economy, logging, and server operations.

👑 Owner: watispro1
⌨️ Commands: see #commands
🎫 Need help? Open a ticket in #ticket-panel
📊 Service status: see #status
🔐 Never share bot tokens, API keys, cookies, or database URLs.
```

### 📜 `#rules`

Purpose:

- Public rules.
- Safety expectations.
- Support expectations.

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

### 🧭 `#start-here`

Purpose:

- Simple onboarding checklist.

Include:

- 🤖 bot invite link
- ⌨️ commands page
- 🖥️ dashboard page
- 📊 status page
- 🔐 privacy policy
- 📄 terms of service
- 🎫 ticket panel

## 📚 PULSEKEEP INFO

### 📣 `#announcements`

Purpose:

- Major updates
- Outages
- DiscordBotList approval updates
- New feature releases

Permissions:

| Role | View | Send |
| --- | --- | --- |
| 👤 Member | ✅ | ❌ |
| 💬 Support Team | ✅ | ❌ |
| 🧑‍💻 Developer | ✅ | ✅ |
| 🛡️ Administrator | ✅ | ✅ |

### 📝 `#changelog`

Purpose:

- Release notes.
- Website changelog mirrors.
- Command changes.
- Dashboard changes.

### 📊 `#status`

Purpose:

- Public bot and website health updates.
- Maintenance windows.
- Deploy notes.

Suggested pinned message:

```text
📊 PulseKeep Status

This channel is for service health updates, deploy notes, and outage notices.
For live details, use the website status page.
```

### ⌨️ `#commands`

Purpose:

- Public command reference.
- Tell users to run `/help`.
- Link the website command page.

Starter commands:

```text
/help
/ping
/stats
/uptime
/ticketpanel
/daily
/work
/balance
/profile
```

### ❓ `#faq`

Add answers for:

- ❓ Why slash commands are not appearing
- ❓ Why the bot cannot create tickets
- ❓ Why moderation commands fail
- ❓ Why dashboard login fails
- ❓ What permissions PulseKeep needs
- ❓ How to request data deletion
- ❓ How to report bugs
- ❓ Why the status page says offline

## 🎫 SUPPORT

### 🧰 `#support-info`

Purpose:

- Explain how support works.
- Tell users what details to include.

Suggested message:

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

Purpose:

- Public non-private help.
- Quick setup questions.
- General questions about commands.

Permissions:

- 👤 Members can send.
- 💬 Support Team can manage messages.
- 🔨 Moderators can manage messages.

### 🐞 `#bug-reports`

Purpose:

- Public reproducible bug reports.

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

Purpose:

- Community feature suggestions.

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

Purpose:

- Permanent PulseKeep ticket panel.

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

## 🤖 BOT COMMANDS

### 🧭 `#command-menu`

Purpose:

- Let members run `/help`.
- Let members browse commands.
- Let staff test the command menu.

### 🩺 `#bot-status-checks`

Purpose:

- Lightweight health checks.

Test commands:

```text
/ping
/stats
/uptime
/about
```

### 🎰 `#economy`

Purpose:

- Economy usage and economy questions.

Allowed commands:

```text
/daily
/weekly
/work
/balance
/profile
/pay
/coinflip
/slots
/gamble
/blackjack
/fish
/mine
/shop
/buy
/sell
/inventory
/use
/gift
/rich
/lottery
/lottery-claim
```

### 🎟️ `#tickets-demo`

Purpose:

- Demo the ticket flow without touching real support.

## 🧪 TEST LAB

### 🧪 `#slash-command-testing`

Purpose:

- Test all slash commands.

Recommended access:

- 🧑‍💻 Developer
- 🧪 Bot Tester
- 🐞 Bug Hunter
- 🛡️ Administrator
- 🤖 PulseKeep Bot

### 🔨 `#moderation-testing`

Purpose:

- Test moderation safely.

Allowed:

- `/purge`
- `/timeout`
- `/slowmode`
- `/lock`
- `/unlock`
- `/announce`
- `/role`
- `/kick` only against consenting test users
- `/ban` only against consenting test users or alternate test accounts

### 🎰 `#economy-testing`

Purpose:

- Test balances, items, cooldowns, gambling, fishing, mining, lottery, and shop behavior.

### 🎫 `#ticket-testing`

Purpose:

- Test ticket creation, close buttons, private permissions, and ticket logs.

### 🖥️ `#dashboard-testing`

Purpose:

- Test OAuth login.
- Test guild loading.
- Test configuration saving.
- Test log channel IDs.
- Test welcome settings.
- Test permission-denied states.

### 🚧 `#automod-testing`

Purpose:

- Test banned words.
- Test spam detection.
- Test mass mentions.
- Test link spam.
- Test caps detection.

## 💬 COMMUNITY

### 💬 `#general`

Purpose:

- Public community chat.

### 🖼️ `#showcase`

Purpose:

- Users show their PulseKeep setup, ticket panels, command menus, or server layouts.

### 💡 `#suggestions`

Purpose:

- Lightweight suggestions that are not full feature requests.

### ☕ `#off-topic`

Purpose:

- Optional casual chat.

## 🧑‍💻 CONTRIBUTORS

### 🧱 `#contributor-chat`

Purpose:

- Community contributors.
- Documentation ideas.
- Website feedback.

### 📚 `#docs-feedback`

Purpose:

- Improve README, support docs, command descriptions, DiscordBotList text, and website copy.

### 🌐 `#translation-help`

Purpose:

- Optional future localization help.

## 🔒 STAFF

### 🔒 `#staff-chat`

Purpose:

- Private staff coordination.

### ⚖️ `#mod-chat`

Purpose:

- Moderation coordination.

### 🧰 `#support-notes`

Purpose:

- Repeated support issues.
- Known fixes.
- Dashboard/OAuth gotchas.

### 📋 `#review-queue`

Track:

- 🐞 open bugs
- 💡 pending features
- ✅ DiscordBotList tasks
- 🖥️ dashboard issues
- 🎫 support escalations

### 🛠️ `#staff-commands`

Purpose:

- Staff-only command use and moderation command checks.

## 🧾 LOGS

### 🔨 `#mod-logs`

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

Log:

- bot startup
- bot shutdown
- command errors
- runtime warnings
- gateway reconnects

### 🎫 `#ticket-logs`

Log:

- ticket opened
- ticket closed
- ticket creator
- staff actions

### 📊 `#vote-logs`

Log:

- DiscordBotList vote events
- voter announcement
- vote reward summary

### 🚧 `#automod-logs`

Log:

- spam detection
- banned words
- link spam
- caps abuse
- mass mentions

### 🖥️ `#dashboard-logs`

Log:

- config changes
- login issues
- permission errors
- failed saves

### 🚀 `#deploy-logs`

Log:

- deploy started
- deploy completed
- deploy failed
- health check failures

## 🚨 INCIDENTS

### 🚨 `#incident-response`

Purpose:

- Active security or outage response.

Use for:

- leaked tokens
- service outages
- bot offline incidents
- dashboard auth failures
- raids or abuse waves

### 🔐 `#security-reports`

Purpose:

- Private vulnerability reports.
- Token leak reports.
- Abuse reports.

### 🧾 `#audit-review`

Purpose:

- Review suspicious activity.
- Review staff actions.
- Review permission changes.

## 📦 ARCHIVE

### ✅ `#resolved-tickets`

Purpose:

- Optional summaries of resolved tickets.

Do not store private user data unless your privacy policy covers it.

### 🗃️ `#old-announcements`

Purpose:

- Archive older announcements.

### 🐞 `#old-bugs`

Purpose:

- Archive fixed bug reports.

## 🧩 Category Permission Defaults

### 📌 Public Info Categories

Use for:

- START HERE
- PULSEKEEP INFO

Member:

- ✅ View Channels
- ✅ Read Message History
- ❌ Send Messages in announcement/rules channels
- ✅ Send Messages in allowed help/community channels

Staff:

- ✅ Send Messages
- ✅ Manage Messages

PulseKeep Bot:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Read Message History

### 🎫 Support Category

Member:

- ✅ View support info
- ✅ Send in help/report channels
- ✅ Click ticket button
- ❌ Send in `#ticket-panel` after panel is posted

Support Team:

- ✅ View support channels
- ✅ Send Messages
- ✅ Manage Messages
- ✅ Read Message History

PulseKeep Bot:

- ✅ View Channels
- ✅ Send Messages
- ✅ Embed Links
- ✅ Attach Files
- ✅ Manage Channels
- ✅ Read Message History

### 🔒 Staff Category

Member:

- ❌ View Channels

Support Team:

- ✅ View selected staff/support channels

Moderator:

- ✅ View mod channels

Admin:

- ✅ View all

PulseKeep Bot:

- ✅ View log/staff channels it needs

## 🤖 PulseKeep Bot Invite Settings

Scopes:

```text
bot
applications.commands
```

Easy testing permission integer:

```text
8
```

Production granular permissions:

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

## 🎫 Ticket Workflow

1. 👤 User opens `#ticket-panel`.
2. 🎫 User clicks the ticket button.
3. 🤖 PulseKeep creates a private ticket channel.
4. 💬 Support Team joins the conversation.
5. 🧰 Staff helps the user.
6. ✅ Staff closes the ticket.
7. 🧾 PulseKeep deletes or archives the ticket channel.
8. 📌 Ticket action appears in `#ticket-logs`.

Recommended ticket channel names:

```text
ticket-username
ticket-userid
support-username
```

## 🖥️ Dashboard Workflow

1. 👤 User opens the dashboard.
2. 🔐 User logs in with Discord.
3. 🧭 Dashboard loads guilds.
4. 🏠 User selects a guild.
5. ✅ Dashboard checks Manage Server permission.
6. ⚙️ User updates configuration.
7. 💾 Dashboard saves configuration.
8. 🤖 PulseKeep applies the setting.

Troubleshooting:

> [!TIP]
> - ❌ **Login fails**: Check Discord OAuth redirect URI.
> - ❌ **Guilds do not load**: Check `/api/guilds` Worker route.
> - ❌ **Save fails**: Check Worker POST body forwarding.
> - ❌ **Permission denied**: User needs Manage Server permission.
> - ❌ **Config does not persist**: Check database connection.

## 📊 Status Workflow

Status page should show:

- 🩺 API health
- 🗄️ database health
- 🏠 guild count
- 👥 user count
- ⚡ commands run
- ⏱️ API uptime
- 🤖 bot uptime
- 🧬 runtime version
- 🧠 memory usage
- 🧵 goroutines

If status shows offline:

1. Check the bot backend is running.
2. Check the Cloudflare Worker has the backend origin secret.
3. Check `/health` through the public website.
4. Check deploy logs.
5. Check bot logs.

## ✅ DiscordBotList Readiness Checklist

Before submitting:

- ✅ Bot is online.
- ✅ Bot can be invited.
- ✅ Slash commands are registered.
- ✅ `/help` works.
- ✅ `/ping` works.
- ✅ `/stats` works.
- ✅ `/ticketpanel` works.
- ✅ Support server invite works.
- ✅ Website loads.
- ✅ Commands page loads.
- ✅ Status page works.
- ✅ Dashboard login opens.
- ✅ Privacy policy exists.
- ✅ Terms of service exists.
- ✅ Owner is listed as `watispro1`.
- ✅ No private backend origin is shown publicly.
- ✅ Support server has clear rules.
- ✅ Support server has a ticket channel.
- ✅ Support server has a command guide.
- ✅ Support server has a status channel.

## 📣 First Announcement

Post this in `#announcements`:

```text
📣 PulseKeep Support is open!

Welcome to the official PulseKeep support server.

👑 Owner: watispro1
🎫 Need help? Open a ticket in #ticket-panel
⌨️ Need commands? Check #commands
📊 Need service status? Check #status

Thanks for helping test and improve PulseKeep.
```
