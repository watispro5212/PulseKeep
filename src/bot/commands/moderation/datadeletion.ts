import {
  SlashCommandBuilder,
  EmbedBuilder,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ComponentType,
  type MessageComponentInteraction,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp, Ephemeral } from '../../../utils/embed.js';
import { userEconomy, userWarnings, userInventory, commandLogs } from '../../../db/schema.js';
import { eq, sql } from 'drizzle-orm';

export const dataDeletionCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('data-deletion')
    .setDescription('Wipe all data for a user by their ID (owner/co-owner only)')
    .addStringOption((o) =>
      o.setName('user_id').setDescription('Discord user ID to delete').setRequired(true),
    )
    .toJSON(),

  async execute({ bot, db, config }, interaction) {
    if (interaction.user.id !== config.botOwnerID && interaction.user.id !== config.botCoOwnerID) {
      await interaction.reply({ content: '❌ Only the bot owner or co-owner can use this command.', flags: Ephemeral });
      return;
    }

    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: Ephemeral });
      return;
    }

    const userId = interaction.options.getString('user_id', true).trim();
    if (!/^\d{17,19}$/.test(userId)) {
      await interaction.reply({ content: '❌ Invalid Discord user ID. Must be a 17-19 digit numeric ID.', flags: Ephemeral });
      return;
    }

    let economyCount = 0, warningCount = 0, inventoryCount = 0, logCount = 0;
    try {
      const [eco] = await db.select({ count: sql<number>`count(*)::int` }).from(userEconomy).where(eq(userEconomy.userId, userId));
      economyCount = eco?.count ?? 0;
      const [warn] = await db.select({ count: sql<number>`count(*)::int` }).from(userWarnings).where(eq(userWarnings.userId, userId));
      warningCount = warn?.count ?? 0;
      const [inv] = await db.select({ count: sql<number>`count(*)::int` }).from(userInventory).where(eq(userInventory.userId, userId));
      inventoryCount = inv?.count ?? 0;
      const [log] = await db.select({ count: sql<number>`count(*)::int` }).from(commandLogs).where(eq(commandLogs.userId, userId));
      logCount = log?.count ?? 0;
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to query user data: ${err instanceof Error ? err.message : err}`, flags: Ephemeral });
      return;
    }

    if (economyCount + warningCount + inventoryCount + logCount === 0) {
      await interaction.reply({ content: `❌ No data found for user ID \`${userId}\`.`, flags: Ephemeral });
      return;
    }

    const confirm = new ButtonBuilder()
      .setCustomId('confirm_delete')
      .setLabel('Yes, delete all data')
      .setStyle(ButtonStyle.Danger);
    const cancel = new ButtonBuilder()
      .setCustomId('cancel_delete')
      .setLabel('Cancel')
      .setStyle(ButtonStyle.Secondary);
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(confirm, cancel);

    const prompt = new EmbedBuilder()
      .setTitle('⚠️ Confirm Data Deletion')
      .setDescription(`Delete **all** data for user ID \`${userId}\`? This cannot be undone.`)
      .addFields(
        { name: 'Economy', value: `${economyCount} record(s)`, inline: true },
        { name: 'Warnings', value: `${warningCount} record(s)`, inline: true },
        { name: 'Inventory', value: `${inventoryCount} record(s)`, inline: true },
        { name: 'Command Logs', value: `${logCount} record(s)`, inline: true },
      )
      .setColor(Colors.Warning);
    await interaction.reply({ embeds: [footer(timestamp(prompt))], components: [row], flags: Ephemeral });

    const reply = await interaction.fetchReply();
    const collected = await reply.awaitMessageComponent({
      componentType: ComponentType.Button,
      time: 30000,
      filter: (i: MessageComponentInteraction) => i.user.id === interaction.user.id,
    }).catch(() => null);

    if (!collected || collected.customId === 'cancel_delete') {
      await interaction.editReply({ content: '❌ Cancelled.', embeds: [], components: [] });
      return;
    }

    await collected.deferUpdate();

    let deletedEco = 0, deletedWarn = 0, deletedInv = 0, deletedLog = 0;
    try {
      await db.transaction(async (tx: any) => {
        const ecoResult = await tx.delete(userEconomy).where(eq(userEconomy.userId, userId));
        deletedEco = (ecoResult as any)?.rowCount ?? 0;
        const warnResult = await tx.delete(userWarnings).where(eq(userWarnings.userId, userId));
        deletedWarn = (warnResult as any)?.rowCount ?? 0;
        const invResult = await tx.delete(userInventory).where(eq(userInventory.userId, userId));
        deletedInv = (invResult as any)?.rowCount ?? 0;
        const logResult = await tx.delete(commandLogs).where(eq(commandLogs.userId, userId));
        deletedLog = (logResult as any)?.rowCount ?? 0;
      });
    } catch (err) {
      await interaction.editReply({ content: `❌ Failed to delete data: ${err instanceof Error ? err.message : err}`, components: [] });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle('Data Deletion Complete')
      .setDescription(`Data for user ID \`${userId}\` has been wiped.`)
      .addFields(
        { name: 'Economy', value: `${deletedEco} row(s)`, inline: true },
        { name: 'Warnings', value: `${deletedWarn} row(s)`, inline: true },
        { name: 'Inventory', value: `${deletedInv} row(s)`, inline: true },
        { name: 'Command Logs', value: `${deletedLog} row(s)`, inline: true },
      )
      .setColor(Colors.Moderation);
    await interaction.editReply({ embeds: [footer(timestamp(emb))], components: [] });

    const log = new EmbedBuilder()
      .setTitle('Data Deletion')
      .setDescription(`User ID \`${userId}\` data deleted.`)
      .addFields(
        { name: 'Requested by', value: `${interaction.user}`, inline: true },
        { name: 'User ID', value: userId, inline: true },
      )
      .setColor(Colors.Moderation)
      .setTimestamp();
    if (interaction.guildId) await bot.logToChannel(interaction.guildId!, log);
  },
};
