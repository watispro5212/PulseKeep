import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const muteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('mute')
    .setDescription('Timeout a member')
    .addUserOption((o) => o.setName('user').setDescription('User to mute').setRequired(true))
    .addIntegerOption((o) =>
      o.setName('duration').setDescription('Duration in minutes').setRequired(true).setMinValue(1).setMaxValue(40320),
    )
    .addStringOption((o) => o.setName('reason').setDescription('Reason for mute'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ bot }, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const minutes = interaction.options.getInteger('duration', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const member = await interaction.guild.members.fetch(target.id).catch(() => null);

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }
    if (!member.moderatable) {
      await interaction.reply({ content: '❌ I cannot moderate that member. Check role hierarchy.', flags: 64 });
      return;
    }

    let dmFailed = false;
    try {
      const dm = new EmbedBuilder()
        .setTitle(`Muted in ${interaction.guild.name}`)
        .setDescription(`You have been muted in **${interaction.guild.name}** for **${minutes} minute(s)**.`)
        .addFields({ name: 'Reason', value: reason })
        .setColor(Colors.Moderation)
        .setTimestamp();
      await target.send({ embeds: [dm] });
    } catch {
      dmFailed = true;
    }

    try {
      await member.timeout(minutes * 60 * 1000, reason);
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to timeout that member: ${err instanceof Error ? err.message : err}`, flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Member Muted')
      .setDescription(`**${target.username}** has been timed out for **${minutes} minute(s)**.`)
      .addFields(
        { name: 'Reason', value: reason, inline: false },
        { name: 'Duration', value: `${minutes} minute(s)`, inline: true },
        { name: 'Moderator', value: `${interaction.user}`, inline: true },
      )
      .setColor(Colors.Moderation);
    if (dmFailed) emb.setFooter({ text: '⚠️ Could not DM the user (DMs may be closed).' });
    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

    const log = new EmbedBuilder()
      .setTitle('Moderation: Mute')
      .setDescription(`**${target.username}** was muted by ${interaction.user}`)
      .addFields(
        { name: 'Reason', value: reason },
        { name: 'Duration', value: `${minutes} minutes`, inline: true },
        { name: 'User ID', value: target.id, inline: true },
      )
        .setColor(Colors.Moderation)
        .setTimestamp();
    bot.logToChannel(interaction.guildId!, log);
  },
};

export const unmuteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('unmute')
    .setDescription('Remove timeout from a member')
    .addUserOption((o) => o.setName('user').setDescription('User to unmute').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for unmute'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ bot }, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const member = await interaction.guild.members.fetch(target.id).catch(() => null);

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }

    let dmFailed = false;
    try {
      const dm = new EmbedBuilder()
        .setTitle(`Unmuted in ${interaction.guild.name}`)
        .setDescription(`You have been unmuted in **${interaction.guild.name}**.`)
        .addFields({ name: 'Reason', value: reason })
        .setColor(Colors.Moderation)
        .setTimestamp();
      await target.send({ embeds: [dm] });
    } catch {
      dmFailed = true;
    }

    try {
      await member.timeout(null, reason);
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to unmute that member: ${err instanceof Error ? err.message : err}`, flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Member Unmuted')
      .setDescription(`**${target.username}** has been unmuted.`)
      .addFields(
        { name: 'Reason', value: reason, inline: false },
        { name: 'Moderator', value: `${interaction.user}`, inline: true },
      )
      .setColor(Colors.Moderation);
    if (dmFailed) emb.setFooter({ text: '⚠️ Could not DM the user (DMs may be closed).' });
    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

    const log = new EmbedBuilder()
      .setTitle('Moderation: Unmute')
      .setDescription(`**${target.username}** was unmuted by ${interaction.user}`)
      .addFields(
        { name: 'Reason', value: reason },
        { name: 'User ID', value: target.id, inline: true },
      )
        .setColor(Colors.Moderation)
        .setTimestamp();
    await bot.logToChannel(interaction.guildId!, log);
  },
};
