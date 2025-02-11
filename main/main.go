package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/RobloxUSArmyCID/CIDBot/events"
	"github.com/bwmarrin/discordgo"
)

func handleInterrupt() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func main() {
	config, err := config.Parse()
	if err != nil {
		slog.Error("could not parse config", "err", err)
		return
	}

	var logger *slog.Logger

	if config.IsDevelopment {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)

	slog.Info("starting bot")

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		slog.Error("could not create discord session", "err", err)
		return
	}

	defer func() {
		slog.Info("closing discord session")
		err := discord.Close()
		if err != nil {
			slog.Warn("couldn't close session gracefully", "err", err)
		}
	}()

	discord.AddHandler(events.Ready(config))
	discord.AddHandler(events.InteractionCreate(config))

	slog.Info("opening discord session")
	err = discord.Open()
	if err != nil {
		slog.Error("could not open discord session", "err", err)
		return
	}

	handleInterrupt()
}
