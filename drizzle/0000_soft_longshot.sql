CREATE TABLE "bot_stats" (
	"id" integer PRIMARY KEY NOT NULL,
	"guild_count" integer DEFAULT 0 NOT NULL,
	"user_count" integer DEFAULT 0 NOT NULL,
	"commands_run" integer DEFAULT 0 NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "command_logs" (
	"id" serial PRIMARY KEY NOT NULL,
	"guild_id" varchar(20),
	"user_id" varchar(20) NOT NULL,
	"command_name" varchar(50) NOT NULL,
	"executed_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "guild_configs" (
	"guild_id" varchar(20) PRIMARY KEY NOT NULL,
	"prefix" varchar(10) DEFAULT '!' NOT NULL,
	"log_channel_id" varchar(20),
	"premium" boolean DEFAULT false NOT NULL,
	"economy_enabled" boolean DEFAULT true NOT NULL,
	"tickets_enabled" boolean DEFAULT true NOT NULL,
	"modlogs_enabled" boolean DEFAULT true NOT NULL,
	"welcome_enabled" boolean DEFAULT false NOT NULL,
	"ticket_category_id" varchar(20),
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "user_economy" (
	"user_id" varchar(20) PRIMARY KEY NOT NULL,
	"balance" integer DEFAULT 0 NOT NULL,
	"last_daily_claim" timestamp,
	"last_weekly_claim" timestamp,
	"last_work" timestamp,
	"last_rob" timestamp,
	"last_fish" timestamp,
	"last_mine" timestamp,
	"total_earned" integer DEFAULT 0 NOT NULL,
	"total_gambled" integer DEFAULT 0 NOT NULL,
	"transactions" integer DEFAULT 0 NOT NULL,
	"streak" integer DEFAULT 0 NOT NULL
);
--> statement-breakpoint
CREATE TABLE "user_inventory" (
	"id" serial PRIMARY KEY NOT NULL,
	"user_id" varchar(20) NOT NULL,
	"item_id" varchar(50) NOT NULL,
	"item_name" varchar(100) NOT NULL,
	"quantity" integer DEFAULT 1 NOT NULL
);
--> statement-breakpoint
CREATE TABLE "user_warnings" (
	"id" serial PRIMARY KEY NOT NULL,
	"guild_id" varchar(20) NOT NULL,
	"user_id" varchar(20) NOT NULL,
	"moderator_id" varchar(20) NOT NULL,
	"reason" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
