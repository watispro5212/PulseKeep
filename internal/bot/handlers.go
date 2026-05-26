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
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/watispro5212/PulseKeep/internal/bot/commands"
	"github.com/watispro5212/PulseKeep/internal/bot/economy"
	"github.com/watispro5212/PulseKeep/internal/cache"
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
			case "about":
				if err := e.CreateMessage(aboutMessage(e)); err != nil {
					log.Printf("failed to send about response: %v", err)
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
				if err := e.CreateMessage(announceMessage(e, data)); err != nil {
					log.Printf("failed to send announcement response: %v", err)
				}
			case "purge":
				if err := e.CreateMessage(handlePurge(e, data)); err != nil {
					log.Printf("failed to send purge response: %v", err)
				}
			case "kick":
				if err := e.CreateMessage(handleKick(e, data)); err != nil {
					log.Printf("failed to send kick response: %v", err)
				}
			case "ban":
				if err := e.CreateMessage(handleBan(e, data)); err != nil {
					log.Printf("failed to send ban response: %v", err)
				}
			case "poll":
				handlePoll(e, data)
			case "role":
				if err := e.CreateMessage(handleRole(e, data)); err != nil {
					log.Printf("failed to send role response: %v", err)
				}
			case "unban":
				if err := e.CreateMessage(handleUnban(e, data)); err != nil {
					log.Printf("failed to send unban response: %v", err)
				}
			case "slowmode":
				if err := e.CreateMessage(handleSlowmode(e, data)); err != nil {
					log.Printf("failed to send slowmode response: %v", err)
				}
			case "nick":
				if err := e.CreateMessage(handleNick(e, data)); err != nil {
					log.Printf("failed to send nick response: %v", err)
				}
			case "timeout":
				if err := e.CreateMessage(handleTimeout(e, data)); err != nil {
					log.Printf("failed to send timeout response: %v", err)
				}
			case "lock":
				if err := e.CreateMessage(handleLock(e)); err != nil {
					log.Printf("failed to send lock response: %v", err)
				}
			case "unlock":
				if err := e.CreateMessage(handleUnlock(e)); err != nil {
					log.Printf("failed to send unlock response: %v", err)
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
			startStatusRotation(context.Background(), e.Client())
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
			AddField("Categories", "Utility · Moderation · Economy · Tickets · Gambling", false).
			AddField("Get started", "Use `/help` to browse all commands.", false).
			WithFooterText("PulseKeep v5.4 · Go runtime").
			WithTimestamp(time.Now()))
}

func serverInfoMessage(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Server Info").
				WithDescription("Server details are only available inside a guild.").
				WithColor(commands.UtilityMenuAccent))
	}

	guild, err := e.Client().Rest.GetGuild(*guildID, true)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Server Info").
				WithDescription(fmt.Sprintf("Could not fetch server info: %s", err.Error())).
				WithColor(commands.UtilityMenuAccent))
	}

	createdAt := guild.ID.Time().Format("Jan 02, 2006")
	iconURL := ""
	if url := guild.IconURL(); url != nil {
		iconURL = *url
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle(fmt.Sprintf("🏠 %s", guild.Name)).
			WithDescription(fmt.Sprintf("Server overview for **%s**.", guild.Name)).
			WithColor(commands.UtilityMenuAccent).
			AddField("Owner", fmt.Sprintf("<@%s>", guild.OwnerID), true).
			AddField("Members", fmt.Sprintf("%d", guild.ApproximateMemberCount), true).
			AddField("Boosts", fmt.Sprintf("Tier %d (%d boosts)", guild.PremiumTier, guild.PremiumSubscriptionCount), true).
			AddField("Online", fmt.Sprintf("%d", guild.ApproximatePresenceCount), true).
			AddField("Roles", fmt.Sprintf("%d", len(guild.Roles)), true).
			AddField("Server ID", guild.ID.String(), false).
			AddField("Created", createdAt, false).
			WithThumbnail(iconURL).
			WithFooterText("PulseKeep Utility").
			WithTimestamp(time.Now()))
}

func userInfoMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	user, ok := data.OptUser("user")
	if !ok {
		user = e.User()
	}

	botBadge := "User"
	if user.Bot {
		botBadge = "🤖 Bot"
	}

	createdAt := user.ID.Time().Format("Jan 02, 2006")

	embed := discord.NewEmbed().
		WithTitle(user.EffectiveName()).
		WithDescription(fmt.Sprintf("%s · %s", user.Tag(), botBadge)).
		WithColor(commands.UtilityMenuAccent).
		AddField("User ID", user.ID.String(), false).
		AddField("Joined Discord", createdAt, true).
		WithThumbnail(user.EffectiveAvatarURL()).
		WithFooterText("PulseKeep Utility").
		WithTimestamp(time.Now())

	if guildID := e.GuildID(); guildID != nil {
		member, err := e.Client().Rest.GetMember(*guildID, user.ID)
		if err == nil {
			joinedAt := member.JoinedAt.Format("Jan 02, 2006")
			rolesCount := len(member.RoleIDs)
			embed = embed.AddField("Joined Server", joinedAt, true)
			if rolesCount > 0 {
				embed = embed.AddField("Roles", fmt.Sprintf("%d roles", rolesCount), true)
			} else {
				embed = embed.AddField("Roles", "None", true)
			}
		}
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(embed)
}

