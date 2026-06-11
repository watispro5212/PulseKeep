import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userInventory } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';

export const inventoryCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('inventory')
    .setDescription('Check your inventory')
    .addUserOption((o) => o.setName('user').setDescription('User to check'))
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
      return;
    }

    const target = interaction.options.getUser('user') ?? interaction.user;
    const rows: any[] = await db
      .select()
      .from(userInventory)
      .where(eq(userInventory.userId, target.id));

    if (rows.length === 0) {
      await interaction.reply({
        embeds: [footer(timestamp(new EmbedBuilder()
          .setTitle(`${target.username}'s Inventory`)
          .setDescription('Nothing here! Visit the **/shop** to buy items.')
          .setColor(Colors.Economy)))],
        ephemeral: true,
      });
      return;
    }

    const items = rows.map(
      (r: any) => `**${r.itemName}** × ${r.quantity}`,
    );

    const emb = new EmbedBuilder()
      .setTitle(`${target.username}'s Inventory`)
      .setDescription(items.join('\n'))
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
