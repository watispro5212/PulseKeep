package bot

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/watispro5212/PulseKeep/internal/bot/commands"
	"github.com/watispro5212/PulseKeep/internal/bot/economy"
)

func handleEconomyCommand(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) (discord.MessageCreate, bool) {
	if msg, unavailable := economyUnavailable(store); unavailable {
		switch data.CommandName() {
		case "balance", "profile", "daily", "work", "pay", "coinflip", "rob", "shop", "buy", "inventory", "slots", "fish", "mine", "gamble", "sell", "use", "blackjack", "lottery", "lottery-claim", "lottery-config", "rich", "weekly", "gift":
			return msg, true
		}
	}
	switch data.CommandName() {
	case "balance":
		return balanceMessage(store, e, data), true
	case "profile":
		return profileMessage(store, e, data), true
	case "daily":
		return dailyMessage(store, e, data), true
	case "work":
		return workMessage(store, e, data), true
	case "pay":
		return payMessage(store, e, data), true
	case "coinflip":
		return coinflipMessage(store, e, data), true
	case "rob":
		return robMessage(store, e, data), true
	case "shop":
		return shopMessage(store, e, data), true
	case "buy":
		return buyMessage(store, e, data), true
	case "inventory":
		return inventoryMessage(store, e, data), true
	case "slots":
		return slotsMessage(store, e, data), true
	case "fish":
		return fishMessage(store, e, data), true
	case "mine":
		return mineMessage(store, e, data), true
	case "gamble":
		return gambleMessage(store, e, data), true
	case "sell":
		return sellMessage(store, e, data), true
	case "use":
		return useItemMessage(store, e, data), true
	case "blackjack":
		return blackjackMessage(store, e, data), true
	case "lottery":
		return lotteryMessage(store, e, data), true
	case "lottery-claim":
		return lotteryClaimMessage(store, e, data), true
	case "lottery-config":
		return lotteryConfigMessage(store, e, data), true
	case "rich":
		return richMessage(store, e, data), true
	case "weekly":
		return weeklyMessage(store, e, data), true
	case "gift":
		return giftMessage(store, e, data), true
	default:
		return discord.MessageCreate{}, false
	}
}

func economyUnavailable(store *economy.Store) (discord.MessageCreate, bool) {
	if store != nil {
		return discord.MessageCreate{}, false
	}
	return economyError("Economy unavailable", "The economy store is not ready yet. Try again in a moment."), true
}

func balanceMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}
	public := data.Bool("public")

	record := store.Balance(user.ID, user.EffectiveName())

	// Assign rank tier based on balance
	rank := "🦧 Broke"
	switch {
	case record.Balance >= 1_000_000:
		rank = "💎 Millionaire"
	case record.Balance >= 500_000:
		rank = "💰 Elite"
	case record.Balance >= 100_000:
		rank = "🏆 Rich"
	case record.Balance >= 50_000:
		rank = "⭐ Wealthy"
	case record.Balance >= 10_000:
		rank = "🟢 Growing"
	case record.Balance >= 1_000:
		rank = "🟡 Starter"
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("Pulse Balance", fmt.Sprintf("💳 %s has **%s Pulses** in their wallet.", user.String(), formatPulses(record.Balance))).
			AddField("Rank Tier", rank, true).
			AddField("Daily Streak", fmt.Sprintf("%d day(s) 🔥", record.DailyStreak), true).
			AddField("Net Worth", formatPulses(record.NetWorth()), true).
			AddField("Lifetime Earned", formatPulses(record.Earned), true).
			AddField("Lifetime Spent", formatPulses(record.Spent), true).
			AddField("Quick Commands", "`/pay` · `/profile` · `/shop` · `/daily`", false).
			WithThumbnail(user.EffectiveAvatarURL()))
}

func profileItemValue(items []economy.InventoryEntry) int {
	total := 0
	for _, item := range items {
		for _, shopItem := range economy.ShopItems {
			if shopItem.ID == item.ItemID {
				total += shopItem.Price * item.Quantity
				break
			}
		}
	}
	return total
}

func profileMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}
	public := data.Bool("public")

	record := store.Balance(user.ID, user.EffectiveName())
	items := store.Inventory(user.ID, user.EffectiveName())
	itemVal := profileItemValue(items)
	netWorth := record.Balance + itemVal
	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("Economy Profile", fmt.Sprintf("📊 Full economy breakdown for **%s**.", user.String())).
			AddField("Wallet", formatPulses(record.Balance), true).
			AddField("Net worth", fmt.Sprintf("%s (items: %s)", formatPulses(netWorth), formatPulses(itemVal)), true).
			AddField("Earned / Spent", fmt.Sprintf("%s / %s", formatPulses(record.Earned), formatPulses(record.Spent)), true).
			AddField("Daily streak", fmt.Sprintf("%d day(s)", record.DailyStreak), true).
			AddField("Weekly streak", fmt.Sprintf("%d week(s)", record.WeeklyStreak), true).
			AddField("Coinflip", fmt.Sprintf("%dW / %dL", record.FlipWins, record.FlipLosses), true).
			AddField("Gambling", fmt.Sprintf("%dW / %dL", record.GambleWins, record.GambleTotal-record.GambleWins), true).
			AddField("Blackjack", fmt.Sprintf("%dW / %dL", record.BlackjackWins, record.BlackjackLosses), true).
			AddField("Fishing", fmt.Sprintf("%d caught / %d total", record.FishCaught, record.FishTotal), true).
			AddField("Mining", fmt.Sprintf("%d mined / %d total", record.MineMined, record.MineTotal), true).
			AddField("Account started", discordTimestamp(record.CreatedAt), true).
			WithThumbnail(user.EffectiveAvatarURL()))
}

