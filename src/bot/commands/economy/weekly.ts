import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS, hasXpBoost, applyXpBoost } from '../../economy/store.js';

export const weeklyCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('weekly')
    .setDescription('Claim your weekly Pulses')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;
    const now = new Date();

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    const rec = rows[0];

    if (rec?.lastWeeklyClaim) {
      const elapsed = now.getTime() - new Date(rec.lastWeeklyClaim).getTime();
      if (elapsed < COOLDOWNS.weekly) {
        const remaining = Math.ceil((COOLDOWNS.weekly - elapsed) / 86400000);
        await interaction.reply({ content: `⏳ Come back in **${remaining}d** for your weekly reward.`, flags: 64 });
        return;
      }
    }

    const baseReward = 1000 + Math.floor(Math.random() * 1001);
    const reward = hasXpBoost(rec) ? applyXpBoost(baseReward, rec) : baseReward;
    const boosted = reward !== baseReward;
    const publicReply = !!interaction.options.getBoolean('public');

    await db.transaction(async (tx: any) => {
      const reRec = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
      if (reRec[0]) {
        await tx
          .update(userEconomy)
          .set({
            balance: reRec[0].balance + reward,
            lastWeeklyClaim: now,
            totalEarned: (reRec[0].totalEarned ?? 0) + reward,
            transactions: (reRec[0].transactions ?? 0) + 1,
          })
          .where(eq(userEconomy.userId, userId));
      } else {
        await tx
          .insert(userEconomy)
          .values({ userId, balance: reward, lastWeeklyClaim: now, totalEarned: reward, transactions: 1 });
      }
    });

    const emb = new EmbedBuilder()
      .setTitle(boosted ? 'Weekly Reward ⚡ (x2)' : 'Weekly Reward')
      .setDescription(`📅 You claimed **${reward.toLocaleString()}** Pulses${boosted ? ' (with XP Boost!)' : '!'}`)
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
