package roblox

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
	"golang.org/x/sync/errgroup"
)

type User struct {
	Name    string    `json:"name"`
	ID      uint64    `json:"id"`
	Created time.Time `json:"created"`

	DaysFromCreation             int
	Groups                       []*Group
	SuspiciousGroups             []*Group
	Badges                       []*Badge
	Friends                      []*User
	SuspiciousFriends            []*User
	UsernamesOfSuspiciousFriends []string
	ThumbnailURL                 string

	IsE1     bool
	IsInUsar bool
	UsarRank string

	mu sync.Mutex
}

type GetUsersByUsernameRequest struct {
	Usernames          []string `json:"usernames"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

type GetUsersByIDRequest struct {
	UserIDs            []uint64 `json:"userIds"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

func NewUser(username string) (*User, error) {

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	u := &User{
		Name: username,
	}

	userID, err := getUserIDByUsername(u.Name)
	if err != nil {
		return nil, err
	}

	u.ID = userID

	eg.Go(func() error {
		const (
			UsarGroupID = 3108077
			RankE1      = 5
		)

		groups, err := GetUserGroups(u.ID)
		if err != nil {
			return err
		}

		var isE1, isInUsar bool
		var usarRank string

		for _, group := range groups {
			if group.Group.ID == UsarGroupID {
				isInUsar = true
				usarRank = group.Role.Name
				isE1 = group.Role.Rank == RankE1
			}
		}

		susGroups := GetSuspiciousGroups(groups)

		u.mu.Lock()
		defer u.mu.Unlock()

		u.Groups = groups
		u.IsE1 = isE1
		u.IsInUsar = isInUsar
		u.UsarRank = usarRank
		u.SuspiciousGroups = susGroups
		return nil
	})

	eg.Go(func() error {
		created, err := GetUserCreationDate(u.ID)
		if err != nil {
			return err
		}
		u.mu.Lock()
		defer u.mu.Unlock()
		u.Created = created
		u.DaysFromCreation = int(time.Since(u.Created).Hours() / 24)
		return nil
	})

	eg.Go(func() error {
		badges, err := GetUserBadges(u.ID)
		if err != nil {
			return err
		}
		u.mu.Lock()
		defer u.mu.Unlock()
		u.Badges = badges
		return nil
	})

	eg.Go(func() error {
		friends, err := GetUserFriends(u.ID)
		if err != nil {
			return err
		}

		susFriends := GetSuspiciousFriends(u, u.Friends)

		var usernamesOfSusFriends []string
		for _, susFriend := range susFriends {
			usernamesOfSusFriends = append(usernamesOfSusFriends, susFriend.Name)
		}

		u.mu.Lock()
		defer u.mu.Unlock()

		u.Friends = friends
		u.SuspiciousFriends = susFriends
		u.UsernamesOfSuspiciousFriends = usernamesOfSusFriends
		return nil
	})

	eg.Go(func() error {
		thumbnailURL, err := GetUserThumbnail(u.ID)
		if err != nil {
			return err
		}
		u.mu.Lock()
		defer u.mu.Unlock()
		u.ThumbnailURL = *thumbnailURL
		return nil
	})

	err = eg.Wait()
	return u, err
}

func getUserIDByUsername(username string) (uint64, error) {
	users, err := GetUsersByUsernames([]string{username})
	return users[0].ID, err
}

func GetUsersByUsernames(names []string) ([]*User, error) {
	requestUrl := "https://users.roblox.com/v1/usernames/users"

	requestData := GetUsersByUsernameRequest{
		Usernames:          names,
		ExcludeBannedUsers: true,
	}

	response, err := requests.Post[requests.ResponseData[*User]](requestUrl, requestData)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}

func GetUsersByID(ids []uint64) ([]*User, error) {
	requestUrl := "https://users.roblox.com/v1/users"

	requestData := GetUsersByIDRequest{
		UserIDs:            ids,
		ExcludeBannedUsers: true,
	}

	response, err := requests.Post[requests.ResponseData[*User]](requestUrl, requestData)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}

func GetUserByID(id uint64) (*User, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", id)
	return requests.Get[User](requestUrl)
}

func GetUserCreationDate(id uint64) (time.Time, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return time.Time{}, err
	}

	return user.Created, nil
}
