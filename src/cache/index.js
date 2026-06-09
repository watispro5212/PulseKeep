"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Cache = void 0;
class Cache {
    guildsCount = 0;
    totalUserCount = 0;
    commandsRun = 0;
    startedAt = new Date();
    latencies = [];
    constructor() { }
    setGuildsCount(count) {
        this.guildsCount = count;
    }
    getGuildsCount() {
        return this.guildsCount;
    }
    setTotalUserCount(count) {
        this.totalUserCount = count;
    }
    getTotalUserCount() {
        return this.totalUserCount;
    }
    incrementCommandsRun() {
        this.commandsRun++;
    }
    getCommandsRun() {
        return this.commandsRun;
    }
    getStartedAt() {
        return this.startedAt;
    }
    addLatency(latency) {
        this.latencies.push(latency);
        if (this.latencies.length > 100) {
            this.latencies.shift();
        }
    }
    getAvgLatency() {
        if (this.latencies.length === 0)
            return 0;
        const sum = this.latencies.reduce((a, b) => a + b, 0);
        return sum / this.latencies.length;
    }
}
exports.Cache = Cache;
//# sourceMappingURL=index.js.map