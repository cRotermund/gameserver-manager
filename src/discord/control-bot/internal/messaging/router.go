package messaging

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name    string
	Args    []string
	RawText string
}

type Context struct {
	Logger  *slog.Logger
	Sender  string
	Command *Command
	Client  MessagingClient
}

type DispatchResult struct {
	WasRouted bool
}

type Handler func(ctx *Context) error

type Middleware func(next Handler) Handler

type Router struct {
	prefix     string
	logger     *slog.Logger
	middleware []Middleware
	handlers   map[string]Handler
}

func NewRouter(prefix string, logger *slog.Logger) *Router {
	return &Router{
		prefix:   prefix,
		logger:   logger,
		handlers: make(map[string]Handler),
	}
}

func (r *Router) Handle(command string, handler Handler) {
	r.handlers[strings.ToLower(command)] = handler
}

func (r *Router) Use(m Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *Router) Dispatch(session *discordgo.Session, message *discordgo.Message, msgClient MessagingClient) (result DispatchResult, err error) {
	res := DispatchResult{
		WasRouted: false,
	}

	sanitized := cleanBotMention(message.Content, session)

	//Check to see if it matches our command pattern ("!")
	if !strings.HasPrefix(sanitized, r.prefix) {
		return res, nil
	}

	//Parse out the fields and the command content
	sanitized = strings.TrimPrefix(sanitized, r.prefix)
	fields := strings.Fields(sanitized)

	if len(fields) == 0 {
		return res, nil
	}

	cmdName := strings.ToLower(fields[0])
	args := fields[1:]

	handler, exists := r.handlers[cmdName]
	if !exists {
		return res, nil
	}

	res.WasRouted = true

	ctx := &Context{
		Logger: r.logger,
		Sender: message.Author.ID,
		Client: msgClient,
		Command: &Command{
			Name:    cmdName,
			Args:    args,
			RawText: message.Content,
		},
	}

	// Build the execution chain backwards, so that the first
	// func reference is the first executed
	chain := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		chain = r.middleware[i](chain)
	}

	err = chain(ctx)

	if err != nil {
		return res, err
	}

	return res, nil
}

func cleanBotMention(content string, s *discordgo.Session) string {
	if s.State.User == nil {
		return strings.TrimSpace(content)
	}

	botID := s.State.User.ID

	// Discord uses two variants for user/bot mentions
	mentionStandard := "<@" + botID + ">"
	mentionNick := "<@!" + botID + ">"

	trimmed := content
	if strings.HasPrefix(trimmed, mentionStandard) {
		trimmed = strings.TrimPrefix(trimmed, mentionStandard)
	} else if strings.HasPrefix(trimmed, mentionNick) {
		trimmed = strings.TrimPrefix(trimmed, mentionNick)
	}

	// Remove any extra leading spaces left after stripping the tag (e.g., "@Bot hello" -> "hello")
	return strings.TrimSpace(trimmed)
}
