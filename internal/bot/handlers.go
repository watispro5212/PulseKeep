package bot

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
	"github.com/disgoorg/snowflake/v2"
	"github.com/watispro5212/PulseKeep/internal/bot/automod"
	"github.com/watispro5212/PulseKeep/internal/bot/commands"
	"github.com/watispro5212/PulseKeep/internal/bot/economy"
	"github.com/watispro5212/PulseKeep/internal/cache"
)

const Version = "v6.1.0"

type Bot struct {
	Client       *bot.Client
	cache        *cache.Cache
	db           *sql.DB
	automod      *automod.Engine
	cfgStore     *automod.ConfigStore
	webhookURL   string
	guildCount   int64
	economyStore *economy.Store
	httpClient   *http.Client
	webhookMu    sync.Mutex
	lastWebhook  map[string]time.Time
	statusCtx    context.Context
	statusCancel context.CancelFunc
	ticketCtx    context.Context
	ticketCancel context.CancelFunc
}

func New(token string, memCache *cache.Cache, database *sql.DB, webhookURL string) *Bot {
	startedAt := time.Now()
	economyStore := economy.NewStore(database)
	economyStore.StartLotteryAutoDraw(context.Background())
	cfgStore := automod.NewConfigStore(database)
	am := automod.NewEngine(cfgStore)

	botCtx, botCancel := context.WithCancel(context.Background())
	ticketCtx, ticketCancel := context.WithCancel(botCtx)

	b := &Bot{
		cache:        memCache,
		db:           database,
		automod:      am,
		cfgStore:     cfgStore,
		webhookURL:   webhookURL,
		economyStore: economyStore,
		httpClient:   &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{MaxIdleConns: 10, IdleConnTimeout: 30 * time.Second}},
		lastWebhook:  make(map[string]time.Time),
		statusCtx:    botCtx,
		statusCancel: botCancel,
		ticketCtx:    ticketCtx,
		ticketCancel: ticketCancel,
	}

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(b.onMessageCreate),
		bot.WithEventListenerFunc(b.onGuildJoin),
		bot.WithEventListenerFunc(b.onGuildLeave),
		bot.WithEventListenerFunc(func(e *events.ApplicationCommandInteractionCreate) {
			b.onSlashCommand(e, startedAt)
		}),
		bot.WithEventListenerFunc(b.onComponentInteraction),
		bot.WithEventListenerFunc(b.onReady),
	)
	if err != nil {
		log.Printf("error while building disgo instance: %s", err)
		return nil
	}

	b.Client = client
	return b
}

func (b *Bot) onMessageCreate(e *events.MessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if e.Message.Content == "!ping" {
		_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().WithContent("Pong from Go PulseKeep!"))
	}
	if e.Message.Content == "!menu" || e.Message.Content == "!help" {
		_, _ = e.Client().Rest.CreateMessage(e.ChannelID, commands.MenuMessage("", false).WithMessageReferenceByID(e.Message.ID))
	}
	if e.GuildID != nil && b.automod != nil {
		b.checkAutoMod(e)
	}
}

func (b *Bot) checkAutoMod(e *events.MessageCreate) {
	guildID := e.GuildID.String()
	result := b.automod.CheckMessage(guildID, e.Message.Author.ID.String(), e.Message.Content)
	if result.Action == automod.ActionNone {
		return
	}
	cfg := b.cfgStore.Get(guildID)
	if result.DeleteMsg {
		_ = e.Client().Rest.DeleteMessage(e.ChannelID, e.MessageID)
	}
	if cfg != nil && cfg.LogChannelID != "" {
		logEmbed := discord.NewEmbed().
			WithTitle("Auto-mod action").
			WithDescription(fmt.Sprintf("**User:** <@%s>\n**Action:** %s\n**Reason:** %s\n**Content:** %s", e.Message.Author.ID.String(), string(result.Action), result.Reason, e.Message.Content)).
			WithColor(0xfc8181).
			WithFooterText("PulseKeep " + Version + " · Auto-mod").
			WithTimestamp(time.Now())
		if logChanID, err := snowflake.Parse(cfg.LogChannelID); err == nil {
			_, _ = e.Client().Rest.CreateMessage(logChanID, discord.NewMessageCreate().AddEmbeds(logEmbed))
		}
	}
	if result.Action == automod.ActionWarn || result.Action == automod.ActionTimeout {
		warnEmbed := discord.NewEmbed().
			WithTitle("Auto-mod notice").
			WithDescription(fmt.Sprintf("Your message was removed: **%s**\nPlease follow the server rules.", result.Reason)).
			WithColor(0xf5bd4f).
			WithFooterText("PulseKeep " + Version + " · Auto-mod").
			WithTimestamp(time.Now())
		_, _ = e.Client().Rest.CreateMessage(e.ChannelID, discord.NewMessageCreate().AddEmbeds(warnEmbed))
	}
}

func (b *Bot) onGuildJoin(e *events.GuildJoin) {
	b.cache.AddGuild(e.Guild.ID.String(), e.Guild.Name)
	b.cache.UserCount.Add(int64(e.Guild.MemberCount))
	b.guildCount = b.cache.GuildCount.Load()
	log.Printf("Joined guild: %s (%d members)", e.Guild.Name, e.Guild.MemberCount)
	b.sendWebhook(discord.NewEmbed().
		WithTitle("Joined a new server").
		WithDescription(fmt.Sprintf("PulseKeep was added to **%s**.", e.Guild.Name)).
		AddField("Members", fmt.Sprintf("%d", e.Guild.MemberCount), true).
		AddField("Total servers", fmt.Sprintf("%d", b.guildCount), true).
		AddField("Owner", fmt.Sprintf("<@%s>", e.Guild.OwnerID), true).
		WithColor(0x4f8cff).
		WithFooterText("PulseKeep " + Version + " · Status").
		WithTimestamp(time.Now()))
}

func (b *Bot) onGuildLeave(e *events.GuildLeave) {
	b.cache.UserCount.Add(-int64(e.Guild.MemberCount))
	b.cache.RemoveGuild(e.Guild.ID.String())
}

