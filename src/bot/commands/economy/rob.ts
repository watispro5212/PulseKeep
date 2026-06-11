import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS, robSuccess } from '../../economy/store.js';

export const robCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('rob')
    .setDescription('Attempt to rob another user')
    .addUserOption((o) => o.setName('user').setDescription('Target').setRequired(true))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;
    const target = interaction.options.getUser('user', true);
    if (target.id === userId) {
      await interaction.reply({ content: '❌ You cannot rob yourself.', flags: 64 });
      return;
    }
    if (target.bot) {
      await interaction.reply({ content: '❌ You cannot rob a bot.', flags: 64 });
      return;
    }

    const now = new Date();
    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);
    const rec = rows[0];

    if (rec?.lastRob) {
      const elapsed = now.getTime() - new Date(rec.lastRob).getTime();
      if (elapsed < COOLDOWNS.rob) {
        const remaining = Math.ceil((COOLDOWNS.rob - elapsed) / 60000);
        await interaction.reply({ content: `⏳ Wait **${remaining}m** before robbing again.`, flags: 64 });
        return;
      }
    }

    const targetRows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, target.id))
      .limit(1);
    const targetRec = targetRows[0];
    if (!targetRec || targetRec.balance < 50) {
      await interaction.reply({ content: '❌ That user has barely any Pulses.', flags: 64 });
      return;
    }

    const stealAmount = Math.min(targetRec.balance, 100 + Math.floor(Math.random() * 401));

    if (robSuccess()) {
      await db
        .update(userEconomy)
        .set({
          balance: rec ? rec.balance + stealAmount : stealAmount,
          lastRob: now,
          totalEarned: (rec?.totalEarned ?? 0) + stealAmount,
          transactions: (rec?.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));

      await db
        .update(userEconomy)
        .set({
          balance: targetRec.balance - stealAmount,
          transactions: (targetRec.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, target.id));

      const emb = new EmbedBuilder()
        .setTitle(`🦹 Robbery Successful!`)
        .setDescription(`You stole **${stealAmount.toLocaleString()}** Pulses from ${target.username}!`)
        .setColor(Colors.Economy);

      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    } else {
      const fine = Math.min(rec?.balance ?? 0, 200);
      if (fine > 0 && rec) {
        await db
          .update(userEconomy)
          .set({ balance: rec.balance - fine, lastRob: now })
          .where(eq(userEconomy.userId, userId));
      }

      const emb = new EmbedBuilder()
        .setTitle(`🚔 Robbery Failed!`)
        .setDescription(`You got caught and fined **${fine.toLocaleString()}** Pulses.`)
        .setColor(Colors.Error);

      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
