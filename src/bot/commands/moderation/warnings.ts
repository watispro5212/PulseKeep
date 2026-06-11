import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userWarnings } from '../../../db/schema.js';
import { eq, desc } from 'drizzle-orm';

export const warningsCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('warnings')
    .setDescription('Check a user\'s warnings')
    .addUserOption((o) => o.setName('user').setDescription('User to check').setRequired(true))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const target = interaction.options.getUser('user', true);
    const rows: any[] = await db
      .select()
      .from(userWarnings)
      .where(eq(userWarnings.userId, target.id))
      .where(eq(userWarnings.guildId, interaction.guildId!))
      .orderBy(desc(userWarnings.createdAt));

    if (rows.length === 0) {
      await interaction.reply({
        embeds: [footer(timestamp(new EmbedBuilder()
          .setTitle('Warnings')
          .setDescription(`**${target.tag}** has no warnings.`)
          .setColor(Colors.Moderation)))],
        flags: 64,
      });
      return;
    }

    const list = rows.slice(0, 10).map(
      (w: any) => `**#${w.id}** — <@${w.moderatorId}> • ${new Date(w.createdAt).toLocaleDateString()}\n└ ${w.reason}`,
    );

    if (rows.length > 10) {
      list.push(`... and ${rows.length - 10} more`);
    }

    const emb = new EmbedBuilder()
      .setTitle(`Warnings — ${target.tag}`)
      .setDescription(list.join('\n'))
      .addFields({ name: 'Total', value: `${rows.length} warning(s)`, inline: true })
      .setColor(Colors.Moderation);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
