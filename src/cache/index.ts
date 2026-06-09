export class Cache {
    private guildsCount: number = 0;
    private totalUserCount: number = 0;
    private commandsRun: number = 0;
    private startedAt: Date = new Date();
    private latencies: number[] = [];

    constructor() {}

    public setGuildsCount(count: number) {
        this.guildsCount = count;
    }

    public getGuildsCount(): number {
        return this.guildsCount;
    }

    public setTotalUserCount(count: number) {
        this.totalUserCount = count;
    }

    public getTotalUserCount(): number {
        return this.totalUserCount;
    }

    public incrementCommandsRun() {
        this.commandsRun++;
    }

    public getCommandsRun(): number {
        return this.commandsRun;
    }

    public getStartedAt(): Date {
        return this.startedAt;
    }

    public addLatency(latency: number) {
        this.latencies.push(latency);
        if (this.latencies.length > 100) {
            this.latencies.shift();
        }
    }

    public getAvgLatency(): number {
        if (this.latencies.length === 0) return 0;
        const sum = this.latencies.reduce((a, b) => a + b, 0);
        return sum / this.latencies.length;
    }
}
