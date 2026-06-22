import { EmbedBuilder } from 'discord.js';
import type { Bot } from '../client.js';
import { eq } from 'drizzle-orm';

const SPAM_WINDOW_MS = 5000;
const SPAM_LIMIT = 5;
const MASS_MENTION_LIMIT = 5;
const CAPS_RATIO = 0.7;
const CAPS_MIN_LENGTH = 10;

const URL_PATTERN = /https?:\/\/[^\s]+/gi;

interface SpamTracker {
  count: number;
  firstMsg: number;
  lastMsg: number;
}

const spamTracking = new Map<string, SpamTracker>();
let spamCleanupInterval: ReturnType<typeof setInterval> | null = null;

export function startSpamCleanup() {
  if (spamCleanupInterval) return;
  spamCleanupInterval = setInterval(() => {
    const now = Date.now();
    for (const [key, val] of spamTracking) {
      if (now - val.lastMsg > 30000) spamTracking.delete(key);
    }
  }, 60000);
  if (spamCleanupInterval.unref) spamCleanupInterval.unref();
}

export function stopSpamCleanup() {
  if (spamCleanupInterval) {
    clearInterval(spamCleanupInterval);
    spamCleanupInterval = null;
  }
}

function checkSpam(guildId: string, userId: string): boolean {
  const key = `${guildId}:${userId}`;
  const now = Date.now();
  const cur = spamTracking.get(key);

  if (!cur || now - cur.firstMsg > SPAM_WINDOW_MS) {
    spamTracking.set(key, { count: 1, firstMsg: now, lastMsg: now });
    return false;
  }

  cur.count++;
  cur.lastMsg = now;
  return cur.count > SPAM_LIMIT;
}

function countMentions(content: string): number {
  const userMentions = (content.match(/<@!?\d{17,19}>/g) || []).length;
  const roleMentions = (content.match(/<@&\d{17,19}>/g) || []).length;
  const everyoneMention = content.includes('@everyone') || content.includes('@here') ? 3 : 0;
  return userMentions + roleMentions + everyoneMention;
}

function hasExcessiveCaps(content: string): boolean {
  const letters = content.replace(/[^a-zA-Z]/g, '');
  if (letters.length < CAPS_MIN_LENGTH) return false;
  const upper = letters.replace(/[^A-Z]/g, '').length;
  return upper / letters.length >= CAPS_RATIO;
}

function hasBannedWord(content: string, bannedWords: string[]): boolean {
  const lower = content.toLowerCase();
  return bannedWords.some(w => lower.includes(w));
}

function hasLink(content: string): boolean {
  return URL_PATTERN.test(content);
}

export async function runAutomod(
  bot: Bot,
  guildId: string,
  channelId: string,
  userId: string,
  content: string,
  memberRoles: string[],
): Promise<{ action: string | null; reason: string | null }> {
  if (!bot.db || !guildId) return { action: null, reason: null };

  const toggles = await bot.getGuildToggles(guildId);
  if (toggles.automodEnabled === false) return { action: null, reason: null };

  try {
    const { guildConfigs } = await import('../../db/schema.js');
    const rows: any[] = await bot.db
      .select()
      .from(guildConfigs)
      .where(eq(guildConfigs.guildId, guildId))
      .limit(1);

    const cfg = rows[0];
    if (!cfg) return { action: null, reason: null };

    const bannedWords = (cfg.automodBannedWords || '').split(',').map((w: string) => w.trim().toLowerCase()).filter(Boolean);

    const checks: { enabled: boolean; check: () => boolean; action: string; reason: string }[] = [
      { enabled: cfg.automodSpamEnabled !== false, check: () => checkSpam(guildId, userId), action: 'warn', reason: 'Spam detected (excessive messages)' },
      { enabled: cfg.automodMentionEnabled !== false, check: () => countMentions(content) > MASS_MENTION_LIMIT, action: 'timeout', reason: 'Mass mention detected' },
      { enabled: cfg.automodCapsEnabled !== false, check: () => hasExcessiveCaps(content), action: 'delete', reason: 'Excessive caps lock' },
      { enabled: cfg.automodLinkEnabled !== false, check: () => hasLink(content), action: 'delete', reason: 'Link blocked by automod' },
      { enabled: cfg.automodWordsEnabled !== false && bannedWords.length > 0, check: () => hasBannedWord(content, bannedWords), action: 'delete', reason: 'Banned word detected' },
    ];

    for (const check of checks) {
      if (!check.enabled) continue;
      if (check.check()) {
        const logEmbed = new EmbedBuilder()
          .setTitle('Auto-Mod Action')
          .setDescription(`**User:** <@${userId}> (\`${userId}\`)\n**Action:** ${check.action}\n**Reason:** ${check.reason}\n**Message:** ${content.slice(0, 500)}`)
          .setColor(0xF04747)
          .setTimestamp();
        await bot.logToChannel(guildId, logEmbed);

        return { action: check.action, reason: check.reason };
      }
    }
  } catch (err) {
    console.error(`[Automod] Error in guild ${guildId}:`, err);
  }

  return { action: null, reason: null };
}
