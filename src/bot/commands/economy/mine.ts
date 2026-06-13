import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq, and } from 'drizzle-orm';
import { COOLDOWNS, mine, hasXpBoost, applyXpBoost, formatCooldown } from '../../economy/store.js';

export const mineCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('mine')
    .setDescription('Go mining to earn Pulses')
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
      .where(and(eq(userInventory.userId, userId), eq(userInventory.itemId, 'mining_pick')))
      .limit(1);

    if (inv.length === 0 || inv[0].quantity <= 0) {
      await interaction.reply({ content: '❌ You need a **Mining Pick**! Buy one from /shop.', flags: 64 });
      return;
    }

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);
    const rec = rows[0];
    const now = new Date();

    if (rec?.lastMine) {
      const elapsed = now.getTime() - new Date(rec.lastMine).getTime();
      if (elapsed < COOLDOWNS.mine) {
        const remaining = formatCooldown(elapsed, COOLDOWNS.mine);
        await interaction.reply({ content: `⏳ Mine again in **${remaining}**.`, flags: 64 });
        return;
      }
    }

    const { name, value } = mine();
    const boosted = hasXpBoost(rec) ? applyXpBoost(value, rec) : value;
    const publicReply = !!interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + boosted,
          lastMine: now,
          totalEarned: (rec.totalEarned ?? 0) + boosted,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({ userId, balance: boosted, lastMine: now, totalEarned: boosted, transactions: 1 });
    }

    const newBalance = rec ? rec.balance + boosted : boosted;
    const desc = boosted !== value
      ? `⚡ XP Boost! You found **${name}** worth **${boosted.toLocaleString()}** Pulses (${value.toLocaleString()} ×2)!`
      : `You found **${name}** worth **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle('⛏️ Mining')
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
