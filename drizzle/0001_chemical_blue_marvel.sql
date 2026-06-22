ALTER TABLE "guild_configs" ADD COLUMN "welcome_channel_id" varchar(20);--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "vote_channel_id" varchar(20);--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_spam_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_mention_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_link_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_caps_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_words_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD COLUMN "automod_banned_words" text;--> statement-breakpoint
ALTER TABLE "user_economy" ADD COLUMN "last_search" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD COLUMN "last_vote" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD COLUMN "xp_boost_expiry" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD COLUMN "lucky_clover_active" integer DEFAULT 0 NOT NULL;