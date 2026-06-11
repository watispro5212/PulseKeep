import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS, WORK_JOBS, WORK_FLAVOR } from '../../economy/store.js';

export const workCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('work')
    .setDescription('Work to earn Pulses')
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
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

    if (rec?.lastWork) {
      const elapsed = now.getTime() - new Date(rec.lastWork).getTime();
      if (elapsed < COOLDOWNS.work) {
        const remaining = Math.ceil((COOLDOWNS.work - elapsed) / 60000);
        await interaction.reply({ content: `⏳ You're tired! Come back in **${remaining}m**.`, ephemeral: true });
        return;
      }
    }

    const job = WORK_JOBS[Math.floor(Math.random() * WORK_JOBS.length)]!;
    const earned = job.pay[0]! + Math.floor(Math.random() * (job.pay[1]! - job.pay[0]!));
    const flavor = WORK_FLAVOR[Math.floor(Math.random() * WORK_FLAVOR.length)];
    const ephemeral = !interaction.options.getBoolean('public');

    if (rec) {
      await db
        .update(userEconomy)
        .set({
          balance: rec.balance + earned,
          lastWork: now,
          totalEarned: (rec.totalEarned ?? 0) + earned,
          transactions: (rec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    } else {
      await db
        .insert(userEconomy)
        .values({ userId, balance: earned, lastWork: now, totalEarned: earned, transactions: 1 });
    }

    const emb = new EmbedBuilder()
      .setTitle(`💼 ${job.title}`)
      .setDescription(`You ${flavor} **${earned.toLocaleString()}** Pulses!`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral });
  },
};
