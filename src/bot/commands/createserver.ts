import { SlashCommandBuilder, PermissionFlagsBits, ChannelType } from 'discord.js';
import type { SlashCommand } from '../types.js';

const ROLES = [
  { name: '👑 Owner', color: 0x7c5cfc, reason: 'Owner of PulseKeep', permissions: [PermissionFlagsBits.Administrator] },
  { name: '💎 Co-Owner', color: 0x4363d8, reason: 'Co-owner backup authority', permissions: [PermissionFlagsBits.Administrator] },
  { name: '🛡️ Administrator', color: 0xe6194b, reason: 'Server structure and permissions', permissions: [PermissionFlagsBits.Administrator] },
  { name: '🧑‍💻 Developer', color: 0x4b0082, reason: 'Code and bot debugging', permissions: [PermissionFlagsBits.ViewAuditLog] },
  { name: '🔐 Security Lead', color: 0x8b0000, reason: 'Abuse and exploit reports', permissions: [PermissionFlagsBits.ViewAuditLog, PermissionFlagsBits.KickMembers, PermissionFlagsBits.ModerateMembers] },
  { name: '🧰 Support Manager', color: 0x00d2d3, reason: 'Support team lead', permissions: [PermissionFlagsBits.ManageThreads] },
  { name: '🎧 Senior Support', color: 0x2ca02c, reason: 'Experienced support', permissions: [] },
  { name: '💬 Support Team', color: 0x98fb98, reason: 'Support helpers', permissions: [] },
  { name: '⚖️ Head Moderator', color: 0xffbf00, reason: 'Moderation lead', permissions: [PermissionFlagsBits.ModerateMembers, PermissionFlagsBits.KickMembers, PermissionFlagsBits.BanMembers, PermissionFlagsBits.ManageMessages] },
  { name: '🔨 Moderator', color: 0xff8c00, reason: 'Public channel moderation', permissions: [PermissionFlagsBits.ModerateMembers, PermissionFlagsBits.KickMembers, PermissionFlagsBits.ManageMessages] },
  { name: '🧹 Trial Moderator', color: 0xffff00, reason: 'Trainee moderator', permissions: [PermissionFlagsBits.ManageMessages] },
  { name: '🧪 Bot Tester', color: 0x9370db, reason: 'Feature testing', permissions: [] },
  { name: '🐞 Bug Hunter', color: 0xff69b4, reason: 'Bug reports', permissions: [] },
  { name: '🧱 Contributor', color: 0x008080, reason: 'Documentation and ideas', permissions: [] },
  { name: '🤝 Partner', color: 0x87ceeb, reason: 'Partner representatives', permissions: [] },
  { name: '⭐ VIP', color: 0xffd700, reason: 'Trusted community members', permissions: [] },
  { name: '🙋 Active Helper', color: 0x66c266, reason: 'Helpful community members', permissions: [] },
  { name: '👤 Member', color: 0x999999, reason: 'Verified members', permissions: [] },
  { name: '🔇 Muted', color: 0x666666, reason: 'Restricted members', permissions: [] },
];

