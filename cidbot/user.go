package cidbot

import (
	"fmt"
	"time"
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

type PastUsername struct {
	Name string `json:"name"`
}

func GetUsersByUsername(name string) ([]*User, error) {
	requestUrl := "https://users.roblox.com/v1/usernames/users"

	requestData := GetUsersByUsernameRequest{
		Usernames:          []string{name},
		ExcludeBannedUsers: true,
	}

	response, err := PostRequest[ResponseData[*User]](requestUrl, requestData)
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

	response, err := PostRequest[ResponseData[*User]](requestUrl, requestData)
	if err != nil {
		return nil, err
	}

	return response.Data, err
}

func GetUserByID(id uint64) (*User, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", id)
	return GetRequest[User](requestUrl)
}

func GetUserPastUsernames(userID uint64) (list []string, err error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d/username-history?limit=100", userID)
	response, err := GetRequest[ResponseData[PastUsername]](requestUrl)
	if err != nil {
		return nil, err
	}

	for _, username := range response.Data {
		list = append(list, username.Name)
	}

	return list, nil
}
