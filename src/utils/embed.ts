import { EmbedBuilder } from 'discord.js';

export const Colors = {
  Economy: 0x43B581,
  Moderation: 0xF04747,
  Tickets: 0x5865F2,
  Utility: 0x0099FF,
  Success: 0x57F287,
  Error: 0xED4245,
  Warning: 0xFEE75C,
  Configure: 0x99AAFF,
} as const;

export function footer(emb: EmbedBuilder): EmbedBuilder {
  return emb.setFooter({ text: `PulseKeep v7.4.0` });
}

export function timestamp(emb: EmbedBuilder): EmbedBuilder {
  return emb.setTimestamp(new Date());
}

export function baseEmbed(): EmbedBuilder {
  return timestamp(footer(new EmbedBuilder()));
}

export const Ephemeral = 64;

export function formatNumber(n: number): string {
  return n.toLocaleString('en-US');
}

export function formatCompact(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M';
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K';
  return String(n);
}

export function formatCooldown(seconds: number): string {
  if (seconds >= 86400) return `${Math.round(seconds / 86400)}d`;
  if (seconds >= 3600) return `${Math.round(seconds / 3600)}h`;
  if (seconds >= 60) return `${Math.floor(seconds / 60)}m ${Math.round(seconds % 60)}s`;
  return `${Math.round(seconds)}s`;
}
