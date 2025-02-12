package commands

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/RobloxUSArmyCID/CIDBot/config"
	"github.com/RobloxUSArmyCID/CIDBot/embeds"
	"github.com/RobloxUSArmyCID/CIDBot/roblox"
	"github.com/bwmarrin/discordgo"
)

func backgroundCheckCommand(discord *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions, config *config.Config) {
	slog.Debug("opening whitelist file")
	whitelistBytes, err := os.ReadFile(config.WhitelistPath)
	if err != nil {
		interactionFailed(discord, interaction, "couldn't open the whitelist file", errUnauthorized)
		return
	}

	var invoker *discordgo.User
	if interaction.Member.User == nil {
		invoker = interaction.User
	} else {
		invoker = interaction.Member.User
	}

	slog.Debug("checking if invoker is whitelisted", "id", invoker.ID)
	whitelist := string(whitelistBytes)
	if !strings.Contains(whitelist, invoker.ID) {
		interactionFailed(discord, interaction, "You are not allowed to run this command", errUnauthorized)
		return
	}

	username := options["username"].StringValue()

	user, err := roblox.NewUser(username)

	if err != nil {
		interactionFailed(discord, interaction, "", err)
		return
	}

	var descriptionBuilder strings.Builder
	failedBackgroundCheck := false

	if user.UsarRank == "N/A" {
		descriptionBuilder.WriteString("- ⚠ Not in USAR ⚠\n")
		failedBackgroundCheck = true
	}

	if user.IsE1 {
		descriptionBuilder.WriteString("- ⚠ E1 ⚠\n")
		failedBackgroundCheck = true
	}

	if len(user.Badges) < 100 {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Less than 100 badges (%d) ⚠\n", len(user.Badges)))
		failedBackgroundCheck = true
	}

	if user.DaysFromCreation < 365 {
		if user.DaysFromCreation < 90 {
			descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Account age under 90 days old (%d) (failing) ⚠\n", user.DaysFromCreation))
		} else {
			descriptionBuilder.WriteString("- ⚠ Account age under 365 days old (suspicious, not failing) ⚠\n")
		}
	}

	if len(user.Friends) <= 3 {
		descriptionBuilder.WriteString("- ⚠ 3 or less friends. ⚠\n")
		failedBackgroundCheck = true
	}

	if len(user.Groups) <= 15 {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ In 15 or less groups (%d) ⚠\n", len(user.Friends)))
	}

	for _, group := range user.SuspiciousGroups {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Suspicious group: %s (%d members) ⚠\n", group.Group.Name, group.Group.MemberCount))
	}

	if len(user.SuspiciousFriends) > 0 {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Suspected alts in friends list: %s\n", strings.Join(user.UsernamesOfSuspiciousFriends, ", ")))
	}

	if descriptionBuilder.Len() == 0 {
		descriptionBuilder.WriteString("+ No concerns found! (verify criteria not checked by the bot)\n")
	}

	descriptionBuilder.WriteString("Check past usernames!")

	description := descriptionBuilder.String()

	if len(description) > 4096 {
		description = "- ❌ Embed description too long! Please manually background check the user! ❌\n"
	}

	failed := "+ No"
	if failedBackgroundCheck {
		failed = "- Yes"
	}

	profileURL := fmt.Sprintf("https://www.roblox.com/users/%d/profile", user.ID)

	embed := embeds.NewBuilder().
		SetAuthorUser(invoker).
		SetColor(embeds.ColorGopherBlue).
		SetCurrentTimestamp().
		SetTitle(":white_check_mark: | Background check finished!").
		SetCodeBlockDescription(description).
		SetURL(profileURL).
		AddCodeBlockField("Username:", user.Name, true).
		AddCodeBlockField("ID:", fmt.Sprintf("%d", user.ID), true).
		AddCodeBlockField("Failed:", failed, true).
		AddCodeBlockField("USAR Rank:", user.UsarRank, true).
		AddCodeBlockField("Account age:", fmt.Sprintf("%d days old", user.DaysFromCreation), true).
		Build()

	discord.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	})
}
