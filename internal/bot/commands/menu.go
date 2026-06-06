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
	CommandMenuAccent    = 0x8b5cf6
	ModerationMenuAccent = 0xf43f5e
	UtilityMenuAccent    = 0x8b5cf6
	EconomyMenuAccent    = 0x10b981
	EconomyWarningAccent = 0xf59e0b
	TicketMenuAccent     = 0x0ea5e9
	PulseKeepVersion     = "v6.3.0"
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
			{Name: "/warn", Description: "Warn a user for a rule violation.", Usage: "/warn user:@member reason:Spamming in chat"},
			{Name: "/warnings", Description: "View all warnings for a user.", Usage: "/warnings user:@member"},
			{Name: "/clearwarns", Description: "Clear all warnings for a user.", Usage: "/clearwarns user:@member"},
			{Name: "/move", Description: "Move a member to a different voice channel.", Usage: "/move user:@member channel:General"},
			{Name: "/vckick", Description: "Disconnect a member from voice chat.", Usage: "/vckick user:@member"},
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
			{Name: "/servericon", Description: "Show the server's icon in full resolution.", Usage: "/servericon"},
			{Name: "/roleinfo", Description: "Show detailed information about a role.", Usage: "/roleinfo role:Staff"},
			{Name: "/channelinfo", Description: "Show information about the current or specified channel.", Usage: "/channelinfo"},
			{Name: "/invite", Description: "Get invite links for PulseKeep and the support server.", Usage: "/invite"},
			{Name: "/poll", Description: "Create a multi-option reaction poll.", Usage: "/poll question:\"Best color?\" option1:Red option2:Blue"},
		},
	},
	{
		ID:          "economy",
		Label:       "Economy",
		Description: "Member engagement commands for rewards, balances, and light social play.",
		Color:       EconomyMenuAccent,
		Commands: []CommandInfo{
			{Name: "/balance", Description: "Check a wallet, lifetime totals, streaks, and net worth.", Usage: "/balance user:@member public:false"},
			{Name: "/profile", Description: "View wallet, net worth, streaks, blackjack, fishing, mining, and gambling stats.", Usage: "/profile user:@member public:false"},
			{Name: "/daily", Description: "Claim your daily reward with transparent streak bonus display.", Usage: "/daily"},
			{Name: "/work", Description: "Work a shift with base pay, rare bonus events, and XP Boost support.", Usage: "/work"},
			{Name: "/coinflip", Description: "Wager on heads or tails with guided side choices and Lucky Pickaxe support.", Usage: "/coinflip amount:100 side:heads"},
			{Name: "/pay", Description: "Send Pulses to another member with transfer caps.", Usage: "/pay recipient:@member amount:100"},
			{Name: "/blackjack", Description: "Play button-based blackjack with guided difficulty choices and wager caps.", Usage: "/blackjack amount:100 difficulty:easy"},
			{Name: "/lottery", Description: "Check the weekly lottery jackpot and status privately or publicly.", Usage: "/lottery public:false"},
			{Name: "/lottery-claim", Description: "Claim your prize if you won the weekly draw.", Usage: "/lottery-claim"},
			{Name: "/rich", Description: "Show the richest members in the economy privately or publicly.", Usage: "/rich public:false"},
			{Name: "/rob", Description: "Attempt to steal Pulses with clearer shield and fine outcomes.", Usage: "/rob user:@member"},
			{Name: "/shop", Description: "Browse item prices, effects, usability, and refunds.", Usage: "/shop public:false"},
			{Name: "/buy", Description: "Buy an item from the PulseKeep shop.", Usage: "/buy item:lucky_pickaxe"},
			{Name: "/inventory", Description: "View owned items plus estimated item value.", Usage: "/inventory public:false"},
			{Name: "/slots", Description: "Spin the slot machine with wager caps and clearer payout text.", Usage: "/slots amount:100"},
			{Name: "/fish", Description: "Catch fish with rarity bonuses and Treasure Map doubling.", Usage: "/fish"},
			{Name: "/mine", Description: "Mine ores with rarity bonuses and Treasure Map doubling.", Usage: "/mine"},
			{Name: "/gamble", Description: "Roll 1-100 with push, multiplier, and Lucky Clover outcomes.", Usage: "/gamble amount:100"},
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
			{Name: "/ticket add", Description: "Add a user to the current ticket channel.", Usage: "/ticket add user:@member"},
			{Name: "/ticket remove", Description: "Remove a user from the current ticket channel.", Usage: "/ticket remove user:@member"},
			{Name: "/ticket close", Description: "Close the current ticket channel.", Usage: "/ticket close"},
			{Name: "/ticket rename", Description: "Rename the current ticket channel.", Usage: "/ticket rename name:new-name"},
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
			WithFooterText("PulseKeep " + PulseKeepVersion + " · Support Tickets")).
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
		WithFooterText("PulseKeep " + PulseKeepVersion + " · Use the menu below")

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
		WithFooterText("PulseKeep " + PulseKeepVersion + " · Command browser")
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
