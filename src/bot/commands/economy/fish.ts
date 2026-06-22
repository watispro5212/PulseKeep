import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq, and } from 'drizzle-orm';
import { COOLDOWNS, fish, hasXpBoost, applyXpBoost, formatCooldown } from '../../economy/store.js';

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

    if (inv.length === 0 || inv[0].quantity <= 0) {
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
        const remaining = formatCooldown(elapsed, COOLDOWNS.fish);
        await interaction.reply({ content: `⏳ Fish again in **${remaining}**.`, flags: 64 });
        return;
      }
    }

    const { name, value } = fish();
    const boosted = hasXpBoost(rec) ? applyXpBoost(value, rec) : value;
    const publicReply = !!interaction.options.getBoolean('public');

    await db.transaction(async (tx: any) => {
      const reRec = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
      if (reRec[0]) {
        await tx
          .update(userEconomy)
          .set({
            balance: reRec[0].balance + boosted,
            lastFish: now,
            totalEarned: (reRec[0].totalEarned ?? 0) + boosted,
            transactions: (reRec[0].transactions ?? 0) + 1,
          })
          .where(eq(userEconomy.userId, userId));
      } else {
        await tx
          .insert(userEconomy)
          .values({ userId, balance: boosted, lastFish: now, totalEarned: boosted, transactions: 1 });
      }
    });

    const newBalance = (rec?.balance ?? 0) + boosted;
    const desc = boosted !== value
      ? `⚡ XP Boost! You caught a **${name}** worth **${boosted.toLocaleString()}** Pulses (${value.toLocaleString()} ×2)!`
      : `You caught a **${name}** worth **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle('🎣 Fishing')
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
