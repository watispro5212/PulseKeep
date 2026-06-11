import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq, and } from 'drizzle-orm';
import { COOLDOWNS, fish, hasXpBoost, applyXpBoost } from '../../economy/store.js';

export const fishCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('fish')
    .setDescription('Go fishing to earn Pulses')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;

    const inv = await db
      .select()
      .from(userInventory)
      .where(and(eq(userInventory.userId, userId), eq(userInventory.itemId, 'fishing_rod')))
      .limit(1);

    if (inv.length === 0) {
      await interaction.reply({ content: '❌ You need a **Fishing Rod**! Buy one from /shop.', flags: 64 });
      return;
    }

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);
    const rec = rows[0];
    const now = new Date();

    if (rec?.lastFish) {
      const elapsed = now.getTime() - new Date(rec.lastFish).getTime();
      if (elapsed < COOLDOWNS.fish) {
        await interaction.reply({ content: '⏳ Wait a moment before fishing again.', flags: 64 });
        return;
      }
    }

    const { name, value } = fish();
    const boosted = hasXpBoost(rec) ? applyXpBoost(value, rec) : value;
    const ephemeral = !interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + boosted,
          lastFish: now,
          totalEarned: (rec.totalEarned ?? 0) + boosted,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({ userId, balance: boosted, lastFish: now, totalEarned: boosted, transactions: 1 });
    }

    const desc = boosted !== value
      ? `⚡ XP Boost! You caught a **${name}** worth **${boosted.toLocaleString()}** Pulses (${value.toLocaleString()} ×2)!`
      : `You caught a **${name}** worth **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle('🎣 Fishing')
      .setDescription(desc)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral });
  },
};
