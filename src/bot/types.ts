import type { Cache } from '../cache/index.js';
import type { Bot } from './client.js';

export interface CommandContext {
  bot: Bot;
  cache: Cache;
  db: any;
}

export interface SlashCommand {
  data: any;
  execute(ctx: CommandContext, interaction: any): Promise<void>;
}
