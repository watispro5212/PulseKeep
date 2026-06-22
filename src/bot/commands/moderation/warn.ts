import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userWarnings } from '../../../db/schema.js';

export const warnCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('warn')
    .setDescription('Warn a user')
    .addUserOption((o) => o.setName('user').setDescription('User to warn').setRequired(true))
    .addStringOption((o) => o.setName('reason').setDescription('Reason for the warning'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
    .toJSON(),

  async execute({ db, bot }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }

    const user = interaction.options.getUser('user', true);
    const reason = interaction.options.getString('reason') ?? 'No reason provided';
    const moderatorId = interaction.user.id;
    const guildName = interaction.guild?.name ?? 'the server';

    await db.insert(userWarnings).values({
      guildId: interaction.guildId!,
      userId: user.id,
      moderatorId,
      reason,
    });

    let dmFailed = false;
    try {
      const dm = new EmbedBuilder()
        .setTitle(`Warning from ${guildName}`)
        .setDescription(`You have been warned in **${guildName}**.`)
        .addFields({ name: 'Reason', value: reason })
        .setColor(Colors.Moderation)
        .setTimestamp();
      await user.send({ embeds: [dm] });
    } catch {
      dmFailed = true;
    }

    const emb = new EmbedBuilder()
      .setTitle('User Warned')
      .setDescription(`**${user.username}** has been warned.`)
      .addFields(
        { name: 'Reason', value: reason, inline: false },
        { name: 'Moderator', value: `<@${moderatorId}>`, inline: true },
        { name: 'User', value: `${user}`, inline: true },
      )
      .setColor(Colors.Moderation);
    if (dmFailed) emb.setFooter({ text: '⚠️ Could not DM the user (DMs may be closed).' });

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

    // log it
    const log = new EmbedBuilder()
      .setTitle('Moderation: Warn')
      .setDescription(`**${user.username}** was warned by ${interaction.user}`)
      .addFields(
        { name: 'Reason', value: reason },
        { name: 'User ID', value: user.id, inline: true },
        { name: 'Moderator ID', value: moderatorId, inline: true },
      )
      .setColor(Colors.Moderation)
      .setTimestamp();
    await bot.logToChannel(interaction.guildId!, log);
  },
};
