package events

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Ready(session *discordgo.Session, readyEvent *discordgo.Ready) {
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

	file, err := os.OpenFile(Configuration.WhitelistPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("couldn't create the whitelist file: %s", err)
	}
	defer file.Close()
}


