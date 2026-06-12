import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { search, hasXpBoost, applyXpBoost } from '../../economy/store.js';

export const searchCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('search')
    .setDescription('Search for hidden Pulses around town')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;
    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    const rec = rows[0];
    const { name, value } = search();
    const boosted = hasXpBoost(rec) ? applyXpBoost(value, rec) : value;
    const publicReply = !!interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + boosted,
          totalEarned: (rec.totalEarned ?? 0) + boosted,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({ userId, balance: boosted, totalEarned: boosted, transactions: 1 });
    }

    const desc = boosted !== value
      ? `⚡ XP Boost! You searched **${name}** and found **${boosted.toLocaleString()}** Pulses (${value.toLocaleString()} ×2)!`
      : `You searched **${name}** and found **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle('🔍 Search')
      .setDescription(desc)
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