func (b *Bot) onReady(e *events.Ready) {
	log.Printf("Bot is ready as %s", e.User.EffectiveName())
	gc := int64(len(e.Guilds))
	b.guildCount = gc
	b.cache.ResetGuilds()
	b.cache.GuildCount.Store(gc)
	for _, g := range e.Guilds {
		b.cache.AddGuild(g.ID.String(), g.ID.String())
	}
	if _, err := e.Client().Rest.SetGlobalCommands(e.Client().ApplicationID, commands.Register()); err != nil {
		log.Printf("failed to register global slash commands: %v", err)
	}
	if b.statusCancel != nil {
		b.statusCancel()
	}
	b.statusCtx, b.statusCancel = context.WithCancel(context.Background())
	startStatusRotation(b.statusCtx, e.Client(), b.cache)
	log.Printf("PulseKeep online — %d guilds", gc)
	b.sendWebhook(discord.NewEmbed().
		WithTitle("PulseKeep is online").
		WithDescription(fmt.Sprintf("**%s** is ready across **%d** guilds.", e.User.EffectiveName(), gc)).
		AddField("Users", fmt.Sprintf("%d", b.cache.UserCount.Load()), true).
		WithColor(0x36d399).
		WithFooterText("PulseKeep " + Version + " · Status").
		WithTimestamp(time.Now()))
}

func (b *Bot) onSlashCommand(e *events.ApplicationCommandInteractionCreate, startedAt time.Time) {
	t0 := time.Now()
	b.cache.IncrCommands()
	data := e.SlashCommandInteractionData()
	cmdName := data.CommandName()

	if guildID := e.GuildID(); guildID != nil {
		cfg := b.cfgStore.Get(guildID.String())
		if cfg != nil {
			if isEconomyCommand(cmdName) && !cfg.EconomyEnabled {
				b.respond(e, featureDisabledEmbed("Economy", "Economy commands are disabled in this server."))
				return
			}
			if cmdName == "ticketpanel" && !cfg.TicketsEnabled {
				b.respond(e, featureDisabledEmbed("Tickets", "The ticket system is disabled in this server."))
				return
			}
			if cmdName == "ticket" && !cfg.TicketsEnabled {
				b.respond(e, featureDisabledEmbed("Tickets", "The ticket system is disabled in this server."))
				return
			}
		}
	}

	if b.db != nil {
		gID := ""
		if e.GuildID() != nil {
			gID = e.GuildID().String()
		}
		_, _ = b.db.ExecContext(context.Background(), "INSERT INTO command_logs (guild_id, user_id, command_name) VALUES ($1, $2, $3)", gID, e.User().ID.String(), cmdName)
	}

	if response, ok := handleEconomyCommand(b.economyStore, e, data); ok {
		if err := e.CreateMessage(response); err != nil {
			log.Printf("failed to send economy response: %v", err)
		}
		b.cache.AddLatency(time.Since(t0))
		return
	}

	switch cmdName {
	case "magic8ball":
		question := data.String("question")
		answers := []string{
			"It is certain.", "It is decidedly so.", "Without a doubt.", "Yes, definitely.",
			"You may rely on it.", "As I see it, yes.", "Most likely.", "Outlook good.",
			"Yes.", "Signs point to yes.", "Reply hazy, try again.", "Ask again later.",
			"Better not tell you now.", "Cannot predict now.", "Concentrate and ask again.",
			"Don't count on it.", "My reply is no.", "My sources say no.",
			"Outlook not so good.", "Very doubtful.",
		}
		answer := answers[rand.Intn(len(answers))]
		embed := discord.NewEmbed().
			WithTitle("🎱 Magic 8-Ball").
			WithDescription(fmt.Sprintf("**Question:** %s\n**Answer:** %s", question, answer)).
			WithColor(commands.UtilityMenuAccent)
		b.respond(e, discord.NewMessageCreate().AddEmbeds(embed))
	case "help":
		b.respond(e, commands.MenuMessage("", true))
	case "ticketpanel":
		b.respond(e, commands.TicketPanelMessage(false))
	case "about":
		b.respond(e, aboutMessage(e))
	case "ping":
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("🏓 Pong!").
				WithDescription("PulseKeep is online and responding to commands.").
				AddField("WebSocket", "Connected", true).
				AddField("API", "Reachable", true).
				AddField("Latency", fmt.Sprintf("%dms", time.Since(t0).Milliseconds()), true).
				WithColor(commands.UtilityMenuAccent).
				WithFooterText("PulseKeep "+Version+" · Utility").
				WithTimestamp(time.Now()),
		))
	case "stats":
		b.respond(e, statsMessage(startedAt))
	case "uptime":
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContentf("PulseKeep has been online for `%s`.", formatBotDuration(time.Since(startedAt))))
	case "serverinfo":
		b.respond(e, serverInfoMessage(e))
	case "userinfo":
		b.respond(e, userInfoMessage(e, data))
	case "avatar":
		b.respond(e, avatarMessage(e, data))
	case "announce":
		b.respond(e, announceMessage(e, data))
	case "purge":
		b.respond(e, handlePurge(e, data))
	case "kick":
		b.respond(e, handleKick(e, data))
	case "ban":
		b.respond(e, handleBan(e, data))
	case "poll":
		handlePoll(e, data)
	case "role":
		b.respond(e, handleRole(e, data))
	case "unban":
		b.respond(e, handleUnban(e, data))
	case "slowmode":
		b.respond(e, handleSlowmode(e, data))
	case "nick":
		b.respond(e, handleNick(e, data))
	case "timeout":
		b.respond(e, handleTimeout(e, data))
	case "lock":
		b.respond(e, handleLock(e))
	case "unlock":
		b.respond(e, handleUnlock(e))
	case "warn":
		b.handleWarn(e, data)
	case "warnings":
		b.handleWarnings(e, data)
	case "clearwarns":
		b.handleClearWarns(e, data)
	case "move":
		b.handleMove(e, data)
	case "vckick":
		b.handleVCKick(e, data)
	case "ticket":
		b.handleTicketCmd(e, data)
	case "servericon":
		b.respond(e, serverIconMessage(e))
	case "roleinfo":
		b.respond(e, roleInfoMessage(e, data))
	case "channelinfo":
		b.respond(e, channelInfoMessage(e, data))
	case "invite":
		b.respond(e, inviteMessage())
	case "tip":
		b.respond(e, tipMessage())
	case "vote":
		b.respond(e, voteMessage())
	default:
		b.respond(e, commands.MenuMessage("", true))
	}
	b.cache.AddLatency(time.Since(t0))
}

