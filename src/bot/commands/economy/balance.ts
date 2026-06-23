import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { STARTING_BALANCE, hasXpBoost, COOLDOWNS } from '../../economy/store.js';

function formatRemaining(ms: number): string {
  if (ms <= 0) return '✅ Ready';
  if (ms < 0) return '✅ Ready';
  const secs = Math.ceil(ms / 1000);
  if (secs < 60) return `⏳ ${secs}s`;
  const mins = Math.ceil(secs / 60);
  if (mins < 60) return `⏳ ${mins}m`;
  const hrs = Math.floor(mins / 60);
  return `⏳ ${hrs}h ${mins % 60}m`;
}

export const balanceCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('balance')
    .setDescription('Check your or another user\'s balance')
    .addUserOption((o) => o.setName('user').setDescription('User to check'))
    .addBooleanOption((o) => o.setName('public').setDescription('Show publicly'))
    .toJSON(),

  async execute({ db }, interaction) {
    const publicReply = !!interaction.options.getBoolean('public');
    const target = interaction.options.getUser('user') ?? interaction.user;
    const isSelf = target.id === interaction.user.id;
    let balance = STARTING_BALANCE;
    const extras: string[] = [];

    if (db) {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, target.id))
        .limit(1);
      if (rows.length > 0 && rows[0]) {
        const rec = rows[0];
        balance = rec.balance;
        if (hasXpBoost(rec)) {
          const remaining = Math.ceil((new Date(rec.xpBoostExpiry).getTime() - Date.now()) / 60000);
          extras.push(`⚡ XP Boost active — **${remaining}m** remaining`);
        }
        if ((rec.luckyCloverActive ?? 0) > 0) {
          extras.push(`🍀 ${rec.luckyCloverActive} Lucky Clover(s) ready`);
        }
        if (isSelf) {
          const now = Date.now();
          const cdChecks: [string, number, Date | null | undefined][] = [
            ['Daily', COOLDOWNS.daily, rec.lastDailyClaim],
            ['Weekly', COOLDOWNS.weekly, rec.lastWeeklyClaim],
            ['Work', COOLDOWNS.work, rec.lastWork],
          ];
          for (const [name, cooldown, last] of cdChecks) {
            if (last) {
              const elapsed = now - new Date(last).getTime();
              const remaining = cooldown - elapsed;
              if (remaining > 0) {
                extras.push(`${formatRemaining(remaining)} until next **/${name.toLowerCase()}**`);
              }
            }
          }
        }
      }
    }

    const emb = new EmbedBuilder()
      .setTitle(`${target.username}'s Balance`)
      .setDescription(`💰 **${balance.toLocaleString()}** Pulses${extras.length > 0 ? '\n' + extras.join('\n') : ''}`)
      .setColor(Colors.Economy);

    if (publicReply) {
      await interaction.reply({ embeds: [footer(timestamp(emb))] });
    } else {
      await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
    }
  },
};
