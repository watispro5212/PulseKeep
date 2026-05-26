package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/omit"
)

// Register returns all slash commands PulseKeep exposes to Discord.
func Register() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Check if PulseKeep is alive and measure gateway latency",
		},
		discord.SlashCommandCreate{
			Name:        "help",
			Description: "Open the interactive PulseKeep command browser",
		},
		discord.SlashCommandCreate{
			Name:        "stats",
			Description: "Display PulseKeep operational statistics",
		},
		discord.SlashCommandCreate{
			Name:        "serverinfo",
			Description: "Show detailed information about the current server",
		},
		discord.SlashCommandCreate{
			Name:        "userinfo",
			Description: "Show detailed information about a user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to look up, defaults to yourself",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
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
		},
		discord.SlashCommandCreate{
			Name:        "leaderboard",
			Description: "Show the richest PulseKeep economy members",
		},
		discord.SlashCommandCreate{
			Name:        "rob",
			Description: "Attempt to rob another user for their Pulses",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to rob",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "shop",
			Description: "Browse the PulseKeep item shop",
		},
		discord.SlashCommandCreate{
			Name:        "buy",
			Description: "Purchase an item from the PulseKeep shop",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item to purchase",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "inventory",
			Description: "View your purchased items",
		},
		discord.SlashCommandCreate{
			Name:        "slots",
			Description: "Spin the slot machine and wager your Pulses",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "How many Pulses to wager",
					Required:    true,
					MinValue:    ptrInt(1),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "ticketpanel",
			Description:              "Post the interactive PulseKeep ticket panel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
		},
	}
}

func ptrPermissions(p discord.Permissions) omit.Omit[*discord.Permissions] {
	return omit.NewPtr(p)
}

func ptrInt(i int) *int {
	return &i
}
