import { EmbedBuilder } from 'discord.js';

export const Colors = {
  Economy: 0x43B581,
  Moderation: 0xF04747,
  Tickets: 0x5865F2,
  Utility: 0x0099FF,
  Success: 0x57F287,
  Error: 0xED4245,
  Warning: 0xFEE75C,
} as const;

export function footer(emb: EmbedBuilder): EmbedBuilder {
  return emb.setFooter({ text: `PulseKeep v7.0.0` });
}

export function timestamp(emb: EmbedBuilder): EmbedBuilder {
  return emb.setTimestamp(new Date());
}

export function baseEmbed(): EmbedBuilder {
  return timestamp(footer(new EmbedBuilder()));
}
