package events

func InteractionCreate(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
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