func (b *Bot) onComponentInteraction(e *events.ComponentInteractionCreate) {
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
		b.handleTicketOpen(e)
	case commands.TicketCloseButtonID:
		b.handleTicketClose(e)
	case economy.BlackjackHitCustomID, economy.BlackjackStandCustomID:
		handleBlackjackButton(b.economyStore, e)
	}
}

func (b *Bot) respond(e *events.ApplicationCommandInteractionCreate, msg discord.MessageCreate) {
	if err := e.CreateMessage(msg); err != nil {
		log.Printf("failed to send response for %s: %v", e.SlashCommandInteractionData().CommandName(), err)
	}
}

func (b *Bot) handleTicketOpen(e *events.ComponentInteractionCreate) {
	if err := e.DeferCreateMessage(true); err != nil {
		log.Printf("failed to defer ticket interaction: %v", err)
		return
	}

	guildID := e.GuildID()
	if guildID == nil {
		_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent("Tickets can only be opened inside a server.").WithEphemeral(true))
		return
	}

	cfg := b.cfgStore.Get(guildID.String())
	if cfg != nil && !cfg.TicketsEnabled {
		_, _ = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent("The ticket system is disabled in this server.").WithEphemeral(true))
		return
	}

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
			discord.MemberPermissionOverwrite{
				UserID: e.Client().ApplicationID,
				Allow: discord.PermissionViewChannel |
					discord.PermissionSendMessages |
					discord.PermissionReadMessageHistory |
					discord.PermissionManageChannels |
					discord.PermissionManageMessages,
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

func (b *Bot) handleTicketClose(e *events.ComponentInteractionCreate) {
	if err := e.DeferUpdateMessage(); err != nil {
		log.Printf("failed to defer ticket close: %v", err)
		return
	}

	_, _ = e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(),
		discord.NewMessageUpdate().
			WithContentf("Ticket closed by %s. Channel will be deleted shortly.", e.User().Mention()).
			WithComponents())

	b.doTicketClose(e.Client().Rest, e.Message.ChannelID)
}

func (b *Bot) doTicketClose(restClient rest.Rest, channelID snowflake.ID) {
	go func() {
		select {
		case <-b.ticketCtx.Done():
			return
		case <-time.After(2 * time.Second):
		}

		if err := restClient.DeleteChannel(channelID); err != nil {
			log.Printf("failed to delete ticket channel %s: %v", channelID, err)
		}
	}()
}

func (b *Bot) GetConfigStore() *automod.ConfigStore {
	return b.cfgStore
}

func (b *Bot) GetEconomyStore() *economy.Store {
	return b.economyStore
}

func (b *Bot) Start(ctx context.Context) error {
	return b.Client.OpenGateway(ctx)
}

func (b *Bot) Stop(ctx context.Context) {
	if b.economyStore != nil {
		b.economyStore.FlushAll()
	}
	b.sendWebhook(discord.NewEmbed().
		WithTitle("PulseKeep went offline").
		WithDescription(fmt.Sprintf("Bot process is shutting down. Served **%d** guilds.", b.guildCount)).
		WithColor(0xfb7185).
		WithFooterText("PulseKeep " + Version + " · Status").
		WithTimestamp(time.Now()))
	if b.ticketCancel != nil {
		b.ticketCancel()
	}
	if b.statusCancel != nil {
		b.statusCancel()
	}
	if ctx != nil && b.Client != nil {
		b.Client.Close(ctx)
	}
}

