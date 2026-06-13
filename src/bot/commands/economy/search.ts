import { SlashCommandBuilder, EmbedBuilder, Collection } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { search, hasXpBoost, applyXpBoost } from '../../economy/store.js';

const SEARCH_COOLDOWN = 15 * 60 * 1000;
const searchCooldowns = new Collection<string, number>();

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
    const now = Date.now();
    const lastSearch = searchCooldowns.get(userId);
    if (lastSearch && now - lastSearch < SEARCH_COOLDOWN) {
      const remaining = Math.ceil((SEARCH_COOLDOWN - (now - lastSearch)) / 60000);
      await interaction.reply({ content: `⏳ Still searching! Come back in **${remaining}m**.`, flags: 64 });
      return;
    }

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

    searchCooldowns.set(userId, now);

    const newBalance = rec ? rec.balance + boosted : boosted;
    const desc = boosted !== value
      ? `⚡ XP Boost! You searched **${name}** and found **${boosted.toLocaleString()}** Pulses (${value.toLocaleString()} ×2)!`
      : `You searched **${name}** and found **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle('🔍 Search')
      .setDescription(desc)
      .addFields({ name: 'Balance', value: `💰 **${newBalance.toLocaleString()}** Pulses`, inline: false })
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
