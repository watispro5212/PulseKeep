import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const kickCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('kick')
    .setDescription('Kick a member from the server')
    .addUserOption((o) => o.setName('user').setDescription('User to kick').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for kick'))
    .setDefaultMemberPermissions(PermissionFlagsBits.KickMembers)
    .toJSON(),

  async execute({ bot }, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const target = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const member = await interaction.guild.members.fetch(target.id).catch(() => null);

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }
    if (!member.kickable) {
      await interaction.reply({ content: '❌ I cannot kick that member. Check role hierarchy.', flags: 64 });
      return;
    }

    try {
      try {
        const dm = new EmbedBuilder()
          .setTitle(`Kicked from ${interaction.guild.name}`)
          .setDescription(`You have been kicked from **${interaction.guild.name}**.`)
          .addFields({ name: 'Reason', value: reason })
          .setColor(Colors.Moderation)
          .setTimestamp();
        await target.send({ embeds: [dm] });
      } catch {}

      await member.kick(reason);
      const emb = new EmbedBuilder()
        .setTitle('Member Kicked')
        .setDescription(`**${target.username}** has been kicked.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

      const log = new EmbedBuilder()
        .setTitle('Moderation: Kick')
        .setDescription(`**${target.username}** was kicked by ${interaction.user}`)
        .addFields(
          { name: 'Reason', value: reason },
          { name: 'User ID', value: target.id, inline: true },
        )
        .setColor(Colors.Moderation)
        .setTimestamp();
      bot.logToChannel(interaction.guildId!, log);
    } catch {
      await interaction.reply({ content: '❌ Failed to kick that member.', flags: 64 });
    }
  },
};
