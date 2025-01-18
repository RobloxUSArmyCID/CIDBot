package roblox

import (
	"fmt"
	"time"

	"github.com/RobloxUSArmyCID/CIDBot/requests"
)

type User struct {
	Name string `json:"name"`
	ID   uint64 `json:"id"`

	// Is not given with GetUserInfoByUsername or GetUserInfoByID
	Created time.Time `json:"created"`
}

type GetUsersByUsernameRequest struct {
	Usernames          []string `json:"usernames"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

type GetUsersByIDRequest struct {
	UserIDs            []uint64 `json:"userIds"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

func GetUsersByUsernames(names []string) ([]*User, error) {
	requestUrl := "https://users.roblox.com/v1/usernames/users"

	requestData := GetUsersByUsernameRequest{
		Usernames:          names,
		ExcludeBannedUsers: true,
	}

	response, err := requests.PostRequest[requests.ResponseData[*User]](requestUrl, requestData)
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

	response, err := requests.PostRequest[requests.ResponseData[*User]](requestUrl, requestData)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}

func GetUserByID(id uint64) (*User, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", id)
	return requests.Get[User](requestUrl)
}
