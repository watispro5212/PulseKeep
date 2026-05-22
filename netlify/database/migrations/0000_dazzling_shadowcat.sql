CREATE TABLE IF NOT EXISTS "bot_stats" (
	"id" integer PRIMARY KEY NOT NULL,
	"guild_count" integer DEFAULT 0 NOT NULL,
	"user_count" integer DEFAULT 0 NOT NULL,
	"commands_run" integer DEFAULT 0 NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "command_logs" (
	"id" serial PRIMARY KEY NOT NULL,
	"guild_id" varchar(20),
	"user_id" varchar(20) NOT NULL,
	"command_name" varchar(50) NOT NULL,
	"executed_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "guild_configs" (
	"guild_id" varchar(20) PRIMARY KEY NOT NULL,
	"prefix" varchar(10) DEFAULT '!' NOT NULL,
	"log_channel_id" varchar(20),
	"premium" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE IF NOT EXISTS "user_economy" (
	"user_id" varchar(20) PRIMARY KEY NOT NULL,
	"balance" integer DEFAULT 0 NOT NULL,
	"last_daily_claim" timestamp
);
