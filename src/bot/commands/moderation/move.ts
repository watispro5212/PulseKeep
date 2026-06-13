import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
  ChannelType,
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
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const channel = interaction.options.getChannel('channel', true);

    if (channel.type !== ChannelType.GuildVoice) {
      await interaction.reply({ content: '❌ Please select a voice channel.', flags: 64 });
      return;
    }

    const member = await interaction.guild.members.fetch(target.id).catch(() => null);
    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }

    const oldChannelId = member.voice.channelId;
    try {
      await member.voice.setChannel(channel.id);
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to move the member: ${err instanceof Error ? err.message : err}`, flags: 64 });
      return;
    }

    const oldChannel = oldChannelId ? `<#${oldChannelId}>` : 'nowhere (not connected)';

    const emb = new EmbedBuilder()
      .setTitle('Member Moved')
      .setDescription(`Moved **${target.username}** to ${channel}.`)
      .addFields(
        { name: 'From', value: oldChannel, inline: true },
        { name: 'To', value: `${channel}`, inline: true },
      )
      .setColor(Colors.Moderation);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
