package bot

import (
	"context"
	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type Bot struct {
	Client *bot.Client
}

func New(token string) *Bot {
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
		}),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			log.Printf("Bot is ready as %s#%s", e.User.Username, e.User.Discriminator)
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
