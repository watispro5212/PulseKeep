package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/bot/economy"
	"github.com/watispro/pulsekeep/internal/cache"
)

type Bot struct {
	Client     *bot.Client
	cache      *cache.Cache
}

func New(token string, memCache *cache.Cache) *Bot {
	startedAt := time.Now()
	economyStore := economy.NewStore()
	_ = startedAt

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
		bot.WithEventListenerFunc(func(e *events.GuildJoin) {
			memCache.AddGuild(e.Guild.ID.String(), e.Guild.Name)
			memCache.UserCount.Add(int64(e.Guild.MemberCount))
			log.Printf("Joined guild: %s (%d members)", e.Guild.Name, e.Guild.MemberCount)
		}),
		bot.WithEventListenerFunc(func(e *events.GuildLeave) {
			memCache.UserCount.Add(-int64(e.Guild.MemberCount))
			memCache.RemoveGuild(e.Guild.ID.String())
		}),
		bot.WithEventListenerFunc(func(e *events.ApplicationCommandInteractionCreate) {
			memCache.IncrCommands()
			data := e.SlashCommandInteractionData()
			if response, ok := handleEconomyCommand(economyStore, e, data); ok {
				if err := e.CreateMessage(response); err != nil {
					log.Printf("failed to send economy response: %v", err)
				}
				return
			}

			switch data.CommandName() {
			case "help":
				if err := e.CreateMessage(commands.MenuMessage("", true)); err != nil {
					log.Printf("failed to send command menu: %v", err)
				}
			case "ticketpanel":
				if err := e.CreateMessage(commands.TicketPanelMessage(false)); err != nil {
					log.Printf("failed to send ticket panel: %v", err)
				}
		case "ping":
			if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
				discord.NewEmbed().
					WithTitle("🏓 Pong!").
					WithDescription("PulseKeep is online and responding to commands.").
					AddField("WebSocket", "Connected", true).
					AddField("API", "Reachable", true).
					WithColor(commands.UtilityMenuAccent).
					WithFooterText("PulseKeep Utility").
					WithTimestamp(time.Now()),
			)); err != nil {
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
				handleTicketOpen(e)
			case commands.TicketCloseButtonID:
				handleTicketClose(e)
			}
		}),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			log.Printf("Bot is ready as %s#%s", e.User.Username, e.User.Discriminator)
			memCache.GuildCount.Store(int64(len(e.Guilds)))
			if _, err := e.Client().Rest.SetGlobalCommands(e.Client().ApplicationID, commands.Register()); err != nil {
				log.Printf("failed to register global slash commands: %v", err)
			}
		}),
	)

	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
	}

	return &Bot{Client: client, cache: memCache}
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
			WithDescription("Current runtime snapshot for this bot process.").
			WithColor(commands.CommandMenuAccent).
			AddField("Status", "🟢 Online", true).
			AddField("Uptime", formatBotDuration(time.Since(startedAt)), true).
			AddField("Latency", "Real-time in /ping", true).
			AddField("Categories", "Utility · Moderation · Economy · Tickets", false).
			AddField("Get started", "Use `/help` to browse all commands.", false).
			WithFooterText("PulseKeep v5.0 · Go runtime").
			WithTimestamp(time.Now()))
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
			WithTitle(fmt.Sprintf("⏳ /%s", commandName)).
			WithDescription("This command is registered and visible in Discord. The full action handler is the next implementation step.").
			WithColor(commands.CommandMenuAccent).
			AddField("Try now", "Use `/help` to browse all finished command groups.", false).
			WithFooterText("PulseKeep · Coming soon").
			WithTimestamp(time.Now()))
}

func handleTicketOpen(e *events.ComponentInteractionCreate) {
	if err := e.DeferCreateMessage(true); err != nil {
		log.Printf("failed to defer ticket interaction: %v", err)
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent("Tickets can only be opened inside a server.").WithEphemeral(true))
		return
	}

	// check if user already has an open ticket
	existingChannels, err := e.Client().Rest.GetGuildChannels(*guildID)
	if err == nil {
		prefix := fmt.Sprintf("ticket-%s", strings.ToLower(e.User().Username))
		for _, ch := range existingChannels {
			if strings.HasPrefix(ch.Name(), prefix) {
				_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContentf("You already have an open ticket: %s. Please use that channel.", ch.Mention()).WithEphemeral(true))
				return
			}
		}
	}

	safeName := strings.ToLower(e.User().Username)
	safeName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, safeName)
	if len(safeName) > 24 {
		safeName = safeName[:24]
	}
	channelName := fmt.Sprintf("ticket-%s", safeName)

	gChannel, err := e.Client().Rest.CreateGuildChannel(*guildID, discord.GuildTextChannelCreate{
		Name:  channelName,
		Topic: fmt.Sprintf("Support ticket for %s | User ID: %s", e.User().Tag(), e.User().ID),
		PermissionOverwrites: []discord.PermissionOverwrite{
			discord.RolePermissionOverwrite{
				RoleID: *guildID,
				Deny:   discord.PermissionViewChannel,
			},
			discord.MemberPermissionOverwrite{
				UserID: e.User().ID,
				Allow:  discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory,
			},
		},
	})
	if err != nil {
		log.Printf("failed to create ticket channel: %v", err)
		_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent("Failed to create ticket. Please contact staff directly.").WithEphemeral(true))
		return
	}

	_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContentf("Your ticket has been created: %s", gChannel.Mention()).WithEphemeral(true))

	if _, err := e.Client().Rest.CreateMessage(gChannel.ID(), discord.NewMessageCreate().
		WithContentf("Welcome %s! A staff member will be with you shortly. Please describe your issue.", e.User().Mention()).
		AddActionRow(
			discord.NewDangerButton("Close Ticket", commands.TicketCloseButtonID),
		),
	); err != nil {
		log.Printf("failed to send welcome message in ticket: %v", err)
	}
}

func handleTicketClose(e *events.ComponentInteractionCreate) {
	if err := e.UpdateMessage(discord.NewMessageUpdate().WithContent("Ticket closed by user. Channel will be deleted shortly.")); err != nil {
		log.Printf("failed to update ticket close message: %v", err)
	}

	channelID := e.Channel().ID()
	go func() {
		time.Sleep(5 * time.Second)
		if err := e.Client().Rest.DeleteChannel(channelID); err != nil {
			log.Printf("failed to delete ticket channel: %v", err)
		}
	}()
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
