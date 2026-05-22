package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/omit"
)

// CommandDef wraps a slash command definition with its metadata.
type CommandDef struct {
	Create  discord.ApplicationCommandCreate
	Handler CommandHandler
}

// CommandHandler is the function signature every command must implement.
type CommandHandler func(ctx *CommandContext) error

// CommandContext holds the data passed to every command handler.
type CommandContext struct {
	Event interface{} // The raw interaction event
}

// Register returns all slash commands PulseKeep exposes to Discord.
func Register() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		// ── /ping ──────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Check if PulseKeep is alive and measure gateway latency",
		},

		// ── /help ──────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "help",
			Description: "Show all available PulseKeep commands",
		},

		// ── /stats ─────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "stats",
			Description: "Display PulseKeep operational statistics",
		},

		// ── /serverinfo ────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "serverinfo",
			Description: "Show detailed information about the current server",
		},

		// ── /userinfo ──────────────────────────────────────
		discord.SlashCommandCreate{
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

		// ── /avatar ────────────────────────────────────────
		discord.SlashCommandCreate{
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

		// ── /purge ─────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:                     "purge",
			Description:              "Bulk delete messages from the current channel",
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

		// ── /kick ──────────────────────────────────────────
		discord.SlashCommandCreate{
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

		// ── /ban ───────────────────────────────────────────
		discord.SlashCommandCreate{
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

		// ── /announce ──────────────────────────────────────
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

		// ── /uptime ────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "uptime",
			Description: "Show how long PulseKeep has been running since last restart",
		},

		// ── /balance ───────────────────────────────────────
		discord.SlashCommandCreate{
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

		// ── /daily ─────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "daily",
			Description: "Claim your daily Pulses reward",
		},

		// ── /work ──────────────────────────────────────────
		discord.SlashCommandCreate{
			Name:        "work",
			Description: "Work a shift to earn some Pulses",
		},

		// ── /pay ───────────────────────────────────────────
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
		},
	}
}

// ── Helpers ───────────────────────────────────────────────

// ptrPermissions converts a discord.Permissions into a omit.Omit[*discord.Permissions].
func ptrPermissions(p discord.Permissions) omit.Omit[*discord.Permissions] {
	return omit.NewPtr(p)
}

// ptrInt converts an int to *int for option min/max values.
func ptrInt(i int) *int {
	return &i
}

// ptrSnowflake converts a snowflake.ID into a pointer — used for optional IDs.
func ptrSnowflake(id snowflake.ID) *snowflake.ID {
	return &id
}
