export class Cache {
  private guildsCount = 0;
  private totalUserCount = 0;
  private commandsRun = 0;
  private avgLatencyValue = 0;
  private startedAt = new Date();
  private botGuilds: { id: string; name: string }[] = [];
  private shardStats = new Map<number, { guilds: number; users: number; commandsRun: number; avgLatency: number; uptime: number }>();

  setGuildsCount(count: number) { this.guildsCount = count; }
  getGuildsCount(): number { return this.guildsCount; }

  setTotalUserCount(count: number) { this.totalUserCount = count; }
  getTotalUserCount(): number { return this.totalUserCount; }

  setCommandsRun(count: number) { this.commandsRun = count; }
  getCommandsRun(): number { return this.commandsRun; }
  incrementCommandsRun() { this.commandsRun++; }

  setAvgLatency(ms: number) { this.avgLatencyValue = ms; }
  getAvgLatency(): number { return this.avgLatencyValue; }

  getStartedAt(): Date { return this.startedAt; }

  setBotGuilds(guilds: { id: string; name: string }[]) { this.botGuilds = guilds; }
  getBotGuilds(): { id: string; name: string }[] { return this.botGuilds; }

  getShardCount(): number { return this.shardStats.size; }

  setShardStats(shardId: number, stats: { guilds: number; users: number; commandsRun: number; avgLatency: number; uptime: number }) {
    this.shardStats.set(shardId, stats);
  }

  getAllShardStats(): { guilds: number; users: number; commandsRun: number; avgLatency: number; uptime: number }[] {
    return Array.from(this.shardStats.values());
  }
}
