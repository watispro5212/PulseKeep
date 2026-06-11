import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { STARTING_BALANCE, hasXpBoost } from '../../economy/store.js';

export const balanceCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('balance')
    .setDescription('Check your or another user\'s balance')
    .addUserOption((o) => o.setName('user').setDescription('User to check'))
    .toJSON(),

  async execute({ db }, interaction) {
    const target = interaction.options.getUser('user') ?? interaction.user;
    let balance = STARTING_BALANCE;
    let extras = '';

    if (db) {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, target.id))
        .limit(1);
      if (rows.length > 0 && rows[0]) {
        const rec = rows[0];
        balance = rec.balance;
        if (hasXpBoost(rec)) {
          const remaining = Math.ceil((new Date(rec.xpBoostExpiry).getTime() - Date.now()) / 60000);
          extras += `\n⚡ XP Boost active — **${remaining}m** remaining`;
        }
        if ((rec.luckyCloverActive ?? 0) > 0) {
          extras += `\n🍀 ${rec.luckyCloverActive} Lucky Clover(s) ready`;
        }
      }
    }

    const emb = new EmbedBuilder()
      .setTitle(`${target.username}'s Balance`)
      .setDescription(`💰 **${balance.toLocaleString()}** Pulses${extras}`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
