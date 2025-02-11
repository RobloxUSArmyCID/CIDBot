package events

import (
	"log/slog"
	"os"

	"github.com/RobloxUSArmyCID/CIDBot/commands"
	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/bwmarrin/discordgo"
)

func Ready(config *config.Config) func(*discordgo.Session, *discordgo.Ready) {
	return func(discord *discordgo.Session, event *discordgo.Ready) {
		status := "Background checking..."

		slog.Debug("updating status", "status", status)
		err := discord.UpdateCustomStatus(status)
		if err != nil {
			slog.Error("could not set custom status", "err", err)
			os.Exit(1)
		}

		slog.Debug("registering global commands")
		_, err = discord.ApplicationCommandBulkOverwrite(event.Application.ID, "", commands.Commands)
		if err != nil {
			slog.Error("could not register global commands", "err", err)
			os.Exit(1)
		}

		slog.Debug("registering server commands")
		_, err = discord.ApplicationCommandBulkOverwrite(event.Application.ID, config.AdminServerID, commands.AdminServerCommands)
		if err != nil {
			slog.Error("could not register server commands", "err", err)
			os.Exit(1)
		}

		slog.Debug("creating whitelist file")
		file, err := os.OpenFile(config.WhitelistPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			slog.Error("couldn't create the whitelist file", "err", err)
			os.Exit(1)
		}

		slog.Info("ready")
		defer file.Close()
	}
}
