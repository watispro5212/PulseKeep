import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const servericonCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('servericon')
    .setDescription('Get the server\'s icon')
    .toJSON(),

  async execute(_ctx, interaction) {
    const guild = interaction.guild;
    if (!guild) {
      await interaction.reply({ content: 'This command can only be used in a server.', ephemeral: true });
      return;
    }

    const iconURL = guild.iconURL({ size: 4096 });
    if (!iconURL) {
      await interaction.reply({ content: 'This server does not have a custom icon.', ephemeral: true });
      return;
    }

    const emb = new EmbedBuilder()
      .setTitle(`${guild.name} — Server Icon`)
      .setImage(iconURL)
      .setDescription(`[Open in browser](${iconURL})`)
      .setColor(Colors.Utility);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
