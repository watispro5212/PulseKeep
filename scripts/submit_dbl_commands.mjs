// Posts PulseKeep's slash command list to DiscordBotList
// Usage: node scripts/submit_dbl_commands.mjs
// Requires DBL_API_TOKEN env var

const commands = [
  // Utility
  { name: 'help', description: 'Show all commands filtered by category', type: 1 },
  { name: 'ping', description: 'Check bot latency', type: 1 },
  { name: 'stats', description: 'Bot statistics and metrics', type: 1 },
  { name: 'invite', description: 'Invite links and support server', type: 1 },
  { name: 'vote', description: 'Vote on DiscordBotList to earn Pulses', type: 1 },
  { name: 'servericon', description: 'Get the server icon', type: 1 },
  { name: 'roleinfo', description: 'Get role information', type: 1 },
  { name: 'channelinfo', description: 'Get channel information', type: 1 },
  { name: 'tip', description: 'Get an economy tip', type: 1 },

  // Economy
  { name: 'balance', description: 'Check your or another user\'s balance', type: 1 },
  { name: 'daily', description: 'Claim daily Pulses', type: 1 },
  { name: 'weekly', description: 'Claim weekly Pulses', type: 1 },
  { name: 'work', description: 'Work to earn Pulses', type: 1 },
  { name: 'gamble', description: 'Gamble your Pulses with a chance to win', type: 1 },
  { name: 'blackjack', description: 'Play blackjack against the dealer', type: 1 },
  { name: 'slots', description: 'Spin the slot machine', type: 1 },
  { name: 'rob', description: 'Attempt to rob another user', type: 1 },
  { name: 'pay', description: 'Send Pulses to another user', type: 1 },
  { name: 'search', description: 'Search for hidden Pulses', type: 1 },
  { name: 'fish', description: 'Go fishing (needs Fishing Rod)', type: 1 },
  { name: 'mine', description: 'Go mining (needs Mining Pick)', type: 1 },
  { name: 'shop', description: 'View the item shop', type: 1 },
  { name: 'buy', description: 'Buy an item from the shop', type: 1 },
  { name: 'inventory', description: 'Check your item inventory', type: 1 },
  { name: 'use', description: 'Use an item from your inventory', type: 1 },
  { name: 'leaderboard', description: 'View the richest users', type: 1 },

  // Moderation
  { name: 'warn', description: 'Warn a user with a reason', type: 1 },
  { name: 'warnings', description: 'Check a user\'s warnings', type: 1 },
  { name: 'clearwarns', description: 'Clear all warnings for a user', type: 1 },
  { name: 'history', description: 'View moderation history for a user', type: 1 },
  { name: 'mute', description: 'Timeout a member', type: 1 },
  { name: 'unmute', description: 'Remove timeout from a member', type: 1 },
  { name: 'kick', description: 'Kick a member from the server', type: 1 },
  { name: 'ban', description: 'Ban a user from the server', type: 1 },
  { name: 'softban', description: 'Ban then immediately unban to clear messages', type: 1 },
  { name: 'purge', description: 'Bulk delete messages', type: 1 },
  { name: 'clean', description: 'Delete messages with filters (bots, links, etc)', type: 1 },
  { name: 'slowmode', description: 'Set channel slowmode', type: 1 },
  { name: 'lock', description: 'Lock the current channel', type: 1 },
  { name: 'unlock', description: 'Unlock the current channel', type: 1 },
  { name: 'nick', description: 'Change a member\'s nickname', type: 1 },
  { name: 'role', description: 'Add or remove a role from a member', type: 1 },
  { name: 'move', description: 'Move a user to another voice channel', type: 1 },
  { name: 'vckick', description: 'Disconnect a user from voice chat', type: 1 },
  { name: 'announce', description: 'Send an announcement embed', type: 1 },

  // Tickets
  { name: 'ticketpanel', description: 'Create a ticket panel with a button', type: 1 },
  { name: 'ticket', description: 'Manage tickets (add, remove, close, rename)', type: 1 },

  // Configuration
  { name: 'configure', description: 'Configure bot settings for this server', type: 1 },
];

const BOT_ID = '1507498795569512598';
const API_TOKEN = process.env.DBL_API_TOKEN;

if (!API_TOKEN) {
  console.error('❌ DBL_API_TOKEN environment variable is not set.');
  console.error('   Set it with: $env:DBL_API_TOKEN="your-token"');
  process.exit(1);
}

console.log(`Submitting ${commands.length} commands to DiscordBotList...`);

const response = await fetch(`https://discordbotlist.com/api/v1/bots/${BOT_ID}/commands`, {
  method: 'POST',
  headers: {
    'Authorization': API_TOKEN,
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(commands),
});

const text = await response.text();
if (response.ok) {
  console.log(`✅ Successfully submitted ${commands.length} commands to DiscordBotList.`);
} else {
  console.error(`❌ Failed (${response.status}):`, text);
  process.exit(1);
}
