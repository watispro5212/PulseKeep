import { SlashCommandBuilder } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors } from '../../utils/embed.js';

export const pingCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('ping')
    .setDescription('Check the bot\'s latency')
    .toJSON(),

  async execute(ctx, interaction) {
    const latency = Date.now() - interaction.createdTimestamp;
    ctx.cache.addLatency(latency);

    await interaction.reply({
      embeds: [{
        title: 'Pong!',
        description: `**Latency:** ${latency}ms\n**API Latency:** ${Math.round(interaction.client.ws.ping)}ms`,
        color: Colors.Utility,
        timestamp: new Date().toISOString(),
      }],
      ephemeral: true,
    });
  },
};
