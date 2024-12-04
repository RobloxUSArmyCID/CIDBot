package cidbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type User struct {
	Name string `json:"name"`
	ID   uint64 `json:"id"`

	// Is not given with GetUserInfoByUsername
	Created time.Time `json:"created"`
}

type GetUserRequest struct {
	Usernames          []string `json:"usernames"`
	ExcludeBannedUsers bool     `json:"excludeBannedUsers"`
}

func GetUserInfoByUsername(name string) (*User, error) {
	requestUrl := "https://users.roblox.com/v1/usernames/users"

	requestData := GetUserRequest{
		Usernames:          []string{name},
		ExcludeBannedUsers: true,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code - %d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseBody ResponseData[User]

	err = json.Unmarshal(responseBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	if len(responseBody.Data) == 0 {
		return nil, errors.New("roblox did not return a user")
	}

	return &responseBody.Data[0], nil
}

func GetUserInfoByID(id uint64) (*User, error) {
	requestUrl := fmt.Sprintf("https://users.roblox.com/v1/users/%d", id)

	response, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code - %d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseBody ResponseData[User]

	err = json.Unmarshal(responseBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	if len(responseBody.Data) == 0 {
		return nil, errors.New("roblox did not return a user")
	}

	return &responseBody.Data[0], nil
}