const CATEGORIES: { name: string; channels: { name: string; type: ChannelType; topic: string; sendable?: boolean }[] }[] = [
  {
    name: '📌 START HERE', channels: [
      { name: 'welcome', type: ChannelType.GuildText, topic: 'Welcome to PulseKeep Support! First channel you see.', sendable: false },
      { name: 'rules', type: ChannelType.GuildText, topic: 'Server rules and behavior expectations.', sendable: false },
      { name: 'start-here', type: ChannelType.GuildText, topic: 'New here? Start with this onboarding checklist.', sendable: false },
    ],
  },
  {
    name: '📚 PULSEKEEP INFO', channels: [
      { name: 'announcements', type: ChannelType.GuildText, topic: 'Official PulseKeep announcements.', sendable: false },
      { name: 'changelog', type: ChannelType.GuildText, topic: 'Release notes and version history.', sendable: false },
      { name: 'status', type: ChannelType.GuildText, topic: 'Live service health updates.', sendable: false },
      { name: 'commands', type: ChannelType.GuildText, topic: 'Full command reference.', sendable: false },
      { name: 'faq', type: ChannelType.GuildText, topic: 'Frequently asked questions.', sendable: false },
    ],
  },
  {
    name: '🎫 SUPPORT', channels: [
      { name: 'support-info', type: ChannelType.GuildText, topic: 'How PulseKeep support works.' },
      { name: 'help-chat', type: ChannelType.GuildText, topic: 'Public help channel for quick questions.' },
      { name: 'bug-reports', type: ChannelType.GuildText, topic: 'Report bugs you find in PulseKeep.' },
      { name: 'feature-requests', type: ChannelType.GuildText, topic: 'Suggest new features.' },
      { name: 'ticket-panel', type: ChannelType.GuildText, topic: 'Open a support ticket.', sendable: false },
    ],
  },
  {
    name: '🤖 BOT COMMANDS', channels: [
      { name: 'command-menu', type: ChannelType.GuildText, topic: 'Browse and test commands here.' },
      { name: 'bot-status-checks', type: ChannelType.GuildText, topic: 'Quick bot health checks.' },
      { name: 'economy', type: ChannelType.GuildText, topic: 'Economy gameplay and discussion.' },
      { name: 'tickets-demo', type: ChannelType.GuildText, topic: 'Try the ticket system here.' },
    ],
  },
  {
    name: '🧪 TEST LAB', channels: [
      { name: 'slash-command-testing', type: ChannelType.GuildText, topic: 'Safe space to test any slash command.' },
      { name: 'moderation-testing', type: ChannelType.GuildText, topic: 'Test moderation commands safely.' },
      { name: 'economy-testing', type: ChannelType.GuildText, topic: 'Test economy features.' },
      { name: 'ticket-testing', type: ChannelType.GuildText, topic: 'Test ticket creation and closing.' },
      { name: 'dashboard-testing', type: ChannelType.GuildText, topic: 'Test the web dashboard.' },
      { name: 'automod-testing', type: ChannelType.GuildText, topic: 'Test auto-moderation filters.' },
    ],
  },
  {
    name: '💬 COMMUNITY', channels: [
      { name: 'general', type: ChannelType.GuildText, topic: 'General PulseKeep discussion.' },
      { name: 'showcase', type: ChannelType.GuildText, topic: 'Share your PulseKeep setup!' },
      { name: 'suggestions', type: ChannelType.GuildText, topic: 'Quick ideas for PulseKeep.' },
      { name: 'off-topic', type: ChannelType.GuildText, topic: 'Casual chat.' },
    ],
  },
  {
    name: '🧑‍💻 CONTRIBUTORS', channels: [
      { name: 'contributor-chat', type: ChannelType.GuildText, topic: 'Chat for contributors.' },
      { name: 'docs-feedback', type: ChannelType.GuildText, topic: 'Documentation feedback.' },
      { name: 'translation-help', type: ChannelType.GuildText, topic: 'Translation coordination.' },
    ],
  },
  {
    name: '🔒 STAFF', channels: [
      { name: 'staff-chat', type: ChannelType.GuildText, topic: 'Private staff coordination.' },
      { name: 'mod-chat', type: ChannelType.GuildText, topic: 'Moderator discussion.' },
      { name: 'support-notes', type: ChannelType.GuildText, topic: 'Support case notes.' },
      { name: 'review-queue', type: ChannelType.GuildText, topic: 'Items awaiting staff review.' },
      { name: 'staff-commands', type: ChannelType.GuildText, topic: 'Staff-only command testing.' },
    ],
  },
  {
    name: '🧾 LOGS', channels: [
      { name: 'mod-logs', type: ChannelType.GuildText, topic: 'Moderation action logs.' },
      { name: 'bot-logs', type: ChannelType.GuildText, topic: 'Bot event logs.' },
      { name: 'ticket-logs', type: ChannelType.GuildText, topic: 'Ticket close summaries.' },
      { name: 'vote-logs', type: ChannelType.GuildText, topic: 'Vote reward logs.' },
      { name: 'automod-logs', type: ChannelType.GuildText, topic: 'Auto-moderation action logs.' },
      { name: 'deploy-logs', type: ChannelType.GuildText, topic: 'Deployment logs.' },
    ],
  },
  {
    name: '🚨 INCIDENTS', channels: [
      { name: 'incident-response', type: ChannelType.GuildText, topic: 'Incident coordination.' },
      { name: 'security-reports', type: ChannelType.GuildText, topic: 'Security issue reports.' },
      { name: 'audit-review', type: ChannelType.GuildText, topic: 'Audit log review.' },
    ],
  },
  {
    name: '📦 ARCHIVE', channels: [
      { name: 'resolved-tickets', type: ChannelType.GuildText, topic: 'Resolved ticket records.' },
      { name: 'old-announcements', type: ChannelType.GuildText, topic: 'Past announcements.' },
      { name: 'old-bugs', type: ChannelType.GuildText, topic: 'Resolved bug reports.' },
    ],
  },
];

