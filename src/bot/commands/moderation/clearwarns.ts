import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
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
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const target = interaction.options.getUser('user', true);

    const result = await db
      .delete(userWarnings)
      .where(and(eq(userWarnings.guildId, interaction.guildId!), eq(userWarnings.userId, target.id)));

    const count = (result as any)?.rowCount ?? 0;

    const emb = new EmbedBuilder()
      .setTitle('Warnings Cleared')
      .setDescription(`Cleared **${count}** warning(s) for ${target}.`)
      .setColor(Colors.Moderation);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

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
