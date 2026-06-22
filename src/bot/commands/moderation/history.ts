import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userWarnings } from '../../../db/schema.js';
import { eq, desc, and } from 'drizzle-orm';

export const historyCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('history')
    .setDescription('View moderation history for a user')
    .addUserOption((o) => o.setName('user').setDescription('User to check').setRequired(true))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }

    const target = interaction.options.getUser('user', true);
    const guildId = interaction.guildId!;

    const warns: any[] = await db
      .select()
      .from(userWarnings)
      .where(and(eq(userWarnings.guildId, guildId), eq(userWarnings.userId, target.id)))
      .orderBy(desc(userWarnings.createdAt));

    const warnCount = warns.length;
    const lastWarn = warns[0]
      ? new Date(warns[0].createdAt).toLocaleDateString()
      : 'None';

    const emb = new EmbedBuilder()
      .setTitle(`Moderation History — ${target.username}`)
      .setColor(Colors.Moderation)
      .addFields(
        { name: 'Total Warnings', value: `${warnCount}`, inline: true },
        { name: 'Last Warning', value: lastWarn, inline: true },
      );

    if (warns.length > 0) {
      const recent = warns.slice(0, 5).map(
        (w: any) => `**#${w.id}** — <@${w.moderatorId}> • ${new Date(w.createdAt).toLocaleDateString()}\n└ ${w.reason}`,
      );
      emb.addFields({ name: 'Recent Warnings', value: recent.join('\n'), inline: false });
    }

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
