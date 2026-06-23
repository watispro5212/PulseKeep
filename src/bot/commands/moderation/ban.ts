import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ComponentType,
  type MessageComponentInteraction,
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

    const member = await interaction.guild.members.fetch(target.id).catch(() => null);
    if (member && !member.bannable) {
      await interaction.reply({ content: '❌ I cannot ban that user. Check role hierarchy.', flags: 64 });
      return;
    }

    const confirmId = 'confirm_ban';
    const cancelId = 'cancel_ban';
    const confirm = new ButtonBuilder()
      .setCustomId(confirmId)
      .setLabel(`Ban ${target.username}`)
      .setStyle(ButtonStyle.Danger);
    const cancel = new ButtonBuilder()
      .setCustomId(cancelId)
      .setLabel('Cancel')
      .setStyle(ButtonStyle.Secondary);
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(confirm, cancel);

    const prompt = new EmbedBuilder()
      .setTitle('⚠️ Confirm Ban')
      .setDescription(`Ban **${target.username}** (\`${target.id}\`)?\n**Reason:** ${reason}${days > 0 ? `\n**Delete messages from last ${days} day(s)**` : ''}`)
      .setColor(Colors.Warning);

    await interaction.reply({ embeds: [footer(timestamp(prompt))], components: [row], flags: 64 });

    const reply = await interaction.fetchReply();
    const collected = await reply.awaitMessageComponent({
      componentType: ComponentType.Button,
      time: 15000,
      filter: (i: MessageComponentInteraction) => i.user.id === interaction.user.id,
    }).catch(() => null);

    if (!collected || collected.customId === cancelId) {
      await interaction.editReply({ content: '❌ Cancelled.', embeds: [], components: [] });
      return;
    }

    await collected.deferUpdate();

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
      await interaction.editReply({ content: `❌ Failed to ban that user: ${err instanceof Error ? err.message : err}`, components: [] });
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
    await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });

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
    await bot.logToChannel(interaction.guildId!, log);
  },
};
