import { SlashCommandBuilder, EmbedBuilder, PermissionFlagsBits } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, footer, timestamp } from '../../utils/embed.js';
import { eq } from 'drizzle-orm';

export const configureCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('configure')
    .setDescription('Configure server settings for PulseKeep')
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
    .addSubcommand((s) =>
      s.setName('economy').setDescription('Toggle the economy system')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable economy').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('tickets').setDescription('Toggle the ticket system')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable tickets').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('modlogs').setDescription('Toggle moderation logging')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable mod logs').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('welcome').setDescription('Toggle welcome messages')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable welcome messages').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('log_channel').setDescription('Set the moderation log channel')
        .addChannelOption((o) => o.setName('channel').setDescription('The channel for logs').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('ticket_category').setDescription('Set the ticket category')
        .addChannelOption((o) => o.setName('category').setDescription('The category for tickets').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('welcome_channel').setDescription('Set the welcome message channel')
        .addChannelOption((o) => o.setName('channel').setDescription('The channel for welcome messages').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('vote_channel').setDescription('Set the vote announcement channel')
        .addChannelOption((o) => o.setName('channel').setDescription('The channel for vote announcements').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('show').setDescription('Show current server configuration'),
    )
    .toJSON(),

  async execute({ db }, interaction) {
    if (!interaction.guildId) {
      await interaction.reply({ content: 'This command can only be used in a server.', flags: 64 });
      return;
    }

    const sub = interaction.options.getSubcommand();
    const guildId = interaction.guildId;

    if (sub === 'show') {
      if (!db) {
        await interaction.reply({ content: 'Database unavailable.', flags: 64 });
        return;
      }
      const { guildConfigs } = await import('../../db/schema.js');
      const rows: any[] = await db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);

      const cfg = rows[0] || {};
      const emb = new EmbedBuilder()
        .setTitle('Server Configuration')
        .setDescription(`Settings for **${interaction.guild?.name || guildId}**`)
        .addFields(
          { name: 'Economy', value: cfg.economyEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Tickets', value: cfg.ticketsEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Mod Logs', value: cfg.modlogsEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Welcome', value: cfg.welcomeEnabled === true ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Vote Announcements', value: cfg.voteChannelId ? `✅ #${cfg.voteChannelId}` : '❌ Not set', inline: true },
          { name: 'Welcome Channel', value: cfg.welcomeChannelId ? `<#${cfg.welcomeChannelId}>` : 'Not set', inline: true },
          { name: 'Vote Channel', value: cfg.voteChannelId ? `<#${cfg.voteChannelId}>` : 'Not set', inline: true },
          { name: 'Log Channel', value: cfg.logChannelId ? `<#${cfg.logChannelId}>` : 'Not set', inline: true },
          { name: 'Ticket Category', value: cfg.ticketCategoryId ? `<#${cfg.ticketCategoryId}>` : 'Not set', inline: true },
        )
        .setColor(Colors.Utility);
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
      return;
    }

    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const { guildConfigs } = await import('../../db/schema.js');

    const existing: any[] = await db
      .select()
      .from(guildConfigs)
      .where(eq(guildConfigs.guildId, guildId))
      .limit(1);

    let updateData: any = {};

    switch (sub) {
      case 'economy':
        updateData.economyEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'tickets':
        updateData.ticketsEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'modlogs':
        updateData.modlogsEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'welcome':
        updateData.welcomeEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'log_channel': {
        const channel = interaction.options.getChannel('channel', true);
        updateData.logChannelId = channel.id;
        break;
      }
      case 'ticket_category': {
        const category = interaction.options.getChannel('category', true);
        updateData.ticketCategoryId = category.id;
        break;
      }
      case 'welcome_channel': {
        const wc = interaction.options.getChannel('channel', true);
        updateData.welcomeChannelId = wc.id;
        break;
      }
      case 'vote_channel': {
        const vc = interaction.options.getChannel('channel', true);
        updateData.voteChannelId = vc.id;
        break;
      }
    }

    if (existing.length > 0) {
      await db
        .update(guildConfigs)
        .set({ ...updateData, updatedAt: new Date() })
        .where(eq(guildConfigs.guildId, guildId));
    } else {
      await db
        .insert(guildConfigs)
        .values({ guildId, ...updateData });
    }

    await interaction.reply({ content: '✅ Configuration updated. Use `/configure show` to view current settings.', flags: 64 });
  },
};
