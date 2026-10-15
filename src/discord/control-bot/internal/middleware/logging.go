package middleware

import (
	"time"

	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
)

func LoggingMiddleware(next messaging.Handler) messaging.Handler {
	return func(ctx *messaging.RoutingContext) error {
		start := time.Now()

		ctx.Logger.Debug("handling message")

		err := next(ctx)

		logCtx := []any{
			"duration_ms", time.Since(start).Milliseconds(),
		}

		if err != nil {
			ctx.Logger.Error(
				"message handling failed",
				append(logCtx, "err", err)...)

		} else {
			ctx.Logger.Info(
				"message handling completed",
				logCtx...)
		}

		return err
	}
}