func (b *Bot) sendWebhook(embed discord.Embed) {
	if b.webhookURL == "" {
		return
	}
	key := fmt.Sprintf("%s|%s|%d", embed.Title, embed.Description, embed.Color)
	now := time.Now()
	b.webhookMu.Lock()
	for existingKey, sentAt := range b.lastWebhook {
		if now.Sub(sentAt) > 10*time.Minute {
			delete(b.lastWebhook, existingKey)
		}
	}
	if sentAt, ok := b.lastWebhook[key]; ok && now.Sub(sentAt) < 2*time.Minute {
		b.webhookMu.Unlock()
		log.Printf("skipping duplicate status webhook: %s", embed.Title)
		return
	}
	b.lastWebhook[key] = now
	b.webhookMu.Unlock()

	payload, err := json.Marshal(map[string]any{
		"embeds": []discord.Embed{embed},
	})
	if err != nil {
		log.Printf("failed to marshal webhook payload: %v", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, b.webhookURL, bytes.NewReader(payload))
	if err != nil {
		log.Printf("failed to create webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := b.httpClient.Do(req)
	if err != nil {
		log.Printf("failed to send webhook: %v", err)
		return
	}
	resp.Body.Close()
}

func statsMessage(startedAt time.Time) discord.MessageCreate {
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle("PulseKeep Stats").
			WithDescription("Current runtime snapshot for this bot process.").
			WithColor(commands.CommandMenuAccent).
			AddField("Status", "🟢 Online", true).
			AddField("Version", Version, true).
			AddField("Uptime", formatBotDuration(time.Since(startedAt)), true).
			AddField("Language", "Go (disgo v0.19)", true).
			AddField("Database", "PostgreSQL (Neon)", true).
			AddField("Commands", "40+ slash commands", true).
			AddField("Categories", "Utility · Moderation · Economy · Tickets", false).
			AddField("Get started", "Use `/help` to browse all commands.", false).
			WithFooterText("PulseKeep " + Version + " · Stats").
			WithTimestamp(time.Now()))
}

func voteMessage() discord.MessageCreate {
	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Vote for PulseKeep").
			WithDescription("Support PulseKeep by voting on Top.gg! Each vote helps us grow and reach more communities.").
			AddField("Top.gg", "https://top.gg/bot/1507498795569512598/vote", false).
			AddField("Rewards", "Voting helps us rank higher and attract new users. Thank you for your support!", false).
			WithColor(commands.EconomyMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Voting").
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
				WithDescription("Could not fetch server info. Make sure I have access to this server.").
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
			WithFooterText("PulseKeep " + Version + " · Utility").
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
		WithFooterText("PulseKeep " + Version + " · Utility").
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
	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("About PulseKeep").
			WithDescription("A modern, Go-powered Discord bot built for staff teams and community engagement.").
			AddField("Version", Version, true).
			AddField("Language", "Go (disgo)", true).
			AddField("Database", "PostgreSQL (Neon)", true).
			AddField("Commands", "40+ slash commands across 4 categories", true).
			AddField("Creator", "watispro1 (Discord)", true).
			AddField("Open source", "Yes (MIT)", true).
			AddField("Links", "[Commands](https://pulsekeep.williamdelilah3.workers.dev/commands) · [Status](https://pulsekeep.williamdelilah3.workers.dev/status) · [Changelog](https://pulsekeep.williamdelilah3.workers.dev/changelog)", false).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep " + Version + " · About").
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
			WithFooterText("PulseKeep " + Version + " · Utility").
			WithTimestamp(time.Now()))
}

func announceMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	title := data.String("title")
	body := data.String("message")
	ping := data.Bool("ping")
	if title == "" {
		title = "PulseKeep Announcement"
	}
	if body == "" {
		body = "No announcement message was provided."
	}

	member := e.Member()
	if ping && (member == nil || (member.Permissions&discord.PermissionMentionEveryone) == 0) {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Missing permission").
				WithDescription("You need **Mention Everyone** permission to ping @everyone.").
				WithColor(commands.ModerationMenuAccent))
	}

	user := e.User()
	msg := discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithAuthorName(user.EffectiveName()).
			WithAuthorIcon(user.EffectiveAvatarURL()).
			WithTitle(title).
			WithDescription(body).
			WithColor(commands.CommandMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Announcements").
			WithTimestamp(time.Now()))
	if ping {
		msg = msg.WithContent("@everyone")
	} else {
		msg = msg.WithAllowedMentions(&discord.AllowedMentions{})
	}
	return msg
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
				WithDescription("Could not fetch messages. Make sure I have **Read Message History** permission.").
				WithColor(commands.ModerationMenuAccent))
	}

	if len(msgs) == 0 {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Nothing to purge").
				WithDescription("No messages found to delete.").
				WithColor(commands.UtilityMenuAccent))
	}

	ids := make([]snowflake.ID, len(msgs))
	for i, m := range msgs {
		ids[i] = m.ID
	}

	if err := e.Client().Rest.BulkDeleteMessages(channelID, ids); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Purge failed").
				WithDescription("Could not delete messages. Make sure I have **Manage Messages** permission.").
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Messages purged").
			WithDescription(fmt.Sprintf("Successfully deleted **%d** messages.", len(ids))).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
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

	if user.ID == e.User().ID {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You cannot kick yourself.")
	}

	reason := data.String("reason")
	if reason == "" {
		reason = "No reason provided"
	}

	if err := e.Client().Rest.RemoveMember(*guildID, user.ID, rest.WithReason(reason)); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Kick failed").
				WithDescription(fmt.Sprintf("Could not kick %s. Check that I have **Kick Members** permission and my role is above theirs.", user.Tag())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member kicked").
			WithDescription(fmt.Sprintf("**%s** has been kicked.", user.Tag())).
			AddField("Reason", reason, false).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
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

	if user.ID == e.User().ID {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You cannot ban yourself.")
	}

	reason := data.String("reason")
	if reason == "" {
		reason = "No reason provided"
	}

	if err := e.Client().Rest.AddBan(*guildID, user.ID, 0, rest.WithReason(reason)); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Ban failed").
				WithDescription(fmt.Sprintf("Could not ban %s. Check that I have **Ban Members** permission and my role is above theirs.", user.Tag())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member banned").
			WithDescription(fmt.Sprintf("**%s** has been banned.", user.Tag())).
			AddField("Reason", reason, false).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now()))
}

