import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp, Ephemeral } from '../../../utils/embed.js';
import { userInventory, userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

export const inventoryCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('inventory')
    .setDescription('Check your inventory')
    .addUserOption((o) => o.setName('user').setDescription('User to check'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: Ephemeral });
      return;
    }

    const target = interaction.options.getUser('user') ?? interaction.user;

    const [invRows, econRows] = await Promise.all([
      db.select().from(userInventory).where(eq(userInventory.userId, target.id)),
      db.select().from(userEconomy).where(eq(userEconomy.userId, target.id)).limit(1),
    ]);

    const econ = econRows[0];
    const activeStatus: string[] = [];

    if (econ?.xpBoostExpiry) {
      const remaining = new Date(econ.xpBoostExpiry).getTime() - Date.now();
      if (remaining > 0) {
        const mins = Math.ceil(remaining / 60000);
        activeStatus.push(`⚡ **XP Boost** — ${mins}m remaining`);
      }
    }
    if ((econ?.luckyCloverActive ?? 0) > 0) {
      activeStatus.push(`🍀 **Lucky Clover** — ${econ.luckyCloverActive} ready`);
    }

    if (invRows.length === 0 && activeStatus.length === 0) {
      await interaction.reply({
        embeds: [footer(timestamp(new EmbedBuilder()
          .setTitle(`${target.username}'s Inventory`)
          .setDescription('Nothing here! Visit the **/shop** to buy items.')
          .setColor(Colors.Economy)))],
        flags: Ephemeral,
      });
      return;
    }

    const items = invRows.map(
      (r: any) => `**${r.itemName}** × ${r.quantity}`,
    );

    const emb = new EmbedBuilder()
      .setTitle(`${target.username}'s Inventory`)
      .setDescription(activeStatus.length > 0
        ? `${activeStatus.join('\n')}\n\n**Items:**\n${items.join('\n')}`
        : items.join('\n'))
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: Ephemeral });
  },
};
