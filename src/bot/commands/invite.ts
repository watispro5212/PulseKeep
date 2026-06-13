import { SlashCommandBuilder } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors } from '../../utils/embed.js';

export const inviteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('invite')
    .setDescription('Get invite links for PulseKeep')
    .toJSON(),

  async execute(_ctx, interaction) {
    await interaction.reply({
      embeds: [{
        title: 'Invite PulseKeep',
        description: 'Add PulseKeep to your server or join the community.',
        color: Colors.Utility,
        fields: [
          { name: 'Invite Bot', value: '[Click here](https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands)', inline: false },
          { name: 'Support Server', value: '[Click here](https://discord.gg/b9HBphyeuP)', inline: false },
          { name: 'DiscordBotList', value: '[Vote for us](https://discordbotlist.com/bots/1507498795569512598)', inline: false },
          { name: 'Discords.com', value: '[Vote for us](https://discords.com/bots/bot/1507498795569512598)', inline: false },
        ],
        timestamp: new Date().toISOString(),
      }],
      flags: 64,
    });
  },
};
