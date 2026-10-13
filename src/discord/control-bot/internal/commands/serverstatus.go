package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func ServerStatusCommand(ctx *messaging.Context) error {
	ctx.Client.Reply("Received the server status message!")
	return nil
}
