package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/watispro5212/PulseKeep/internal/bot/commands"
	"github.com/watispro5212/PulseKeep/internal/bot/economy"
)

func handleEconomyCommand(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) (discord.MessageCreate, bool) {
	switch data.CommandName() {
	case "balance":
		return balanceMessage(store, e, data), true
	case "profile":
		return profileMessage(store, e, data), true
	case "daily":
		return dailyMessage(store, e), true
	case "work":
		return workMessage(store, e), true
	case "pay":
		return payMessage(store, e, data), true
	case "coinflip":
		return coinflipMessage(store, e, data), true
	case "rob":
		return robMessage(store, e, data), true
	case "shop":
		return shopMessage(store), true
	case "buy":
		return buyMessage(store, e, data), true
	case "inventory":
		return inventoryMessage(store, e), true
	case "slots":
		return slotsMessage(store, e, data), true
	case "fish":
		return fishMessage(store, e), true
	case "mine":
		return mineMessage(store, e), true
	case "gamble":
		return gambleMessage(store, e, data), true
	case "sell":
		return sellMessage(store, e, data), true
	case "use":
		return useItemMessage(store, e, data), true
	case "leaderboard":
		return leaderboardMessage(store), true
	case "rich":
		return richMessage(store), true
	case "weekly":
		return weeklyMessage(store, e), true
	case "gift":
		return giftMessage(store, e, data), true
	default:
		return discord.MessageCreate{}, false
	}
}

func balanceMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}

	record := store.Balance(user.ID, user.EffectiveName())
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(economyEmbed("Pulse Balance", fmt.Sprintf("%s has **%s Pulses**.", user.String(), formatPulses(record.Balance))).
			AddField("Lifetime earned", formatPulses(record.Earned), true).
			AddField("Lifetime spent", formatPulses(record.Spent), true).
			AddField("Daily streak", fmt.Sprintf("%d day(s)", record.DailyStreak), true).
			WithThumbnail(user.EffectiveAvatarURL()))
}

func profileMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}

	record := store.Balance(user.ID, user.EffectiveName())
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(economyEmbed("Economy Profile", fmt.Sprintf("Economy snapshot for %s.", user.String())).
			AddField("Wallet", formatPulses(record.Balance), true).
			AddField("Earned", formatPulses(record.Earned), true).
			AddField("Spent", formatPulses(record.Spent), true).
			AddField("Daily streak", fmt.Sprintf("%d day(s)", record.DailyStreak), true).
			AddField("Coinflip", fmt.Sprintf("%dW / %dL", record.FlipWins, record.FlipLosses), true).
			AddField("Gambling", fmt.Sprintf("%dW / %dL", record.GambleWins, record.GambleTotal-record.GambleWins), true).
			AddField("Fishing", fmt.Sprintf("%d caught / %d total", record.FishCaught, record.FishTotal), true).
			AddField("Mining", fmt.Sprintf("%d mined / %d total", record.MineMined, record.MineTotal), true).
			AddField("Account started", discordTimestamp(record.CreatedAt), true).
			WithThumbnail(user.EffectiveAvatarURL()))
}

func dailyMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	result := store.Daily(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Daily reward already claimed", result.NextAvailable)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Daily Reward Claimed", fmt.Sprintf("%s claimed **%s Pulses**.", e.User().String(), formatPulses(result.Reward))).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Streak", fmt.Sprintf("%d day(s)", result.Streak), true).
			AddField("Next claim", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()))
}

func workMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	result := store.Work(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Work is on cooldown", result.NextAvailable)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Shift Complete", fmt.Sprintf("%s %s and earned **%s Pulses**.", e.User().String(), result.Job, formatPulses(result.Reward))).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Next shift", discordTimestamp(result.NextAvailable), true).
			AddField("Tip", "Use `/profile` to view your economy stats.", false).
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
		AddEmbeds(economyEmbed("Payment Sent", fmt.Sprintf("%s sent **%s Pulses** to %s.", e.User().String(), formatPulses(result.Amount), recipient.String())).
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

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle("Pulse Coinflip").
			WithDescription(fmt.Sprintf("%s picked **%s**. The coin landed on **%s** and they **%s %s Pulses**.", e.User().String(), result.Side, result.Result, outcome, formatPulses(result.Wager))).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.FlipWins, result.Record.FlipLosses), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("Coinflip uses virtual Pulses only."))
}

