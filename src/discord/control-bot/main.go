package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/commands"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/messaging"
	"github.com/cRotermund/gameserver-manager/src/discord/control-bot/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.Default()

	// Load the .env file. By default, it looks for a file named ".env" in the current directory.
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using existing environment", "Error", err)
	}

	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		logger.Error("Error: DISCORD_TOKEN environment variable is required")
		os.Exit(1)
	}

	//Initialize the router, the middleware, and the command handlers
	router := messaging.NewRouter("!", logger)
	router.Use(middleware.LoggingMiddelware)
	commands.RegisterCommands(router)

	//Connect to discord
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session", "err", err)
		os.Exit(1)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		onMessageCreate(s, m, router)
	})

	sess.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds)

	err = sess.Open()
	if err != nil {
		logger.Error("Error opening connection to Discord", "err", err)
		os.Exit(1)
	}

	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Shutting down gracefully...")
	sess.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, r *messaging.Router) {
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

	resp, err := r.Dispatch(s, m.Message)

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
