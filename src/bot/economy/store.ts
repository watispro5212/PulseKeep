export const STARTING_BALANCE = 500;

export const COOLDOWNS = {
  daily: 24 * 60 * 60 * 1000,
  weekly: 7 * 24 * 60 * 60 * 1000,
  work: 60 * 60 * 1000,
  rob: 4 * 60 * 60 * 1000,
  fish: 30 * 1000,
  mine: 30 * 1000,
} as const;

export const WORK_JOBS = [
  { title: 'Software Engineer', pay: [200, 500] },
  { title: 'Chef', pay: [150, 400] },
  { title: 'Teacher', pay: [100, 350] },
  { title: 'Doctor', pay: [300, 600] },
  { title: 'Artist', pay: [100, 300] },
  { title: 'Musician', pay: [150, 350] },
  { title: 'Writer', pay: [120, 320] },
  { title: 'Pilot', pay: [250, 550] },
  { title: 'Architect', pay: [200, 500] },
  { title: 'Detective', pay: [180, 450] },
  { title: 'Astronaut', pay: [350, 700] },
  { title: 'Farmer', pay: [80, 250] },
  { title: 'Mechanic', pay: [140, 380] },
  { title: 'Electrician', pay: [160, 400] },
  { title: 'Lawyer', pay: [250, 600] },
  { title: 'Journalist', pay: [130, 350] },
  { title: 'Photographer', pay: [110, 310] },
  { title: 'Scientist', pay: [280, 580] },
  { title: 'Nurse', pay: [180, 420] },
  { title: 'Streamer', pay: [100, 800] },
  { title: 'Truck Driver', pay: [150, 380] },
  { title: 'Bartender', pay: [100, 280] },
];

export const WORK_FLAVOR = [
  'finished a shift at work and earned',
  'completed a freelance project and earned',
  'did some overtime and earned',
  'found some extra work and earned',
  'helped a colleague and was tipped',
  'solved a tough problem and earned',
  'wrapped up a project and got paid',
  'picked up a side hustle and made',
];

export const SHOP_ITEMS = [
  { id: 'fishing_rod', name: 'Fishing Rod 🎣', price: 2500, description: 'Lets you fish for treasures' },
  { id: 'mining_pick', name: 'Mining Pick ⛏️', price: 3000, description: 'Lets you mine for valuables' },
  { id: 'lucky_clover', name: 'Lucky Clover 🍀', price: 5000, description: 'Reroll a gamble loss once' },
  { id: 'xp_boost', name: 'XP Boost ⚡', price: 4000, description: 'Double earnings for 30 minutes' },
  { id: 'treasure_map', name: 'Treasure Map 🗺️', price: 7500, description: 'Claim 2500–7500 Pulses instantly' },
];

export const STREAK_MILESTONES = [7, 14, 21, 30, 50, 100] as const;

export const ECONOMY_TIPS = [
  'Use /daily every 24h for free Pulses!',
  'Stack your daily streak for bigger bonuses at 7, 14, 21, 30, 50, and 100 days.',
  '/work lets you earn 100–800 Pulses per hour depending on your job.',
  '/fish and /mine require a Fishing Rod or Mining Pick from the shop.',
  '/gamble gives you a 55% chance to win 2x–10x your bet!',
  'Use /rob to steal from other users — 45% success rate.',
  'The /slots machine has a 10% chance to upgrade your payout tier.',
  'Lucky Clover lets you reroll one gamble loss — buy it from /shop.',
  'XP Boost doubles all earnings for 30 minutes.',
  'Treasure Map gives 2500–7500 instant Pulses!',
  'Check /leaderboard to see the richest users.',
  'Your net worth = balance + inventory item value.',
];

const FISH_CATCHES = [
  { name: 'Common Carp', min: 50, max: 150, weight: 40 },
  { name: 'Rainbow Trout', min: 100, max: 300, weight: 25 },
  { name: 'Bass', min: 150, max: 400, weight: 15 },
  { name: 'Salmon', min: 200, max: 500, weight: 10 },
  { name: 'Golden Fish', min: 500, max: 1500, weight: 5 },
  { name: 'Junk Boot', min: 5, max: 25, weight: 5 },
];

const MINE_FINDS = [
  { name: 'Iron Ore', min: 50, max: 150, weight: 35 },
  { name: 'Copper Ore', min: 80, max: 200, weight: 25 },
  { name: 'Silver Ore', min: 150, max: 350, weight: 18 },
  { name: 'Gold Ore', min: 300, max: 600, weight: 12 },
  { name: 'Gemstone', min: 500, max: 2000, weight: 5 },
  { name: 'Rusty Nail', min: 5, max: 25, weight: 5 },
];

function weightedRandom(items: { weight: number }[]): number {
  const total = items.reduce((s, i) => s + i.weight, 0);
  let r = Math.random() * total;
  for (let i = 0; i < items.length; i++) {
    const item = items[i];
    if (!item) continue;
    r -= item.weight;
    if (r <= 0) return i;
  }
  return items.length - 1;
}

