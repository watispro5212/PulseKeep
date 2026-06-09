export interface Config {
    port: string;
    databaseURL: string;
    discordToken: string;
    botDisabled: boolean;
    statusWebhookURL: string;
    discordClientID: string;
    discordClientSecret: string;
    discordRedirectURI: string;
    allowedOrigins: string;
}
export declare const loadConfig: () => Config;
//# sourceMappingURL=index.d.ts.map