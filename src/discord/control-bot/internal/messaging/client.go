package messaging

import "context"

type MessagingClient interface {
	Reply(ctx context.Context, msg string) error
}
