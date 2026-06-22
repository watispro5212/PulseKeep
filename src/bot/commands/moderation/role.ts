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

  async execute({ bot }, interaction) {
    if (!interaction.guild) {
      await interaction.reply({ content: '❌ This command must be used in a server.', flags: 64 });
      return;
    }
    const sub = interaction.options.getSubcommand();
    const target = interaction.options.getUser('user', true);
    const role = interaction.options.getRole('role', true);
    const member = await interaction.guild.members.fetch(target.id).catch(() => null);
    const botMember = interaction.guild.members.me;

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
    if (botMember && role.comparePositionTo(botMember.roles.highest) >= 0) {
      await interaction.reply({ content: '❌ That role is higher than my highest role.', flags: 64 });
      return;
    }

    if (sub === 'add' && member.roles.cache.has(role.id)) {
      await interaction.reply({ content: `❌ ${target} already has the ${role} role.`, flags: 64 });
      return;
    }
    if (sub === 'remove' && !member.roles.cache.has(role.id)) {
      await interaction.reply({ content: `❌ ${target} does not have the ${role} role.`, flags: 64 });
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
        const log = new EmbedBuilder()
          .setTitle('Moderation: Role Added')
          .setDescription(`${interaction.user} added ${role} to ${target}`)
          .setColor(Colors.Moderation).setTimestamp();
        await bot.logToChannel(interaction.guildId!, log);
      } else {
        await member.roles.remove(role);
        const emb = new EmbedBuilder()
          .setTitle('Role Removed')
          .setDescription(`Removed ${role} from ${target}.`)
          .setColor(Colors.Moderation);
        await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
        const log = new EmbedBuilder()
          .setTitle('Moderation: Role Removed')
          .setDescription(`${interaction.user} removed ${role} from ${target}`)
          .setColor(Colors.Moderation).setTimestamp();
        await bot.logToChannel(interaction.guildId!, log);
      }
    } catch (err) {
      await interaction.reply({ content: `❌ Failed to ${sub} role: ${err instanceof Error ? err.message : err}`, flags: 64 });
    }
  },
};
