package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/apihandling"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func ServerStatusCommand(gsm *apihandling.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the server status message!")
		return nil
	}
}

func StartServerCommand(gsm *apihandling.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the start server message!")
		return nil
	}
}

func StopServerCommand(gsm *apihandling.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the stop server message!")
		return nil
	}
}
