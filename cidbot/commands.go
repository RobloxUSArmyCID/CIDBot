package cidbot

import (
	"sync"

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
	groups    []*Group
	badges    []*Badge
	friends   []*User
	user      *User
	thumbnail *string

	mu sync.Mutex
)

const ( 
	USAR_GROUP_ID = 3108077 
	USAR_E1_RANK = 5
	THIRTY_REQUIRED_MEMBERS = 30
)

func BackgroundCheckCommand(session *discordgo.Session, interaction *discordgo.Interaction, options CommandOptions) {
	username := options["username"].StringValue()

	user, err := GetUserInfoByUsername(username)

	if err != nil {
		InteractionFailed(session, interaction, "could not get user info by username", err)
		return
	}

	doUserInfoCalls(user.ID)

	usarRank := ""
	isE1 := false

	for _, group := range groups {
		var groupsUnder30Members []*Group

		if group.Group.MemberCount <= THIRTY_REQUIRED_MEMBERS {
			groupsUnder30Members = append(groupsUnder30Members, group)
		}

		if group.Group.ID == USAR_GROUP_ID {
			usarRank = group.Role.Name
			isE1 = group.Role.Rank == USAR_E1_RANK
		}
	}
}

func doUserInfoCalls(userID uint64) {
	concurrentCalls := errgroup.Group{}

	concurrentCalls.Go(func() error {
		data, err := GetUserInfoByID(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		user = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := GetUserGroups(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		groups = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := GetUserBadges(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		badges = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := GetUserFriends(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		friends = data
		mu.Unlock()
		return nil
	})

	concurrentCalls.Go(func() error {
		data, err := GetUserThumbnail(userID)
		if err != nil {
			return err
		}
		mu.Lock()
		thumbnail = data
		mu.Unlock()
		return nil
	})
	concurrentCalls.Wait()
}
