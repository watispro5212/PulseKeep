import {
  SlashCommandBuilder,
  EmbedBuilder,
  ActionRowBuilder,
  StringSelectMenuBuilder,
  StringSelectMenuInteraction,
  ComponentType,
} from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, footer, timestamp, Ephemeral } from '../../utils/embed.js';

interface CommandInfo {
  name: string;
  desc: string;
  usage?: string;
}

interface CommandCategory {
  id: string;
  label: string;
  emoji: string;
  color: number;
  blurb: string;
  commands: CommandInfo[];
}

const categories: CommandCategory[] = [
  {
    id: 'economy',
    label: 'Economy',
    emoji: '💸',
    color: Colors.Economy,
    blurb: 'Earn Pulses, play games, buy items, climb the leaderboard.',
    commands: [
      { name: 'balance', desc: "Check your or another user's Pulses balance.", usage: '/balance [user]' },
      { name: 'daily', desc: 'Claim your 24h daily reward — streak bonuses at 7, 14, 21, 30, 50, 100 days.', usage: '/daily [public]' },
      { name: 'weekly', desc: 'Claim a bigger weekly reward.', usage: '/weekly' },
      { name: 'work', desc: 'Pick a job and earn 80–800 Pulses/hour.', usage: '/work' },
      { name: 'gamble', desc: 'Risk Pulses on a 55%-win coin flip. 2x–10x multipliers.', usage: '/gamble <amount>' },
      { name: 'blackjack', desc: 'Beat the dealer in classic 21.', usage: '/blackjack <amount>' },
      { name: 'slots', desc: 'Spin for tiered payouts and 10% upgrade chance.', usage: '/slots <amount>' },
      { name: 'rob', desc: 'Steal from another user — 45% success rate, 4h cooldown.', usage: '/rob <user>' },
      { name: 'pay', desc: 'Send Pulses to a friend.', usage: '/pay <user> <amount>' },
      { name: 'tip', desc: 'Show a random economy tip.' },
      { name: 'fish', desc: 'Catch fish for Pulses (needs a Fishing Rod).', usage: '/fish' },
      { name: 'mine', desc: 'Mine for ore and gemstones (needs a Pick).', usage: '/mine' },
      { name: 'search', desc: 'Hunt for hidden Pulses in 12 locations. 15m cooldown.' },
      { name: 'shop', desc: 'List shop items — rods, picks, boosts, clovers, maps.' },
      { name: 'buy', desc: 'Buy a shop item.', usage: '/buy <item>' },
      { name: 'inventory', desc: 'Show your items, active XP Boost, and Lucky Clover.' },
      { name: 'use', desc: 'Use an item from your inventory.', usage: '/use <item>' },
      { name: 'vote', desc: 'Vote on DiscordBotList / Discords.com for Pulses.' },
      { name: 'leaderboard', desc: 'Top 10 richest users in this server.', usage: '/leaderboard [page]' },
    ],
  },
  {
    id: 'moderation',
    label: 'Moderation',
    emoji: '🛡️',
    color: Colors.Moderation,
    blurb: 'Keep the peace. Every action is logged and can be DMed to the user.',
    commands: [
      { name: 'warn', desc: 'Issue a warning. Shows in /history and /warnings.', usage: '/warn <user> <reason>' },
      { name: 'warnings', desc: "Show a user's warning history.", usage: '/warnings <user>' },
      { name: 'clearwarns', desc: 'Clear all warnings for a user (confirmation buttons).', usage: '/clearwarns <user>' },
      { name: 'history', desc: "View a user's full mod history (warns, kicks, bans).", usage: '/history <user>' },
      { name: 'mute', desc: 'Timeout a member (max 28 days).', usage: '/mute <user> <duration> [reason]' },
      { name: 'unmute', desc: 'Lift an active timeout.', usage: '/unmute <user>' },
      { name: 'kick', desc: 'Kick a member. Logs the reason.', usage: '/kick <user> [reason]' },
      { name: 'ban', desc: 'Ban a user. Optionally delete recent messages.', usage: '/ban <user> [days] [reason]' },
      { name: 'softban', desc: 'Ban + unban to clear recent messages.', usage: '/softban <user> [reason]' },
      { name: 'purge', desc: 'Bulk delete 1–100 messages.', usage: '/purge <count>' },
      { name: 'clean', desc: 'Delete messages matching a filter (bots, links, attachments).', usage: '/clean <count> <filter>' },
      { name: 'slowmode', desc: 'Set channel slowmode (0–21600s).', usage: '/slowmode <seconds>' },
      { name: 'lock / unlock', desc: 'Lock or unlock a channel (@everyone can\'t send).' },
      { name: 'nick', desc: "Change a member's nickname.", usage: '/nick <user> <name>' },
      { name: 'role add / remove', desc: 'Add or remove a role from a member.', usage: '/role <add|remove> <user> <role>' },
      { name: 'move', desc: 'Move a member to another voice channel.', usage: '/move <user> <channel>' },
      { name: 'vckick', desc: 'Disconnect a member from voice.', usage: '/vckick <user>' },
      { name: 'announce', desc: 'Send an embed announcement to a channel.', usage: '/announce <channel> <title> <body>' },
    ],
  },
  {
    id: 'tickets',
    label: 'Tickets',
    emoji: '🎫',
    color: Colors.Tickets,
    blurb: 'Button-based support panels with private channels.',
    commands: [
      { name: 'ticketpanel', desc: 'Post a ticket-creation panel in the current channel.' },
      { name: 'ticket add', desc: 'Add a user to the current ticket.', usage: '/ticket add <user>' },
      { name: 'ticket remove', desc: 'Remove a user from the current ticket.', usage: '/ticket remove <user>' },
      { name: 'ticket close', desc: 'Archive and lock the current ticket.' },
      { name: 'ticket rename', desc: 'Rename the current ticket channel.', usage: '/ticket rename <name>' },
    ],
  },
  {
    id: 'config',
    label: 'Configuration',
    emoji: '⚙️',
    color: Colors.Configure,
    blurb: 'Toggle features and pick channels. Admin only.',
    commands: [
      { name: 'configure show', desc: 'Show the current server config.' },
      { name: 'configure economy <on|off>', desc: 'Enable or disable the economy system.' },
      { name: 'configure tickets <on|off>', desc: 'Enable or disable tickets.' },
      { name: 'configure modlogs <on|off>', desc: 'Enable or disable moderation logging.' },
      { name: 'configure welcome <on|off>', desc: 'Enable or disable welcome messages.' },
      { name: 'configure welcome_channel', desc: 'Set the channel new-member welcomes go to.' },
      { name: 'configure log_channel', desc: 'Set the channel mod actions log to.' },
      { name: 'configure vote_channel', desc: 'Set the channel that vote announcements go to.' },
      { name: 'configure ticket_category', desc: 'Set the category new tickets spawn in.' },
    ],
  },
  {
    id: 'utility',
    label: 'Utility',
    emoji: '🔧',
    color: Colors.Utility,
    blurb: 'Lookups, stats, and bot info.',
    commands: [
      { name: 'help', desc: 'Open this menu.' },
      { name: 'ping', desc: 'Round-trip + WebSocket latency.' },
      { name: 'stats', desc: 'Servers, users, commands run, uptime, version, shards.' },
      { name: 'invite', desc: 'Bot invite + support + DBL/Discords.com vote links.' },
      { name: 'about', desc: 'About PulseKeep — version, license, links.' },
      { name: 'userinfo', desc: "Inspect a user: avatar, joined, roles, permissions.", usage: '/userinfo [user]' },
      { name: 'serverinfo', desc: 'Server details: members, channels, roles, boosts, verification.' },
      { name: 'roleinfo', desc: 'Inspect a role: color, permissions, position.', usage: '/roleinfo <role>' },
      { name: 'channelinfo', desc: 'Inspect a channel: type, created, NSFW.', usage: '/channelinfo [channel]' },
      { name: 'servericon', desc: 'High-res version of the server icon.' },
    ],
  },
];

