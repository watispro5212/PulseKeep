import type { Cache } from '../cache/index.js';
import type { Bot } from './client.js';
import type { Config } from '../config.js';

export type DB = any;

export interface CommandContext {
  bot: Bot;
  cache: Cache;
  db: DB;
  config: Config;
}

export interface SlashCommand {
  data: any;
  execute(ctx: CommandContext, interaction: any): Promise<void>;
}
