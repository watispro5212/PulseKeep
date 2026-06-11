import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, footer, timestamp } from '../../utils/embed.js';

const categories = [
  {
    name: '🎮 Economy',
    accent: Colors.Economy,
    commands: [
      { name: 'balance', desc: 'Check your or another user\'s balance' },
      { name: 'daily', desc: 'Claim daily Pulses' },
      { name: 'weekly', desc: 'Claim weekly Pulses' },
      { name: 'work', desc: 'Work to earn Pulses' },
      { name: 'gamble', desc: 'Gamble your Pulses' },
      { name: 'blackjack', desc: 'Play blackjack' },
      { name: 'slots', desc: 'Spin the slot machine' },
      { name: 'rob', desc: 'Rob another user' },
      { name: 'pay', desc: 'Send Pulses to a user' },
      { name: 'fish', desc: 'Go fishing' },
      { name: 'mine', desc: 'Go mining' },
      { name: 'shop', desc: 'View the shop' },
      { name: 'buy', desc: 'Buy an item' },
      { name: 'inventory', desc: 'Check your inventory' },
      { name: 'leaderboard', desc: 'Richest users' },
      { name: 'tip', desc: 'Get an economy tip' },
    ],
  },
  {
    name: '🛡️ Moderation',
    accent: Colors.Moderation,
    commands: [
      { name: 'warn', desc: 'Warn a user' },
      { name: 'warnings', desc: 'Check warnings' },
      { name: 'clearwarns', desc: 'Clear warnings' },
      { name: 'mute', desc: 'Timeout a member' },
      { name: 'unmute', desc: 'Remove timeout' },
      { name: 'kick', desc: 'Kick a member' },
      { name: 'ban', desc: 'Ban a user' },
      { name: 'softban', desc: 'Ban + unban to clear messages' },
      { name: 'purge', desc: 'Bulk delete messages' },
      { name: 'clean', desc: 'Delete messages with filters' },
      { name: 'move', desc: 'Move user to another VC' },
      { name: 'vckick', desc: 'Disconnect from VC' },
      { name: 'nick', desc: 'Change a nickname' },
      { name: 'role add', desc: 'Add role to member' },
      { name: 'role remove', desc: 'Remove role from member' },
      { name: 'slowmode', desc: 'Set channel slowmode' },
      { name: 'lock', desc: 'Lock a channel' },
      { name: 'unlock', desc: 'Unlock a channel' },
      { name: 'announce', desc: 'Send an announcement' },
      { name: 'history', desc: 'View moderation history' },
    ],
  },
  {
    name: '🎫 Tickets',
    accent: Colors.Tickets,
    commands: [
      { name: 'ticketpanel', desc: 'Create a ticket panel' },
      { name: 'ticket add', desc: 'Add user to ticket' },
      { name: 'ticket remove', desc: 'Remove user from ticket' },
      { name: 'ticket close', desc: 'Close the ticket' },
      { name: 'ticket rename', desc: 'Rename the ticket' },
    ],
  },
  {
    name: '🔧 Utility',
    accent: Colors.Utility,
    commands: [
      { name: 'ping', desc: 'Check bot latency' },
      { name: 'stats', desc: 'Bot statistics' },
      { name: 'invite', desc: 'Invite links' },
      { name: 'vote', desc: 'Vote on Top.gg' },
      { name: 'servericon', desc: 'Get server icon' },
      { name: 'roleinfo', desc: 'Role information' },
      { name: 'channelinfo', desc: 'Channel information' },
    ],
  },
];

export const helpCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('help')
    .setDescription('Show all commands')
    .addStringOption((o) =>
      o.setName('category').setDescription('Filter by category')
        .addChoices(
          { name: 'Economy', value: 'economy' },
          { name: 'Moderation', value: 'moderation' },
          { name: 'Tickets', value: 'tickets' },
          { name: 'Utility', value: 'utility' },
        ),
    )
    .toJSON(),

  async execute(_ctx, interaction) {
    const filter = interaction.options.getString('category');

    const emb = new EmbedBuilder()
      .setTitle('PulseKeep Commands')
      .setDescription('Use `/help <category>` to filter by category.')
      .setColor(Colors.Utility);

    for (const cat of categories) {
      if (filter && cat.name.toLowerCase().replace(/[^a-z]/g, '') !== filter) continue;
      const cmds = cat.commands.map((c) => `\`/${c.name}\` — ${c.desc}`).join('\n');
      emb.addFields({ name: cat.name, value: cmds, inline: false });
    }

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
