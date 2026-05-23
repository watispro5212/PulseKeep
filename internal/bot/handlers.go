package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/bot/economy"
)

type Bot struct {
	Client *bot.Client
}

func New(token string) *Bot {
	startedAt := time.Now()
	economyStore := economy.NewStore()

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(func(e *events.MessageCreate) {
			if e.Message.Author.Bot {
				return
			}
			if e.Message.Content == "!ping" {
				_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().WithContent("Pong from Go PulseKeep!"))
			}
			if e.Message.Content == "!menu" || e.Message.Content == "!help" {
				_, _ = e.Client().Rest.CreateMessage(e.ChannelID, commands.MenuMessage("", false).WithMessageReferenceByID(e.Message.ID))
			}
		}),
		bot.WithEventListenerFunc(func(e *events.ApplicationCommandInteractionCreate) {
			data := e.SlashCommandInteractionData()
			if response, ok := handleEconomyCommand(economyStore, e, data); ok {
				if err := e.CreateMessage(response); err != nil {
					log.Printf("failed to send economy response: %v", err)
				}
				return
			}

			switch data.CommandName() {
			case "help", "menu":
				if err := e.CreateMessage(commands.MenuMessage("", true)); err != nil {
					log.Printf("failed to send command menu: %v", err)
				}
			case "ticketpanel":
				if err := e.CreateMessage(commands.TicketPanelMessage(false)); err != nil {
					log.Printf("failed to send ticket panel: %v", err)
				}
			case "ping":
				if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("Pong from Go PulseKeep!")); err != nil {
					log.Printf("failed to send ping response: %v", err)
				}
			case "stats":
				if err := e.CreateMessage(statsMessage(startedAt)); err != nil {
					log.Printf("failed to send stats response: %v", err)
				}
			case "uptime":
				if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContentf("PulseKeep has been online for `%s`.", formatBotDuration(time.Since(startedAt)))); err != nil {
					log.Printf("failed to send uptime response: %v", err)
				}
			case "serverinfo":
				if err := e.CreateMessage(serverInfoMessage(e)); err != nil {
					log.Printf("failed to send server info response: %v", err)
				}
			case "userinfo":
				if err := e.CreateMessage(userInfoMessage(e, data)); err != nil {
					log.Printf("failed to send user info response: %v", err)
				}
			case "avatar":
				if err := e.CreateMessage(avatarMessage(e, data)); err != nil {
					log.Printf("failed to send avatar response: %v", err)
				}
			case "announce":
				if err := e.CreateMessage(announceMessage(data)); err != nil {
					log.Printf("failed to send announcement response: %v", err)
				}
			case "purge", "kick", "ban":
				if err := e.CreateMessage(comingSoonMessage(data.CommandName())); err != nil {
					log.Printf("failed to send command placeholder: %v", err)
				}
			default:
				if err := e.CreateMessage(commands.MenuMessage("", true)); err != nil {
					log.Printf("failed to send fallback command menu: %v", err)
				}
			}
		}),
		bot.WithEventListenerFunc(func(e *events.ComponentInteractionCreate) {
			switch e.Data.CustomID() {
			case commands.MenuSelectID:
				data := e.StringSelectMenuInteractionData()
				selected := "overview"
				if len(data.Values) > 0 {
					selected = data.Values[0]
				}
				if err := e.UpdateMessage(commands.MenuUpdate(selected)); err != nil {
					log.Printf("failed to update command menu: %v", err)
				}
			case commands.MenuOverviewButtonID:
				if err := e.UpdateMessage(commands.MenuUpdate("overview")); err != nil {
					log.Printf("failed to return command menu to overview: %v", err)
				}
			case commands.TicketPanelButtonID:
				if err := e.CreateMessage(commands.TicketPlaceholderMessage()); err != nil {
					log.Printf("failed to send ticket placeholder: %v", err)
				}
			}
		}),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			log.Printf("Bot is ready as %s#%s", e.User.Username, e.User.Discriminator)
			if _, err := e.Client().Rest.SetGlobalCommands(e.Client().ApplicationID, commands.Register()); err != nil {
				log.Printf("failed to register global slash commands: %v", err)
			}
		}),
	)

	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
	}

	return &Bot{Client: client}
}

func (b *Bot) Start(ctx context.Context) error {
	return b.Client.OpenGateway(ctx)
}

func (b *Bot) Stop(ctx context.Context) {
	b.Client.Close(ctx)
}

func statsMessage(startedAt time.Time) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle("PulseKeep Stats").
			WithDescription("Current runtime snapshot for this bot process. Database-backed counters can replace these values once persistence is wired into command execution.").
			WithColor(commands.CommandMenuAccent).
			AddField("Status", "Online", true).
			AddField("Uptime", formatBotDuration(time.Since(startedAt)), true).
			AddField("Command Menu", "`/help` or `/menu`", true).
			AddField("Registered Groups", "Utility, Moderation, Economy, Tickets", false))
}

func serverInfoMessage(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	embed := discord.NewEmbed().
		WithTitle("Server Info").
		WithColor(commands.UtilityMenuAccent)

	if guild, ok := e.Guild(); ok {
		embed = embed.
			WithDescription(fmt.Sprintf("Overview for **%s**.", guild.Name)).
			AddField("Server ID", guild.ID.String(), true).
			AddField("Members", fmt.Sprintf("%d", guild.MemberCount), true).
			AddField("Boost Tier", fmt.Sprintf("%d", guild.PremiumTier), true)
	} else {
		embed = embed.WithDescription("Server details are only available when this command is used inside a guild.")
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(embed)
}

func userInfoMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}

	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle("User Info").
			WithDescription(user.String()).
			WithColor(commands.UtilityMenuAccent).
			AddField("Display Name", user.EffectiveName(), true).
			AddField("Username", user.Tag(), true).
			AddField("User ID", user.ID.String(), false).
			WithThumbnail(user.EffectiveAvatarURL()))
}

func avatarMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}

	avatarURL := user.EffectiveAvatarURL()
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(fmt.Sprintf("%s's avatar", user.EffectiveName())).
			WithColor(commands.UtilityMenuAccent).
			WithImage(avatarURL).
			AddField("Direct link", avatarURL, false))
}

func announceMessage(data discord.SlashCommandInteractionData) discord.MessageCreate {
	title := data.String("title")
	body := data.String("message")
	if title == "" {
		title = "PulseKeep Announcement"
	}
	if body == "" {
		body = "No announcement message was provided."
	}

	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithTitle(title).
			WithDescription(body).
			WithColor(commands.CommandMenuAccent).
			WithFooterText("Sent with PulseKeep"))
}

func comingSoonMessage(commandName string) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle(fmt.Sprintf("/%s is registered", commandName)).
			WithDescription("This command is visible in Discord and documented in the interactive menu. Its full action handler is the next implementation step.").
			WithColor(commands.CommandMenuAccent).
			AddField("Try now", "Use `/help` or `/menu` to browse finished command groups.", false))
}

func formatBotDuration(d time.Duration) string {
	d = d.Round(time.Second)
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
