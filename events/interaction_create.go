package events

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func InteractionCreate(discord *discordgo.Session, event *discordgo.InteractionCreate) {
	slog.Debug("interaction received", "interaction", event.Interaction)
	switch event.Type {
	case discordgo.InteractionApplicationCommand:
		commands.Executed(discord, event.Interaction)
	default:
		slog.Warn("unhandled interaction type", "type", event.Type, "id", event.Interaction.ID)
	}
}
