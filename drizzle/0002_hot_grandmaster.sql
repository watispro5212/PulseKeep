ALTER TABLE "guild_configs" ADD "automod_spam_limit" integer DEFAULT 5;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_spam_window" integer DEFAULT 5000;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_mention_limit" integer DEFAULT 5;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_caps_ratio" integer DEFAULT 70;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_caps_min_length" integer DEFAULT 10;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_timeout_duration" integer DEFAULT 10;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_spam_action" varchar(10) DEFAULT 'warn';--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_mention_action" varchar(10) DEFAULT 'timeout';--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_caps_action" varchar(10) DEFAULT 'delete';--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_link_action" varchar(10) DEFAULT 'delete';--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_words_action" varchar(10) DEFAULT 'delete';--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_exempt_roles" text;--> statement-breakpoint
ALTER TABLE "guild_configs" ADD "automod_allowed_domains" text;