func leaderboardMessage(store *economy.Store) discord.MessageCreate {
	records := store.Leaderboard(10)
	if len(records) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Pulse Leaderboard", "No economy accounts exist yet. Use `/daily` or `/work` to start earning."))
	}

	var rows strings.Builder
	for i, record := range records {
		name := record.Name
		if name == "" {
			name = record.UserID.String()
		}
		rows.WriteString(fmt.Sprintf("**%d. %s** - %s Pulses\n", i+1, name, formatPulses(record.Balance)))
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Pulse Leaderboard", strings.TrimSpace(rows.String())).
			AddField("How to climb", "Claim `/daily`, run `/work`, win `/coinflip`, and trade with `/pay`.", false))
}

func richMessage(store *economy.Store) discord.MessageCreate {
	records := store.Leaderboard(10)
	if len(records) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Pulse Wealth Leaderboard", "No economy accounts exist yet. Use `/daily` or `/work` to start earning."))
	}

	medals := []string{"🥇", "🥈", "🥉"}
	var rows strings.Builder
	for i, record := range records {
		name := record.Name
		if name == "" {
			name = record.UserID.String()
		}
		prefix := fmt.Sprintf("%d.", i+1)
		if i < 3 {
			prefix = medals[i]
		}
		rows.WriteString(fmt.Sprintf("%s **%s** — %s Pulses\n", prefix, name, formatPulses(record.Balance)))
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Pulse Wealth Leaderboard", strings.TrimSpace(rows.String())).
			AddField("How to climb", "Claim `/daily`, run `/work`, win bets, and keep earning!", false))
}

func weeklyMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	result := store.Weekly(e.User().ID, e.User().EffectiveName(), time.Now())
	if result.OnCooldown {
		return cooldownMessage("Weekly reward already claimed", result.NextAvailable)
	}

	return discord.NewMessageCreate().
		AddEmbeds(economyEmbed("Weekly Reward Claimed", fmt.Sprintf("%s claimed **%s Pulses**!", e.User().String(), formatPulses(result.Reward))).
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

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle("Robbery Failed!").
			WithDescription(fmt.Sprintf("%s got caught and paid a **%s Pulses** fine!", e.User().String(), formatPulses(result.Fine))).
			WithColor(commands.ModerationMenuAccent).
			AddField("Your balance", formatPulses(result.Record.Balance), true).
			AddField("Next robbery", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("Better luck next time."))
}

func shopMessage(store *economy.Store) discord.MessageCreate {
	items := store.Shop()
	if len(items) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("PulseKeep Shop", "The shop is currently empty. Check back later!"))
	}

	embed := discord.NewEmbed().
		WithTitle("PulseKeep Shop").
		WithColor(commands.EconomyMenuAccent).
		WithDescription("Use `/buy item:<id>` to purchase an item.\nUse `/inventory` to view your items.")

	for _, item := range items {
		embed = embed.AddField(fmt.Sprintf("%s — %s Pulses", item.Name, formatPulses(item.Price)),
			fmt.Sprintf("`%s` — %s", item.ID, item.Description), false)
	}

	return discord.NewMessageCreate().
		WithEphemeral(true).
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

func inventoryMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	items := store.Inventory(e.User().ID, e.User().EffectiveName())
	if len(items) == 0 {
		return discord.NewMessageCreate().
			WithEphemeral(true).
			AddEmbeds(economyEmbed("Inventory", "You don't own any items yet. Use `/shop` to browse what's available."))
	}

	var rows strings.Builder
	for _, item := range items {
		rows.WriteString(fmt.Sprintf("**%s** x%d\n", item.ItemName, item.Quantity))
	}

	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(economyEmbed(fmt.Sprintf("%s's Inventory", e.User().EffectiveName()), strings.TrimSpace(rows.String())).
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
			WithDescription(fmt.Sprintf("%s spun `%s` and **%s %s Pulses** (x%d).", e.User().String(), reels, outcome, formatPulses(result.Wager), result.Multiplier)).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.SlotWins, result.Record.SlotLosses), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep Slots"))
}

