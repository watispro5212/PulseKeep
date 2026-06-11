import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy, userInventory } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { SHOP_ITEMS } from '../../economy/store.js';

export const buyCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('buy')
    .setDescription('Buy an item from the shop')
    .addStringOption((o) =>
      o.setName('item').setDescription('Item to buy').setRequired(true)
        .addChoices(...SHOP_ITEMS.map((i) => ({ name: i.name, value: i.id }))),
    )
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', ephemeral: true });
      return;
    }

    const itemId = interaction.options.getString('item', true);
    const item = SHOP_ITEMS.find((i) => i.id === itemId);
    if (!item) {
      await interaction.reply({ content: '❌ Item not found.', ephemeral: true });
      return;
    }

    const userId = interaction.user.id;
    const rows = await db
      .select()
      .from(userEconomy)
      .where(eq(userEconomy.userId, userId))
      .limit(1);
    const rec = rows[0];

    if (!rec || rec.balance < item.price) {
      await interaction.reply({ content: `❌ You need **${item.price.toLocaleString()}** Pulses for ${item.name}.`, ephemeral: true });
      return;
    }

    await db
      .update(userEconomy)
      .set({
        balance: rec.balance - item.price,
        transactions: (rec.transactions ?? 0) + 1,
      })
      .where(eq(userEconomy.userId, userId));

    const existing = await db
      .select()
      .from(userInventory)
      .where(eq(userInventory.userId, userId))
      .where(eq(userInventory.itemId, itemId))
      .limit(1);

    if (existing.length > 0 && existing[0]) {
      await db
        .update(userInventory)
        .set({ quantity: existing[0].quantity + 1 })
        .where(eq(userInventory.id, existing[0].id));
    } else {
      await db
        .insert(userInventory)
        .values({ userId, itemId: item.id, itemName: item.name, quantity: 1 });
    }

    const emb = new EmbedBuilder()
      .setTitle('Purchase Successful')
      .setDescription(`You bought **${item.name}** for **${item.price.toLocaleString()}** Pulses!`)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
