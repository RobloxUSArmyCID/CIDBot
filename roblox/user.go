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

	CanViewInventory             bool
	DaysFromCreation             int
	Groups                       []*GroupAndRole
	SuspiciousGroups             []*GroupAndRole
	UsarUnits                    []string
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
	// automatically adds a WithCancelCause trait to the context
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

	eg.Go(u.GetGroups)
	eg.Go(u.GetFriends)
	eg.Go(u.GetBadges)
	eg.Go(u.GetThumbnail)
	eg.Go(u.GetCreationTime)

	if err = eg.Wait(); err != nil {
		return nil, err
	}

	return u, nil
}

func getUserNameAndIDByName(username string) (string, uint64, error) {
	users, err := getUsersByNames([]string{username})
	if len(users) != 1 {
		return "", 0, fmt.Errorf("no such user exists: %s", username)
	}
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

func (u *User) GetCreationTime() error {
	created, err := getUserCreationDate(u.ID)
	if err != nil {
		return err
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Created = created
	u.DaysFromCreation = int(time.Since(u.Created).Hours() / 24)
	return nil

}
