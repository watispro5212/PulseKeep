import { SlashCommandBuilder, PermissionFlagsBits, ChannelType } from 'discord.js';
import type { SlashCommand } from '../types.js';

const ROLES = [
  { name: '👑 Owner', color: 0x7c5cfc, permissions: [PermissionFlagsBits.Administrator] },
  { name: '💎 Co-Owner', color: 0x4363d8, permissions: [PermissionFlagsBits.Administrator] },
  { name: '🛡️ Admin', color: 0xe6194b, permissions: [PermissionFlagsBits.Administrator] },
  { name: '🔨 Moderator', color: 0xff8c00, permissions: [PermissionFlagsBits.ModerateMembers, PermissionFlagsBits.KickMembers, PermissionFlagsBits.BanMembers, PermissionFlagsBits.ManageMessages] },
  { name: '🧑‍💻 Developer', color: 0x4b0082, permissions: [PermissionFlagsBits.ViewAuditLog] },
  { name: '🤝 Support', color: 0x2ca02c, permissions: [PermissionFlagsBits.ManageThreads] },
  { name: '⭐ VIP', color: 0xffd700, permissions: [] },
  { name: '👤 Member', color: 0x999999, permissions: [] },
  { name: '🔇 Muted', color: 0x666666, permissions: [] },
];

const STAFF_ROLE_NAMES = ['👑 Owner', '💎 Co-Owner', '🛡️ Admin', '🔨 Moderator', '🤝 Support'];

interface CatConfig {
  name: string;
  staffOnly?: boolean;
  channels: { name: string; topic: string; sendable?: boolean }[];
}

const CATEGORIES: CatConfig[] = [
  {
    name: '📌 INFORMATION', channels: [
      { name: 'welcome', topic: 'Welcome to PulseKeep Support!', sendable: false },
      { name: 'rules', topic: 'Server rules and expectations.', sendable: false },
      { name: 'announcements', topic: 'Official PulseKeep announcements.', sendable: false },
    ],
  },
  {
    name: '💬 COMMUNITY', channels: [
      { name: 'general', topic: 'General PulseKeep discussion.' },
      { name: 'support', topic: 'Get help with PulseKeep.' },
      { name: 'showcase', topic: 'Share your PulseKeep setup!' },
      { name: 'off-topic', topic: 'Casual chat.' },
    ],
  },
  {
    name: '🤖 PULSEKEEP', channels: [
      { name: 'bot-commands', topic: 'Use PulseKeep commands here.' },
      { name: 'status', topic: 'Bot service status updates.', sendable: false },
      { name: 'testing', topic: 'Test commands safely.' },
    ],
  },
  {
    name: '🔒 STAFF', staffOnly: true, channels: [
      { name: 'staff-chat', topic: 'Private staff coordination.' },
      { name: 'mod-logs', topic: 'Moderation and bot action logs.', sendable: false },
      { name: 'admin', topic: 'Server administration.' },
    ],
  },
  {
    name: '🎫 TICKETS', channels: [
      { name: 'open-ticket', topic: 'Open a support ticket using the button below.', sendable: false },
      { name: 'ticket-logs', topic: 'Closed ticket transcripts.', sendable: false },
    ],
  },
];

const GUILD_FULL_SEND: any = [
  PermissionFlagsBits.ViewChannel,
  PermissionFlagsBits.SendMessages,
  PermissionFlagsBits.ReadMessageHistory,
  PermissionFlagsBits.AddReactions,
  PermissionFlagsBits.EmbedLinks,
  PermissionFlagsBits.AttachFiles,
  PermissionFlagsBits.UseApplicationCommands,
];

const STAFF_CHANNEL_VIEW: any = [
  PermissionFlagsBits.ViewChannel,
  PermissionFlagsBits.SendMessages,
  PermissionFlagsBits.ReadMessageHistory,
];

export const createServerCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('createserver')
    .setDescription('Wipe and rebuild the entire server layout (owner only)')
    .toJSON(),

  async execute({ config }, interaction) {
    if (interaction.user.id !== config.botOwnerID && interaction.user.id !== config.botCoOwnerID) {
      await interaction.reply({ content: `This command is restricted to the bot owner. Your ID: \`${interaction.user.id}\` — expected: \`${config.botOwnerID}\``, flags: 64 });
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
    const botMember = guild.members.me;
    const botRoleId = botMember?.roles.highest.id;

    results.push('**Phase 1: Deleting existing roles...**');
    const deletableRoles = [...guild.roles.cache.values()]
      .filter(r => r.id !== everyone.id && r.id !== botRoleId && r.id !== guild.roles.premiumSubscriberRole?.id && r.managed === false);
    for (const role of deletableRoles) {
      try {
        await role.delete('PulseKeep server rebuild — roles first');
        results.push(`  ✅ Deleted role ${role.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to delete role ${role.name}: ${err.message}`);
      }
    }

    results.push('');
    results.push('**Phase 2: Deleting existing channels...**');
    const existingChannels = [...guild.channels.cache.values()];
    for (const ch of existingChannels) {
      try {
        await ch.delete('PulseKeep server rebuild');
        results.push(`  ✅ Deleted #${ch.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to delete #${ch.name}: ${err.message}`);
      }
    }

    results.push('');
    results.push('**Phase 3: Creating roles...**');
    const createdRoles = new Map<string, string>();
    for (const r of ROLES) {
      try {
        const role = await guild.roles.create({
          name: r.name,
          color: r.color,
          permissions: r.permissions,
          reason: 'PulseKeep support server layout',
        });
        createdRoles.set(r.name, role.id);
        results.push(`  ✅ Created role ${r.name}`);
      } catch (err: any) {
        errors.push(`  ❌ Failed to create role ${r.name}: ${err.message}`);
      }
    }

    results.push('');
    results.push('**Phase 4: Creating categories and channels...**');
    for (const cat of CATEGORIES) {
      let categoryId: string | null = null;
      try {
        const catOverwrites: any[] = [{ id: everyone.id, deny: [PermissionFlagsBits.SendMessages] }];
        if (cat.staffOnly) {
          catOverwrites[0] = { id: everyone.id, deny: [PermissionFlagsBits.ViewChannel] };
          for (const roleName of STAFF_ROLE_NAMES) {
            const roleId = createdRoles.get(roleName);
            if (roleId) catOverwrites.push({ id: roleId, allow: STAFF_CHANNEL_VIEW });
          }
        }
        catOverwrites.push({ id: interaction.client.user.id, allow: GUILD_FULL_SEND });

        const category = await guild.channels.create({
          name: cat.name,
          type: ChannelType.GuildCategory,
          permissionOverwrites: catOverwrites,
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
            { id: interaction.client.user.id, allow: GUILD_FULL_SEND },
          ];
          if (cat.staffOnly) {
            for (const roleName of STAFF_ROLE_NAMES) {
              const roleId = createdRoles.get(roleName);
              if (roleId) overwrites.push({ id: roleId, allow: STAFF_CHANNEL_VIEW });
            }
          } else if (ch.sendable === false) {
            overwrites.push({ id: everyone.id, deny: [PermissionFlagsBits.SendMessages] });
          } else {
            overwrites.push({ id: everyone.id, allow: [PermissionFlagsBits.SendMessages] });
          }

          await guild.channels.create({
            name: ch.name,
            type: ChannelType.GuildText,
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
