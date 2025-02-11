package events

import (
	"log/slog"

	"github.com/RobloxUSArmyCID/CIDBot/commands"
	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(config *config.Config) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(discord *discordgo.Session, event *discordgo.InteractionCreate) {
		slog.Debug("interaction received", "interaction", event.Interaction)
		switch event.Type {
		case discordgo.InteractionApplicationCommand:
			commands.Executed(discord, event.Interaction, config)
		default:
			slog.Warn("unhandled interaction type", "type", event.Type, "id", event.Interaction.ID)
		}

	}
}
