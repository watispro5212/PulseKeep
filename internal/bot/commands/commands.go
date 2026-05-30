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
				discord.ApplicationCommandOptionBool{
					Name:        "ping",
					Description: "Whether to ping @everyone",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "about",
			Description: "Show PulseKeep version, tech stack, and links",
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
			Name:        "fish",
			Description: "Cast a line and catch fish to sell for Pulses",
		},
		discord.SlashCommandCreate{
			Name:        "mine",
			Description: "Mine for valuable ores and minerals",
		},
		discord.SlashCommandCreate{
			Name:        "gamble",
			Description: "Roll a dice (1-100) and wager your Pulses",
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
			Name:        "sell",
			Description: "Sell an item from your inventory for a 60% refund",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item to sell",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "use",
			Description: "Use a usable item from your inventory",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item to use",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "unban",
			Description:              "Unban a user from the server",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionBanMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "user_id",
					Description: "The ID of the user to unban",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "slowmode",
			Description:              "Set the slowmode delay in the current channel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageChannels),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "seconds",
					Description: "Slowmode delay in seconds (0 to disable, max 21600)",
					Required:    true,
					MinValue:    ptrInt(0),
					MaxValue:    ptrInt(21600),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "nick",
			Description:              "Change a member's nickname",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageNicknames),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to rename",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "nickname",
					Description: "The new nickname (leave empty to reset)",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "timeout",
			Description:              "Timeout a member for a specified duration",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionModerateMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to timeout",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "duration",
					Description: "Timeout duration in minutes (max 40320)",
					Required:    true,
					MinValue:    ptrInt(1),
					MaxValue:    ptrInt(40320),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "lock",
			Description:              "Lock the current channel to prevent @everyone from sending messages",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageChannels),
		},
		discord.SlashCommandCreate{
			Name:                     "unlock",
			Description:              "Unlock the current channel to allow @everyone to send messages",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageChannels),
		},
		discord.SlashCommandCreate{
			Name:                     "ticketpanel",
			Description:              "Post the interactive PulseKeep ticket panel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
		},
		discord.SlashCommandCreate{
			Name:        "poll",
			Description:  "Create a poll for members to vote on",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "question",
					Description: "The poll question",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "option1",
					Description: "First option",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "option2",
					Description: "Second option",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "option3",
					Description: "Third option (optional)",
					Required:    false,
				},
				discord.ApplicationCommandOptionString{
					Name:        "option4",
					Description: "Fourth option (optional)",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "rich",
			Description: "Show the top 10 richest PulseKeep members with rank badges",
		},
		discord.SlashCommandCreate{
			Name:        "weekly",
			Description: "Claim your weekly Pulses reward (7-day cooldown)",
		},
		discord.SlashCommandCreate{
			Name:        "blackjack",
			Description: "Play blackjack against the CPU and wager your Pulses",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "How many Pulses to wager (min 100)",
					Required:    true,
					MinValue:    ptrInt(100),
				},
				discord.ApplicationCommandOptionString{
					Name:        "difficulty",
					Description: "Choose difficulty: easy, normal, hard, expert",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "lottery",
			Description: "Check the weekly lottery jackpot and status",
		},
		discord.SlashCommandCreate{
			Name:        "lottery-claim",
			Description: "Claim your prize if you won the weekly lottery draw",
		},
		discord.SlashCommandCreate{
			Name:        "gift",
			Description: "Give an item from your inventory to another user",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to give the item to",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item ID to give (use /inventory to see your items)",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "role",
			Description:              "Add or remove a role from a member",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageRoles),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to update",
					Required:    true,
				},
				discord.ApplicationCommandOptionRole{
					Name:        "role",
					Description: "The role to add or remove",
					Required:    true,
				},
			},
		},
	}
}

func ptrPermissions(p discord.Permissions) omit.Omit[*discord.Permissions] {
	return omit.NewPtr(p)
}

func ptrInt(i int) *int {
	return &i
}
