package commands

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
)

const (
	MenuSelectID        = "pulsekeep:menu:category"
	TicketPanelButtonID = "pulsekeep:tickets:open"
	TicketCloseButtonID = "pulsekeep:tickets:close"
	CommandMenuAccent    = 0x4f8cff
	ModerationMenuAccent = 0xfb7185
	UtilityMenuAccent    = 0x38d5c8
	EconomyMenuAccent    = 0xf5bd4f
	TicketMenuAccent     = 0x36d399
)

type CommandInfo struct {
	Name        string
	Description string
	Usage       string
}

type CommandCategory struct {
	ID          string
	Label       string
	Description string
	Color       int
	Commands    []CommandInfo
}

var Categories = []CommandCategory{
	{
		ID:          "moderation",
		Label:       "Moderation",
		Description: "Staff controls for keeping busy servers calm and accountable.",
		Color:       ModerationMenuAccent,
		Commands: []CommandInfo{
			{Name: "/purge", Description: "Bulk delete recent messages from the current channel.", Usage: "/purge amount:25"},
			{Name: "/kick", Description: "Remove a member from the server with an optional reason.", Usage: "/kick user:@member reason:Repeated spam"},
			{Name: "/ban", Description: "Ban a member from the server with an optional reason.", Usage: "/ban user:@member reason:Severe rule break"},
			{Name: "/unban", Description: "Unban a user by their ID.", Usage: "/unban user_id:123456789"},
			{Name: "/timeout", Description: "Timeout a member for a duration in minutes.", Usage: "/timeout user:@member duration:60"},
			{Name: "/nick", Description: "Change a member's nickname.", Usage: "/nick user:@member nickname:NewName"},
			{Name: "/slowmode", Description: "Set the slowmode delay in seconds.", Usage: "/slowmode seconds:5"},
			{Name: "/lock", Description: "Lock the current channel for @everyone.", Usage: "/lock"},
			{Name: "/unlock", Description: "Unlock the current channel for @everyone.", Usage: "/unlock"},
			{Name: "/announce", Description: "Send a clean embedded announcement with optional @everyone ping.", Usage: "/announce title:Update message:Patch notes are live ping:true"},
			{Name: "/role", Description: "Toggle a role on a member (add or remove).", Usage: "/role user:@member role:Staff"},
		},
	},
	{
		ID:          "utility",
		Label:       "Utility",
		Description: "Fast lookups, health checks, profile tools, and bot status commands.",
		Color:       UtilityMenuAccent,
		Commands: []CommandInfo{
			{Name: "/ping", Description: "Check whether PulseKeep is responding.", Usage: "/ping"},
			{Name: "/help", Description: "Open this interactive command browser.", Usage: "/help"},
			{Name: "/about", Description: "Show PulseKeep version, tech stack, and links.", Usage: "/about"},
			{Name: "/stats", Description: "Show operational bot statistics.", Usage: "/stats"},
			{Name: "/uptime", Description: "Show how long PulseKeep has been running.", Usage: "/uptime"},
			{Name: "/serverinfo", Description: "Show details about the current server.", Usage: "/serverinfo"},
			{Name: "/userinfo", Description: "Show details about a user, defaulting to yourself.", Usage: "/userinfo user:@member"},
			{Name: "/avatar", Description: "Display a user's avatar at full resolution.", Usage: "/avatar user:@member"},
			{Name: "/poll", Description: "Create a multi-option reaction poll.", Usage: "/poll question:\"Best color?\" option1:Red option2:Blue"},
		},
	},
	{
		ID:          "economy",
		Label:       "Economy",
		Description: "Member engagement commands for rewards, balances, and light social play.",
		Color:       EconomyMenuAccent,
		Commands: []CommandInfo{
			{Name: "/balance", Description: "Check your Pulse balance or another user's balance.", Usage: "/balance user:@member"},
			{Name: "/profile", Description: "View wallet, earned/spent totals, streaks, and coinflip record.", Usage: "/profile user:@member"},
			{Name: "/daily", Description: "Claim your daily Pulses reward.", Usage: "/daily"},
			{Name: "/work", Description: "Work a shift to earn Pulses.", Usage: "/work"},
			{Name: "/coinflip", Description: "Wager Pulses on heads or tails.", Usage: "/coinflip amount:100 side:heads"},
			{Name: "/pay", Description: "Send Pulses to another member.", Usage: "/pay recipient:@member amount:100"},
			{Name: "/rich", Description: "Show the richest members in the economy.", Usage: "/rich"},
			{Name: "/rob", Description: "Attempt to steal Pulses from another member.", Usage: "/rob user:@member"},
			{Name: "/shop", Description: "Browse the PulseKeep item shop.", Usage: "/shop"},
			{Name: "/buy", Description: "Buy an item from the PulseKeep shop.", Usage: "/buy item:lucky_pickaxe"},
			{Name: "/inventory", Description: "View items you own.", Usage: "/inventory"},
			{Name: "/slots", Description: "Spin the slot machine to win big.", Usage: "/slots amount:100"},
			{Name: "/fish", Description: "Cast a line to catch and sell fish.", Usage: "/fish"},
			{Name: "/mine", Description: "Mine for valuable ores and minerals.", Usage: "/mine"},
			{Name: "/gamble", Description: "Roll 1-100 and wager Pulses on the outcome.", Usage: "/gamble amount:100"},
			{Name: "/sell", Description: "Sell an item for a 60% refund.", Usage: "/sell item:lucky_pickaxe"},
			{Name: "/use", Description: "Use a usable item from your inventory.", Usage: "/use item:shield_token"},
			{Name: "/weekly", Description: "Claim your weekly reward (7-day cooldown).", Usage: "/weekly"},
			{Name: "/gift", Description: "Give an item to another user.", Usage: "/gift user:@member item:fishing_rod"},
		},
	},
	{
		ID:          "tickets",
		Label:       "Tickets",
		Description: "Support panel flows for private help channels and staff handoff.",
		Color:       TicketMenuAccent,
		Commands: []CommandInfo{
			{Name: "/ticketpanel", Description: "Post the interactive ticket opener panel.", Usage: "/ticketpanel"},
			{Name: "Open Ticket", Description: "Button flow for users to request help from staff.", Usage: "Click Open Ticket in the support panel"},
		},
	},
}