func handlePoll(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	question := data.String("question")
	options := make([]string, 0, 4)
	emojis := []string{"1\uFE0F\u20E3", "2\uFE0F\u20E3", "3\uFE0F\u20E3", "4\uFE0F\u20E3"}

	if o, ok := data.OptString("option1"); ok {
		options = append(options, o)
	}
	if o, ok := data.OptString("option2"); ok {
		options = append(options, o)
	}
	if o, ok := data.OptString("option3"); ok {
		options = append(options, o)
	}
	if o, ok := data.OptString("option4"); ok {
		options = append(options, o)
	}

	if len(options) < 2 {
		if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("You need at least 2 options for a poll.")); err != nil {
			log.Printf("failed to send poll error: %v", err)
		}
		return
	}
	for i, opt := range options {
		options[i] = strings.TrimSpace(opt)
		if options[i] == "" {
			if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("Poll options cannot be empty.")); err != nil {
				log.Printf("failed to send poll error: %v", err)
			}
			return
		}
		for j := 0; j < i; j++ {
			if strings.EqualFold(options[j], options[i]) {
				if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("Poll options must be unique.")); err != nil {
					log.Printf("failed to send poll error: %v", err)
				}
				return
			}
		}
	}

	question = strings.TrimSpace(question)
	if question == "" {
		if err := e.CreateMessage(discord.NewMessageCreate().WithEphemeral(true).WithContent("Poll question cannot be empty.")); err != nil {
			log.Printf("failed to send poll error: %v", err)
		}
		return
	}

	if err := e.DeferCreateMessage(false); err != nil {
		log.Printf("failed to defer poll: %v", err)
		return
	}

	var desc strings.Builder
	desc.Grow(len(question) + len(options)*16)
	desc.WriteString("## ")
	desc.WriteString(question)
	desc.WriteString("\n\n")

	for i, opt := range options {
		desc.WriteString(emojis[i])
		desc.WriteString(" ")
		desc.WriteString(opt)
		desc.WriteString("\n")
	}

	msg, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().AddEmbeds(
		discord.NewEmbed().
			WithTitle("\U0001F4CA Poll").
			WithDescription(desc.String()).
			WithColor(commands.CommandMenuAccent).
			WithFooterText(fmt.Sprintf("PulseKeep %s · Poll by %s", Version, e.User().EffectiveName())).
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

	member, err := e.Client().Rest.GetMember(*guildID, user.ID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Role update failed").
				WithDescription("Could not fetch member info. Make sure I can see them.").
				WithColor(commands.UtilityMenuAccent))
	}

	hasRole := false
	for _, rid := range member.RoleIDs {
		if rid == role.ID {
			hasRole = true
			break
		}
	}

	if hasRole {
		if err := e.Client().Rest.RemoveMemberRole(*guildID, user.ID, role.ID); err != nil {
			return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
				discord.NewEmbed().
					WithTitle("Role update failed").
					WithDescription("Could not remove role. Check that I have **Manage Roles** permission and my role is above theirs.").
					WithColor(commands.UtilityMenuAccent))
		}
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Role removed").
				WithDescription(fmt.Sprintf("Removed <@&%s> from **%s**.", role.ID, user.Tag())).
				WithColor(commands.UtilityMenuAccent).
				WithFooterText("PulseKeep " + Version + " · Utility").
				WithTimestamp(time.Now()))
	}

	if err := e.Client().Rest.AddMemberRole(*guildID, user.ID, role.ID); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Role update failed").
				WithDescription("Could not add role. Check that I have **Manage Roles** permission and my role is above theirs.").
				WithColor(commands.UtilityMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Role added").
			WithDescription(fmt.Sprintf("Added <@&%s> to **%s**.", role.ID, user.Tag())).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Utility").
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

	if err := e.Client().Rest.DeleteBan(*guildID, userID, rest.WithReason("Unbanned via /unban")); err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Unban failed").
				WithDescription("Could not unban that user. Check that the ID is valid and the user is currently banned.").
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("User unbanned").
			WithDescription(fmt.Sprintf("Successfully unbanned `<@%s>`.", userID)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
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
				WithDescription("Could not set slowmode. Check that I have **Manage Channels** permission.").
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
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
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
				WithDescription("Could not change nickname. Check that I have **Manage Nicknames** permission and my role is above theirs.").
				WithColor(commands.ModerationMenuAccent))
	}

	if nickname == "" {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Nickname reset").
				WithDescription(fmt.Sprintf("Reset nickname for **%s**.", user.Tag())).
				WithColor(commands.ModerationMenuAccent).
				WithFooterText("PulseKeep " + Version + " · Moderation").
				WithTimestamp(time.Now()))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Nickname changed").
			WithDescription(fmt.Sprintf("Changed **%s**'s nickname to **%s**.", user.Tag(), nickname)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
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
				WithDescription(fmt.Sprintf("Could not timeout %s. Check that I have **Moderate Members** permission and my role is above theirs.", user.Tag())).
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member timed out").
			WithDescription(fmt.Sprintf("**%s** has been timed out for **%d minutes**.", user.Tag(), duration)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now()))
}

func requireManageChannels(e *events.ApplicationCommandInteractionCreate) (bool, string) {
	member := e.Member()
	if member == nil {
		return false, "This command can only be used in a server."
	}
	if member.Permissions.Has(discord.PermissionAdministrator) {
		return true, ""
	}
	if !member.Permissions.Has(discord.PermissionManageChannels) {
		return false, "You need the **Manage Channels** permission to use this command."
	}
	botPerms := e.AppPermissions()
	if botPerms == nil || !botPerms.Has(discord.PermissionManageChannels) {
		return false, "I need the **Manage Channels** permission to lock or unlock channels."
	}
	return true, ""
}

func handleLock(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	ok, msg := requireManageChannels(e)
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent(msg)
	}

	channelID := e.Channel().ID()

	channel, err := e.Client().Rest.GetChannel(channelID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Lock failed").
				WithDescription("Could not fetch channel info. Make sure I have access to this channel.").
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
				WithDescription("Could not lock the channel. Check that I have **Manage Channels** permission and try again.").
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("\U0001F512 Channel locked").
			WithDescription(fmt.Sprintf("<#%s> has been locked.", channelID)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now()))
}

func handleUnlock(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server.")
	}

	ok, msg := requireManageChannels(e)
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent(msg)
	}

	channelID := e.Channel().ID()

	channel, err := e.Client().Rest.GetChannel(channelID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Unlock failed").
				WithDescription("Could not fetch channel info. Make sure I have access to this channel.").
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
				WithDescription("Could not unlock the channel. Check that I have **Manage Channels** permission and try again.").
				WithColor(commands.ModerationMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("\U0001F513 Channel unlocked").
			WithDescription(fmt.Sprintf("<#%s> has been unlocked.", channelID)).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now()))
}

