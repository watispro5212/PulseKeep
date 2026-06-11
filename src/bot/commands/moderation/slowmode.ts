import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const slowmodeCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('slowmode')
    .setDescription('Set slowmode in this channel')
    .addIntegerOption((o) =>
      o.setName('seconds').setDescription('Slowmode in seconds (0 to disable)').setRequired(true).setMinValue(0).setMaxValue(21600),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageChannels)
    .toJSON(),

  async execute(_ctx, interaction) {
    const seconds = interaction.options.getInteger('seconds', true);
    const channel = interaction.channel;
    if (!channel || !('setRateLimitPerUser' in channel)) {
      await interaction.reply({ content: 'This channel does not support slowmode.', ephemeral: true });
      return;
    }

    try {
      await (channel as any).setRateLimitPerUser(seconds);
      const emb = new EmbedBuilder()
        .setTitle('Slowmode Updated')
        .setDescription(seconds > 0
          ? `Slowmode set to **${seconds}s** in ${channel}.`
          : `Slowmode **disabled** in ${channel}.`)
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
    } catch {
      await interaction.reply({ content: '❌ Failed to set slowmode.', ephemeral: true });
    }
  },
};
