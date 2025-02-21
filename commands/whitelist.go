package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/bwmarrin/discordgo"
)

func whitelist(discord *discordgo.Session, interaction *discordgo.Interaction, config *config.Config) {
	subcommand := interaction.ApplicationCommandData().Options[0]
	subcommandOptions := ParseOptions(subcommand.Options)
	switch subcommand.Name {
	case "add":
		whitelistAdd(discord, interaction, subcommandOptions, config)
	case "view":
		whitelistView(discord, interaction, config)
	case "remove":
		whitelistRemove(discord, interaction, subcommandOptions, config)
	default:
		log.Printf("incorrect whitelist command used: %s", subcommand.Name)
	}
}

func whitelistAdd(discord *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions, config *config.Config) {
	userID := options["user_id"].StringValue()
	userIDBytes := []byte(userID + "\n")

	user, err := discord.User(userID)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}

	file, err := os.OpenFile(config.WhitelistPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}
	defer file.Close()

	_, err = file.Write(userIDBytes)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}

	discord.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: "Succesfully added " + user.Username + " to the whitelist.",
	})
}

func whitelistView(discord *discordgo.Session, interaction *discordgo.Interaction, config *config.Config) {
	fileContentsBytes, err := os.ReadFile(config.WhitelistPath)
	if err != nil {
		interactionFailed(discord, interaction, err)
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
			Text:    discord.State.User.Username,
			IconURL: discord.State.User.AvatarURL(""),
		},
	}

	discord.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{responseEmbed},
	})
}

func whitelistRemove(discord *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions, config *config.Config) {
	userID := strings.TrimSpace(options["user_id"].StringValue())

	user, err := discord.User(userID)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}

	fileContentsBytes, err := os.ReadFile(config.WhitelistPath)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}

	fileContents := string(fileContentsBytes)

	var newContents string
	for _, whitelistedUserID := range strings.Split(strings.TrimSpace(fileContents), "\n") {
		if whitelistedUserID != userID {
			newContents += whitelistedUserID + "\n"
		}
	}

	file, err := os.OpenFile(config.WhitelistPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(newContents))
	if err != nil {
		interactionFailed(discord, interaction, err)
		return
	}

	discord.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: "Successfully removed " + user.Username + " from the whitelist.",
	})
}