func tipMessage() discord.MessageCreate {
	tips := []string{
		"Use `/daily` every day to build your streak! Each day adds bonus Pulses up to 750.",
		"Buy a **Fishing Rod** from `/shop` to unlock the `/fish` command for passive income.",
		"Buy an **Iron Pickaxe** from `/shop` to unlock the `/mine` command for ore profits.",
		"Use `/profile` to track your earnings, streaks, and item value.",
		"The `/daily` streak resets if you miss a day. Stay consistent!",
		"Buy a **Shield Token** and use it to protect your balance from robbers.",
		"Roll 46+ in `/gamble` to win! 36-45 is a push (refund). 85+ = 2x, 95+ = 4x, 100 = 10x!",
		"Use `/weekly` every 7 days for a big Pulses bonus with streak rewards.",
		"Sell unwanted items with `/sell` for a 60% refund.",
		"Use `/rich` to see the wealthiest members on the leaderboard.",
		"Multiple income sources = fastest growth: daily + work + fish/mine + gambling.",
		"Keep an eye on `/lottery` for the weekly jackpot draw!",
		"Play `/blackjack` on Easy difficulty for better odds against the dealer.",
		"Use `/inventory` to see all your items and their estimated total value.",
		"Run `/work` every 45 minutes for steady income between daily claims.",
		"The economy has a built-in 0.1% interest applied every 6 hours on your wallet balance.",
		"Buy a **XP Boost** from `/shop` to double your work earnings for 30 minutes.",
		"Use `/balance` with `public: true` to show off your Pulse balance in chat.",
		"Server admins can configure auto-mod rules from the web dashboard.",
		"Ticket support is available via `/ticketpanel` — users can open private support channels.",
	}
	tip := tips[rand.Intn(len(tips))]
	return discord.NewMessageCreate().
		WithEphemeral(true).
		AddEmbeds(discord.NewEmbed().
			WithTitle("PulseKeep Tip").
			WithDescription(tip).
			WithColor(commands.EconomyMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Tips").
			WithTimestamp(time.Now()))
}

func snowflakeID(i int) *snowflake.ID {
	id := snowflake.ID(i)
	return &id
}

func startStatusRotation(ctx context.Context, client *bot.Client, memCache *cache.Cache) {
	statuses := []struct {
		text func(gc, uc int64) string
		kind string
	}{
		{func(_, _ int64) string { return "/help to browse commands" }, "playing"},
		{func(gc, _ int64) string { return fmt.Sprintf("over %d servers", gc) }, "watching"},
		{func(_, _ int64) string { return "PulseKeep Economy" }, "playing"},
		{func(_, _ int64) string { return "/daily | /work | /slots" }, "playing"},
		{func(_, _ int64) string { return "/gamble | /fish | /mine" }, "playing"},
		{func(gc, uc int64) string { return fmt.Sprintf("%d servers · %d users", gc, uc) }, "watching"},
		{func(_, _ int64) string { return "/purge | /kick | /ban" }, "playing"},
		{func(_, _ int64) string { return "/slowmode | /lock | /timeout" }, "playing"},
		{func(_, _ int64) string { return "member activity" }, "watching"},
		{func(_, _ int64) string { return "/poll | /role | /announce" }, "playing"},
		{func(_, _ int64) string { return "PulseKeep " + Version }, "competing"},
		{func(_, _ int64) string { return "support tickets" }, "listening"},
		{func(_, _ int64) string { return "/shop | /rich | /weekly" }, "playing"},
		{func(_, _ int64) string { return "automod" }, "watching"},
	}

	ticker := time.NewTicker(2 * time.Minute)
	idx := 0

	applyStatus := func() {
		s := statuses[idx%len(statuses)]
		idx++
		gc, uc := int64(0), int64(0)
		if memCache != nil {
			gc = memCache.GuildCount.Load()
			uc = memCache.GetTotalUserCount()
		}
		text := s.text(gc, uc)
		var opts []gateway.PresenceOpt
		switch s.kind {
		case "watching":
			opts = []gateway.PresenceOpt{gateway.WithWatchingActivity(text), gateway.WithOnlineStatus(discord.OnlineStatusOnline)}
		case "listening":
			opts = []gateway.PresenceOpt{gateway.WithListeningActivity(text), gateway.WithOnlineStatus(discord.OnlineStatusOnline)}
		case "competing":
			opts = []gateway.PresenceOpt{gateway.WithCompetingActivity(text), gateway.WithOnlineStatus(discord.OnlineStatusOnline)}
		default:
			opts = []gateway.PresenceOpt{gateway.WithPlayingActivity(text), gateway.WithOnlineStatus(discord.OnlineStatusOnline)}
		}
		if err := client.SetPresence(ctx, opts...); err != nil {
			log.Printf("failed to set presence: %v", err)
		}
	}

	applyStatus()
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				applyStatus()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func isEconomyCommand(name string) bool {
	switch name {
	case "balance", "profile", "daily", "work", "pay", "coinflip", "rob", "shop", "buy", "inventory", "slots", "fish", "mine", "gamble", "sell", "use", "blackjack", "lottery", "lottery-claim", "rich", "weekly", "gift":
		return true
	}
	return false
}

func featureDisabledEmbed(feature, description string) discord.MessageCreate {
	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("🔒 " + feature + " Disabled").
			WithDescription(description).
			WithColor(commands.EconomyWarningAccent).
			WithFooterText("PulseKeep " + Version + " · Configuration").
			WithTimestamp(time.Now()))
}

func (b *Bot) handleWarn(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	reason, _ := data.OptString("reason")
	if reason == "" {
		reason = "No reason provided"
	}

	if b.db != nil {
		_, err := b.db.ExecContext(context.Background(),
			"INSERT INTO user_warnings (guild_id, user_id, moderator_id, reason) VALUES ($1, $2, $3, $4)",
			guildID.String(), user.ID.String(), e.User().ID.String(), reason)
		if err != nil {
			log.Printf("failed to save warning: %v", err)
		}
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("User Warned").
			WithDescription(fmt.Sprintf("**%s** has been warned.", user.Mention())).
			AddField("Reason", reason, false).
			AddField("Moderator", e.User().Mention(), true).
			AddField("User", user.Mention(), true).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleWarnings(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	if b.db == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Database is not available."))
		return
	}

	rows, err := b.db.QueryContext(context.Background(),
		"SELECT id, moderator_id, reason, created_at FROM user_warnings WHERE guild_id = $1 AND user_id = $2 ORDER BY created_at DESC",
		guildID.String(), user.ID.String())
	if err != nil {
		log.Printf("failed to query warnings: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to fetch warnings."))
		return
	}
	defer rows.Close()

	type warnEntry struct {
		ID          int
		ModeratorID string
		Reason      string
		CreatedAt   time.Time
	}
	var warns []warnEntry
	for rows.Next() {
		var w warnEntry
		if err := rows.Scan(&w.ID, &w.ModeratorID, &w.Reason, &w.CreatedAt); err != nil {
			log.Printf("failed to scan warning row: %v", err)
			continue
		}
		warns = append(warns, w)
	}

	if len(warns) == 0 {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Warnings").
				WithDescription(fmt.Sprintf("**%s** has no warnings.", user.Mention())).
				WithColor(commands.ModerationMenuAccent).
				WithFooterText("PulseKeep " + Version + " · Moderation").
				WithTimestamp(time.Now())))
		return
	}

	var list strings.Builder
	for i, w := range warns {
		if i >= 10 {
			list.WriteString(fmt.Sprintf("\n... and %d more", len(warns)-10))
			break
		}
		list.WriteString(fmt.Sprintf("**#%d** — <@%s> • %s\n└ %s\n", w.ID, w.ModeratorID, w.CreatedAt.Format("Jan 02"), w.Reason))
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle(fmt.Sprintf("Warnings — %s", user.EffectiveName())).
			WithDescription(list.String()).
			AddField("Total", fmt.Sprintf("%d warning(s)", len(warns)), true).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleClearWarns(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	if b.db == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Database is not available."))
		return
	}

	res, err := b.db.ExecContext(context.Background(),
		"DELETE FROM user_warnings WHERE guild_id = $1 AND user_id = $2",
		guildID.String(), user.ID.String())
	if err != nil {
		log.Printf("failed to delete warnings: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to clear warnings."))
		return
	}

	count, _ := res.RowsAffected()
	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Warnings Cleared").
			WithDescription(fmt.Sprintf("Cleared **%d** warning(s) for %s.", count, user.Mention())).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleMove(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	targetChan, ok := data.OptChannel("channel")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a voice channel."))
		return
	}

	if _, err := e.Client().Rest.UpdateMember(*guildID, user.ID, discord.MemberUpdate{ChannelID: &targetChan.ID}); err != nil {
		log.Printf("failed to move member: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to move the member. Make sure they are connected to a voice channel."))
		return
	}

	oldChannel := "Unknown (voice state not checked)"

	chanMention := fmt.Sprintf("<#%s>", targetChan.ID.String())

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member Moved").
			WithDescription(fmt.Sprintf("Moved **%s** to %s.", user.Mention(), chanMention)).
			AddField("From", oldChannel, true).
			AddField("To", chanMention, true).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleVCKick(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	if _, err := e.Client().Rest.UpdateMember(*guildID, user.ID, discord.MemberUpdate{ChannelID: nil}); err != nil {
		log.Printf("failed to disconnect member: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to disconnect the member from voice."))
		return
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Member Disconnected").
			WithDescription(fmt.Sprintf("Disconnected **%s** from voice chat.", user.Mention())).
			WithColor(commands.ModerationMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Moderation").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleTicketCmd(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	guildID := e.GuildID()
	if guildID == nil {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command can only be used in a server."))
		return
	}

	channelID := e.Channel().ID()
	ch, err := e.Client().Rest.GetChannel(channelID)
	if err != nil || !strings.HasPrefix(ch.Name(), "ticket-") {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command must be used inside a ticket channel."))
		return
	}

	guildChannel, ok := ch.(discord.GuildTextChannel)
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("This command must be used inside a ticket channel."))
		return
	}

	sub := ""
	if data.SubCommandName != nil {
		sub = *data.SubCommandName
	}
	switch sub {
	case "add":
		b.handleTicketAdd(e, data, guildChannel)
	case "remove":
		b.handleTicketRemove(e, data, guildChannel)
	case "close":
		b.respond(e, discord.NewMessageCreate().WithContent("Ticket closed. Channel will be deleted shortly."))
		b.doTicketClose(e.Client().Rest, e.Channel().ID())
	case "rename":
		b.handleTicketRenameCmd(e, data, guildChannel)
	default:
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Unknown ticket subcommand."))
	}
}

