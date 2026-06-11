import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
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

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    let rec = rows[0];
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

    const ephemeral = !interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + total,
          lastDailyClaim: now,
          streak: newStreak,
          totalEarned: (rec.totalEarned ?? 0) + total,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({
          userId,
          balance: total,
          lastDailyClaim: now,
          streak: newStreak,
          totalEarned: total,
          transactions: 1,
        });
    }

    const progress = getStreakMilestoneProgress(newStreak);
    const boosted = hasXpBoost(rec);
    const emb = new EmbedBuilder()
      .setTitle(boosted ? 'Daily Reward ⚡ (x2)' : 'Daily Reward')
      .setDescription(`🎉 You claimed **${total.toLocaleString()}** Pulses${boosted ? ' (with XP Boost!)' : '!'}`)
      .addFields(
        { name: 'Base', value: `${base.toLocaleString()}`, inline: true },
        { name: 'Streak Bonus', value: `+${bonus.toLocaleString()}`, inline: true },
        { name: 'Streak', value: `🔥 **${newStreak}** days (${progress})`, inline: true },
      )
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral });
  },
};
