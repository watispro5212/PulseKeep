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
  };
}