const GUILD_FULL_SEND = [
  PermissionFlagsBits.ViewChannel,
  PermissionFlagsBits.SendMessages,
  PermissionFlagsBits.ReadMessageHistory,
  PermissionFlagsBits.AddReactions,
  PermissionFlagsBits.EmbedLinks,
  PermissionFlagsBits.AttachFiles,
  PermissionFlagsBits.UseApplicationCommands,
];

export const createServerCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('createserver')
    .setDescription('Wipe and rebuild the entire server layout (owner only)')
    .toJSON(),

  async execute({ config }, interaction) {
    if (interaction.user.id !== config.botOwnerID && interaction.user.id !== config.botCoOwnerID) {
      await interaction.reply({ content: 'This command is restricted to the bot owner.', flags: 64 });
      return;
    }

    await interaction.deferReply();

    const guild = interaction.guild;
    if (!guild) {
      await interaction.editReply({ content: 'This can only be used in a server.' });
      return;
    }

    const results: string[] = [];
    const errors: string[] = [];
    const everyone = guild.roles.everyone;

    // Step 0: Delete everything
    results.push('**Phase 1: Deleting existing channels and roles...**');

    const existingChannels = [...guild.channels.cache.values()].filter(c => c.id !== guild.rulesChannelId && c.id !== guild.publicUpdatesChannelId);
    for (const ch of existingChannels) {
      try {
        await ch.delete('PulseKeep server rebuild');
        results.push(`  ✅ Deleted ${ch.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to delete channel #${ch.name}: ${err.message}`);
      }
    }

    const botMember = guild.members.me;
    const botRoleId = botMember?.roles.highest.id;
    const deletableRoles = [...guild.roles.cache.values()]
      .filter(r => r.id !== everyone.id && r.id !== botRoleId && r.id !== guild.roles.premiumSubscriberRole?.id && r.managed === false);
    for (const role of deletableRoles) {
      try {
        await role.delete('PulseKeep server rebuild');
        results.push(`  ✅ Deleted role ${role.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to delete role ${role.name}: ${err.message}`);
      }
    }

    results.push('');
    results.push('**Phase 2: Creating roles...**');

    const createdRoles = new Map<string, string>();
    for (const r of ROLES) {
      try {
        const role = await guild.roles.create({
          name: r.name,
          color: r.color,
          permissions: r.permissions,
          reason: r.reason,
        });
        createdRoles.set(r.name, role.id);
        results.push(`  ✅ Created role ${r.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to create role ${r.name}: ${err.message}`);
      }
    }

    results.push('');
    results.push('**Phase 3: Creating categories and channels...**');

    for (const cat of CATEGORIES) {
      let categoryId: string | null = null;
      try {
        const category = await guild.channels.create({
          name: cat.name,
          type: ChannelType.GuildCategory,
          reason: 'PulseKeep support server layout',
        });
        categoryId = category.id;
        results.push(`  ✅ Created category ${cat.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to create category ${cat.name}: ${err.message}`);
        continue;
      }

      for (const ch of cat.channels) {
        try {
          const overwrites: any[] = [
            { id: everyone.id, deny: [PermissionFlagsBits.SendMessages] },
            { id: interaction.client.user.id, allow: GUILD_FULL_SEND },
          ];

          if (ch.sendable === false) {
            overwrites[0] = { id: everyone.id, deny: [PermissionFlagsBits.SendMessages] };
          }

          await guild.channels.create({
            name: ch.name,
            type: ch.type,
            parent: categoryId,
            topic: ch.topic,
            permissionOverwrites: overwrites,
            reason: 'PulseKeep support server layout',
          });
          results.push(`    ✅ Created #${ch.name}`);
        } catch (err: any) {
          errors.push(`    ❌ Failed to create #${ch.name}: ${err.message}`);
        }
      }
    }

    const summary = [
      `**Server rebuild complete!**`,
      ``,
      ...results,
      ``,
      errors.length > 0 ? `**Errors (${errors.length}):**\n${errors.join('\n')}` : '**No errors.**',
    ].join('\n');

    await interaction.editReply({ content: summary.slice(0, 1900) });
    if (summary.length > 1900) {
      await interaction.followUp({ content: summary.slice(1900, 3900) });
    }
  },
};
