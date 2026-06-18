# Contributing to PulseKeep

Thank you for your interest in contributing to PulseKeep.

## Code of Conduct

By participating, you agree to uphold the [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

> [!NOTE]
> Check the [issues](https://github.com/watispro5212/PulseKeep/issues) for duplicates before reporting a new bug.

1. Use the **Bug Report** issue template.
2. Include steps to reproduce, expected behavior, and actual behavior.
3. Attach logs or screenshots if relevant.

### Suggesting Features

1. Open a **Feature Request** issue using the template
2. Describe the problem you're solving and your proposed solution
3. Explain how the feature benefits PulseKeep users

### Submitting Changes

1. Fork the repository.
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes following the coding conventions below.
4. Run checks: `npm run typecheck`
5. Commit with a descriptive message.
6. Open a Pull Request using the PR template.

## Development Setup

### Prerequisites

- Node.js 20+
- npm
- A Discord bot token (for full bot mode)
- (Optional) A PostgreSQL database

### Local Development

```bash
git clone https://github.com/watispro5212/PulseKeep.git
cd PulseKeep
cp .env.example .env
# Fill in DISCORD_TOKEN and DATABASE_URL in .env
npm install
npm run dev
```

### Project Structure

```text
src/
  api/               Express API server (stats, webhooks, health)
  bot/
    commands/        Slash command implementations
      economy/       Economy commands (balance, daily, work, etc.)
      moderation/    Moderation commands (warn, ban, kick, etc.)
      tickets/       Ticket system commands
    automod/         Auto-moderation engine (spam, mentions, caps, links, words)
    client.ts        Bot class with gateway events and cooldown system
    types.ts         Discord.js type extensions
  cache/             In-memory stats cache
  config.ts          Environment variable loader
  db/                Drizzle ORM schema and database pool
web/                 Static website (HTML, CSS, JS)
```

## Coding Conventions

> [!IMPORTANT]
> Every database error path must degrade gracefully — the bot must keep running.

- Run `npm run typecheck` before committing (must pass with zero errors).
- Use `EmbedBuilder` from discord.js for all message embeds (see `src/utils/embed.ts` for helpers like `formatNumber`, `formatCooldown`).
- Register slash commands in `src/bot/commands/index.ts` with their category.
- Keep economy logic in `src/bot/economy/store.ts`, not in command handlers.
- Prefix component custom IDs with `pulsekeep:namespace:action`.
- Use descriptive variable names; avoid single-letter names outside loops.
- Handle all errors; use `console.error` for non-critical failures (never `process.exit`).
- Use `Ephemeral` constant from `src/utils/embed.ts` for ephemeral responses.

## Testing

- No test framework is configured yet — manual testing via Discord is the current approach.
- Run `npm run typecheck` to verify TypeScript compilation before submitting changes.

## Pull Request Process

1. Ensure `npm run typecheck` passes with zero errors
2. Update documentation if you add or change commands
3. Update `web/commands.html` for new commands
4. Update `web/changelog.html` with a brief entry under the next version
5. Update command registration in `src/bot/commands/index.ts`
6. PRs require at least one review before merging

## Questions?

Open a [Discussion](https://github.com/watispro5212/PulseKeep/discussions) or join the support server.
