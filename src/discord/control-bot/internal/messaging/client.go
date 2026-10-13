package messaging

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

type MessagingClient interface {
	Reply(msg string)
}

type DiscordMessagingClient struct {
	logger  *slog.Logger
	session *discordgo.Session
	rootMsg *discordgo.Message
}

func NewDiscordMessagingClient(logger *slog.Logger, session *discordgo.Session, rootMsg *discordgo.Message) MessagingClient {

	return &DiscordMessagingClient{
		logger:  logger,
		session: session,
		rootMsg: rootMsg,
	}
}

func (c *DiscordMessagingClient) Reply(msg string) {
	c.session.ChannelMessageSend(c.rootMsg.ChannelID, msg)
}
