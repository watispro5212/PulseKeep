package commands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/omit"
	"github.com/watispro/pulsekeep/internal/db"
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

	// ── /ping ──────────────────────────────────────────
	cmds["ping"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Check if PulseKeep is alive",
		},
		Handler: func(ctx *CommandContext) error {
			latency := ctx.Event.Client().Gateway.Latency()
			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Pong! 🏓\nGateway Latency: %s", latency),
			})
		},
	}

	// ── /uptime ────────────────────────────────────────
	cmds["uptime"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "uptime",
			Description: "Show how long PulseKeep has been running since last restart",
		},
		Handler: func(ctx *CommandContext) error {
			uptime := time.Since(startTime).Round(time.Second)
			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("PulseKeep has been up for: %s", uptime),
			})
		},
	}

	// ── /help ──────────────────────────────────────────
	cmds["help"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "help",
			Description: "Show all available PulseKeep commands",
		},
		Handler: func(ctx *CommandContext) error {
			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: "Available commands: /ping, /uptime, /help, /stats, /serverinfo, /userinfo, /avatar, /balance, /daily, /work, /pay, /kick, /ban, /purge, /setlog",
			})
		},
	}

	// ── /stats ─────────────────────────────────────────
	cmds["stats"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "stats",
			Description: "Display PulseKeep operational statistics",
		},
		Handler: func(ctx *CommandContext) error {
			servers := 0
			for range ctx.Event.Client().Caches.Guilds() {
				servers++
			}
			uptime := time.Since(startTime).Round(time.Second)

			embed := discord.NewEmbed().
				WithTitle("PulseKeep Statistics").
				AddField("Servers", fmt.Sprintf("%d", servers), true).
				AddField("Uptime", uptime.String(), true).
				AddField("Library", "Disgo", true).
				WithColor(0x7289da)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ── /serverinfo ────────────────────────────────────
	cmds["serverinfo"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "serverinfo",
			Description: "Show detailed information about the current server",
		},
		Handler: func(ctx *CommandContext) error {
			if ctx.Event.GuildID() == nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "This command can only be used in a server.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			guild, _ := ctx.Event.Guild()
			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					discord.NewEmbed().
						WithTitle("Server Information").
						WithDescription(fmt.Sprintf("Name: %s\nID: %s\nOwner ID: %s\nMembers: %d", guild.Name, guild.ID, guild.OwnerID, guild.MemberCount)).
						WithColor(0x00ff00),
				},
			})
		},
	}

	// ── /userinfo ──────────────────────────────────────
	cmds["userinfo"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "userinfo",
			Description: "Show detailed information about a user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to look up (defaults to yourself)",
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

			embed := discord.NewEmbed().
				WithTitle("User Information").
				AddField("Username", user.Username, true).
				AddField("ID", user.ID.String(), true).
				WithColor(0x0000ff)

			if user.AvatarURL() != nil {
				embed = embed.WithThumbnail(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ── /avatar ────────────────────────────────────────
	cmds["avatar"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "avatar",
			Description: "Display a user's avatar in full resolution",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user whose avatar to display (defaults to yourself)",
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

			embed := discord.NewEmbed().
				WithTitle(fmt.Sprintf("%s's Avatar", user.Username)).
				WithColor(0x00ffff)

			if user.AvatarURL() != nil {
				embed = embed.WithImage(*user.AvatarURL())
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed},
			})
		},
	}

	// ── ECONOMY COMMANDS ──────────────────────────────────

	// ── /balance ───────────────────────────────────────
	cmds["balance"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "balance",
			Description: "Check your current PulseKeep balance",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to check the balance of (defaults to yourself)",
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
					Content: "Failed to fetch balance. Please try again later.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					discord.NewEmbed().
						WithTitle(fmt.Sprintf("%s's Balance", user.Username)).
						WithDescription(fmt.Sprintf("Balance: **%d Pulses** ⚛️", balance)).
						WithColor(0xf1c40f),
				},
			})
		},
	}

	// ── /daily ─────────────────────────────────────────
	cmds["daily"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "daily",
			Description: "Claim your daily Pulses reward",
		},
		Handler: func(ctx *CommandContext) error {
			reward := 100
			balance, success, err := ctx.DB.ClaimDaily(context.Background(), ctx.Event.User().ID.String(), reward)
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "Failed to claim daily reward. Please try again later.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			if !success {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "You have already claimed your daily reward today! Come back tomorrow.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You claimed your daily reward of **%d Pulses**! ⚛️\nNew Balance: **%d Pulses**", reward, balance),
			})
		},
	}

	// ── /work ──────────────────────────────────────────
	cmds["work"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:        "work",
			Description: "Work a shift to earn some Pulses",
		},
		Handler: func(ctx *CommandContext) error {
			reward := rand.Intn(50) + 10
			balance, err := ctx.DB.AddBalance(context.Background(), ctx.Event.User().ID.String(), reward)
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "Failed to process work shift.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("You worked a shift and earned **%d Pulses**! ⚛️\nNew Balance: **%d Pulses**", reward, balance),
			})
		},
	}

	// ── /pay ───────────────────────────────────────────
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
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "The amount of Pulses to send",
					Required:    true,
					MinValue:    ptrInt(1),
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			recipient := data.User("recipient")
			amount := data.Int("amount")

			if recipient.ID == ctx.Event.User().ID {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "You cannot pay yourself!",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			err := ctx.DB.Transfer(context.Background(), ctx.Event.User().ID.String(), recipient.ID.String(), amount)
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "Transaction failed. Check your balance and try again.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Successfully sent **%d Pulses** to **%s**! 💸", amount, recipient.Username),
			})
		},
	}

	// ── MODERATION COMMANDS ──────────────────────────────────

	// ── /kick ──────────────────────────────────────────
	cmds["kick"] = CommandDef{
		Create: discord.SlashCommandCreate{
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
					Description: "Reason for kicking the member",
					Required:    false,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			user := data.User("user")
			reason, _ := data.OptString("reason")
			if reason == "" {
				reason = "No reason provided"
			}

			err := ctx.Event.Client().Rest.RemoveMember(*ctx.Event.GuildID(), user.ID, rest.WithReason(reason))
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf("Failed to kick user: %v", err),
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			logAction(ctx, "Member Kicked", fmt.Sprintf("User: %s (%s)\nReason: %s\nModerator: %s", user.Tag(), user.ID, reason, ctx.Event.User().Tag()), 0xff0000)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Successfully kicked **%s**.", user.Username),
			})
		},
	}

	// ── /ban ───────────────────────────────────────────
	cmds["ban"] = CommandDef{
		Create: discord.SlashCommandCreate{
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
					Description: "Reason for banning the member",
					Required:    false,
				},
			},
		},
		Handler: func(ctx *CommandContext) error {
			data := ctx.Event.Data.(discord.SlashCommandInteractionData)
			user := data.User("user")
			reason, _ := data.OptString("reason")
			if reason == "" {
				reason = "No reason provided"
			}

			err := ctx.Event.Client().Rest.AddBan(*ctx.Event.GuildID(), user.ID, 0, rest.WithReason(reason))
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf("Failed to ban user: %v", err),
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			logAction(ctx, "Member Banned", fmt.Sprintf("User: %s (%s)\nReason: %s\nModerator: %s", user.Tag(), user.ID, reason, ctx.Event.User().Tag()), 0x8b0000)

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Successfully banned **%s**.", user.Username),
			})
		},
	}

	// ── /purge ─────────────────────────────────────────
	cmds["purge"] = CommandDef{
		Create: discord.SlashCommandCreate{
			Name:                     "purge",
			Description:              "Bulk delete messages",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
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
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf("Failed to fetch messages: %v", err),
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			messageIDs := make([]snowflake.ID, len(messages))
			for i, m := range messages {
				messageIDs[i] = m.ID
			}

			err = ctx.Event.Client().Rest.BulkDeleteMessages(ctx.Event.Channel().ID(), messageIDs)
			if err != nil {
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: fmt.Sprintf("Failed to delete messages: %v", err),
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Successfully deleted **%d** messages.", len(messageIDs)),
				Flags:   discord.MessageFlagEphemeral,
			})
		},
	}

	// ── /setlog ─────────────────────────────────────────
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
				return ctx.Event.CreateMessage(discord.MessageCreate{
					Content: "Failed to set logging channel.",
					Flags:   discord.MessageFlagEphemeral,
				})
			}

			return ctx.Event.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Logging channel has been set to <#%s>.", channel.ID),
			})
		},
	}

	return cmds
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
		WithColor(color)

	_, _ = ctx.Event.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
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

func ptrPermissions(p discord.Permissions) omit.Omit[*discord.Permissions] {
	return omit.NewPtr(p)
}

func ptrInt(i int) *int {
	return &i
}

func ptrSnowflake(id snowflake.ID) *snowflake.ID {
	return &id
}
