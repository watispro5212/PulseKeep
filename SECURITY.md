# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 5.x     | :white_check_mark: |
| < 5.0   | :x:                |

## Reporting a Vulnerability

PulseKeep takes security seriously. If you discover a security vulnerability, please **do not** open a public issue.

Instead, send a private report to the maintainers:

1. **GitHub Security Advisories**: Navigate to the repository's **Security** tab and use the "Report a vulnerability" feature
2. **Email**: Contact the maintainers directly (check repository commit history for contact information)

You should receive a response within **48 hours**. If you don't, please follow up.

### What to include

- Type of vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if applicable)

### What to expect

- Confirmation of receipt within 48 hours
- Regular updates on the fix timeline
- Credit in release notes once the fix is published

## Scope

This policy covers:

- The PulseKeep Discord bot and its source code
- The PulseKeep API server
- The PulseKeep web interface

Out of scope:

- Discord's platform and API
- Third-party libraries and dependencies (report to their maintainers)

## Best Practices

- The bot uses the minimum required Discord intents and permissions
- All moderation commands verify permissions at runtime
- Economy operations use mutex locks to prevent race conditions
- API endpoints return only aggregate statistics, not user-specific data
- No secrets or tokens are logged or exposed in responses

## Disclosure

We follow a 90-day disclosure timeline for fixed vulnerabilities to allow users time to update.