func dailyStreakMessage(streak int) string {
	switch streak {
	case 7:
		return "🔥 **7-day streak!** You're on fire! Keep it going!"
	case 14:
		return "⭐ **14-day streak!** You're unstoppable! Two weeks strong!"
	case 21:
		return "💪 **21-day streak!** Three weeks of dedication!"
	case 30:
		return "🏆 **30-day streak!** A full month! Legendary dedication!"
	case 50:
		return "👑 **50-day streak!** Half a century! Absolutely royal!"
	case 100:
		return "💯 **100-day streak!** PulseKeep royalty! You're a legend!"
	}
	return ""
}

func dailyMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	public := data.Bool("public")

	result := store.Daily(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Daily reward already claimed", result.NextAvailable)
	}

	desc := fmt.Sprintf("☀️ %s claimed their daily reward of **%s Pulses**!", e.User().String(), formatPulses(result.Reward))
	if milestone := dailyStreakMessage(result.Streak); milestone != "" {
		desc += "\n\n" + milestone
	}

	nextMilestone := ""
	switch {
	case result.Streak < 7:
		nextMilestone = fmt.Sprintf("%d/7 days", result.Streak)
	case result.Streak < 14:
		nextMilestone = fmt.Sprintf("%d/14 days", result.Streak)
	case result.Streak < 21:
		nextMilestone = fmt.Sprintf("%d/21 days", result.Streak)
	case result.Streak < 30:
		nextMilestone = fmt.Sprintf("%d/30 days", result.Streak)
	default:
		nextMilestone = "Max streak!"
	}

	streakBar := ""
	for i := 0; i < min(result.Streak, 10); i++ {
		streakBar += "🔥"
	}
	if result.Streak > 10 {
		streakBar += fmt.Sprintf(" (+%d)", result.Streak-10)
	}
	if streakBar == "" {
		streakBar = "☐ Start your streak!"
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("☀️ Daily Reward Claimed", desc).
			AddField("Streak", streakBar, false).
			AddField("Streak Bonus", fmt.Sprintf("+%s Pulses", formatPulses(result.StreakBonus)), true).
			AddField("New Balance", formatPulses(result.Record.Balance), true).
			AddField("Progress", nextMilestone, true).
			AddField("Next Claim", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func workMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	public := data.Bool("public")

	result := store.Work(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Work is on cooldown", result.NextAvailable)
	}

	desc := fmt.Sprintf("%s %s and earned **%s Pulses**.", e.User().String(), result.Job, formatPulses(result.Reward))
	if result.RareEvent != "" {
		desc += fmt.Sprintf("\nBonus event: **%s** (+%s).", result.RareEvent, formatPulses(result.Bonus))
	}
	if result.Boosted {
		desc += "\nXP Boost doubled the final payout."
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("Shift Complete", desc).
			AddField("Base pay", formatPulses(result.BaseReward), true).
			AddField("Pulses earned", formatPulses(result.Reward), true).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Next shift", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func payMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	recipient, ok := data.OptUser("recipient")
	if !ok {
		return economyError("Missing recipient", "Choose a member to pay.")
	}

	amount := data.Int("amount")
	result, err := store.Pay(e.User().ID, e.User().EffectiveName(), recipient.ID, recipient.EffectiveName(), amount, time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Payment Sent", fmt.Sprintf("💸 %s sent **%s Pulses** to %s.", e.User().String(), formatPulses(result.Amount), recipient.String())).
			AddField("Your balance", formatPulses(result.Sender.Balance), true).
			AddField("Recipient balance", formatPulses(result.Receiver.Balance), true).
			WithThumbnail(recipient.EffectiveAvatarURL()))
}

func coinflipMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	side := strings.ToLower(data.String("side"))
	wager := data.Int("amount")

	result, err := store.Coinflip(e.User().ID, e.User().EffectiveName(), side, wager, time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	outcome := "lost"
	color := commands.ModerationMenuAccent
	if result.Won {
		outcome = "won"
		color = commands.EconomyMenuAccent
	}
	boostText := ""
	if result.Boosted {
		boostText = "\nLucky Pickaxe was consumed for a second-chance flip."
	}

	coinEmoji := "🪙"
	outcomeEmoji := "🟥"
	if result.Won {
		outcomeEmoji = "🟢"
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(fmt.Sprintf("%s Pulse Coinflip", coinEmoji)).
			WithDescription(fmt.Sprintf("%s picked **%s** %s. The coin landed on **%s** — they **%s %s Pulses**.%s",
				e.User().String(), result.Side, outcomeEmoji, result.Result, outcome, formatPulses(result.Wager), boostText)).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("W/L Record", fmt.Sprintf("%dW / %dL", result.Record.FlipWins, result.Record.FlipLosses), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Coinflip"))
}

func richMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	records := store.Leaderboard(10)
	if len(records) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Pulse Wealth Leaderboard", "No economy accounts exist yet. Use `/daily` or `/work` to start earning."))
	}
	public := data.Bool("public")

	medals := []string{"🥇", "🥈", "🥉"}
	var rows strings.Builder
	for i, record := range records {
		name := record.Name
		if name == "" {
			name = record.UserID.String()
		}
		prefix := fmt.Sprintf("`%2d.`", i+1)
		if i < 3 {
			prefix = medals[i]
		}
		rows.WriteString(fmt.Sprintf("%s **%s** — `%s` Pulses\n", prefix, name, formatPulses(record.Balance)))
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("Pulse Wealth Leaderboard", strings.TrimSpace(rows.String())).
			AddField("How to climb", "Claim `/daily`, run `/work`, win bets, and keep earning!", false))
}

func weeklyMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	public := data.Bool("public")

	result := store.Weekly(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Weekly reward already claimed", result.NextAvailable)
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed("Weekly Reward Claimed", fmt.Sprintf("%s claimed **%s Pulses**!", e.User().String(), formatPulses(result.Reward))).
			AddField("Streak bonus", formatPulses(result.StreakBonus), true).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Weekly streak", fmt.Sprintf("%d week(s)", result.Streak), true).
			AddField("Next claim", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func giftMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	receiver, ok := data.OptUser("user")
	if !ok {
		return economyError("Missing recipient", "Choose a member to give the item to.")
	}

	itemID := strings.ToLower(data.String("item"))

	result, err := store.GiftItem(e.User().ID, e.User().EffectiveName(), receiver.ID, receiver.EffectiveName(), itemID, time.Now())
	if err != nil {
		if errors.Is(err, economy.ErrNotOwned) {
			return economyError("Item not owned", "You don't have that item. Use `/inventory` to see your items.")
		}
		if errors.Is(err, economy.ErrSelfPayment) {
			return economyError("Invalid target", "You cannot give items to yourself.")
		}
		return economyCommandError(err)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Item Gifted", fmt.Sprintf("%s gave **%s** to %s.", e.User().String(), result.Item.ItemName, receiver.String())).
			AddField("Tip", "Use `/inventory` to see your remaining items.", false).
			WithThumbnail(receiver.EffectiveAvatarURL()))
}

func robMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	target, ok := data.OptUser("user")
	if !ok {
		return economyError("Missing target", "Choose a member to rob.")
	}

	if target.ID == e.User().ID {
		return economyError("Invalid target", "You cannot rob yourself.")
	}

	result, err := store.Rob(e.User().ID, e.User().EffectiveName(), target.ID, target.EffectiveName(), time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	if result.OnCooldown {
		return cooldownMessage("Rob is on cooldown", result.NextAvailable)
	}

	if result.Success {
		return discord.NewMessageCreate().
			AddEmbeds(economyEmbed("Robbery Succeeded", fmt.Sprintf("%s robbed **%s Pulses** from %s!", e.User().String(), formatPulses(result.Stolen), target.String())).
				AddField("Your balance", formatPulses(result.Record.Balance), true).
				AddField("Target balance", formatPulses(result.Target.Balance), true).
				AddField("Next robbery", discordTimestamp(result.NextAvailable), true).
				WithThumbnail(e.User().EffectiveAvatarURL()))
	}

	if result.Shielded {
		return discord.NewMessageCreate().
			AddEmbeds(discord.NewEmbed().
				WithTitle("🛡️ Robbery Blocked").
				WithDescription(fmt.Sprintf("%s tried to rob %s, but a **Shield Token** blocked it!\nFine paid: **%s Pulses**.", e.User().String(), target.String(), formatPulses(result.Fine))).
				WithColor(commands.EconomyWarningAccent).
				AddField("Your balance", formatPulses(result.Record.Balance), true).
				AddField("Next robbery", discordTimestamp(result.NextAvailable), true).
				WithThumbnail(e.User().EffectiveAvatarURL()).
				WithFooterText("PulseKeep " + Version + " · Economy · Robbery"))
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle("🔫 Robbery Failed!").
			WithDescription(fmt.Sprintf("%s got caught trying to rob %s and paid a **%s Pulses** fine!", e.User().String(), target.String(), formatPulses(result.Fine))).
			WithColor(commands.ModerationMenuAccent).
			AddField("Your balance", formatPulses(result.Record.Balance), true).
			AddField("Next robbery", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Robbery"))
}

func shopMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	items := store.Shop()
	if len(items) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("PulseKeep Shop", "The shop is currently empty. Check back later!"))
	}
	record := store.Balance(e.User().ID, e.User().EffectiveName())
	public := data.Bool("public")

	embed := discord.NewEmbed().
		WithTitle("PulseKeep Shop").
		WithColor(commands.EconomyMenuAccent).
		WithDescription(fmt.Sprintf("Your balance: **%s Pulses**\nUse `/buy item:<id>` to purchase an item.\nUse `/inventory` to view your items.", formatPulses(record.Balance))).
		WithFooterText("PulseKeep " + Version + " · Economy · Shop").
		WithTimestamp(time.Now())

	for _, item := range items {
		embed = embed.AddField(fmt.Sprintf("%s — %s Pulses", item.Name, formatPulses(item.Price)),
			fmt.Sprintf("`%s` — %s", item.ID, item.Description), false)
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(embed)
}

func buyMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	itemID := strings.ToLower(data.String("item"))

	result, err := store.Buy(e.User().ID, e.User().EffectiveName(), itemID, time.Now())
	if err != nil {
		if errors.Is(err, economy.ErrItemNotFound) {
			return economyError("Item not found", "Use `/shop` to see available items.")
		}
		return economyCommandError(err)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Purchase Complete", fmt.Sprintf("%s bought **%s** for **%s Pulses**.", e.User().String(), result.Item.Name, formatPulses(result.Item.Price))).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Tip", "Use `/inventory` to see your items.", false).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func inventoryMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	items := store.Inventory(e.User().ID, e.User().EffectiveName())
	public := data.Bool("public")
	if len(items) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("🎒 Inventory", "You don't own any items yet. Use `/shop` to browse what's available."))
	}

	var rows strings.Builder
	totalValue := 0
	for _, item := range items {
		for _, shopItem := range economy.ShopItems {
			if shopItem.ID == item.ItemID {
				totalValue += shopItem.Price * item.Quantity
				break
			}
		}
		rows.WriteString(fmt.Sprintf("**%s** x%d\n", item.ItemName, item.Quantity))
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(economyEmbed(fmt.Sprintf("%s's Inventory", e.User().EffectiveName()), strings.TrimSpace(rows.String())).
			AddField("Estimated item value", formatPulses(totalValue), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func slotsMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	wager := data.Int("amount")

	result, err := store.Slots(e.User().ID, e.User().EffectiveName(), wager, time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	reels := fmt.Sprintf("%s %s %s", result.Symbols[0], result.Symbols[1], result.Symbols[2])
	title := "Pulse Slots"
	color := commands.ModerationMenuAccent
	outcome := "Lost"

	if result.Won {
		title = "Pulse Slots — You Won!"
		color = commands.EconomyMenuAccent
		outcome = "Won"
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("%s spun `%s` and **%s %s Pulses** (x%d).", e.User().String(), reels, outcome, formatPulses(max(result.Payout, result.Wager)), result.Multiplier)).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.SlotWins, result.Record.SlotLosses), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Slots"))
}

func economyCommandError(err error) discord.MessageCreate {
	switch {
	case errors.Is(err, economy.ErrInvalidAmount):
		return economyError("Invalid amount", fmt.Sprintf("Amount must be greater than zero. Max wager: **%s Pulses**. Max transfer: **%s Pulses**.", formatPulses(economy.MaxWager), formatPulses(economy.MaxTransfer)))
	case errors.Is(err, economy.ErrInsufficientFund):
		return economyError("Not enough Pulses", "You don't have enough Pulses for that action. Use `/daily`, `/work`, or `/fish` to earn more.")
	case errors.Is(err, economy.ErrSelfPayment):
		return economyError("Cannot pay yourself", "You cannot send Pulses or items to yourself. Choose a different recipient.")
	case errors.Is(err, economy.ErrItemNotFound):
		return economyError("Item not found", "That item doesn't exist. Use `/shop` to see available items.")
	case errors.Is(err, economy.ErrNotOwned):
		return economyError("You don't own that", "You don't have that item. Use `/inventory` to see your items.")
	case errors.Is(err, economy.ErrCannotUse):
		return economyError("Cannot use", "That item cannot be used. Try `/sell` for a refund instead.")
	case errors.Is(err, economy.ErrCooldown):
		return economyError("On cooldown", "That command is still on cooldown. Please wait before using it again.")
	default:
		return economyError("Something went wrong", "An unexpected error occurred. Try again or contact support if the issue persists.")
	}
}

func economyError(title string, description string) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(description).
			WithFooterText("PulseKeep " + Version + " · Economy Error").
			WithColor(commands.ModerationMenuAccent).
			WithTimestamp(time.Now()))
}

func cooldownMessage(title string, nextAvailable time.Time) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("Try again %s.", discordTimestamp(nextAvailable))).
			WithFooterText("PulseKeep "+Version+" · Economy").
			WithColor(commands.EconomyWarningAccent).
			AddField("Remaining", formatBotDuration(time.Until(nextAvailable)), true).
			WithTimestamp(time.Now()))
}

