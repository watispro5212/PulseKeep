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
      s.setName('automod_spam_limit').setDescription('Max messages in spam window')
        .addIntegerOption((o) => o.setName('limit').setDescription('Number of messages').setRequired(true).setMinValue(2).setMaxValue(50)),
    )
    .addSubcommand((s) =>
      s.setName('automod_spam_window').setDescription('Spam detection window in seconds')
        .addIntegerOption((o) => o.setName('seconds').setDescription('Time window').setRequired(true).setMinValue(1).setMaxValue(120)),
    )
    .addSubcommand((s) =>
      s.setName('automod_mention_limit').setDescription('Max mentions before action')
        .addIntegerOption((o) => o.setName('limit').setDescription('Number of mentions').setRequired(true).setMinValue(1).setMaxValue(50)),
    )
    .addSubcommand((s) =>
      s.setName('automod_caps_ratio').setDescription('Caps percentage threshold (0-100)')
        .addIntegerOption((o) => o.setName('percent').setDescription('Percentage of caps').setRequired(true).setMinValue(10).setMaxValue(100)),
    )
    .addSubcommand((s) =>
      s.setName('automod_caps_min_length').setDescription('Min message length for caps check')
        .addIntegerOption((o) => o.setName('length').setDescription('Min characters').setRequired(true).setMinValue(3).setMaxValue(100)),
    )
    .addSubcommand((s) =>
      s.setName('automod_timeout_duration').setDescription('Timeout duration in minutes')
        .addIntegerOption((o) => o.setName('minutes').setDescription('Minutes').setRequired(true).setMinValue(1).setMaxValue(40320)),
    )
    .addSubcommand((s) =>
      s.setName('automod_action').setDescription('Set action for a specific rule')
        .addStringOption((o) => o.setName('rule').setDescription('Rule to configure').setRequired(true)
          .addChoices(
            { name: 'Spam', value: 'spam' },
            { name: 'Mentions', value: 'mentions' },
            { name: 'Caps', value: 'caps' },
            { name: 'Links', value: 'links' },
            { name: 'Banned Words', value: 'words' },
          ))
        .addStringOption((o) => o.setName('action').setDescription('Action to take').setRequired(true)
          .addChoices(
            { name: 'Warn', value: 'warn' },
            { name: 'Delete', value: 'delete' },
            { name: 'Timeout', value: 'timeout' },
          )),
    )
    .addSubcommand((s) =>
      s.setName('automod_exempt_roles').setDescription('Role IDs exempt from automod (comma-sep)')
        .addStringOption((o) => o.setName('roles').setDescription('Comma-separated role IDs').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('automod_allowed_domains').setDescription('Domains allowed through link filter (comma-sep)')
        .addStringOption((o) => o.setName('domains').setDescription('Comma-separated domains (e.g. discord.com,github.com)').setRequired(true)),
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
      let cfg: any = {};
      try {
        const rows: any[] = await db
          .select()
          .from(guildConfigs)
          .where(eq(guildConfigs.guildId, guildId))
          .limit(1);
        cfg = rows[0] || {};
      } catch {
        cfg = {};
      }

      const autoToggles = [
        `Spam: ${cfg.automodSpamEnabled !== false ? '✅' : '❌'}`,
        `Mentions: ${cfg.automodMentionEnabled !== false ? '✅' : '❌'}`,
        `Caps: ${cfg.automodCapsEnabled !== false ? '✅' : '❌'}`,
        `Links: ${cfg.automodLinkEnabled !== false ? '✅' : '❌'}`,
        `Words: ${cfg.automodWordsEnabled !== false ? '✅' : '❌'}`,
        cfg.automodBannedWords ? `Banned: \`${cfg.automodBannedWords.slice(0, 60)}${cfg.automodBannedWords.length > 60 ? '...' : ''}\`` : null,
      ].filter(Boolean).join('\n');
      const autoThresholds = [
        `Spam limit: ${cfg.automodSpamLimit ?? 5} msgs / ${(cfg.automodSpamWindow ?? 5000) / 1000}s`,
        `Mention limit: ${cfg.automodMentionLimit ?? 5}`,
        `Caps: ≥${cfg.automodCapsRatio ?? 70}% / min ${cfg.automodCapsMinLength ?? 10} chars`,
        `Timeout: ${cfg.automodTimeoutDuration ?? 10} min`,
        cfg.automodExemptRoles ? `Exempt roles: ${cfg.automodExemptRoles}` : null,
        cfg.automodAllowedDomains ? `Allowed domains: \`${cfg.automodAllowedDomains.slice(0, 60)}${cfg.automodAllowedDomains.length > 60 ? '...' : ''}\`` : null,
      ].filter(Boolean).join('\n');
      const autoActions = [
        `Spam → ${cfg.automodSpamAction || 'warn'}`,
        `Mentions → ${cfg.automodMentionAction || 'timeout'}`,
        `Caps → ${cfg.automodCapsAction || 'delete'}`,
        `Links → ${cfg.automodLinkAction || 'delete'}`,
        `Words → ${cfg.automodWordsAction || 'delete'}`,
      ].join('\n');
      const emb = new EmbedBuilder()
        .setTitle('Server Configuration')
        .setDescription(`Settings for **${interaction.guild?.name || guildId}**`)
        .addFields(
          { name: 'Economy', value: cfg.economyEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Tickets', value: cfg.ticketsEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Mod Logs', value: cfg.modlogsEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Welcome', value: cfg.welcomeEnabled === true ? '✅ Enabled' : '❌ Disabled', inline: true },
          { name: 'Auto-Mod', value: cfg.automodEnabled !== false ? '✅ Enabled' : '❌ Disabled', inline: false },
          { name: 'Auto-Mod Filters', value: autoToggles || 'None set', inline: false },
          { name: 'Auto-Mod Thresholds', value: autoThresholds || 'Default', inline: false },
          { name: 'Auto-Mod Actions', value: autoActions, inline: false },
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

    let existing: any[] = [];
    try {
      existing = await db
        .select()
        .from(guildConfigs)
        .where(eq(guildConfigs.guildId, guildId))
        .limit(1);
    } catch {
      await interaction.reply({ content: '❌ Database schema outdated. Ask the bot owner to run `npm run db:migrate`.', flags: 64 });
      return;
    }

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
      case 'automod_spam_limit':
        updateData.automodSpamLimit = interaction.options.getInteger('limit', true);
        break;
      case 'automod_spam_window':
        updateData.automodSpamWindow = interaction.options.getInteger('seconds', true) * 1000;
        break;
      case 'automod_mention_limit':
        updateData.automodMentionLimit = interaction.options.getInteger('limit', true);
        break;
      case 'automod_caps_ratio':
        updateData.automodCapsRatio = interaction.options.getInteger('percent', true);
        break;
      case 'automod_caps_min_length':
        updateData.automodCapsMinLength = interaction.options.getInteger('length', true);
        break;
      case 'automod_timeout_duration':
        updateData.automodTimeoutDuration = interaction.options.getInteger('minutes', true);
        break;
      case 'automod_action': {
        const rule = interaction.options.getString('rule', true) as keyof typeof colMap;
        const action = interaction.options.getString('action', true);
        const colMap = {
          spam: 'automodSpamAction',
          mentions: 'automodMentionAction',
          caps: 'automodCapsAction',
          links: 'automodLinkAction',
          words: 'automodWordsAction',
        } as const;
        const col = colMap[rule];
        if (col) updateData[col] = action;
        break;
      }
      case 'automod_exempt_roles':
        updateData.automodExemptRoles = interaction.options.getString('roles', true);
        break;
      case 'automod_allowed_domains':
        updateData.automodAllowedDomains = interaction.options.getString('domains', true);
        break;
      case 'log_channel': {
        const channel = interaction.options.getChannel('channel', true);
        if (!channel.isTextBased() || channel.isVoiceBased()) {
          await interaction.reply({ content: '❌ Please select a text channel for logs.', flags: 64 });
          return;
        }
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
        if (!wc.isTextBased() || wc.isVoiceBased()) {
          await interaction.reply({ content: '❌ Please select a text channel for welcome messages.', flags: 64 });
          return;
        }
        updateData.welcomeChannelId = wc.id;
        break;
      }
      case 'vote_channel': {
        const vc = interaction.options.getChannel('channel', true);
        if (!vc.isTextBased() || vc.isVoiceBased()) {
          await interaction.reply({ content: '❌ Please select a text channel for vote announcements.', flags: 64 });
          return;
        }
        updateData.voteChannelId = vc.id;
        break;
      }
    }

    try {
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
    } catch {
      await interaction.reply({ content: '❌ Failed to save configuration. Database error.', flags: 64 });
      return;
    }

    bot.invalidateGuildToggles(guildId);
    await interaction.reply({ content: '✅ Configuration updated. Use `/configure show` to view current settings.', flags: 64 });
  },
};
