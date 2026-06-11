import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { SHOP_ITEMS } from '../../economy/store.js';

export const shopCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('shop')
    .setDescription('View the PulseKeep shop')
    .toJSON(),

  async execute({ db }, interaction) {
    let balance = 500;
    if (db) {
      const { userEconomy } = await import('../../../db/schema.js');
      const { eq } = await import('drizzle-orm');
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, interaction.user.id))
        .limit(1);
      if (rows.length > 0 && rows[0]) balance = rows[0].balance;
    }

    const items = SHOP_ITEMS.map(
      (item) => `**${item.name}** — ${item.price.toLocaleString()} Pulses\n└ ${item.description}`,
    );

    const emb = new EmbedBuilder()
      .setTitle('🛒 PulseKeep Shop')
      .setDescription(`**Your Balance:** 💰 ${balance.toLocaleString()} Pulses\n\n${items.join('\n')}`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
