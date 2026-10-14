package main

import (
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
	gsmClient      *controlplaneclient.GsmApiClient
	router         *messaging.Router
	discordSession *discordgo.Session
}

func (s *Runtime) destroy() {
	s.gsmClient.Close()
	s.discordSession.Close()
}

func main() {
	logger := slog.Default()

	config := config.NewBotConfig(logger)
	err := config.Load()
	if err != nil {
		logger.Error("Fatal error loading configuration, shutting down")
		os.Exit(1)
	}

	services := initRuntime(config, logger)

	services.discordSession.AddHandler(
		func(session *discordgo.Session, message *discordgo.MessageCreate) {
			onMessageCreate(
				session,
				message,
				services.router,
				logger)
		})

	services.discordSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds)

	err = services.discordSession.Open()
	if err != nil {
		logger.Error("Error opening connection to Discord", "err", err)
		os.Exit(1)
	}

	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Shutting down gracefully...")

	services.destroy()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, router *messaging.Router, logger *slog.Logger) {
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

	if !isMentioned {
		return
	}

	msgClient := discord.NewDiscordMessagingClient(logger, s, m.Message)
	resp, err := router.Dispatch(s, m.Message, msgClient)

	if err != nil {
		slog.Error("Error sending message", "err", err)
		sysErrorResponse(s, m)
		return
	}

	if !resp.WasRouted {
		noRouteResponse(s, m)
		return
	}
}

func noRouteResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	defaultResp := fmt.Sprintf("Hello <@%s>! You rang?  Give me a command, prefixed with '!', like !help to learn more.", m.Author.ID)
	_, err := s.ChannelMessageSend(m.ChannelID, defaultResp)

	if err != nil {
		slog.Error("Error sending message", "err", err)
		sysErrorResponse(s, m)
	}
}

func sysErrorResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Something went wrong trying to process your message. I'm sorry.")
}

func initRuntime(cfg *config.DiscordBotConfig, logger *slog.Logger) *Runtime {
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
	router.Use(middleware.LoggingMiddelware)
	router.Handle("startserver", commands.StartServerCommand(gsmApi))
	router.Handle("stopserver", commands.StopServerCommand(gsmApi))
	router.Handle("serverstatus", commands.ServerStatusCommand(gsmApi))
	router.Handle("listservers", commands.ListServersCommand(gsmApi))
	router.Handle("help", commands.HelpCommand())

	//Connect to discord and capture the session information
	sess, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		logger.Error("Error creating Discord session", "err", err)
		os.Exit(1)
	}

	//Return these in a struct for later reference
	return &Runtime{
		gsmClient:      gsmApi,
		router:         router,
		discordSession: sess,
	}
}
