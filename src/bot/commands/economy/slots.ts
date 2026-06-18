import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { rollSlots } from '../../economy/store.js';

export const slotsCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('slots')
    .setDescription('Spin the slot machine')
    .addIntegerOption((o) =>
      o.setName('amount').setDescription('Amount to bet').setRequired(true).setMinValue(10),
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
    const bet = interaction.options.getInteger('amount', true);

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    const rec = rows[0];
    if (!rec || rec.balance < bet) {
      await interaction.reply({ content: '❌ You don\'t have enough Pulses.', flags: 64 });
      return;
    }

    const { payout, emojis } = rollSlots();
    const emojiStr = emojis.join(' ');
    let change: number;
    let title: string;
    let color: number;

    if (payout > 0) {
      const winnings = Math.floor(bet * payout);
      change = winnings - bet;
      title = `🎰 **${emojiStr}**\nYou won **${winnings.toLocaleString()}** Pulses! (${payout}x)`;
      color = Colors.Economy;
    } else {
      change = -bet;
      title = `🎰 **${emojiStr}**\nNo luck this time. Lost **${bet.toLocaleString()}** Pulses.`;
      color = Colors.Error;
    }

    const newBalance = rec.balance + change;

    await db
      .update(userEconomy)
      .set({
        balance: newBalance,
        totalGambled: (rec.totalGambled ?? 0) + bet,
        transactions: (rec.transactions ?? 0) + 1,
      })
      .where(eq(userEconomy.userId, userId));

    const emb = new EmbedBuilder()
      .setTitle('Slot Machine')
      .setDescription(title)
      .addFields({ name: 'New Balance', value: `💰 ${newBalance.toLocaleString()}`, inline: true })
      .setColor(color);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