func economyCommandError(err error) discord.MessageCreate {
	switch {
	case errors.Is(err, economy.ErrInvalidAmount):
		return economyError("Invalid amount", "Use an amount greater than zero.")
	case errors.Is(err, economy.ErrInsufficientFund):
		return economyError("Not enough Pulses", "Your wallet does not have enough Pulses for that action.")
	case errors.Is(err, economy.ErrSelfPayment):
		return economyError("Payment blocked", "You cannot pay yourself.")
	case errors.Is(err, economy.ErrItemNotFound):
		return economyError("Item not found", "That item does not exist. Use `/shop` to see available items.")
	case errors.Is(err, economy.ErrNotOwned):
		return economyError("Item not owned", "You don't have that item. Use `/inventory` to see your items.")
	case errors.Is(err, economy.ErrCannotUse):
		return economyError("Cannot use", "That item cannot be used. Try `/sell` to get a refund instead.")
	case errors.Is(err, economy.ErrCooldown):
		return economyError("On cooldown", "Please wait before using that command again.")
	default:
		return economyError("Economy action failed", "Something went wrong while processing that command.")
	}
}

func economyError(title string, description string) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(description).
			WithColor(commands.ModerationMenuAccent))
}

func cooldownMessage(title string, nextAvailable time.Time) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("Try again %s.", discordTimestamp(nextAvailable))).
			WithColor(commands.EconomyMenuAccent).
			AddField("Remaining", formatBotDuration(time.Until(nextAvailable)), true))
}

func economyEmbed(title string, description string) discord.Embed {
	return discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithColor(commands.EconomyMenuAccent).
		WithFooterText("PulseKeep Economy · Use /help for all commands").
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

func fishMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
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

	desc := fmt.Sprintf("%s cast a line and caught... **%s %s** (%s, %s)!\nThey sold it for **%s Pulses**.",
		e.User().String(), result.Fish.Emoji, result.Fish.Name, result.Fish.Weight, result.Fish.Rarity, formatPulses(result.Reward))
	if interest > 0 {
		desc += fmt.Sprintf("\n💤 Passive interest: **+%s Pulses**", formatPulses(interest))
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(fmt.Sprintf("🎣 Fishing — %s", result.Fish.Rarity)).
			WithDescription(desc).
			WithColor(rarityColor).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Fish caught", fmt.Sprintf("%d", result.Record.FishCaught), true).
			AddField("Cast again", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep Fishing"))
}

func mineMessage(store *economy.Store, e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
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

	desc := fmt.Sprintf("%s swung their pickaxe and found... **%s %s** (%s)!\nThey sold it for **%s Pulses**.",
		e.User().String(), result.Ore.Emoji, result.Ore.Name, result.Ore.Rarity, formatPulses(result.Reward))
	if interest > 0 {
		desc += fmt.Sprintf("\n💤 Passive interest: **+%s Pulses**", formatPulses(interest))
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(fmt.Sprintf("⛏️ Mining — %s", result.Ore.Rarity)).
			WithDescription(desc).
			WithColor(rarityColor).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Ores mined", fmt.Sprintf("%d", result.Record.MineMined), true).
			AddField("Mine again", discordTimestamp(result.NextAvailable), true).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText("PulseKeep Mining"))
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

	if result.Won {
		payout := result.Wager * result.Multiplier
		title = "Pulse Gamble — Won!"
		color = commands.EconomyMenuAccent
		outcome = "won"
		payoutStr = fmt.Sprintf("won **%s Pulses** (x%d)", formatPulses(payout), result.Multiplier)
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(fmt.Sprintf("%s rolled **%d** (1-100) and %s!", e.User().String(), result.Roll, payoutStr)).
			WithColor(color).
			AddField("New balance", formatPulses(result.Record.Balance), true).
			AddField("Record", fmt.Sprintf("%dW / %dL", result.Record.GambleWins, result.Record.GambleTotal-result.Record.GambleWins), true).
			AddField("Tip", "Roll 60+ to win! 85+ = 2x, 95+ = 4x, 100 = 10x!", false).
			WithThumbnail(e.User().EffectiveAvatarURL()).
			WithFooterText(fmt.Sprintf("PulseKeep Gambling · %s", outcome)))
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

func discordTimestamp(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return fmt.Sprintf("<t:%d:R>", t.Unix())
}
