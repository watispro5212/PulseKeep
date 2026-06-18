import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { gamble } from '../../economy/store.js';

export const gambleCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('gamble')
    .setDescription('Gamble your Pulses for a chance to win big')
    .addIntegerOption((o) =>
      o.setName('amount').setDescription('Amount to gamble').setRequired(true).setMinValue(10),
    )
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const publicReply = !!interaction.options.getBoolean('public');
    const userId = interaction.user.id;
    const amount = interaction.options.getInteger('amount', true);

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    const rec = rows[0];
    if (!rec || rec.balance < amount) {
      await interaction.reply({ content: '❌ You don\'t have enough Pulses.', flags: 64 });
      return;
    }

    const { result, payout, multiplier } = gamble(amount);
    let change = 0;
    let title: string;
    let color: number;
    let extraFields: { name: string; value: string; inline: boolean }[] = [];

    if (result === 'win') {
      change = payout - amount;
      title = `🎉 You won **${change.toLocaleString()}** Pulses!`;
      color = Colors.Economy;
      extraFields.push({ name: 'Multiplier', value: `×${multiplier}`, inline: true });
    } else if (result === 'push') {
      change = 0;
      title = `🤝 Push! You got your **${amount.toLocaleString()}** back.`;
      color = Colors.Warning;
    } else {
      if ((rec.luckyCloverActive ?? 0) > 0) {
        change = 0;
        title = `🍀 Lucky Clover saved you! Your **${amount.toLocaleString()}** was refunded.`;
        color = Colors.Economy;
        extraFields.push({ name: 'Clovers Left', value: `${(rec.luckyCloverActive ?? 0) - 1}`, inline: true });
      } else {
        change = -amount;
        title = `💸 You lost **${amount.toLocaleString()}** Pulses.`;
        color = Colors.Error;
      }
    }

    const newBalance = rec.balance + change;

    await db.transaction(async (tx: any) => {
      const reRec = await tx
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, userId))
        .limit(1);
      if (result === 'lose' && (rec.luckyCloverActive ?? 0) > 0) {
        await tx
          .update(userEconomy)
          .set({ luckyCloverActive: (reRec[0]?.luckyCloverActive ?? 0) - 1 })
          .where(eq(userEconomy.userId, userId));
      }
      await tx
        .update(userEconomy)
        .set({
          balance: newBalance,
          totalGambled: (reRec[0]?.totalGambled ?? 0) + amount,
          transactions: (reRec[0]?.transactions ?? 0) + 1,
        })
        .where(eq(userEconomy.userId, userId));
    });

    const emb = new EmbedBuilder()
      .setTitle(title)
      .addFields(
        { name: 'New Balance', value: `💰 ${newBalance.toLocaleString()}`, inline: true },
        { name: 'Bet', value: `${amount.toLocaleString()}`, inline: true },
        ...extraFields,
      )
      .setColor(color);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
