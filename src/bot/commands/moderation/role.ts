import {
  SlashCommandBuilder,
  EmbedBuilder,
  PermissionFlagsBits,
} from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';

export const roleCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('role')
    .setDescription('Manage a member\'s roles')
    .addSubcommand((s) =>
      s.setName('add').setDescription('Add a role to a member')
        .addUserOption((o) => o.setName('user').setDescription('User').setRequired(true))
        .addRoleOption((o) => o.setName('role').setDescription('Role to add').setRequired(true)),
    )
    .addSubcommand((s) =>
      s.setName('remove').setDescription('Remove a role from a member')
        .addUserOption((o) => o.setName('user').setDescription('User').setRequired(true))
        .addRoleOption((o) => o.setName('role').setDescription('Role to remove').setRequired(true)),
    )
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageRoles)
    .toJSON(),

  async execute(_ctx, interaction) {
    const sub = interaction.options.getSubcommand();
    const target = interaction.options.getUser('user', true);
    const role = interaction.options.getRole('role', true);
    const member = interaction.guild?.members.cache.get(target.id);

    if (!member) {
      await interaction.reply({ content: '❌ Could not find that member.', flags: 64 });
      return;
    }
    if (!member.moderatable) {
      await interaction.reply({ content: '❌ I cannot manage roles for that member.', flags: 64 });
      return;
    }
    if (role.managed) {
      await interaction.reply({ content: '❌ That role is managed by an integration and cannot be manually assigned.', flags: 64 });
      return;
    }
    if (role.comparePositionTo(interaction.guild?.members.me?.roles.highest ?? role) >= 0) {
      await interaction.reply({ content: '❌ That role is higher than my highest role.', flags: 64 });
      return;
    }

    try {
      if (sub === 'add') {
        await member.roles.add(role);
        const emb = new EmbedBuilder()
          .setTitle('Role Added')
          .setDescription(`Added ${role} to ${target}.`)
          .setColor(Colors.Moderation);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      } else {
        await member.roles.remove(role);
        const emb = new EmbedBuilder()
          .setTitle('Role Removed')
          .setDescription(`Removed ${role} from ${target}.`)
          .setColor(Colors.Moderation);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
      }
    } catch {
      await interaction.reply({ content: `❌ Failed to ${sub} role.`, flags: 64 });
    }
  },
};
