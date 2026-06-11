import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq, and } from 'drizzle-orm';
import { COOLDOWNS, mine } from '../../economy/store.js';

export const mineCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('mine')
    .setDescription('Go mining to earn Pulses')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
      return;
    }

    const userId = interaction.user.id;

    const inv = await db
      .select()
      .from(userInventory)
      .where(and(eq(userInventory.userId, userId), eq(userInventory.itemId, 'mining_pick')))
      .limit(1);

    if (inv.length === 0) {
      await interaction.reply({ content: '❌ You need a **Mining Pick**! Buy one from /shop.', ephemeral: true });
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
        await interaction.reply({ content: '⏳ Wait a moment before mining again.', ephemeral: true });
        return;
      }
    }

    const { name, value } = mine();
    const ephemeral = !interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + value,
          lastMine: now,
          totalEarned: (rec.totalEarned ?? 0) + value,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({ userId, balance: value, lastMine: now, totalEarned: value, transactions: 1 });
    }

    const emb = new EmbedBuilder()
      .setTitle('⛏️ Mining')
      .setDescription(`You found **${name}** worth **${value.toLocaleString()}** Pulses!`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral });
  },
};
