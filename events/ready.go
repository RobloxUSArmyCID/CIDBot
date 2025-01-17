package events

import (
	"log/slog"
	"os"

	"github.com/RobloxUSArmyCID/CIDBot/commands"
	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/bwmarrin/discordgo"
)

func Ready(discord *discordgo.Session, event *discordgo.Ready) {
	err := discord.UpdateCustomStatus("Background checking...")
	if err != nil {
		slog.Error("could not set custom status", "err", err)
		os.Exit(1)
	}

	_, err = discord.ApplicationCommandBulkOverwrite(event.Application.ID, "", commands.Commands)
	if err != nil {
		slog.Error("could not register commands", "err", err)
		os.Exit(1)
	}

	_, err = discord.ApplicationCommandBulkOverwrite(event.Application.ID, config.Configuration.AdminServerID, commands.ServerCommands)
	if err != nil {
		slog.Error("could not register admin commands", "err", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(config.Configuration.WhitelistPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("couldn't create the whitelist file", "err", err)
		os.Exit(1)
	}
	defer file.Close()
}
