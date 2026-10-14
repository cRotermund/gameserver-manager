package messaging

type MessagingClient interface {
	Reply(msg string)
}
