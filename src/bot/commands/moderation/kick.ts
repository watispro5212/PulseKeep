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

export const kickCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('kick')
    .setDescription('Kick a member from the server')
    .addUserOption((o) => o.setName('user').setDescription('User to kick').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for kick'))
    .setDefaultMemberPermissions(PermissionFlagsBits.KickMembers)
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
    if (!member.kickable) {
      await interaction.reply({ content: '❌ I cannot kick that member. Check role hierarchy.', flags: 64 });
      return;
    }

    const confirmId = 'confirm_kick';
    const cancelId = 'cancel_kick';
    const confirm = new ButtonBuilder()
      .setCustomId(confirmId)
      .setLabel(`Kick ${target.username}`)
      .setStyle(ButtonStyle.Danger);
    const cancel = new ButtonBuilder()
      .setCustomId(cancelId)
      .setLabel('Cancel')
      .setStyle(ButtonStyle.Secondary);
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(confirm, cancel);

    const prompt = new EmbedBuilder()
      .setTitle('⚠️ Confirm Kick')
      .setDescription(`Kick **${target.username}** (\`${target.id}\`)?\n**Reason:** ${reason}`)
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
        .setTitle(`Kicked from ${interaction.guild.name}`)
        .setDescription(`You have been kicked from **${interaction.guild.name}**.`)
        .addFields({ name: 'Reason', value: reason })
        .setColor(Colors.Moderation)
        .setTimestamp();
      await target.send({ embeds: [dm] });
    } catch {
      dmFailed = true;
    }

    try {
      await member.kick(reason);
    } catch (err) {
      await interaction.editReply({ content: `❌ Failed to kick that member: ${err instanceof Error ? err.message : err}`, components: [] });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Member Kicked')
      .setDescription(`**${target.username}** has been kicked.`)
      .addFields(
        { name: 'Reason', value: reason, inline: false },
        { name: 'Moderator', value: `${interaction.user}`, inline: true },
      )
      .setColor(Colors.Moderation);
    if (dmFailed) emb.setFooter({ text: '⚠️ Could not DM the user (DMs may be closed).' });
    await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });

    const log = new EmbedBuilder()
      .setTitle('Moderation: Kick')
      .setDescription(`**${target.username}** was kicked by ${interaction.user}`)
      .addFields(
        { name: 'Reason', value: reason },
        { name: 'User ID', value: target.id, inline: true },
      )
        .setColor(Colors.Moderation)
        .setTimestamp();
    await bot.logToChannel(interaction.guildId!, log);
  },
};
