import dotenv from 'dotenv';
dotenv.config();

export interface Config {
  port: string;
  databaseURL: string;
  discordToken: string;
  discordClientID: string;
  discordClientSecret: string;
  discordRedirectURI: string;
  statusWebhookURL: string;
  allowedOrigins: string;
  botDisabled: boolean;
  dblWebhookSecret: string;
  dblApiToken: string;
  discordBotID: string;
  discordsWebhookSecret: string;
  botOwnerID: string;
  botCoOwnerID: string;
  cooldownEconomy: number;
  cooldownModeration: number;
}

export function loadConfig(): Config {
  return {
    port: process.env.PORT || '8080',
    databaseURL: process.env.DATABASE_URL || '',
    discordToken: process.env.DISCORD_TOKEN || '',
    discordClientID: process.env.DISCORD_CLIENT_ID || '',
    discordClientSecret: process.env.DISCORD_CLIENT_SECRET || '',
    discordRedirectURI: process.env.DISCORD_REDIRECT_URI || '',
    statusWebhookURL: process.env.STATUS_WEBHOOK_URL || '',
    allowedOrigins: process.env.ALLOWED_ORIGINS || '',
    botDisabled: process.env.BOT_DISABLED === 'true',
    dblWebhookSecret: process.env.DBL_WEBHOOK_SECRET || '',
    dblApiToken: process.env.DBL_API_TOKEN || '',
    discordBotID: process.env.DISCORD_BOT_ID || '',
    discordsWebhookSecret: process.env.DISCORDS_WEBHOOK_SECRET || '',
    botOwnerID: process.env.BOT_OWNER_ID || '',
    botCoOwnerID: process.env.BOT_CO_OWNER_ID || '',
    cooldownEconomy: Number(process.env.COOLDOWN_ECONOMY) || 3,
    cooldownModeration: Number(process.env.COOLDOWN_MODERATION) || 2,
  };
}
