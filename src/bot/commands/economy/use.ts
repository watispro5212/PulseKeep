import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

const XP_BOOST_DURATION_MS = 30 * 60 * 1000;

const USABLE_ITEMS: Record<string, { label: string; min: number; max: number }> = {
  treasure_map: { label: 'Treasure Map 🗺️', min: 2500, max: 7500 },
  lucky_clover: { label: 'Lucky Clover 🍀', min: 0, max: 0 },
  xp_boost: { label: 'XP Boost ⚡', min: 0, max: 0 },
};

export const useCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('use')
    .setDescription('Use an item from your inventory')
    .addStringOption((o) =>
      o.setName('item').setDescription('Item to use').setRequired(true)
        .addChoices(
          { name: 'Treasure Map 🗺️', value: 'treasure_map' },
          { name: 'Lucky Clover 🍀', value: 'lucky_clover' },
          { name: 'XP Boost ⚡', value: 'xp_boost' },
        ),
    )
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const itemId = interaction.options.getString('item', true);
    const userId = interaction.user.id;

    const inv = await db
      .select()
      .from(userInventory)
      .where(eq(userInventory.userId, userId))
      .where(eq(userInventory.itemId, itemId))
      .limit(1);

    if (!inv.length || !inv[0] || (inv[0].quantity ?? 0) < 1) {
      await interaction.reply({ content: '❌ You don\'t own this item.', flags: 64 });
      return;
    }

    const def = USABLE_ITEMS[itemId];
    if (!def) {
      await interaction.reply({ content: '❌ This item cannot be used.', flags: 64 });
      return;
    }

    if (itemId === 'treasure_map') {
      const value = def.min + Math.floor(Math.random() * (def.max - def.min));

      await db.transaction(async (tx: any) => {
        const rows = await tx
          .select()
          .from(userEconomy)
          .where(eq(userEconomy.userId, userId))
          .limit(1);

        const rec = rows[0];
        if (rec) {
          await tx
            .update(userEconomy)
            .set({
              balance: (rec.balance ?? 0) + value,
              totalEarned: (rec.totalEarned ?? 0) + value,
              transactions: (rec.transactions ?? 0) + 1,
            })
            .where(eq(userEconomy.userId, userId));
        } else {
          await tx
            .insert(userEconomy)
            .values({ userId, balance: value, totalEarned: value, transactions: 1 });
        }

        if ((inv[0]?.quantity ?? 0) > 1) {
          await tx
            .update(userInventory)
            .set({ quantity: (inv[0]?.quantity ?? 1) - 1 })
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        } else {
          await tx
            .delete(userInventory)
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        }
      });

      const emb = new EmbedBuilder()
        .setTitle('Treasure Found!')
        .setDescription(`You used **${def.label}** and found **${value.toLocaleString()}** Pulses! 💰`)
        .setColor(Colors.Economy);

      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      return;
    }

    if (itemId === 'lucky_clover') {
      await db.transaction(async (tx: any) => {
        const rows = await tx
          .select()
          .from(userEconomy)
          .where(eq(userEconomy.userId, userId))
          .limit(1);
        const rec = rows[0];
        if (rec) {
          await tx
            .update(userEconomy)
            .set({ luckyCloverActive: (rec.luckyCloverActive ?? 0) + 1 })
            .where(eq(userEconomy.userId, userId));
        } else {
          await tx
            .insert(userEconomy)
            .values({ userId, luckyCloverActive: 1 });
        }
        if ((inv[0]?.quantity ?? 0) > 1) {
          await tx
            .update(userInventory)
            .set({ quantity: (inv[0]?.quantity ?? 1) - 1 })
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        } else {
          await tx
            .delete(userInventory)
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        }
      });
      await interaction.reply({ content: '🍀 Lucky Clover activated! Your next gamble loss will be refunded.', flags: 64 });
      return;
    }

    if (itemId === 'xp_boost') {
      const expiry = new Date(Date.now() + XP_BOOST_DURATION_MS);
      await db.transaction(async (tx: any) => {
        const rows = await tx
          .select()
          .from(userEconomy)
          .where(eq(userEconomy.userId, userId))
          .limit(1);
        const rec = rows[0];
        if (rec) {
          await tx
            .update(userEconomy)
            .set({ xpBoostExpiry: expiry })
            .where(eq(userEconomy.userId, userId));
        } else {
          await tx
            .insert(userEconomy)
            .values({ userId, xpBoostExpiry: expiry });
        }
        if ((inv[0]?.quantity ?? 0) > 1) {
          await tx
            .update(userInventory)
            .set({ quantity: (inv[0]?.quantity ?? 1) - 1 })
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        } else {
          await tx
            .delete(userInventory)
            .where(eq(userInventory.id, inv[0]?.id ?? 0));
        }
      });
      await interaction.reply({ content: `⚡ XP Boost activated! Your earnings are doubled for **30 minutes**.`, flags: 64 });
      return;
    }
  },
};
