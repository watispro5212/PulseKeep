import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const nickCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('nick')
    .setDescription('Change a member\'s nickname')
    .addUserOption((o) => o.setName('user').setDescription('User').setRequired(true))
    .addStringOption((o) => o.setName('nickname').setDescription('New nickname (leave empty to reset)'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageNicknames)
    .toJSON(),

  async execute(_ctx, interaction) {
    const target = interaction.options.getUser('user', true);
    const nickname = interaction.options.getString('nickname');
    const member = interaction.guild?.members.cache.get(target.id);

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }
    if (!member.moderatable) {
      await interaction.reply({ content: '❌ I cannot change that member\'s nickname.', flags: 64 });
      return;
    }

    try {
      await member.setNickname(nickname ?? null);

      const emb = new EmbedBuilder()
        .setTitle('Nickname Updated')
        .setDescription(nickname
          ? `Changed **${target.tag}**'s nickname to **${nickname}**.`
          : `Reset **${target.tag}**'s nickname.`,
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    } catch {
      await interaction.reply({ content: '❌ Failed to change nickname.', flags: 64 });
    }
  },
};
