package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/RobloxUSArmyCID/CIDBot/roblox"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/time/rate"
)

type CommandOptions map[string]*discordgo.ApplicationCommandInteractionDataOption

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "bgcheck",
		Description: "Background check a Roblox user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "The username of the user you wish to background check",
				Required:    true,
			},
		},
		IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
			discordgo.ApplicationIntegrationUserInstall,
			discordgo.ApplicationIntegrationGuildInstall,
		},
	},
}

var whitelistPermissions int64 = discordgo.PermissionAdministrator

var AdminServerCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "whitelist",
		Description: "CID Bot whitelist commands",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Adds a user to the CID Bot whitelist",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "user_id",
						Description: "The Discord ID of the user who needs to be whitelisted",
						Required:    true,
					},
				},
			},

			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "Lists all users allowed to use the CID Bot",
			},

			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Remove a user from the CID Bot",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "user_id",
						Description: "The Discord ID of the user who needs to be removed from the whitelist",
						Required:    true,
					},
				},
			},
		},
		DefaultMemberPermissions: &whitelistPermissions,
	},
}

func ParseCommandOptions(opts []*discordgo.ApplicationCommandInteractionDataOption) CommandOptions {
	co := make(CommandOptions)
	for _, option := range opts {
		co[option.Name] = option
	}
	return co
}

var (
	groups    []*roblox.Group
	badges    []*roblox.Badge
	friends   []*roblox.User
	user      *roblox.User
	thumbnail *string

	mu sync.Mutex
)

var limiter = rate.NewLimiter(10, 1)

const (
	USAR_GROUP_ID           = 3108077
	USAR_E1_RANK            = 5
	THIRTY_REQUIRED_MEMBERS = 30
)

// error constants
var (
	errUnauthorized = errors.New("you are unauthorized to run this command")
)

func Executed(discord *discordgo.Session, interaction *discordgo.Interaction) {
	command := interaction.ApplicationCommandData()
	options := ParseCommandOptions(command.Options)

	slog.Info("command executed",
		"command", command.Name,
		"id", interaction.ID,
		"user", interaction.Member.User.ID,
		"guild", interaction.GuildID)

	err := deferInteraction(discord, interaction)
	if err != nil {
		slog.Warn("could not defer interaction", "err", err)
		return
	}

	switch command.Name {
	case "bgcheck":
		backgroundCheckCommand(discord, interaction, options)
	case "whitelist":
		whitelistCommand(discord, interaction)
	default:
		slog.Warn("incorrect command used", "command", command.Name, "id", interaction.ID)
	}
}

func interactionFailed(discord *discordgo.Session, interaction *discordgo.Interaction, content string, err error) error {
	slog.Warn("interaction failed", "id", interaction.ID, "guild", interaction.GuildID, "err", err)
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    interaction.Member.User.Username,
			IconURL: interaction.Member.User.AvatarURL(""),
		},
		Title:       ":x: | An error occurred!",
		Description: fmt.Sprintf("Error contents:\n```%s: %s```", content, err),
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x8b0000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "If you believe this is an error, contact the Investigatory Director.",
		},
	}

	_, err = discord.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})

	return err

}

func deferInteraction(discord *discordgo.Session, interaction *discordgo.Interaction) error {
	return discord.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}
