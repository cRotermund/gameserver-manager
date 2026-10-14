package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
	"github.com/cRotermund/gameserver-manager/src/pkg/controlplaneclient"
)

func ServerStatusCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the server status message!")
		return nil
	}
}

func StartServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the start server message!")
		return nil
	}
}

func StopServerCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the stop server message!")
		return nil
	}
}

func ListServersCommand(gsm *controlplaneclient.GsmApiClient) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Client.Reply("Received the list servers command!")
		return nil
	}
}
