package commands

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
<<<<<<< HEAD
=======
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
	"github.com/disgoorg/omit"
	"github.com/watispro/pulsekeep/internal/db"
)

<<<<<<< HEAD
// Register returns all slash commands PulseKeep exposes to Discord.
func Register() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
=======
// Embed colors
const (
	ColorPrimary   = 0x7058ff
	ColorSuccess   = 0x10b981
	ColorWarning   = 0xf59e0b
	ColorDanger    = 0xef4444
	ColorInfo      = 0x00f2fe
	ColorGold      = 0xffd700
	ColorPurple    = 0x9b59b6
)

var startTime = time.Now()

// CommandContext holds the data passed to every command handler.
type CommandContext struct {
	Event *events.ApplicationCommandInteractionCreate
	DB    *db.Database
}

// CommandHandler is the function signature every command must implement.
type CommandHandler func(ctx *CommandContext) error

// CommandDef wraps a slash command definition with its handler.
type CommandDef struct {
	Create  discord.ApplicationCommandCreate
	Handler CommandHandler
}

// GetCommands returns all slash command definitions.
func GetCommands(database *db.Database) map[string]CommandDef {
	cmds := make(map[string]CommandDef)

	// ═══════════════════════════════════════════════════════════════
	// GENERAL COMMANDS
	// ═══════════════════════════════════════════════════════════════

	// PING
	cmds["ping"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:        "ping",
			Description: "Check PulseKeep's response time and status",
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
			Name:        "help",
			Description: "Open the interactive PulseKeep command browser",
		},
		discord.SlashCommandCreate{
			Name:        "menu",
			Description: "Open the PulseKeep interactive command menu",
		},
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			latency := ctx.Event.Client().Gateway.Latency()
			embed := discord.NewEmbed().
				WithTitle("🏓 Pong!").
				WithDescription("PulseKeep is alive and kicking!").
				AddField("Gateway Latency", fmt.Sprintf("`%s`", latency), true).
				AddField("Status", "🟢 Online", true).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// UPTIME
	cmds["uptime"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "uptime",
			Description: "Check how long PulseKeep has been running",
		},
		Handler: func(ctx *CommandContext) error {
			uptime := time.Since(startTime)
			days := int(uptime.Hours() / 24)
			hours := int(uptime.Hours()) % 24
			minutes := int(uptime.Minutes()) % 60
			seconds := int(uptime.Seconds()) % 60

			uptimeStr := fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)

			embed := discord.NewEmbed().
				WithTitle("⏱️ Bot Uptime").
				WithDescription("PulseKeep has been running continuously.").
				AddField("Current Uptime", fmt.Sprintf("`%s`", uptimeStr), true).
				AddField("Started", fmt.Sprintf("<t:%d:R>", startTime.Unix()), true).
				WithColor(ColorInfo)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// STATS
	cmds["stats"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:        "stats",
			Description: "View PulseKeep's operational statistics",
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			servers := 0
			for range ctx.Event.Client().Caches.Guilds() {
				servers++
			}
			uptime := time.Since(startTime)
			latency := ctx.Event.Client().Gateway.Latency()

			embed := discord.NewEmbed().
				WithTitle("📊 PulseKeep Statistics").
				WithColor(ColorPrimary).
				AddField("🖥️ Servers", fmt.Sprintf("**%d**", servers), true).
				AddField("⏱️ Uptime", formatDuration(uptime), true).
				AddField("📡 Latency", fmt.Sprintf("`%s`", latency), true).
				AddField("⚡ Library", "Disgo v14", true).
				WithFooter("Powered by Go 1.26 + Gin", "")

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// HELP
	cmds["help"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "help",
			Description: "Get help with PulseKeep commands",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "category",
					Description: "Command category to view",
					Required:    false,
					Choices: []discord.ApplicationCommandOptionChoice{
						{Name: "General", Value: "general"},
						{Name: "Economy", Value: "economy"},
						{Name: "Moderation", Value: "moderation"},
					},
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			category, _ := data.OptString("category")

			var embed discord.Embed

			switch category {
			case "economy":
				embed = discord.NewEmbed().
					WithTitle("💰 Economy Commands").
					WithDescription("Manage your Pulses and play economy games!").
					WithColor(ColorGold).
					AddField("/balance [user]", "Check your or another user's balance", true).
					AddField("/daily", "Claim your daily reward", true).
					AddField("/work", "Work to earn Pulses", true).
					AddField("/gamble <amount>", "Try your luck at gambling", true).
					AddField("/rob <user>", "Attempt to rob another user", true).
					AddField("/pay <user> <amount>", "Send Pulses to another user", true).
					AddField("/leaderboard", "View the richest users", true).
					AddField("/profile [user]", "View economy profile", true).
					AddField("/shop", "Browse the item shop", true).
					AddField("/inventory [user]", "View your items", true)
			case "moderation":
				embed = discord.NewEmbed().
					WithTitle("🛡️ Moderation Commands").
					WithDescription("Keep your server safe and organized!").
					WithColor(ColorDanger).
					AddField("/kick <user> [reason]", "Kick a member from the server", true).
					AddField("/ban <user> [reason]", "Ban a member from the server", true).
					AddField("/purge <amount>", "Delete multiple messages", true).
					AddField("/setlog <channel>", "Set the logging channel", true)
			default:
				embed = discord.NewEmbed().
					WithTitle("📚 PulseKeep Help").
					WithDescription("Use the category option to filter commands.").
					WithColor(ColorPrimary).
					AddField("📌 General", "`/ping` `/uptime` `/stats` `/help`", false).
					AddField("💰 Economy", "`/balance` `/daily` `/work` `/gamble` `/pay` `/leaderboard` `/shop`", false).
					AddField("🛡️ Moderation", "`/kick` `/ban` `/purge` `/setlog`", false).
					AddField("🔍 Utility", "`/serverinfo` `/userinfo` `/avatar`", false)
			}

			embed = embed.WithFooter("Use /help category:<name> to filter", "")

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ═══════════════════════════════════════════════════════════════
	// UTILITY COMMANDS
	// ═══════════════════════════════════════════════════════════════

	cmds["serverinfo"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:        "serverinfo",
			Description: "View detailed information about the current server",
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			if ctx.Event.GuildID() == nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ This command can only be used in a server.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			guild, _ := ctx.Event.Guild()
			if guild == nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ Could not fetch server information.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			createdAt := time.Unix(int64(guild.ID.Time()), 0)

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("🏰 %s", guild.Name)).
				WithColor(ColorInfo).
				AddField("👑 Owner", fmt.Sprintf("<@!%s>", guild.OwnerID), true).
				AddField("👥 Members", fmt.Sprintf("**%d**", guild.MemberCount), true).
				AddField("📅 Created", fmt.Sprintf("<t:%d>", createdAt.Unix()), true).
				AddField("🆔 Server ID", guild.ID.String(), true).
				WithFooter("PulseKeep Server Info", "")

			if guild.IconURL() != "" {
				embed = embed.WithThumbnail(guild.IconURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["userinfo"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:        "userinfo",
			Description: "View detailed information about a user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to look up, defaults to yourself",
					Required:    false,
				},
			},
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			user := ctx.Event.User()
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			if u, ok := data.OptUser("user"); ok {
				user = u
			}

			createdAt := time.Unix(int64(user.ID.Time()), 0)

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("👤 %s", user.Username)).
				WithColor(ColorPrimary).
				AddField("🏷️ Tag", user.Username+"#"+user.Discriminator, true).
				AddField("🆔 ID", user.ID.String(), true).
				AddField("📅 Joined Discord", fmt.Sprintf("<t:%d>", createdAt.Unix()), true).
				WithFooter("Requested by "+ctx.Event.User().Username, "")

			if user.AvatarURL() != nil {
				embed = embed.WithThumbnail(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["avatar"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:        "avatar",
			Description: "Display a user's avatar in full resolution",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user whose avatar to display, defaults to yourself",
					Required:    false,
				},
			},
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
			Name:                     "purge",
			Description:              "Bulk delete messages from the current channel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "Number of messages to delete, from 1 to 100",
					Required:    true,
					MinValue:    ptrInt(1),
					MaxValue:    ptrInt(100),
				},
			},
		},
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			user := ctx.Event.User()
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			if u, ok := data.OptUser("user"); ok {
				user = u
			}

			avatarURL := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=512", user.ID, user.AvatarID())

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("%s's Avatar", user.Username)).
				WithDescription(fmt.Sprintf("[Click here to download](%s)", avatarURL)).
				WithColor(ColorInfo)

			if user.AvatarURL() != nil {
				embed = embed.WithImage(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ═══════════════════════════════════════════════════════════════
	// ECONOMY COMMANDS
	// ═══════════════════════════════════════════════════════════════

	cmds["balance"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "balance",
			Description: "Check your current PulseKeep balance",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to check balance for (defaults to yourself)",
					Required:    false,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			user := ctx.Event.User()
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			if u, ok := data.OptUser("user"); ok {
				user = u
			}

			balance, err := ctx.DB.GetBalance(context.Background(), user.ID.String())
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ Failed to fetch balance. Please try again later.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			rank, _ := ctx.DB.GetUserRank(context.Background(), user.ID.String())

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("💰 %s's Balance", user.Username)).
				WithColor(ColorGold).
				AddField("⚛️ Balance", fmt.Sprintf("**%,d** Pulses", balance), true).
				AddField("🏆 Rank", fmt.Sprintf("#%d", rank), true).
				WithFooter("PulseKeep Economy", "")

			if user.AvatarURL() != nil {
				embed = embed.WithThumbnail(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["daily"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "daily",
			Description: "Claim your daily Pulses reward",
		},
		Handler: func(ctx *CommandContext) error {
			reward := 100 + rand.Intn(51) // 100-150 Pulses

			balance, success, err := ctx.DB.ClaimDaily(context.Background(), ctx.Event.User().ID.String(), reward)
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ Failed to claim daily reward. Please try again later.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			if !success {
				nextClaim, _ := ctx.DB.GetNextDailyTime(context.Background(), ctx.Event.User().ID.String())
				timeUntil := time.Until(nextClaim)
				hours := int(timeUntil.Hours())
				minutes := int(timeUntil.Minutes()) % 60

				embed := discord.NewEmbed().
					WithTitle("⏰ Daily Already Claimed").
					WithDescription("You've already claimed your daily reward today!").
					AddField("Next Claim", fmt.Sprintf("in **%dh %dm**", hours, minutes), true).
					WithColor(ColorWarning)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			embed := discord.NewEmbed().
				WithTitle("🎉 Daily Reward Claimed!").
				WithDescription(fmt.Sprintf("You claimed your daily reward of **%,d Pulses** ⚛️", reward)).
				AddField("💰 New Balance", fmt.Sprintf("**%,d** Pulses", balance), true).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["work"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "work",
			Description: "Work a shift to earn some Pulses",
		},
		Handler: func(ctx *CommandContext) error {
			canWork, remaining := ctx.DB.CanWork(context.Background(), ctx.Event.User().ID.String())
			if !canWork {
				minutes := int(remaining.Minutes())
				embed := discord.NewEmbed().
					WithTitle("⏰ Work Cooldown").
					WithDescription("You're still on a work cooldown.").
					AddField("Next Shift", fmt.Sprintf("in **%d minutes**", minutes), true).
					WithColor(ColorWarning)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			jobs := []struct {
				name   string
				reward int
			}{
				{"Developer", 50 + rand.Intn(101)},
				{"Designer", 40 + rand.Intn(81)},
				{"Writer", 30 + rand.Intn(71)},
				{"Moderator", 60 + rand.Intn(91)},
				{"DJ", 45 + rand.Intn(86)},
				{"Streamer", 55 + rand.Intn(96)},
			}

			job := jobs[rand.Intn(len(jobs))]
			ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), job.reward)
			ctx.DB.SetWorkCooldown(context.Background(), ctx.Event.User().ID.String())
			balance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("💼 Worked as %s", job.name)).
				WithDescription(fmt.Sprintf("You completed your shift and earned **%,d Pulses**!", job.reward)).
				AddField("💰 New Balance", fmt.Sprintf("**%,d** Pulses", balance), true).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["gamble"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "gamble",
			Description: "Gamble your Pulses for a chance to win big!",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInteger{
					Name:        "amount",
					Description: "Amount of Pulses to gamble (10-10000)",
					Required:    true,
					MinValue:    ptrInt(10),
					MaxValue:    ptrInt(10000),
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			amount := int(data.Int("amount"))

			balance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())
			if balance < amount {
				embed := discord.NewEmbed().
					WithTitle("❌ Insufficient Funds").
					WithDescription(fmt.Sprintf("You need **%,d** Pulses to gamble, but you only have **%,d**.", amount, balance)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			winChance := rand.Float64()
			var result string
			var color int

			if winChance < 0.45 {
				// Lose
				ctx.DB.RemoveBalance(context.Background(), ctx.Event.User().ID.String(), amount)
				newBalance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())
				result = fmt.Sprintf("💸 You lost **%,d** Pulses! Better luck next time.", amount)
				color = ColorDanger
				_ = newBalance
			} else if winChance < 0.7 {
				// Small win (1.5x)
				won := amount
				ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), won)
				result = fmt.Sprintf("🎉 You won **%,d** Pulses! (1.5x multiplier)", won)
				color = ColorSuccess
			} else if winChance < 0.9 {
				// Medium win (2x)
				won := amount * 2
				ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), won)
				result = fmt.Sprintf("🎊 JACKPOT! You won **%,d** Pulses! (2x multiplier)", won)
				color = ColorGold
			} else {
				// Big win (3x)
				won := amount * 3
				ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), won)
				result = fmt.Sprintf("🌟 MEGA JACKPOT! You won **%,d** Pulses! (3x multiplier)", won)
				color = ColorPurple
			}

			newBalance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())

			embed := discord.NewEmbed().
				WithTitle("🎰 Gambling Results").
				WithDescription(result).
				AddField("💰 New Balance", fmt.Sprintf("**%,d** Pulses", newBalance), true).
				WithColor(color)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["rob"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "rob",
			Description: "Attempt to rob another user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to rob",
					Required:    true,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			victim := data.User("user")

			if victim.ID == ctx.Event.User().ID {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ You can't rob yourself!",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			canRob, remaining := ctx.DB.CanRob(context.Background(), ctx.Event.User().ID.String())
			if !canRob {
				minutes := int(remaining.Minutes())
				embed := discord.NewEmbed().
					WithTitle("⏰ Rob Cooldown").
					WithDescription("You're still on a rob cooldown.").
					AddField("Next Attempt", fmt.Sprintf("in **%d minutes**", minutes), true).
					WithColor(ColorWarning)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			victimBalance, _ := ctx.DB.GetBalance(context.Background(), victim.ID.String())
			if victimBalance < 50 {
				embed := discord.NewEmbed().
					WithTitle("💸 Poor Target").
					WithDescription(fmt.Sprintf("%s only has **%,d** Pulses. Find a richer target!", victim.Username, victimBalance)).
					WithColor(ColorWarning)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			success := rand.Float64() < 0.35

			if success {
				stealPercent := 10 + rand.Intn(41)
				stealAmount := (victimBalance * stealPercent) / 100
				if stealAmount < 10 {
					stealAmount = 10
				}
				if stealAmount > 5000 {
					stealAmount = 5000
				}

				ctx.DB.RemoveBalance(context.Background(), victim.ID.String(), stealAmount)
				ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), stealAmount)
				ctx.DB.SetRobCooldown(context.Background(), ctx.Event.User().ID.String())

				embed := discord.NewEmbed().
					WithTitle("🤑 Successful Heist!").
					WithDescription(fmt.Sprintf("You successfully robbed **%,d** Pulses from **%s**!", stealAmount, victim.Username)).
					WithColor(ColorSuccess)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			} else {
				fine := 50 + rand.Intn(101)
				userBalance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())
				if userBalance < fine {
					fine = userBalance / 2
				}
				if fine > 0 {
					ctx.DB.RemoveBalance(context.Background(), ctx.Event.User().ID.String(), fine)
				}
				ctx.DB.SetRobCooldown(context.Background(), ctx.Event.User().ID.String())

				embed := discord.NewEmbed().
					WithTitle("🚔 Caught!").
					WithDescription(fmt.Sprintf("You got caught and paid a fine of **%,d** Pulses!", fine)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}
		},
	}

	cmds["pay"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "pay",
			Description: "Send Pulses to another user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "recipient",
					Description: "The user to send Pulses to",
					Required:    true,
				},
				discord.ApplicationCommandOptionInteger{
					Name:        "amount",
					Description: "The amount of Pulses to send (minimum 1)",
					Required:    true,
					MinValue:    ptrInt(1),
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			recipient := data.User("recipient")
			amount := int(data.Int("amount"))

			if recipient.ID == ctx.Event.User().ID {
				embed := discord.NewEmbed().
					WithTitle("❌ Invalid Transfer").
					WithDescription("You cannot send Pulses to yourself!").
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			balance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())
			if balance < amount {
				embed := discord.NewEmbed().
					WithTitle("❌ Insufficient Funds").
					WithDescription(fmt.Sprintf("You need **%,d** Pulses but only have **%,d**.", amount, balance)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			err := ctx.DB.Transfer(context.Background(), ctx.Event.User().ID.String(), recipient.ID.String(), amount)
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Transfer Failed").
					WithDescription("Transaction failed. Please try again.").
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			embed := discord.NewEmbed().
				WithTitle("💸 Transfer Successful").
				WithDescription(fmt.Sprintf("You sent **%,d** Pulses to **%s**!", amount, recipient.Username)).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["leaderboard"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "leaderboard",
			Description: "View the richest users on PulseKeep",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInteger{
					Name:        "page",
					Description: "Page number (default 1)",
					Required:    false,
					MinValue:    ptrInt(1),
					MaxValue:    ptrInt(10),
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			page, _ := data.OptInt("page")
			if page < 1 {
				page = 1
			}

			leaders, err := ctx.DB.GetLeaderboard(context.Background(), 10, int(page))
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ Failed to fetch leaderboard.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			if len(leaders) == 0 {
				embed := discord.NewEmbed().
					WithTitle("🏆 PulseKeep Leaderboard").
					WithDescription("No users yet! Be the first to earn Pulses.").
					WithColor(ColorGold)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			var description strings.Builder
			for i, entry := range leaders {
				rank := (page-1)*10 + i + 1
				medal := ""
				switch rank {
				case 1:
					medal = "🥇"
				case 2:
					medal = "🥈"
				case 3:
					medal = "🥉"
				default:
					medal = fmt.Sprintf("**%d.**", rank)
				}
				description.WriteString(fmt.Sprintf("%s <@!%s>: **%,d** Pulses\n", medal, entry.UserID, entry.Balance))
			}

			embed := discord.NewEmbed().
				WithTitle("🏆 PulseKeep Leaderboard").
				WithDescription(description.String()).
				AddField("📄 Page", fmt.Sprintf("%d", page), true).
				WithColor(ColorGold).
				WithFooter("Use /leaderboard page:<number> to see more", "")

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["profile"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "profile",
			Description: "View your or another user's economy profile",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to view profile for (defaults to yourself)",
					Required:    false,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			user := ctx.Event.User()
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			if u, ok := data.OptUser("user"); ok {
				user = u
			}

			balance, _ := ctx.DB.GetBalance(context.Background(), user.ID.String())
			rank, _ := ctx.DB.GetUserRank(context.Background(), user.ID.String())
			stats, _ := ctx.DB.GetUserStats(context.Background(), user.ID.String())

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("👤 %s's Profile", user.Username)).
				WithColor(ColorPrimary).
				AddField("💰 Balance", fmt.Sprintf("**%,d** Pulses", balance), true).
				AddField("🏆 Rank", fmt.Sprintf("#%d", rank), true).
				AddField("📊 Total Earned", fmt.Sprintf("**%,d**", stats.TotalEarned), true).
				AddField("🎰 Total Gambled", fmt.Sprintf("**%,d**", stats.TotalGambled), true).
				AddField("🔄 Transactions", fmt.Sprintf("**%d**", stats.Transactions), true).
				WithFooter("PulseKeep Economy", "")

			if user.AvatarURL() != nil {
				embed = embed.WithThumbnail(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["shop"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "shop",
			Description: "Browse and buy items from the PulseKeep shop",
		},
		Handler: func(ctx *CommandContext) error {
			embed := discord.NewEmbed().
				WithTitle("🛒 PulseKeep Shop").
				WithColor(ColorGold).
				AddField("🍀 Lucky Charm (500 Pulses)", "Increases gambling win chance by 10% for 24h", true).
				AddField("⭐ VIP Pass (1,000 Pulses)", "2x rewards from /work for 7 days", true).
				AddField("🛡️ Shield (300 Pulses)", "Protects from being robbed for 24h", true).
				AddField("🔢 2x Multiplier (750 Pulses)", "Doubles /daily reward for 24h", true).
				AddField("📝 Bank Note (2,000 Pulses)", "Safely store up to 10,000 Pulses offline", true).
				AddField("", "", false).
				AddField("💡 Use /buy <item>", "to purchase an item", false).
				WithFooter("PulseKeep Shop", "")

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["buy"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "buy",
			Description: "Buy an item from the shop",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item name to buy",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoice{
						{Name: "Lucky Charm (500)", Value: "lucky charm"},
						{Name: "VIP Pass (1,000)", Value: "vip pass"},
						{Name: "Shield (300)", Value: "shield"},
						{Name: "2x Multiplier (750)", Value: "2x multiplier"},
						{Name: "Bank Note (2,000)", Value: "bank note"},
					},
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			itemName := data.String("item")

			items := map[string]struct {
				price int
				emoji string
			}{
				"lucky charm":    {500, "🍀"},
				"vip pass":        {1000, "⭐"},
				"shield":          {300, "🛡️"},
				"2x multiplier":   {750, "🔢"},
				"bank note":       {2000, "📝"},
			}

			item, exists := items[itemName]
			if !exists {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "❌ Invalid item. Use /shop to see available items.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			balance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())
			if balance < item.price {
				embed := discord.NewEmbed().
					WithTitle("❌ Insufficient Funds").
					WithDescription(fmt.Sprintf("You need **%,d** Pulses but only have **%,d**.", item.price, balance)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			ctx.DB.RemoveBalance(context.Background(), ctx.Event.User().ID.String(), item.price)
			ctx.DB.AddItem(context.Background(), ctx.Event.User().ID.String(), strings.ReplaceAll(itemName, " ", "_"), itemName, 1)
			newBalance, _ := ctx.DB.GetBalance(context.Background(), ctx.Event.User().ID.String())

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("✅ Purchased %s %s", item.emoji, itemName)).
				WithDescription(fmt.Sprintf("You bought **%s** for **%,d** Pulses!", itemName, item.price)).
				AddField("💰 New Balance", fmt.Sprintf("**%,d** Pulses", newBalance), true).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["inventory"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "inventory",
			Description: "View your purchased items",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user whose inventory to view (defaults to yourself)",
					Required:    false,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			user := ctx.Event.User()
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			if u, ok := data.OptUser("user"); ok {
				user = u
			}

			items, _ := ctx.DB.GetInventory(context.Background(), user.ID.String())

			if len(items) == 0 {
				embed := discord.NewEmbed().
					WithTitle(fmt.Sprintf("🎒 %s's Inventory", user.Username)).
					WithDescription("Your inventory is empty! Use `/shop` to buy items.").
					WithColor(ColorInfo)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			var description strings.Builder
			for _, item := range items {
				emoji := getItemEmoji(item.Name)
				description.WriteString(fmt.Sprintf("%s **%s** x%d\n", emoji, item.Name, item.Quantity))
			}

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("🎒 %s's Inventory", user.Username)).
				WithDescription(description.String()).
				WithColor(ColorInfo)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ═══════════════════════════════════════════════════════════════
	// MODERATION COMMANDS
	// ═══════════════════════════════════════════════════════════════

	cmds["kick"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:                     "kick",
			Description:              "Kick a member from the server",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionKickMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to kick",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "Reason for kicking (optional)",
					Required:    false,
				},
			},
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
=======
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			user := data.User("user")
			reason, _ := data.OptString("reason")
			if reason == "" {
				reason = "No reason provided"
			}

			err := ctx.Event.Client().Rest.RemoveMember(*ctx.Event.GuildID(), user.ID, rest.WithReason(reason))
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Kick Failed").
					WithDescription(fmt.Sprintf("Failed to kick %s: %v", user.Username, err)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			logAction(ctx, "Member Kicked", fmt.Sprintf("**User:** %s (%s)\n**Reason:** %s\n**Moderator:** %s",
				user.Tag(), user.ID, reason, ctx.Event.User().Tag()), ColorDanger)

			embed := discord.NewEmbed().
				WithTitle("👢 Member Kicked").
				WithDescription(fmt.Sprintf("**%s** has been kicked from the server.", user.Username)).
				AddField("📋 Reason", reason, true).
				WithColor(ColorDanger)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["ban"] = CommandDef{
		Create: discord.SlashCommandCreate{
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
			Name:                     "ban",
			Description:              "Ban a member from the server",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionBanMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to ban",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "Reason for banning (optional)",
					Required:    false,
				},
			},
		},
<<<<<<< HEAD
		discord.SlashCommandCreate{
			Name:                     "announce",
			Description:              "Send an embedded announcement to the current channel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "title",
					Description: "Announcement title",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "Announcement body text",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "uptime",
			Description: "Show how long PulseKeep has been running since last restart",
		},
		discord.SlashCommandCreate{
			Name:        "balance",
			Description: "Check your current PulseKeep balance",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to check the balance of, defaults to yourself",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "profile",
			Description: "View a PulseKeep economy profile",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to inspect, defaults to yourself",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "daily",
			Description: "Claim your daily Pulses reward",
		},
		discord.SlashCommandCreate{
			Name:        "work",
			Description: "Work a shift to earn some Pulses",
		},
		discord.SlashCommandCreate{
			Name:        "coinflip",
			Description: "Wager Pulses on a heads-or-tails coinflip",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "How many Pulses to wager",
					Required:    true,
					MinValue:    ptrInt(1),
				},
				discord.ApplicationCommandOptionString{
					Name:        "side",
					Description: "Pick heads or tails",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "pay",
			Description: "Send Pulses to another user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "recipient",
					Description: "The user to send Pulses to",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "The amount of Pulses to send",
					Required:    true,
					MinValue:    ptrInt(1),
				},
			},
=======
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			user := data.User("user")
			reason, _ := data.OptString("reason")
			if reason == "" {
				reason = "No reason provided"
			}

			err := ctx.Event.Client().Rest.AddBan(*ctx.Event.GuildID(), user.ID, 0, rest.WithReason(reason))
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Ban Failed").
					WithDescription(fmt.Sprintf("Failed to ban %s: %v", user.Username, err)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			logAction(ctx, "Member Banned", fmt.Sprintf("**User:** %s (%s)\n**Reason:** %s\n**Moderator:** %s",
				user.Tag(), user.ID, reason, ctx.Event.User().Tag()), ColorDanger)

			embed := discord.NewEmbed().
				WithTitle("🔨 Member Banned").
				WithDescription(fmt.Sprintf("**%s** has been banned from the server.", user.Username)).
				AddField("📋 Reason", reason, true).
				WithColor(ColorDanger)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
		},
		discord.SlashCommandCreate{
			Name:        "leaderboard",
			Description: "Show the richest PulseKeep economy members",
		},
		discord.SlashCommandCreate{
			Name:                     "ticketpanel",
			Description:              "Post the interactive PulseKeep ticket panel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
		},
	}

	cmds["purge"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:                     "purge",
			Description:              "Bulk delete messages in the current channel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInteger{
					Name:        "amount",
					Description: "Number of messages to delete (1-100)",
					Required:    true,
					MinValue:    ptrInt(1),
					MaxValue:    ptrInt(100),
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			amount := ctx.Event.Data.(discord.SlashCommandInteractionData).Int("amount")

			messages, err := ctx.Event.Client().Rest.GetMessages(ctx.Event.Channel().ID(), 0, 0, 0, amount)
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Purge Failed").
					WithDescription(fmt.Sprintf("Failed to fetch messages: %v", err)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			messageIDs := make([]snowflake.ID, len(messages))
			for i, m := range messages {
				messageIDs[i] = m.ID
			}

			err = ctx.Event.Client().Rest.BulkDeleteMessages(ctx.Event.Channel().ID(), messageIDs)
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Purge Failed").
					WithDescription(fmt.Sprintf("Failed to delete messages: %v", err)).
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			embed := discord.NewEmbed().
				WithTitle("🗑️ Messages Purged").
				WithDescription(fmt.Sprintf("Successfully deleted **%d** messages.", len(messageIDs))).
				AddField("Moderator", ctx.Event.User().Username, true).
				WithColor(ColorWarning)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	cmds["setlog"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:                     "setlog",
			Description:              "Set the logging channel for this server",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageGuild),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionChannel{
					Name:        "channel",
					Description: "The channel to send logs to",
					Required:    true,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			channel := data.Channel("channel")

			err := ctx.DB.SetLogChannel(context.Background(), ctx.Event.GuildID().String(), channel.ID.String())
			if err != nil {
				embed := discord.NewEmbed().
					WithTitle("❌ Setup Failed").
					WithDescription("Failed to set logging channel. Please try again.").
					WithColor(ColorDanger)

				return ctx.Event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{embed},
				})
			}

			embed := discord.NewEmbed().
				WithTitle("✅ Logging Channel Set").
				WithDescription(fmt.Sprintf("All moderation logs will now be sent to <#%s>.", channel.ID)).
				WithColor(ColorSuccess)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	return cmds
}

<<<<<<< HEAD
=======
// Helper functions

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func logAction(ctx *CommandContext, title, description string, color int) {
	cfg, err := ctx.DB.GetGuildConfig(context.Background(), ctx.Event.GuildID().String())
	if err != nil || cfg.LogChannelID == nil {
		return
	}

	logChanID, _ := snowflake.Parse(*cfg.LogChannelID)
	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithColor(color).
		WithFooter("Moderation Action", "")

	_, _ = ctx.Event.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
}

func getItemEmoji(name string) string {
	emojis := map[string]string{
		"lucky charm":    "🍀",
		"vip pass":        "⭐",
		"shield":          "🛡️",
		"2x multiplier":   "🔢",
		"bank note":       "📝",
	}

	lower := strings.ToLower(name)
	if emoji, ok := emojis[lower]; ok {
		return emoji
	}
	return "🎁"
}

// Register returns all slash commands PulseKeep exposes to Discord.
func Register(database *db.Database) []discord.ApplicationCommandCreate {
	commandsMap := GetCommands(database)
	creates := make([]discord.ApplicationCommandCreate, 0, len(commandsMap))
	for _, cmd := range commandsMap {
		creates = append(creates, cmd.Create)
	}
	return creates
}

>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
func ptrPermissions(p discord.Permissions) omit.Omit[*discord.Permissions] {
	return omit.NewPtr(p)
}

func ptrInt(i int) *int {
	return &i
}
<<<<<<< HEAD
=======

// LeaderboardEntry represents a user on the leaderboard
type LeaderboardEntry struct {
	UserID  string
	Balance int
}

// UserStats represents a user's economy statistics
type UserStats struct {
	TotalEarned  int
	TotalGambled int
	Transactions int
}

// InventoryItem represents an item in a user's inventory
type InventoryItem struct {
	Name     string
	Quantity int
}
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
