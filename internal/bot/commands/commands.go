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
					MaxLength:   ptrInt(180),
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
					MaxLength:   ptrInt(180),
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
					MaxLength:   ptrInt(120),
				},
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "Announcement body text",
					Required:    true,
					MaxLength:   ptrInt(1800),
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
			Name:        "tip",
			Description: "Get a random helpful tip about using PulseKeep",
		},
		discord.SlashCommandCreate{
			Name:        "vote",
			Description: "Get the link to vote for PulseKeep on Top.gg",
		},
		discord.SlashCommandCreate{
			Name:        "balance",
			Description: "Check a PulseKeep wallet balance and quick economy stats",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to check the balance of, defaults to yourself",
					Required:    false,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the balance publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "profile",
			Description: "View a full PulseKeep economy profile and activity record",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to inspect, defaults to yourself",
					Required:    false,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the profile publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "daily",
			Description: "Claim your daily Pulses reward",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the reward publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "work",
			Description: "Work a shift to earn some Pulses",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show your earnings publicly instead of only to you",
					Required:    false,
				},
			},
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
					MaxValue:    ptrInt(50000),
				},
				discord.ApplicationCommandOptionString{
					Name:        "side",
					Description: "Pick heads or tails",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{Name: "Heads", Value: "heads"},
						{Name: "Tails", Value: "tails"},
					},
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
					MaxValue:    ptrInt(250000),
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
			Description: "Browse the PulseKeep item shop with prices and effects",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the shop publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "buy",
			Description: "Purchase an item from the PulseKeep shop",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "item",
					Description: "The item to purchase",
					Required:    true,
					Choices:     itemChoices(),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "inventory",
			Description: "View your purchased items and item value",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show your inventory publicly instead of only to you",
					Required:    false,
				},
			},
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
					MaxValue:    ptrInt(50000),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "fish",
			Description: "Cast a line and catch fish to sell for Pulses",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the catch publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "mine",
			Description: "Mine for valuable ores and minerals",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the haul publicly instead of only to you",
					Required:    false,
				},
			},
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
					MaxValue:    ptrInt(50000),
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
					Choices:     itemChoices(),
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
					Choices:     itemChoices(),
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
					MinLength:   ptrInt(17),
					MaxLength:   ptrInt(20),
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
					MaxLength:   ptrInt(32),
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
			Description: "Create a poll for members to vote on",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "question",
					Description: "The poll question",
					Required:    true,
					MaxLength:   ptrInt(240),
				},
				discord.ApplicationCommandOptionString{
					Name:        "option1",
					Description: "First option",
					Required:    true,
					MaxLength:   ptrInt(80),
				},
				discord.ApplicationCommandOptionString{
					Name:        "option2",
					Description: "Second option",
					Required:    true,
					MaxLength:   ptrInt(80),
				},
				discord.ApplicationCommandOptionString{
					Name:        "option3",
					Description: "Third option (optional)",
					Required:    false,
					MaxLength:   ptrInt(80),
				},
				discord.ApplicationCommandOptionString{
					Name:        "option4",
					Description: "Fourth option (optional)",
					Required:    false,
					MaxLength:   ptrInt(80),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "rich",
			Description: "Show the top 10 richest PulseKeep members with rank badges",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the leaderboard publicly instead of only to you",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "weekly",
			Description: "Claim your weekly Pulses reward (7-day cooldown)",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show the reward publicly instead of only to you",
					Required:    false,
				},
			},
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
					MaxValue:    ptrInt(50000),
				},
				discord.ApplicationCommandOptionString{
					Name:        "difficulty",
					Description: "Choose difficulty: easy, normal, hard, expert",
					Required:    false,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{Name: "Easy", Value: "easy"},
						{Name: "Normal", Value: "normal"},
						{Name: "Hard", Value: "hard"},
						{Name: "Expert", Value: "expert"},
					},
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "lottery",
			Description: "Check the weekly lottery jackpot and status",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "public",
					Description: "Show lottery status publicly instead of only to you",
					Required:    false,
				},
			},
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
					Choices:     itemChoices(),
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
		discord.SlashCommandCreate{
			Name:                     "warn",
			Description:              "Warn a user for a rule violation",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionModerateMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to warn",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "Reason for the warning",
					Required:    false,
					MaxLength:   ptrInt(500),
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "warnings",
			Description:              "View all warnings for a user",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionModerateMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to check warnings for",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "clearwarns",
			Description:              "Clear all warnings for a user",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionModerateMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to clear warnings for",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "move",
			Description:              "Move a member to a different voice channel",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionMoveMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to move",
					Required:    true,
				},
				discord.ApplicationCommandOptionChannel{
					Name:        "channel",
					Description: "The target voice channel",
					Required:    true,
					ChannelTypes: []discord.ChannelType{discord.ChannelTypeGuildVoice},
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "vckick",
			Description:              "Disconnect a member from voice chat",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionMoveMembers),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The member to disconnect from voice",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:                     "ticket",
			Description:              "Manage support tickets — add, remove, close, rename",
			DefaultMemberPermissions: ptrPermissions(discord.PermissionManageMessages),
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "add",
					Description: "Add a user to the current ticket channel",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "user",
							Description: "The user to add",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove a user from the current ticket channel",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "user",
							Description: "The user to remove",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "close",
					Description: "Close the current ticket channel",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "rename",
					Description: "Rename the current ticket channel",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "The new channel name",
							Required:    true,
							MaxLength:   ptrInt(80),
						},
					},
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "servericon",
			Description: "Show the current server's icon in full resolution",
		},
		discord.SlashCommandCreate{
			Name:        "roleinfo",
			Description: "Show detailed information about a role",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionRole{
					Name:        "role",
					Description: "The role to inspect",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "channelinfo",
			Description: "Show information about the current or a specified channel",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionChannel{
					Name:        "channel",
					Description: "The channel to inspect (defaults to current)",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "invite",
			Description: "Get invite links for PulseKeep and support server",
		},
		discord.SlashCommandCreate{
			Name:        "magic8ball",
			Description: "Ask the magic 8-ball a yes/no question",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "question",
					Description: "The question you want to ask",
					Required:    true,
					MaxLength:   ptrInt(200),
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

func itemChoices() []discord.ApplicationCommandOptionChoiceString {
	return []discord.ApplicationCommandOptionChoiceString{
		{Name: "Lucky Pickaxe", Value: "lucky_pickaxe"},
		{Name: "XP Boost", Value: "xp_boost"},
		{Name: "Golden Watch", Value: "golden_watch"},
		{Name: "Shield Token", Value: "shield_token"},
		{Name: "Lucky Clover", Value: "lucky_clover"},
		{Name: "Fishing Rod", Value: "fishing_rod"},
		{Name: "Iron Pickaxe", Value: "iron_pickaxe"},
		{Name: "Lottery Ticket", Value: "lottery_ticket"},
		{Name: "Health Potion", Value: "health_potion"},
		{Name: "Treasure Map", Value: "treasure_map"},
	}
}
