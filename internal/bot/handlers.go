package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/watispro/pulsekeep/internal/bot/commands"
	"github.com/watispro/pulsekeep/internal/db"
)

type Bot struct {
	Client bot.Client
	DB     *db.Database
}

func New(token string, database *db.Database) *Bot {
	cmdDefs := commands.GetCommands(database)

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
				if err := cmd.Handler(ctx); err != nil {
					log.Printf("Error handling command %s: %v", data.CommandName(), err)
				}
			}
		}),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			log.Printf("Bot is ready as %s#%s", e.User.Username, e.User.Discriminator)
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
		}),
	)
	if err != nil {
		log.Fatalf("error while building disgo instance: %s", err)
	}
	return &Bot{Client: *client, DB: database}
}

func (b *Bot) Start(ctx context.Context) error {
	return b.Client.OpenShardManager(ctx)
}

func (b *Bot) Stop(ctx context.Context) {
	b.Client.Close(ctx)
}
