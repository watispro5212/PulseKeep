import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const banCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('ban')
    .setDescription('Ban a user from the server')
    .addUserOption((o) => o.setName('user').setDescription('User to ban').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for ban'))
    .addIntegerOption((o) => o.setName('days').setDescription('Days of messages to delete (0-7)').setMinValue(0).setMaxValue(7))
    .setDefaultMemberPermissions(PermissionFlagsBits.BanMembers)
    .toJSON(),

  async execute(_ctx, interaction) {
    const target = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const days = interaction.options.getInteger('days') ?? 0;
    const guildName = interaction.guild?.name ?? 'the server';

    try {
      try {
        const dm = new EmbedBuilder()
          .setTitle(`Banned from ${guildName}`)
          .setDescription(`You have been banned from **${guildName}**.`)
          .addFields({ name: 'Reason', value: reason })
          .setColor(Colors.Moderation)
          .setTimestamp();
        await target.send({ embeds: [dm] });
      } catch {}

      await interaction.guild?.members.ban(target, { reason, deleteMessageSeconds: days * 86400 });
      const emb = new EmbedBuilder()
        .setTitle('User Banned')
        .setDescription(`**${target.tag}** has been banned.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
    } catch {
      await interaction.reply({ content: '❌ Failed to ban that user.', ephemeral: true });
    }
  },
};