func (b *Bot) handleTicketAdd(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData, ch discord.GuildTextChannel) {
	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	overwrites := ch.PermissionOverwrites()
	overwrites = append(overwrites, discord.MemberPermissionOverwrite{
		UserID: user.ID,
		Allow:  discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory,
	})

	name := ch.Name()
	owSlice := []discord.PermissionOverwrite(overwrites)
	if _, err := e.Client().Rest.UpdateChannel(ch.ID(), discord.GuildTextChannelUpdate{
		Name:                 &name,
		PermissionOverwrites: &owSlice,
	}); err != nil {
		log.Printf("failed to add user to ticket: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to add user to ticket."))
		return
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("User Added").
			WithDescription(fmt.Sprintf("Added %s to the ticket.", user.Mention())).
			WithColor(commands.TicketMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Tickets").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleTicketRemove(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData, ch discord.GuildTextChannel) {
	user, ok := data.OptUser("user")
	if !ok {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a user."))
		return
	}

	var filtered []discord.PermissionOverwrite
	for _, ow := range ch.PermissionOverwrites() {
		if ow.Type() == discord.PermissionOverwriteTypeMember && ow.ID() == user.ID {
			continue
		}
		filtered = append(filtered, ow)
	}

	if len(filtered) == len(ch.PermissionOverwrites()) {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("That user does not have access to this ticket."))
		return
	}

	name := ch.Name()
	if _, err := e.Client().Rest.UpdateChannel(ch.ID(), discord.GuildTextChannelUpdate{
		Name:                 &name,
		PermissionOverwrites: &filtered,
	}); err != nil {
		log.Printf("failed to remove user from ticket: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to remove user from ticket."))
		return
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("User Removed").
			WithDescription(fmt.Sprintf("Removed %s from the ticket.", user.Mention())).
			WithColor(commands.TicketMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Tickets").
			WithTimestamp(time.Now())))
}

func (b *Bot) handleTicketRenameCmd(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData, ch discord.GuildTextChannel) {
	newName, ok := data.OptString("name")
	if !ok || newName == "" {
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("You must provide a new name."))
		return
	}

	if _, err := e.Client().Rest.UpdateChannel(ch.ID(), discord.GuildTextChannelUpdate{Name: &newName}); err != nil {
		log.Printf("failed to rename ticket channel: %v", err)
		b.respond(e, discord.NewMessageCreate().WithEphemeral(true).WithContent("Failed to rename the ticket channel."))
		return
	}

	b.respond(e, discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Ticket Renamed").
			WithDescription(fmt.Sprintf("Renamed to **%s**.", newName)).
			WithColor(commands.TicketMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Tickets").
			WithTimestamp(time.Now())))
}

