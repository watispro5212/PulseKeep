# PulseKeep DiscordBotList Submission Packet

Use this as the source copy when submitting PulseKeep to DiscordBotList.

## Core Details

- **Bot name:** PulseKeep
- **Website:** https://pulsekeep.williamdelilah3.workers.dev
- **Command page:** https://pulsekeep.williamdelilah3.workers.dev/commands.html
- **Status page:** https://pulsekeep.williamdelilah3.workers.dev/status.html
- **Privacy policy:** https://pulsekeep.williamdelilah3.workers.dev/privacy.html
- **Terms of service:** https://pulsekeep.williamdelilah3.workers.dev/terms.html
- **Prefix:** Slash commands
- **Primary command for reviewers:** `/help`
- **Backup test commands:** `/ping`, `/stats`

## Short Description

PulseKeep is a Discord bot for moderation, audit logs, private support tickets, economy commands, and live server analytics.

## Long Description

PulseKeep helps Discord staff teams keep servers organized with practical slash commands.

## Full Command List

### Moderation
| Command | Description |
|---------|-------------|
| `/ban` | Ban a user from the server |
| `/kick` | Kick a member |
| `/mute` | Timeout a member |
| `/unmute` | Remove timeout from a member |
| `/warn` | Warn a user with a reason |
| `/warnings` | Check a user's warnings |
| `/clearwarns` | Clear all warnings for a user |
| `/history` | View moderation history |
| `/softban` | Ban + unban to clear messages |
| `/purge` | Bulk delete messages |
| `/clean` | Delete messages with filters (bots, links, etc.) |
| `/lock` | Lock a channel |
| `/unlock` | Unlock a channel |
| `/slowmode` | Set channel slowmode |
| `/move` | Move user to another voice channel |
| `/vckick` | Disconnect user from voice channel |
| `/nick` | Change a member's nickname |
| `/role add` | Add a role to a member |
| `/role remove` | Remove a role from a member |
| `/announce` | Send an announcement embed |

### Economy
| Command | Description |
|---------|-------------|
| `/balance` | Check your or another user's balance |
| `/daily` | Claim daily Pulses |
| `/weekly` | Claim weekly Pulses |
| `/work` | Work to earn Pulses |
| `/gamble` | Gamble your Pulses (2x-10x win, 55% win rate) |
| `/blackjack` | Play blackjack against the dealer |
| `/slots` | Spin the slot machine |
| `/rob` | Attempt to rob another user |
| `/pay` | Send Pulses to another user |
| `/search` | Search for hidden Pulses |
| `/fish` | Go fishing (needs Fishing Rod) |
| `/mine` | Go mining (needs Mining Pick) |
| `/shop` | View the shop |
| `/buy` | Buy an item from the shop |
| `/use` | Use an item from your inventory |
| `/inventory` | Check your inventory |
| `/leaderboard` | View richest users |
| `/tip` | Get an economy tip |
| `/vote` | Vote on DiscordBotList and earn Pulses |

### Tickets
| Command | Description |
|---------|-------------|
| `/ticketpanel` | Create a ticket panel with button |
| `/ticket add` | Add user to a ticket |
| `/ticket remove` | Remove user from a ticket |
| `/ticket close` | Close the current ticket |
| `/ticket rename` | Rename the ticket channel |

### Configuration
| Command | Description |
|---------|-------------|
| `/configure economy` | Toggle economy system |
| `/configure tickets` | Toggle ticket system |
| `/configure modlogs` | Toggle moderation logging |
| `/configure welcome` | Toggle welcome messages |
| `/configure log_channel` | Set moderation log channel |
| `/configure welcome_channel` | Set welcome message channel |
| `/configure vote_channel` | Set vote announcement channel |
| `/configure ticket_category` | Set ticket category |
| `/configure show` | View current configuration |

### Utility
| Command | Description |
|---------|-------------|
| `/help` | Show all commands |
| `/ping` | Check bot latency |
| `/stats` | Bot statistics |
| `/invite` | Invite links and support |
| `/vote` | Vote for PulseKeep |
| `/servericon` | Get server icon |
| `/roleinfo` | Role information |
| `/channelinfo` | Channel information |

### Total: 50+ slash commands

## Categories

- Moderation
- Economy
- Tickets
- Configuration
- Utility
- Logging

## Webhook Setup

Set your DBL webhook in DiscordBotList dashboard → Bot Settings → Webhooks:
- **URL:** `https://pulsekeep.fly.dev/api/dbl/webhook`
- **Authorization header:** Your chosen secret (must match `DBL_WEBHOOK_SECRET` env var)

When a user votes, PulseKeep automatically:
1. Credits 500-750 Pulses to the voter
2. Posts a celebratory embed in all shared guilds with a `voteChannelId` configured
