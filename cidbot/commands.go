package cidbot

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/RobloxUSArmyCID/CIDBot/roblox"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/errgroup"
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

var (
	groups        []*roblox.Group
	badges        []*roblox.Badge
	friends       []*roblox.User
	user          *roblox.User
	thumbnail     *string
	pastUsernames []string

	mu sync.Mutex
)

const (
	USAR_GROUP_ID           = 3108077
	USAR_E1_RANK            = 5
	THIRTY_REQUIRED_MEMBERS = 30
)

func BackgroundCheckCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	username := options["username"].StringValue()

	temp_user, err := roblox.GetUsersByUsernames([]string{username})
	if len(temp_user) == 0 {
		InteractionFailed(session, interaction, "no such user exists", err)
		return
	}

	if err != nil {
		InteractionFailed(session, interaction, "could not get user info by username", err)
		return
	}

	err = doUserInfoCalls(temp_user[0].ID)
	if err != nil {
		InteractionFailed(session, interaction, "error happened when doing one of the requests to roblox", err)
		return
	}

	friendsIDs := []uint64{}
	for _, friend := range friends {
		friendsIDs = append(friendsIDs, friend.ID)
	}

	friendsWithNames, err := roblox.GetUsersByID(friendsIDs)
	if err != nil {
		InteractionFailed(session, interaction, "error happened when getting one of the user's friends' information", err)
		return
	}

	usarRank := "N/A"
	isE1 := false
	var groupsUnder30Members []*roblox.Group

	for _, group := range groups {
		if group.Group.MemberCount <= THIRTY_REQUIRED_MEMBERS {
			groupsUnder30Members = append(groupsUnder30Members, group)
		}

		if group.Group.ID == USAR_GROUP_ID {
			usarRank = group.Role.Name
			isE1 = group.Role.Rank == USAR_E1_RANK
		}
	}

	daysFromAccountCreation := int(time.Since(user.Created).Hours() / 24)

	suspiciousFriends := roblox.GetSuspiciousFriends(user, friendsWithNames)

	amountOfFriends := len(friends)
	amountOfBadges := len(badges)
	amountOfGroups := len(groups)
	amountOfPastUsernames := len(pastUsernames)
	amountOfSuspiciousFriends := len(suspiciousFriends)

	var descriptionBuilder strings.Builder
	failedBackgroundCheck := false

	if usarRank == "N/A" {
		descriptionBuilder.WriteString("- ⚠ Not in USAR ⚠\n")
		failedBackgroundCheck = true
	}

	if isE1 {
		descriptionBuilder.WriteString("- ⚠ E1 ⚠\n")
		failedBackgroundCheck = true
	}

	if amountOfBadges < 100 {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Less than 100 badges (%d) ⚠\n", amountOfBadges))
		failedBackgroundCheck = true
	}

	if daysFromAccountCreation < 365 {
		if daysFromAccountCreation < 90 {
			descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Account age under 90 days old (%d) (failing) ⚠\n", daysFromAccountCreation))
		} else {
			descriptionBuilder.WriteString("- ⚠ Account age under 365 days old (suspicious, not failing) ⚠\n")
		}
	}

	if amountOfFriends <= 3 {
		descriptionBuilder.WriteString("- ⚠ 3 or less friends. ⚠\n")
		failedBackgroundCheck = true
	}

	if amountOfGroups <= 15 {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ In 15 or less groups (%d) ⚠\n", amountOfGroups))
	}

	for _, group := range groupsUnder30Members {
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Suspicious group: %s (%d members) ⚠\n", group.Group.Name, group.Group.MemberCount))
	}

	if amountOfSuspiciousFriends > 0 {
		usernamesOfSuspiciousFriends := []string{}
		for _, friend := range suspiciousFriends {
			usernamesOfSuspiciousFriends = append(usernamesOfSuspiciousFriends, friend.Name)
		}
		descriptionBuilder.WriteString(fmt.Sprintf("- ⚠ Suspected alts in friends list: %s\n", strings.Join(usernamesOfSuspiciousFriends, ", ")))
	}

	if amountOfPastUsernames > 0 {
		descriptionBuilder.WriteString(fmt.Sprintf("• Past usernames: %s\n", strings.Join(pastUsernames, ", ")))
	}

	if descriptionBuilder.Len() == 0 {
		descriptionBuilder.WriteString("+ No concerns found! (verify criteria not checked by the bot)\n")
	}

	description := descriptionBuilder.String()

	if len(description) > 4096 {
		description = "- ❌ Embed description too long! Please manually background check the user! ❌\n"
	}

	failed := "+ No"
	if failedBackgroundCheck {
		failed = "- Yes"
	}

	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    interaction.Member.User.Username,
			IconURL: interaction.Member.User.AvatarURL(""),
		},
		Title:       ":white_check_mark: | Background check finished!",
		URL:         fmt.Sprintf("https://www.roblox.com/users/%d/profile", user.ID),
		Description: fmt.Sprintf("```diff\n%s```", description),
		Color:       0x00ADD8, // gopher blue
		Timestamp:   time.Now().Format(time.RFC3339),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: *thumbnail,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username:",
				Value:  fmt.Sprintf("```%s```", user.Name),
				Inline: true,
			},
			{
				Name:   "ID:",
				Value:  fmt.Sprintf("```%d```", user.ID),
				Inline: true,
			},
			{
				Name:   "Failed:",
				Value:  fmt.Sprintf("```diff\n%s```", failed),
				Inline: true,
			},
			{
				Name:   "USAR Rank:",
				Value:  fmt.Sprintf("```%s```", usarRank),
				Inline: true,
			},
			{
				Name:   "Account Age:",
				Value:  fmt.Sprintf("```%d days old```", daysFromAccountCreation),
				Inline: true,
			},
		},
	}

	_, err = session.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			&embed,
		},
	})
	if err != nil {
		InteractionFailed(session, interaction, "could not send message", err)
	}

}

func doUserInfoCalls(userID uint64) error {
	concurrentCalls := errgroup.Group{}

	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserByID(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		user = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserGroups(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		groups = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserBadges(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		badges = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserFriends(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		friends = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserThumbnail(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		thumbnail = data
		mu.Unlock()
		return nil
	})
	concurrentCalls.Go(func() error {
		data, err := roblox.GetUserPastUsernames(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		pastUsernames = data
		mu.Unlock()
		return nil
	})

	return concurrentCalls.Wait()
}
