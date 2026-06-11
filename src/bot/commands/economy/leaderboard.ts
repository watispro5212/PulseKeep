import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { desc } from 'drizzle-orm';

export const leaderboardCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('leaderboard')
    .setDescription('Show the richest users')
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
      return;
    }

    const rows: any[] = await db
      .select()
      .from(userEconomy)
      .orderBy(desc(userEconomy.balance))
      .limit(10);

    if (rows.length === 0) {
      await interaction.reply({
        embeds: [footer(timestamp(new EmbedBuilder()
          .setTitle('🏆 Leaderboard')
          .setDescription('No data yet! Start earning with /daily and /work.')
          .setColor(Colors.Economy)))],
        ephemeral: true,
      });
      return;
    }

    const entries = rows.map((r: any, i: number) => `**#${i + 1}** <@${r.userId}> — 💰 **${r.balance.toLocaleString()}**`);

    const emb = new EmbedBuilder()
      .setTitle('🏆 PulseKeep Leaderboard')
      .setDescription(entries.join('\n'))
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
