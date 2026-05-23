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
<<<<<<< HEAD
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/bot/economy"
=======
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/db"
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
)

type Bot struct {
	Client    bot.Client
	DB        *db.Database
	StartTime time.Time
}

<<<<<<< HEAD
func New(token string) *Bot {
	startedAt := time.Now()
	economyStore := economy.NewStore()
=======
func New(token string, database *db.Database) *Bot {
	cmdDefs := commands.GetCommands(database)

	// Track start time (will be updated when bot is ready)
	startTime := time.Now()
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuildMembers,
			),
		),
		bot.WithShardManagerConfigOpts(
			sharding.WithShardCount(1),
			sharding.WithShardIDs(0),
		),
		bot.WithEventListenerFunc(func(e *events.ApplicationCommandInteractionCreate) {
			data := e.Data
			if cmd, ok := cmdDefs[data.CommandName()]; ok {
				ctx := &commands.CommandContext{
					Event: e,
					DB:    database,
				}
				// Log command execution
				guildID := "0"
				if e.GuildID() != nil {
					guildID = e.GuildID().String()
				}
				database.IncrementCommandRun(context.Background(), guildID, e.User().ID.String(), data.CommandName())

				if err := cmd.Handler(ctx); err != nil {
					log.Printf("Error handling command %s: %v", data.CommandName(), err)
				}
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
<<<<<<< HEAD
			if _, err := e.Client().Rest.SetGlobalCommands(e.Client().ApplicationID, commands.Register()); err != nil {
				log.Printf("failed to register global slash commands: %v", err)
			}
=======
			// Update start time to actual ready time
			startTime = time.Now()
			if _, err := e.Client().Rest.SetGlobalCommands(e.Client().ApplicationID, commands.Register(database)); err != nil {
				log.Printf("Failed to register global commands: %v", err)
			} else {
				log.Println("Successfully registered global slash commands")
			}
		}),
		bot.WithEventListenerFunc(func(e *events.MessageDelete) {
			if e.GuildID == nil {
				return
			}
			cfg, err := database.GetGuildConfig(context.Background(), e.GuildID.String())
			if err != nil || cfg.LogChannelID == nil {
				return
			}
			logChanID, _ := snowflake.Parse(*cfg.LogChannelID)
			embed := discord.NewEmbed().
				WithTitle("Message Deleted").
				WithDescription(fmt.Sprintf("A message was deleted in <#%s>", e.ChannelID)).
				WithColor(0xff0000)
			_, _ = e.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
		}),
		bot.WithEventListenerFunc(func(e *events.MessageUpdate) {
			if e.GuildID == nil || e.Message.Author.Bot {
				return
			}
			cfg, err := database.GetGuildConfig(context.Background(), e.GuildID.String())
			if err != nil || cfg.LogChannelID == nil {
				return
			}
			logChanID, _ := snowflake.Parse(*cfg.LogChannelID)
			embed := discord.NewEmbed().
				WithTitle("Message Edited").
				WithDescription(fmt.Sprintf("A message was edited in <#%s>", e.ChannelID)).
				AddField("Old Content", e.OldMessage.Content, false).
				AddField("New Content", e.Message.Content, false).
				WithColor(0xffff00)
			_, _ = e.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
		}),
		bot.WithEventListenerFunc(func(e *events.GuildMemberJoin) {
			cfg, err := database.GetGuildConfig(context.Background(), e.GuildID.String())
			if err != nil || cfg.LogChannelID == nil {
				return
			}
			logChanID, _ := snowflake.Parse(*cfg.LogChannelID)
			embed := discord.NewEmbed().
				WithTitle("Member Joined").
				WithDescription(fmt.Sprintf("%s (%s) joined the server.", e.Member.User.Tag(), e.Member.User.ID)).
				WithColor(0x00ff00)
			_, _ = e.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
		}),
		bot.WithEventListenerFunc(func(e *events.GuildMemberLeave) {
			cfg, err := database.GetGuildConfig(context.Background(), e.GuildID.String())
			if err != nil || cfg.LogChannelID == nil {
				return
			}
			logChanID, _ := snowflake.Parse(*cfg.LogChannelID)
			embed := discord.NewEmbed().
				WithTitle("Member Left").
				WithDescription(fmt.Sprintf("%s (%s) left the server.", e.User.Tag(), e.User.ID)).
				WithColor(0xffa500)
			_, _ = e.Client().Rest.CreateMessage(logChanID, discord.MessageCreate{Embeds: []discord.Embed{embed}})
>>>>>>> 3e3c91af610d865bc7a95aea623510c44dcca715
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
	}
	return &Bot{Client: *client, DB: database, StartTime: startTime}
}

func (b *Bot) Start(ctx context.Context) error {
	return b.Client.OpenShardManager(ctx)
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
