import { SlashCommandBuilder, EmbedBuilder, version as djsVersion } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer } from '../../../utils/embed.js';

export const statsCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('stats')
    .setDescription('View bot statistics')
    .toJSON(),

  async execute({ cache }, interaction) {
    const uptimeSeconds = Math.floor((Date.now() - cache.getStartedAt().getTime()) / 1000);
    const days = Math.floor(uptimeSeconds / 86400);
    const hours = Math.floor((uptimeSeconds % 86400) / 3600);
    const minutes = Math.floor((uptimeSeconds % 3600) / 60);

    const emb = new EmbedBuilder()
      .setTitle('PulseKeep Statistics')
      .setColor(Colors.Utility)
      .addFields(
        { name: 'Servers', value: `${cache.getGuildsCount()}`, inline: true },
        { name: 'Users', value: `${cache.getTotalUserCount()}`, inline: true },
        { name: 'Commands Run', value: `${cache.getCommandsRun()}`, inline: true },
        { name: 'Avg Latency', value: `${Math.round(cache.getAvgLatency())}ms`, inline: true },
        { name: 'Uptime', value: `${days}d ${hours}h ${minutes}m`, inline: true },
        { name: 'Version', value: 'v7.0.0', inline: true },
        { name: 'Library', value: `discord.js v${djsVersion}`, inline: true },
        { name: 'Runtime', value: `Node.js ${process.version}`, inline: true },
      );

    await interaction.reply({ embeds: [footer(emb).setTimestamp(new Date())], flags: 64 });
  },
};