var economyTips = []string{
	"Use /daily every day to build your streak!",
	"Buy a Fishing Rod from /shop to unlock /fish!",
	"Roll 46+ in /gamble to win! 36-45 = push (refund).",
	"Buy a Shield Token to protect from robbers.",
	"Sell unwanted items with /sell for 60% refund.",
	"Use /weekly every 7 days for a big bonus!",
	"Multiple income sources = fastest growth!",
	"Play /blackjack on Easy for better odds!",
	"Run /work every 45 minutes for steady income.",
	"Buy XP Boost to double work earnings!",
	"Use /profile to track your net worth.",
	"Interest of 0.1% applies every 6 hours.",
}

var economyTipMu sync.Mutex
var economyTipIndex int

func nextEconomyTip() string {
	economyTipMu.Lock()
	defer economyTipMu.Unlock()
	tip := economyTips[economyTipIndex]
	economyTipIndex = (economyTipIndex + 1) % len(economyTips)
	return tip
}

func economyEmbed(title string, description string) discord.Embed {
	return discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithColor(commands.EconomyMenuAccent).
		WithFooterText("PulseKeep " + Version + " · Economy · 💡 " + nextEconomyTip()).
		WithTimestamp(time.Now())
}

func formatPulses(amount int) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	raw := strconv.Itoa(amount)
	for i := len(raw) - 3; i > 0; i -= 3 {
		raw = raw[:i] + "," + raw[i:]
	}
	return sign + raw
}

func fishMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	public := data.Bool("public")
	interest := store.ApplyInterest(e.User().ID, e.User().EffectiveName(), time.Now())

	result := store.Fish(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Fishing is on cooldown", result.NextAvailable)
	}

	if result.Reward == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Fishing — No Rod!", "You need a **Fishing Rod** to catch fish!\nBuy one from `/shop` (`fishing_rod` — 1,500 Pulses).").
				WithThumbnail(e.User().EffectiveAvatarURL()))
	}

	rarityColor := commands.EconomyMenuAccent
	switch result.Fish.Rarity {
	case "Common":
		rarityColor = 0x999999
	case "Uncommon":
		rarityColor = 0x4ade80 // green
	case "Rare":
		rarityColor = 0x60a5fa // blue
	case "Epic":
		rarityColor = 0xa78bfa // purple
	case "Legendary":
		rarityColor = 0xfbbf24 // gold
	case "Mythic":
		rarityColor = 0xfb7185 // pink/red
	case "Junk":
		rarityColor = 0x78716c // gray/brown
	}

	title := fmt.Sprintf("🎣 Fishing — %s %s", result.Fish.Emoji, result.Fish.Rarity)
	if result.Fish.Rarity == "Mythic" || result.Fish.Rarity == "Legendary" {
		title = fmt.Sprintf("🌟 INCREDIBLE CATCH! 🌟 — %s %s", result.Fish.Emoji, result.Fish.Rarity)
	} else if result.Fish.Rarity == "Junk" {
		title = fmt.Sprintf("🗑️ Trashed! — %s %s", result.Fish.Emoji, result.Fish.Rarity)
	}

	desc := fmt.Sprintf("%s cast a line and reeled in a...\n\n### %s %s\n> **Weight:** %s\n> **Sold for:** **%s Pulses**",
		e.User().String(), result.Fish.Emoji, result.Fish.Name, result.Fish.Weight, formatPulses(result.Reward))

	if result.Bonus > 0 {
		desc += fmt.Sprintf("\n> 💎 **Rarity Bonus:** +%s Pulses", formatPulses(result.Bonus))
	}
	if result.Boosted {
		desc += "\n> 🗺️ **Treasure Map:** Payout doubled!"
	}
	if interest > 0 {
		desc += fmt.Sprintf("\n> 💤 **Passive Interest:** +%s Pulses", formatPulses(interest))
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(desc).
			WithColor(rarityColor).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Total Caught", fmt.Sprintf("%d", result.Record.FishCaught), true).
			AddField("Cast again", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Fishing"))
}

func mineMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	public := data.Bool("public")
	interest := store.ApplyInterest(e.User().ID, e.User().EffectiveName(), time.Now())

	result := store.Mine(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Mining is on cooldown", result.NextAvailable)
	}

	if result.Reward == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Mining — No Pickaxe!", "You need an **Iron Pickaxe** to mine ore!\nBuy one from `/shop` (`iron_pickaxe` — 2,000 Pulses).").
				WithThumbnail(e.User().EffectiveAvatarURL()))
	}

	rarityColor := commands.EconomyMenuAccent
	switch result.Ore.Rarity {
	case "Common":
		rarityColor = 0x999999
	case "Uncommon":
		rarityColor = 0x4ade80
	case "Rare":
		rarityColor = 0x60a5fa
	case "Epic":
		rarityColor = 0xa78bfa
	case "Legendary":
		rarityColor = 0xfbbf24
	case "Mythic":
		rarityColor = 0xfb7185
	case "Junk":
		rarityColor = 0x78716c
	}

	title := fmt.Sprintf("⛏️ Mining — %s %s", result.Ore.Emoji, result.Ore.Rarity)
	if result.Ore.Rarity == "Mythic" || result.Ore.Rarity == "Legendary" {
		title = fmt.Sprintf("💎 RARE FIND! 💎 — %s %s", result.Ore.Emoji, result.Ore.Rarity)
	} else if result.Ore.Rarity == "Junk" {
		title = fmt.Sprintf("🪨 Just rubble... — %s %s", result.Ore.Emoji, result.Ore.Rarity)
	}

	desc := fmt.Sprintf("%s swung their pickaxe and uncovered...\n\n### %s %s\n> **Sold for:** **%s Pulses**",
		e.User().String(), result.Ore.Emoji, result.Ore.Name, formatPulses(result.Reward))

	if result.Bonus > 0 {
		desc += fmt.Sprintf("\n> ✨ **Rarity Bonus:** +%s Pulses", formatPulses(result.Bonus))
	}
	if result.Boosted {
		desc += "\n> 🗺️ **Treasure Map:** Payout doubled!"
	}
	if interest > 0 {
		desc += fmt.Sprintf("\n> 💤 **Passive Interest:** +%s Pulses", formatPulses(interest))
	}

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(desc).
			WithColor(rarityColor).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Ores mined", fmt.Sprintf("%d", result.Record.MineMined), true).
			AddField("Mine again", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Mining"))
}

func gambleMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	wager := data.Int("amount")

	result, err := store.Gamble(e.User().ID, e.User().EffectiveName(), wager, time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	if result.OnCooldown {
		return cooldownMessage("Gamble is on cooldown", result.NextAvailable)
	}

	title := "Pulse Gamble — Lost"
	color := commands.ModerationMenuAccent
	outcome := "lost"
	payoutStr := fmt.Sprintf("lost **%s Pulses**", formatPulses(result.Wager))

	if result.Push {
		title = "Pulse Gamble — Push"
		color = commands.EconomyWarningAccent
		outcome = "pushed"
		payoutStr = "pushed! Your wager is returned."
	} else if result.Won {
		payout := result.Payout
		title = "Pulse Gamble — Won!"
		color = commands.EconomyMenuAccent
		outcome = "won"
		payoutStr = fmt.Sprintf("won **%s Pulses** (x%d)", formatPulses(payout), result.Multiplier)
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("%s rolled **%d** (1-100) and %s!%s", e.User().String(), result.Roll, payoutStr, gambleBoostText(result.Boosted))).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.GambleWins, result.Record.GambleTotal-result.Record.GambleWins), true).
			AddField("Tip", "Roll 46+ to win! 36–45 = push (refund). 85+ = 2x, 95+ = 4x, 100 = 10x!", false).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText(fmt.Sprintf("PulseKeep "+Version+" · Economy · Gambling · %s", outcome)))
}

