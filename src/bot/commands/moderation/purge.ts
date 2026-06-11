import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const purgeCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('purge')
    .setDescription('Bulk delete messages in this channel')
    .addIntegerOption((o) =>
      o.setName('amount').setDescription('Messages to delete (1-100)').setRequired(true).setMinValue(1).setMaxValue(100),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages)
    .toJSON(),

  async execute(_ctx, interaction) {
    const amount = interaction.options.getInteger('amount', true);
    const channel = interaction.channel;
    if (!channel || !channel.isTextBased()) {
      await interaction.reply({ content: 'This command can only be used in a text channel.', flags: 64 });
      return;
    }

    try {
      const messages = await channel.bulkDelete(amount, true);
      const emb = new EmbedBuilder()
        .setTitle('Messages Purged')
        .setDescription(`Deleted **${messages.size}** messages in ${channel}.`)
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    } catch {
      await interaction.reply({ content: '❌ Failed to purge messages. Messages may be too old (older than 14 days).', flags: 64 });
    }
  },
};
