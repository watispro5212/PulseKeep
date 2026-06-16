import {
  SlashCommandBuilder,
  EmbedBuilder,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, Ephemeral, timestamp } from '../../utils/embed.js';

// Minimal permission set — everything PulseKeep actually needs.
// (Not 8/Administrator — that scares people.)
const RECOMMENDED_PERMS = [
  PermissionFlagsBits.ViewChannel,
  PermissionFlagsBits.SendMessages,
  PermissionFlagsBits.EmbedLinks,
  PermissionFlagsBits.AttachFiles,
  PermissionFlagsBits.ReadMessageHistory,
  PermissionFlagsBits.ManageMessages,
  PermissionFlagsBits.ManageChannels,
  PermissionFlagsBits.ManageRoles,
  PermissionFlagsBits.KickMembers,
  PermissionFlagsBits.BanMembers,
  PermissionFlagsBits.ModerateMembers,
  PermissionFlagsBits.MentionEveryone,
  PermissionFlagsBits.Connect,
  PermissionFlagsBits.Speak,
  PermissionFlagsBits.MoveMembers,
  PermissionFlagsBits.UseExternalEmojis,
  PermissionFlagsBits.AddReactions,
  PermissionFlagsBits.UseApplicationCommands,
].reduce((acc, p) => acc | p, 0n).toString();

const ADMIN_PERMS = PermissionFlagsBits.Administrator.toString();

const BOT_ID = '1507498795569512598';
const inviteRecommended = `https://discord.com/oauth2/authorize?client_id=${BOT_ID}&permissions=${RECOMMENDED_PERMS}&scope=bot%20applications.commands`;
const inviteAdmin = `https://discord.com/oauth2/authorize?client_id=${BOT_ID}&permissions=${ADMIN_PERMS}&scope=bot%20applications.commands`;

export const inviteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('invite')
    .setDescription('Get invite links for PulseKeep')
    .toJSON(),

  async execute(_ctx, interaction) {
    const emb = new EmbedBuilder()
      .setTitle('Invite PulseKeep')
      .setDescription(
        'Add PulseKeep to your server, or join the community. The **Recommended** link requests the ' +
        'minimum permissions PulseKeep actually uses — pick **Administrator** only if you want the ' +
        'full feature set (lockdown, slowmode, full automod) without re-granting perms later.'
      )
      .setColor(Colors.Utility)
      .addFields(
        {
          name: '🔗 Recommended permissions',
          value: `Picks the ~18 perms PulseKeep uses. **Most servers want this.**\n${inviteRecommended}`,
          inline: false,
        },
        {
          name: '🔐 Administrator (full access)',
          value: `Grants every permission. Use only on servers where you fully trust the bot.\n${inviteAdmin}`,
          inline: false,
        },
        { name: '💬 Support Server', value: 'https://discord.gg/b9HBphyeuP', inline: true },
        { name: '⭐ DiscordBotList', value: 'https://discordbotlist.com/bots/' + BOT_ID, inline: true },
        { name: '🤖 Discords.com', value: 'https://discords.com/bots/bot/' + BOT_ID, inline: true },
        { name: '🌐 Website', value: 'https://pulsekeep.fly.dev', inline: true },
        { name: '📦 GitHub', value: 'https://github.com/watispro5212/PulseKeep', inline: true },
      )
      .setFooter({ text: 'Tip: vote on DBL or Discords.com once a day for free Pulses.' });

    const row = new ActionRowBuilder().addComponents(
      new ButtonBuilder().setStyle(ButtonStyle.Link).setURL(inviteRecommended).setLabel('Recommended perms').setEmoji('🔗'),
      new ButtonBuilder().setStyle(ButtonStyle.Link).setURL(inviteAdmin).setLabel('Administrator').setEmoji('🔐'),
      new ButtonBuilder().setStyle(ButtonStyle.Link).setURL('https://discord.gg/b9HBphyeuP').setLabel('Support').setEmoji('💬'),
    );

    await interaction.reply({ embeds: [timestamp(emb)], components: [row], flags: Ephemeral });
  },
};
