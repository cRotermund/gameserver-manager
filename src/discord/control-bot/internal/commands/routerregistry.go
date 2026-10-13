package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func RegisterCommands(router *messaging.Router) {
	router.Handle("status", ServerStatusCommand)
	router.Handle("startserver", StartServerCommand)
	router.Handle("stopserver", StopServerCommand)
	router.Handle("help", HelpCommand)
}
