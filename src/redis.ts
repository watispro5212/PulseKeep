let redis: any = null;

export function getRedis(url?: string) {
  if (redis) return redis;
  if (!url) return null;
  try {
    const Redis = require('ioredis');
    redis = new Redis(url, {
      maxRetriesPerRequest: 3,
      retryStrategy(times: number) { return Math.min(times * 200, 3000); },
      lazyConnect: true,
    });
    redis.on('error', (err: Error) => console.error('[Redis] Error:', err.message));
    return redis;
  } catch {
    return null;
  }
}

export function closeRedis() {
  if (redis) { redis.disconnect(); redis = null; }
}
