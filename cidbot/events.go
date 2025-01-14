package cidbot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func OnReady(session *discordgo.Session, readyEvent *discordgo.Ready) {
	err := session.UpdateCustomStatus("Background checking...")
	if err != nil {
		log.Fatalf("could not set custom status: %s", err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(readyEvent.Application.ID, "", Commands)
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(readyEvent.Application.ID, Configuration.AdminServerID, ServerCommands)
	if err != nil {
		log.Fatalf("could not register admin commands: %s", err)
	}

	_, err = os.OpenFile("./whitelist", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("couldn't create the whitelist file: %s", err)
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
	case "whitelist":
		WhitelistCommand(session, interaction.Interaction, ParseCommandOptions(data.Options))
	default:
		log.Printf("invalid command \"%s\" selected", data.Name)
	}

}
