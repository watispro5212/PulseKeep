import { SlashCommandBuilder, EmbedBuilder } from 'discord.js';
import type { SlashCommand } from '../../types.js';
import { Colors, footer, timestamp } from '../../../utils/embed.js';
import { userEconomy } from '../../../db/schema.js';
import { eq } from 'drizzle-orm';
import { COOLDOWNS } from '../../economy/store.js';

const EXTRA_COOLDOWNS = {
  search: 15 * 60 * 1000,
  vote: 12 * 60 * 60 * 1000,
} as const;

function formatRemaining(ms: number): string {
  if (ms <= 0) return '**✅ Ready**';
  const secs = Math.ceil(ms / 1000);
  if (secs < 60) return `⏳ **${secs}s**`;
  const mins = Math.ceil(secs / 60);
  if (mins < 60) return `⏳ **${mins}m**`;
  const hrs = Math.floor(mins / 60);
  return `⏳ **${hrs}h ${mins % 60}m**`;
}

export const cooldownsCommand: SlashCommand = {
  data: new SlashCommandBuilder()
    .setName('cooldowns')
    .setDescription('Check your economy cooldown timers')
    .toJSON(),

  async execute({ db }, interaction) {
    if (!db) {
      await interaction.reply({ content: 'Database unavailable.', flags: 64 });
      return;
    }

    const userId = interaction.user.id;
    const now = Date.now();

    let rec: any;
    try {
      const rows = await db
        .select()
        .from(userEconomy)
        .where(eq(userEconomy.userId, userId))
        .limit(1);
      rec = rows[0];
    } catch {
      await interaction.reply({ content: '❌ Database error. Please try again.', flags: 64 });
      return;
    }

    const fields: { name: string; value: string; inline: boolean }[] = [];

    const checks: [string, number, Date | null | undefined][] = [
      ['Daily', COOLDOWNS.daily, rec?.lastDailyClaim],
      ['Weekly', COOLDOWNS.weekly, rec?.lastWeeklyClaim],
      ['Work', COOLDOWNS.work, rec?.lastWork],
      ['Rob', COOLDOWNS.rob, rec?.lastRob],
      ['Fish', COOLDOWNS.fish, rec?.lastFish],
      ['Mine', COOLDOWNS.mine, rec?.lastMine],
      ['Search', EXTRA_COOLDOWNS.search, rec?.lastSearch],
      ['Vote', EXTRA_COOLDOWNS.vote, rec?.lastVote],
    ];

    for (const [name, cooldown, lastClaim] of checks) {
      if (!lastClaim) {
        fields.push({ name, value: '**✅ Ready**', inline: true });
        continue;
      }
      const elapsed = now - new Date(lastClaim).getTime();
      fields.push({ name, value: formatRemaining(cooldown - elapsed), inline: true });
    }

    const emb = new EmbedBuilder()
      .setTitle('⏱️ Your Cooldowns')
      .setDescription('All your economy timers at a glance')
      .addFields(fields)
      .setColor(Colors.Economy);

    await interaction.reply({ embeds: [footer(timestamp(emb))], flags: 64 });
  },
};
