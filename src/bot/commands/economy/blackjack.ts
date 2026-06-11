import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { blackjackDealerShouldStand } from '../../economy/store.js';

function drawCard(): number {
  return Math.min(10, 1 + Math.floor(Math.random() * 13));
}

function handValue(hand: number[]): number {
  let total = hand.reduce((s, c) => s + c, 0);
  let aces = hand.filter((c) => c === 1).length;
  while (total <= 11 && aces > 0) {
    total += 10;
    aces--;
  }
  return total;
}

export const blackjackCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('blackjack')
    .setDescription('Play blackjack against the dealer')
    .addIntegerOption((o) =>
      o.setName('bet').setDescription('Your bet').setRequired(true).setMinValue(10),
    )
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
      return;
    }

    const userId = interaction.user.id;
    const bet = interaction.options.getInteger('bet', true);

    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);

    const rec = rows[0];
    if (!rec || rec.balance < bet) {
      await interaction.reply({ content: '❌ You don\'t have enough Pulses.', ephemeral: true });
      return;
    }

    const playerHand = [drawCard(), drawCard()];
    const dealerHand = [drawCard(), drawCard()];

    const playerTotal = handValue(playerHand);
    const dealerTotal = handValue(dealerHand);

    let result: string;
    let change: number;

    if (playerTotal === 21 && dealerTotal !== 21) {
      result = '🎉 **Blackjack!** You got 21!';
      change = Math.floor(bet * 2.5);
    } else if (playerTotal > 21) {
      result = '💥 **Bust!** You went over 21.';
      change = -bet;
    } else if (dealerTotal > 21) {
      result = '🎉 **Dealer busts!** You win!';
      change = bet;
    } else if (playerTotal > dealerTotal) {
      result = '🎉 **You win!** Better hand than the dealer.';
      change = bet;
    } else if (playerTotal < dealerTotal) {
      result = '💸 **You lose!** Dealer has a better hand.';
      change = -bet;
    } else {
      result = '🤝 **Push!** Same hand.';
      change = 0;
    }

    let dealerPlays = '';
    if (change >= 0 || playerTotal > 21) {
      dealerPlays = ' (dealer didn\'t need to draw)';
    } else {
      const finalDealer = handValue(dealerHand);
      dealerPlays = finalDealer > 21 ? ' — Dealer busted!' : '';
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
      .setTitle('♠️ Blackjack')
      .setDescription(result)
      .addFields(
        {
          name: `Your Hand (${playerTotal})`,
          value: playerHand.join(' · '),
          inline: true,
        },
        {
          name: `Dealer's Hand (${dealerTotal})`,
          value: `${dealerHand.join(' · ')}${dealerPlays}`,
          inline: true,
        },
        { name: 'Payout', value: `${change >= 0 ? '+' : ''}${change.toLocaleString()}`, inline: true },
      )
      .setColor(change >= 0 ? Colors.Economy : Colors.Error);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
