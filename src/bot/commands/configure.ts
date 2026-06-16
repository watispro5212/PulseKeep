import { SlashCommandBuilder, EmbedBuilder, PermissionFlagsBits, ChannelType } from 'discord.js';
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
      s.setName('automod').setDescription('Toggle the auto-moderation system')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable auto-mod').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_spam').setDescription('Toggle spam detection')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable spam detection').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_mentions').setDescription('Toggle mass mention detection')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable mass mention protection').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_caps').setDescription('Toggle excessive caps detection')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable caps enforcement').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_links').setDescription('Toggle link blocking')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable link blocking').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_words').setDescription('Toggle banned words filtering')
        .addBooleanOption((o) => o.setName('enabled').setDescription('Enable/disable banned words filter').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_banned_words').setDescription('Set banned words (comma-separated)')
        .addStringOption((o) => o.setName('words').setDescription('Comma-separated banned words').setRequired(true)),
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

  async execute({ bot, db }, interaction) {
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
          { name: 'Auto-Mod', value: cfg.automodEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Welcome Channel', value: cfg.welcomeChannelId ? `<#${cfg.welcomeChannelId}>` : 'Not set', inline: true },
          { name: 'Vote Channel', value: cfg.voteChannelId ? `<#${cfg.voteChannelId}>` : 'Not set', inline: true },
          { name: 'Log Channel', value: cfg.logChannelId ? `<#${cfg.logChannelId}>` : 'Not set', inline: true },
          { name: 'Ticket Category', value: cfg.ticketCategoryId ? `<#${cfg.ticketCategoryId}>` : 'Not set', inline: true },
        )
        .setColor(Colors.Utility);
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
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
      case 'automod':
        updateData.automodEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_spam':
        updateData.automodSpamEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_mentions':
        updateData.automodMentionEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_caps':
        updateData.automodCapsEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_links':
        updateData.automodLinkEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_words':
        updateData.automodWordsEnabled = interaction.options.getBoolean('enabled', true);
        break;
      case 'automod_banned_words':
        updateData.automodBannedWords = interaction.options.getString('words', true);
        break;
      case 'log_channel': {
        const channel = interaction.options.getChannel('channel', true);
        updateData.logChannelId = channel.id;
        break;
      }
      case 'ticket_category': {
        const category = interaction.options.getChannel('category', true);
        if (category.type !== ChannelType.GuildCategory) {
          await interaction.reply({ content: '❌ You must select a **category** channel type, not a text/voice channel.', flags: 64 });
          return;
        }
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

    bot.invalidateGuildToggles(guildId);
    await interaction.reply({ content: '✅ Configuration updated. Use `/configure show` to view current settings.', flags: 64 });
  },
};
