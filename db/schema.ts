import { pgTable, varchar, boolean, timestamp, integer, serial } from 'drizzle-orm/pg-core'

export const guildConfigs = pgTable('guild_configs', {
  guildId: varchar('guild_id', { length: 20 }).primaryKey(),
  prefix: varchar('prefix', { length: 10 }).default('!').notNull(),
  logChannelId: varchar('log_channel_id', { length: 20 }),
  premium: boolean('premium').default(false).notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
})

export const botStats = pgTable('bot_stats', {
  id: integer('id').primaryKey(),
  guildCount: integer('guild_count').default(0).notNull(),
  userCount: integer('user_count').default(0).notNull(),
  commandsRun: integer('commands_run').default(0).notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
})

export const commandLogs = pgTable('command_logs', {
  id: serial('id').primaryKey(),
  guildId: varchar('guild_id', { length: 20 }),
  userId: varchar('user_id', { length: 20 }).notNull(),
  commandName: varchar('command_name', { length: 50 }).notNull(),
  executedAt: timestamp('executed_at').defaultNow().notNull(),
})
