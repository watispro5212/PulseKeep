import { SlashCommandBuilder, EmbedBuilder, ActionRowBuilder, ButtonBuilder, ButtonStyle } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

const VOTE_REWARD = 500;
const VOTE_COOLDOWN = 12 * 60 * 60 * 1000;
const DBL_URL = 'https://discordbotlist.com/bots/1507498795569512598';
const DISCORDS_URL = 'https://discords.com/bots/bots/1507498795569512598';

async function checkDblVote(userId: string, botId: string, apiToken: string): Promise<boolean | null> {
  try {
    if (!botId) return null;
    const res = await fetch(`https://discordbotlist.com/api/v1/bots/${botId}/votes/${userId}`, {
      headers: { Authorization: `Bot ${apiToken}` },
    });
    if (!res.ok) return false;
    const data: any = await res.json();
    return data?.voted === true;
  } catch {
    return null;
  }
}

export const voteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('vote')
    .setDescription('Vote for PulseKeep on bot lists and earn Pulses')
    .toJSON(),

  async execute({ db, config }, interaction) {
    const now = new Date();
    const linksRow = new ActionRowBuilder<ButtonBuilder>().addComponents(
      new ButtonBuilder().setLabel('DiscordBotList').setStyle(ButtonStyle.Link).setURL(DBL_URL),
      new ButtonBuilder().setLabel('Discords.com').setStyle(ButtonStyle.Link).setURL(DISCORDS_URL),
    );

    if (!db) {
      const emb = new EmbedBuilder()
        .setTitle('📊 Vote for PulseKeep')
        .setDescription('Support us by voting on bot lists! Vote every 12 hours to earn **500+ Pulses**.')
        .setColor(Colors.Economy);
      await interaction.reply({ embeds: [footer(timestamp(emb))], components: [linksRow], flags: 64 });
      return;
    }

    let rec: any;
    try {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, interaction.user.id))
        .limit(1);
      rec = rows[0];
    } catch {
      await interaction.reply({ content: '❌ Database error checking vote status.', flags: 64 });
      return;
    }

    if (rec?.lastVote) {
      const elapsed = now.getTime() - new Date(rec.lastVote).getTime();
      if (elapsed < VOTE_COOLDOWN) {
        const remaining = Math.ceil((VOTE_COOLDOWN - elapsed) / 3600000);
        const emb = new EmbedBuilder()
          .setTitle('📊 Vote for PulseKeep')
          .setDescription(`You can claim again in **${remaining}h**. Vote to earn **${VOTE_REWARD.toLocaleString()}** Pulses!`)
          .setColor(Colors.Utility);
        await interaction.reply({ embeds: [footer(timestamp(emb))], components: [linksRow], flags: 64 });
        return;
      }
    }

    if (config.dblApiToken) {
      const voted = await checkDblVote(interaction.user.id, config.discordBotID, config.dblApiToken);
      if (voted === null) {
        const emb = new EmbedBuilder()
          .setTitle('📊 Vote Check Unavailable')
          .setDescription('Could not verify vote status right now. Try claiming anyway!')
          .setColor(Colors.Warning);
        await interaction.reply({ embeds: [footer(timestamp(emb))], components: [linksRow], flags: 64 });
        return;
      }
      if (!voted) {
        const emb = new EmbedBuilder()
          .setTitle('📊 Vote First')
          .setDescription("You haven't voted yet! Vote on either site, then use `/vote` again to claim.")
          .setColor(Colors.Warning);
        await interaction.reply({ embeds: [footer(timestamp(emb))], components: [linksRow], flags: 64 });
        return;
      }
    }

    const reward = VOTE_REWARD + Math.floor(Math.random() * 251);
    try {
      await db.transaction(async (tx: any) => {
        const reRec = await tx.select().from(userEconomy).where(eq(userEconomy.userId, interaction.user.id)).limit(1);
        if (reRec[0]) {
          await tx
            .update(userEconomy)
            .set({
              balance: (reRec[0].balance ?? 0) + reward,
              lastVote: now,
              totalEarned: (reRec[0].totalEarned ?? 0) + reward,
              transactions: (reRec[0].transactions ?? 0) + 1,
            })
            .where(eq(userEconomy.userId, interaction.user.id));
        } else {
          await tx
            .insert(userEconomy)
            .values({ userId: interaction.user.id, balance: reward, lastVote: now, totalEarned: reward, transactions: 1 });
        }
      });
    } catch {
      await interaction.reply({ content: '❌ Failed to save vote reward.', flags: 64 });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('📊 Vote Reward Claimed!')
      .setDescription(`Thanks for voting! You earned **${reward.toLocaleString()}** Pulses.`)
      .setColor(Colors.Economy);
    await interaction.reply({ embeds: [footer(timestamp(emb))], components: [linksRow], flags: 64 });
  },
};
