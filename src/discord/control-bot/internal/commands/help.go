package commands

import (
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func HelpCommand() messaging.Handler {
	return func(ctx *messaging.Context) error {
		message := `Hi, I'm your friendly game server control bot.  Use me to manage the gameservers used in this discord.
		
		Usage:
			@{tagme} ![command] [flags]

		Available Commands:
			startserver			Start the game server
			stopserver			Stop the game server
			status				Show the status of the game server
			help				Display this help message

		Use "@{tagme} ![command] --help" for more information about a specific command.`

		ctx.Client.Reply(message)
		return nil
	}
}
