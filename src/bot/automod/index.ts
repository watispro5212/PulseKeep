export enum ActionType {
    None = 'none',
    Warn = 'warn',
    Delete = 'delete',
    Timeout = 'timeout',
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

export class AutomodEngine {
    private spamTimes: Map<string, number[]> = new Map();
    private configs: Map<string, GuildConfig> = new Map();

    constructor() {}

    public updateConfig(config: GuildConfig) {
        this.configs.set(config.guildID, config);
    }

    public getConfig(guildID: string): GuildConfig | undefined {
        return this.configs.get(guildID);
    }

    public checkMessage(guildID: string, userID: string, content: string): RuleResult {
        const cfg = this.configs.get(guildID);
        if (!cfg || !cfg.automodEnabled) {
            return { action: ActionType.None, reason: '', deleteMsg: false };
        }

        // 1. Spam detection
        if (cfg.spamEnabled) {
            const spamResult = this.checkSpam(guildID, userID, cfg);
            if (spamResult.action !== ActionType.None) return spamResult;
        }

        // 2. Mass mention protection
        if (cfg.mentionEnabled) {
            const mentionCount = (content.match(/<@!?\d+>|<@&\d+>|@everyone|@here/g) || []).length;
            if (mentionCount > cfg.mentionMax) {
                const action = this.parseAction(cfg.mentionAction);
                return { action, reason: 'Mass mention detected', deleteMsg: action === ActionType.Delete || action === ActionType.Timeout };
            }
        }

        // 3. Keyword filter
        if (cfg.bannedWords) {
            const words = cfg.bannedWords.split(',').map(w => w.trim().toLowerCase()).filter(w => w !== '');
            const lowerContent = content.toLowerCase();
            for (const word of words) {
                const regex = new RegExp(`\\b${this.escapeRegExp(word)}\\b`, 'i');
                if (regex.test(lowerContent)) {
                    return { action: ActionType.Delete, reason: `Message contained banned word: ${word}`, deleteMsg: true };
                }
            }
        }

        // 4. Link protection
        if (cfg.linksEnabled && /https?:\/\/[^\s]+/.test(content)) {
            const action = this.parseAction(cfg.linksAction);
            return { action, reason: 'Links are not allowed', deleteMsg: action === ActionType.Delete || action === ActionType.Timeout };
        }

        // 5. Caps lock protection
        if (cfg.capsEnabled && content.length >= cfg.capsMinLength) {
            if (this.isMostlyCaps(content, cfg.capsPercent)) {
                const action = this.parseAction(cfg.capsAction);
                return { action, reason: 'Excessive caps lock', deleteMsg: action === ActionType.Delete || action === ActionType.Timeout };
            }
        }

        return { action: ActionType.None, reason: '', deleteMsg: false };
    }

    private checkSpam(guildID: string, userID: string, cfg: GuildConfig): RuleResult {
        const key = `${guildID}:${userID}`;
        const now = Date.now();
        const times = this.spamTimes.get(key) || [];
        
        const cutoff = now - (cfg.spamWindowSecs * 1000);
        const recent = times.filter(t => t > cutoff);
        recent.push(now);
        this.spamTimes.set(key, recent);

        if (recent.length > cfg.spamMaxMessages) {
            const action = this.parseAction(cfg.spamAction);
            return { action, reason: 'Spam detected', deleteMsg: action === ActionType.Delete || action === ActionType.Timeout };
        }

        return { action: ActionType.None, reason: '', deleteMsg: false };
    }

    private parseAction(s: string): ActionType {
        switch (s.toLowerCase()) {
            case 'warn': return ActionType.Warn;
            case 'delete': return ActionType.Delete;
            case 'timeout': return ActionType.Timeout;
            default: return ActionType.Warn;
        }
    }

    private isMostlyCaps(s: string, percent: number): boolean {
        let upper = 0;
        let letters = 0;
        for (const char of s) {
            if (/[a-zA-Z]/.test(char)) {
                letters++;
                if (char === char.toUpperCase()) {
                    upper++;
                }
            }
        }
        if (letters === 0) return false;
        return (upper * 100 / letters) >= percent;
    }

    private escapeRegExp(string: string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }
}
