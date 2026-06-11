import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const muteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('mute')
    .setDescription('Timeout a member')
    .addUserOption((o) => o.setName('user').setDescription('User to mute').setRequired(true))
    .addIntegerOption((o) =>
      o.setName('duration').setDescription('Duration in minutes').setRequired(true).setMinValue(1).setMaxValue(40320),
    )
    .addStringOption((o) => o.setName('reason').setDescription('Reason for mute'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute(_ctx, interaction) {
    const target = interaction.options.getUser('user', true);
    const minutes = interaction.options.getInteger('duration', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const member = interaction.guild?.members.cache.get(target.id);
    const guildName = interaction.guild?.name ?? 'the server';

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', ephemeral: true });
      return;
    }
    if (!member.moderatable) {
      await interaction.reply({ content: '❌ I cannot moderate that member. Check role hierarchy.', ephemeral: true });
      return;
    }

    try {
      try {
        const dm = new EmbedBuilder()
          .setTitle(`Muted in ${guildName}`)
          .setDescription(`You have been muted in **${guildName}** for **${minutes} minute(s)**.`)
          .addFields({ name: 'Reason', value: reason })
          .setColor(Colors.Moderation)
          .setTimestamp();
        await target.send({ embeds: [dm] });
      } catch {}

      await member.timeout(minutes * 60 * 1000, reason);
      const emb = new EmbedBuilder()
        .setTitle('Member Muted')
        .setDescription(`**${target.tag}** has been timed out for **${minutes} minute(s)**.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Duration', value: `${minutes} minute(s)`, inline: true },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
    } catch {
      await interaction.reply({ content: '❌ Failed to timeout that member.', ephemeral: true });
    }
  },
};

export const unmuteCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('unmute')
    .setDescription('Remove timeout from a member')
    .addUserOption((o) => o.setName('user').setDescription('User to unmute').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for unmute'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute(_ctx, interaction) {
    const target = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const member = interaction.guild?.members.cache.get(target.id);
    const guildName = interaction.guild?.name ?? 'the server';

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', ephemeral: true });
      return;
    }

    try {
      try {
        const dm = new EmbedBuilder()
          .setTitle(`Unmuted in ${guildName}`)
          .setDescription(`You have been unmuted in **${guildName}**.`)
          .addFields({ name: 'Reason', value: reason })
          .setColor(Colors.Moderation)
          .setTimestamp();
        await target.send({ embeds: [dm] });
      } catch {}

      await member.timeout(null, reason);
      const emb = new EmbedBuilder()
        .setTitle('Member Unmuted')
        .setDescription(`**${target.tag}** has been unmuted.`)
        .addFields(
          { name: 'Reason', value: reason, inline: false },
          { name: 'Moderator', value: `${interaction.user}`, inline: true },
        )
        .setColor(Colors.Moderation);
      await interaction.reply({ embeds: [footer(timestamp(emb))], ephemeral: true });
    } catch {
      await interaction.reply({ content: '❌ Failed to unmute that member.', ephemeral: true });
    }
  },
};
