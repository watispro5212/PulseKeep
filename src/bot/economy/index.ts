import type { Snowflake } from 'discord.js';

export interface InventoryEntry {
    itemID: string;
    itemName: string;
    quantity: number;
}

export interface Record {
    userID: string;
    name: string;
    balance: number;
    lastInterest: Date;
    inventory: Map<string, InventoryEntry>;
    createdAt: Date;
    updatedAt: Date;
    flipWins: number;
    flipLosses: number;
}

export const STARTING_BALANCE = 500;

export class EconomyStore {
    private records: Map<string, Record> = new Map();
    private db: any; // Will be the Drizzle instance

    constructor(db: any) {
        this.db = db;
    }

    public async getRecord(userID: string, name: string = ''): Promise<Record> {
        let record = this.records.get(userID);
        if (!record) {
            // In a real app, we would fetch from DB here
            record = {
                userID,
                name,
                balance: STARTING_BALANCE,
                lastInterest: new Date(),
                inventory: new Map(),
                createdAt: new Date(),
                updatedAt: new Date(),
                flipWins: 0,
                flipLosses: 0,
            };
            this.records.set(userID, record);
        } else if (name) {
            record.name = name;
        }
        return record;
    }

    public async addBalance(userID: string, amount: number): Promise<number> {
        const record = await this.getRecord(userID);
        record.balance += amount;
        record.updatedAt = new Date();
        return record.balance;
    }

    // Additional economy methods would go here (lottery, store, etc.)
}