func aboutMessage(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	version := "v5.4"
	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("About PulseKeep").
			WithDescription("A modern, Go-powered Discord bot built for staff teams and community engagement.").
			AddField("Version", version, true).
			AddField("Language", "Go (disgo)", true).
			AddField("Database", "PostgreSQL (Neon)", true).
			AddField("Commands", "40+ slash commands across 4 categories", true).
			AddField("Creator", "watispro1 (Discord)", true).
			AddField("Open source", "Yes (MIT)", true).
			AddField("Links", "[Commands](https://pulsekeep.xyz/commands) · [Status](https://pulsekeep.xyz/status) · [Changelog](https://pulsekeep.xyz/changelog)", false).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep · Built for real Discord staff work").
			WithTimestamp(time.Now()))
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
			WithTitle(fmt.Sprintf("%s's avatar", user.Tag())).
			WithDescription(avatarURL).
			WithColor(commands.UtilityMenuAccent).
			WithImage(avatarURL).
			WithFooterText("PulseKeep Utility").
			WithTimestamp(time.Now()))
}

func announceMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	title := data.String("title")
	body := data.String("message")
	if title == "" {
		title = "PulseKeep Announcement"
	}
	if body == "" {
		body = "No announcement message was provided."
	}

	user := e.User()
	return discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithAuthorName(user.EffectiveName()).
			WithAuthorIcon(user.EffectiveAvatarURL()).
			WithTitle(title).
			WithDescription(body).
			WithColor(commands.CommandMenuAccent).
			WithFooterText("PulseKeep Announcements").
			WithTimestamp(time.Now()))
}

func handlePurge(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	amount := data.Int("amount")
	if amount < 1 || amount > 100 {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Invalid amount").
				WithDescription("Please specify an amount between 1 and 100.").
				WithColor(commands.ModerationMenuAccent))
	}

	channelID := e.Channel().ID()
	msgs, err := e.Client().Rest.GetMessages(channelID, *snowflakeID(0), *snowflakeID(0), *snowflakeID(0), amount)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Purge failed").
				WithDescription(fmt.Sprintf("Could not fetch messages: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	ids := make([]snowflake.ID, len(msgs))
	for i, m := range msgs {
		ids[i] = m.ID
	}

	if err := e.Client().Rest.BulkDeleteMessages(channelID, ids); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Purge failed").
				WithDescription(fmt.Sprintf("Could not delete messages: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Messages purged").
			WithDescription(fmt.Sprintf("Successfully deleted **%d** messages.", len(ids))).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleKick(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	user, ok := data.OptUser("user")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a member to kick.")
	}

	reason := data.String("reason")
	if reason == "" {
		reason = "No reason provided"
	}

	if err := e.Client().Rest.RemoveMember(*guildID, user.ID); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Kick failed").
				WithDescription(fmt.Sprintf("Could not kick %s: %s", user.Tag(), err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member kicked").
			WithDescription(fmt.Sprintf("**%s** has been kicked.", user.Tag())).
			AddField("Reason", reason, false).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleBan(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	user, ok := data.OptUser("user")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a member to ban.")
	}

	reason := data.String("reason")
	if reason == "" {
		reason = "No reason provided"
	}

	if err := e.Client().Rest.AddBan(*guildID, user.ID, 0); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Ban failed").
				WithDescription(fmt.Sprintf("Could not ban %s: %s", user.Tag(), err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member banned").
			WithDescription(fmt.Sprintf("**%s** has been banned.", user.Tag())).
			AddField("Reason", reason, false).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handlePoll(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	question := data.String("question")
	options := make([]string, 0, 4)
	emojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣"}

	if o, ok := data.OptString("option1"); ok { options = append(options, o) }
	if o, ok := data.OptString("option2"); ok { options = append(options, o) }
	if o, ok := data.OptString("option3"); ok { options = append(options, o) }
	if o, ok := data.OptString("option4"); ok { options = append(options, o) }

	if len(options) < 2 {
		if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("You need at least 2 options for a poll.")); err != nil {
			log.Printf("failed to send poll error: %v", err)
		}
		return
	}

	if err := e.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer poll: %v", err)
		return
	}

	var desc strings.Builder
	desc.WriteString(fmt.Sprintf("## %s\n\n", question))
	for i, opt := range options {
		desc.WriteString(fmt.Sprintf("%s %s\n", emojis[i], opt))
	}

	msg, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().AddEmbeds(
		discord.NewEmbed().
			WithTitle("📊 Poll").
			WithDescription(desc.String()).
			WithColor(commands.CommandMenuAccent).
			WithFooterText(fmt.Sprintf("Poll by %s", e.User().EffectiveName())).
			WithTimestamp(time.Now())))
	if err != nil {
		log.Printf("failed to send poll: %v", err)
		return
	}

	for i := 0; i < len(options); i++ {
		if err := e.Client().Rest.AddReaction(msg.ChannelID, msg.ID, emojis[i]); err != nil {
			log.Printf("failed to add reaction to poll: %v", err)
		}
	}
}

func handleRole(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	user, ok := data.OptUser("user")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a member.")
	}

	role, ok := data.OptRole("role")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a role.")
	}

	// try to remove first, if member already has the role
	err := e.Client().Rest.RemoveMemberRole(*guildID, user.ID, role.ID)
	if err == nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Role removed").
				WithDescription(fmt.Sprintf("Removed <@&%s> from **%s**.", role.ID, user.Tag())).
				WithColor(commands.UtilityMenuAccent).
				WithFooterText("PulseKeep Utility").
				WithTimestamp(time.Now()))
	}

	// member doesn't have the role, add it
	if err := e.Client().Rest.AddMemberRole(*guildID, user.ID, role.ID); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Role update failed").
				WithDescription(fmt.Sprintf("Could not update role for %s: %s", user.Tag(), err.Error())).
				WithColor(commands.UtilityMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Role added").
			WithDescription(fmt.Sprintf("Added <@&%s> to **%s**.", role.ID, user.Tag())).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Utility").
			WithTimestamp(time.Now()))
}

