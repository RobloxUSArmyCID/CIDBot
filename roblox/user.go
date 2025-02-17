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
	Name    string    `json:"name" json:"username"`
	ID      uint64    `json:"id" json:"userId"`
	Created time.Time `json:"created"`

	DaysFromCreation             int
	Groups                       []*GroupAndRole
	SuspiciousGroups             []*GroupAndRole
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

	username, userID, err := getUserNameAndIDByName(u.Name)
	if err != nil {
		return nil, err
	}

	u.ID = userID
	u.Name = username

	eg.Go(func() error {
		const (
			UsarGroupID = 3108077
			RankE1      = 5
		)

		groups, err := getUserGroups(u.ID)
		if err != nil {
			return err
		}

		isE1 := false
		isInUsar := false
		usarRank := "N/A"

		for _, group := range groups {
			if group.Group.ID == UsarGroupID {
				isInUsar = true
				usarRank = group.Role.Name
				isE1 = group.Role.Rank == RankE1
			}
		}

		if !isInUsar {
			usarRank = "N/A"
		}

		// O(n) with n-max = 100
		// MaxO(100) -- not worth of optimization
		susGroups := getSuspiciousGroups(groups)

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
		created, err := getUserCreationDate(u.ID)
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
		badges, err := getUserBadges(u.ID)
		if err != nil {
			return err
		}
		u.mu.Lock()
		defer u.mu.Unlock()
		u.Badges = badges
		return nil
	})

	eg.Go(func() error {
		friends, err := getUserFriends(u.ID)
		if err != nil {
			return err
		}
		susFriends := getSuspiciousFriends(u.Name, friends)

		usernamesOfSusFriends := []string{}
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
		thumbnailURL, err := getUserThumbnail(u.ID)
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

func getUserNameAndIDByName(username string) (string, uint64, error) {
	users, err := getUsersByNames([]string{username})
	return users[0].Name, users[0].ID, err
}

func getUsersByNames(names []string) ([]*User, error) {
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

func getUsersByIDs(ids []uint64) ([]*User, error) {
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

func getUserByID(id uint64) (*User, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", id)
	return requests.Get[User](requestUrl)
}

func getUserCreationDate(id uint64) (time.Time, error) {
	user, err := getUserByID(id)
	if err != nil {
		return time.Time{}, err
	}

	return user.Created, nil
}
