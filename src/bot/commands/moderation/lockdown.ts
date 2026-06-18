import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const lockCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('lock')
    .setDescription('Lock this channel (prevent @everyone from sending)')
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageChannels)
    .toJSON(),

  async execute(_ctx, interaction) {
    const channel = interaction.channel;
    if (!channel || !('permissionOverwrites' in channel)) {
      await interaction.reply({ content: 'This channel cannot be locked.', flags: 64 });
      return;
    }

    try {
      await channel.permissionOverwrites.edit(interaction.guild.roles.everyone, {
        SendMessages: false,
      });
      const emb = new EmbedBuilder()
        .setTitle('Channel Locked')
        .setDescription(`🔒 ${channel} has been locked.`)
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to lock channel: ${err instanceof Error ? err.message : err}`, flags: 64 });
    }
  },
};

export const unlockCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('unlock')
    .setDescription('Unlock this channel')
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageChannels)
    .toJSON(),

  async execute(_ctx, interaction) {
    const channel = interaction.channel;
    if (!channel || !('permissionOverwrites' in channel)) {
      await interaction.reply({ content: 'This channel cannot be unlocked.', flags: 64 });
      return;
    }

    try {
      await channel.permissionOverwrites.edit(interaction.guild.roles.everyone, {
        SendMessages: null,
      });
      const emb = new EmbedBuilder()
        .setTitle('Channel Unlocked')
        .setDescription(`🔓 ${channel} has been unlocked.`)
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to unlock channel: ${err instanceof Error ? err.message : err}`, flags: 64 });
    }
  },
};
