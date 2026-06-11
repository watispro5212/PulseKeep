import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { ECONOMY_TIPS } from '../../economy/store.js';

export const tipCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('tip')
    .setDescription('Get a random economy tip')
    .toJSON(),

  async execute(_ctx, interaction) {
    const tip = ECONOMY_TIPS[Math.floor(Math.random() * ECONOMY_TIPS.length)] ?? 'Use /daily every 24h for free Pulses!';

    const emb = new EmbedBuilder()
      .setTitle('💡 Economy Tip')
      .setDescription(tip)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
