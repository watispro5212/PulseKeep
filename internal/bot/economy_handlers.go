package bot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/bot/economy"
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
	case "leaderboard":
		return leaderboardMessage(store), true
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
			AddField("Coinflip record", fmt.Sprintf("%dW / %dL", record.FlipWins, record.FlipLosses), true).
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

func economyCommandError(err error) discord.MessageCreate {
	switch {
	case errors.Is(err, economy.ErrInvalidAmount):
		return economyError("Invalid amount", "Use an amount greater than zero.")
	case errors.Is(err, economy.ErrInsufficientFund):
		return economyError("Not enough Pulses", "Your wallet does not have enough Pulses for that action.")
	case errors.Is(err, economy.ErrSelfPayment):
		return economyError("Payment blocked", "You cannot pay yourself.")
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
		WithFooterText("PulseKeep Economy")
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

func discordTimestamp(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return fmt.Sprintf("<t:%d:R>", t.Unix())
}
