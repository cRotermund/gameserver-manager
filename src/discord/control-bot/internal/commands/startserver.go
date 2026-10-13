package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func StartServerCommand(ctx *messaging.Context) error {
	ctx.Client.Reply("Received the start server message!")
	return nil
}
