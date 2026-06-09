"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AutomodEngine = exports.ActionType = void 0;
var ActionType;
(function (ActionType) {
    ActionType["None"] = "none";
    ActionType["Warn"] = "warn";
    ActionType["Delete"] = "delete";
    ActionType["Timeout"] = "timeout";
})(ActionType || (exports.ActionType = ActionType = {}));
class AutomodEngine {
    spamTimes = new Map();
    configs = new Map();
    constructor() { }
    updateConfig(config) {
        this.configs.set(config.guildID, config);
    }
    getConfig(guildID) {
        return this.configs.get(guildID);
    }
    checkMessage(guildID, userID, content) {
        const cfg = this.configs.get(guildID);
        if (!cfg || !cfg.automodEnabled) {
            return { action: ActionType.None, reason: '', deleteMsg: false };
        }
        // 1. Spam detection
        if (cfg.spamEnabled) {
            const spamResult = this.checkSpam(guildID, userID, cfg);
            if (spamResult.action !== ActionType.None)
                return spamResult;
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
    checkSpam(guildID, userID, cfg) {
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
    parseAction(s) {
        switch (s.toLowerCase()) {
            case 'warn': return ActionType.Warn;
            case 'delete': return ActionType.Delete;
            case 'timeout': return ActionType.Timeout;
            default: return ActionType.Warn;
        }
    }
    isMostlyCaps(s, percent) {
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
        if (letters === 0)
            return false;
        return (upper * 100 / letters) >= percent;
    }
    escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }
}
exports.AutomodEngine = AutomodEngine;
//# sourceMappingURL=index.js.map