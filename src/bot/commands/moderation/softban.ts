import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const softbanCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('softban')
    .setDescription('Ban then unban a user to clear their messages')
    .addUserOption((o) => o.setName('user').setDescription('User to softban').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for softban'))
    .setDefaultMemberPermissions(PermissionFlagsBits.BanMembers)
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
    if (!member.bannable) {
      await interaction.reply({ content: '❌ I cannot ban that member. Check role hierarchy.', flags: 64 });
      return;
    }

    try {
      await target.send(`You have been softbanned from **${interaction.guild.name}**.\nReason: ${reason}`).catch(() => {});
      await member.ban({ reason: `[Softban] ${reason}`, deleteMessageSeconds: 86400 });
      await interaction.guild.members.unban(target.id, `[Softban completed] ${reason}`);

      const emb = new EmbedBuilder()
        .setTitle('Member Softbanned')
        .setDescription(`**${target.username}** was banned and immediately unbanned. Their recent messages have been removed.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

      const log = new EmbedBuilder()
        .setTitle('Member Softbanned')
        .setDescription(`**${target.username}** was softbanned.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation)
        .setTimestamp();
      await bot.logToChannel(interaction.guildId!, log);
    } catch {
      await interaction.reply({ content: '❌ Failed to softban that member.', flags: 64 });
    }
  },
};
