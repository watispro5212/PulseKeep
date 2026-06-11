export class Cache {
  private guildsCount = 0;
  private totalUserCount = 0;
  private commandsRun = 0;
  private startedAt = new Date();
  private latencies: number[] = [];

  setGuildsCount(count: number) { this.guildsCount = count; }
  getGuildsCount(): number { return this.guildsCount; }

  setTotalUserCount(count: number) { this.totalUserCount = count; }
  getTotalUserCount(): number { return this.totalUserCount; }

  incrementCommandsRun() { this.commandsRun++; }
  getCommandsRun(): number { return this.commandsRun; }

  getStartedAt(): Date { return this.startedAt; }

  addLatency(ms: number) {
    this.latencies.push(ms);
    if (this.latencies.length > 100) this.latencies.shift();
  }

  getAvgLatency(): number {
    if (this.latencies.length === 0) return 0;
    return this.latencies.reduce((a, b) => a + b, 0) / this.latencies.length;
  }
}
