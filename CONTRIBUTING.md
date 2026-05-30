# Contributing to PulseKeep

Thank you for your interest in contributing to PulseKeep.

## Code of Conduct

By participating, you agree to uphold the [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

1. Check the [issues](https://github.com/watispro5212/PulseKeep/issues) for duplicates
2. Use the **Bug Report** issue template
3. Include steps to reproduce, expected behavior, and actual behavior
4. Attach logs or screenshots if relevant

### Suggesting Features

1. Open a **Feature Request** issue using the template
2. Describe the problem you're solving and your proposed solution
3. Explain how the feature benefits PulseKeep users

### Submitting Changes

1. Fork the repository
2. Create a branch: `git checkout -b feat/my-feature`
3. Make your changes following the coding conventions below
4. Run checks: `go build ./... && go vet ./... && go test ./...`
5. Commit with a descriptive message
6. Open a Pull Request using the PR template

## Development Setup

### Prerequisites

- Go 1.26+
- A Discord bot token (for full bot mode)
- (Optional) A Neon PostgreSQL database

### Local Development

```bash
git clone https://github.com/watispro5212/PulseKeep.git
cd PulseKeep
cp .env.example .env
# Fill in DISCORD_TOKEN in .env
go run ./cmd/pulsekeep
```

### Project Structure

```
cmd/pulsekeep/        — Application entrypoint
internal/api/         — HTTP server (Gin) for health/stats/OAuth
internal/auth/        — Discord OAuth2 token exchange and API calls
internal/bot/         — Discord gateway client and command handlers
  commands/           — Slash command registration
  economy/            — In-memory economy store with PostgreSQL persistence
  automod/            — Configurable auto-moderation engine
  handlers.go         — Event listeners and command dispatch
  economy_handlers.go — Economy embed builders
internal/cache/       — Atomic counters for live stats
internal/config/      — Environment configuration loader
internal/db/          — PostgreSQL connection and migrations
web/                  — Cloudflare Pages static website (11 pages)
```

## Coding Conventions

- Run `gofmt` before committing
- Use `discord.NewEmbed()` builder for all message embeds
- Register slash commands in `commands/commands.go`
- Add descriptions to `commands/menu.go` categories
- Keep economy logic in `economy/store.go`, not in handlers
- Prefix component custom IDs with `pulsekeep:namespace:action`
- Use descriptive variable names; avoid single-letter names outside loops
- Handle all errors; use `log.Printf` for non-critical failures (never `log.Fatalf`)
- Use ephemeral responses via `WithEphemeral(true)` for confirmations
- Every database error path must degrade gracefully — the web server must keep running

## Testing

- Tests use the standard `testing` package
- Run all tests: `go test ./...`
- Economy store tests are in `internal/bot/economy/store_test.go`
- Add tests for new store methods and edge cases

## Pull Request Process

1. Ensure `go build ./...` and `go vet ./...` pass
2. Update documentation if you add or change commands
3. Update `web/commands.html` for new commands
4. Update `web/changelog.html` with a brief entry under the next version
5. Update `commands/menu.go` with new command descriptions
6. PRs require at least one review before merging

## Questions?

Open a [Discussion](https://github.com/watispro5212/PulseKeep/discussions) or join the support server.
