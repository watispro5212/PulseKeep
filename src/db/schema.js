"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.userInventory = exports.userEconomy = exports.commandLogs = exports.botStats = exports.guildConfigs = void 0;
const pg_core_1 = require("drizzle-orm/pg-core");
exports.guildConfigs = (0, pg_core_1.pgTable)('guild_configs', {
    guildId: (0, pg_core_1.varchar)('guild_id', { length: 20 }).primaryKey(),
    prefix: (0, pg_core_1.varchar)('prefix', { length: 10 }).default('!').notNull(),
    logChannelId: (0, pg_core_1.varchar)('log_channel_id', { length: 20 }),
    premium: (0, pg_core_1.boolean)('premium').default(false).notNull(),
    createdAt: (0, pg_core_1.timestamp)('created_at').defaultNow().notNull(),
    updatedAt: (0, pg_core_1.timestamp)('updated_at').defaultNow().notNull(),
});
exports.botStats = (0, pg_core_1.pgTable)('bot_stats', {
    id: (0, pg_core_1.integer)('id').primaryKey(),
    guildCount: (0, pg_core_1.integer)('guild_count').default(0).notNull(),
    userCount: (0, pg_core_1.integer)('user_count').default(0).notNull(),
    commandsRun: (0, pg_core_1.integer)('commands_run').default(0).notNull(),
    updatedAt: (0, pg_core_1.timestamp)('updated_at').defaultNow().notNull(),
});
exports.commandLogs = (0, pg_core_1.pgTable)('command_logs', {
    id: (0, pg_core_1.serial)('id').primaryKey(),
    guildId: (0, pg_core_1.varchar)('guild_id', { length: 20 }),
    userId: (0, pg_core_1.varchar)('user_id', { length: 20 }).notNull(),
    commandName: (0, pg_core_1.varchar)('command_name', { length: 50 }).notNull(),
    executedAt: (0, pg_core_1.timestamp)('executed_at').defaultNow().notNull(),
});
exports.userEconomy = (0, pg_core_1.pgTable)('user_economy', {
    userId: (0, pg_core_1.varchar)('user_id', { length: 20 }).primaryKey(),
    balance: (0, pg_core_1.integer)('balance').default(0).notNull(),
    lastDailyClaim: (0, pg_core_1.timestamp)('last_daily_claim'),
    lastWork: (0, pg_core_1.timestamp)('last_work'),
    lastRob: (0, pg_core_1.timestamp)('last_rob'),
    totalEarned: (0, pg_core_1.integer)('total_earned').default(0).notNull(),
    totalGambled: (0, pg_core_1.integer)('total_gambled').default(0).notNull(),
    transactions: (0, pg_core_1.integer)('transactions').default(0).notNull(),
});
exports.userInventory = (0, pg_core_1.pgTable)('user_inventory', {
    id: (0, pg_core_1.serial)('id').primaryKey(),
    userId: (0, pg_core_1.varchar)('user_id', { length: 20 }).notNull(),
    itemId: (0, pg_core_1.varchar)('item_id', { length: 50 }).notNull(),
    itemName: (0, pg_core_1.varchar)('item_name', { length: 100 }).notNull(),
    quantity: (0, pg_core_1.integer)('quantity').default(1).notNull(),
});
//# sourceMappingURL=schema.js.map