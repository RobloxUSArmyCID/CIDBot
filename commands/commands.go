package commands

import (
	"errors"
	"log/slog"

	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/RobloxUSArmyCID/CIDBot/embeds"
	"github.com/bwmarrin/discordgo"
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
		Contexts: &[]discordgo.InteractionContextType{
			discordgo.InteractionContextBotDM,
			discordgo.InteractionContextGuild,
			discordgo.InteractionContextPrivateChannel,
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
		IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
			discordgo.ApplicationIntegrationGuildInstall,
		},
	},
}

func ParseOptions(opts []*discordgo.ApplicationCommandInteractionDataOption) CommandOptions {
	co := make(CommandOptions)
	for _, option := range opts {
		co[option.Name] = option
	}
	return co
}

// error constants
var (
	errUnauthorized = errors.New("you are unauthorized to run this command")
)

func Executed(discord *discordgo.Session, interaction *discordgo.Interaction, config *config.Config) {
	command := interaction.ApplicationCommandData()
	options := ParseOptions(command.Options)

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
		backgroundCheckCommand(discord, interaction, options, config)
	case "whitelist":
		whitelist(discord, interaction, config)
	default:
		slog.Warn("incorrect command used", "command", command.Name, "id", interaction.ID)
	}
}

func getCommandInvoker(interaction *discordgo.Interaction) *discordgo.User {
	if interaction.User != nil {
		return interaction.User
	}

	return interaction.Member.User
}

func interactionFailed(discord *discordgo.Session, interaction *discordgo.Interaction, err error) {
	slog.Warn("interaction failed", "id", interaction.ID, "guild", interaction.GuildID, "err", err)

	invoker := getCommandInvoker(interaction)

	embed := embeds.NewBuilder().
		SetAuthorUser(invoker).
		SetColor(embeds.ColorError).
		SetCurrentTimestamp().
		SetTitle(":x: | An error occured!").
		SetCodeBlockDescription(err.Error()).
		SetFooter("If you believe this is a bug, contact f_o1oo or the CID ID.", "").
		Build()

	discord.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})

}

func deferInteraction(discord *discordgo.Session, interaction *discordgo.Interaction) error {
	return discord.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}
