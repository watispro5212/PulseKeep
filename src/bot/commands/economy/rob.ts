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
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const publicReply = !!interaction.options.getBoolean('public');
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
      await db.transaction(async (tx: any) => {
        const attRows = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
        const tgtRows = await tx.select().from(userEconomy).where(eq(userEconomy.userId, target.id)).limit(1);
        const attRec = attRows[0];
        const tgtRec = tgtRows[0];
        if (!tgtRec || tgtRec.balance < 50) {
          throw new Error('Target has no pulses');
        }
        const actualSteal = Math.min(tgtRec.balance, stealAmount);
        await tx
          .update(userEconomy)
          .set({ balance: (attRec?.balance ?? 0) + actualSteal, lastRob: now, totalEarned: (attRec?.totalEarned ?? 0) + actualSteal, transactions: (attRec?.transactions ?? 0) + 1 })
          .where(eq(userEconomy.userId, userId));
        await tx
          .update(userEconomy)
          .set({ balance: tgtRec.balance - actualSteal, transactions: (tgtRec.transactions ?? 0) + 1 })
          .where(eq(userEconomy.userId, target.id));
      });

      const emb = new EmbedBuilder()
        .setTitle(`🦹 Robbery Successful!`)
        .setDescription(`You stole **${stealAmount.toLocaleString()}** Pulses from ${target.username}!`)
        .setColor(Colors.Economy);

      if (publicReply) {
        await interaction.reply({ embeds: [footer(timestamp(emb))] });
      } else {
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      }
    } else {
      let fineAmount = 0;
      await db.transaction(async (tx: any) => {
        const attRows = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
        const attRec = attRows[0];
        if (attRec) {
          fineAmount = Math.min(attRec.balance, 200);
          if (fineAmount > 0) {
            await tx
              .update(userEconomy)
              .set({ balance: attRec.balance - fineAmount, lastRob: now })
              .where(eq(userEconomy.userId, userId));
          }
        }
      });

      const emb = new EmbedBuilder()
        .setTitle(`🚔 Robbery Failed!`)
        .setDescription(`You got caught and fined **${fineAmount.toLocaleString()}** Pulses.`)
        .setColor(Colors.Error);

      if (publicReply) {
        await interaction.reply({ embeds: [footer(timestamp(emb))] });
      } else {
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      }
    }
  },
};