func handleUnban(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	userIDStr := data.String("user_id")
	userID, err := snowflake.Parse(userIDStr)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Invalid ID").
				WithDescription("Please provide a valid user ID.").
				WithColor(commands.ModerationMenuAccent))
	}

	if err := e.Client().Rest.DeleteBan(*guildID, userID); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Unban failed").
				WithDescription(fmt.Sprintf("Could not unban that user: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("User unbanned").
			WithDescription(fmt.Sprintf("Successfully unbanned `<@%s>`.", userID)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleSlowmode(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	seconds := data.Int("seconds")
	if seconds < 0 || seconds > 21600 {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Invalid slowmode").
				WithDescription("Slowmode must be between 0 and 21600 seconds.").
				WithColor(commands.ModerationMenuAccent))
	}

	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	channelID := e.Channel().ID()
	_, err := e.Client().Rest.UpdateChannel(channelID, discord.GuildTextChannelUpdate{
		RateLimitPerUser: &seconds,
	})
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Slowmode failed").
				WithDescription(fmt.Sprintf("Could not set slowmode: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	msg := fmt.Sprintf("Slowmode set to **%d seconds**.", seconds)
	if seconds == 0 {
		msg = "Slowmode has been **disabled**."
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Slowmode updated").
			WithDescription(msg).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleNick(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	user, ok := data.OptUser("user")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a member.")
	}

	nickname, _ := data.OptString("nickname")

	if _, err := e.Client().Rest.UpdateMember(*guildID, user.ID, discord.MemberUpdate{
		Nick: &nickname,
	}); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Nickname failed").
				WithDescription(fmt.Sprintf("Could not change nickname: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	if nickname == "" {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Nickname reset").
				WithDescription(fmt.Sprintf("Reset nickname for **%s**.", user.Tag())).
				WithColor(commands.UtilityMenuAccent).
				WithFooterText("PulseKeep Moderation").
				WithTimestamp(time.Now()))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Nickname changed").
			WithDescription(fmt.Sprintf("Changed **%s**'s nickname to **%s**.", user.Tag(), nickname)).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleTimeout(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	user, ok := data.OptUser("user")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a member.")
	}

	duration := data.Int("duration")
	if duration < 1 || duration > 40320 {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Invalid duration").
				WithDescription("Duration must be between 1 and 40320 minutes.").
				WithColor(commands.ModerationMenuAccent))
	}

	until := time.Now().Add(time.Duration(duration) * time.Minute)
	if _, err := e.Client().Rest.UpdateMember(*guildID, user.ID, discord.MemberUpdate{
		CommunicationDisabledUntil: omit.NewPtr(until),
	}); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Timeout failed").
				WithDescription(fmt.Sprintf("Could not timeout %s: %s", user.Tag(), err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member timed out").
			WithDescription(fmt.Sprintf("**%s** has been timed out for **%d minutes**.", user.Tag(), duration)).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleLock(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	channelID := e.Channel().ID()

	channel, err := e.Client().Rest.GetChannel(channelID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Lock failed").
				WithDescription(fmt.Sprintf("Could not fetch channel: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	guildChannel, ok := channel.(discord.GuildTextChannel)
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in text channels.")
	}

	existing := guildChannel.PermissionOverwrites()
	var overwrites []discord.PermissionOverwrite
	found := false
	for _, ov := range existing {
		if roleOv, ok := ov.(discord.RolePermissionOverwrite); ok && roleOv.RoleID == *guildID {
			overwrites = append(overwrites, discord.RolePermissionOverwrite{
				RoleID: *guildID,
				Deny:   roleOv.Deny | discord.PermissionSendMessages,
				Allow:  roleOv.Allow &^ discord.PermissionSendMessages,
			})
			found = true
		} else {
			overwrites = append(overwrites, ov)
		}
	}
	if !found {
		overwrites = append(overwrites, discord.RolePermissionOverwrite{
			RoleID: *guildID,
			Deny:   discord.PermissionSendMessages,
		})
	}

	if _, err := e.Client().Rest.UpdateChannel(channelID, discord.GuildTextChannelUpdate{
		PermissionOverwrites: &overwrites,
	}); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Lock failed").
				WithDescription(fmt.Sprintf("Could not lock channel: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("🔒 Channel locked").
			WithDescription(fmt.Sprintf("<#%s> has been locked.", channelID)).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func handleUnlock(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	channelID := e.Channel().ID()

	channel, err := e.Client().Rest.GetChannel(channelID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Unlock failed").
				WithDescription(fmt.Sprintf("Could not fetch channel: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	guildChannel, ok := channel.(discord.GuildTextChannel)
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in text channels.")
	}

	existing := guildChannel.PermissionOverwrites()
	var overwrites []discord.PermissionOverwrite
	for _, ov := range existing {
		if roleOv, ok := ov.(discord.RolePermissionOverwrite); ok && roleOv.RoleID == *guildID {
			deny := roleOv.Deny &^ discord.PermissionSendMessages
			allow := roleOv.Allow
			if deny != 0 || allow != 0 {
				overwrites = append(overwrites, discord.RolePermissionOverwrite{
					RoleID: *guildID,
					Deny:   deny,
					Allow:  allow,
				})
			}
		} else {
			overwrites = append(overwrites, ov)
		}
	}

	if _, err := e.Client().Rest.UpdateChannel(channelID, discord.GuildTextChannelUpdate{
		PermissionOverwrites: &overwrites,
	}); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Unlock failed").
				WithDescription(fmt.Sprintf("Could not unlock channel: %s", err.Error())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("🔓 Channel unlocked").
			WithDescription(fmt.Sprintf("<#%s> has been unlocked.", channelID)).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep Moderation").
			WithTimestamp(time.Now()))
}

func snowflakeID(i int) *snowflake.ID {
	id := snowflake.ID(i)
	return &id
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

func startStatusRotation(ctx context.Context, client *bot.Client) {
	statuses := []struct {
		text string
		kind string
	}{
		{"/help to browse commands", "playing"},
		{"over %d servers", "watching"},
		{"PulseKeep Economy", "playing"},
		{"/daily | /work | /slots", "playing"},
		{"/gamble | /fish | /mine", "playing"},
		{"moderation tickets", "watching"},
		{"/purge | /kick | /ban", "playing"},
		{"/slowmode | /lock | /timeout", "playing"},
		{"member activity", "watching"},
		{"/poll | /role | /announce", "playing"},
		{"PulseKeep v5.4", "playing"},
	}

	ticker := time.NewTicker(3 * time.Minute)
	idx := 0

	applyStatus := func() {
		s := statuses[idx%len(statuses)]
		idx++

		text := s.text
		if s.kind != "watching" && s.kind != "playing" {
			s.kind = "playing"
		}

		var opts []gateway.PresenceOpt
		switch s.kind {
		case "watching":
			opts = []gateway.PresenceOpt{
				gateway.WithWatchingActivity(text),
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
			}
		default:
			opts = []gateway.PresenceOpt{
				gateway.WithPlayingActivity(text),
				gateway.WithOnlineStatus(discord.OnlineStatusOnline),
			}
		}

		if err := client.SetPresence(ctx, opts...); err != nil {
			log.Printf("failed to set presence: %v", err)
		}
	}

	applyStatus()

	go func() {
		for {
			select {
			case <-ticker.C:
				applyStatus()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
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
