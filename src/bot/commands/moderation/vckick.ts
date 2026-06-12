import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const vckickCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('vckick')
    .setDescription('Disconnect a user from voice chat')
    .addUserOption((o) => o.setName('user').setDescription('User to disconnect').setRequired(true))
    .setDefaultMemberPermissions(PermissionFlagsBits.MoveMembers)
    .toJSON(),

  async execute(_ctx, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const member = await interaction.guild.members.fetch(target.id).catch(() => null);
    if (!member?.voice.channelId) {
      await interaction.reply({ content: '❌ That user is not in a voice channel.', flags: 64 });
      return;
    }

    try {
      await member.voice.disconnect();
    } catch {
      await interaction.reply({ content: '❌ Failed to disconnect the user.', flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Member Disconnected')
      .setDescription(`Disconnected **${target.username}** from voice chat.`)
      .setColor(Colors.Moderation);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
