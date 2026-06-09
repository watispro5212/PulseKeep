"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EconomyStore = exports.STARTING_BALANCE = void 0;
const discord_js_1 = require("discord.js");
exports.STARTING_BALANCE = 500;
class EconomyStore {
    records = new Map();
    db; // Will be the Drizzle instance
    constructor(db) {
        this.db = db;
    }
    async getRecord(userID, name = '') {
        let record = this.records.get(userID);
        if (!record) {
            // In a real app, we would fetch from DB here
            record = {
                userID,
                name,
                balance: exports.STARTING_BALANCE,
                lastInterest: new Date(),
                inventory: new Map(),
                createdAt: new Date(),
                updatedAt: new Date(),
                flipWins: 0,
                flipLosses: 0,
            };
            this.records.set(userID, record);
        }
        else if (name) {
            record.name = name;
        }
        return record;
    }
    async addBalance(userID, amount) {
        const record = await this.getRecord(userID);
        record.balance += amount;
        record.updatedAt = new Date();
        return record.balance;
    }
}
exports.EconomyStore = EconomyStore;
//# sourceMappingURL=index.js.map