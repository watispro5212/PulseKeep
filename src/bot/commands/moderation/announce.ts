import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors } from '../../../utils/embed.js';

export const announceCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('announce')
    .setDescription('Send an announcement to this channel')
    .addStringOption((o) => o.setName('title').setDescription('Announcement title').setRequired(true))
    .addStringOption((o) => o.setName('message').setDescription('Announcement content').setRequired(true))
    .addStringOption((o) =>
      o.setName('color').setDescription('Embed color')
        .addChoices(
          { name: 'Purple', value: 'purple' },
          { name: 'Blue', value: 'blue' },
          { name: 'Green', value: 'green' },
          { name: 'Red', value: 'red' },
          { name: 'Yellow', value: 'yellow' },
        ),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
    .toJSON(),

  async execute(_ctx, interaction) {
    const title = interaction.options.getString('title', true);
    const message = interaction.options.getString('message', true);
    const colorChoice = interaction.options.getString('color') ?? 'purple';

    const colors: Record<string, number> = {
      purple: Colors.Moderation,
      blue: 0x3b82f6,
      green: Colors.Economy,
      red: Colors.Error,
      yellow: Colors.Warning,
    };

    try {
      const emb = new EmbedBuilder()
        .setTitle(title)
        .setDescription(message)
        .setColor(colors[colorChoice] ?? Colors.Moderation)
        .setFooter({ text: `Announcement by ${interaction.user.tag}` })
        .setTimestamp(new Date());

      await interaction.reply({ embeds: [emb] });
    } catch {
      await interaction.reply({ content: '❌ Failed to send announcement.', flags: 64 });
    }
  },
};
