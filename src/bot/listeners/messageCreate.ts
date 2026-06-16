import { Events, EmbedBuilder, Colors } from 'discord.js';
import type { Bot } from '../client.js';
import { runAutomod } from '../automod/index.js';

export function registerMessageCreateListener(bot: Bot) {
  bot.client.on(Events.MessageCreate, async (message) => {
    if (message.author.bot) return;
    if (!message.guildId || !message.guild) return;

    const member = message.member;
    if (!member) return;

    const isMod = member.permissions.has('ManageMessages') || member.permissions.has('Administrator');
    if (isMod) return;

    const { action, reason } = await runAutomod(
      bot,
      message.guildId,
      message.channelId,
      message.author.id,
      message.content,
      member.roles.cache.map(r => r.id),
    );

    if (!action) return;

    switch (action) {
      case 'delete':
        await message.delete().catch(() => {});
        break;
      case 'warn': {
        const warnEmb = new EmbedBuilder()
          .setTitle('⚠️ Auto-Mod Warning')
          .setDescription(`**Reason:** ${reason}`)
          .setColor(Colors.Warning)
          .setTimestamp();
        await message.channel.send({ content: `<@${message.author.id}>`, embeds: [warnEmb] }).then(m => {
          setTimeout(() => m.delete().catch(() => {}), 5000);
        });
        break;
      }
      case 'timeout': {
        await message.delete().catch(() => {});
        await member.timeout(10 * 60 * 1000, reason).catch(() => {});
        const timeoutEmb = new EmbedBuilder()
          .setTitle('🔇 Auto-Mod Timeout')
          .setDescription(`<@${message.author.id}> has been timed out for **10 minutes**.\n**Reason:** ${reason}`)
          .setColor(Colors.Moderation)
          .setTimestamp();
        await message.channel.send({ embeds: [timeoutEmb] }).then(m => {
          setTimeout(() => m.delete().catch(() => {}), 8000);
        });
        break;
      }
    }
  });
}
