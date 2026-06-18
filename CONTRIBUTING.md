# Contributing

By contributing you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).

## Bugs

Search [issues](https://github.com/watispro5212/PulseKeep/issues) first. Use the bug report template, include steps to reproduce.

## Features

Open a feature request. Say what problem you're solving and how.

## Pull Requests

1. Fork, branch: `git checkout -b feat/my-thing`
2. Make changes
3. `npm run typecheck` — must pass
4. Commit, open PR

## Setup

- Node.js 20+, npm
- Discord bot token (for full mode)
- PostgreSQL (optional)

```bash
git clone https://github.com/watispro5212/PulseKeep.git
cd PulseKeep
cp .env.example .env
# edit .env
npm install
npm run dev
```

## Structure

```
src/
  api/               Express server
  bot/
    commands/        Slash commands
      economy/       Economy stuff
      moderation/    Moderation stuff
      tickets/       Ticket stuff
    automod/         Auto-mod engine
    client.ts        Bot class, events, cooldowns
    types.ts         Type extensions
  cache/             In-memory cache
  config.ts          Env loader
  db/                Drizzle schema + pool
web/                 Static site
```

## Conventions

- DB errors must never crash the bot
- `npm run typecheck` before committing
- Use `EmbedBuilder` for embeds (see `src/utils/embed.ts` for helpers)
- Register commands in `src/bot/commands/index.ts`
- Economy logic in `src/bot/economy/store.ts`
- Custom IDs: `pulsekeep:namespace:action`
- Handle errors with `console.error`, never `process.exit`
- Use `Ephemeral` from `src/utils/embed.ts`

## Testing

No tests yet — manual testing for now.

## PR Checklist

- [ ] `npm run typecheck` passes
- [ ] Updated `web/commands.html` for new commands
- [ ] Updated `web/changelog.html`
- [ ] Updated command registration
- [ ] At least one review

## Questions?

Open a [Discussion](https://github.com/watispro5212/PulseKeep/discussions) or join the support server.
