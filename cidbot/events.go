package cidbot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(session *discordgo.Session, readyEvent *discordgo.Ready) {
	_, err := session.ApplicationCommandBulkOverwrite(readyEvent.Application.ID, "", Commands)

	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}
}

func OnInteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := interaction.ApplicationCommandData()

	DeferInteraction(session, interaction.Interaction)

	switch data.Name {
	case "bgcheck":
		BackgroundCheckCommand(session, interaction.Interaction, ParseCommandOptions(data.Options))
	default:
		log.Fatalf("invalid command \"%s\" selected", data.Name)
	}

}
