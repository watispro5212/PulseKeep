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

  async execute({ bot }, interaction) {
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
          .setDescription(
            `Ticket closed by ${interaction.user}.\n\n` +
            'Saving transcript and deleting this channel in **10 seconds**…'
          )
          .setColor(Colors.Tickets);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });

        // try to send a transcript to the mod channel
        try {
          if ('messages' in channel && interaction.guildId) {
            const messages = await (channel as any).messages.fetch({ limit: 100 }).catch(() => null);
            if (messages && messages.size > 0) {
              const transcript = messages
                .sort((a: any, b: any) => a.createdTimestamp - b.createdTimestamp)
                .map((m: any) => `[${new Date(m.createdTimestamp).toISOString()}] ${m.author?.username ?? 'unknown'}: ${m.content ?? ''}`)
                .join('\n')
                .slice(0, 6000); // hard cap so we don't explode a single embed field
              const transcriptEmb = new EmbedBuilder()
                .setTitle(`Ticket transcript — #${channel.name}`)
                .setDescription('```\n' + transcript + '\n```')
                .setColor(Colors.Tickets)
                .setFooter({ text: `${messages.size} message(s) • closed by ${interaction.user.username}` });
              if (bot) await bot.logToChannel(interaction.guildId, transcriptEmb);
            }
          }
        } catch (err) {
          console.error('[ticket close] transcript failed:', err);
        }

        setTimeout(() => channel.delete().catch(() => {}), 10_000);
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
