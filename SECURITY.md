# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 7.x     | :white_check_mark: |
| < 7.0   | :x:                |

## Reporting a Vulnerability

Found something bad? Don't post it publicly.

Send a private report:

1. **GitHub Security Advisories** — repo's Security tab, "Report a vulnerability"
2. **Email** — check commit history for maintainer contact

Expect a response within ~48 hours.

### What to include

- What kind of vulnerability
- How to reproduce it
- What it could do
- A fix if you have one

## Scope

In scope:
- The bot, the API server, the website

Out of scope:
- Discord itself
- Third-party libs (report to them)

## Best Practices

- Minimum required intents and permissions
- Runtime permission checks on all mod commands
- Mutex locks on economy operations
- API returns aggregates only, no user-specific data
- No secrets logged or exposed
- Graceful error handling everywhere — no crash-on-fail

## Disclosure

90-day timeline after fix to give people time to update.
