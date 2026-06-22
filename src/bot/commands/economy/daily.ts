import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp, formatNumber } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS, getStreakBonus, getStreakMilestoneProgress, hasXpBoost, applyXpBoost } from '../../economy/store.js';

export const dailyCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('daily')
    .setDescription('Claim your daily Pulses')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;
    const now = new Date();

    let rec: any;
    try {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, userId))
        .limit(1);
      rec = rows[0];
    } catch (err) {
      await interaction.reply({ content: '❌ Database error. Please try again.', flags: 64 });
      return;
    }
    if (rec?.lastDailyClaim) {
      const elapsed = now.getTime() - new Date(rec.lastDailyClaim).getTime();
      if (elapsed < COOLDOWNS.daily) {
        const remaining = Math.ceil((COOLDOWNS.daily - elapsed) / 3600000);
        await interaction.reply({ content: `⏳ Come back in **${remaining}h** for your daily reward.`, flags: 64 });
        return;
      }
    }

    const lastClaim = rec?.lastDailyClaim ? new Date(rec.lastDailyClaim) : null;
    const isConsecutive = lastClaim ? (now.getTime() - lastClaim.getTime()) < 48 * 60 * 60 * 1000 : false;
    const streak = rec?.streak ?? 0;
    const newStreak = lastClaim ? (isConsecutive ? streak + 1 : 1) : 1;
    const base = 250 + Math.floor(Math.random() * 251);
    const bonus = getStreakBonus(newStreak);
    const total = hasXpBoost(rec) ? applyXpBoost(base + bonus, rec) : base + bonus;

    const publicReply = !!interaction.options.getBoolean('public');

    try {
      await db.transaction(async (tx: any) => {
        const reRec = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
        if (reRec[0]) {
          await tx
            .update(userEconomy)
            .set({
              balance: reRec[0].balance + total,
              lastDailyClaim: now,
              streak: newStreak,
              totalEarned: (reRec[0].totalEarned ?? 0) + total,
              transactions: (reRec[0].transactions ?? 0) + 1,
            })
            .where(eq(userEconomy.userId, userId));
        } else {
          await tx
            .insert(userEconomy)
            .values({ userId, balance: total, lastDailyClaim: now, streak: newStreak, totalEarned: total, transactions: 1 });
        }
      });
    } catch (err) {
      await interaction.reply({ content: '❌ Failed to save daily reward. Please try again.', flags: 64 });
      return;
    }

    const newBalance = (rec?.balance ?? 0) + total;
    const progress = getStreakMilestoneProgress(newStreak);
    const boosted = hasXpBoost(rec);
    const emb = new EmbedBuilder()
      .setTitle(boosted ? 'Daily Reward ⚡ (x2)' : 'Daily Reward')
      .setDescription(`🎉 You claimed **${formatNumber(total)}** Pulses${boosted ? ' (with XP Boost!)' : '!'}`)
      .addFields(
        { name: 'Base', value: `${formatNumber(base)}`, inline: true },
        { name: 'Streak Bonus', value: `+${formatNumber(bonus)}`, inline: true },
        { name: 'Streak', value: `🔥 **${newStreak}** days (${progress})`, inline: true },
        { name: 'Balance', value: `💰 **${formatNumber(newBalance)}** Pulses`, inline: false },
      )
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
