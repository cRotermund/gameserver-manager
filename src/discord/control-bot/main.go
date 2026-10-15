package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/commands"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/config"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/discord"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/middleware"
	"github.com/cRotermund/gameserver-manager/src/pkg/controlplaneclient"
	"resty.dev/v3"
)

type Runtime struct {
	ctx            context.Context
	gsmClient      *controlplaneclient.GsmApiClient
	router         *messaging.Router
	discordSession *discordgo.Session
}

func (s *Runtime) destroy() {
	s.gsmClient.Close()
	s.discordSession.Close()
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
	slog.SetDefault(logger)

	config := config.NewBotConfig(logger)
	err := config.Load()
	if err != nil {
		logger.Error(
			"failed to load configuration, shutting down",
			"err", err)
		os.Exit(1)
	}

	runtime := initRuntime(ctx, config, logger)

	runtime.discordSession.AddHandler(
		func(session *discordgo.Session, message *discordgo.MessageCreate) {
			onMessageCreate(
				runtime.ctx,
				session,
				message,
				runtime.router,
				logger)
		})

	runtime.discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds)

	err = runtime.discordSession.Open()
	if err != nil {
		logger.Error("failed to open Discord connection", "err", err)
		os.Exit(1)
	}

	logger.Info("bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	logger.Info("shutting down gracefully")

	runtime.destroy()
}

func onMessageCreate(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, router *messaging.Router, logger *slog.Logger) {
	// ignore messages from self, to avoid infinite loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	// ensure we have been tagged - otherwise ignore.
	isMentioned := false
	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID {
			isMentioned = true
			break
		}
	}

	// ignore, if there is no metntion.
	if !isMentioned {
		return
	}

	msgClient := discord.NewDiscordMessagingClient(logger, s, m.Message)
	resp, err := router.Dispatch(ctx, s, m.Message, msgClient)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			logger.Debug("message handling cancelled")
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			logger.Error(
				"message handling timed out",
				"err", err)
			sysErrorResponse(ctx, s, m)
			return
		}

		logger.Error("failed to handle message", "err", err)
		sysErrorResponse(ctx, s, m)
		return
	}

	if !resp.WasRouted {
		noRouteResponse(ctx, logger, s, m)
		return
	}
}

func noRouteResponse(ctx context.Context, logger *slog.Logger, s *discordgo.Session, m *discordgo.MessageCreate) {
	defaultResp := fmt.Sprintf("Hello <@%s>! You rang?  Give me a command, prefixed with '!', like !help to learn more.", m.Author.ID)

	_, err := s.ChannelMessageSend(
		m.ChannelID,
		defaultResp,
		discordgo.WithContext(ctx))

	if err != nil {
		logger.Error("failed to send no-route response", "err", err)

		if !errors.Is(err, context.Canceled) {
			sysErrorResponse(ctx, s, m)
		}
	}
}

func sysErrorResponse(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(
		m.ChannelID,
		"Something went wrong trying to process your message. I'm sorry.",
		discordgo.WithContext(ctx))
}

func initRuntime(ctx context.Context, cfg *config.DiscordBotConfig, logger *slog.Logger) *Runtime {
	//Initialize the resty protocol client for the gsm api
	restyClient := resty.New()
	restyClient.SetTimeout(time.Duration(cfg.GsmApiConfig.TimeoutMs) * time.Millisecond)

	//Initialize the gsm api handler
	gsmApi := controlplaneclient.NewClient(
		logger,
		cfg.GsmApiConfig.Host,
		cfg.GsmApiConfig.Port,
		restyClient)

	//Initialize the router, the middleware, and the command handlers
	router := messaging.NewRouter("!", logger)
	router.Use(middleware.LoggingMiddleware)
	router.Handle("startserver", commands.StartServerCommand(gsmApi))
	router.Handle("stopserver", commands.StopServerCommand(gsmApi))
	router.Handle("serverstatus", commands.ServerStatusCommand(gsmApi))
	router.Handle("listservers", commands.ListServersCommand(gsmApi))
	router.Handle("help", commands.HelpCommand())

	//Connect to discord and capture the session information
	sess, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		logger.Error("failed to create Discord session", "err", err)
		os.Exit(1)
	}

	//Return these in a struct for later reference
	return &Runtime{
		ctx:            ctx,
		gsmClient:      gsmApi,
		router:         router,
		discordSession: sess,
	}
}
