package cidbot

import (
	"fmt"

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
	},
}

func ParseCommandOptions(opts []*discordgo.ApplicationCommandInteractionDataOption) CommandOptions {
	co := make(CommandOptions)
	for _, option := range opts {
		co[option.Name] = option
	}
	return co
}

func BackgroundCheckCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	//	g := errgroup.Group{}
	username := options["username"].StringValue()

	robloxUser, err := GetUserInfoByUsername(username)

	if err != nil {
		InteractionFailed(session, interaction, "could not get user info", err)
		return
	}

	session.FollowupMessageCreate(interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s, %d, %s", robloxUser.Name, robloxUser.ID, robloxUser.Created.String()),
	})

}
