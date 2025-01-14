package cidbot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func WhitelistCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	subcommand := interaction.ApplicationCommandData().Options[0]
	subcommandOptions := ParseCommandOptions(subcommand.Options)
	switch subcommand.Name {
	case "add":
		addCommand(session, interaction, subcommandOptions)
	case "view":
		viewCommand(session, interaction, subcommandOptions)
	case "remove":
		removeCommand(session, interaction, subcommandOptions)
	default:
		log.Printf("incorrect whitelist command used: %s", subcommand.Name)
	}
}

func addCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	userID := options["user_id"].StringValue()
	userIDBytes := []byte(userID + "\n")

	user, err := session.User(userID)
	if err != nil {
		InteractionFailed(session, interaction, "user doesn't exist or another error has occured", err)
		return
	}

	file, err := os.OpenFile("./whitelist", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		InteractionFailed(session, interaction, "couldn't open whitelist file", err)
		return
	}
	defer file.Close()

	_, err = file.Write(userIDBytes)
	if err != nil {
		InteractionFailed(session, interaction, "couldn't write the user ID to the whitelist file", err)
		return
	}

	session.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: "Succesfully added " + user.Username + " to the whitelist.",
	})
}

func viewCommand(session *discordgo.Session, interaction *discordgo.Interaction, _ CommandOptions) {
	fileContentsBytes, err := os.ReadFile("./whitelist")
	if err != nil {
		InteractionFailed(session, interaction, "couldn't open whitelist file", err)
		return
	}

	fileContents := string(fileContentsBytes)

	var mentions string

	for _, userID := range strings.Split(strings.TrimSpace(fileContents), "\n") {
		mentions += fmt.Sprintf("- <@%s>\n", userID)
	}

	responseEmbed := &discordgo.MessageEmbed{
		Title:       "Users on the whitelist:",
		Description: mentions,
		Color:       0x00ADD8,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    session.State.User.Username,
			IconURL: session.State.User.AvatarURL(""),
		},
	}

	session.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{responseEmbed},
	})
}

func removeCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	userID := strings.TrimSpace(options["user_id"].StringValue())

	user, err := session.User(userID)
	if err != nil {
		InteractionFailed(session, interaction, "user doesn't exist or another error has occured", err)
		return
	}

	fileContentsBytes, err := os.ReadFile("./whitelist")
	if err != nil {
		InteractionFailed(session, interaction, "couldn't read whitelist file", err)
		return
	}

	fileContents := string(fileContentsBytes)

	var newContents string
	for _, whitelistedUserID := range strings.Split(strings.TrimSpace(fileContents), "\n") {
		if whitelistedUserID != userID {
			newContents += whitelistedUserID + "\n"
		}
	}

	file, err := os.OpenFile("./whitelist", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		InteractionFailed(session, interaction, "couldn't open whitelist file", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(newContents))
	if err != nil {
		InteractionFailed(session, interaction, "couldn't write to the whitelist file", err)
		return
	}

	session.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: "Successfully removed " + user.Username + " from the whitelist.",
	})
}
