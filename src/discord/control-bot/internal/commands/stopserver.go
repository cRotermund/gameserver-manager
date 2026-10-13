package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func StopServerCommand(ctx *messaging.Context) error {
	ctx.Client.Reply("Received the stop server message!")
	return nil
}
