ALTER TABLE "guild_configs" ADD "welcome_channel_id" varchar(20);--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "vote_channel_id" varchar(20);--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_spam_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_mention_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_link_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_caps_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_words_enabled" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_banned_words" text;--> statement-breakpoint
ALTER TABLE "user_economy" ADD "last_search" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD "last_vote" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD "xp_boost_expiry" timestamp;--> statement-breakpoint
ALTER TABLE "user_economy" ADD "lucky_clover_active" integer DEFAULT 0 NOT NULL;