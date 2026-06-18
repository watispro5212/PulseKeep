import { SlashCommandBuilder, EmbedBuilder, ChannelType, GuildVerificationLevel } from 'discord.js';
import type { GuildMember, GuildBasedChannel } from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, Ephemeral, timestamp } from '../../utils/embed.js';

const verificationNames: Record<number, string> = {
  [GuildVerificationLevel.None]: 'None',
  [GuildVerificationLevel.Low]: 'Low — email verified',
  [GuildVerificationLevel.Medium]: 'Medium — registered ≥5 min',
  [GuildVerificationLevel.High]: 'High — member ≥10 min',
  [GuildVerificationLevel.VeryHigh]: 'Highest — verified phone',
};

function formatDate(d: Date | null | undefined): string {
  if (!d) return 'Unknown';
  return `<t:${Math.floor(d.getTime() / 1000)}:D>`;
}

function relativeDate(d: Date | null | undefined): string {
  if (!d) return 'Unknown';
  return `<t:${Math.floor(d.getTime() / 1000)}:R>`;
}

export const serverinfoCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('serverinfo')
    .setDescription('Server details: members, channels, roles, boosts, verification')
    .toJSON(),

  async execute(_ctx, interaction) {
    const guild = interaction.guild;
    if (!guild) {
      await interaction.reply({ content: 'This command can only be used in a server.', flags: Ephemeral });
      return;
    }

    // try to get real member counts — huge servers get the cached version
    let memberCount = guild.memberCount;
    let onlineCount: number | null = null;
    try {
      await guild.members.fetch({ withPresences: true }).catch(() => null);
      memberCount = guild.members.cache.size || memberCount;
      const online = guild.members.cache.filter((m: GuildMember) => m.presence?.status === 'online').length;
      onlineCount = online;
    } catch { /* fall back */ }

    const channels = await guild.channels.fetch().catch(() => null);
    const textCount = channels?.filter((c: GuildBasedChannel | null) => c?.type === ChannelType.GuildText || c?.type === ChannelType.GuildAnnouncement).size ?? 0;
    const voiceCount = channels?.filter((c: GuildBasedChannel | null) => c?.type === ChannelType.GuildVoice || c?.type === ChannelType.GuildStageVoice).size ?? 0;
    const categoryCount = channels?.filter((c: GuildBasedChannel | null) => c?.type === ChannelType.GuildCategory).size ?? 0;
    const forumCount = channels?.filter((c: GuildBasedChannel | null) => c?.type === ChannelType.GuildForum).size ?? 0;
    const totalChannels = channels?.size ?? 0;

    const roleCount = guild.roles.cache.size - 1; // exclude @everyone
    const emojiCount = guild.emojis.cache.size;
    const stickerCount = guild.stickers.cache.size;

    const owner = await guild.fetchOwner().catch(() => null);

    const emb = new EmbedBuilder()
      .setTitle(`${guild.name}`)
      .setColor(Colors.Utility)
      .setThumbnail(guild.iconURL({ size: 256 }) ?? null)
      .setImage(guild.bannerURL({ size: 512 }) ?? null)
      .addFields(
        { name: 'Server ID', value: guild.id, inline: true },
        { name: 'Owner', value: owner ? `<@${owner.id}>` : 'Unknown', inline: true },
        { name: 'Created', value: `${formatDate(guild.createdAt)} (${relativeDate(guild.createdAt)})`, inline: true },
        { name: 'Members', value: `${memberCount.toLocaleString()}${onlineCount !== null ? `\n(${onlineCount.toLocaleString()} online)` : ''}`, inline: true },
        { name: 'Channels', value: `${totalChannels}\n📝 ${textCount} text • 🔊 ${voiceCount} voice\n📁 ${categoryCount} cat • 💬 ${forumCount} forum`, inline: true },
        { name: 'Roles', value: `${roleCount}`, inline: true },
        { name: 'Emojis', value: `${emojiCount}`, inline: true },
        { name: 'Stickers', value: `${stickerCount}`, inline: true },
        { name: 'Boosts', value: guild.premiumSubscriptionCount ? `✨ ${guild.premiumSubscriptionCount} (Tier ${guild.premiumTier})` : 'None', inline: true },
        { name: 'Verification', value: verificationNames[guild.verificationLevel] ?? 'Unknown', inline: true },
        { name: 'MFA / 2FA', value: guild.mfaLevel ? 'Required for mods' : 'Not required', inline: true },
        { name: 'NSFW Filter', value: guild.nsfwLevel === 0 ? 'Default' : guild.nsfwLevel === 1 ? 'Less restrictive' : 'More restrictive', inline: true },
        { name: 'Preferred Locale', value: guild.preferredLocale ?? 'en-US', inline: true },
      )
      .setFooter({ text: `Requested by ${interaction.user.tag}` });

    if (guild.description) emb.setDescription(guild.description);

    await interaction.reply({ embeds: [timestamp(emb)], flags: Ephemeral });
  },
};