func serverIconMessage(e *events.ApplicationCommandInteractionCreate) discord.MessageCreate {
	guildID := e.GuildID()
	if guildID == nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Server Icon").
				WithDescription("Server icon is only available inside a guild.").
				WithColor(commands.UtilityMenuAccent))
	}

	guild, err := e.Client().Rest.GetGuild(*guildID, true)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Server Icon").
				WithDescription("Could not fetch server info.").
				WithColor(commands.UtilityMenuAccent))
	}

	iconURL := guild.IconURL()
	if iconURL == nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Server Icon").
				WithDescription("This server does not have a custom icon.").
				WithColor(commands.UtilityMenuAccent))
	}

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle(fmt.Sprintf("🏠 %s — Server Icon", guild.Name)).
			WithColor(commands.UtilityMenuAccent).
			WithImage(*iconURL).
			WithDescription(fmt.Sprintf("[Open in browser](%s)", *iconURL)).
			WithFooterText("PulseKeep " + Version + " · Utility").
			WithTimestamp(time.Now()))
}

func roleInfoMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	role, ok := data.OptRole("role")
	if !ok {
		return discord.NewMessageCreate().WithEphemeral(true).WithContent("You must specify a role.")
	}

	colorStr := fmt.Sprintf("#%06X", role.Color)
	createdAt := role.ID.Time().Format("Jan 02, 2006")

	permList := "None"
	if role.Permissions != 0 {
		permNames := []string{}
		if role.Permissions.Has(discord.PermissionAdministrator) {
			permNames = append(permNames, "Administrator")
		} else {
			if role.Permissions.Has(discord.PermissionManageGuild) { permNames = append(permNames, "Manage Server") }
			if role.Permissions.Has(discord.PermissionManageRoles) { permNames = append(permNames, "Manage Roles") }
			if role.Permissions.Has(discord.PermissionManageChannels) { permNames = append(permNames, "Manage Channels") }
			if role.Permissions.Has(discord.PermissionKickMembers) { permNames = append(permNames, "Kick Members") }
			if role.Permissions.Has(discord.PermissionBanMembers) { permNames = append(permNames, "Ban Members") }
			if role.Permissions.Has(discord.PermissionModerateMembers) { permNames = append(permNames, "Timeout Members") }
			if role.Permissions.Has(discord.PermissionMentionEveryone) { permNames = append(permNames, "Mention Everyone") }
			if role.Permissions.Has(discord.PermissionManageMessages) { permNames = append(permNames, "Manage Messages") }
		}
		if len(permNames) > 0 {
			permList = strings.Join(permNames, ", ")
		}
	}

	embed := discord.NewEmbed().
		WithTitle(fmt.Sprintf("Role — %s", role.Name)).
		WithColor(int(role.Color)).
		AddField("Role ID", role.ID.String(), true).
		AddField("Color", colorStr, true).
		AddField("Position", fmt.Sprintf("%d", role.Position), true).
		AddField("Mentionable", fmt.Sprintf("%t", role.Mentionable), true).
		AddField("Displayed Separately", fmt.Sprintf("%t", role.Hoist), true).
		AddField("Created", createdAt, true).
		AddField("Key Permissions", permList, false).
		WithFooterText("PulseKeep " + Version + " · Utility").
		WithTimestamp(time.Now())

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(embed)
}

func channelInfoMessage(e *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) discord.MessageCreate {
	guildID := e.GuildID()

	var channelID snowflake.ID
	if ch, ok := data.OptChannel("channel"); ok {
		channelID = ch.ID
	} else {
		channelID = e.Channel().ID()
	}

	ch, err := e.Client().Rest.GetChannel(channelID)
	if err != nil {
		return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
			discord.NewEmbed().
				WithTitle("Channel Info").
				WithDescription("Could not fetch channel information.").
				WithColor(commands.UtilityMenuAccent))
	}

	typeName := map[discord.ChannelType]string{
		discord.ChannelTypeGuildText:  "Text",
		discord.ChannelTypeGuildVoice: "Voice",
		discord.ChannelTypeGuildCategory: "Category",
		discord.ChannelTypeGuildNews: "Announcement",
		discord.ChannelTypeGuildForum: "Forum",
		discord.ChannelTypeGuildStageVoice: "Stage",
	}[ch.Type()]
	if typeName == "" {
		typeName = "Unknown"
	}

	createdAt := ch.ID().Time().Format("Jan 02, 2006")

	embed := discord.NewEmbed().
		WithTitle(fmt.Sprintf("#%s — Channel Info", ch.Name())).
		WithDescription(fmt.Sprintf("Information about <#%s>.", ch.ID())).
		WithColor(commands.UtilityMenuAccent).
		AddField("Type", typeName, true).
		AddField("Channel ID", ch.ID().String(), true).
		AddField("Created", createdAt, true)

	if guildID != nil {
		embed = embed.AddField("Server", fmt.Sprintf("<@%s>", guildID.String()), true)
	}

	embed = embed.WithFooterText("PulseKeep " + Version + " · Utility").
		WithTimestamp(time.Now())

	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(embed)
}

func inviteMessage() discord.MessageCreate {
	return discord.NewMessageCreate().WithEphemeral(true).AddEmbeds(
		discord.NewEmbed().
			WithTitle("Invite PulseKeep").
			WithDescription("Add PulseKeep to your server or join the support community.").
			AddField("Invite Bot", "[Click here](https://discord.com/oauth2/authorize?client_id=1507498795569512598&permissions=8&scope=bot%20applications.commands)", false).
			AddField("Support Server", "[Click here](https://discord.gg/pulsekeep)", false).
			AddField("Top.gg", "[Vote for us](https://top.gg/bot/1507498795569512598/vote)", false).
			WithColor(commands.UtilityMenuAccent).
			WithFooterText("PulseKeep " + Version + " · Utility").
			WithTimestamp(time.Now()))
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
