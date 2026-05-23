# PulseKeep Top.gg Submission Packet

Use this as the source copy when submitting PulseKeep to Top.gg.

## Core Details

Bot name: PulseKeep

Website: https://pulsekeep.netlify.app

API status: https://pulsekeep.fly.dev/health

Command page: https://pulsekeep.netlify.app/commands.html

Status page: https://pulsekeep.netlify.app/status.html

Privacy policy: https://pulsekeep.netlify.app/privacy.html

Terms of service: https://pulsekeep.netlify.app/terms.html

Support page: https://pulsekeep.netlify.app/support.html

Invite URL:
https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands

Prefix: Slash commands

Primary command for reviewers: `/help`

Backup test commands: `/ping`, `/menu`, `/stats`, `/uptime`

## Short Description

PulseKeep is a Go-powered Discord bot for moderation, audit logs, private support tickets, economy commands, and live server analytics.

## Long Description

PulseKeep helps Discord staff teams keep servers organized with practical slash commands and lightweight operations tooling.

Key features:

- Moderation commands including purge, kick, ban, and announce.
- Audit-friendly logging for moderation and server activity.
- Private ticket panel for support workflows.
- Economy commands including balance, profile, daily, work, pay, coinflip, and leaderboard.
- Utility commands including ping, uptime, stats, help, menu, server info, user info, and avatar.
- Live website status backed by the Fly.io API at `https://pulsekeep.fly.dev`.

The fastest way to test PulseKeep after inviting it is to run `/help`, `/menu`, or `/ping`.

## Categories

Suggested categories:

- Moderation
- Utility
- Economy
- Logging

## Review Checklist

- Keep the Fly.io bot service online during review.
- Confirm `https://pulsekeep.fly.dev/health` returns JSON before submitting.
- Confirm the OAuth invite works in a test server.
- Confirm `/help`, `/menu`, `/ping`, `/stats`, and `/uptime` are usable by reviewers.
- Make sure the public website is deployed after the latest `web/` changes.
