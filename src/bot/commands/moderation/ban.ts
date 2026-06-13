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

  async execute({ bot }, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const days = interaction.options.getInteger('days') ?? 0;

    let dmFailed = false;
    try {
      const dm = new EmbedBuilder()
        .setTitle(`Banned from ${interaction.guild.name}`)
        .setDescription(`You have been banned from **${interaction.guild.name}**.`)
        .addFields({ name: 'Reason', value: reason })
        .setColor(Colors.Moderation)
        .setTimestamp();
      await target.send({ embeds: [dm] });
    } catch {
      dmFailed = true;
    }

    try {
      await interaction.guild.members.ban(target, { reason, deleteMessageSeconds: days * 86400 });
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to ban that user: ${err instanceof Error ? err.message : err}`, flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('User Banned')
      .setDescription(`**${target.username}** has been banned.`)
      .addFields(
        { name: 'Reason', value: reason, inline: false },
        { name: 'Moderator', value: `${interaction.user}`, inline: true },
      )
      .setColor(Colors.Moderation);
    if (dmFailed) emb.setFooter({ text: '⚠️ Could not DM the user (DMs may be closed).' });
    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

    const log = new EmbedBuilder()
      .setTitle('Moderation: Ban')
      .setDescription(`**${target.username}** was banned by ${interaction.user}`)
      .addFields(
        { name: 'Reason', value: reason },
        { name: 'User ID', value: target.id, inline: true },
        { name: 'Deleted Days', value: `${days}`, inline: true },
      )
      .setColor(Colors.Moderation)
      .setTimestamp();
    bot.logToChannel(interaction.guildId!, log);
  },
};
