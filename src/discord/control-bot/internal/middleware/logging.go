package middleware

import (
	"time"

	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func LoggingMiddelware(next messaging.Handler) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		start := time.Now()

		ctx.Logger.Info("Message handling start")

		err := next(ctx)

		logCtx := []any{
			"duration_ms", time.Since(start).Milliseconds(),
		}

		if err != nil {
			ctx.Logger.Error(
				"Message handling failed with error",
				append(logCtx, "err", err)...)

		} else {
			ctx.Logger.Info(
				"Message handling complete",
				logCtx...)
		}

		return err
	}
}