export function gamble(amount: number): { result: 'win' | 'lose' | 'push'; payout: number; multiplier: number } {
  const roll = Math.floor(Math.random() * 100) + 1;
  if (roll <= 35) return { result: 'lose', payout: 0, multiplier: 0 };
  if (roll <= 45) return { result: 'push', payout: amount, multiplier: 1 };
  const multiplier = 2 + Math.floor(Math.random() * 9);
  return { result: 'win', payout: amount * multiplier, multiplier };
}

export function blackjackDealerShouldStand(hand: number): boolean {
  return hand >= 16;
}

export function robSuccess(): boolean {
  return Math.random() < 0.45;
}

export function rollSlots(): { payout: number; emojis: string[] } {
  const emojis = ['🍒', '🍋', '🍊', '🍇', '💎', '7️⃣'];
  const result: string[] = [];
  for (let i = 0; i < 3; i++) {
    const e = emojis[Math.floor(Math.random() * emojis.length)];
    result.push(e ?? '🍒');
  }

  const c0 = result[0]!;
  const c1 = result[1]!;
  const c2 = result[2]!;
  const allSame = c0 === c1 && c1 === c2;
  const twoMatch = c0 === c1 || c1 === c2;

  if (allSame) {
    const sym = c0;
    const upgrade = Math.random() < 0.1;
    const base = sym === '7️⃣' ? 10 : sym === '💎' ? 7 : sym === '🍇' ? 5 : sym === '🍊' ? 4 : sym === '🍋' ? 3 : 2;
    return { payout: upgrade ? base * 2 : base, emojis: result };
  }

  if (twoMatch && (c0 === c1 || c1 === c2)) {
    return { payout: 1.5, emojis: result };
  }

  return { payout: 0, emojis: result };
}

export function fish(): { name: string; value: number } {
  const idx = weightedRandom(FISH_CATCHES);
  const c = FISH_CATCHES[idx];
  if (!c) return { name: 'Something Weird', value: 0 };
  const value = c.min + Math.floor(Math.random() * (c.max - c.min + 1));
  if (c.name === 'Junk Boot' && Math.random() < 0.15) return fish();
  return { name: c.name, value };
}

export function mine(): { name: string; value: number } {
  const idx = weightedRandom(MINE_FINDS);
  const c = MINE_FINDS[idx];
  if (!c) return { name: 'Something Weird', value: 0 };
  const value = c.min + Math.floor(Math.random() * (c.max - c.min + 1));
  if (c.name === 'Rusty Nail' && Math.random() < 0.15) return mine();
  return { name: c.name, value };
}

export function getStreakBonus(days: number): number {
  let bonus = 0;
  for (const m of STREAK_MILESTONES) {
    if (days >= m) bonus += m * 10;
  }
  return bonus;
}

export const SEARCH_PLACES = [
  { name: 'the old library', min: 50, max: 300, weight: 30 },
  { name: 'the abandoned mansion', min: 100, max: 500, weight: 20 },
  { name: 'the dark alley', min: 30, max: 200, weight: 25 },
  { name: 'the ancient temple', min: 200, max: 800, weight: 10 },
  { name: 'the sewer tunnels', min: 20, max: 150, weight: 15 },
  { name: 'the rooftop garden', min: 150, max: 600, weight: 12 },
  { name: 'the hidden bunker', min: 300, max: 1200, weight: 5 },
  { name: 'the arcade', min: 80, max: 350, weight: 18 },
  { name: 'the flea market', min: 40, max: 250, weight: 22 },
  { name: 'the construction site', min: 100, max: 400, weight: 15 },
  { name: 'the beach shore', min: 60, max: 280, weight: 20 },
  { name: 'the haunted forest', min: 120, max: 700, weight: 8 },
];

export function search(): { name: string; value: number } {
  const idx = weightedRandom(SEARCH_PLACES);
  const p = SEARCH_PLACES[idx];
  if (!p) return { name: 'nowhere', value: 0 };
  const value = p.min + Math.floor(Math.random() * (p.max - p.min + 1));
  return { name: p.name, value };
}

export function getStreakMilestoneProgress(days: number): string {
  for (const m of STREAK_MILESTONES) {
    if (days < m) return `${days}/${m} days`;
  }
  return 'Maxed out!';
}

const XP_BOOST_DURATION_MS = 30 * 60 * 1000;

export function hasXpBoost(rec: any): boolean {
  if (!rec?.xpBoostExpiry) return false;
  return new Date(rec.xpBoostExpiry).getTime() > Date.now();
}

export function applyXpBoost(earned: number, rec: any): number {
  return hasXpBoost(rec) ? earned * 2 : earned;
}

export function formatCooldown(elapsed: number, cooldown: number): string {
  const remaining = cooldown - elapsed;
  if (remaining <= 0) return 'Available now';
  const mins = Math.ceil(remaining / 60000);
  if (mins >= 60) {
    const hrs = Math.floor(mins / 60);
    return `${hrs}h ${mins % 60}m`;
  }
  return `${mins}m`;
}
