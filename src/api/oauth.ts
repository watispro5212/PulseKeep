import type { Config } from '../config.js';

const DISCORD_API = 'https://discord.com/api/v10';

export function getOAuthURL(config: Config, state: string): string {
  const url = new URL(`${DISCORD_API}/oauth2/authorize`);
  url.searchParams.set('client_id', config.discordClientID);
  url.searchParams.set('redirect_uri', config.discordRedirectURI);
  url.searchParams.set('response_type', 'code');
  url.searchParams.set('scope', 'identify guilds');
  url.searchParams.set('state', state);
  return url.toString();
}

export async function exchangeCode(config: Config, code: string): Promise<any> {
  const body = new URLSearchParams({
    client_id: config.discordClientID,
    client_secret: config.discordClientSecret,
    grant_type: 'authorization_code',
    code,
    redirect_uri: config.discordRedirectURI,
  });

  const res = await fetch(`${DISCORD_API}/oauth2/token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Token exchange failed: ${text}`);
  }

  return res.json();
}

export async function fetchUser(accessToken: string): Promise<any> {
  const res = await fetch(`${DISCORD_API}/users/@me`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  });
  if (!res.ok) throw new Error('Failed to fetch user');
  return res.json();
}

export async function fetchGuilds(accessToken: string): Promise<any[]> {
  const res = await fetch(`${DISCORD_API}/users/@me/guilds`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  });
  if (!res.ok) throw new Error('Failed to fetch guilds');
  return res.json();
}

export async function fetchMutualGuilds(userGuilds: any[], botGuilds: any[]): Promise<any[]> {
  const botGuildIds = new Set(botGuilds.map((g: any) => g.id));
  return userGuilds
    .filter((g: any) => botGuildIds.has(g.id))
    .map((g: any) => ({
      id: g.id,
      name: g.name,
      icon: g.icon,
      owner: g.owner,
      permissions: g.permissions,
    }));
}
