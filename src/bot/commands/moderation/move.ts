import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const moveCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('move')
    .setDescription('Move a user to another voice channel')
    .addUserOption((o) => o.setName('user').setDescription('User to move').setRequired(true))
    .addChannelOption((o) =>
      o.setName('channel').setDescription('Target voice channel').setRequired(true),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.MoveMembers)
    .toJSON(),

  async execute(_ctx, interaction) {
    const target = interaction.options.getUser('user', true);
    const channel = interaction.options.getChannel('channel', true);

    if (channel.type !== 2) {
      await interaction.reply({ content: '❌ Please select a voice channel.', ephemeral: true });
      return;
    }

    const member = interaction.guild?.members.cache.get(target.id);
    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', ephemeral: true });
      return;
    }

    const oldChannelId = member.voice.channelId;
    try {
      await member.voice.setChannel(channel.id);
    } catch {
      await interaction.reply({ content: '❌ Failed to move the member.', ephemeral: true });
      return;
    }

    const oldChannel = oldChannelId ? `<#${oldChannelId}>` : 'nowhere (not connected)';

    const emb = new EmbedBuilder()
      .setTitle('Member Moved')
      .setDescription(`Moved **${target.tag}** to ${channel}.`)
      .addFields(
        { name: 'From', value: oldChannel, inline: true },
        { name: 'To', value: `${channel}`, inline: true },
      )
      .setColor(Colors.Moderation);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
