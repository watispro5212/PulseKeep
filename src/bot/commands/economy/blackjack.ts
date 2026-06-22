import { SlashCommandBuilder, EmbedBuilder, ActionRowBuilder, ButtonBuilder, ButtonStyle, ComponentType, ButtonInteraction } from 'discord.js';
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

function handDisplay(hand: number[]): string {
  return hand.join(' · ');
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

    if (handValue(playerHand) === 21) {
      const change = Math.floor(bet * 1.5);
      const newBalance = rec.balance + change;
      await db.transaction(async (tx: any) => {
        await tx.update(userEconomy)
          .set({ balance: newBalance, totalGambled: (rec.totalGambled ?? 0) + bet, transactions: (rec.transactions ?? 0) + 1 })
          .where(eq(userEconomy.userId, userId));
      });
      const emb = new EmbedBuilder()
        .setTitle('♠️ Blackjack')
        .setDescription('🎉 **Blackjack!** You got 21!')
        .addFields(
          { name: `Your Hand (21)`, value: handDisplay(playerHand), inline: true },
          { name: `Dealer's Hand (${handValue(dealerHand)})`, value: handDisplay(dealerHand), inline: true },
          { name: 'Payout', value: `+${change.toLocaleString()}`, inline: true },
          { name: 'Balance', value: `💰 **${newBalance.toLocaleString()}** Pulses`, inline: false },
        )
        .setColor(Colors.Economy);
      await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
      return;
    }

    const hitButton = new ButtonBuilder().setCustomId('hit').setLabel('Hit').setStyle(ButtonStyle.Primary).setEmoji('👋');
    const standButton = new ButtonBuilder().setCustomId('stand').setLabel('Stand').setStyle(ButtonStyle.Secondary).setEmoji('✋');
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(hitButton, standButton);

    const initialEmbed = () => new EmbedBuilder()
      .setTitle('♠️ Blackjack')
      .setDescription(`Bet: **${bet.toLocaleString()}** Pulses`)
      .addFields(
        { name: `Your Hand (${handValue(playerHand)})`, value: handDisplay(playerHand), inline: true },
        { name: `Dealer's Hand (${handValue([dealerHand[0]!])} + ??)`, value: `${dealerHand[0]} · ??`, inline: true },
      )
      .setColor(Colors.Economy);

    const reply = await interaction.reply({
      embeds: [footer(timestamp(initialEmbed()))],
      components: [row],
      fetchReply: true,
    });

    let bust = false;
    let stood = false;

    const collector = reply.createMessageComponentCollector({
      componentType: ComponentType.Button,
      filter: (i: ButtonInteraction) => i.user.id === userId,
      time: 30000,
    });

    collector.on('collect', async (i: ButtonInteraction) => {
      if (i.customId === 'hit') {
        playerHand.push(drawCard());
        const total = handValue(playerHand);
        if (total > 21) {
          bust = true;
          collector.stop();
          const change = -bet;
          await db.transaction(async (tx: any) => {
            await tx.update(userEconomy)
              .set({ balance: (rec?.balance ?? 0) + change, totalGambled: ((rec?.totalGambled ?? 0) + bet), transactions: ((rec?.transactions ?? 0) + 1) })
              .where(eq(userEconomy.userId, userId));
          });
          const newBalance = rec ? rec.balance + change : bet;
          const emb = new EmbedBuilder()
            .setTitle('♠️ Blackjack — Bust!')
            .setDescription('💥 **Bust!** You went over 21.')
            .addFields(
              { name: `Your Hand (${total})`, value: handDisplay(playerHand), inline: true },
              { name: `Dealer's Hand (${handValue(dealerHand)})`, value: handDisplay(dealerHand), inline: true },
              { name: 'Payout', value: `-${bet.toLocaleString()}`, inline: true },
              { name: 'Balance', value: `💰 **${newBalance.toLocaleString()}** Pulses`, inline: false },
            )
            .setColor(Colors.Error);
          await i.update({ embeds: [footer(timestamp(emb))], components: [] });
          return;
        }
        const embed = initialEmbed()
          .spliceFields(0, 2,
            { name: `Your Hand (${total})`, value: handDisplay(playerHand), inline: true },
            { name: `Dealer's Hand (${handValue([dealerHand[0]!])} + ??)`, value: `${dealerHand[0]} · ??`, inline: true },
          );
        await i.update({ embeds: [footer(timestamp(embed))], components: [row] });
      } else if (i.customId === 'stand') {
        stood = true;
        collector.stop();
      }
    });

    collector.on('end', async () => {
      if (bust) return;
      while (!blackjackDealerShouldStand(handValue(dealerHand))) {
        dealerHand.push(drawCard());
      }
      const playerTotal = handValue(playerHand);
      const dealerTotal = handValue(dealerHand);
      let result: string;
      let change: number;
      if (dealerTotal > 21) {
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
      await db.transaction(async (tx: any) => {
        const reRec = await tx.select().from(userEconomy).where(eq(userEconomy.userId, userId)).limit(1);
        if (reRec[0]) {
          await tx.update(userEconomy)
            .set({ balance: reRec[0].balance + change, totalGambled: (reRec[0].totalGambled ?? 0) + bet, transactions: (reRec[0].transactions ?? 0) + 1 })
            .where(eq(userEconomy.userId, userId));
        }
      });
      const finalBalance = rec ? rec.balance + change : bet;
      const emb = new EmbedBuilder()
        .setTitle('♠️ Blackjack')
        .setDescription(result)
        .addFields(
          { name: `Your Hand (${playerTotal})`, value: handDisplay(playerHand), inline: true },
          { name: `Dealer's Hand (${dealerTotal})`, value: handDisplay(dealerHand), inline: true },
          { name: 'Payout', value: `${change >= 0 ? '+' : ''}${change.toLocaleString()}`, inline: true },
          { name: 'Balance', value: `💰 **${finalBalance.toLocaleString()}** Pulses`, inline: false },
        )
        .setColor(change >= 0 ? Colors.Economy : Colors.Error);

      if (stood) {
        try {
          await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });
        } catch {
          await interaction.followUp({ embeds: [footer(timestamp(emb))], ephemeral: true });
        }
      } else {
        try {
          await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });
        } catch {
          await interaction.followUp({ embeds: [footer(timestamp(emb))], ephemeral: true });
        }
      }
    });
  },
};