func gambleBoostText(boosted bool) string {
	if boosted {
		return "\nLucky Clover was consumed for a second-chance roll."
	}
	return ""
}

func sellMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	itemID := strings.ToLower(data.String("item"))

	result, err := store.Sell(e.User().ID, e.User().EffectiveName(), itemID, time.Now())
	if err != nil {
		if errors.Is(err, economy.ErrNotOwned) {
			return economyError("Item not owned", "You don't have that item. Use `/inventory` to see your items.")
		}
		return economyCommandError(err)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Item Sold", fmt.Sprintf("%s sold **%s** for **%s Pulses** (60%% refund).", e.User().String(), result.Item.ItemName, formatPulses(result.Reward))).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func useItemMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	itemID := strings.ToLower(data.String("item"))

	result, err := store.UseItem(e.User().ID, e.User().EffectiveName(), itemID, time.Now())
	if err != nil {
		if errors.Is(err, economy.ErrNotOwned) {
			return economyError("Item not owned", "You don't have that item. Use `/inventory` to see your items.")
		}
		if errors.Is(err, economy.ErrCannotUse) {
			return economyError("Cannot use", "That item cannot be used. Try `/sell` to get a refund instead.")
		}
		return economyCommandError(err)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed(fmt.Sprintf("Used: %s", result.ItemName), result.Description).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func blackjackMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	wager := data.Int("amount")
	if wager < 100 {
		return economyError("Minimum bet", "You must wager at least **100 Pulses** to play blackjack.")
	}

	difficultyStr, hasDifficulty := data.OptString("difficulty")
	difficulty := economy.BlackjackNormal
	if hasDifficulty {
		switch strings.ToLower(difficultyStr) {
		case "easy":
			difficulty = economy.BlackjackEasy
		case "normal":
			difficulty = economy.BlackjackNormal
		case "hard":
			difficulty = economy.BlackjackHard
		case "expert":
			difficulty = economy.BlackjackExpert
		}
	}

	result, err := store.BlackjackStart(e.User().ID, e.User().EffectiveName(), wager, difficulty, time.Now())
	if err != nil {
		return economyCommandError(err)
	}

	if result.GameOver {
		return finalBlackjackMessage(e, result.Player, result.Dealer, result.Wager, result.Won, result.Payout, result.Push, result.Natural, result.Record.Balance, result.Record.BlackjackWins, result.Record.BlackjackLosses, difficulty)
	}

	return ongoingBlackjackMessage(e, result.Player, result.Dealer, result.Wager, result.Record.Balance, difficulty)
}

func finalBlackjackMessage(e *events.ApplicationCommandInteractionCreate, player, dealer economy.BlackjackHand, wager int, won bool, payout int, push bool, natural bool, balance, wins, losses int, difficulty economy.BlackjackDifficulty) discord.MessageCreate {
	title := "♠️ Blackjack — Lost"
	color := commands.ModerationMenuAccent
	outcome := "lost"
	payoutStr := fmt.Sprintf("lost **%s Pulses**", formatPulses(wager))

	if push {
		title = "♠️ Blackjack — Push"
		color = commands.UtilityMenuAccent
		outcome = "pushed"
		payoutStr = "it's a tie! Your wager is returned."
	} else if natural {
		title = "♠️ Blackjack — Natural 21! 🎉"
		color = commands.EconomyMenuAccent
		outcome = "won"
		payoutStr = fmt.Sprintf("won **%s Pulses** (3:2)", formatPulses(payout))
	} else if won {
		title = "♠️ Blackjack — Won!"
		color = commands.EconomyMenuAccent
		outcome = "won"
		payoutStr = fmt.Sprintf("won **%s Pulses**", formatPulses(payout))
	}

	difficultyName := "Normal"
	switch difficulty {
	case economy.BlackjackEasy:
		difficultyName = "Easy"
	case economy.BlackjackHard:
		difficultyName = "Hard"
	case economy.BlackjackExpert:
		difficultyName = "Expert"
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("%s %s playing against **%s** CPU!", e.User().String(), payoutStr, difficultyName)).
			AddField("Your hand", fmt.Sprintf("`%s` = **%d**", player.CardStr, player.Value), true).
			AddField("Dealer hand", fmt.Sprintf("`%s` = **%d**", dealer.CardStr, dealer.Value), true).
			AddField("New balance", formatPulses(balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", wins, losses), true).
			WithColor(color).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText(fmt.Sprintf("PulseKeep "+Version+" · Economy · Blackjack · %s", outcome)))
}

func ongoingBlackjackMessage(e *events.ApplicationCommandInteractionCreate, player, dealer economy.BlackjackHand, wager int, balance int, difficulty economy.BlackjackDifficulty) discord.MessageCreate {
	difficultyName := "Normal"
	switch difficulty {
	case economy.BlackjackEasy:
		difficultyName = "Easy"
	case economy.BlackjackHard:
		difficultyName = "Hard"
	case economy.BlackjackExpert:
		difficultyName = "Expert"
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle("♠️ Blackjack — In Progress").
			WithDescription(fmt.Sprintf("%s is playing against **%s** CPU! Hit or stand?", e.User().String(), difficultyName)).
			AddField("Your hand", fmt.Sprintf("`%s` = **%d**", player.CardStr, player.Value), true).
			AddField("Dealer showing", fmt.Sprintf("`%s`", dealer.CardStr), true).
			AddField("Wager", formatPulses(wager), true).
			AddField("Balance", formatPulses(balance), true).
			WithColor(commands.EconomyMenuAccent).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep "+Version+" · Economy · Blackjack · Choose an action below")).
		AddActionRow(
			discord.NewPrimaryButton("Hit", economy.BlackjackHitCustomID),
			discord.NewDangerButton("Stand", economy.BlackjackStandCustomID),
		)
}

func handleBlackjackButton(store *economy.Store, e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()
	userID := e.User().ID

	var result economy.BlackjackTurnResult
	var err error

	switch customID {
	case economy.BlackjackHitCustomID:
		result, err = store.BlackjackHit(userID, time.Now())
	case economy.BlackjackStandCustomID:
		result, err = store.BlackjackStand(userID, time.Now())
	}

	if err != nil {
		if err2 := e.UpdateMessage(discord.NewMessageUpdate().WithContent("Error: " + err.Error())); err2 != nil {
			log.Printf("failed to update blackjack message: %v", err2)
		}
		return
	}

	if result.GameOver {
		title := "♠️ Blackjack — Lost"
		color := commands.ModerationMenuAccent
		outcome := "lost"
		payoutStr := fmt.Sprintf("lost **%s Pulses**", formatPulses(result.Wager))

		if result.Push {
			title = "♠️ Blackjack — Push"
			color = commands.UtilityMenuAccent
			outcome = "pushed"
			payoutStr = "it's a tie! Your wager is returned."
		} else if result.Won {
			title = "♠️ Blackjack — Won!"
			color = commands.EconomyMenuAccent
			outcome = "won"
			payoutStr = fmt.Sprintf("won **%s Pulses**", formatPulses(result.Payout))
		}

		embed := discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("You %s!", payoutStr)).
			AddField("Your hand", fmt.Sprintf("`%s` = **%d**", result.Player.CardStr, result.Player.Value), true)
		if result.Player.Bust {
			embed = embed.AddField("Result", "**Bust!** You went over 21.", true)
		}
		embed = embed.
			AddField("Dealer hand", fmt.Sprintf("`%s` = **%d**", result.Dealer.CardStr, result.Dealer.Value), true).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.BlackjackWins, result.Record.BlackjackLosses), true).
			WithColor(color).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText(fmt.Sprintf("PulseKeep "+Version+" · Economy · Blackjack · %s", outcome))

		if err := e.UpdateMessage(discord.NewMessageUpdate().WithEmbeds(embed).WithComponents()); err != nil {
			log.Printf("failed to update blackjack message: %v", err)
		}
	} else {
		embed := discord.NewEmbed().
			WithTitle("♠️ Blackjack — In Progress").
			WithDescription("Hit or stand?").
			AddField("Your hand", fmt.Sprintf("`%s` = **%d**", result.Player.CardStr, result.Player.Value), true).
			AddField("Dealer showing", fmt.Sprintf("`%s`", result.Dealer.CardStr), true).
			AddField("Wager", formatPulses(result.Wager), true).
			AddField("Balance", formatPulses(result.Record.Balance), true).
			WithColor(commands.EconomyMenuAccent).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep " + Version + " · Economy · Blackjack · Choose an action below")

		if err := e.UpdateMessage(discord.NewMessageUpdate().
			WithEmbeds(embed).
			AddActionRow(
				discord.NewPrimaryButton("Hit", economy.BlackjackHitCustomID),
				discord.NewDangerButton("Stand", economy.BlackjackStandCustomID),
			)); err != nil {
			log.Printf("failed to update blackjack message: %v", err)
		}
	}
}

func lotteryMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	status := store.LotteryStatus(time.Now())
	public := data.Bool("public")

	jackpotLowStr := formatPulses(status.JackpotLow)
	jackpotHighStr := formatPulses(status.JackpotHigh)

	modeLow := "Auto"
	if !status.AutoDrawLow {
		modeLow = "Manual"
	}
	modeHigh := "Auto"
	if !status.AutoDrawHigh {
		modeHigh = "Manual"
	}

	lastWinnerLow := "None yet"
	if status.LastWinnerIDLow != "" {
		lastWinnerLow = fmt.Sprintf("<@%s>", status.LastWinnerIDLow)
	}

	lastWinnerHigh := "None yet"
	if status.LastWinnerIDHigh != "" {
		lastWinnerHigh = fmt.Sprintf("<@%s>", status.LastWinnerIDHigh)
	}

	embed := economyEmbed("🎰 Weekly Lottery", "PulseKeep's dual-tier weekly lottery draw.").
		AddField("🟡 Low Tier (500 Pulses)", fmt.Sprintf("Jackpot: **%s**\nEntries: **%d**\nMode: %s\nLast Winner: %s", jackpotLowStr, status.TotalEntriesLow, modeLow, lastWinnerLow), true).
		AddField("💎 High Tier (5,000 Pulses)", fmt.Sprintf("Jackpot: **%s**\nEntries: **%d**\nMode: %s\nLast Winner: %s", jackpotHighStr, status.TotalEntriesHigh, modeHigh, lastWinnerHigh), true).
		AddField("How to enter", "Buy a **Low Lottery Ticket** or **High Lottery Ticket** from `/shop`, then use `/use`.", false).
		AddField("Claim", "Winners are drawn each week. Use `/lottery-claim tier:low/high` to claim.", false)

	return discord.NewMessageCreate().
		WithEphemeral(!public).
		AddEmbeds(embed)
}

func lotteryClaimMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	tier := data.String("tier")
	result, err := store.LotteryClaim(e.User().ID, e.User().EffectiveName(), tier, time.Now())
	if err != nil {
		return economyError("Lottery Claim", err.Error())
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed(fmt.Sprintf("🎉 %s Winner!", result.ItemName), result.Description).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func lotteryConfigMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	tier := data.String("tier")
	autodraw := data.Bool("autodraw")

	err := store.LotteryToggleAutoDraw(tier, autodraw)
	if err != nil {
		return economyError("Lottery Config", "Failed to update configuration.")
	}

	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(economyEmbed("⚙️ Lottery Configured", fmt.Sprintf("The **%s tier** lottery auto-draw is now set to: **%t**.", tier, autodraw)))
}

func discordTimestamp(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return fmt.Sprintf("<t:%d:R>", t.Unix())
}
