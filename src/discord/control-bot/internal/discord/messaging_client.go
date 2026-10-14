package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

type DiscordMessagingClient struct {
	logger  *slog.Logger
	session *discordgo.Session
	rootMsg *discordgo.Message
}

func NewDiscordMessagingClient(
	logger *slog.Logger,
	session *discordgo.Session,
	rootMsg *discordgo.Message) messaging.MessagingClient {

	return &DiscordMessagingClient{
		logger:  logger,
		session: session,
		rootMsg: rootMsg,
	}
}

func (c *DiscordMessagingClient) Reply(msg string) {
	c.session.ChannelMessageSend(c.rootMsg.ChannelID, msg)
}