func CategoryByID(id string) (CommandCategory, bool) {
	for _, category := range Categories {
		if category.ID == id {
			return category, true
		}
	}
	return CommandCategory{}, false
}

func MenuMessage(selectedCategoryID string, ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(ephemeral).
		AddEmbeds(menuEmbed(selectedCategoryID)).
		AddActionRow(categorySelect(selectedCategoryID)).
		AddActionRow(discord.NewSuccessButton("Open Ticket", TicketPanelButtonID))
}

func MenuUpdate(selectedCategoryID string) discord.MessageUpdate {
	return discord.NewMessageUpdate().
		WithEmbeds(menuEmbed(selectedCategoryID)).
		AddActionRow(categorySelect(selectedCategoryID)).
		AddActionRow(discord.NewSuccessButton("Open Ticket", TicketPanelButtonID))
}

func TicketPanelMessage(ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(ephemeral).
		AddEmbeds(discord.NewEmbed().
			WithTitle("PulseKeep Support Tickets").
			WithDescription("Use this panel when you need setup help, billing questions, bug triage, or private support. A staff member can pick up the request from the support queue.").
			WithColor(TicketMenuAccent).
			AddField("Best for", "Setup help, deploy problems, moderation appeals, billing questions, and bug reports.", false).
			AddField("Before opening", "Share the affected server, the command or feature involved, and any error text you saw.", false).
			WithFooterText("PulseKeep Support Tickets")).
		AddActionRow(discord.NewSuccessButton("Open Ticket", TicketPanelButtonID))
}

func menuEmbed(selectedCategoryID string) discord.Embed {
	if selectedCategoryID == "" || selectedCategoryID == "overview" {
		return overviewEmbed()
	}

	category, ok := CategoryByID(selectedCategoryID)
	if !ok {
		return overviewEmbed()
	}

	embed := discord.NewEmbed().
		WithTitle(fmt.Sprintf("PulseKeep Commands - %s", category.Label)).
		WithDescription(category.Description).
		WithColor(category.Color).
		WithFooterText("Use the menu below to switch command groups.")

	for _, command := range category.Commands {
		embed = embed.AddField(command.Name, fmt.Sprintf("%s\n`%s`", command.Description, command.Usage), false)
	}

	return embed
}

func overviewEmbed() discord.Embed {
	var summary strings.Builder
	for _, category := range Categories {
		summary.WriteString(fmt.Sprintf("**%s** - %d commands\n%s\n\n", category.Label, len(category.Commands), category.Description))
	}

	return discord.NewEmbed().
		WithTitle("PulseKeep Interactive Menu").
		WithDescription(strings.TrimSpace(summary.String())).
		WithColor(CommandMenuAccent).
		AddField("How to use it", "Pick a category from the dropdown or press Open Ticket for support. Staff-only actions still rely on Discord permissions.", false).
		WithFooterText("PulseKeep command browser")
}

func categorySelect(selectedCategoryID string) discord.StringSelectMenuComponent {
	options := make([]discord.StringSelectMenuOption, 0, len(Categories)+1)
	options = append(options, discord.NewStringSelectMenuOption("Overview", "overview").
		WithDescription("Start page and command group summary").
		WithDefault(selectedCategoryID == "" || selectedCategoryID == "overview"))

	for _, category := range Categories {
		options = append(options, discord.NewStringSelectMenuOption(category.Label, category.ID).
			WithDescription(category.Description).
			WithDefault(category.ID == selectedCategoryID))
	}

	return discord.NewStringSelectMenu(MenuSelectID, "Choose a PulseKeep command category", options...)
}
