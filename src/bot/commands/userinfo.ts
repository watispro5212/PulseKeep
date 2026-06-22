import { SlashCommandBuilder, EmbedBuilder, GuildMember } from 'discord.js';
import type { Role } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, Ephemeral, timestamp } from '../../utils/embed.js';

function formatDate(d: Date | number | null | undefined): string {
  if (!d) return 'Unknown';
  const date = d instanceof Date ? d : new Date(d);
  return `<t:${Math.floor(date.getTime() / 1000)}:D>`;
}

function relativeDate(d: Date | number | null | undefined): string {
  if (!d) return 'Unknown';
  const date = d instanceof Date ? d : new Date(d);
  return `<t:${Math.floor(date.getTime() / 1000)}:R>`;
}

function statusEmoji(status: string): string {
  switch (status) {
    case 'online': return '🟢';
    case 'idle': return '🟡';
    case 'dnd': return '🔴';
    case 'offline':
    case 'invisible':
    default: return '⚫';
  }
}

export const userinfoCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('userinfo')
    .setDescription('Inspect a user: avatar, join date, roles, key permissions')
    .addUserOption((o) => o.setName('user').setDescription('User to inspect (defaults to you)'))
    .toJSON(),

  async execute({}, interaction) {
    const user = interaction.options.getUser('user') ?? interaction.user;
    const member = (interaction.guild?.members.cache.get(user.id) as GuildMember | undefined)
      ?? (interaction.guild ? await interaction.guild.members.fetch(user.id).catch(() => null) : null);

    if (!interaction.guild) {
      const emb = new EmbedBuilder()
        .setTitle(`${user.username}`)
        .setColor(Colors.Utility)
        .setThumbnail(user.displayAvatarURL({ size: 256 }))
        .addFields(
          { name: 'User ID', value: user.id, inline: true },
          { name: 'Created', value: formatDate(user.createdAt), inline: true },
          { name: 'Bot?', value: user.bot ? 'Yes' : 'No', inline: true },
        )
        .setFooter({ text: 'Some details are only available in servers.' });
      await interaction.reply({ embeds: [timestamp(emb)], flags: Ephemeral });
      return;
    }

    if (!member) {
      await interaction.reply({ content: 'User not found in this server.', flags: Ephemeral });
      return;
    }

    const roles = member.roles.cache
      .filter((r: Role) => r.id !== interaction.guild!.id)
      .sort((a: Role, b: Role) => b.position - a.position)
      .map((r: Role) => `<@&${r.id}>`);

    const topRole = member.roles.highest;
    const keyPerms: string[] = [];
    const p = member.permissions;
    if (p.has('Administrator')) keyPerms.push('Administrator');
    else {
      if (p.has('ManageGuild')) keyPerms.push('Manage Server');
      if (p.has('ManageRoles')) keyPerms.push('Manage Roles');
      if (p.has('ManageChannels')) keyPerms.push('Manage Channels');
      if (p.has('KickMembers')) keyPerms.push('Kick');
      if (p.has('BanMembers')) keyPerms.push('Ban');
      if (p.has('ModerateMembers')) keyPerms.push('Timeout');
      if (p.has('ManageMessages')) keyPerms.push('Manage Messages');
      if (p.has('MentionEveryone')) keyPerms.push('Mention @everyone');
    }

    const presenceStatus = (member.presence?.status ?? 'offline') as string;
    const statusLabel = presenceStatus === 'dnd' ? 'Do Not Disturb' : presenceStatus.charAt(0).toUpperCase() + presenceStatus.slice(1);

    const emb = new EmbedBuilder()
      .setTitle(`${user.bot ? '🤖 ' : ''}${user.displayName || user.username}`)
      .setColor(topRole?.color || Colors.Utility)
      .setThumbnail(user.displayAvatarURL({ size: 256 }))
      .addFields(
        { name: 'Username', value: `${user.tag}`, inline: true },
        { name: 'User ID', value: user.id, inline: true },
        { name: 'Status', value: `${statusEmoji(presenceStatus)} ${statusLabel}`, inline: true },
        { name: 'Account Created', value: `${formatDate(user.createdAt)} (${relativeDate(user.createdAt)})`, inline: true },
        { name: 'Joined Server', value: `${formatDate(member.joinedAt)} (${relativeDate(member.joinedAt)})`, inline: true },
        { name: 'Top Role', value: topRole ? `<@&${topRole.id}>` : 'None', inline: true },
        { name: `Roles [${roles.length}]`, value: roles.length > 0 ? roles.slice(0, 20).join(' ') + (roles.length > 20 ? ` …and ${roles.length - 20} more` : '') : 'None', inline: false },
        { name: 'Key Permissions', value: keyPerms.length > 0 ? keyPerms.join(', ') : 'None', inline: false },
        { name: 'Boosting Since', value: member.premiumSince ? formatDate(member.premiumSince) : 'Not boosting', inline: true },
      )
      .setFooter({ text: `Requested by ${interaction.user.tag}` });

    await interaction.reply({ embeds: [timestamp(emb)], flags: Ephemeral });
  },
};
