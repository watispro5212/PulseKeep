import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const voteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('vote')
    .setDescription('Vote for PulseKeep on Top.gg')
    .toJSON(),

  async execute(_ctx, interaction) {
    const emb = new EmbedBuilder()
      .setTitle('📊 Vote for PulseKeep')
      .setDescription('Support us by voting on Top.gg!')
      .addFields(
        { name: 'Vote', value: '[Click here to vote](https://top.gg/bot/1507498795569512598/vote)', inline: false },
      )
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
  },
};
