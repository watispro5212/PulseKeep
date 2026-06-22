import {
  SlashCommandBuilder,
  EmbedBuilder,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const ticketpanelCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('ticketpanel')
    .setDescription('Create a ticket panel')
    .addStringOption((o) =>
      o.setName('title').setDescription('Panel title').setRequired(true),
    )
    .addStringOption((o) =>
      o.setName('description').setDescription('Panel description'),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
    .toJSON(),

  async execute({}, interaction) {
    const title = interaction.options.getString('title', true);
    const description = interaction.options.getString('description') ?? 'Click the button below to open a support ticket.';

    const embed = new EmbedBuilder()
      .setTitle(title)
      .setDescription(description)
      .setColor(Colors.Tickets);

    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(
      new ButtonBuilder()
        .setCustomId('ticket_open')
        .setLabel('Open Ticket')
        .setStyle(ButtonStyle.Primary)
        .setEmoji('🎫'),
    );

    await interaction.reply({ embeds: [footer(timestamp(embed))], components: [row] });
  },
};
