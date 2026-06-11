import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const roleinfoCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('roleinfo')
    .setDescription('Get information about a role')
    .addRoleOption((o) => o.setName('role').setDescription('Role to inspect').setRequired(true))
    .toJSON(),

  async execute(_ctx, interaction) {
    const role = interaction.options.getRole('role', true);

    const perms: string[] = [];
    if (role.permissions.has('Administrator')) {
      perms.push('Administrator');
    } else {
      if (role.permissions.has('ManageGuild')) perms.push('Manage Server');
      if (role.permissions.has('ManageRoles')) perms.push('Manage Roles');
      if (role.permissions.has('ManageChannels')) perms.push('Manage Channels');
      if (role.permissions.has('KickMembers')) perms.push('Kick Members');
      if (role.permissions.has('BanMembers')) perms.push('Ban Members');
      if (role.permissions.has('ModerateMembers')) perms.push('Timeout Members');
      if (role.permissions.has('MentionEveryone')) perms.push('Mention Everyone');
      if (role.permissions.has('ManageMessages')) perms.push('Manage Messages');
    }

    const colorHex = role.hexColor || '#000000';
    const created = new Date(role.createdTimestamp).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
    });

    const emb = new EmbedBuilder()
      .setTitle(`Role — ${role.name}`)
      .setColor(role.color || Colors.Utility)
      .addFields(
        { name: 'Role ID', value: role.id, inline: true },
        { name: 'Color', value: colorHex, inline: true },
        { name: 'Position', value: `${role.position}`, inline: true },
        { name: 'Mentionable', value: role.mentionable ? 'Yes' : 'No', inline: true },
        { name: 'Displayed Separately', value: role.hoist ? 'Yes' : 'No', inline: true },
        { name: 'Created', value: created, inline: true },
        { name: 'Key Permissions', value: perms.length > 0 ? perms.join(', ') : 'None', inline: false },
      );

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
