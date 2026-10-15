package discord

import (
	"context"
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

func (c *DiscordMessagingClient) Reply(ctx context.Context, msg string) error {
	_, err := c.session.ChannelMessageSend(
		c.rootMsg.ChannelID,
		msg,
		discordgo.WithContext(ctx))

	if err != nil {
		c.logger.Error(
			"Error sending message to discord.",
			"err", err)

		return err
	}

	return nil
}
