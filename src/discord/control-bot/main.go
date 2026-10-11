package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
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

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session: $v", err)
		os.Exit(1)
	}

	sess.AddHandler(messageCreate)

	sess.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds)

	err = sess.Open()
	if err != nil {
		logger.Error("Error opening connection to Discord $v", err)
		os.Exit(1)
	}

	fmt.Println("Bot is now running.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Shutting down gracefully...")
	sess.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//ignore messages from self, to avoid infinite loops
	if m.Author.ID == s.State.User.ID {
		return
	}

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

	responseMessage := fmt.Sprintf("Hello <@%s>! You rang?", m.Author.ID)
	_, err := s.ChannelMessageSend(m.ChannelID, responseMessage)

	if err != nil {
		slog.Error("Error sending message: %v", err)
	}
}
