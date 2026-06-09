"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.loadConfig = void 0;
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
const loadConfig = () => {
    return {
        port: process.env.PORT || '8080',
        databaseURL: process.env.DATABASE_URL || '',
        discordToken: process.env.DISCORD_TOKEN || '',
        botDisabled: process.env.BOT_DISABLED === 'true',
        statusWebhookURL: process.env.STATUS_WEBHOOK_URL || '',
        discordClientID: process.env.DISCORD_CLIENT_ID || '',
        discordClientSecret: process.env.DISCORD_CLIENT_SECRET || '',
        discordRedirectURI: process.env.DISCORD_REDIRECT_URI || '',
        allowedOrigins: process.env.ALLOWED_ORIGINS || '',
    };
};
exports.loadConfig = loadConfig;
//# sourceMappingURL=index.js.map