# Contributing to PulseKeep

Thank you for your interest in contributing to PulseKeep! This document outlines the process for contributing to the project.

## Code of Conduct

By participating, you agree to uphold the [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

1. Check the [issues](https://github.com/watispro/pulsekeep/issues) to avoid duplicates
2. Use the **Bug Report** issue template
3. Include clear steps to reproduce, expected behavior, and actual behavior
4. Attach logs or screenshots if relevant

### Suggesting Features

1. Open a **Feature Request** issue using the template
2. Describe the problem you're solving and your proposed solution
3. Explain how the feature benefits the broader PulseKeep community

### Submitting Changes

1. Fork the repository
2. Create a branch: `git checkout -b feat/my-feature`
3. Make your changes following the coding conventions below
4. Run tests: `go test ./...`
5. Run vet: `go vet ./...`
6. Commit with a descriptive message
7. Open a Pull Request using the PR template

## Development Setup

### Prerequisites

- Go 1.22+
- A Discord bot token
- (Optional) A Neon PostgreSQL database for persistence

### Local Development

1. Clone the repo: `git clone https://github.com/watispro/pulsekeep.git`
2. Copy the example config: `cp .env.example .env`
3. Fill in your `DISCORD_BOT_TOKEN` in `.env`
4. Run: `go run ./cmd/pulsekeep`

### Project Structure

```
cmd/pulsekeep/        — Application entrypoint
internal/api/         — HTTP server (Gin) for health/stats endpoints
internal/bot/         — Discord bot logic
  commands/           — Slash command registration and menu system
  economy/            — In-memory economy store and shop system
  handlers.go         — Event handlers and command implementations
  economy_handlers.go — Economy command response builders
internal/cache/       — Atomic counters for live stats
internal/config/      — Environment configuration loading
internal/db/          — PostgreSQL database layer (optional)
web/                  — Static website (HTML/CSS/JS)
```

## Coding Conventions

- Follow standard Go formatting (`gofmt` / `gofumpt`)
- Use `discord.NewEmbed()` builder pattern for all message embeds
- Register all slash commands in `commands/commands.go`
- Add command descriptions and usage to `commands/menu.go` categories
- Keep economy logic in `economy/store.go`, not in handlers
- Prefix component custom IDs with `pulsekeep:namespace:action`
- Use descriptive variable names; avoid single-letter names outside loops
- Handle errors; log with `log.Printf` for non-critical failures
- Use ephemeral responses (`WithEphemeral(true)`) for confirmation messages

## Testing

- Tests use the standard `testing` package
- Run all tests: `go test ./...`
- Economy store tests are in `internal/bot/economy/store_test.go`
- Add tests for new store methods and edge cases

## Pull Request Process

1. Ensure all tests pass and `go vet` is clean
2. Update documentation if you add or change commands
3. Update `web/commands.html` and `commands/menu.go` for new commands
4. Update `web/changelog.html` with a brief entry under the next version
5. PRs require at least one review before merging

## Questions?

Open a [Discussion](https://github.com/watispro/pulsekeep/discussions) or join the PulseKeep support server.
