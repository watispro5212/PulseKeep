import { SlashCommandBuilder, EmbedBuilder, ChannelType } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

const typeNames: Record<number, string> = {
  [ChannelType.GuildText]: 'Text',
  [ChannelType.GuildVoice]: 'Voice',
  [ChannelType.GuildCategory]: 'Category',
  [ChannelType.GuildAnnouncement]: 'Announcement',
  [ChannelType.GuildForum]: 'Forum',
  [ChannelType.GuildStageVoice]: 'Stage',
};

export const channelinfoCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('channelinfo')
    .setDescription('Get information about a channel')
    .addChannelOption((o) => o.setName('channel').setDescription('Channel to inspect'))
    .toJSON(),

  async execute({}, interaction) {
    const channel = interaction.options.getChannel('channel') ?? interaction.channel;
    if (!channel) {
      await interaction.reply({ content: 'Could not find channel.', flags: 64 });
      return;
    }

    const typeName = typeNames[channel.type] ?? 'Unknown';
    const created = new Date(channel.createdTimestamp ?? Date.now()).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
    });

    const emb = new EmbedBuilder()
      .setTitle(`#${channel.name} — Channel Info`)
      .setDescription(`Information about ${channel}.`)
      .setColor(Colors.Utility)
      .addFields(
        { name: 'Type', value: typeName, inline: true },
        { name: 'Channel ID', value: channel.id, inline: true },
        { name: 'Created', value: created, inline: true },
        { name: 'NSFW', value: ('nsfw' in channel && channel.nsfw) ? 'Yes' : 'No', inline: true },
      );

    if (interaction.guild) {
      emb.addFields({ name: 'Server', value: interaction.guild.name, inline: true });
    }

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
