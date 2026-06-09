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
export declare const STARTING_BALANCE = 500;
export declare class EconomyStore {
    private records;
    private db;
    constructor(db: any);
    getRecord(userID: string, name?: string): Promise<Record>;
    addBalance(userID: string, amount: number): Promise<number>;
}
//# sourceMappingURL=index.d.ts.map