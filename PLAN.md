# PulseKeep TypeScript Rewrite Plan

## 1. Project Structure

```
PulseKeep-TS/
├── src/
│   ├── index.ts              // Main entry point
│   ├── config/               // Configuration loading
│   │   └── index.ts
│   ├── db/                   // Database connection and schema
│   │   ├── index.ts
│   │   └── schema.ts
│   ├── cache/                // In-memory cache
│   │   └── index.ts
│   ├── bot/                  // Discord bot logic
│   │   ├── index.ts          // Bot initialization and event handling
│   │   ├── commands/         // Command definitions and handlers
│   │   │   └── index.ts
│   │   ├── automod/          // Automod rules and engine
│   │   │   └── index.ts
│   │   ├── economy/          // Economy features
│   │   │   └── index.ts
│   │   └── utils/            // Utility functions for the bot
│   │       └── index.ts
│   └── api/                  // Web server and API endpoints
│       └── index.ts
├── .env.example
├── package.json
├── tsconfig.json
└── README.md
```

## 2. Module Mapping (Go to TypeScript)

| Go Module/File               | TypeScript Equivalent (Proposed)       | Description                                                               |
| :--------------------------- | :------------------------------------- | :------------------------------------------------------------------------ |
| `cmd/pulsekeep/main.go`      | `src/index.ts`                         | Main application entry point, orchestrating config, DB, cache, bot, and API. |
| `internal/config/config.go`  | `src/config/index.ts`                  | Handles loading and managing application configuration.                   |
| `internal/db/postgres.go`    | `src/db/index.ts`                      | Manages PostgreSQL database connection and migrations.                    |
| `db/schema.ts`               | `src/db/schema.ts`                     | Drizzle ORM schema definitions. (Already TypeScript)                      |
| `internal/cache/cache.go`    | `src/cache/index.ts`                   | In-memory caching mechanisms.                                             |
| `internal/bot/handlers.go`   | `src/bot/index.ts`                     | Discord bot event handlers and core logic.                                |
| `internal/bot/commands/commands.go` | `src/bot/commands/index.ts`            | Definitions and implementations of Discord commands.                      |
| `internal/bot/automod/rules.go` | `src/bot/automod/index.ts`             | Automoderation rules and logic.                                           |
| `internal/bot/economy/*.go`  | `src/bot/economy/index.ts`             | Economy-related features (blackjack, store, lottery).                     |
| `internal/api/server.go`     | `src/api/index.ts`                     | Web server setup and API endpoint definitions.                            |
| `internal/auth/discord.go`   | `src/api/auth.ts` (or similar)         | Discord OAuth2 authentication logic for the web API.                      |

## 3. Key Libraries and Technologies

*   **TypeScript**: Primary language for the rewrite.
*   **Node.js**: Runtime environment.
*   **Discord.js**: For interacting with the Discord API.
*   **Drizzle ORM**: For database interactions (already in use for schema).
*   **Express.js / Fastify**: For building the web API (replacing Gin Gonic).
*   **Dotenv**: For environment variable management.

## 4. Migration Strategy

1.  **Configuration**: Translate `config.go` into `src/config/index.ts`.
2.  **Database**: Adapt `internal/db/postgres.go` to `src/db/index.ts` using Drizzle ORM for queries.
3.  **Cache**: Rewrite `internal/cache/cache.go` to `src/cache/index.ts`.
4.  **Bot Core**: Translate `internal/bot/handlers.go` and `internal/bot/commands/*.go` into `src/bot/index.ts` and `src/bot/commands/index.ts` using Discord.js.
5.  **Automod**: Rewrite `internal/bot/automod/rules.go` to `src/bot/automod/index.ts`.
6.  **Economy**: Translate `internal/bot/economy/*.go` to `src/bot/economy/index.ts`.
7.  **API Server**: Rewrite `internal/api/server.go` to `src/api/index.ts` using a Node.js web framework.
8.  **Main Entry Point**: Create `src/index.ts` to tie everything together, similar to `cmd/pulsekeep/main.go`.

## 5. Next Steps

*   Implement configuration loading.
*   Set up basic database connection with Drizzle.
*   Initialize Discord.js client.
*   Start translating core bot functionalities. 
