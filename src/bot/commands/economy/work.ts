import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS, WORK_JOBS, WORK_FLAVOR, hasXpBoost, applyXpBoost } from '../../economy/store.js';

export const workCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('work')
    .setDescription('Work to earn Pulses')
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

    if (rec?.lastWork) {
      const elapsed = now.getTime() - new Date(rec.lastWork).getTime();
      if (elapsed < COOLDOWNS.work) {
        const remaining = Math.ceil((COOLDOWNS.work - elapsed) / 60000);
        await interaction.reply({ content: `⏳ You're tired! Come back in **${remaining}m**.`, flags: 64 });
        return;
      }
    }

    const job = WORK_JOBS[Math.floor(Math.random() * WORK_JOBS.length)]!;
    const earned = job.pay[0]! + Math.floor(Math.random() * (job.pay[1]! - job.pay[0]!));
    const boosted = hasXpBoost(rec) ? applyXpBoost(earned, rec) : earned;
    const flavor = WORK_FLAVOR[Math.floor(Math.random() * WORK_FLAVOR.length)];
    const publicReply = !!interaction.options.getBoolean('public');

    try {
      if (rec) {
        await db
          .update(userEconomy)
          .set({
            balance: rec.balance + boosted,
            lastWork: now,
            totalEarned: (rec.totalEarned ?? 0) + boosted,
            transactions: (rec.transactions ?? 0) + 1,
          })
          .where(eq(userEconomy.userId, userId));
      } else {
        await db
          .insert(userEconomy)
          .values({ userId, balance: boosted, lastWork: now, totalEarned: boosted, transactions: 1 });
      }
    } catch (err) {
      await interaction.reply({ content: '❌ Failed to save earnings. Please try again.', flags: 64 });
      return;
    }

    const newBalance = rec ? rec.balance + boosted : boosted;
    const titleDesc = boosted !== earned
      ? `⚡ XP Boost active! You ${flavor} **${boosted.toLocaleString()}** Pulses (${earned.toLocaleString()} ×2)!`
      : `You ${flavor} **${boosted.toLocaleString()}** Pulses!`;

    const emb = new EmbedBuilder()
      .setTitle(`💼 ${job.title}`)
      .setDescription(titleDesc)
      .addFields({ name: 'Balance', value: `💰 **${newBalance.toLocaleString()}** Pulses`, inline: false })
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
