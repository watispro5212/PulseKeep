import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

const VOTE_REWARD = 500;
const VOTE_COOLDOWN = 12 * 60 * 60 * 1000;

async function checkDblVote(userId: string, botId: string, apiToken: string): Promise<boolean> {
  try {
    const res = await fetch(`https://discordbotlist.com/api/v1/bots/${botId}/votes/${userId}`, {
      headers: { Authorization: `Bot ${apiToken}` },
    });
    if (!res.ok) return false;
    const data: any = await res.json();
    return data?.voted === true;
  } catch {
    return false;
  }
}

export const voteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('vote')
    .setDescription('Vote for PulseKeep on DiscordBotList and earn Pulses')
    .toJSON(),

  async execute({ db, config }, interaction) {
    const voteUrl = 'https://discordbotlist.com/bots/1507498795569512598';
    const now = new Date();

    if (db) {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, interaction.user.id))
        .limit(1);
      const rec = rows[0];

      if (rec?.lastVote) {
        const elapsed = now.getTime() - new Date(rec.lastVote).getTime();
        if (elapsed < VOTE_COOLDOWN) {
          const remaining = Math.ceil((VOTE_COOLDOWN - elapsed) / 3600000);
          const emb = new EmbedBuilder()
            .setTitle('📊 Vote for PulseKeep')
            .setDescription(`You can claim again in **${remaining}h**. Vote on DiscordBotList to earn **${VOTE_REWARD.toLocaleString()}** Pulses!`)
            .addFields({ name: 'Vote', value: `[Click here to vote](${voteUrl})`, inline: false })
            .setColor(Colors.Utility);
          await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
          return;
        }
      }

      if (config.dblApiToken) {
        const voted = await checkDblVote(interaction.user.id, config.discordBotID, config.dblApiToken);
        if (!voted) {
          const emb = new EmbedBuilder()
            .setTitle('📊 Vote First')
            .setDescription(`You haven't voted yet! Vote for PulseKeep on DiscordBotList, then use this command again to claim your reward.`)
            .addFields({ name: 'Vote', value: `[Click here to vote](${voteUrl})`, inline: false })
            .setColor(Colors.Warning);
          await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
          return;
        }
      }

      const reward = VOTE_REWARD + Math.floor(Math.random() * 251);
      if (rec) {
        await db
          .update(userEconomy)
          .set({
            balance: (rec.balance ?? 0) + reward,
            lastVote: now,
            totalEarned: (rec.totalEarned ?? 0) + reward,
            transactions: (rec.transactions ?? 0) + 1,
          })
          .where(eq(userEconomy.userId, interaction.user.id));
      } else {
        await db
          .insert(userEconomy)
          .values({ userId: interaction.user.id, balance: reward, lastVote: now, totalEarned: reward, transactions: 1 });
      }

      const emb = new EmbedBuilder()
        .setTitle('📊 Vote Reward Claimed!')
        .setDescription(`Thanks for voting! You earned **${reward.toLocaleString()}** Pulses.`)
        .addFields(
          { name: 'Vote Again', value: `[Click here to vote](${voteUrl}) — ${VOTE_COOLDOWN / 3600000}h cooldown`, inline: false },
        )
        .setColor(Colors.Economy);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('📊 Vote for PulseKeep')
      .setDescription('Support us by voting on DiscordBotList!')
      .addFields({ name: 'Vote', value: `[Click here to vote](${voteUrl})`, inline: false })
      .setColor(Colors.Economy);
    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
