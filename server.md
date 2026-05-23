# PulseKeep Discord Server Blueprint

This blueprint is the recommended layout for a PulseKeep support and testing server. It keeps permissions category-first, gives staff a clear support workflow, and gives members a clean place to discover commands through PulseKeep's interactive Discord menu.

## Core Goals

- Keep public channels readable and low-noise.
- Give support staff one obvious ticket workflow.
- Keep bot testing contained to sandbox channels.
- Use category permissions first, with channel overrides only when they are truly needed.
- Make `/help`, `/menu`, and the ticket panel the main way users interact with PulseKeep.

## Role Hierarchy

Create roles in this order from highest to lowest.

| Role | Color | Purpose | Key Permissions |
| --- | --- | --- | --- |
| Founder | `#E74C3C` | Project owner and final authority. | Administrator |
| Administrator | `#E67E22` | Senior operations and server configuration. | Manage Server, Manage Channels, Manage Roles, View Audit Log, Ban Members, Kick Members, Moderate Members |
| Moderator | `#2ECC71` | Community safety and chat moderation. | Manage Messages, Moderate Members, Kick Members, Ban Members, View Audit Log |
| Support Team | `#3498DB` | Handles tickets, setup help, and customer questions. | Manage Messages, Send Messages, Attach Files, Embed Links, Read Message History |
| PulseKeep Bot | `#9B59B6` | The bot account. | Manage Messages, Send Messages, Embed Links, Attach Files, Read Message History, Use Slash Commands |
| Server Booster | `#FD79A8` | Trusted supporters with small perks. | Attach Files, Embed Links, Use External Emojis |
| Verified Member | `#95A5A6` | Standard trusted community member. | View Channels, Send Messages, Read Message History, Use Slash Commands |
| @everyone | Default | Pre-verification visitors. | View only the welcome/rules area |

## Channel Layout

### 1. Information

Use this category for read-only server information.

Channels:

- `rules-and-info`: Rules, useful links, bot invite, and support expectations.
- `announcements`: Product updates, releases, outages, and important changes.
- `status-logs`: Automated deploy, uptime, and incident notices.
- `welcome`: New member greeting and first steps.

Permissions:

- Staff can view, send, and manage messages.
- PulseKeep Bot can view, send, and embed links.
- Verified Member and @everyone can view and read history.
- Verified Member and @everyone cannot send messages.

Recommended pinned message for `rules-and-info`:

```text
Welcome to PulseKeep.

Start here:
1. Read the rules.
2. Use /help or /menu to browse commands.
3. Use the ticket panel in #open-a-ticket when you need private setup help.

Do not post bot tokens, database URLs, or private server logs in public channels.
```

### 2. Community

Use this category for normal conversation.

Channels:

- `general-chat`: Main community chat.
- `bot-discussion`: Usage questions and feature ideas.
- `showcase`: Server setups, panels, and PulseKeep configurations.
- `economy-chat`: Member economy commands such as `/daily`, `/work`, `/balance`, `/profile`, `/coinflip`, `/leaderboard`, and `/pay`.

Permissions:

- Verified Member can view, send messages, use slash commands, react, and read history.
- Server Booster can also attach files and embed links.
- Staff can manage messages.
- @everyone cannot view this category.

Recommended settings:

- Slowmode: 3 seconds in `general-chat`.
- Slowmode: 5 seconds in `economy-chat`.
- Disable link embeds for normal members unless you trust the community.

### 3. Command Center

Use this category to make PulseKeep's interactive command experience easy to find.

Channels:

- `command-menu`: The permanent PulseKeep command menu.
- `bot-sandbox`: General testing for slash commands.
- `moderation-lab`: Staff-only testing for moderation commands.

Permissions:

- Verified Member can view `command-menu` and `bot-sandbox`.
- Verified Member can use slash commands in both.
- Moderator and above can view and use `moderation-lab`.
- PulseKeep Bot can send messages and embed links everywhere in the category.

Setup:

1. In `command-menu`, run `/menu`.
2. Pin the bot's interactive command browser message.
3. Tell users to use the dropdown to switch between Utility, Moderation, Economy, and Tickets.
4. Keep normal chatting out of this channel so the menu is always visible.

Command aliases:

- `/help`: Opens the interactive command browser privately.
- `/menu`: Opens the same interactive command browser privately.
- `!menu`: Posts the menu publicly in a text channel.
- `!help`: Posts the menu publicly in a text channel.

### 4. Client Support

Use this category for support entry points.

Channels:

- `support-faq`: Answers to common setup, deploy, and permission issues.
- `open-a-ticket`: The permanent ticket panel.
- `pre-sales`: Questions about premium features, setup help, and custom work.

