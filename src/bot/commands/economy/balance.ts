import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { STARTING_BALANCE } from '../../economy/store.js';

export const balanceCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('balance')
    .setDescription('Check your or another user\'s balance')
    .addUserOption((o) => o.setName('user').setDescription('User to check'))
    .toJSON(),

  async execute({ db }, interaction) {
    const target = interaction.options.getUser('user') ?? interaction.user;
    let balance = STARTING_BALANCE;

    if (db) {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, target.id))
        .limit(1);
      if (rows.length > 0 && rows[0]) balance = rows[0].balance;
    }

    const emb = new EmbedBuilder()
      .setTitle(`${target.username}'s Balance`)
      .setDescription(`💰 **${balance.toLocaleString()}** Pulses`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
