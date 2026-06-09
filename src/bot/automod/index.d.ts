export declare enum ActionType {
    None = "none",
    Warn = "warn",
    Delete = "delete",
    Timeout = "timeout"
}
export interface RuleResult {
    action: ActionType;
    reason: string;
    deleteMsg: boolean;
}
export interface GuildConfig {
    guildID: string;
    automodEnabled: boolean;
    spamEnabled: boolean;
    spamMaxMessages: number;
    spamWindowSecs: number;
    spamAction: string;
    mentionEnabled: boolean;
    mentionMax: number;
    mentionAction: string;
    bannedWords: string;
    linksEnabled: boolean;
    linksAction: string;
    capsEnabled: boolean;
    capsMinLength: number;
    capsPercent: number;
    capsAction: string;
}
export declare class AutomodEngine {
    private spamTimes;
    private configs;
    constructor();
    updateConfig(config: GuildConfig): void;
    getConfig(guildID: string): GuildConfig | undefined;
    checkMessage(guildID: string, userID: string, content: string): RuleResult;
    private checkSpam;
    private parseAction;
    private isMostlyCaps;
    private escapeRegExp;
}
//# sourceMappingURL=index.d.ts.map