Permissions:

- Verified Member can view, send messages, read history, and use slash commands.
- Support Team, Moderator, Administrator, and Founder can manage messages.
- @everyone cannot view this category.
- PulseKeep Bot can send messages, embed links, and attach files.

Setup:

1. In `open-a-ticket`, run `/ticketpanel`.
2. Pin the ticket panel.
3. Keep the channel locked to support questions only.
4. Ask users to include server name, command or feature, expected behavior, actual behavior, and any safe error text.

Recommended `support-faq` topics:

- Fly.io deploy checklist.
- Netlify website checklist.
- Missing Discord permissions.
- Bot token and secret safety.
- Database connection troubleshooting.
- Slash commands not appearing.

### 5. Active Tickets

Use this category for private ticket channels.

Channels:

- `ticket-0001`, `ticket-0002`, and so on, generated by the bot once ticket creation is fully implemented.

Base category permissions:

- Founder, Administrator, and Support Team can view, send, attach files, embed links, and read history.
- PulseKeep Bot can view, send, embed links, attach files, and manage channels if ticket creation is enabled.
- @everyone and Verified Member cannot view.

Dynamic ticket override:

- The ticket creator gets View Channel, Send Messages, Attach Files, and Read Message History.
- Do not add broad member-role overrides inside ticket channels.
- Archive closed ticket transcripts into `ticket-archives`.

### 6. Staff Operations

Use this category for private internal work.

Channels:

- `staff-chat`: Coordination and escalation.
- `mod-logs`: Moderation actions, audit events, and security notes.
- `ticket-archives`: Closed ticket transcripts and summaries.
- `deploy-logs`: Fly.io, Netlify, database migration, and incident notes.

Permissions:

- Founder, Administrator, Moderator, and Support Team can view and send.
- Only Administrator and Founder should manage channels and roles.
- PulseKeep Bot can send logs and embeds.
- @everyone, Verified Member, and Server Booster cannot view.

## Interactive Command Menu

PulseKeep now includes a command-first interactive menu for Discord. It is designed to reduce support questions and make command discovery feel built into the bot.

Menu entry points:

- `/help`: Private interactive command menu.
- `/menu`: Private interactive command menu.
- `!help`: Public menu reply for legacy text-command users.
- `!menu`: Public menu reply for legacy text-command users.

Menu categories:

- Overview: Summary of every command group.
- Utility: `/ping`, `/help`, `/menu`, `/stats`, `/uptime`, `/serverinfo`, `/userinfo`, `/avatar`.
- Moderation: `/purge`, `/kick`, `/ban`, `/announce`.
- Economy: `/balance`, `/profile`, `/daily`, `/work`, `/coinflip`, `/pay`, `/leaderboard`.
- Tickets: `/ticketpanel` and the Open Ticket button.

Recommended server setup:

1. Put `/menu` in `command-menu`.
2. Put `/ticketpanel` in `open-a-ticket`.
3. Pin both bot messages.
4. Mention `/help` in the welcome/rules channel.
5. Keep sandbox testing in `bot-sandbox` so support channels stay clean.

## Permission Checklist

Before launching the server, confirm:

- Members cannot see staff channels.
- Members cannot see active tickets unless they own that ticket.
- Members can use slash commands in command and sandbox channels.
- PulseKeep Bot has Send Messages, Embed Links, Attach Files, Read Message History, and Use Slash Commands.
- PulseKeep Bot has Manage Channels only if automated ticket channel creation is enabled.
- Staff moderation commands are protected by Discord permissions.

## Launch Checklist

1. Create roles and categories.
2. Apply category permissions.
3. Create channels.
4. Invite PulseKeep with the required scopes: `bot` and `applications.commands`.
5. Run `/menu` in `command-menu`.
6. Run `/ticketpanel` in `open-a-ticket`.
7. Pin the generated bot messages.
8. Test `/help`, `/menu`, `!menu`, and the ticket button.
9. Test member visibility with a non-staff account.
10. Publish the website and link it in `rules-and-info`.

## Webhook Recommendations

Create separate webhooks for clean operational messages:

| Channel | Webhook Name | Use |
| --- | --- | --- |
| `announcements` | PulseKeep News | Releases and major updates |
| `status-logs` | PulseKeep Status | Uptime, deploys, incidents |
| `mod-logs` | PulseKeep Audit | Moderation and audit events |
| `deploy-logs` | PulseKeep Deploy | Fly.io, Netlify, and database deployment notes |

Do not post secrets, full database URLs, bot tokens, or private user data through webhooks.
