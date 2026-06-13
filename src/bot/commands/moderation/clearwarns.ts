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
import { Colors, footer, timestamp, Ephemeral } from '../../../utils/embed.js';
import { userWarnings } from '../../../db/schema.js';
import { eq, and } from 'drizzle-orm';

export const clearwarnsCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('clearwarns')
    .setDescription('Clear all warnings for a user')
    .addUserOption((o) => o.setName('user').setDescription('User to clear').setRequired(true))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ bot, db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: Ephemeral });
      return;
    }

    const target = interaction.options.getUser('user', true);

    const confirm = new ButtonBuilder()
      .setCustomId('confirm_clear')
      .setLabel('Yes, clear all warnings')
      .setStyle(ButtonStyle.Danger);
    const cancel = new ButtonBuilder()
      .setCustomId('cancel_clear')
      .setLabel('Cancel')
      .setStyle(ButtonStyle.Secondary);
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(confirm, cancel);

    const prompt = new EmbedBuilder()
      .setTitle('⚠️ Confirm Clear Warnings')
      .setDescription(`Are you sure you want to clear **all** warnings for ${target}? This cannot be undone.`)
      .setColor(Colors.Warning);
    await interaction.reply({ embeds: [footer(timestamp(prompt))], components: [row], flags: Ephemeral });

    const reply = await interaction.fetchReply();
    const collected = await reply.awaitMessageComponent({
      componentType: ComponentType.Button,
      time: 30000,
      filter: (i: MessageComponentInteraction) => i.user.id === interaction.user.id,
    }).catch(() => null);

    if (!collected || collected.customId === 'cancel_clear') {
      await interaction.editReply({ content: '❌ Cancelled.', embeds: [], components: [] });
      return;
    }

    await collected.deferUpdate();

    let count = 0;
    try {
      const result = await db
        .delete(userWarnings)
        .where(and(eq(userWarnings.guildId, interaction.guildId!), eq(userWarnings.userId, target.id)));
      count = (result as any)?.rowCount ?? 0;
    } catch (err) {
      await interaction.editReply({ content: `❌ Failed to clear warnings: ${err instanceof Error ? err.message : err}`, components: [] });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Warnings Cleared')
      .setDescription(`Cleared **${count}** warning(s) for ${target}.`)
      .setColor(Colors.Moderation);

    await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });

    const log = new EmbedBuilder()
      .setTitle('Warnings Cleared')
      .setDescription(`**${target.username}** had **${count}** warning(s) cleared.`)
      .addFields(
        { name: 'Moderator', value: `${interaction.user}`, inline: true },
        { name: 'User', value: `${target}`, inline: true },
      )
      .setColor(Colors.Moderation)
      .setTimestamp();
    await bot.logToChannel(interaction.guildId!, log);
  },
};
