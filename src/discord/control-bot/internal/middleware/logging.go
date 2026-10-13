package middleware

import (
	"strings"

	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func LoggingMiddelware(next messaging.Handler) messaging.Handler {
	return func(ctx *messaging.Context) error {
		ctx.Logger.Info("Message handling start",
			"sender", ctx.Sender,
			"command", ctx.Command.Name,
			"args", strings.Join(ctx.Command.Args, ","),
			"raw", ctx.Command.RawText)

		err := next(ctx)

		ctx.Logger.Info("Message hhandling complete",
			"sender", ctx.Sender,
			"command", ctx.Command.Name)

		return err
	}
}
