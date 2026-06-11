import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const ticketCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('ticket')
    .setDescription('Manage tickets')
    .addSubcommand((s) =>
      s.setName('add').setDescription('Add a user to the ticket')
        .addUserOption((o) => o.setName('user').setDescription('User to add').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('remove').setDescription('Remove a user from the ticket')
        .addUserOption((o) => o.setName('user').setDescription('User to remove').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('close').setDescription('Close this ticket'),
    )
    .addSubcommand((s) =>
      s.setName('rename').setDescription('Rename this ticket')
        .addStringOption((o) => o.setName('name').setDescription('New name').setRequired(true)),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageChannels)
    .toJSON(),

  async execute(_ctx, interaction) {
    const sub = interaction.options.getSubcommand();

    const channel = interaction.channel;
    if (!channel || !channel.name.startsWith('ticket-')) {
      await interaction.reply({ content: '❌ This command must be used inside a ticket channel.', flags: 64 });
      return;
    }
    if (!channel.isTextBased() || !('permissionOverwrites' in channel)) {
      await interaction.reply({ content: '❌ Invalid channel type.', flags: 64 });
      return;
    }

    switch (sub) {
      case 'add': {
        const user = interaction.options.getUser('user', true);
        await channel.permissionOverwrites.create(user.id, {
          ViewChannel: true,
          SendMessages: true,
          ReadMessageHistory: true,
        });
        const emb = new EmbedBuilder()
          .setTitle('User Added')
          .setDescription(`Added ${user} to the ticket.`)
          .setColor(Colors.Tickets);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
        break;
      }
      case 'remove': {
        const user = interaction.options.getUser('user', true);
        const overwrite = channel.permissionOverwrites.cache.get(user.id);
        if (!overwrite) {
          await interaction.reply({ content: '❌ That user does not have access to this ticket.', flags: 64 });
          return;
        }
        await overwrite.delete();
        const emb = new EmbedBuilder()
          .setTitle('User Removed')
          .setDescription(`Removed ${user} from the ticket.`)
          .setColor(Colors.Tickets);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
        break;
      }
      case 'close': {
        const emb = new EmbedBuilder()
          .setTitle('Ticket Closed')
          .setDescription(`Ticket closed by ${interaction.user}. This channel will be deleted shortly.`)
          .setColor(Colors.Tickets);
        await interaction.reply({ embeds: [footer(timestamp(emb))] });
        setTimeout(() => channel.delete().catch(() => {}), 3000);
        break;
      }
      case 'rename': {
        const name = interaction.options.getString('name', true);
        await channel.setName(name);
        const emb = new EmbedBuilder()
          .setTitle('Ticket Renamed')
          .setDescription(`Renamed to **${name}**.`)
          .setColor(Colors.Tickets);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
        break;
      }
    }
  },
};
