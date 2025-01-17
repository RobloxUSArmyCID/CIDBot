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
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("starting bot")
	slog.Debug("parsing config")
	if err := config.Parse(); err != nil {
		slog.Error("could not parse config", "err", err)
		return
	}

	discord, err := discordgo.New("Bot " + config.Configuration.Token)
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

	discord.AddHandler(events.Ready)
	discord.AddHandler(events.InteractionCreate)

	slog.Info("opening discord session")
	err = discord.Open()
	if err != nil {
		slog.Error("could not open discord session", "err", err)
		return
	}

	handleInterrupt()
}