function buildCategoryEmbed(cat: CommandCategory, page = 0): EmbedBuilder {
  const total = cat.commands.length;
  const perPage = 12;
  const slice = cat.commands.slice(page * perPage, (page + 1) * perPage);
  const lines = slice.map((c) => {
    const cmd = `\`/${c.name}\``;
    return `${cmd} — ${c.desc}` + (c.usage ? `\n  ↳ \`${c.usage}\`` : '');
  });

  const emb = new EmbedBuilder()
    .setTitle(`${cat.emoji} ${cat.label} commands`)
    .setDescription(cat.blurb + (total > perPage ? `\n_Page ${page + 1} of ${Math.ceil(total / perPage)} — ${total} total_` : `\n_${total} commands total_`))
    .setColor(cat.color)
    .addFields({ name: '\u200b', value: lines.join('\n\n') })
    .setFooter({ text: 'PulseKeep v7.0.0 • /help <category> to filter' });

  return emb;
}

function buildOverviewEmbed(): EmbedBuilder {
  const lines = categories
    .map((c) => `${c.emoji} **${c.label}** — ${c.commands.length} commands\n  _${c.blurb}_`)
    .join('\n\n');
  return new EmbedBuilder()
    .setTitle('PulseKeep — Command Center')
    .setDescription(
      `Welcome! Use the dropdown below to jump to a category, or type \`/help <category>\`.\n\n${lines}\n\n` +
      `**Tip:** every command is a slash command — start typing \`/\` in any channel to see them all.`
    )
    .setColor(Colors.Utility)
    .setFooter({ text: 'PulseKeep v7.0.0 • Select a category below' });
}

