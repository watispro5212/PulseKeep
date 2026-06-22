import { EmbedBuilder } from 'discord.js';
import type { Bot } from '../client.js';
import { eq } from 'drizzle-orm';

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

function checkSpam(guildId: string, userId: string, limit: number, windowMs: number): boolean {
  const key = `${guildId}:${userId}`;
  const now = Date.now();
  const cur = spamTracking.get(key);

  if (!cur || now - cur.firstMsg > windowMs) {
    spamTracking.set(key, { count: 1, firstMsg: now, lastMsg: now });
    return false;
  }

  cur.count++;
  cur.lastMsg = now;
  return cur.count > limit;
}

function countMentions(content: string): number {
  const userMentions = (content.match(/<@!?\d{17,19}>/g) || []).length;
  const roleMentions = (content.match(/<@&\d{17,19}>/g) || []).length;
  const everyoneMention = content.includes('@everyone') || content.includes('@here') ? 3 : 0;
  return userMentions + roleMentions + everyoneMention;
}

function hasExcessiveCaps(content: string, ratio: number, minLength: number): boolean {
  const letters = content.replace(/[^a-zA-Z]/g, '');
  if (letters.length < minLength) return false;
  const upper = letters.replace(/[^A-Z]/g, '').length;
  return upper / letters.length >= ratio;
}

function hasBannedWord(content: string, bannedWords: string[]): boolean {
  const lower = content.toLowerCase();
  return bannedWords.some(w => lower.includes(w));
}

function hasLink(content: string, allowedDomains: string[]): boolean {
  const matches = content.match(URL_PATTERN);
  if (!matches) return false;
  if (allowedDomains.length === 0) return true;
  return matches.some(url => {
    try {
      const host = new URL(url).hostname.replace(/^www\./, '');
      return !allowedDomains.some(d => host === d || host.endsWith('.' + d));
    } catch {
      return true;
    }
  });
}

function isExempt(memberRoles: string[], exemptRoleIds: string[]): boolean {
  if (exemptRoleIds.length === 0) return false;
  return memberRoles.some(r => exemptRoleIds.includes(r));
}

export async function runAutomod(
  bot: Bot,
  guildId: string,
  channelId: string,
  userId: string,
  content: string,
  memberRoles: string[],
): Promise<{ action: string | null; reason: string | null; timeoutMinutes?: number }> {
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

    const spamLimit = cfg.automodSpamLimit ?? 5;
    const spamWindow = cfg.automodSpamWindow ?? 5000;
    const mentionLimit = cfg.automodMentionLimit ?? 5;
    const capsRatio = (cfg.automodCapsRatio ?? 70) / 100;
    const capsMinLength = cfg.automodCapsMinLength ?? 10;
    const bannedWords = (cfg.automodBannedWords || '').split(',').map((w: string) => w.trim().toLowerCase()).filter(Boolean);
    const allowedDomains = (cfg.automodAllowedDomains || '').split(',').map((d: string) => d.trim().toLowerCase().replace(/^www\./, '')).filter(Boolean);
    const exemptRoles = (cfg.automodExemptRoles || '').split(',').map((r: string) => r.trim()).filter(Boolean);
    const spamAction = cfg.automodSpamAction || 'warn';
    const mentionAction = cfg.automodMentionAction || 'timeout';
    const capsAction = cfg.automodCapsAction || 'delete';
    const linkAction = cfg.automodLinkAction || 'delete';
    const wordsAction = cfg.automodWordsAction || 'delete';
    const timeoutDuration = cfg.automodTimeoutDuration ?? 10;

    if (exemptRoles.length > 0 && isExempt(memberRoles, exemptRoles)) {
      return { action: null, reason: null };
    }

    const checks: { enabled: boolean; check: () => boolean; action: string; reason: string }[] = [
      { enabled: cfg.automodSpamEnabled !== false, check: () => checkSpam(guildId, userId, spamLimit, spamWindow), action: spamAction, reason: `Spam detected (${spamLimit} msgs in ${spamWindow / 1000}s)` },
      { enabled: cfg.automodMentionEnabled !== false, check: () => countMentions(content) > mentionLimit, action: mentionAction, reason: `Mass mention detected (${countMentions(content)} mentions, limit ${mentionLimit})` },
      { enabled: cfg.automodCapsEnabled !== false, check: () => hasExcessiveCaps(content, capsRatio, capsMinLength), action: capsAction, reason: 'Excessive caps lock' },
      { enabled: cfg.automodLinkEnabled !== false, check: () => hasLink(content, allowedDomains), action: linkAction, reason: allowedDomains.length > 0 ? 'Link not in allowed domains' : 'Link blocked by automod' },
      { enabled: cfg.automodWordsEnabled !== false && bannedWords.length > 0, check: () => hasBannedWord(content, bannedWords), action: wordsAction, reason: 'Banned word detected' },
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

        return { action: check.action, reason: check.reason, timeoutMinutes: check.action === 'timeout' ? timeoutDuration : undefined };
      }
    }
  } catch (err) {
    console.error(`[Automod] Error in guild ${guildId}:`, err);
  }

  return { action: null, reason: null };
}
