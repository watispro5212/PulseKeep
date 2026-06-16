import { SlashCommandBuilder, EmbedBuilder, version as djsVersion } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, Ephemeral, timestamp } from '../../utils/embed.js';

export const aboutCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('about')
    .setDescription('About PulseKeep — version, license, and links')
    .toJSON(),

  async execute({ cache, config }, interaction) {
    const uptime = Math.floor((Date.now() - cache.getStartedAt().getTime()) / 1000);
    const days = Math.floor(uptime / 86400);
    const hours = Math.floor((uptime % 86400) / 3600);
    const minutes = Math.floor((uptime % 3600) / 60);

    const emb = new EmbedBuilder()
      .setTitle('About PulseKeep')
      .setDescription(
        'A TypeScript Discord bot for moderation, audit logging, support tickets, economy games, and server analytics — all through clean slash commands.\n\n' +
        'Open source under the MIT License. Issues and feature requests welcome on GitHub.'
      )
      .setColor(Colors.Utility)
      .setThumbnail(interaction.client.user?.displayAvatarURL({ size: 256 }) ?? null)
      .addFields(
        { name: 'Version', value: 'v7.0.0', inline: true },
        { name: 'Library', value: `discord.js v${djsVersion}`, inline: true },
        { name: 'Runtime', value: `Node.js ${process.version}`, inline: true },
        { name: 'Servers', value: `${cache.getGuildsCount()}`, inline: true },
        { name: 'Users', value: `${cache.getTotalUserCount().toLocaleString()}`, inline: true },
        { name: 'Commands Run', value: `${cache.getCommandsRun().toLocaleString()}`, inline: true },
        { name: 'Uptime', value: `${days}d ${hours}h ${minutes}m`, inline: true },
        { name: 'Shards', value: `${cache.getShardCount()}`, inline: true },
        { name: 'License', value: 'MIT', inline: true },
        { name: 'Website', value: '[pulsekeep.fly.dev](https://pulsekeep.fly.dev)', inline: true },
        { name: 'Support', value: '[discord.gg/b9HBphyeuP](https://discord.gg/b9HBphyeuP)', inline: true },
        { name: 'Source', value: '[github.com/watispro5212/PulseKeep](https://github.com/watispro5212/PulseKeep)', inline: true },
      )
      .setFooter({ text: 'Built with care by the PulseKeep team.' });

    void config;
    await interaction.reply({ embeds: [timestamp(emb)], flags: Ephemeral });
  },
};
