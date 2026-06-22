import {
  SlashCommandBuilder,
  EmbedBuilder,
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ButtonInteraction,
} from 'discord.js';
import type { SlashCommand } from '../types.js';
import { Colors, Ephemeral, timestamp } from '../../utils/embed.js';

function healthEmoji(ms: number): string {
  if (ms < 100) return '🟢 Excellent';
  if (ms < 200) return '🟡 Good';
  if (ms < 400) return '🟠 Okay';
  return '🔴 Slow';
}

export const pingCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('ping')
    .setDescription("Check the bot's latency")
    .toJSON(),

  async execute(ctx, interaction) {
    const sentTs = Date.now();
    const replyTs = interaction.createdTimestamp;
    const roundTrip = sentTs - replyTs;
    const ws = Math.round(interaction.client.ws.ping);

    const emb = new EmbedBuilder()
      .setTitle('🏓 Pong!')
      .setColor(Colors.Utility)
      .addFields(
        { name: 'Round-trip', value: `${healthEmoji(roundTrip)} **${roundTrip}ms**`, inline: true },
        { name: 'WebSocket', value: `${healthEmoji(ws)} **${ws}ms**`, inline: true },
        { name: 'Avg Latency', value: `${healthEmoji(Math.round(ctx.cache.getAvgLatency()))} **${Math.round(ctx.cache.getAvgLatency())}ms**`, inline: true },
      )
      .setFooter({ text: 'Times are measured from gateway → bot → response' });

    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(
      new ButtonBuilder()
        .setCustomId('ping_refresh')
        .setLabel('Refresh')
        .setStyle(ButtonStyle.Secondary)
        .setEmoji('🔄'),
    );

    const reply = await interaction.reply({ embeds: [timestamp(emb)], components: [row], flags: Ephemeral, fetchReply: true });

    // refresh button, only works for whoever ran the command, dies after 30s
    const collector = reply.createMessageComponentCollector({ time: 30_000 });
    collector.on('collect', async (i: ButtonInteraction) => {
      if (i.user.id !== interaction.user.id) {
        await i.reply({ content: 'That button isn\'t yours.', flags: Ephemeral });
        return;
      }
      const rt = Date.now() - i.createdTimestamp;
      const ws2 = Math.round(i.client.ws.ping);
      const emb2 = new EmbedBuilder()
        .setTitle('🏓 Pong!')
        .setColor(Colors.Utility)
        .addFields(
          { name: 'Round-trip', value: `${healthEmoji(rt)} **${rt}ms**`, inline: true },
          { name: 'WebSocket', value: `${healthEmoji(ws2)} **${ws2}ms**`, inline: true },
          { name: 'Avg Latency', value: `${healthEmoji(Math.round(ctx.cache.getAvgLatency()))} **${Math.round(ctx.cache.getAvgLatency())}ms**`, inline: true },
        )
        .setFooter({ text: 'Refreshed' });
      await i.update({ embeds: [timestamp(emb2)], components: [row] });
    });
  },
};