export const helpCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('help')
    .setDescription('Show all commands, or filter by category')
    .addStringOption((o) =>
      o.setName('category').setDescription('Filter by category')
        .addChoices(
          { name: 'Economy', value: 'economy' },
          { name: 'Moderation', value: 'moderation' },
          { name: 'Tickets', value: 'tickets' },
          { name: 'Configuration', value: 'configuration' },
          { name: 'Utility', value: 'utility' },
          { name: 'Overview', value: 'overview' },
        ),
    )
    .toJSON(),

  async execute({}, interaction) {
    const filter = interaction.options.getString('category');

    if (filter) {
      if (filter === 'overview' || filter === 'configuration') {
        const emb = filter === 'overview' ? buildOverviewEmbed() : buildCategoryEmbed(categories.find(c => c.id === 'config') ?? categories[0]!);
        await interaction.reply({ embeds: [timestamp(emb)], flags: Ephemeral });
        return;
      }
      const cat = categories.find((c) => c.id === filter);
      if (!cat) {
        await interaction.reply({ content: 'Unknown category.', flags: Ephemeral });
        return;
      }
      await interaction.reply({ embeds: [timestamp(buildCategoryEmbed(cat))], flags: Ephemeral });
      return;
    }

    // dropdown menu
    const select = new StringSelectMenuBuilder()
      .setCustomId('help_category')
      .setPlaceholder('Pick a category…')
      .addOptions(
        categories.map((c) => ({
          label: `${c.label} (${c.commands.length})`,
          description: c.blurb.slice(0, 100),
          value: c.id,
          emoji: c.emoji,
        })),
      );

    const row = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(select);

    const reply = await interaction.reply({
      embeds: [timestamp(buildOverviewEmbed())],
      components: [row],
      flags: Ephemeral,
      fetchReply: true,
    });

    const collector = reply.createMessageComponentCollector({
      componentType: ComponentType.StringSelect,
      time: 120_000,
    });

    collector.on('collect', async (i: StringSelectMenuInteraction) => {
      if (i.user.id !== interaction.user.id) {
        await i.reply({ content: 'This menu isn\'t yours — run `/help` to get your own.', flags: Ephemeral });
        return;
      }
      const cat = categories.find((c) => c.id === i.values[0]);
      if (!cat) return;
      await i.update({ embeds: [timestamp(buildCategoryEmbed(cat))], components: [row] });
    });

    collector.on('end', () => {
      // best-effort disable (just in case)
    });
  },
};
