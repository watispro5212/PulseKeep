import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

export const payCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('pay')
    .setDescription('Send Pulses to another user')
    .addUserOption((o) => o.setName('user').setDescription('Recipient').setRequired(true))
    .addIntegerOption((o) =>
      o.setName('amount').setDescription('Amount to send').setRequired(true).setMinValue(1),
    )
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const senderId = interaction.user.id;
    const recipient = interaction.options.getUser('user', true);
    const amount = interaction.options.getInteger('amount', true);

    if (recipient.id === senderId) {
      await interaction.reply({ content: '❌ You cannot pay yourself.', flags: 64 });
      return;
    }
    if (recipient.bot) {
      await interaction.reply({ content: '❌ You cannot pay a bot.', flags: 64 });
      return;
    }

    const senderRows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, senderId))
      .limit(1);
    const sender = senderRows[0];
    if (!sender || sender.balance < amount) {
      await interaction.reply({ content: '❌ You don\'t have enough Pulses.', flags: 64 });
      return;
    }

    const recvRows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, recipient.id))
      .limit(1);
    const recv = recvRows[0];

    await db.transaction(async (tx: any) => {
      await tx
        .update(userEconomy)
        .set({ balance: sender.balance - amount, transactions: (sender.transactions ?? 0) + 1 })
        .where(eq(userEconomy.userId, senderId));

      if (recv) {
        await tx
          .update(userEconomy)
          .set({ balance: recv.balance + amount, totalEarned: (recv.totalEarned ?? 0) + amount, transactions: (recv.transactions ?? 0) + 1 })
          .where(eq(userEconomy.userId, recipient.id));
      } else {
        await tx
          .insert(userEconomy)
          .values({ userId: recipient.id, balance: amount, totalEarned: amount, transactions: 1 });
      }
    });

    const emb = new EmbedBuilder()
      .setTitle('Payment Sent')
      .setDescription(`Sent **${amount.toLocaleString()}** Pulses to ${recipient}.`)
      .addFields(
        { name: 'New Balance', value: `💰 ${(sender.balance - amount).toLocaleString()}`, inline: true },
      )
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
