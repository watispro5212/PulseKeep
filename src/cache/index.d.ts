export declare class Cache {
    private guildsCount;
    private totalUserCount;
    private commandsRun;
    private startedAt;
    private latencies;
    constructor();
    setGuildsCount(count: number): void;
    getGuildsCount(): number;
    setTotalUserCount(count: number): void;
    getTotalUserCount(): number;
    incrementCommandsRun(): void;
    getCommandsRun(): number;
    getStartedAt(): Date;
    addLatency(latency: number): void;
    getAvgLatency(): number;
}
//# sourceMappingURL=index.d.ts.map