import { pgTable, varchar, boolean, timestamp, integer, serial, text } from 'drizzle-orm/pg-core';

export const guildConfigs = pgTable('guild_configs', {
  guildId: varchar('guild_id', { length: 20 }).primaryKey(),
  prefix: varchar('prefix', { length: 10 }).default('!').notNull(),
  logChannelId: varchar('log_channel_id', { length: 20 }),
  premium: boolean('premium').default(false).notNull(),
  economyEnabled: boolean('economy_enabled').default(true).notNull(),
  ticketsEnabled: boolean('tickets_enabled').default(true).notNull(),
  modlogsEnabled: boolean('modlogs_enabled').default(true).notNull(),
  welcomeEnabled: boolean('welcome_enabled').default(false).notNull(),
  welcomeChannelId: varchar('welcome_channel_id', { length: 20 }),
  ticketCategoryId: varchar('ticket_category_id', { length: 20 }),
  voteChannelId: varchar('vote_channel_id', { length: 20 }),
  automodEnabled: boolean('automod_enabled').default(true).notNull(),
  automodSpamEnabled: boolean('automod_spam_enabled').default(true).notNull(),
  automodMentionEnabled: boolean('automod_mention_enabled').default(true).notNull(),
  automodLinkEnabled: boolean('automod_link_enabled').default(true).notNull(),
  automodCapsEnabled: boolean('automod_caps_enabled').default(true).notNull(),
  automodWordsEnabled: boolean('automod_words_enabled').default(true).notNull(),
  automodBannedWords: text('automod_banned_words'),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});

export const userEconomy = pgTable('user_economy', {
  userId: varchar('user_id', { length: 20 }).primaryKey(),
  balance: integer('balance').default(0).notNull(),
  lastDailyClaim: timestamp('last_daily_claim'),
  lastWeeklyClaim: timestamp('last_weekly_claim'),
  lastWork: timestamp('last_work'),
  lastRob: timestamp('last_rob'),
  lastFish: timestamp('last_fish'),
  lastMine: timestamp('last_mine'),
  totalEarned: integer('total_earned').default(0).notNull(),
  totalGambled: integer('total_gambled').default(0).notNull(),
  transactions: integer('transactions').default(0).notNull(),
  streak: integer('streak').default(0).notNull(),
  lastVote: timestamp('last_vote'),
  xpBoostExpiry: timestamp('xp_boost_expiry'),
  luckyCloverActive: integer('lucky_clover_active').default(0).notNull(),
});

export const userInventory = pgTable('user_inventory', {
  id: serial('id').primaryKey(),
  userId: varchar('user_id', { length: 20 }).notNull(),
  itemId: varchar('item_id', { length: 50 }).notNull(),
  itemName: varchar('item_name', { length: 100 }).notNull(),
  quantity: integer('quantity').default(1).notNull(),
});

export const userWarnings = pgTable('user_warnings', {
  id: serial('id').primaryKey(),
  guildId: varchar('guild_id', { length: 20 }).notNull(),
  userId: varchar('user_id', { length: 20 }).notNull(),
  moderatorId: varchar('moderator_id', { length: 20 }).notNull(),
  reason: text('reason').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
});

export const commandLogs = pgTable('command_logs', {
  id: serial('id').primaryKey(),
  guildId: varchar('guild_id', { length: 20 }),
  userId: varchar('user_id', { length: 20 }).notNull(),
  commandName: varchar('command_name', { length: 50 }).notNull(),
  executedAt: timestamp('executed_at').defaultNow().notNull(),
});

export const botStats = pgTable('bot_stats', {
  id: integer('id').primaryKey(),
  guildCount: integer('guild_count').default(0).notNull(),
  userCount: integer('user_count').default(0).notNull(),
  commandsRun: integer('commands_run').default(0).notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});
