import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
  type Message,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const cleanCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('clean')
    .setDescription('Clean messages with filters')
    .addIntegerOption((o) =>
      o.setName('amount').setDescription('Messages to check (1-100)').setRequired(true).setMinValue(1).setMaxValue(100),
    )
    .addUserOption((o) => o.setName('user').setDescription('Only delete messages from this user'))
    .addBooleanOption((o) => o.setName('bots').setDescription('Only delete bot messages'))
    .addBooleanOption((o) => o.setName('attachments').setDescription('Only delete messages with attachments'))
    .addStringOption((o) => o.setName('contains').setDescription('Only delete messages containing this text'))
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages)
    .toJSON(),

  async execute({ bot }, interaction) {
    const amount = interaction.options.getInteger('amount', true);
    const targetUser = interaction.options.getUser('user');
    const botsOnly = interaction.options.getBoolean('bots');
    const attachmentsOnly = interaction.options.getBoolean('attachments');
    const contains = interaction.options.getString('contains');

    const channel = interaction.channel;
    if (!channel || !channel.isTextBased()) {
      await interaction.reply({ content: 'This command can only be used in a text channel.', flags: 64 });
      return;
    }

    await interaction.deferReply({ flags: 64 });

    try {
      const messages = await channel.messages.fetch({ limit: amount });
      let filtered: typeof messages = messages;

      if (targetUser) {
        filtered = filtered.filter((m: Message) => m.author.id === targetUser.id);
      }
      if (botsOnly) {
        filtered = filtered.filter((m: Message) => m.author.bot);
      }
      if (attachmentsOnly) {
        filtered = filtered.filter((m: Message) => m.attachments.size > 0);
      }
      if (contains) {
        filtered = filtered.filter((m: Message) => m.content.toLowerCase().includes(contains.toLowerCase()));
      }

      const deleted = await channel.bulkDelete(filtered, true);

      const emb = new EmbedBuilder()
        .setTitle('Messages Cleaned')
        .setDescription(`Deleted **${deleted.size}**/${filtered.size} matching messages in ${channel}.`)
        .addFields(
          ...(targetUser ? [{ name: 'User', value: `${targetUser}`, inline: true }] : []),
          ...(botsOnly ? [{ name: 'Filter', value: 'Bot messages', inline: true }] : []),
          ...(attachmentsOnly ? [{ name: 'Filter', value: 'With attachments', inline: true }] : []),
        )
        .setColor(Colors.Moderation);

      await interaction.editReply({ embeds: [footer(timestamp(emb))] });
      const log = new EmbedBuilder()
        .setTitle('Moderation: Clean')
        .setDescription(`**${deleted.size}** messages cleaned in ${channel} by ${interaction.user}${targetUser ? ` (user: ${targetUser})` : ''}${botsOnly ? ' (bots)' : ''}${attachmentsOnly ? ' (attachments)' : ''}`)
        .setColor(Colors.Moderation).setTimestamp();
      bot.logToChannel(interaction.guildId!, log);
    } catch (err) {
      await interaction.editReply({ content: `❌ Failed to clean messages: ${err instanceof Error ? err.message : err}` });
    }
  },
};
